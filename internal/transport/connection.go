// transport/connection.go
package transport

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/usercanal/sdk-go/internal/logger"
	"github.com/usercanal/sdk-go/types"
)

const defaultTCPPort = "9000"

var ErrConnectionClosed = types.NewValidationError("connection", "is closed")

// ConnectionState represents the TCP connection state
type ConnectionState struct {
	State       string
	LastChanged time.Time
	Endpoint    string
}

type ConnManager struct {
	// Core connection
	conn     net.Conn
	endpoint string

	// State management
	currentState ConnectionState
	stateChange  chan ConnectionState

	// DNS management
	resolvedIPs    []string
	currentIPIndex int
	mu             sync.RWMutex

	// Retry handling
	backoff     backoff.BackOff
	attempts    int64 // using atomic for thread safety
	retrying    int32 // atomic flag for retry status
	retrySignal chan struct{}

	// Lifecycle management
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewConnManager(endpoint string) *ConnManager {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize exponential backoff
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 1 * time.Second
	b.Multiplier = 1.5
	b.RandomizationFactor = 0.2
	b.MaxInterval = 30 * time.Second
	b.MaxElapsedTime = 0 // Never stop retrying

	cm := &ConnManager{
		endpoint:    endpoint,
		ctx:         ctx,
		cancel:      cancel,
		stateChange: make(chan ConnectionState, 1),
		backoff:     b,
		retrySignal: make(chan struct{}, 1),
	}

	// Initialize state
	cm.currentState = ConnectionState{
		State:       "Idle",
		LastChanged: time.Now(),
		Endpoint:    endpoint,
	}

	// Initial DNS resolution
	if err := cm.resolveEndpoint(); err != nil {
		logger.Warn("Initial DNS resolution failed: %v", err)
	}

	// Start retry handler
	cm.wg.Add(1)
	go cm.handleRetries()

	return cm
}

func (cm *ConnManager) resolveEndpoint() error {
	host := cm.endpoint
	port := defaultTCPPort
	if h, p, err := net.SplitHostPort(cm.endpoint); err == nil {
		host = h
		port = p
	}

	ips, err := net.LookupHost(host)
	if err != nil {
		return fmt.Errorf("DNS resolution failed: %w", err)
	}

	cm.mu.Lock()
	cm.resolvedIPs = make([]string, len(ips))
	for i, ip := range ips {
		cm.resolvedIPs[i] = net.JoinHostPort(ip, port)
	}
	cm.mu.Unlock()

	logger.Debug("Resolved %s to %d endpoints", host, len(ips))
	return nil
}

func (cm *ConnManager) getNextEndpoint() string {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if len(cm.resolvedIPs) == 0 {
		return cm.endpoint
	}

	endpoint := cm.resolvedIPs[cm.currentIPIndex]
	cm.currentIPIndex = (cm.currentIPIndex + 1) % len(cm.resolvedIPs)
	return endpoint
}

func (cm *ConnManager) Connect(ctx context.Context) error {
	if cm.isClosed() {
		return ErrConnectionClosed
	}

	attempt := atomic.AddInt64(&cm.attempts, 1)
	endpoint := cm.getNextEndpoint()

	logger.Debug("Starting connection attempt %d to %s", attempt, endpoint)
	cm.updateState("Connecting")

	// Create TCP connection with timeout
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", endpoint)
	if err != nil {
		cm.updateState("Failed")
		logger.Error("Connection attempt %d failed: %v", attempt, err)
		cm.signalRetry()
		return fmt.Errorf("failed to connect to %s: %w", endpoint, err)
	}

	// Configure TCP connection
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetNoDelay(true)
	tcpConn.SetWriteBuffer(256 * 1024)

	cm.mu.Lock()
	if cm.conn != nil {
		cm.conn.Close()
	}
	cm.conn = conn
	cm.mu.Unlock()

	cm.updateState("Connected")
	logger.Debug("Connection established on attempt %d", attempt)
	return nil
}

func (cm *ConnManager) handleRetries() {
	defer cm.wg.Done()

	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-cm.retrySignal:
			if !atomic.CompareAndSwapInt32(&cm.retrying, 0, 1) {
				continue // Already retrying
			}

			operation := func() error {
				return cm.Connect(cm.ctx)
			}

			retryCount := 0
			notify := func(err error, d time.Duration) {
				retryCount++
				logger.Info("Connection retry %d scheduled in %v (error: %v)",
					retryCount, d, err)
			}

			if err := backoff.RetryNotify(operation, cm.backoff, notify); err != nil {
				logger.Error("Retry sequence failed after %d attempts: %v",
					retryCount, err)
			}

			atomic.StoreInt32(&cm.retrying, 0)
		}
	}
}

func (cm *ConnManager) signalRetry() {
	select {
	case cm.retrySignal <- struct{}{}:
	default:
	}
}

func (cm *ConnManager) GetConn() net.Conn {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.conn
}

func (cm *ConnManager) isClosed() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.ctx.Err() != nil
}

func (cm *ConnManager) Close() error {
	cm.mu.Lock()
	if cm.ctx.Err() != nil {
		cm.mu.Unlock()
		return nil
	}
	conn := cm.conn
	cm.conn = nil
	cm.mu.Unlock()

	// Cancel context first to stop all goroutines
	cm.cancel()

	// Wait for retry handler to finish
	cm.wg.Wait()

	// Close connection if it exists
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func (cm *ConnManager) GetState() ConnectionState {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.currentState
}

func (cm *ConnManager) StateChanges() <-chan ConnectionState {
	return cm.stateChange
}

func (cm *ConnManager) updateState(state string) {
	cm.mu.Lock()
	oldState := cm.currentState.State
	cm.currentState = ConnectionState{
		State:       state,
		LastChanged: time.Now(),
		Endpoint:    cm.endpoint,
	}
	cm.mu.Unlock()

	if oldState != state {
		logger.Debug("Connection state changed from %s to %s", oldState, state)
		select {
		case cm.stateChange <- cm.currentState:
		case <-cm.ctx.Done():
			return
		default:
			if cm.ctx.Err() == nil {
				logger.Warn("State change notification dropped - channel full")
			}
		}
	}
}

func (cm *ConnManager) GetAttempts() int64 {
	return atomic.LoadInt64(&cm.attempts)
}

func (cm *ConnManager) IsRetrying() bool {
	return atomic.LoadInt32(&cm.retrying) == 1
}

func (cm *ConnManager) ResetBackoff() {
	if b, ok := cm.backoff.(*backoff.ExponentialBackOff); ok {
		b.Reset()
	}
}
