// sdk-go/examples/string_flexibility.go
package main

import (
	"context"
	"fmt"
	"log"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Initialize client
	client, err := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
		Debug: true,
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

	fmt.Println("=== String Flexibility Demo ===")

	// 1. Using predefined constants (recommended for common events)
	fmt.Println("\n1. Using predefined constants:")
	
	err = client.Event(ctx, "user_123", usercanal.UserSignedUp, usercanal.Properties{
		"auth_method": usercanal.AuthMethodGoogle,
		"source":      "landing_page",
	})
	if err != nil {
		log.Printf("Failed to track signup: %v", err)
	}
	fmt.Printf("Tracked: %s\n", usercanal.UserSignedUp)

	err = client.Event(ctx, "user_123", usercanal.FeatureUsed, usercanal.Properties{
		"feature_name": "export",
		"duration_ms":  1500,
	})
	if err != nil {
		log.Printf("Failed to track feature: %v", err)
	}
	fmt.Printf("Tracked: %s\n", usercanal.FeatureUsed)

	// 2. Using custom strings (full flexibility for domain-specific events)
	fmt.Println("\n2. Using custom strings:")
	
	err = client.Event(ctx, "user_123", "video.viewed", usercanal.Properties{
		"video_id":   "vid_123",
		"duration":   120,
		"completion": 0.85,
	})
	if err != nil {
		log.Printf("Failed to track video: %v", err)
	}
	fmt.Println("Tracked: video.viewed")

	err = client.Event(ctx, "user_123", "report.generated", usercanal.Properties{
		"report_type": "monthly_summary",
		"format":      "pdf",
		"size_mb":     2.3,
	})
	if err != nil {
		log.Printf("Failed to track report: %v", err)
	}
	fmt.Println("Tracked: report.generated")

	// 3. Mixed usage in real scenarios
	fmt.Println("\n3. Mixed usage (realistic scenario):")
	
	// Standard e-commerce events use constants
	err = client.Event(ctx, "user_123", usercanal.CartViewed, usercanal.Properties{
		"cart_value": 149.99,
		"item_count": 3,
	})
	if err != nil {
		log.Printf("Failed to track cart: %v", err)
	}

	// Domain-specific events use custom strings
	err = client.Event(ctx, "user_123", "ai.prompt.submitted", usercanal.Properties{
		"prompt_length": 45,
		"model":         "gpt-4",
		"tokens_used":   120,
	})
	if err != nil {
		log.Printf("Failed to track AI prompt: %v", err)
	}

	// 4. Demonstrate type conversion
	fmt.Println("\n4. EventName type flexibility:")
	
	// EventName is just a string type, so you can convert freely
	customEventName := usercanal.EventName("custom.workflow.completed")
	
	err = client.Event(ctx, "user_123", customEventName, usercanal.Properties{
		"workflow_id": "wf_789",
		"steps":       5,
		"duration":    "2m30s",
	})
	if err != nil {
		log.Printf("Failed to track workflow: %v", err)
	}
	fmt.Printf("Tracked: %s\n", customEventName)

	// 5. Check if event is standard or custom
	fmt.Println("\n5. Standard vs Custom detection:")
	
	events := []usercanal.EventName{
		usercanal.UserSignedUp,     // Standard
		"video.viewed",             // Custom
		usercanal.OrderCompleted,   // Standard
		"ai.prompt.submitted",      // Custom
	}
	
	for _, event := range events {
		isStandard := event.IsStandardEvent()
		fmt.Printf("%s -> %s\n", event, map[bool]string{true: "Standard", false: "Custom"}[isStandard])
	}

	// Flush all events
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}

	fmt.Println("\n✅ All events tracked successfully!")
	fmt.Println("\nKey Benefits:")
	fmt.Println("• Predefined constants provide consistency and IDE autocomplete")
	fmt.Println("• Custom strings give full flexibility for domain-specific events")
	fmt.Println("• Same API works for both - no need to choose upfront")
	fmt.Println("• Dashboard shows human-readable names (e.g., 'User Signed Up')")
}