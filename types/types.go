// types/types.go

package types

import (
	"time"
)

// Properties represents a map of property values
type Properties map[string]interface{}

// Event represents a tracking event
type Event struct {
	ID         string
	UserId     string
	Name       EventName // Using EventName type instead of string
	Properties Properties
	Timestamp  time.Time
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
	OrderID    string
	Amount     float64
	Currency   Currency
	Type       RevenueType // Using RevenueType instead of string
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
