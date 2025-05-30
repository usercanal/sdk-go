package types

import "time"

type TransportMetrics struct {
	// Separate counters
	EventsSent       int64
	EventBatchesSent int64
	LogsSent         int64
	LogBatchesSent   int64

	// Combined totals
	TotalBatchesSent int64
	BytesSent        int64
	FailedAttempts   int64

	// Timing
	LastSendTime     time.Time
	LastFailureTime  time.Time
	ConnectionUptime time.Duration
	ReconnectCount   int64

	// Calculated fields
	AverageEventBatchSize float64
	AverageLogBatchSize   float64
}
