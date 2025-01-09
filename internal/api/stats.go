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

	metrics := c.sender.GetMetrics()

	return types.Stats{
		EventsInQueue:    int64(c.batcher.QueueSize()),
		EventsSent:       metrics.EventsSent,
		EventsFailed:     metrics.FailedAttempts,
		ConnectionState:  c.sender.State(),
		ConnectionUptime: c.sender.Uptime(),
		LastFlushTime:    c.batcher.LastFlushTime(),
		LastFailureTime:  metrics.LastFailureTime,
		AverageBatchSize: metrics.AverageBatchSize,
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
