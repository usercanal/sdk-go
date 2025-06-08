// Package config provides centralized configuration defaults for the UserCanal SDK.
package config

import "time"

const (
	// DefaultEndpoint is the canonical production endpoint for UserCanal
	DefaultEndpoint = "collect.usercanal.com:50000"
	
	// DefaultBatchSize is the default number of events/logs per batch
	DefaultBatchSize = 100
	
	// DefaultFlushInterval is the default time between batch sends
	DefaultFlushInterval = 10 * time.Second
	
	// DefaultMaxRetries is the default number of retry attempts
	DefaultMaxRetries = 3
	
	// DefaultCloseTimeout is the default timeout for graceful shutdown
	DefaultCloseTimeout = 5 * time.Second
	
	// DefaultDebug is the default debug logging state
	DefaultDebug = false
)

// Defaults returns a map of all default configuration values
// This is useful for documentation and testing
func Defaults() map[string]interface{} {
	return map[string]interface{}{
		"endpoint":       DefaultEndpoint,
		"batch_size":     DefaultBatchSize,
		"flush_interval": DefaultFlushInterval,
		"max_retries":    DefaultMaxRetries,
		"close_timeout":  DefaultCloseTimeout,
		"debug":          DefaultDebug,
	}
}