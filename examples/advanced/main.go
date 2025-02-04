// examples/main.go
package main

import (
	"context"
	"log"
	"time"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Initialize with proper endpoint
	client, err := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
		Endpoint:      "collect.usercanal.com:9000",
		BatchSize:     100,
		FlushInterval: 5 * time.Second,
		MaxRetries:    3,
		Debug:         true,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Add a small delay to ensure connection is established
	time.Sleep(2 * time.Second)

	// Test connection with a simple event
	ctx := context.Background()
	err = client.Track(ctx, usercanal.Event{
		UserId: "test_user",
		Name:   usercanal.FeatureUsed,
		Properties: usercanal.Properties{
			"test": true,
		},
	})
	if err != nil {
		log.Fatalf("Failed to send test event: %v", err)
	}

	// Track revenue
	err = client.Track(context.Background(), usercanal.Event{
		UserId: "user123",
		Name:   usercanal.OrderCompleted,
		Properties: usercanal.Properties{
			"revenue":        99.99,
			"currency":       usercanal.CurrencyUSD,
			"product_id":     "prod_123",
			"quantity":       1,
			"payment_method": usercanal.PaymentMethodCard,
			"type":           usercanal.RevenueTypeOneTime,
		},
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Printf("Failed to track revenue: %v", err)
	}

	// Identify user
	err = client.Identify(context.Background(), usercanal.Identity{
		UserId: "user123",
		Properties: usercanal.Properties{
			"name":        "John Doe",
			"email":       "john@example.com",
			"auth_method": usercanal.AuthMethodEmail,
		},
	})
	if err != nil {
		log.Printf("Failed to identify user: %v", err)
	}

	// Flush and wait to ensure events are sent
	if err := client.Flush(context.Background()); err != nil {
		log.Printf("Failed to flush: %v", err)
	}

	// Print stats
	client.DumpStatus()
}
