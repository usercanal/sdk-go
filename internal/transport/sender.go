// sdk-go/internal/transport/sender.go
package transport

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/usercanal/sdk-go/internal/logger"
	schema_common "github.com/usercanal/sdk-go/internal/schema/common"
	"github.com/usercanal/sdk-go/types"
)

// Sender handles data sending and metrics
type Sender struct {
	connMgr   *ConnManager
	apiKey    []byte
	startTime time.Time
	metrics   types.TransportMetrics
	mu        sync.RWMutex

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func generateBatchID() uint64 {
	var id uint64
	binary.Read(rand.Reader, binary.BigEndian, &id)
	return id
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

func (s *Sender) sendBatch(ctx context.Context, schemaType schema_common.SchemaType, data []byte) error {
	// Size validation for critical environments
	if len(data) > MaxBatchSize {
		return types.NewValidationError("batch", fmt.Sprintf("batch size %d exceeds limit %d", len(data), MaxBatchSize))
	}

	builder := flatbuffers.NewBuilder(1024)

	batchID := generateBatchID()
	apiKeyOffset := builder.CreateByteVector(s.apiKey)
	dataOffset := builder.CreateByteVector(data)

	schema_common.BatchStart(builder)
	schema_common.BatchAddApiKey(builder, apiKeyOffset)
	schema_common.BatchAddBatchId(builder, batchID)
	schema_common.BatchAddSchemaType(builder, schemaType)
	schema_common.BatchAddData(builder, dataOffset)
	batchOffset := schema_common.BatchEnd(builder)

	builder.Finish(batchOffset)
	finalData := builder.FinishedBytes()

	return s.sendFrame(ctx, finalData)
}

func (s *Sender) sendFrame(ctx context.Context, data []byte) error {
	// Send length-prefixed message
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))

	frame := make([]byte, len(lenBuf)+len(data))
	copy(frame, lenBuf)
	copy(frame[len(lenBuf):], data)

	// Get connection and send with graceful retry
	conn := s.connMgr.GetConn()
	if conn == nil {
		// Try to reconnect once for immediate recovery
		logger.Debug("No connection available, attempting immediate reconnect")
		if err := s.connMgr.Connect(ctx); err != nil {
			s.recordFailure()
			return &types.NetworkError{
				Operation: "Send",
				Message:   "no active connection and reconnect failed: " + err.Error(),
			}
		}
		conn = s.connMgr.GetConn()
		if conn == nil {
			s.recordFailure()
			return &types.NetworkError{
				Operation: "Send",
				Message:   "connection still unavailable after reconnect",
			}
		}
	}

	if deadline, ok := ctx.Deadline(); ok {
		conn.SetWriteDeadline(deadline)
	}

	_, err := conn.Write(frame)
	if err != nil {
		s.recordFailure()
		// Signal retry for connection issues
		s.connMgr.signalRetry()
		return &types.NetworkError{
			Operation: "Send",
			Message:   err.Error(),
		}
	}

	// Record bytes sent for metrics
	s.recordBytesSent(len(frame))
	return nil
}

func (s *Sender) recordEventSuccess(eventCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metrics.EventsSent += int64(eventCount)
	s.metrics.EventBatchesSent++
	s.metrics.TotalBatchesSent++
	s.metrics.LastSendTime = time.Now()
	s.metrics.ConnectionUptime = s.Uptime()
	s.metrics.ReconnectCount = s.connMgr.GetReconnectCount()

	// Calculate separate averages
	if s.metrics.EventBatchesSent > 0 {
		s.metrics.AverageEventBatchSize = float64(s.metrics.EventsSent) / float64(s.metrics.EventBatchesSent)
	}
}

func (s *Sender) recordLogSuccess(logCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metrics.LogsSent += int64(logCount)
	s.metrics.LogBatchesSent++
	s.metrics.TotalBatchesSent++
	s.metrics.LastSendTime = time.Now()
	s.metrics.ConnectionUptime = s.Uptime()
	s.metrics.ReconnectCount = s.connMgr.GetReconnectCount()

	// Calculate separate averages
	if s.metrics.LogBatchesSent > 0 {
		s.metrics.AverageLogBatchSize = float64(s.metrics.LogsSent) / float64(s.metrics.LogBatchesSent)
	}
}

func (s *Sender) recordBytesSent(bytes int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics.BytesSent += int64(bytes)
}

func (s *Sender) recordFailure() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metrics.FailedAttempts++
	s.metrics.LastFailureTime = time.Now()
}

func (s *Sender) GetMetrics() types.TransportMetrics {
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

// HealthCheck performs connection health check
func (s *Sender) HealthCheck() error {
	return s.connMgr.HealthCheck()
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
