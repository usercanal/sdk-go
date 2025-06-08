// sdk-go/internal/api/logs.go
package api

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/usercanal/sdk-go/internal/convert"
	"github.com/usercanal/sdk-go/types"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
}

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
	if source == "" {
		source = hostname
	}
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
	if source == "" {
		source = hostname
	}
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
	if source == "" {
		source = hostname
	}
	return c.Log(ctx, types.LogEntry{
		Level:     types.LogDebug,
		EventType: types.LogCollect,
		Service:   service,
		Source:    source,
		Message:   message,
		Data:      data,
	})
}

// LogWarning sends a warning-level log entry
func (c *Client) LogWarning(ctx context.Context, service, source, message string, data map[string]interface{}) error {
	if source == "" {
		source = hostname
	}
	return c.Log(ctx, types.LogEntry{
		Level:     types.LogWarning,
		EventType: types.LogCollect,
		Service:   service,
		Source:    source,
		Message:   message,
		Data:      data,
	})
}

// LogCritical sends a critical-level log entry
func (c *Client) LogCritical(ctx context.Context, service, source, message string, data map[string]interface{}) error {
	if source == "" {
		source = hostname
	}
	return c.Log(ctx, types.LogEntry{
		Level:     types.LogCritical,
		EventType: types.LogCollect,
		Service:   service,
		Source:    source,
		Message:   message,
		Data:      data,
	})
}

// LogAlert sends an alert-level log entry
func (c *Client) LogAlert(ctx context.Context, service, source, message string, data map[string]interface{}) error {
	if source == "" {
		source = hostname
	}
	return c.Log(ctx, types.LogEntry{
		Level:     types.LogAlert,
		EventType: types.LogCollect,
		Service:   service,
		Source:    source,
		Message:   message,
		Data:      data,
	})
}

// LogEmergency sends an emergency-level log entry
func (c *Client) LogEmergency(ctx context.Context, service, source, message string, data map[string]interface{}) error {
	if source == "" {
		source = hostname
	}
	return c.Log(ctx, types.LogEntry{
		Level:     types.LogEmergency,
		EventType: types.LogCollect,
		Service:   service,
		Source:    source,
		Message:   message,
		Data:      data,
	})
}

// LogNotice sends a notice-level log entry
func (c *Client) LogNotice(ctx context.Context, service, source, message string, data map[string]interface{}) error {
	if source == "" {
		source = hostname
	}
	return c.Log(ctx, types.LogEntry{
		Level:     types.LogNotice,
		EventType: types.LogCollect,
		Service:   service,
		Source:    source,
		Message:   message,
		Data:      data,
	})
}

// LogTrace sends a trace-level log entry
func (c *Client) LogTrace(ctx context.Context, service, source, message string, data map[string]interface{}) error {
	if source == "" {
		source = hostname
	}
	return c.Log(ctx, types.LogEntry{
		Level:     types.LogTrace,
		EventType: types.LogCollect,
		Service:   service,
		Source:    source,
		Message:   message,
		Data:      data,
	})
}
