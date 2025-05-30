// api/stats.go
package api

import (
	"github.com/usercanal/sdk-go/internal/logger"
	"github.com/usercanal/sdk-go/types"
)

// GetStats returns current client statistics
func (c *Client) GetStats() types.Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Get transport-level metrics
	transportMetrics := c.sender.GetMetrics()

	// Get connection state (would need to add this method)
	// connInfo := c.sender.GetConnectionInfo()

	// Compose client-level stats from multiple sources
	return types.Stats{
		// Queue info from batch managers
		EventsInQueue: int64(c.eventBatcher.QueueSize()),
		LogsInQueue:   int64(c.logBatcher.QueueSize()),

		// Summary from transport metrics
		EventsSent:   transportMetrics.EventsSent,
		LogsSent:     transportMetrics.LogsSent,
		EventsFailed: transportMetrics.FailedAttempts,

		// Connection from transport
		ConnectionState:  c.sender.State(),
		ConnectionUptime: transportMetrics.ConnectionUptime,

		// Timing from various sources
		LastFlushTime:    c.eventBatcher.LastFlushTime(),   // Client-level timing
		LastFailureTime:  transportMetrics.LastFailureTime, // Transport-level timing
		AverageBatchSize: (transportMetrics.AverageEventBatchSize + transportMetrics.AverageLogBatchSize) / 2,

		// Client config
		ActiveEndpoint: c.cfg.endpoint,
		// DNS info would come from connection manager when available
	}
}

// DumpStatus prints detailed status information
func (c *Client) DumpStatus() {
	stats := c.GetStats()

	logger.Info("UserCanal Status Report")
	logger.Info("=====================")
	logger.Info("Connection State: %s", stats.ConnectionState)
	logger.Info("Connection Uptime: %v", stats.ConnectionUptime)
	logger.Info("Events in Queue: %d", stats.EventsInQueue)
	logger.Info("Events Sent: %d", stats.EventsSent)
	logger.Info("Failed Events: %d", stats.EventsFailed)
	logger.Info("Average Batch Size: %.2f", stats.AverageBatchSize)
	logger.Info("Last Flush: %v", stats.LastFlushTime)
	logger.Info("Last Failure: %v", stats.LastFailureTime)
}
