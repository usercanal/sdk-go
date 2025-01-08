// client/main.go
package main

import (
	"context"
	"log"
	"time"

	"github.com/usercanal/sdk-go/batch"
	"github.com/usercanal/sdk-go/convert"
	"github.com/usercanal/sdk-go/transport"
	"github.com/usercanal/sdk-go/types"
)

type Client struct {
	sender  *transport.Sender
	batcher *batch.Manager
}

func NewClient(endpoint string) (*Client, error) {
	sender, err := transport.NewSender("test-key", endpoint, 3)
	if err != nil {
		return nil, err
	}

	client := &Client{
		sender: sender,
	}

	// Create batcher with send function
	client.batcher = batch.NewManager(2, 5*time.Second, sender.Send)

	return client, nil
}

func (c *Client) Track(ctx context.Context, event types.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	protoEvent, err := convert.EventToProto(&event)
	if err != nil {
		return err
	}

	return c.batcher.Add(ctx, protoEvent)
}

func (c *Client) Close() error {
	if err := c.batcher.Flush(context.Background()); err != nil {
		return err
	}
	return c.sender.Close()
}

func main() {
	client, err := NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Send test events including some non-ASCII characters
	events := []struct {
		name       string
		properties types.Properties
	}{
		{
			name: "test_ascii",
			properties: types.Properties{
				"text": "Hello, World!",
			},
		},
		{
			name: "test_utf8",
			properties: types.Properties{
				"text": "Hello, 世界!", // Japanese
			},
		},
		{
			name: "test_emoji",
			properties: types.Properties{
				"text": "Hello! 👋 🌍", // Emojis
			},
		},
	}

	for i, e := range events {
		event := types.Event{
			UserId:     "test123",
			Name:       e.name,
			Timestamp:  time.Now(),
			Properties: e.properties,
		}

		if err := client.Track(context.Background(), event); err != nil {
			log.Fatalf("Failed to track event %d: %v", i, err)
		}

		log.Printf("Sent event %d: %s", i, e.name)
	}

	// Wait a bit to ensure events are sent
	time.Sleep(2 * time.Second)

	// Print metrics
	metrics := client.sender.GetMetrics()
	log.Printf("Metrics: Events sent: %d, Batches sent: %d, Failed attempts: %d",
		metrics.EventsSent, metrics.BatchesSent, metrics.FailedAttempts)
}
