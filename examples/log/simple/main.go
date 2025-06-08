// sdk-go/examples/log/simple/main.go
package main

import (
	"context"
	"log"
	"time"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Create client
	client, err := usercanal.NewClient("000102030405060708090a0b0c0d0e0f", usercanal.Config{
		Endpoint:      "localhost:50000",
		Debug:         true,
		FlushInterval: 1 * time.Second, // Shorter flush interval for testing
		BatchSize:     1,               // Send immediately for testing
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() {
		if err := client.Close(context.Background()); err != nil {
			log.Printf("Failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	// Super simple logging
	client.LogInfo(ctx, "my-app", "Application started", nil)

	client.LogError(ctx, "my-app", "Login failed", map[string]interface{}{
		"user_id": "123",
		"reason":  "invalid_password",
	})

	client.LogDebug(ctx, "my-app", "Processing request", map[string]interface{}{
		"request_id": "req_456",
		"duration":   "45ms",
	})

	// IMPORTANT: Flush to ensure logs are sent before program exits
	log.Println("Flushing logs...")
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}

	// Give a moment for the data to be processed
	time.Sleep(2 * time.Second)

	log.Println("✅ Logs sent to UserCanal!")
}
