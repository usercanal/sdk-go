// sdk-go/internal/api/logs.go
package api

import (
	"context"
	"fmt"
	"time"

	"github.com/usercanal/sdk-go/internal/convert"
	"github.com/usercanal/sdk-go/types"
)

// Log sends a single log entry
func (c *Client) Log(ctx context.Context, entry types.LogEntry) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := entry.Validate(); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	// Set timestamp if not set
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	transportLog, err := convert.LogToInternal(&entry)
	if err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	if err := c.logBatcher.Add(ctx, transportLog); err != nil {
		return fmt.Errorf("failed to add log entry: %w", err)
	}

	return nil
}

// LogBatch sends multiple log entries
func (c *Client) LogBatch(ctx context.Context, entries []types.LogEntry) error {
	for i, entry := range entries {
		if err := c.Log(ctx, entry); err != nil {
			return fmt.Errorf("failed to add log entry[%d]: %w", i, err)
		}
	}
	return nil
}

// Additional convenience log methods could go here:

// LogInfo sends an info-level log entry
func (c *Client) LogInfo(ctx context.Context, service, source, message string, data map[string]interface{}) error {
	return c.Log(ctx, types.LogEntry{
		Level:     types.LogInfo,
		EventType: types.LogCollect,
		Service:   service,
		Source:    source,
		Message:   message,
		Data:      data,
	})
}

// LogError sends an error-level log entry
func (c *Client) LogError(ctx context.Context, service, source, message string, data map[string]interface{}) error {
	return c.Log(ctx, types.LogEntry{
		Level:     types.LogError,
		EventType: types.LogCollect,
		Service:   service,
		Source:    source,
		Message:   message,
		Data:      data,
	})
}

// LogDebug sends a debug-level log entry
func (c *Client) LogDebug(ctx context.Context, service, source, message string, data map[string]interface{}) error {
	return c.Log(ctx, types.LogEntry{
		Level:     types.LogDebug,
		EventType: types.LogCollect,
		Service:   service,
		Source:    source,
		Message:   message,
		Data:      data,
	})
}
