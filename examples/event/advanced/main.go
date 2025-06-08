// sdk-go/examples/events/advanced/main.go
package main

import (
	"context"
	"log"
	"time"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Initialize with advanced configuration for high-volume analytics
	client, err := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
		Endpoint:      "collect.usercanal.com:50000",
		BatchSize:     100,
		FlushInterval: 5 * time.Second,
		MaxRetries:    3,
		Debug:         true,
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

	// Track feature usage
	err = client.Event(ctx, "user123", usercanal.FeatureUsed, usercanal.Properties{
		"feature_name": "data_export",
		"export_type":  "csv",
		"file_size":    1024,
	})
	if err != nil {
		log.Printf("Failed to track feature usage: %v", err)
	}

	// Track revenue event
	err = client.EventRevenue(ctx, "user123", "order_456", 99.99, usercanal.CurrencyUSD, usercanal.Properties{
		"payment_method": usercanal.PaymentMethodCard,
		"discount_code":  "SAVE10",
		"type":           usercanal.RevenueTypeOneTime,
		"products": []usercanal.Product{
			{
				ID:       "prod_123",
				Name:     "Premium Plan",
				Price:    99.99,
				Quantity: 1,
			},
		},
	})
	if err != nil {
		log.Printf("Failed to track revenue: %v", err)
	}

	// Identify user with traits
	err = client.EventIdentify(ctx, "user123", usercanal.Properties{
		"name":        "John Doe",
		"email":       "john@example.com",
		"plan":        "premium",
		"signup_date": time.Now().AddDate(0, -1, 0),
		"auth_method": usercanal.AuthMethodEmail,
	})
	if err != nil {
		log.Printf("Failed to identify user: %v", err)
	}

	// Associate user with a group/organization
	err = client.EventGroup(ctx, "user123", "org_789", usercanal.Properties{
		"organization_name": "Acme Corp",
		"plan":              "enterprise",
		"seat_count":        50,
	})
	if err != nil {
		log.Printf("Failed to assign group: %v", err)
	}

	// Flush and wait to ensure all events are sent
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}

	// Print analytics stats
	stats := client.GetStats()
	log.Printf("Client Stats: Events queued: %d, Connection: %s", stats.EventsInQueue, stats.ConnectionState)
}
