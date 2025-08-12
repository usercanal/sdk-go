// sdk-go/types/events.go
package types

import "time"

// Event represents a tracking event
type Event struct {
	ID         string
	UserId     string
	SessionID  []byte // Optional session ID override (16-byte binary)
	Name       EventName
	Properties Properties
	Timestamp  time.Time `json:"timestamp,omitempty"`
}

// Identity represents a user identification event
type Identity struct {
	UserId     string
	SessionID  []byte // Optional session ID override (16-byte binary)
	Properties Properties
}

// GroupInfo represents a group event
type GroupInfo struct {
	UserId     string
	GroupId    string
	SessionID  []byte // Optional session ID override (16-byte binary)
	Properties Properties
}

// Revenue represents a revenue event
type Revenue struct {
	UserID     string
	OrderID    string
	SessionID  []byte // Optional session ID override (16-byte binary)
	Amount     float64
	Currency   Currency
	Type       RevenueType
	Products   []Product
	Properties Properties
}

// EventAdvanced represents an advanced tracking event with optional overrides
type EventAdvanced struct {
	UserId     string     // Required - user identifier
	Name       EventName  // Required - event name
	Properties Properties // Optional - event properties

	// Advanced optional overrides
	DeviceID  *[]byte    // Optional - override device_id (16-byte UUID)
	SessionID *[]byte    // Optional - override session_id (16-byte UUID)
	Timestamp *time.Time // Optional - custom timestamp
}

// Product represents a product in a revenue event
type Product struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
}
