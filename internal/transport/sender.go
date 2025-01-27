package transport

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"

	"github.com/usercanal/sdk-go/internal/logger"
	pb "github.com/usercanal/sdk-go/proto"
	"github.com/usercanal/sdk-go/types"
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

// Sender handles event sending and metrics
type Sender struct {
	connMgr   *ConnManager
	client    pb.EventServiceClient
	apiKey    string
	startTime time.Time
	metrics   Metrics
	mu        sync.RWMutex

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewSender(apiKey, endpoint string) (*Sender, error) {
	if apiKey == "" {
		return nil, types.NewValidationError("apiKey", "cannot be empty")
	}

	if endpoint == "" {
		return nil, types.NewValidationError("endpoint", "cannot be empty")
	}

	logger.Debug("Creating new sender for endpoint: %s", endpoint)

	ctx, cancel := context.WithCancel(context.Background())

	// Create connection manager
	connMgr := NewConnManager(endpoint)

	s := &Sender{
		connMgr:   connMgr,
		apiKey:    apiKey,
		startTime: time.Now(),
		ctx:       ctx,
		cancel:    cancel,
	}

	// Attempt initial connection
	if err := connMgr.Connect(ctx); err != nil {
		cancel()
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

func (s *Sender) Send(ctx context.Context, events []*pb.Event) error {
	if len(events) == 0 {
		return nil
	}

	select {
	case <-s.ctx.Done():
		return types.NewValidationError("sender", "is shutting down")
	default:
	}

	// Add API key to context
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"x-api-key": s.apiKey,
	}))

	// Create batch request
	req := &pb.BatchRequest{
		Events: events,
	}

	// Get client
	s.mu.RLock()
	client := s.client
	s.mu.RUnlock()

	if client == nil {
		return types.NewValidationError("client", "not initialized")
	}

	// Send batch
	_, err := client.SendBatch(ctx, req)
	if err != nil {
		s.recordFailure()
		return &types.NetworkError{
			Operation: "Send",
			Message:   err.Error(),
		}
	}

	s.recordSuccess(len(events))
	return nil
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
	return s.connMgr.GetState().State.String()
}

func (s *Sender) Uptime() time.Duration {
	return time.Since(s.startTime)
}

func (s *Sender) Close() error {
	s.cancel()
	s.wg.Wait()

	if err := s.connMgr.Close(); err != nil {
		return &types.NetworkError{
			Operation: "Close",
			Message:   err.Error(),
		}
	}
	return nil
}
