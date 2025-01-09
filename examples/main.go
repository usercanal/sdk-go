// examples/main.go
package main

import (
	"context"
	"log"
	"time"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Initialize with struct config
	client, err := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
		Endpoint:  "127.0.0.1:50051",
		BatchSize: 100,
		Debug:     true,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Track a feature usage
	err = client.Track(context.Background(), usercanal.Event{
		UserId: "user123",
		Name:   usercanal.FeatureUsed,
		Properties: usercanal.Properties{
			"feature_name": "search",
			"duration_ms":  1500,
			"results":      42,
		},
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Printf("Failed to track event: %v", err)
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

	// Flush any pending events
	err = client.Flush(context.Background())
	if err != nil {
		log.Printf("Failed to flush events: %v", err)
	}

	// Wait a bit to ensure events are sent
	time.Sleep(2 * time.Second)

	// Debug, Print detailed status
	client.DumpStatus()
}
