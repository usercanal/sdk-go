// transport/connection.go
package transport

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/usercanal/sdk-go/internal/logger"
	"github.com/usercanal/sdk-go/types"
)

type ConnectionState string

var (
	ErrConnectionClosed = types.NewValidationError("connection", "is closed")
)

const (
	defaultBaseDelay          = 1 * time.Second
	defaultMultiplier         = 1.5
	defaultJitter             = 0.2
	defaultDNSRefreshInterval = 10 * time.Minute // Match DNS TTL
	defaultKeepAliveTime      = 10 * time.Second
	defaultKeepAliveTimeout   = 3 * time.Second
	defaultGRPCPort           = "50051"
	maxDNSRetries             = 3
	dnsRetryBaseDelay         = 500 * time.Millisecond
	monitorStateTimeout       = 5 * time.Second
)

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
	attempts int
	mu       sync.RWMutex
	closed   bool

	// Control
	done chan struct{}

	// New fields for cleanup
	dnsRefreshTicker *time.Ticker
	cancelCtx        context.Context
	cancelFunc       context.CancelFunc
	wg               sync.WaitGroup // track goroutines
}

func NewConnManager(endpoint string) *ConnManager {
	ctx, cancel := context.WithCancel(context.Background())

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                defaultKeepAliveTime,
			Timeout:             defaultKeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	}

	cm := &ConnManager{
		endpoint:   endpoint,
		opts:       opts,
		baseDelay:  defaultBaseDelay,
		multiplier: defaultMultiplier,
		jitter:     defaultJitter,
		cancelCtx:  ctx,
		cancelFunc: cancel,
	}

	// Start DNS refresh routine with WaitGroup tracking
	cm.wg.Add(1)
	go cm.startDNSRefresh()

	return cm
}

func (cm *ConnManager) SetDialOptions(opts []grpc.DialOption) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.opts = opts
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
		if time.Since(cm.lastDNSResolution) < defaultDNSRefreshInterval {
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
	cm.mu.Lock()
	if cm.closed {
		cm.mu.Unlock()
		return ErrConnectionClosed
	}

	if cm.conn != nil {
		cm.conn.Close()
		cm.conn = nil
	}

	cm.attempts++
	currentAttempt := cm.attempts
	cm.mu.Unlock()

	// Resolve DNS with retries
	if err := cm.resolveEndpoint(); err != nil {
		logger.Warn("DNS resolution failed, using original endpoint: %v", err)
	}

	backoff := cm.calculateBackoff(currentAttempt)
	if backoff > 0 {
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
	logger.Debug("Attempting connection to %s (attempt %d)", endpoint, currentAttempt)

	conn, err := grpc.DialContext(ctx, endpoint, cm.opts...)
	if err != nil {
		logger.Warn("Connection attempt %d to %s failed: %v", currentAttempt, endpoint, err)
		return err
	}

	cm.mu.Lock()
	cm.conn = conn
	cm.mu.Unlock()

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

func (cm *ConnManager) startDNSRefresh() {
	defer cm.wg.Done()

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
			if state == connectivity.TransientFailure {
				logger.Warn("Connection in TransientFailure, attempting reconnect")
				go cm.reconnect()
			}

			stateCtx, cancel := context.WithTimeout(cm.cancelCtx, monitorStateTimeout)
			conn.WaitForStateChange(stateCtx, state)
			cancel()
		}
	}
}

func (cm *ConnManager) reconnect() {
	cm.wg.Add(1)
	defer cm.wg.Done()

	ctx, cancel := context.WithTimeout(cm.cancelCtx, 30*time.Second)
	defer cancel()

	if err := cm.Connect(ctx); err != nil {
		logger.Error("Reconnection failed: %v", err)
	}
}

func (cm *ConnManager) GetConn() *grpc.ClientConn {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.conn
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

	// Wait for all goroutines to finish
	cm.wg.Wait()

	// Stop DNS refresh ticker if it exists
	if cm.dnsRefreshTicker != nil {
		cm.dnsRefreshTicker.Stop()
	}

	// Close connection if it exists
	if cm.conn != nil {
		return cm.conn.Close()
	}
	return nil
}
