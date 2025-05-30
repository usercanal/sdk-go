// sdk-go/examples/unified/main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Single client for both analytics events and structured logging
	client, err := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
		Endpoint:      "collect.usercanal.com:9000",
		BatchSize:     200,
		FlushInterval: 5 * time.Second,
		MaxRetries:    3,
		Debug:         true,
	})
	if err != nil {
		// This is the only place we use standard log since client creation failed
		fmt.Printf("Failed to create unified client: %v\n", err)
		return
	}
	defer client.Close()

	ctx := context.Background()
	hostname, _ := os.Hostname()
	userID := "user123"

	// 1. Log application event
	client.LogInfo(ctx, "user-service", "User initiated data export")

	// 2. Track analytics event for the same action
	if err := client.TrackEvent(ctx, usercanal.Event{
		UserId: userID,
		Name:   usercanal.FeatureUsed,
		Properties: usercanal.Properties{
			"feature":     "data_export",
			"export_type": "csv",
		},
	}); err != nil {
		// Log SDK errors through the SDK itself!
		client.LogError(ctx, "sdk-events", fmt.Errorf("failed to track feature usage: %w", err))
	}

	// 3. Detailed structured logging with correlation
	if err := client.SendLog(ctx, usercanal.LogEntry{
		EventType: usercanal.LogCollect,
		Level:     usercanal.LogInfo,
		Service:   "export-service",
		Source:    hostname,
		Message:   "Data export processing started",
		Data: map[string]interface{}{
			"user_id":      userID,
			"export_id":    "exp_789",
			"record_count": 1500,
			"format":       "csv",
		},
	}); err != nil {
		client.LogError(ctx, "sdk-logs", fmt.Errorf("failed to send structured log: %w", err))
	}

	// Simulate processing time
	time.Sleep(2 * time.Second)

	// 4. Log successful completion
	client.LogInfo(ctx, "export-service", "Data export completed successfully")

	// 5. Track completion event for analytics
	if err := client.TrackEvent(ctx, usercanal.Event{
		UserId: userID,
		Name:   usercanal.FeatureUsed,
		Properties: usercanal.Properties{
			"feature":         "data_export",
			"status":          "completed",
			"processing_time": 2000,
			"file_size":       1024,
		},
	}); err != nil {
		client.LogError(ctx, "sdk-events", fmt.Errorf("failed to track completion: %w", err))
	}

	// 6. Handle an error scenario
	simulatedError := errors.New("temporary service unavailable")

	// Log the application error
	if err := client.SendLog(ctx, usercanal.LogEntry{
		EventType: usercanal.LogCollect,
		Level:     usercanal.LogError,
		Service:   "notification-service",
		Source:    hostname,
		Message:   "Failed to send notification",
		Data: map[string]interface{}{
			"user_id":     userID,
			"error":       simulatedError.Error(),
			"retry_count": 3,
		},
	}); err != nil {
		// Even SDK logging errors get logged through the SDK!
		client.LogCritical(ctx, "sdk-logs", fmt.Sprintf("Failed to log application error: %v", err))
	}

	// Track the error for analytics (optional - for error rate metrics)
	if err := client.TrackEvent(ctx, usercanal.Event{
		UserId: userID,
		Name:   "error_occurred",
		Properties: usercanal.Properties{
			"error_type":  "notification_failure",
			"service":     "notification-service",
			"recoverable": true,
		},
	}); err != nil {
		client.LogError(ctx, "sdk-events", fmt.Errorf("failed to track error event: %w", err))
	}

	// 7. User identification for analytics
	if err := client.IdentifyUser(ctx, usercanal.Identity{
		UserId: userID,
		Properties: usercanal.Properties{
			"name":          "John Doe",
			"email":         "john@example.com",
			"last_activity": time.Now(),
			"feature_usage": map[string]int{"data_export": 5},
		},
	}); err != nil {
		client.LogError(ctx, "sdk-events", fmt.Errorf("failed to identify user: %w", err))
	}

	// 8. Demonstrate SDK error handling - what happens when we have connection issues?
	// Force a flush to show how even flush errors get logged
	if err := client.Flush(ctx); err != nil {
		// Log flush failures through emergency logging
		client.SendLog(ctx, usercanal.LogEntry{
			EventType: usercanal.LogCollect,
			Level:     usercanal.LogEmergency,
			Service:   "sdk-transport",
			Source:    hostname,
			Message:   "Critical: Failed to flush data to UserCanal",
			Data: map[string]interface{}{
				"error":           err.Error(),
				"pending_events":  "unknown", // Could add stats here
				"timestamp":       time.Now(),
				"recovery_action": "data may be retried on next flush",
			},
		})
	}

	// Show unified statistics
	client.DumpStatus()

	// Final log message
	client.LogInfo(ctx, "example-app", "Unified events + logs example completed successfully")

	// Final flush attempt
	client.Flush(ctx)
}
