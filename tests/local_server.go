package main

import (
	"context"
	"log"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Configure client for local server
	config := usercanal.Config{
		Endpoint: "localhost:50000",
		Debug:    true,
	}

	client, err := usercanal.NewClient("000102030405060708090a0b0c0d0e0f", config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close(context.Background())

	ctx := context.Background()

	// Send a simple test event
	err = client.Event(ctx, "go_sdk_test_user", usercanal.UserSignedUp, usercanal.Properties{
		"test": true,
		"sdk":  "go",
	})
	if err != nil {
		log.Printf("Failed to send event: %v", err)
		return
	}

	// Flush to ensure event is sent
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
		return
	}

	log.Println("✅ Go SDK test event sent successfully!")
	log.Println("💡 Check collector logs for: user_id='go_sdk_test_user'")
}
