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

func (c *Client) GetStats() Stats {
	return c.internal.GetStats()
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

// Event protocol
// Simplified parameter approach for better developer experience
func (c *Client) Event(ctx context.Context, userID string, eventName EventName, properties Properties) error {
	event := Event{
		UserId:     userID,
		Name:       eventName,
		Properties: properties,
	}
	return c.internal.Track(ctx, event)
}

func (c *Client) EventIdentify(ctx context.Context, userID string, traits Properties) error {
	identity := Identity{
		UserId:     userID,
		Properties: traits,
	}
	return c.internal.Identify(ctx, identity)
}

func (c *Client) EventGroup(ctx context.Context, userID string, groupID string, properties Properties) error {
	group := GroupInfo{
		UserId:     userID,
		GroupId:    groupID,
		Properties: properties,
	}
	return c.internal.Group(ctx, group)
}

func (c *Client) EventRevenue(ctx context.Context, userID string, orderID string, amount float64, currency Currency, properties Properties) error {
	revenue := Revenue{
		UserID:     userID,
		OrderID:    orderID,
		Amount:     amount,
		Currency:   currency,
		Properties: properties,
	}
	return c.internal.Revenue(ctx, revenue)
}

// TODO: EventAdvanced for complex cases (custom timestamps, event IDs, etc.)
// canal.EventAdvanced(ctx, Event{...}) - implement when customers need it

func (c *Client) Flush(ctx context.Context) error {
	return c.internal.Flush(ctx)
}

func (c *Client) Close(ctx context.Context) error {
	return c.internal.Close(ctx)
}



// Re-export types that users need
type (
	Properties           = types.Properties
	Event                = types.Event
	Identity             = types.Identity
	GroupInfo            = types.GroupInfo
	Revenue              = types.Revenue
	Product              = types.Product
	Currency             = types.Currency
	Stats                = types.Stats
	AuthMethod           = types.AuthMethod
	PaymentMethod        = types.PaymentMethod
	RevenueType          = types.RevenueType
	Channel              = types.Channel
	Source               = types.Source
	DeviceType           = types.DeviceType
	OperatingSystem      = types.OperatingSystem
	Browser              = types.Browser
	EventName            = types.EventName
	SubscriptionInterval = types.SubscriptionInterval
	PlanType             = types.PlanType
	UserRole             = types.UserRole
	CompanySize          = types.CompanySize
	Industry             = types.Industry
)

// Re-export constants
const (
	// Authentication & User Management Events
	UserSignedUp         = types.UserSignedUp
	UserSignedIn         = types.UserSignedIn
	UserSignedOut        = types.UserSignedOut
	UserInvited          = types.UserInvited
	UserOnboarded        = types.UserOnboarded
	AuthenticationFailed = types.AuthenticationFailed
	PasswordReset        = types.PasswordReset
	TwoFactorEnabled     = types.TwoFactorEnabled
	TwoFactorDisabled    = types.TwoFactorDisabled

	// Revenue & Billing Events
	OrderCompleted        = types.OrderCompleted
	OrderRefunded         = types.OrderRefunded
	OrderCanceled         = types.OrderCanceled
	PaymentFailed         = types.PaymentFailed
	PaymentMethodAdded    = types.PaymentMethodAdded
	PaymentMethodUpdated  = types.PaymentMethodUpdated
	PaymentMethodRemoved  = types.PaymentMethodRemoved

	// Subscription Management Events
	SubscriptionStarted  = types.SubscriptionStarted
	SubscriptionRenewed  = types.SubscriptionRenewed
	SubscriptionPaused   = types.SubscriptionPaused
	SubscriptionResumed  = types.SubscriptionResumed
	SubscriptionChanged  = types.SubscriptionChanged
	SubscriptionCanceled = types.SubscriptionCanceled

	// Trial & Conversion Events
	TrialStarted     = types.TrialStarted
	TrialEndingSoon  = types.TrialEndingSoon
	TrialEnded       = types.TrialEnded
	TrialConverted   = types.TrialConverted

	// Shopping Experience Events
	CartViewed        = types.CartViewed
	CartUpdated       = types.CartUpdated
	CartAbandoned     = types.CartAbandoned
	CheckoutStarted   = types.CheckoutStarted
	CheckoutCompleted = types.CheckoutCompleted

	// Product Engagement Events
	PageViewed          = types.PageViewed
	FeatureUsed         = types.FeatureUsed
	SearchPerformed     = types.SearchPerformed
	FileUploaded        = types.FileUploaded
	NotificationSent    = types.NotificationSent
	NotificationClicked = types.NotificationClicked

	// Communication Events
	EmailSent               = types.EmailSent
	EmailOpened             = types.EmailOpened
	EmailClicked            = types.EmailClicked
	EmailBounced            = types.EmailBounced
	EmailUnsubscribed       = types.EmailUnsubscribed
	SupportTicketCreated    = types.SupportTicketCreated
	SupportTicketResolved   = types.SupportTicketResolved

	// Authentication Methods
	AuthMethodPassword = types.AuthMethodPassword
	AuthMethodGoogle   = types.AuthMethodGoogle
	AuthMethodGitHub   = types.AuthMethodGitHub
	AuthMethodSSO      = types.AuthMethodSSO
	AuthMethodEmail    = types.AuthMethodEmail

	// Revenue Types
	RevenueTypeOneTime     = types.RevenueTypeOneTime
	RevenueTypeSubscription = types.RevenueTypeSubscription

	// Major Global Currencies
	CurrencyUSD  = types.CurrencyUSD
	CurrencyEUR  = types.CurrencyEUR
	CurrencyGBP  = types.CurrencyGBP
	CurrencyJPY  = types.CurrencyJPY
	CurrencyCAD  = types.CurrencyCAD
	CurrencyAUD  = types.CurrencyAUD
	CurrencyNZD  = types.CurrencyNZD
	CurrencyKRW  = types.CurrencyKRW
	CurrencyCNY  = types.CurrencyCNY
	CurrencyHKD  = types.CurrencyHKD
	CurrencySGD  = types.CurrencySGD
	CurrencyMXN  = types.CurrencyMXN
	CurrencyINR  = types.CurrencyINR
	CurrencyPLN  = types.CurrencyPLN
	CurrencyBRL  = types.CurrencyBRL
	CurrencyRUB  = types.CurrencyRUB
	CurrencyDKK  = types.CurrencyDKK
	CurrencyNOK  = types.CurrencyNOK
	CurrencySEK  = types.CurrencySEK
	CurrencyCHF  = types.CurrencyCHF
	CurrencyTRY  = types.CurrencyTRY
	CurrencyILS  = types.CurrencyILS
	CurrencyTHB  = types.CurrencyTHB
	CurrencyMYR  = types.CurrencyMYR
	CurrencyIDR  = types.CurrencyIDR
	CurrencyVND  = types.CurrencyVND
	CurrencyPHP  = types.CurrencyPHP
	CurrencyCZK  = types.CurrencyCZK
	CurrencyHUF  = types.CurrencyHUF
	CurrencyZAR  = types.CurrencyZAR
	CurrencyARS  = types.CurrencyARS
	CurrencyCLP  = types.CurrencyCLP
	CurrencyCOP  = types.CurrencyCOP
	CurrencyPEN  = types.CurrencyPEN
	CurrencyUYU  = types.CurrencyUYU
	CurrencyEGP  = types.CurrencyEGP
	CurrencyAED  = types.CurrencyAED
	CurrencySAR  = types.CurrencySAR
	CurrencyQAR  = types.CurrencyQAR
	CurrencyBHD  = types.CurrencyBHD
	CurrencyKWD  = types.CurrencyKWD
	CurrencyOMR  = types.CurrencyOMR
	CurrencyJOD  = types.CurrencyJOD
	CurrencyLBP  = types.CurrencyLBP
	CurrencyRON  = types.CurrencyRON
	CurrencyBGN  = types.CurrencyBGN
	CurrencyHRK  = types.CurrencyHRK
	CurrencyRSD  = types.CurrencyRSD
	CurrencyBAM  = types.CurrencyBAM
	CurrencyMKD  = types.CurrencyMKD
	CurrencyALL  = types.CurrencyALL
	CurrencyUAH  = types.CurrencyUAH
	CurrencyBYN  = types.CurrencyBYN
	CurrencyMDL  = types.CurrencyMDL
	CurrencyGEL  = types.CurrencyGEL
	CurrencyAMD  = types.CurrencyAMD
	CurrencyAZN  = types.CurrencyAZN
	CurrencyKZT  = types.CurrencyKZT
	CurrencyUZS  = types.CurrencyUZS
	CurrencyKGS  = types.CurrencyKGS
	CurrencyTJS  = types.CurrencyTJS
	CurrencyTMT  = types.CurrencyTMT
	CurrencyMNT  = types.CurrencyMNT
	CurrencyBTC  = types.CurrencyBTC
	CurrencyETH  = types.CurrencyETH
	CurrencyUSDC = types.CurrencyUSDC
	CurrencyUSDT = types.CurrencyUSDT

	// Payment Methods
	PaymentMethodCard         = types.PaymentMethodCard
	PaymentMethodPayPal       = types.PaymentMethodPayPal
	PaymentMethodWire         = types.PaymentMethodWire
	PaymentMethodApplePay     = types.PaymentMethodApplePay
	PaymentMethodGooglePay    = types.PaymentMethodGooglePay
	PaymentMethodStripe       = types.PaymentMethodStripe
	PaymentMethodSquare       = types.PaymentMethodSquare
	PaymentMethodVenmo        = types.PaymentMethodVenmo
	PaymentMethodZelle        = types.PaymentMethodZelle
	PaymentMethodACH          = types.PaymentMethodACH
	PaymentMethodCheck        = types.PaymentMethodCheck
	PaymentMethodCash         = types.PaymentMethodCash
	PaymentMethodCrypto       = types.PaymentMethodCrypto
	PaymentMethodBankTransfer = types.PaymentMethodBankTransfer
	PaymentMethodGiftCard     = types.PaymentMethodGiftCard
	PaymentMethodStoreCredit  = types.PaymentMethodStoreCredit

	// Channel Types
	ChannelDirect    = types.ChannelDirect
	ChannelOrganic   = types.ChannelOrganic
	ChannelPaid      = types.ChannelPaid
	ChannelSocial    = types.ChannelSocial
	ChannelEmail     = types.ChannelEmail
	ChannelSMS       = types.ChannelSMS
	ChannelPush      = types.ChannelPush
	ChannelReferral  = types.ChannelReferral
	ChannelAffiliate = types.ChannelAffiliate
	ChannelDisplay   = types.ChannelDisplay
	ChannelVideo     = types.ChannelVideo
	ChannelAudio     = types.ChannelAudio
	ChannelPrint     = types.ChannelPrint
	ChannelEvent     = types.ChannelEvent
	ChannelWebinar   = types.ChannelWebinar
	ChannelPodcast   = types.ChannelPodcast

	// Traffic Sources
	SourceGoogle    = types.SourceGoogle
	SourceFacebook  = types.SourceFacebook
	SourceTwitter   = types.SourceTwitter
	SourceLinkedIn  = types.SourceLinkedIn
	SourceInstagram = types.SourceInstagram
	SourceYouTube   = types.SourceYouTube
	SourceTikTok    = types.SourceTikTok
	SourceSnapchat  = types.SourceSnapchat
	SourcePinterest = types.SourcePinterest
	SourceReddit    = types.SourceReddit
	SourceBing      = types.SourceBing
	SourceYahoo     = types.SourceYahoo
	SourceDuckDuckGo = types.SourceDuckDuckGo
	SourceNewsletter = types.SourceNewsletter
	SourceEmail     = types.SourceEmail
	SourceBlog      = types.SourceBlog
	SourcePodcast   = types.SourcePodcast
	SourceWebinar   = types.SourceWebinar
	SourcePartner   = types.SourcePartner
	SourceAffiliate = types.SourceAffiliate
	SourceDirect    = types.SourceDirect
	SourceOrganic   = types.SourceOrganic
	SourceUnknown   = types.SourceUnknown

	// Device Types
	DeviceDesktop = types.DeviceDesktop
	DeviceMobile  = types.DeviceMobile
	DeviceTablet  = types.DeviceTablet
	DeviceTV      = types.DeviceTV
	DeviceWatch   = types.DeviceWatch
	DeviceVR      = types.DeviceVR
	DeviceIoT     = types.DeviceIoT
	DeviceBot     = types.DeviceBot
	DeviceUnknown = types.DeviceUnknown

	// Operating Systems
	OSWindows     = types.OSWindows
	OSMacOS       = types.OSMacOS
	OSLinux       = types.OSLinux
	OSiOS         = types.OSiOS
	OSAndroid     = types.OSAndroid
	OSChromeOS    = types.OSChromeOS
	OSFireOS      = types.OSFireOS
	OSWebOS       = types.OSWebOS
	OSTizen       = types.OSTizen
	OSWatchOS     = types.OSWatchOS
	OStvOS        = types.OStvOS
	OSPlayStation = types.OSPlayStation
	OSXbox        = types.OSXbox
	OSUnknown     = types.OSUnknown

	// Browsers
	BrowserChrome  = types.BrowserChrome
	BrowserSafari  = types.BrowserSafari
	BrowserFirefox = types.BrowserFirefox
	BrowserEdge    = types.BrowserEdge
	BrowserOpera   = types.BrowserOpera
	BrowserIE      = types.BrowserIE
	BrowserSamsung = types.BrowserSamsung
	BrowserUC      = types.BrowserUC
	BrowserOther   = types.BrowserOther
	BrowserUnknown = types.BrowserUnknown

	// Subscription Intervals
	IntervalDaily     = types.IntervalDaily
	IntervalWeekly    = types.IntervalWeekly
	IntervalMonthly   = types.IntervalMonthly
	IntervalQuarterly = types.IntervalQuarterly
	IntervalYearly    = types.IntervalYearly
	IntervalAnnual    = types.IntervalAnnual
	IntervalLifetime  = types.IntervalLifetime
	IntervalCustom    = types.IntervalCustom

	// Plan Types
	PlanFree         = types.PlanFree
	PlanFreemium     = types.PlanFreemium
	PlanBasic        = types.PlanBasic
	PlanStandard     = types.PlanStandard
	PlanProfessional = types.PlanProfessional
	PlanPremium      = types.PlanPremium
	PlanEnterprise   = types.PlanEnterprise
	PlanCustom       = types.PlanCustom
	PlanTrial        = types.PlanTrial
	PlanBeta         = types.PlanBeta

	// User Roles
	RoleOwner     = types.RoleOwner
	RoleAdmin     = types.RoleAdmin
	RoleManager   = types.RoleManager
	RoleUser      = types.RoleUser
	RoleGuest     = types.RoleGuest
	RoleViewer    = types.RoleViewer
	RoleEditor    = types.RoleEditor
	RoleModerator = types.RoleModerator
	RoleSupport   = types.RoleSupport
	RoleDeveloper = types.RoleDeveloper
	RoleAnalyst   = types.RoleAnalyst
	RoleBilling   = types.RoleBilling

	// Company Sizes
	CompanySolopreneur = types.CompanySolopreneur
	CompanySmall       = types.CompanySmall
	CompanyMedium      = types.CompanyMedium
	CompanyLarge       = types.CompanyLarge
	CompanyEnterprise  = types.CompanyEnterprise
	CompanyMegaCorp    = types.CompanyMegaCorp
	CompanyUnknown     = types.CompanyUnknown

	// Industries
	IndustryTechnology    = types.IndustryTechnology
	IndustryFinance       = types.IndustryFinance
	IndustryHealthcare    = types.IndustryHealthcare
	IndustryEducation     = types.IndustryEducation
	IndustryEcommerce     = types.IndustryEcommerce
	IndustryRetail        = types.IndustryRetail
	IndustryManufacturing = types.IndustryManufacturing
	IndustryRealEstate    = types.IndustryRealEstate
	IndustryMedia         = types.IndustryMedia
	IndustryNonProfit     = types.IndustryNonProfit
	IndustryGovernment    = types.IndustryGovernment
	IndustryConsulting    = types.IndustryConsulting
	IndustryLegal         = types.IndustryLegal
	IndustryMarketing     = types.IndustryMarketing
	IndustryOther         = types.IndustryOther
	IndustryUnknown       = types.IndustryUnknown
)

// Logging protocol
func (c *Client) Log(ctx context.Context, entry LogEntry) error {
	return c.internal.Log(ctx, entry)
}

func (c *Client) LogInfo(ctx context.Context, service, message string, data map[string]interface{}) error {
	return c.internal.LogInfo(ctx, service, "", message, data)
}

func (c *Client) LogError(ctx context.Context, service, message string, data map[string]interface{}) error {
	return c.internal.LogError(ctx, service, "", message, data)
}

func (c *Client) LogDebug(ctx context.Context, service, message string, data map[string]interface{}) error {
	return c.internal.LogDebug(ctx, service, "", message, data)
}

func (c *Client) LogWarning(ctx context.Context, service, message string, data map[string]interface{}) error {
	return c.internal.LogWarning(ctx, service, "", message, data)
}

func (c *Client) LogCritical(ctx context.Context, service, message string, data map[string]interface{}) error {
	return c.internal.LogCritical(ctx, service, "", message, data)
}

func (c *Client) LogAlert(ctx context.Context, service, message string, data map[string]interface{}) error {
	return c.internal.LogAlert(ctx, service, "", message, data)
}

func (c *Client) LogEmergency(ctx context.Context, service, message string, data map[string]interface{}) error {
	return c.internal.LogEmergency(ctx, service, "", message, data)
}

func (c *Client) LogNotice(ctx context.Context, service, message string, data map[string]interface{}) error {
	return c.internal.LogNotice(ctx, service, "", message, data)
}

func (c *Client) LogTrace(ctx context.Context, service, message string, data map[string]interface{}) error {
	return c.internal.LogTrace(ctx, service, "", message, data)
}

func (c *Client) LogBatch(ctx context.Context, entries []LogEntry) error {
	return c.internal.LogBatch(ctx, entries)
}

// Re-export log types
type (
	LogEntry     = types.LogEntry
	LogLevel     = types.LogLevel
	LogEventType = types.LogEventType
)

// Re-export log constants
const (
	// Log levels
	LogEmergency = types.LogEmergency
	LogAlert     = types.LogAlert
	LogCritical  = types.LogCritical
	LogError     = types.LogError
	LogWarning   = types.LogWarning
	LogNotice    = types.LogNotice
	LogInfo      = types.LogInfo
	LogDebug     = types.LogDebug
	LogTrace     = types.LogTrace

	// Log event types
	LogCollect = types.LogCollect
	LogEnrich  = types.LogEnrich
)

// Version returns detailed version information
func Version() version.Info {
	return version.Get()
}
