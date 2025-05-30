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

	ctx := context.Background()

	// Track feature usage
	err = client.TrackEvent(ctx, usercanal.Event{
		UserId: "user123",
		Name:   usercanal.FeatureUsed,
		Properties: usercanal.Properties{
			"feature_name": "data_export",
			"export_type":  "csv",
			"file_size":    1024,
		},
	})
	if err != nil {
		log.Printf("Failed to track feature usage: %v", err)
	}

	// Track revenue event
	//	canal.Revenue(ctx, usercanal.Revenue{
	//		UserID:   "user_123", // The actual user
	//		OrderID:  "ord_456",  // Order identifier (goes in payload)
	//		Amount:   99.99,
	//		Currency: usercanal.CurrencyUSD,
	//		Type:     usercanal.RevenueTypeSubscription,
	//		Products: []usercanal.Product{
	//			{
	//				ID:       "prod_123",
	//				Name:     "Pro Plan",
	//				Price:    99.99,
	//				Quantity: 1,
	//			},
	//		},
	//	})

	err = client.Track(ctx, usercanal.Event{
		OrderID:  "order_456",
		Amount:   99.99,
		Currency: usercanal.CurrencyUSD,
		Type:     usercanal.RevenueTypeOneTime,
		Products: []usercanal.Product{
			{
				ID:       "prod_123",
				Name:     "Premium Plan",
				Price:    99.99,
				Quantity: 1,
			},
		},
		Properties: usercanal.Properties{
			"payment_method": usercanal.PaymentMethodCard,
			"discount_code":  "SAVE10",
		},
	})
	if err != nil {
		log.Printf("Failed to track revenue: %v", err)
	}

	// Identify user with traits
	err = client.IdentifyUser(ctx, usercanal.Identity{
		UserId: "user123",
		Properties: usercanal.Properties{
			"name":        "John Doe",
			"email":       "john@example.com",
			"plan":        "premium",
			"signup_date": time.Now().AddDate(0, -1, 0),
			"auth_method": usercanal.AuthMethodEmail,
		},
	})
	if err != nil {
		log.Printf("Failed to identify user: %v", err)
	}

	// Associate user with a group/organization
	err = client.AssignGroup(ctx, usercanal.GroupInfo{
		UserId:  "user123",
		GroupId: "org_789",
		Properties: usercanal.Properties{
			"organization_name": "Acme Corp",
			"plan":              "enterprise",
			"seat_count":        50,
		},
	})
	if err != nil {
		log.Printf("Failed to assign group: %v", err)
	}

	// Flush and wait to ensure all events are sent
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}

	// Print analytics stats
	client.DumpStatus()
}
