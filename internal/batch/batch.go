// sdk-go/internal/batch/batch.go
package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/usercanal/sdk-go/internal/logger"
	"github.com/usercanal/sdk-go/types"
)

const (
	defaultBatchSize     = 100
	defaultFlushInterval = 10 * time.Second
)

// SendFunc is the function type for sending items (generic)
type SendFunc func(context.Context, []interface{}) error

// Manager handles batching and sending of any type of items
type Manager struct {
	size         int
	interval     time.Duration
	send         SendFunc
	items        []interface{} // Changed from []*transport.Event to []interface{}
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
		items:    make([]interface{}, 0, size), // Changed to interface{}
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

// Add accepts any type of item (interface{})
func (m *Manager) Add(ctx context.Context, item interface{}) error {
	if item == nil {
		return types.NewValidationError("item", "cannot be nil")
	}

	select {
	case <-ctx.Done():
		return &types.TimeoutError{
			Operation: "BatchAdd",
			Duration:  ctx.Err().Error(),
		}
	default:
		m.mu.Lock()
		m.items = append(m.items, item)
		needsFlush := len(m.items) >= m.size
		m.mu.Unlock()

		if needsFlush {
			return m.Flush(ctx)
		}

		return nil
	}
}

func (m *Manager) Flush(ctx context.Context) error {
	m.mu.Lock()
	if len(m.items) == 0 {
		m.mu.Unlock()
		return nil
	}

	items := m.items
	m.items = make([]interface{}, 0, m.size) // Changed to interface{}
	m.mu.Unlock()

	if err := m.send(ctx, items); err != nil {
		m.mu.Lock()
		m.failedCount += int64(len(items))
		m.lastFailure = time.Now()
		m.mu.Unlock()

		// Re-queue items on failure if context isn't cancelled
		select {
		case <-ctx.Done():
			return &types.TimeoutError{
				Operation: "Flush",
				Duration:  ctx.Err().Error(),
			}
		default:
			m.mu.Lock()
			m.items = append(m.items, items...)
			m.mu.Unlock()
			return &types.NetworkError{
				Operation: "Flush",
				Message:   err.Error(),
			}
		}
	}

	m.mu.Lock()
	m.successCount += int64(len(items))
	m.lastFlush = time.Now()
	m.mu.Unlock()

	logger.Debug("Flushed %d items successfully", len(items))
	return nil
}

func (m *Manager) QueueSize() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return int64(len(m.items)) // Changed from events to items
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
		logger.Debug("Attempting to flush %d remaining items during shutdown", queueSize)
	}

	if err := m.Flush(ctx); err != nil {
		return &types.NetworkError{
			Operation: "Close",
			Message:   fmt.Sprintf("failed to flush %d items: %v", queueSize, err),
		}
	}

	remainingItems := m.QueueSize()
	if remainingItems > 0 {
		logger.Warn("%d items remained unflushed during shutdown", remainingItems)
	}

	return nil
}
