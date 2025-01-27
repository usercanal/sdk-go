// types/stats.go
package types

import "time"

// Stats represents client statistics
type Stats struct {
	EventsInQueue    int64
	EventsSent       int64
	EventsFailed     int64
	ConnectionState  string
	ConnectionUptime time.Duration
	LastFlushTime    time.Time
	LastFailureTime  time.Time
	AverageBatchSize float64

	ActiveEndpoint    string    // Current endpoint being used
	ResolvedEndpoints []string  // All available endpoints
	LastDNSResolution time.Time // Last successful DNS resolution
	DNSFailures       int64     // Number of DNS resolution failures
}
