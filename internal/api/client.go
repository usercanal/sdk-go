// api/client.go
package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/usercanal/sdk-go/internal/batch"
	"github.com/usercanal/sdk-go/internal/convert"
	"github.com/usercanal/sdk-go/internal/logger"
	"github.com/usercanal/sdk-go/internal/transport"
	"github.com/usercanal/sdk-go/types"
)

const (
	defaultEndpoint      = "collect.usercanal.com:9000"
	defaultBatchSize     = 100
	defaultFlushInterval = 10 * time.Second
	defaultMaxRetries    = 3
	defaultCloseTimeout  = 5 * time.Second
)

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
		debug:         false,
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

// Client represents an analytics client
type Client struct {
	cfg     *config
	sender  *transport.Sender
	batcher *batch.Manager
	mu      sync.RWMutex
	closed  bool
	closing bool
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

	batchMgr := batch.NewManager(cfg.batchSize, cfg.flushInterval, sender.Send)

	client := &Client{
		cfg:     cfg,
		sender:  sender,
		batcher: batchMgr,
	}

	return client, nil
}

// Track sends an analytics event
func (c *Client) Track(ctx context.Context, event types.Event) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	// Set timestamp if not set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	transportEvent, err := convert.EventToInternal(&event)
	if err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	if err := c.batcher.Add(ctx, transportEvent); err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}

	return nil
}

// Identify associates a user with their traits
func (c *Client) Identify(ctx context.Context, identity types.Identity) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := identity.Validate(); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	transportEvent, err := convert.IdentityToInternal(&identity)
	if err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	if err := c.batcher.Add(ctx, transportEvent); err != nil {
		return fmt.Errorf("failed to add identity event: %w", err)
	}

	return nil
}

// Group associates a user with a group
func (c *Client) Group(ctx context.Context, groupInfo types.GroupInfo) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := groupInfo.Validate(); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	transportEvent, err := convert.GroupToInternal(&groupInfo)
	if err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	if err := c.batcher.Add(ctx, transportEvent); err != nil {
		return fmt.Errorf("failed to add group event: %w", err)
	}

	return nil
}

// Revenue tracks a revenue event
func (c *Client) Revenue(ctx context.Context, rev types.Revenue) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := rev.Validate(); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	transportEvent, err := convert.RevenueToInternal(&rev)
	if err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	if err := c.batcher.Add(ctx, transportEvent); err != nil {
		return fmt.Errorf("failed to add revenue event: %w", err)
	}

	return nil
}

// Flush forces a flush of any pending events
func (c *Client) Flush(ctx context.Context) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := c.batcher.Flush(ctx); err != nil {
		return fmt.Errorf("failed to flush events: %w", err)
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

// Close flushes pending events and closes the client
func (c *Client) Close() error {
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

	// Try to flush with timeout
	ctx, cancel := context.WithTimeout(context.Background(), defaultCloseTimeout)
	defer cancel()

	var flushErr error
	if err := c.Flush(ctx); err != nil {
		flushErr = fmt.Errorf("failed to flush events during shutdown: %w", err)
		logger.Warn(flushErr.Error())
	}

	// Close the sender
	if err := c.sender.Close(); err != nil {
		return fmt.Errorf("failed to close sender: %w", err)
	}

	// Now mark as fully closed
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()

	// Return flush error if it occurred
	return flushErr
}
