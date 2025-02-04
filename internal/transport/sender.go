// transport/sender.go
package transport

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"sync"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	event_collector "github.com/usercanal/sdk-go/internal/event"
	"github.com/usercanal/sdk-go/internal/logger"
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
	apiKey    []byte
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

	// Convert hex API key to bytes
	apiKeyBytes, err := hex.DecodeString(apiKey)
	if err != nil {
		return nil, types.NewValidationError("apiKey", "invalid format")
	}

	logger.Debug("Creating new sender for endpoint: %s", endpoint)

	ctx, cancel := context.WithCancel(context.Background())

	// Create connection manager
	connMgr := NewConnManager(endpoint)

	s := &Sender{
		connMgr:   connMgr,
		apiKey:    apiKeyBytes,
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

	// Start state monitoring
	s.wg.Add(1)
	go s.monitorStateChanges()

	return s, nil
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
			logger.Debug("Connection state changed: %s", state.State)
		}
	}
}

// Event represents the internal event format for sending
type Event struct {
	Timestamp uint64
	EventType event_collector.EventType
	UserID    []byte
	Payload   []byte
}

func (s *Sender) Send(ctx context.Context, events []*Event) error {
	if len(events) == 0 {
		return nil
	}

	// Validate events
	for _, evt := range events {
		if evt.Timestamp == 0 {
			return types.NewValidationError("Timestamp", "is required")
		}
		if len(evt.UserID) == 0 {
			return types.NewValidationError("UserID", "is required")
		}
		if len(evt.Payload) == 0 {
			return types.NewValidationError("Payload", "is required")
		}
	}

	select {
	case <-s.ctx.Done():
		return types.NewValidationError("sender", "is shutting down")
	default:
	}

	builder := flatbuffers.NewBuilder(1024 * len(events))

	// Create events vector
	eventOffsets := make([]flatbuffers.UOffsetT, len(events))
	for i := len(events) - 1; i >= 0; i-- {
		evt := events[i]

		payloadOffset := builder.CreateByteVector(evt.Payload)
		userIDOffset := builder.CreateByteVector(evt.UserID)

		event_collector.EventStart(builder)
		event_collector.EventAddTimestamp(builder, evt.Timestamp)
		event_collector.EventAddEventType(builder, evt.EventType)
		event_collector.EventAddUserId(builder, userIDOffset)
		event_collector.EventAddPayload(builder, payloadOffset)
		eventOffsets[i] = event_collector.EventEnd(builder)
	}

	eventsVec := builder.CreateVectorOfTables(eventOffsets)

	// Create API key vector
	apiKeyOffset := builder.CreateByteVector(s.apiKey)

	// Create batch
	event_collector.EventBatchStart(builder)
	event_collector.EventBatchAddApiKey(builder, apiKeyOffset)
	event_collector.EventBatchAddEvents(builder, eventsVec)
	batchEnd := event_collector.EventBatchEnd(builder)

	builder.Finish(batchEnd)
	data := builder.FinishedBytes()

	// Send length-prefixed message
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))

	frame := make([]byte, len(lenBuf)+len(data))
	copy(frame, lenBuf)
	copy(frame[len(lenBuf):], data)

	// Get connection and send
	conn := s.connMgr.GetConn()
	if conn == nil {
		s.recordFailure()
		return &types.NetworkError{
			Operation: "Send",
			Message:   "no active connection",
		}
	}

	if deadline, ok := ctx.Deadline(); ok {
		conn.SetWriteDeadline(deadline)
	}

	_, err := conn.Write(frame)
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
	return s.connMgr.GetState().State
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
