// sdk-go/internal/api/client.go
package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/usercanal/sdk-go/internal/batch"
	configDefaults "github.com/usercanal/sdk-go/internal/config"
	"github.com/usercanal/sdk-go/internal/identity"
	"github.com/usercanal/sdk-go/internal/logger"
	"github.com/usercanal/sdk-go/internal/transport"
	"github.com/usercanal/sdk-go/types"
)

// Use centralized defaults from config package
const (
	defaultEndpoint      = configDefaults.DefaultEndpoint
	defaultBatchSize     = configDefaults.DefaultBatchSize
	defaultFlushInterval = configDefaults.DefaultFlushInterval
	defaultMaxRetries    = configDefaults.DefaultMaxRetries
	defaultCloseTimeout  = configDefaults.DefaultCloseTimeout
)

// Client represents an analytics client
type Client struct {
	cfg          *config
	sender       *transport.Sender
	eventBatcher *batch.Manager
	logBatcher   *batch.Manager
	identityMgr  *identity.Manager
	mu           sync.RWMutex
	closed       bool
	closing      bool
}

// Config represents the external client configuration
type Config struct {
	Endpoint      string        `json:"endpoint"`       // API Endpoint
	BatchSize     int           `json:"batch_size"`     // Events per batch
	FlushInterval time.Duration `json:"flush_interval"` // Max time between sends
	MaxRetries    int           `json:"max_retries"`    // Retry attempts
	Debug         bool          `json:"debug"`          // Enable debug logging
}

// internal config struct
type config struct {
	endpoint      string
	batchSize     int
	flushInterval time.Duration
	maxRetries    int
	debug         bool
}

func defaultConfig() *config {
	return &config{
		endpoint:      defaultEndpoint,
		batchSize:     defaultBatchSize,
		flushInterval: defaultFlushInterval,
		maxRetries:    defaultMaxRetries,
		debug:         configDefaults.DefaultDebug,
	}
}

// Option configures a Client
type Option func(*config)

// WithConfig converts a public Config to functional options
func WithConfig(cfg Config) Option {
	return func(c *config) {
		if cfg.Endpoint != "" {
			c.endpoint = cfg.Endpoint
		}
		if cfg.BatchSize > 0 {
			c.batchSize = cfg.BatchSize
		}
		if cfg.FlushInterval > 0 {
			c.flushInterval = cfg.FlushInterval
		}
		if cfg.MaxRetries > 0 {
			c.maxRetries = cfg.MaxRetries
		}
		c.debug = cfg.Debug
		logger.SetDebug(cfg.Debug)
	}
}

// Individual option functions
func WithEndpoint(endpoint string) Option {
	return func(c *config) {
		if endpoint != "" {
			c.endpoint = endpoint
		}
	}
}

func WithFlushInterval(interval time.Duration) Option {
	return func(c *config) {
		if interval > 0 {
			c.flushInterval = interval
		}
	}
}

func WithMaxRetries(retries int) Option {
	return func(c *config) {
		if retries > 0 {
			c.maxRetries = retries
		}
	}
}

func WithBatchSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.batchSize = size
		}
	}
}

func WithDebug(debug bool) Option {
	return func(c *config) {
		c.debug = debug
		logger.SetDebug(debug)
	}
}

// New creates a new client with the provided API key and options
func New(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, types.NewValidationError("apiKey", "is required")
	}

	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	sender, err := transport.NewSender(apiKey, cfg.endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create sender: %w", err)
	}

	// Create wrapper functions to match batch.SendFunc signature
	eventSendFunc := func(ctx context.Context, items []interface{}) error {
		events := make([]*transport.Event, len(items))
		for i, item := range items {
			if event, ok := item.(*transport.Event); ok {
				events[i] = event
			} else {
				return fmt.Errorf("invalid event type: %T", item)
			}
		}
		return sender.SendEvents(ctx, events)
	}

	logSendFunc := func(ctx context.Context, items []interface{}) error {
		logs := make([]*transport.Log, len(items))
		for i, item := range items {
			if log, ok := item.(*transport.Log); ok {
				logs[i] = log
			} else {
				return fmt.Errorf("invalid log type: %T", item)
			}
		}
		return sender.SendLogs(ctx, logs)
	}

	eventBatchMgr := batch.NewManager(cfg.batchSize, cfg.flushInterval, eventSendFunc)
	logBatchMgr := batch.NewManager(cfg.batchSize, cfg.flushInterval, logSendFunc)

	// Create identity manager for session and device ID management
	identityMgr, err := identity.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create identity manager: %w", err)
	}

	client := &Client{
		cfg:          cfg,
		sender:       sender,
		eventBatcher: eventBatchMgr,
		logBatcher:   logBatchMgr,
		identityMgr:  identityMgr,
	}

	return client, nil
}

// Flush forces a flush of both event and log batchers
func (c *Client) Flush(ctx context.Context) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	// Flush both event and log batchers
	if err := c.eventBatcher.Flush(ctx); err != nil {
		return fmt.Errorf("failed to flush events: %w", err)
	}

	if err := c.logBatcher.Flush(ctx); err != nil {
		return fmt.Errorf("failed to flush logs: %w", err)
	}

	return nil
}

// checkClosed verifies if the client is closed
func (c *Client) checkClosed() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return types.NewValidationError("client", "is closed")
	}
	if c.closing {
		return types.NewValidationError("client", "is shutting down")
	}
	return nil
}

// Close flushes pending data and closes the client
func (c *Client) Close(ctx context.Context) error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return types.NewValidationError("client", "already closed")
	}
	if c.closing {
		c.mu.Unlock()
		return types.NewValidationError("client", "already shutting down")
	}

	// Mark as closing
	c.closing = true
	c.mu.Unlock()

	// Use provided context or create timeout context
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), defaultCloseTimeout)
		defer cancel()
	}

	var flushErr error
	if err := c.Flush(ctx); err != nil {
		flushErr = fmt.Errorf("failed to flush data during shutdown: %w", err)
	}

	// Close batchers
	if err := c.eventBatcher.Close(); err != nil {
		if flushErr == nil {
			flushErr = fmt.Errorf("failed to close event batcher: %w", err)
		}
	}

	if err := c.logBatcher.Close(); err != nil {
		if flushErr == nil {
			flushErr = fmt.Errorf("failed to close log batcher: %w", err)
		}
	}

	// Close the sender
	if err := c.sender.Close(); err != nil {
		if flushErr == nil {
			flushErr = fmt.Errorf("failed to close sender: %w", err)
		}
	}

	// Now mark as fully closed
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()

	return flushErr
}

// GenerateSessionID creates a new session ID
func (c *Client) GenerateSessionID() []byte {
	if c.identityMgr != nil {
		return c.identityMgr.GenerateEventID()
	}
	// Fallback if identity manager is not available
	return make([]byte, 16)
}

// ResetSession creates a new session
func (c *Client) ResetSession() {
	if c.identityMgr != nil {
		c.identityMgr.Reset()
	}
}
