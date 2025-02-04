// convert/convert.go
package convert

import (
	"encoding/json"
	"fmt"
	"time"

	event_collector "github.com/usercanal/sdk-go/internal/event"
	"github.com/usercanal/sdk-go/internal/transport"
	"github.com/usercanal/sdk-go/types"
)

// Map SDK event names to FlatBuffer event types
var eventTypeMap = map[types.EventName]event_collector.EventType{
	types.UserSignedUp:         event_collector.EventTypeTRACK,
	types.UserLoggedIn:         event_collector.EventTypeTRACK,
	types.FeatureUsed:          event_collector.EventTypeTRACK,
	types.OrderCompleted:       event_collector.EventTypeTRACK,
	types.SubscriptionStarted:  event_collector.EventTypeTRACK,
	types.SubscriptionChanged:  event_collector.EventTypeTRACK,
	types.SubscriptionCanceled: event_collector.EventTypeTRACK,
	types.CartViewed:           event_collector.EventTypeTRACK,
	types.CheckoutStarted:      event_collector.EventTypeTRACK,
	types.CheckoutCompleted:    event_collector.EventTypeTRACK,
}

// EventToInternal converts a types.Event to an internal transport.Event
func EventToInternal(e *types.Event) (*transport.Event, error) {
	// Validate required fields
	if e.UserId == "" {
		return nil, types.NewValidationError("UserId", "is required")
	}

	// Always use current time if timestamp is not explicitly set
	timestamp := time.Now()
	if !e.Timestamp.IsZero() {
		timestamp = e.Timestamp
	}

	payload, err := json.Marshal(map[string]interface{}{
		"name":       e.Name.String(),
		"properties": e.Properties,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return &transport.Event{
		Timestamp: uint64(timestamp.UnixMilli()),
		EventType: eventTypeMap[e.Name],
		UserID:    []byte(e.UserId),
		Payload:   payload,
	}, nil
}

// IdentityToInternal converts a types.Identity to an internal transport.Event
func IdentityToInternal(i *types.Identity) (*transport.Event, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"traits": i.Properties,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return &transport.Event{
		Timestamp: uint64(time.Now().UnixMilli()),
		EventType: event_collector.EventTypeIDENTIFY,
		UserID:    []byte(i.UserId),
		Payload:   payload,
	}, nil
}

// GroupToInternal converts a types.GroupInfo to an internal transport.Event
func GroupToInternal(g *types.GroupInfo) (*transport.Event, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"group_id":   g.GroupId,
		"properties": g.Properties,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return &transport.Event{
		Timestamp: uint64(time.Now().UnixMilli()),
		EventType: event_collector.EventTypeGROUP,
		UserID:    []byte(g.UserId),
		Payload:   payload,
	}, nil
}

// RevenueToInternal converts a types.Revenue to an internal transport.Event
func RevenueToInternal(r *types.Revenue) (*transport.Event, error) {
	var products []map[string]interface{}
	if len(r.Products) > 0 {
		products = make([]map[string]interface{}, len(r.Products))
		for i, p := range r.Products {
			products[i] = map[string]interface{}{
				"id":       p.ID,
				"name":     p.Name,
				"price":    p.Price,
				"quantity": p.Quantity,
			}
		}
	}

	properties := map[string]interface{}{
		"order_id": r.OrderID,
		"revenue":  r.Amount,
		"currency": r.Currency,
		"type":     r.Type,
	}

	if products != nil {
		properties["products"] = products
	}

	for k, v := range r.Properties {
		properties[k] = v
	}

	payload, err := json.Marshal(map[string]interface{}{
		"name":       types.OrderCompleted.String(),
		"properties": properties,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return &transport.Event{
		Timestamp: uint64(time.Now().UnixMilli()),
		EventType: event_collector.EventTypeTRACK,
		UserID:    []byte(r.OrderID),
		Payload:   payload,
	}, nil
}
