// usercanal.go
package usercanal

import (
	"context"
	"time"

	"github.com/usercanal/sdk-go/internal/api"
	"github.com/usercanal/sdk-go/internal/version"
	"github.com/usercanal/sdk-go/types"
)

// Config holds client configuration
type Config struct {
	Endpoint      string        // API Endpoint
	BatchSize     int           // Events per batch
	FlushInterval time.Duration // Max time between sends
	MaxRetries    int           // Retry attempts
	Debug         bool          // Enable debug logging
}

// Client is a facade over the internal API client
type Client struct {
	internal *api.Client
}

// NewClient creates a new client with configuration
func NewClient(apiKey string, cfg ...Config) (*Client, error) {
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

	client, err := api.New(apiKey, options...)
	if err != nil {
		return nil, err
	}

	return &Client{internal: client}, nil
}

// Re-export main client methods
func (c *Client) Track(ctx context.Context, event Event) error {
	return c.internal.Track(ctx, event)
}

func (c *Client) Identify(ctx context.Context, identity Identity) error {
	return c.internal.Identify(ctx, identity)
}

func (c *Client) Group(ctx context.Context, group GroupInfo) error {
	return c.internal.Group(ctx, group)
}

func (c *Client) Revenue(ctx context.Context, rev Revenue) error {
	return c.internal.Revenue(ctx, rev)
}

func (c *Client) Flush(ctx context.Context) error {
	return c.internal.Flush(ctx)
}

func (c *Client) Close() error {
	return c.internal.Close()
}

func (c *Client) DumpStatus() {
	c.internal.DumpStatus()
}

// Re-export types that users need
type (
	Properties = types.Properties
	Event      = types.Event
	Identity   = types.Identity
	GroupInfo  = types.GroupInfo
	Revenue    = types.Revenue
	Product    = types.Product
	Currency   = types.Currency
)

// Re-export constants
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
