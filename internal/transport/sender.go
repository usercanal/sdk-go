// transport/sender.go
package transport

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/usercanal/sdk-go/internal/logger"
	pb "github.com/usercanal/sdk-go/proto"
	"github.com/usercanal/sdk-go/types"
)

const (
	defaultMaxRetries     = 3
	defaultDialTimeout    = 5 * time.Second
	defaultMaxMessageSize = 4 * 1024 * 1024 // 4MB
	defaultPingInterval   = 10 * time.Second
	defaultPingTimeout    = 3 * time.Second
)

// Metrics tracks sending statistics
type Metrics struct {
	EventsSent       int64
	BatchesSent      int64
	FailedAttempts   int64
	LastSendTime     time.Time
	LastFailureTime  time.Time
	AverageBatchSize float64
}

// Sender handles the gRPC connection and event sending
type Sender struct {
	client     pb.EventServiceClient
	conn       *grpc.ClientConn
	maxRetries int
	apiKey     string
	state      connectivity.State
	startTime  time.Time
	metrics    Metrics
	mu         sync.RWMutex // Single mutex for both state and metrics
	done       chan struct{}
}

// NewSender creates a new gRPC sender
func NewSender(apiKey, endpoint string, maxRetries int) (*Sender, error) {
	if apiKey == "" {
		return nil, types.NewValidationError("apiKey", "cannot be empty")
	}

	if endpoint == "" {
		return nil, types.NewValidationError("endpoint", "cannot be empty")
	}

	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}

	logger.Debug("Creating new sender for endpoint: %s", endpoint)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                defaultPingInterval,
			Timeout:             defaultPingTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(
			grpc.CallContentSubtype("proto"),
			grpc.ForceCodec(&protoCodec{}),
			grpc.WaitForReady(true), // Changed to true to wait for connection
			grpc.MaxCallRecvMsgSize(defaultMaxMessageSize),
		),
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	if err != nil {
		return nil, &types.NetworkError{
			Operation: "Connect",
			Message:   err.Error(),
		}
	}

	s := &Sender{
		client:     pb.NewEventServiceClient(conn),
		conn:       conn,
		maxRetries: maxRetries,
		apiKey:     apiKey,
		startTime:  time.Now(),
		done:       make(chan struct{}),
	}

	// Initialize starting state
	s.state = conn.GetState()

	// Monitor connection state in background
	go s.monitorConnection()

	return s, nil
}

func (s *Sender) monitorConnection() {
	for {
		select {
		case <-s.done:
			return
		default:
			state := s.conn.GetState()
			if state != s.state {
				s.mu.Lock()
				oldState := s.state
				s.state = state
				s.mu.Unlock()
				logger.Debug("Connection state changed from %s to: %s", oldState, state)
			}
			s.conn.WaitForStateChange(context.Background(), state)
		}
	}
}

type protoCodec struct{}

func (c *protoCodec) Marshal(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, types.NewValidationError("message", "not a proto.Message")
	}
	return proto.Marshal(msg)
}

func (c *protoCodec) Unmarshal(data []byte, v interface{}) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return types.NewValidationError("message", "not a proto.Message")
	}
	return proto.Unmarshal(data, msg)
}

func (c *protoCodec) Name() string {
	return "proto"
}

func calculateBackoff(attempt int) time.Duration {
	if attempt == 0 {
		return 0
	}
	base := 100 * time.Millisecond
	max := time.Duration(1<<uint(attempt)) * base
	jitter := time.Duration(rand.Int63n(int64(max / 2)))
	return max + jitter
}

func (s *Sender) Send(ctx context.Context, events []*pb.Event) error {
	if len(events) == 0 {
		return nil
	}

	// Add API key to context metadata
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"x-api-key": s.apiKey,
	}))

	s.mu.RLock()
	state := s.state
	s.mu.RUnlock()

	// Only log state if not ready, but don't prevent sending
	if state != connectivity.Ready {
		logger.Warn("Attempting to send while connection state is: %s", state)
	}

	req := &pb.BatchRequest{
		Events: events,
	}

	var lastErr error
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return &types.TimeoutError{
				Operation: "Send",
				Duration:  ctx.Err().Error(),
			}
		default:
		}

		if attempt > 0 {
			backoff := calculateBackoff(attempt)
			logger.Debug("Retrying in %v... (attempt %d/%d)",
				backoff, attempt+1, s.maxRetries+1)
			time.Sleep(backoff)
		}

		_, err := s.client.SendBatch(ctx, req)
		if err == nil {
			s.mu.Lock()
			s.metrics.EventsSent += int64(len(events))
			s.metrics.BatchesSent++
			s.metrics.LastSendTime = time.Now()
			s.metrics.AverageBatchSize = float64(s.metrics.EventsSent) / float64(s.metrics.BatchesSent)
			s.mu.Unlock()

			logger.Debug("Successfully sent batch of %d events", len(events))
			return nil
		}

		lastErr = err
		s.mu.Lock()
		s.metrics.FailedAttempts++
		s.metrics.LastFailureTime = time.Now()
		s.mu.Unlock()

		logger.Warn("Failed to send batch (attempt %d/%d): %v",
			attempt+1, s.maxRetries+1, err)
	}

	return &types.NetworkError{
		Operation: "Send",
		Message:   lastErr.Error(),
		Retries:   s.maxRetries,
	}
}

func (s *Sender) GetMetrics() Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.metrics
}

func (s *Sender) State() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state.String()
}

func (s *Sender) Uptime() time.Duration {
	return time.Since(s.startTime)
}

func (s *Sender) Close() error {
	close(s.done) // Signal monitor goroutine to stop
	if err := s.conn.Close(); err != nil {
		return &types.NetworkError{
			Operation: "Close",
			Message:   err.Error(),
		}
	}
	return nil
}
