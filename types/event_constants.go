// sdk-go/types/event_constants.go
package types

// EventName represents a strongly typed event name that also allows custom strings
type EventName string

// String returns the string representation of the event name
func (e EventName) String() string {
	return string(e)
}

// IsStandardEvent checks if the event name is a predefined standard event
func (e EventName) IsStandardEvent() bool {
	switch e {
	// User lifecycle events
	case UserSignedUp, UserSignedIn, UserSignedOut, UserInvited, UserOnboarded,
		AuthenticationFailed, PasswordReset, TwoFactorEnabled, TwoFactorDisabled:
		return true
	// Revenue & Billing
	case OrderCompleted, OrderRefunded, OrderCanceled, PaymentFailed,
		PaymentMethodAdded, PaymentMethodUpdated, PaymentMethodRemoved:
		return true
	// Subscription Management
	case SubscriptionStarted, SubscriptionRenewed, SubscriptionPaused,
		SubscriptionResumed, SubscriptionChanged, SubscriptionCanceled:
		return true
	// Trial & Conversion
	case TrialStarted, TrialEndingSoon, TrialEnded, TrialConverted:
		return true
	// Shopping Experience
	case CartViewed, CartUpdated, CartAbandoned, CheckoutStarted, CheckoutCompleted:
		return true
	// Product Engagement
	case PageViewed, FeatureUsed, SearchPerformed, FileUploaded,
		NotificationSent, NotificationClicked:
		return true
	// Communication
	case EmailSent, EmailOpened, EmailClicked, EmailBounced, EmailUnsubscribed,
		SupportTicketCreated, SupportTicketResolved:
		return true
	}
	return false
}

// Authentication & User Management Events
const (
	UserSignedUp         EventName = "User Signed Up"
	UserSignedIn         EventName = "User Signed In"
	UserSignedOut        EventName = "User Signed Out"
	UserInvited          EventName = "User Invited"
	UserOnboarded        EventName = "User Onboarded"
	AuthenticationFailed EventName = "Authentication Failed"
	PasswordReset        EventName = "Password Reset"
	TwoFactorEnabled     EventName = "Two Factor Enabled"
	TwoFactorDisabled    EventName = "Two Factor Disabled"
)

// Revenue & Billing Events
const (
	OrderCompleted        EventName = "Order Completed"
	OrderRefunded         EventName = "Order Refunded"
	OrderCanceled         EventName = "Order Canceled"
	PaymentFailed         EventName = "Payment Failed"
	PaymentMethodAdded    EventName = "Payment Method Added"
	PaymentMethodUpdated  EventName = "Payment Method Updated"
	PaymentMethodRemoved  EventName = "Payment Method Removed"
)

// Subscription Management Events
const (
	SubscriptionStarted  EventName = "Subscription Started"
	SubscriptionRenewed  EventName = "Subscription Renewed"
	SubscriptionPaused   EventName = "Subscription Paused"
	SubscriptionResumed  EventName = "Subscription Resumed"
	SubscriptionChanged  EventName = "Subscription Changed"
	SubscriptionCanceled EventName = "Subscription Canceled"
)

// Trial & Conversion Events
const (
	TrialStarted     EventName = "Trial Started"
	TrialEndingSoon  EventName = "Trial Ending Soon"
	TrialEnded       EventName = "Trial Ended"
	TrialConverted   EventName = "Trial Converted"
)

// Shopping Experience Events
const (
	CartViewed        EventName = "Cart Viewed"
	CartUpdated       EventName = "Cart Updated"
	CartAbandoned     EventName = "Cart Abandoned"
	CheckoutStarted   EventName = "Checkout Started"
	CheckoutCompleted EventName = "Checkout Completed"
)

// Product Engagement Events
const (
	PageViewed          EventName = "Page Viewed"
	FeatureUsed         EventName = "Feature Used"
	SearchPerformed     EventName = "Search Performed"
	FileUploaded        EventName = "File Uploaded"
	NotificationSent    EventName = "Notification Sent"
	NotificationClicked EventName = "Notification Clicked"
)

// Communication Events
const (
	EmailSent               EventName = "Email Sent"
	EmailOpened             EventName = "Email Opened"
	EmailClicked            EventName = "Email Clicked"
	EmailBounced            EventName = "Email Bounced"
	EmailUnsubscribed       EventName = "Email Unsubscribed"
	SupportTicketCreated    EventName = "Support Ticket Created"
	SupportTicketResolved   EventName = "Support Ticket Resolved"
)

// Authentication Methods
type AuthMethod string

const (
	AuthMethodPassword AuthMethod = "password"
	AuthMethodGoogle   AuthMethod = "google"
	AuthMethodGitHub   AuthMethod = "github"
	AuthMethodSSO      AuthMethod = "sso"
	AuthMethodEmail    AuthMethod = "email"
)

// Revenue Types
type RevenueType string

const (
	RevenueTypeOneTime     RevenueType = "one_time"
	RevenueTypeSubscription RevenueType = "subscription"
)

// Currency Types
type Currency string

const (
	// Major Global Currencies
	CurrencyUSD Currency = "USD" // US Dollar
	CurrencyEUR Currency = "EUR" // Euro
	CurrencyGBP Currency = "GBP" // British Pound
	CurrencyJPY Currency = "JPY" // Japanese Yen
	CurrencyCAD Currency = "CAD" // Canadian Dollar
	CurrencyAUD Currency = "AUD" // Australian Dollar
	CurrencyNZD Currency = "NZD" // New Zealand Dollar
	CurrencyKRW Currency = "KRW" // South Korean Won
	CurrencyCNY Currency = "CNY" // Chinese Yuan
	CurrencyHKD Currency = "HKD" // Hong Kong Dollar
	CurrencySGD Currency = "SGD" // Singapore Dollar
	CurrencyMXN Currency = "MXN" // Mexican Peso
	CurrencyINR Currency = "INR" // Indian Rupee
	CurrencyPLN Currency = "PLN" // Polish Zloty
	CurrencyBRL Currency = "BRL" // Brazilian Real
	CurrencyRUB Currency = "RUB" // Russian Ruble
	CurrencyDKK Currency = "DKK" // Danish Krone
	CurrencyNOK Currency = "NOK" // Norwegian Krone
	CurrencySEK Currency = "SEK" // Swedish Krona
	CurrencyCHF Currency = "CHF" // Swiss Franc
	CurrencyTRY Currency = "TRY" // Turkish Lira
	CurrencyILS Currency = "ILS" // Israeli Shekel
	CurrencyTHB Currency = "THB" // Thai Baht
	CurrencyMYR Currency = "MYR" // Malaysian Ringgit
	CurrencyIDR Currency = "IDR" // Indonesian Rupiah
	CurrencyVND Currency = "VND" // Vietnamese Dong
	CurrencyPHP Currency = "PHP" // Philippine Peso
	CurrencyCZK Currency = "CZK" // Czech Koruna
	CurrencyHUF Currency = "HUF" // Hungarian Forint
	CurrencyZAR Currency = "ZAR" // South African Rand
	CurrencyARS Currency = "ARS" // Argentine Peso
	CurrencyCLP Currency = "CLP" // Chilean Peso
	CurrencyCOP Currency = "COP" // Colombian Peso
	CurrencyPEN Currency = "PEN" // Peruvian Sol
	CurrencyUYU Currency = "UYU" // Uruguayan Peso
	CurrencyEGP Currency = "EGP" // Egyptian Pound
	CurrencyAED Currency = "AED" // UAE Dirham
	CurrencySAR Currency = "SAR" // Saudi Riyal
	CurrencyQAR Currency = "QAR" // Qatari Riyal
	CurrencyBHD Currency = "BHD" // Bahraini Dinar
	CurrencyKWD Currency = "KWD" // Kuwaiti Dinar
	CurrencyOMR Currency = "OMR" // Omani Rial
	CurrencyJOD Currency = "JOD" // Jordanian Dinar
	CurrencyLBP Currency = "LBP" // Lebanese Pound
	CurrencyRON Currency = "RON" // Romanian Leu
	CurrencyBGN Currency = "BGN" // Bulgarian Lev
	CurrencyHRK Currency = "HRK" // Croatian Kuna
	CurrencyRSD Currency = "RSD" // Serbian Dinar
	CurrencyBAM Currency = "BAM" // Bosnia and Herzegovina Mark
	CurrencyMKD Currency = "MKD" // Macedonian Denar
	CurrencyALL Currency = "ALL" // Albanian Lek
	CurrencyUAH Currency = "UAH" // Ukrainian Hryvnia
	CurrencyBYN Currency = "BYN" // Belarusian Ruble
	CurrencyMDL Currency = "MDL" // Moldovan Leu
	CurrencyGEL Currency = "GEL" // Georgian Lari
	CurrencyAMD Currency = "AMD" // Armenian Dram
	CurrencyAZN Currency = "AZN" // Azerbaijani Manat
	CurrencyKZT Currency = "KZT" // Kazakhstani Tenge
	CurrencyUZS Currency = "UZS" // Uzbekistani Som
	CurrencyKGS Currency = "KGS" // Kyrgyzstani Som
	CurrencyTJS Currency = "TJS" // Tajikistani Somoni
	CurrencyTMT Currency = "TMT" // Turkmenistani Manat
	CurrencyMNT Currency = "MNT" // Mongolian Tugrik
	CurrencyBTC Currency = "BTC" // Bitcoin
	CurrencyETH Currency = "ETH" // Ethereum
	CurrencyUSDC Currency = "USDC" // USD Coin
	CurrencyUSDT Currency = "USDT" // Tether
)

// Payment Methods
type PaymentMethod string

const (
	PaymentMethodCard         PaymentMethod = "card"
	PaymentMethodPayPal       PaymentMethod = "paypal"
	PaymentMethodWire         PaymentMethod = "wire"
	PaymentMethodApplePay     PaymentMethod = "apple_pay"
	PaymentMethodGooglePay    PaymentMethod = "google_pay"
	PaymentMethodStripe       PaymentMethod = "stripe"
	PaymentMethodSquare       PaymentMethod = "square"
	PaymentMethodVenmo        PaymentMethod = "venmo"
	PaymentMethodZelle        PaymentMethod = "zelle"
	PaymentMethodACH          PaymentMethod = "ach"
	PaymentMethodCheck        PaymentMethod = "check"
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodCrypto       PaymentMethod = "crypto"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodGiftCard     PaymentMethod = "gift_card"
	PaymentMethodStoreCredit  PaymentMethod = "store_credit"
)

// Channel Types (common in CDP platforms)
type Channel string

const (
	ChannelDirect    Channel = "direct"
	ChannelOrganic   Channel = "organic"
	ChannelPaid      Channel = "paid"
	ChannelSocial    Channel = "social"
	ChannelEmail     Channel = "email"
	ChannelSMS       Channel = "sms"
	ChannelPush      Channel = "push"
	ChannelReferral  Channel = "referral"
	ChannelAffiliate Channel = "affiliate"
	ChannelDisplay   Channel = "display"
	ChannelVideo     Channel = "video"
	ChannelAudio     Channel = "audio"
	ChannelPrint     Channel = "print"
	ChannelEvent     Channel = "event"
	ChannelWebinar   Channel = "webinar"
	ChannelPodcast   Channel = "podcast"
)

// Traffic Sources (standard UTM and attribution)
type Source string

const (
	SourceGoogle    Source = "google"
	SourceFacebook  Source = "facebook"
	SourceTwitter   Source = "twitter"
	SourceLinkedIn  Source = "linkedin"
	SourceInstagram Source = "instagram"
	SourceYouTube   Source = "youtube"
	SourceTikTok    Source = "tiktok"
	SourceSnapchat  Source = "snapchat"
	SourcePinterest Source = "pinterest"
	SourceReddit    Source = "reddit"
	SourceBing      Source = "bing"
	SourceYahoo     Source = "yahoo"
	SourceDuckDuckGo Source = "duckduckgo"
	SourceNewsletter Source = "newsletter"
	SourceEmail     Source = "email"
	SourceBlog      Source = "blog"
	SourcePodcast   Source = "podcast"
	SourceWebinar   Source = "webinar"
	SourcePartner   Source = "partner"
	SourceAffiliate Source = "affiliate"
	SourceDirect    Source = "direct"
	SourceOrganic   Source = "organic"
	SourceUnknown   Source = "unknown"
)

// Device Types (common in analytics)
type DeviceType string

const (
	DeviceDesktop DeviceType = "desktop"
	DeviceMobile  DeviceType = "mobile"
	DeviceTablet  DeviceType = "tablet"
	DeviceTV      DeviceType = "tv"
	DeviceWatch   DeviceType = "watch"
	DeviceVR      DeviceType = "vr"
	DeviceIoT     DeviceType = "iot"
	DeviceBot     DeviceType = "bot"
	DeviceUnknown DeviceType = "unknown"
)

// Operating Systems
type OperatingSystem string

const (
	OSWindows   OperatingSystem = "windows"
	OSMacOS     OperatingSystem = "macos"
	OSLinux     OperatingSystem = "linux"
	OSiOS       OperatingSystem = "ios"
	OSAndroid   OperatingSystem = "android"
	OSChromeOS  OperatingSystem = "chromeos"
	OSFireOS    OperatingSystem = "fireos"
	OSWebOS     OperatingSystem = "webos"
	OSTizen     OperatingSystem = "tizen"
	OSWatchOS   OperatingSystem = "watchos"
	OStvOS      OperatingSystem = "tvos"
	OSPlayStation OperatingSystem = "playstation"
	OSXbox      OperatingSystem = "xbox"
	OSUnknown   OperatingSystem = "unknown"
)

// Browsers
type Browser string

const (
	BrowserChrome  Browser = "chrome"
	BrowserSafari  Browser = "safari"
	BrowserFirefox Browser = "firefox"
	BrowserEdge    Browser = "edge"
	BrowserOpera   Browser = "opera"
	BrowserIE      Browser = "ie"
	BrowserSamsung Browser = "samsung"
	BrowserUC      Browser = "uc"
	BrowserOther   Browser = "other"
	BrowserUnknown Browser = "unknown"
)

func (a AuthMethod) String() string {
	return string(a)
}

func (r RevenueType) String() string {
	return string(r)
}

func (c Currency) String() string {
	return string(c)
}

func (p PaymentMethod) String() string {
	return string(p)
}

func (ch Channel) String() string {
	return string(ch)
}

func (s Source) String() string {
	return string(s)
}

func (d DeviceType) String() string {
	return string(d)
}

func (o OperatingSystem) String() string {
	return string(o)
}

func (b Browser) String() string {
	return string(b)
}

// Subscription Intervals (common in SaaS)
type SubscriptionInterval string

const (
	IntervalDaily    SubscriptionInterval = "daily"
	IntervalWeekly   SubscriptionInterval = "weekly"
	IntervalMonthly  SubscriptionInterval = "monthly"
	IntervalQuarterly SubscriptionInterval = "quarterly"
	IntervalYearly   SubscriptionInterval = "yearly"
	IntervalAnnual   SubscriptionInterval = "annual"
	IntervalLifetime SubscriptionInterval = "lifetime"
	IntervalCustom   SubscriptionInterval = "custom"
)

// Plan Types (common business models)
type PlanType string

const (
	PlanFree       PlanType = "free"
	PlanFreemium   PlanType = "freemium"
	PlanBasic      PlanType = "basic"
	PlanStandard   PlanType = "standard"
	PlanProfessional PlanType = "professional"
	PlanPremium    PlanType = "premium"
	PlanEnterprise PlanType = "enterprise"
	PlanCustom     PlanType = "custom"
	PlanTrial      PlanType = "trial"
	PlanBeta       PlanType = "beta"
)

// User Roles (common in B2B)
type UserRole string

const (
	RoleOwner      UserRole = "owner"
	RoleAdmin      UserRole = "admin"
	RoleManager    UserRole = "manager"
	RoleUser       UserRole = "user"
	RoleGuest      UserRole = "guest"
	RoleViewer     UserRole = "viewer"
	RoleEditor     UserRole = "editor"
	RoleModerator  UserRole = "moderator"
	RoleSupport    UserRole = "support"
	RoleDeveloper  UserRole = "developer"
	RoleAnalyst    UserRole = "analyst"
	RoleBilling    UserRole = "billing"
)

// Company Sizes (common segmentation)
type CompanySize string

const (
	CompanySolopreneur   CompanySize = "solopreneur"
	CompanySmall         CompanySize = "small"        // 1-10
	CompanyMedium        CompanySize = "medium"       // 11-50
	CompanyLarge         CompanySize = "large"        // 51-200
	CompanyEnterprise    CompanySize = "enterprise"   // 201-1000
	CompanyMegaCorp      CompanySize = "mega_corp"    // 1000+
	CompanyUnknown       CompanySize = "unknown"
)

// Industries (common B2B segmentation)
type Industry string

const (
	IndustryTechnology     Industry = "technology"
	IndustryFinance        Industry = "finance"
	IndustryHealthcare     Industry = "healthcare"
	IndustryEducation      Industry = "education"
	IndustryEcommerce      Industry = "ecommerce"
	IndustryRetail         Industry = "retail"
	IndustryManufacturing  Industry = "manufacturing"
	IndustryRealEstate     Industry = "real_estate"
	IndustryMedia          Industry = "media"
	IndustryNonProfit      Industry = "non_profit"
	IndustryGovernment     Industry = "government"
	IndustryConsulting     Industry = "consulting"
	IndustryLegal          Industry = "legal"
	IndustryMarketing      Industry = "marketing"
	IndustryOther          Industry = "other"
	IndustryUnknown        Industry = "unknown"
)

func (si SubscriptionInterval) String() string {
	return string(si)
}

func (pt PlanType) String() string {
	return string(pt)
}

func (ur UserRole) String() string {
	return string(ur)
}

func (cs CompanySize) String() string {
	return string(cs)
}

func (i Industry) String() string {
	return string(i)
}