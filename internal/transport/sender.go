package transport

import (
	"context"
	"fmt"
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
	defaultMaxMessageSize = 4 * 1024 * 1024 // 4MB
	defaultPingInterval   = 10 * time.Second
	defaultPingTimeout    = 3 * time.Second
	defaultWaitTimeout    = 5 * time.Second
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
	connMgr    *ConnManager
	client     pb.EventServiceClient
	maxRetries int
	apiKey     string
	startTime  time.Time
	metrics    Metrics
	mu         sync.RWMutex

	// Control
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

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

	ctx, cancel := context.WithCancel(context.Background())

	// Create connection manager with options
	connMgr := NewConnManager(endpoint)
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
			grpc.WaitForReady(true),
			grpc.MaxCallRecvMsgSize(defaultMaxMessageSize),
		),
	}
	connMgr.SetDialOptions(opts)

	s := &Sender{
		connMgr:    connMgr,
		maxRetries: maxRetries,
		apiKey:     apiKey,
		startTime:  time.Now(),
		ctx:        ctx,
		cancelFunc: cancel,
	}

	// Attempt initial connection
	if err := connMgr.Connect(ctx); err != nil {
		cancel() // Clean up context if connection fails
		return nil, &types.NetworkError{
			Operation: "Connect",
			Message:   err.Error(),
		}
	}

	// Initialize client
	s.updateClient()

	// Start state monitoring
	s.wg.Add(1)
	go s.monitorStateChanges()

	return s, nil
}

func (s *Sender) updateClient() {
	if conn := s.connMgr.GetConn(); conn != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.client = pb.NewEventServiceClient(conn)
	}
}

func (s *Sender) monitorStateChanges() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case state, ok := <-s.connMgr.StateChanges():
			if !ok {
				return
			}
			if state.State == connectivity.Ready {
				s.updateClient()
			}
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

func (s *Sender) Send(ctx context.Context, events []*pb.Event) error {
	if len(events) == 0 {
		return nil
	}

	// Check if sender is shutting down
	select {
	case <-s.ctx.Done():
		return types.NewValidationError("sender", "is shutting down")
	default:
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"x-api-key": s.apiKey,
	}))

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

		// Get current connection state
		state := s.connMgr.GetState()
		if state.State != connectivity.Ready {
			logger.Warn("Attempting to send while connection state is: %s (endpoint: %s)",
				state.State, state.Endpoint)
			if !s.waitForConnection(ctx) {
				continue
			}
		}

		s.mu.RLock()
		client := s.client
		s.mu.RUnlock()

		if client == nil {
			s.updateClient()
			continue
		}

		_, err := client.SendBatch(ctx, req)
		if err == nil {
			s.recordSuccess(len(events))
			return nil
		}

		lastErr = err
		s.recordFailure()
		logger.Warn("Send attempt %d/%d failed: %v", attempt+1, s.maxRetries+1, err)

		// Trigger reconnection through ConnManager
		s.connMgr.Connect(context.Background())
	}

	return &types.NetworkError{
		Operation: "Send",
		Message:   fmt.Sprintf("failed after %d attempts: %v", s.maxRetries, lastErr),
		Retries:   s.maxRetries,
	}
}

func (s *Sender) waitForConnection(ctx context.Context) bool {
	timer := time.NewTimer(defaultWaitTimeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-timer.C:
			return false
		case <-time.After(100 * time.Millisecond):
			if s.connMgr.GetState().State == connectivity.Ready {
				return true
			}
		}
	}
}

func (s *Sender) recordSuccess(eventCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metrics.EventsSent += int64(eventCount)
	s.metrics.BatchesSent++
	s.metrics.LastSendTime = time.Now()
	s.metrics.AverageBatchSize = float64(s.metrics.EventsSent) / float64(s.metrics.BatchesSent)
}

func (s *Sender) recordFailure() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metrics.FailedAttempts++
	s.metrics.LastFailureTime = time.Now()
}

func (s *Sender) GetMetrics() Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.metrics
}

func (s *Sender) State() string {
	state := s.connMgr.GetState()
	return state.State.String()
}

func (s *Sender) Uptime() time.Duration {
	return time.Since(s.startTime)
}

func (s *Sender) Close() error {
	// Signal shutdown
	s.cancelFunc()

	// Wait for goroutines to finish
	s.wg.Wait()

	// Close connection manager
	if err := s.connMgr.Close(); err != nil {
		return &types.NetworkError{
			Operation: "Close",
			Message:   err.Error(),
		}
	}
	return nil
}
