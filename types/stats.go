// sdk-go/types/stats.go - Client-facing stats assembled from multiple sources
package types

import "time"

type Stats struct {
	// Client queue state (from batch managers)
	EventsInQueue int64
	LogsInQueue   int64

	// Summary counters (from transport metrics)
	EventsSent   int64
	LogsSent     int64
	EventsFailed int64

	// Client connection view
	ConnectionState  string
	ConnectionUptime time.Duration

	// Client timing (from batch managers + transport)
	LastFlushTime    time.Time
	LastFailureTime  time.Time
	AverageBatchSize float64

	// Client network view (from connection manager + config)
	ActiveEndpoint    string
	ResolvedEndpoints []string
	LastDNSResolution time.Time
	DNSFailures       int64
}
