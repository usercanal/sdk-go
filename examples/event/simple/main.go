package main

import (
	"context"
	"log"
	"time"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Initialize client with minimal configuration
	client, err := usercanal.NewClient("000102030405060708090a0b0c0d0e0f", usercanal.Config{
		Endpoint:      "localhost:50001",
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

	// Track a signup event using predefined constant
	err = client.Event(ctx, "user_123", usercanal.UserSignedUp, usercanal.Properties{
		"signup_method":   "email",
		"referral_source": "google",
	})
	if err != nil {
		log.Printf("Failed to track signup: %v", err)
	}

	// Track a custom event using string directly
	err = client.Event(ctx, "user_123", "video.viewed", usercanal.Properties{
		"video_id": "vid_123",
		"duration": 120,
		"quality":  "hd",
		"platform": "web",
	})
	if err != nil {
		log.Printf("Failed to track video view: %v", err)
	}

	// Track another predefined event
	err = client.Event(ctx, "user_123", usercanal.FeatureUsed, usercanal.Properties{
		"feature_name": "dashboard",
		"section":      "analytics",
	})
	if err != nil {
		log.Printf("Failed to track feature usage: %v", err)
	}

	// Ensure events are sent before program exits
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}

	// Close the client
	if err := client.Close(ctx); err != nil {
		log.Printf("Failed to close client: %v", err)
	}
}
