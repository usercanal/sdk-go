package transport

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/usercanal/sdk-go/internal/logger"
	"github.com/usercanal/sdk-go/types"
)

const defaultGRPCPort = "50051"

var ErrConnectionClosed = types.NewValidationError("connection", "is closed")

// ConnectionState wraps gRPC connectivity state with additional metadata
type ConnectionState struct {
	State       connectivity.State
	LastChanged time.Time
	Endpoint    string
}

type ConnManager struct {
	// Core connection
	conn     *grpc.ClientConn
	endpoint string
	opts     []grpc.DialOption

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

	// Set default gRPC options
	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	cm := &ConnManager{
		endpoint:    endpoint,
		ctx:         ctx,
		cancel:      cancel,
		stateChange: make(chan ConnectionState, 1),
		backoff:     b,
		retrySignal: make(chan struct{}, 1),
		opts:        defaultOpts, // Set default options
	}

	// Initialize state
	cm.currentState = ConnectionState{
		State:       connectivity.Idle,
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
	port := defaultGRPCPort
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
	cm.updateState(connectivity.Connecting)

	// Define retry policy in service config
	retryPolicy := `{
        "methodConfig": [{
            "name": [{"service": "EventService"}],
            "retryPolicy": {
                "MaxAttempts": 4,
                "InitialBackoff": "1s",
                "MaxBackoff": "30s",
                "BackoffMultiplier": 1.5,
                "RetryableStatusCodes": [ "UNAVAILABLE" ]
            }
        }]}`

	// Add retry policy to dial options
	opts := append([]grpc.DialOption{
		grpc.WithDefaultServiceConfig(retryPolicy),
	}, cm.opts...)

	// Create connection with timeout
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, endpoint, opts...)
	if err != nil {
		cm.updateState(connectivity.TransientFailure)
		logger.Error("Connection attempt %d failed: %v", attempt, err)
		cm.signalRetry()
		return err
	}

	state := conn.GetState()
	logger.Debug("Initial connection state for attempt %d: %s", attempt, state)

	if state == connectivity.TransientFailure {
		conn.Close()
		cm.updateState(connectivity.TransientFailure)
		logger.Error("Connection attempt %d failed immediately", attempt)
		cm.signalRetry()
		return fmt.Errorf("connection failed immediately")
	}

	cm.mu.Lock()
	if cm.conn != nil {
		cm.conn.Close()
	}
	cm.conn = conn
	cm.updateState(state)
	cm.mu.Unlock()

	// Monitor connection state
	go cm.monitorConnection(conn)

	logger.Debug("Connection established on attempt %d with state %s", attempt, state)
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

			// Use backoff retry with better logging
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

func (cm *ConnManager) monitorConnection(conn *grpc.ClientConn) {
	for {
		state := conn.GetState()
		if state == connectivity.TransientFailure {
			cm.signalRetry()
			return
		}
		if !conn.WaitForStateChange(cm.ctx, state) {
			return
		}
		cm.updateState(state)
	}
}

func (cm *ConnManager) SetDialOptions(opts []grpc.DialOption) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.opts = opts
}

func (cm *ConnManager) GetConn() *grpc.ClientConn {
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

func (cm *ConnManager) updateState(state connectivity.State) {
	endpoint := cm.getNextEndpoint()

	cm.mu.Lock()
	oldState := cm.currentState.State
	cm.currentState = ConnectionState{
		State:       state,
		LastChanged: time.Now(),
		Endpoint:    endpoint,
	}
	cm.mu.Unlock()

	if oldState != state {
		logger.Debug("Connection state changed from %s to %s", oldState, state)
		select {
		case cm.stateChange <- cm.currentState:
		case <-cm.ctx.Done():
			return
		default:
			// Don't warn about dropped notifications during shutdown
			if cm.ctx.Err() == nil {
				logger.Warn("State change notification dropped - channel full")
			}
		}
	}
}

// GetAttempts returns the number of connection attempts made
func (cm *ConnManager) GetAttempts() int64 {
	return atomic.LoadInt64(&cm.attempts)
}

// IsRetrying returns whether a retry operation is in progress
func (cm *ConnManager) IsRetrying() bool {
	return atomic.LoadInt32(&cm.retrying) == 1
}

// ResetBackoff resets the backoff to its initial state
func (cm *ConnManager) ResetBackoff() {
	if b, ok := cm.backoff.(*backoff.ExponentialBackOff); ok {
		b.Reset()
	}
}
