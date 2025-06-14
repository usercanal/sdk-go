// sdk-go/types/events.go
package types

import "time"

// Event represents a tracking event
type Event struct {
	ID         string
	UserId     string
	Name       EventName
	Properties Properties
	Timestamp  time.Time `json:"timestamp,omitempty"`
}

// Identity represents a user identification event
type Identity struct {
	UserId     string
	Properties Properties
}

// GroupInfo represents a group event
type GroupInfo struct {
	UserId     string
	GroupId    string
	Properties Properties
}

// Revenue represents a revenue event
type Revenue struct {
	UserID     string
	OrderID    string
	Amount     float64
	Currency   Currency
	Type       RevenueType
	Products   []Product
	Properties Properties
}

// Product represents a product in a revenue event
type Product struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
}


