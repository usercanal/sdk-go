// sdk-go/examples/logs/advanced/main.go
package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Initialize with configuration optimized for high-volume logging
	client, err := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
		Endpoint:      "collect.usercanal.com:9000",
		BatchSize:     500, // Higher batch size for logs
		FlushInterval: 2 * time.Second,
		MaxRetries:    5,
		Debug:         false, // Don't log debug for production logging
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	hostname, _ := os.Hostname()
	sessionID := uint64(time.Now().UnixNano()) // Simple session ID

	// Application startup log
	err = client.SendLog(ctx, usercanal.LogEntry{
		EventType: usercanal.LogCollect,
		ContextID: sessionID,
		Level:     usercanal.LogInfo,
		Service:   "api-gateway",
		Source:    hostname,
		Message:   "Service started successfully",
		Data: map[string]interface{}{
			"version":    "v1.2.3",
			"port":       8080,
			"env":        "production",
			"start_time": time.Now(),
		},
	})
	if err != nil {
		log.Printf("Failed to log startup: %v", err)
	}

	// Request processing logs with correlation
	requestID := "req_789"

	// Request start
	err = client.SendLog(ctx, usercanal.LogEntry{
		EventType: usercanal.LogCollect,
		ContextID: sessionID,
		Level:     usercanal.LogInfo,
		Service:   "api-gateway",
		Source:    hostname,
		Message:   "Processing API request",
		Data: map[string]interface{}{
			"request_id": requestID,
			"method":     "POST",
			"path":       "/api/users",
			"user_id":    "user123",
			"ip":         "192.168.1.100",
			"user_agent": "UserCanal-SDK/1.0",
		},
	})
	if err != nil {
		log.Printf("Failed to log request start: %v", err)
	}

	// Simulate some processing and logging
	time.Sleep(100 * time.Millisecond)

	// Database operation log
	err = client.SendLog(ctx, usercanal.LogEntry{
		EventType: usercanal.LogCollect,
		ContextID: sessionID,
		Level:     usercanal.LogDebug,
		Service:   "database",
		Source:    hostname,
		Message:   "Database query executed",
		Data: map[string]interface{}{
			"request_id":    requestID,
			"query":         "SELECT * FROM users WHERE id = ?",
			"duration_ms":   45,
			"rows_affected": 1,
		},
	})
	if err != nil {
		log.Printf("Failed to log database operation: %v", err)
	}

	// Authentication log (special event type)
	err = client.SendLog(ctx, usercanal.LogEntry{
		EventType: usercanal.LogAuth, // Special routing for security events
		ContextID: sessionID,
		Level:     usercanal.LogNotice,
		Service:   "auth-service",
		Source:    hostname,
		Message:   "User authentication attempt",
		Data: map[string]interface{}{
			"request_id":  requestID,
			"user_id":     "user123",
			"auth_method": "jwt",
			"success":     true,
			"ip":          "192.168.1.100",
			"session_id":  "sess_456",
		},
	})
	if err != nil {
		log.Printf("Failed to log auth event: %v", err)
	}

	// Error handling and logging
	simulatedError := errors.New("connection timeout")
	err = client.LogError(ctx, "external-api", simulatedError)
	if err != nil {
		log.Printf("Failed to log error: %v", err)
	}

	// Request completion
	err = client.SendLog(ctx, usercanal.LogEntry{
		EventType: usercanal.LogCollect,
		ContextID: sessionID,
		Level:     usercanal.LogInfo,
		Service:   "api-gateway",
		Source:    hostname,
		Message:   "Request completed",
		Data: map[string]interface{}{
			"request_id":    requestID,
			"status_code":   200,
			"duration_ms":   150,
			"response_size": 1024,
		},
	})
	if err != nil {
		log.Printf("Failed to log request completion: %v", err)
	}

	// Batch multiple logs for efficiency
	logBatch := []usercanal.LogEntry{
		{
			EventType: usercanal.LogCollect,
			Level:     usercanal.LogInfo,
			Service:   "metrics-collector",
			Source:    hostname,
			Message:   "System metrics collected",
			Data: map[string]interface{}{
				"cpu_usage":    45.2,
				"memory_usage": 67.8,
				"disk_usage":   23.1,
			},
		},
		{
			EventType: usercanal.LogCollect,
			Level:     usercanal.LogInfo,
			Service:   "health-checker",
			Source:    hostname,
			Message:   "Health check passed",
			Data: map[string]interface{}{
				"check_type": "database",
				"latency_ms": 12,
				"status":     "healthy",
			},
		},
	}

	err = client.SendLogs(ctx, logBatch)
	if err != nil {
		log.Printf("Failed to send log batch: %v", err)
	}

	// Final flush
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}

	log.Println("Advanced logging example completed")
}
