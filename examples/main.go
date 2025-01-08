// examples/main.go
package main

import (
	"context"
	"log"
	"time"

	"github.com/usercanal/sdk-go/usercanal"
)

func main() {
	// Initialize with struct config
	canal, err := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
		Endpoint:  "127.0.0.1:50051",
		BatchSize: 100,
		Debug:     true,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer canal.Close()

	// Track a feature usage
	err = canal.Track(context.Background(), usercanal.Event{
		UserId: "user123",
		Name:   usercanal.FeatureUsed, // This is now an EventName type
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
	err = canal.Track(context.Background(), usercanal.Event{
		UserId: "user123",
		Name:   usercanal.OrderCompleted, // This is now an EventName type
		Properties: usercanal.Properties{
			"revenue":        99.99,
			"currency":       usercanal.CurrencyUSD, // Using Currency type
			"product_id":     "prod_123",
			"quantity":       1,
			"payment_method": usercanal.PaymentMethodCard,  // Using PaymentMethod type
			"type":           usercanal.RevenueTypeOneTime, // Using RevenueType type
		},
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Printf("Failed to track revenue: %v", err)
	}

	// Identify user
	err = canal.Identify(context.Background(), usercanal.Identity{
		UserId: "user123",
		Properties: usercanal.Properties{
			"name":        "John Doe",
			"email":       "john@example.com",
			"auth_method": usercanal.AuthMethodEmail, // Using AuthMethod type
		},
	})
	if err != nil {
		log.Printf("Failed to identify user: %v", err)
	}

	canal.Flush(context.Background())

	// Wait a bit to ensure events are sent
	time.Sleep(2 * time.Second)

	// Debug, Print detailed status (includes connection state, event queues, etc.)
	canal.DumpStatus()
}
