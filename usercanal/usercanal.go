// usercanal/usercanal.go
package usercanal

import (
	"time"

	"github.com/usercanal/sdk-go/api"
	"github.com/usercanal/sdk-go/types"
	"github.com/usercanal/sdk-go/version"
)

// Config holds client configuration
type Config struct {
	Endpoint      string        // API Endpoint
	BatchSize     int           // Events per batch
	FlushInterval time.Duration // Max time between sends
	MaxRetries    int           // Retry attempts
	Debug         bool          // Enable debug logging
}

// NewClient creates a new client with configuration
func NewClient(apiKey string, cfg ...Config) (*api.Client, error) {
	var options []api.Option

	if len(cfg) > 0 {
		c := cfg[0]
		options = append(options,
			api.WithEndpoint(c.Endpoint),
			api.WithBatchSize(c.BatchSize),
			api.WithFlushInterval(c.FlushInterval),
			api.WithMaxRetries(c.MaxRetries),
			api.WithDebug(c.Debug),
		)
	}

	return api.New(apiKey, options...)
}

// Re-export only the types from types package that users need
type (
	Properties = types.Properties
	Event      = types.Event
	Identity   = types.Identity
	GroupInfo  = types.GroupInfo
	Revenue    = types.Revenue
	Product    = types.Product
	Currency   = types.Currency
)

// Re-export constants from types package
const (
	// User lifecycle events
	UserSignedUp = types.UserSignedUp
	UserLoggedIn = types.UserLoggedIn
	FeatureUsed  = types.FeatureUsed

	// Revenue & Conversion Events
	OrderCompleted       = types.OrderCompleted
	SubscriptionStarted  = types.SubscriptionStarted
	SubscriptionChanged  = types.SubscriptionChanged
	SubscriptionCanceled = types.SubscriptionCanceled
	CartViewed           = types.CartViewed
	CheckoutStarted      = types.CheckoutStarted
	CheckoutCompleted    = types.CheckoutCompleted

	// Currency codes
	CurrencyUSD = types.CurrencyUSD
	CurrencyEUR = types.CurrencyEUR
	CurrencyGBP = types.CurrencyGBP

	// Revenue types
	RevenueTypeSubscription = types.RevenueTypeSubscription
	RevenueTypeOneTime      = types.RevenueTypeOneTime

	// Auth methods
	AuthMethodGoogle = types.AuthMethodGoogle
	AuthMethodEmail  = types.AuthMethodEmail

	// Payment methods
	PaymentMethodCard = types.PaymentMethodCard
)

// Version returns detailed version information
func Version() version.Info {
	return version.Get()
}
