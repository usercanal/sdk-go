package main

import (
	"context"
	"log"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Initialize client with minimal configuration
	client, err := usercanal.NewClient("000102030405060708090a0b0c0d0e0f")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Track a signup event
	err = client.Track(context.Background(), usercanal.Event{
		UserId: "user_123",
		Name:   usercanal.UserSignedUp,
		Properties: usercanal.Properties{
			"signup_method":   "email",
			"referral_source": "google",
		},
	})
	if err != nil {
		log.Printf("Failed to track signup: %v", err)
	}

	// Ensure event is sent before program exits
	if err := client.Flush(context.Background()); err != nil {
		log.Printf("Failed to flush: %v", err)
	}
}
