// batch/batch.go
package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/usercanal/sdk-go/internal/logger"
	pb "github.com/usercanal/sdk-go/proto"
	"github.com/usercanal/sdk-go/types"
)

const (
	defaultBatchSize     = 100
	defaultFlushInterval = 10 * time.Second
)

// SendFunc is the function type for sending events
type SendFunc func(context.Context, []*pb.Event) error

// Manager handles event batching and sending
type Manager struct {
	size         int
	interval     time.Duration
	send         SendFunc
	events       []*pb.Event
	lastFlush    time.Time
	lastFailure  time.Time
	mu           sync.RWMutex // Changed to RWMutex for better concurrency
	failedCount  int64
	successCount int64
	ticker       *time.Ticker  // Added ticker for periodic flushes
	done         chan struct{} // Added done channel for clean shutdown
}

// NewManager creates a new batch manager
func NewManager(size int, interval time.Duration, send SendFunc) *Manager {
	if send == nil {
		panic("send function cannot be nil")
	}

	if size <= 0 {
		logger.Warn("Invalid batch size %d, using default %d", size, defaultBatchSize)
		size = defaultBatchSize
	}

	if interval <= 0 {
		logger.Warn("Invalid flush interval %v, using default %v", interval, defaultFlushInterval)
		interval = defaultFlushInterval
	}

	m := &Manager{
		size:     size,
		interval: interval,
		send:     send,
		events:   make([]*pb.Event, 0, size),
		done:     make(chan struct{}),
		ticker:   time.NewTicker(interval),
	}

	// Start periodic flush
	go m.periodicFlush()

	return m
}

// periodicFlush runs periodic flush based on the interval
func (m *Manager) periodicFlush() {
	for {
		select {
		case <-m.done:
			return
		case <-m.ticker.C:
			if err := m.Flush(context.Background()); err != nil {
				logger.Warn("Periodic flush failed: %v", err)
			}
		}
	}
}

// Add adds an event to the batch
func (m *Manager) Add(ctx context.Context, event *pb.Event) error {
	if event == nil {
		return types.NewValidationError("event", "cannot be nil")
	}

	select {
	case <-ctx.Done():
		return &types.TimeoutError{
			Operation: "BatchAdd",
			Duration:  ctx.Err().Error(),
		}
	default:
		m.mu.Lock()
		m.events = append(m.events, event)
		needsFlush := len(m.events) >= m.size
		m.mu.Unlock()

		if needsFlush {
			return m.Flush(ctx)
		}

		return nil
	}
}

// Flush sends all pending events
func (m *Manager) Flush(ctx context.Context) error {
	m.mu.Lock()
	if len(m.events) == 0 {
		m.mu.Unlock()
		return nil
	}

	events := m.events
	m.events = make([]*pb.Event, 0, m.size)
	m.mu.Unlock()

	if err := m.send(ctx, events); err != nil {
		m.mu.Lock()
		m.failedCount += int64(len(events))
		m.lastFailure = time.Now()
		m.mu.Unlock()

		// Re-queue events on failure if context isn't cancelled
		select {
		case <-ctx.Done():
			return &types.TimeoutError{
				Operation: "Flush",
				Duration:  ctx.Err().Error(),
			}
		default:
			m.mu.Lock()
			m.events = append(m.events, events...)
			m.mu.Unlock()
			return &types.NetworkError{
				Operation: "Flush",
				Message:   err.Error(),
			}
		}
	}

	m.mu.Lock()
	m.successCount += int64(len(events))
	m.lastFlush = time.Now()
	m.mu.Unlock()

	logger.Debug("Flushed %d events successfully", len(events))
	return nil
}

// QueueSize returns the current number of events in queue
func (m *Manager) QueueSize() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return int64(len(m.events))
}

// FailedCount returns the total number of failed events
func (m *Manager) FailedCount() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.failedCount
}

// SuccessCount returns the total number of successfully sent events
func (m *Manager) SuccessCount() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.successCount
}

// LastFlushTime returns the time of the last successful flush
func (m *Manager) LastFlushTime() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastFlush
}

// LastFailureTime returns the time of the last failure
func (m *Manager) LastFailureTime() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastFailure
}

// Close stops the manager and flushes remaining events
func (m *Manager) Close() error {
	// Stop the ticker first to prevent new flushes
	m.ticker.Stop()

	// Signal monitor goroutine to stop
	close(m.done)

	// Try to flush remaining events with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get current queue size for logging
	queueSize := m.QueueSize()
	if queueSize > 0 {
		logger.Debug("Attempting to flush %d remaining events during shutdown", queueSize)
	}

	if err := m.Flush(ctx); err != nil {
		return &types.NetworkError{
			Operation: "Close",
			Message:   fmt.Sprintf("failed to flush %d events: %v", queueSize, err),
		}
	}

	// Double check if we have any events left (very unlikely but possible)
	remainingEvents := m.QueueSize()
	if remainingEvents > 0 {
		logger.Warn("%d events remained unflushed during shutdown", remainingEvents)
	}

	return nil
}
