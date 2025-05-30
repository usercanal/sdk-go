// sdk-go/types/events.go
package types

import "time"

// EventName represents a strongly typed event name
type EventName string

// Event Types
const (
	// User lifecycle events
	UserSignedUp EventName = "user_signed_up"
	UserLoggedIn EventName = "user_logged_in"
	FeatureUsed  EventName = "feature_used"

	// Revenue & Conversion Events
	OrderCompleted       EventName = "order_completed"
	SubscriptionStarted  EventName = "subscription_started"
	SubscriptionChanged  EventName = "subscription_changed"
	SubscriptionCanceled EventName = "subscription_canceled"
	CartViewed           EventName = "cart_viewed"
	CheckoutStarted      EventName = "checkout_started"
	CheckoutCompleted    EventName = "checkout_completed"
)

// Revenue type represents different types of revenue transactions
type RevenueType string

const (
	RevenueTypeSubscription RevenueType = "subscription"
	RevenueTypeOneTime      RevenueType = "one_time"
)

// AuthMethod represents authentication methods
type AuthMethod string

const (
	AuthMethodGoogle AuthMethod = "google"
	AuthMethodEmail  AuthMethod = "email"
)

// PaymentMethod represents payment methods
type PaymentMethod string

const (
	PaymentMethodCard PaymentMethod = "card"
)

// Currency represents supported currencies
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
)

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

// String methods
func (e EventName) String() string {
	return string(e)
}

func (r RevenueType) String() string {
	return string(r)
}

func (a AuthMethod) String() string {
	return string(a)
}

func (p PaymentMethod) String() string {
	return string(p)
}

func (c Currency) String() string {
	return string(c)
}

// IsStandardEvent checks if the event name is a known standard event
func (e EventName) IsStandardEvent() bool {
	switch e {
	case UserSignedUp, UserLoggedIn, FeatureUsed,
		OrderCompleted, SubscriptionStarted, SubscriptionChanged,
		SubscriptionCanceled, CartViewed, CheckoutStarted,
		CheckoutCompleted:
		return true
	}
	return false
}
