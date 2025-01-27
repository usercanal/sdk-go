// transport/connection.go
// Package transport provides connection management for the Usercanal SDK.
//
// ConnectionManager's Core Responsibilities:
// 1. Maintain stable gRPC connection with automatic recovery
// 2. Support DNS-based high availability and load balancing
// 3. Implement exponential backoff for reconnection attempts
// 4. Provide connection state monitoring and notifications
//
// Product Requirements:
// - Zero data loss during collector upgrades/outages
// - Automatic failover between multiple collectors
// - Graceful handling of network issues
// - Clear status reporting for debugging
//
// Internal behaviors:
// - DNS Resolution: Periodic refresh (10m) with backoff retries on failure
// - Connection States: IDLE -> CONNECTING -> READY/TRANSIENT_FAILURE
// - Backoff Strategy: Exponential with jitter (1s base, 1.5x multiplier)
// - High Availability: Round-robin across resolved endpoints

package transport

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/usercanal/sdk-go/internal/logger"
	"github.com/usercanal/sdk-go/types"
)

const (
	defaultBaseDelay          = 1 * time.Second
	defaultMultiplier         = 1.5
	defaultJitter             = 0.2
	defaultDNSRefreshInterval = 10 * time.Minute
	defaultKeepAliveTime      = 10 * time.Second
	defaultKeepAliveTimeout   = 3 * time.Second
	defaultGRPCPort           = "50051"
	maxDNSRetries             = 3
	dnsRetryBaseDelay         = 500 * time.Millisecond
	monitorStateTimeout       = 5 * time.Second
	reconnectTimeout          = 30 * time.Second
)

var (
	ErrConnectionClosed = types.NewValidationError("connection", "is closed")
)

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

	// Backoff config
	baseDelay  time.Duration
	multiplier float64
	jitter     float64

	// DNS management
	resolvedIPs       []string
	lastDNSResolution time.Time
	currentIPIndex    int

	// State tracking
	attempts     int64 // using atomic for thread safety
	currentState ConnectionState
	mu           sync.RWMutex
	closed       bool
	reconnecting int32 // atomic flag for reconnection status

	// Control and cleanup
	cancelCtx        context.Context
	cancelFunc       context.CancelFunc
	wg               sync.WaitGroup
	dnsRefreshTicker *time.Ticker

	// State change notification
	stateChange chan ConnectionState
}

func NewConnManager(endpoint string) *ConnManager {
	ctx, cancel := context.WithCancel(context.Background())

	cm := &ConnManager{
		endpoint:    endpoint,
		baseDelay:   defaultBaseDelay,
		multiplier:  defaultMultiplier,
		jitter:      defaultJitter,
		cancelCtx:   ctx,
		cancelFunc:  cancel,
		stateChange: make(chan ConnectionState, 1),
	}

	// Initialize state
	cm.currentState = ConnectionState{
		State:       connectivity.Idle,
		LastChanged: time.Now(),
		Endpoint:    endpoint,
	}

	// Start DNS refresh routine
	cm.wg.Add(1)
	go cm.startDNSRefresh()

	return cm
}

func (cm *ConnManager) SetDialOptions(opts []grpc.DialOption) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.opts = opts
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
	cm.mu.Lock()
	oldState := cm.currentState.State
	cm.currentState = ConnectionState{
		State:       state,
		LastChanged: time.Now(),
		Endpoint:    cm.getNextEndpoint(),
	}
	cm.mu.Unlock()

	if oldState != state {
		logger.Debug("Connection state changed from %s to %s", oldState, state)
		select {
		case cm.stateChange <- cm.currentState:
		case <-cm.cancelCtx.Done():
			return
		default:
			logger.Warn("State change notification dropped - channel full")
		}
	}
}

func (cm *ConnManager) resolveEndpoint() error {
	var lastErr error
	for attempt := 0; attempt < maxDNSRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(attempt) * dnsRetryBaseDelay
			timer := time.NewTimer(backoff)
			select {
			case <-cm.cancelCtx.Done():
				timer.Stop()
				return cm.cancelCtx.Err()
			case <-timer.C:
			}
		}

		// Only refresh if TTL expired
		if time.Since(cm.lastDNSResolution) < defaultDNSRefreshInterval && len(cm.resolvedIPs) > 0 {
			return nil
		}

		host := cm.endpoint
		port := defaultGRPCPort
		if h, p, err := net.SplitHostPort(cm.endpoint); err == nil {
			host = h
			port = p
		}

		ips, err := net.LookupHost(host)
		if err != nil {
			lastErr = err
			logger.Warn("DNS resolution attempt %d failed: %v", attempt+1, err)
			continue
		}

		cm.mu.Lock()
		cm.resolvedIPs = make([]string, len(ips))
		for i, ip := range ips {
			cm.resolvedIPs[i] = net.JoinHostPort(ip, port)
		}
		cm.lastDNSResolution = time.Now()
		cm.mu.Unlock()

		logger.Debug("Resolved %s to %d endpoints", host, len(ips))
		return nil
	}

	return fmt.Errorf("DNS resolution failed after %d attempts: %v", maxDNSRetries, lastErr)
}

func (cm *ConnManager) getNextEndpoint() string {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if len(cm.resolvedIPs) == 0 {
		return cm.endpoint // fallback to original endpoint
	}

	endpoint := cm.resolvedIPs[cm.currentIPIndex]
	cm.currentIPIndex = (cm.currentIPIndex + 1) % len(cm.resolvedIPs)
	return endpoint
}

func (cm *ConnManager) Connect(ctx context.Context) error {
	if cm.isClosed() {
		return ErrConnectionClosed
	}

	cm.mu.Lock()
	if cm.conn != nil {
		cm.conn.Close()
		cm.conn = nil
	}
	cm.mu.Unlock()

	attempt := atomic.AddInt64(&cm.attempts, 1)

	// Resolve DNS with retries
	if err := cm.resolveEndpoint(); err != nil {
		logger.Warn("DNS resolution failed, using original endpoint: %v", err)
	}

	backoff := cm.calculateBackoff(int(attempt))
	if backoff > 0 {
		logger.Debug("Backing off for %v (attempt %d)", backoff, attempt)
		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-cm.cancelCtx.Done():
			timer.Stop()
			return ErrConnectionClosed
		case <-timer.C:
		}
	}

	endpoint := cm.getNextEndpoint()
	logger.Debug("Connection attempt %d starting (backoff: %v)", attempt, backoff)

	conn, err := grpc.DialContext(ctx, endpoint, cm.opts...)
	if err != nil {
		logger.Error("Connection attempt %d failed: %v", attempt, err)
		cm.updateState(connectivity.TransientFailure)
		return fmt.Errorf("connection attempt %d failed: %w", attempt, err)
	}

	initialState := conn.GetState()
	logger.Debug("Initial connection state: %s", initialState)

	cm.mu.Lock()
	cm.conn = conn
	cm.updateState(conn.GetState())
	cm.mu.Unlock()

	// Start monitoring if not already running
	cm.wg.Add(1)
	go cm.monitor()

	return nil
}

func (cm *ConnManager) calculateBackoff(attempt int) time.Duration {
	if attempt <= 1 {
		return 0
	}

	backoff := float64(cm.baseDelay) * math.Pow(cm.multiplier, float64(attempt-2))
	delta := cm.jitter * backoff
	min := backoff - delta
	max := backoff + delta
	backoff = min + (rand.Float64() * (max - min))

	return time.Duration(backoff)
}

func (cm *ConnManager) monitor() {
	defer cm.wg.Done()

	for {
		select {
		case <-cm.cancelCtx.Done():
			return
		default:
			cm.mu.RLock()
			if cm.conn == nil {
				cm.mu.RUnlock()
				return
			}
			conn := cm.conn
			cm.mu.RUnlock()

			state := conn.GetState()
			// Add more detailed logging
			logger.Debug("Current connection state: %s (attempt: %d)",
				state, atomic.LoadInt64(&cm.attempts))

			cm.updateState(state)

			if state == connectivity.TransientFailure {
				logger.Warn("Connection failed, initiating reconnect...")
				cm.tryReconnect()
			}

			stateCtx, cancel := context.WithTimeout(cm.cancelCtx, monitorStateTimeout)
			conn.WaitForStateChange(stateCtx, state)
			cancel()
		}
	}
}

func (cm *ConnManager) tryReconnect() {
	if !atomic.CompareAndSwapInt32(&cm.reconnecting, 0, 1) {
		return
	}
	defer atomic.StoreInt32(&cm.reconnecting, 0)

	// Add immediate logging
	logger.Debug("Starting reconnection attempt...")

	ctx, cancel := context.WithTimeout(cm.cancelCtx, reconnectTimeout)
	defer cancel()

	if err := cm.Connect(ctx); err != nil {
		logger.Error("Reconnection failed: %v, will retry in background", err)
		// Force another reconnect attempt
		go func() {
			time.Sleep(time.Second)
			cm.tryReconnect()
		}()
	}
}

func (cm *ConnManager) GetConn() *grpc.ClientConn {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.conn
}

func (cm *ConnManager) isClosed() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.closed
}

func (cm *ConnManager) Close() error {
	cm.mu.Lock()
	if cm.closed {
		cm.mu.Unlock()
		return nil
	}
	cm.closed = true
	cm.mu.Unlock()

	// Cancel context to stop all goroutines
	cm.cancelFunc()

	// Close state change channel after context cancellation
	close(cm.stateChange)

	// Wait for all goroutines to finish
	cm.wg.Wait()

	// Stop DNS refresh ticker if it exists
	if cm.dnsRefreshTicker != nil {
		cm.dnsRefreshTicker.Stop()
	}

	// Close connection if it exists
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.conn != nil {
		return cm.conn.Close()
	}
	return nil
}

func (cm *ConnManager) startDNSRefresh() {
	defer cm.wg.Done()

	// Initial DNS resolution
	if err := cm.resolveEndpoint(); err != nil {
		logger.Warn("Initial DNS resolution failed: %v", err)
	}

	cm.dnsRefreshTicker = time.NewTicker(defaultDNSRefreshInterval)
	defer cm.dnsRefreshTicker.Stop()

	for {
		select {
		case <-cm.cancelCtx.Done():
			return
		case <-cm.dnsRefreshTicker.C:
			if err := cm.resolveEndpoint(); err != nil {
				logger.Warn("DNS refresh failed: %v", err)
			}
		}
	}
}
