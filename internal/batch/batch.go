// batch/batch.go
package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/usercanal/sdk-go/internal/logger"
	"github.com/usercanal/sdk-go/internal/transport"
	"github.com/usercanal/sdk-go/types"
)

const (
	defaultBatchSize     = 100
	defaultFlushInterval = 10 * time.Second
)

// SendFunc is the function type for sending events
type SendFunc func(context.Context, []*transport.Event) error

// Manager handles event batching and sending
type Manager struct {
	size         int
	interval     time.Duration
	send         SendFunc
	events       []*transport.Event
	lastFlush    time.Time
	lastFailure  time.Time
	mu           sync.RWMutex
	failedCount  int64
	successCount int64
	ticker       *time.Ticker
	done         chan struct{}
}

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
		events:   make([]*transport.Event, 0, size),
		done:     make(chan struct{}),
		ticker:   time.NewTicker(interval),
	}

	// Start periodic flush
	go m.periodicFlush()

	return m
}

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

func (m *Manager) Add(ctx context.Context, event *transport.Event) error {
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

func (m *Manager) Flush(ctx context.Context) error {
	m.mu.Lock()
	if len(m.events) == 0 {
		m.mu.Unlock()
		return nil
	}

	events := m.events
	m.events = make([]*transport.Event, 0, m.size)
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

// The rest of the methods remain the same, just with updated comments
func (m *Manager) QueueSize() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return int64(len(m.events))
}

func (m *Manager) FailedCount() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.failedCount
}

func (m *Manager) SuccessCount() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.successCount
}

func (m *Manager) LastFlushTime() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastFlush
}

func (m *Manager) LastFailureTime() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastFailure
}

func (m *Manager) Close() error {
	m.ticker.Stop()
	close(m.done)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	remainingEvents := m.QueueSize()
	if remainingEvents > 0 {
		logger.Warn("%d events remained unflushed during shutdown", remainingEvents)
	}

	return nil
}
