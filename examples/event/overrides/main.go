package main

import (
	"context"
	"log"
	"time"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Initialize client
	client, err := usercanal.NewClient("000102030405060708090a0b0c0d0e0f", usercanal.Config{
		Endpoint:      "localhost:50001",
		Debug:         true,
		FlushInterval: 1 * time.Second,
		BatchSize:     1,
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

	// Server-side scenario 1: Provide device_id (REQUIRED for server SDKs)
	// Server must explicitly provide device_id - no auto-generation
	userDeviceID := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}
	err = client.EventAdvanced(ctx, usercanal.EventAdvanced{
		UserId:   "user_123",
		Name:     usercanal.SubscriptionStarted,
		DeviceID: &userDeviceID, // MUST be provided - no auto-generation for servers
		Properties: usercanal.Properties{
			"plan":   "premium",
			"amount": 29.99,
		},
	})
	if err != nil {
		log.Printf("Failed to track subscription start: %v", err)
	}

	// Server-side scenario 2: No session ID (server events without sessions)
	err = client.EventAdvanced(ctx, usercanal.EventAdvanced{
		UserId:    "user_456",
		Name:      usercanal.FeatureUsed,
		DeviceID:  &userDeviceID,
		SessionID: nil, // No session for server-side events
		Properties: usercanal.Properties{
			"feature_name": "api_data_export",
			"endpoint":     "/api/data/export",
			"method":       "POST",
		},
	})
	if err != nil {
		log.Printf("Failed to track feature usage: %v", err)
	}

	// Server-side scenario 3: Proxy scenario (forwarding client context)
	clientSessionID := []byte{0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
	customTime := time.Now().Add(-5 * time.Minute)
	err = client.EventAdvanced(ctx, usercanal.EventAdvanced{
		UserId:    "user_789",
		Name:      usercanal.CheckoutStarted,
		DeviceID:  &userDeviceID,
		SessionID: &clientSessionID, // Forwarded from client SDK
		Timestamp: &customTime,      // Custom timestamp
		Properties: usercanal.Properties{
			"page":   "/dashboard",
			"source": "mobile_app",
			"via":    "server_proxy",
		},
	})
	if err != nil {
		log.Printf("Failed to track checkout: %v", err)
	}

	// Regular event - will FAIL because device_id is required but not provided
	// Server SDKs don't auto-generate device_id - must use EventAdvanced
	err = client.Event(ctx, "user_123", usercanal.FeatureUsed, usercanal.Properties{
		"feature": "data_export",
	})
	if err != nil {
		log.Printf("Expected failure - regular events need device_id: %v", err)
	}

	// Flush and close
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}

	log.Println("Server-side event tracking examples completed")
}
