// sdk-go/internal/convert/event.go
package convert

import (
	"fmt"
	"time"

	event_collector "github.com/usercanal/sdk-go/internal/schema/event"
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
	if err := validateRequired("UserId", e.UserId); err != nil {
		return nil, err
	}

	if err := validateRequired("Name", string(e.Name)); err != nil {
		return nil, err
	}

	// Validate event type mapping
	eventType, ok := eventTypeMap[e.Name]
	if !ok {
		return nil, fmt.Errorf("unmapped event type: %s", e.Name)
	}

	payload, err := marshalPayload(map[string]interface{}{
		"name":       e.Name.String(),
		"properties": e.Properties,
	})
	if err != nil {
		return nil, err
	}

	return &transport.Event{
		Timestamp: resolveTimestamp(e.Timestamp),
		EventType: eventType,
		UserID:    []byte(e.UserId),
		Payload:   payload,
	}, nil
}

// IdentityToInternal converts a types.Identity to an internal transport.Event
func IdentityToInternal(i *types.Identity) (*transport.Event, error) {
	if err := validateRequired("UserId", i.UserId); err != nil {
		return nil, err
	}

	payload, err := marshalPayload(map[string]interface{}{
		"traits": i.Properties,
	})
	if err != nil {
		return nil, err
	}

	return &transport.Event{
		Timestamp: resolveTimestamp(time.Time{}), // Always use current time
		EventType: event_collector.EventTypeIDENTIFY,
		UserID:    []byte(i.UserId),
		Payload:   payload,
	}, nil
}

// GroupToInternal converts a types.GroupInfo to an internal transport.Event
func GroupToInternal(g *types.GroupInfo) (*transport.Event, error) {
	if err := validateRequired("UserId", g.UserId); err != nil {
		return nil, err
	}

	if err := validateRequired("GroupId", g.GroupId); err != nil {
		return nil, err
	}

	payload, err := marshalPayload(map[string]interface{}{
		"group_id":   g.GroupId,
		"properties": g.Properties,
	})
	if err != nil {
		return nil, err
	}

	return &transport.Event{
		Timestamp: resolveTimestamp(time.Time{}), // Always use current time
		EventType: event_collector.EventTypeGROUP,
		UserID:    []byte(g.UserId),
		Payload:   payload,
	}, nil
}

func RevenueToInternal(r *types.Revenue) (*transport.Event, error) {
	if err := validateRequired("UserID", r.UserID); err != nil {
		return nil, err
	}

	if err := validateRequired("OrderID", r.OrderID); err != nil {
		return nil, err
	}

	if r.Amount <= 0 {
		return nil, types.NewValidationError("Amount", "must be positive")
	}

	if err := validateRequired("Currency", string(r.Currency)); err != nil {
		return nil, err
	}

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
		"order_id": r.OrderID, // OrderID goes in payload where it belongs
		"revenue":  r.Amount,
		"currency": r.Currency,
		"type":     r.Type,
	}

	if products != nil {
		properties["products"] = products
	}

	// Merge custom properties
	for k, v := range r.Properties {
		properties[k] = v
	}

	payload, err := marshalPayload(map[string]interface{}{
		"name":       types.OrderCompleted.String(),
		"properties": properties,
	})
	if err != nil {
		return nil, err
	}

	return &transport.Event{
		Timestamp: resolveTimestamp(time.Time{}),
		EventType: event_collector.EventTypeTRACK,
		UserID:    []byte(r.UserID), // Correct: actual user who made the purchase
		Payload:   payload,          // OrderID is in the payload data
	}, nil
}
