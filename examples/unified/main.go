// sdk-go/examples/unified/main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Single client for both analytics events and structured logging
	client, err := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
		Endpoint:      "collect.usercanal.com:50000",
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
	defer func() {
		if err := client.Close(context.Background()); err != nil {
			log.Printf("Failed to close client: %v", err)
		}
	}()

	ctx := context.Background()
	hostname, _ := os.Hostname()
	userID := "user123"

	// 1. Log application event
	client.LogInfo(ctx, "user-service", "User initiated data export", nil)

	// 2. Track analytics event for the same action
	if err := client.Event(ctx, userID, usercanal.FeatureUsed, usercanal.Properties{
		"feature":     "data_export",
		"export_type": "csv",
	}); err != nil {
		// Log SDK errors through the SDK itself!
		client.LogError(ctx, "sdk-events", fmt.Sprintf("failed to track feature usage: %v", err), nil)
	}

	// 3. Detailed structured logging with correlation
	if err := client.Log(ctx, usercanal.LogEntry{
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
		client.LogError(ctx, "sdk-logs", fmt.Sprintf("failed to send structured log: %v", err), nil)
	}

	// Simulate processing time
	time.Sleep(2 * time.Second)

	// 4. Log successful completion
	client.LogInfo(ctx, "export-service", "Data export completed successfully", nil)

	// 5. Track completion event for analytics
	if err := client.Event(ctx, userID, usercanal.FeatureUsed, usercanal.Properties{
		"feature":         "data_export",
		"status":          "completed",
		"processing_time": 2000,
		"file_size":       1024,
	}); err != nil {
		client.LogError(ctx, "sdk-events", fmt.Sprintf("failed to track completion: %v", err), nil)
	}

	// 6. Handle an error scenario
	simulatedError := errors.New("temporary service unavailable")

	// Log the application error
	if err := client.Log(ctx, usercanal.LogEntry{
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
		client.LogCritical(ctx, "sdk-logs", fmt.Sprintf("Failed to log application error: %v", err), nil)
	}

	// Track the error for analytics (optional - for error rate metrics)
	if err := client.Event(ctx, userID, "error_occurred", usercanal.Properties{
		"error_type":  "notification_failure",
		"service":     "notification-service",
		"recoverable": true,
	}); err != nil {
		client.LogError(ctx, "sdk-events", fmt.Sprintf("failed to track error event: %v", err), nil)
	}

	// 7. User identification for analytics
	if err := client.EventIdentify(ctx, userID, usercanal.Properties{
		"name":          "John Doe",
		"email":         "john@example.com",
		"last_activity": time.Now(),
		"feature_usage": map[string]int{"data_export": 5},
	}); err != nil {
		client.LogError(ctx, "sdk-events", fmt.Sprintf("failed to identify user: %v", err), nil)
	}

	// 8. Demonstrate SDK error handling - what happens when we have connection issues?
	// Force a flush to show how even flush errors get logged
	if err := client.Flush(ctx); err != nil {
		// Log flush failures through emergency logging
		client.Log(ctx, usercanal.LogEntry{
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
	stats := client.GetStats()
	log.Printf("Client Stats: Events queued: %d, Connection: %s", stats.EventsInQueue, stats.ConnectionState)

	// Final log message
	client.LogInfo(ctx, "example-app", "Unified events + logs example completed successfully", nil)

	// Final flush attempt
	client.Flush(ctx)

	// Close the client
	if err := client.Close(ctx); err != nil {
		log.Printf("Failed to close client: %v", err)
	}
}
