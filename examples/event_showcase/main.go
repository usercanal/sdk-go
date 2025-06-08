// sdk-go/examples/cdp_constants_showcase.go
package main

import (
	"context"
	"fmt"
	"log"

	usercanal "github.com/usercanal/sdk-go"
)

func main() {
	// Initialize client with debug mode
	client, err := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
		Debug: true,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() {
		if err := client.Close(context.Background()); err != nil {
			log.Printf("Failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	fmt.Println("=== UserCanal CDP Constants Showcase ===")
	fmt.Println("Demonstrating comprehensive constants for global CDP compatibility")

	// 1. User Registration with Global Context
	fmt.Println("\n1. Global User Registration:")
	err = client.Event(ctx, "user_global_123", usercanal.UserSignedUp, usercanal.Properties{
		"auth_method":     usercanal.AuthMethodGoogle,
		"source":          usercanal.SourceGoogle,
		"channel":         usercanal.ChannelOrganic,
		"device_type":     usercanal.DeviceMobile,
		"operating_system": usercanal.OSAndroid,
		"browser":         usercanal.BrowserChrome,
		"country":         "Japan",
		"plan_type":       usercanal.PlanTrial,
	})
	if err != nil {
		log.Printf("Failed to track global signup: %v", err)
	}

	// 2. Multi-Currency Revenue Tracking
	fmt.Println("\n2. Global Revenue Events:")
	
	// Japanese customer
	err = client.EventRevenue(ctx, "user_japan_456", "ord_jp_789", 12000, usercanal.CurrencyJPY, usercanal.Properties{
		"type":            usercanal.RevenueTypeSubscription,
		"plan_type":       usercanal.PlanPremium,
		"interval":        usercanal.IntervalMonthly,
		"payment_method":  usercanal.PaymentMethodCard,
		"source":          usercanal.SourceLinkedIn,
		"channel":         usercanal.ChannelPaid,
	})
	if err != nil {
		log.Printf("Failed to track Japanese revenue: %v", err)
	}

	// Brazilian customer with crypto payment
	err = client.EventRevenue(ctx, "user_brazil_789", "ord_br_456", 500, usercanal.CurrencyBRL, usercanal.Properties{
		"type":           usercanal.RevenueTypeOneTime,
		"plan_type":      usercanal.PlanEnterprise,
		"payment_method": usercanal.PaymentMethodCrypto,
		"crypto_type":    usercanal.CurrencyBTC,
		"source":         usercanal.SourceDirect,
	})
	if err != nil {
		log.Printf("Failed to track Brazilian revenue: %v", err)
	}

	// 3. B2B Company Identification
	fmt.Println("\n3. B2B Company Segmentation:")
	err = client.EventIdentify(ctx, "user_enterprise_001", usercanal.Properties{
		"name":         "Sarah Chen",
		"email":        "sarah@techcorp.com",
		"role":         usercanal.RoleAdmin,
		"company_size": usercanal.CompanyLarge,
		"industry":     usercanal.IndustryTechnology,
		"plan_type":    usercanal.PlanEnterprise,
		"country":      "Singapore",
		"currency":     usercanal.CurrencySGD,
	})
	if err != nil {
		log.Printf("Failed to identify enterprise user: %v", err)
	}

	// 4. Subscription Management Events
	fmt.Println("\n4. Subscription Lifecycle:")
	
	// Trial conversion
	err = client.Event(ctx, "user_trial_convert", usercanal.TrialConverted, usercanal.Properties{
		"previous_plan": usercanal.PlanTrial,
		"new_plan":      usercanal.PlanProfessional,
		"interval":      usercanal.IntervalYearly,
		"currency":      usercanal.CurrencyEUR,
		"discount":      "ANNUAL20",
		"source":        usercanal.SourceEmail,
		"channel":       usercanal.ChannelEmail,
	})
	if err != nil {
		log.Printf("Failed to track trial conversion: %v", err)
	}

	// Subscription upgrade
	err = client.Event(ctx, "user_upgrade_123", usercanal.SubscriptionChanged, usercanal.Properties{
		"previous_plan":    usercanal.PlanBasic,
		"new_plan":         usercanal.PlanPremium,
		"previous_amount":  29.99,
		"new_amount":       99.99,
		"currency":         usercanal.CurrencyUSD,
		"upgrade_reason":   "need more features",
		"payment_method":   usercanal.PaymentMethodApplePay,
		"device_type":      usercanal.DeviceMobile,
		"operating_system": usercanal.OSiOS,
	})
	if err != nil {
		log.Printf("Failed to track subscription upgrade: %v", err)
	}

	// 5. Multi-Channel Attribution
	fmt.Println("\n5. Multi-Channel Attribution:")
	
	channels := []struct {
		channel usercanal.Channel
		source  usercanal.Source
		event   usercanal.EventName
	}{
		{usercanal.ChannelSocial, usercanal.SourceFacebook, usercanal.PageViewed},
		{usercanal.ChannelEmail, usercanal.SourceNewsletter, usercanal.EmailClicked},
		{usercanal.ChannelPaid, usercanal.SourceGoogle, usercanal.FeatureUsed},
		{usercanal.ChannelWebinar, usercanal.SourceWebinar, usercanal.UserOnboarded},
		{usercanal.ChannelPodcast, usercanal.SourcePodcast, usercanal.TrialStarted},
	}

	for i, attr := range channels {
		err = client.Event(ctx, fmt.Sprintf("user_attribution_%d", i), attr.event, usercanal.Properties{
			"channel":    attr.channel,
			"source":     attr.source,
			"session_id": fmt.Sprintf("sess_%d", i),
			"timestamp":  "2024-01-15T10:30:00Z",
		})
		if err != nil {
			log.Printf("Failed to track attribution %d: %v", i, err)
		}
	}

	// 6. Cross-Platform Device Tracking
	fmt.Println("\n6. Cross-Platform Usage:")
	
	devices := []struct {
		deviceType usercanal.DeviceType
		os         usercanal.OperatingSystem
		browser    usercanal.Browser
	}{
		{usercanal.DeviceDesktop, usercanal.OSWindows, usercanal.BrowserChrome},
		{usercanal.DeviceMobile, usercanal.OSiOS, usercanal.BrowserSafari},
		{usercanal.DeviceTablet, usercanal.OSAndroid, usercanal.BrowserFirefox},
		{usercanal.DeviceTV, usercanal.OStvOS, usercanal.BrowserOther},
		{usercanal.DeviceWatch, usercanal.OSWatchOS, usercanal.BrowserUnknown},
	}

	for i, device := range devices {
		err = client.Event(ctx, "user_crossplatform_123", usercanal.FeatureUsed, usercanal.Properties{
			"feature_name":     "sync_data",
			"device_type":      device.deviceType,
			"operating_system": device.os,
			"browser":          device.browser,
			"session_count":    i + 1,
		})
		if err != nil {
			log.Printf("Failed to track device %d: %v", i, err)
		}
	}

	// 7. Payment Method Diversity
	fmt.Println("\n7. Global Payment Methods:")
	
	payments := []struct {
		method   usercanal.PaymentMethod
		currency usercanal.Currency
		country  string
	}{
		{usercanal.PaymentMethodCard, usercanal.CurrencyUSD, "United States"},
		{usercanal.PaymentMethodPayPal, usercanal.CurrencyEUR, "Germany"},
		{usercanal.PaymentMethodApplePay, usercanal.CurrencyCAD, "Canada"},
		{usercanal.PaymentMethodGooglePay, usercanal.CurrencyINR, "India"},
		{usercanal.PaymentMethodCrypto, usercanal.CurrencyBTC, "Global"},
		{usercanal.PaymentMethodACH, usercanal.CurrencyUSD, "United States"},
		{usercanal.PaymentMethodBankTransfer, usercanal.CurrencyGBP, "United Kingdom"},
	}

	for i, payment := range payments {
		err = client.Event(ctx, fmt.Sprintf("user_payment_%d", i), usercanal.PaymentMethodAdded, usercanal.Properties{
			"payment_method": payment.method,
			"currency":       payment.currency,
			"country":        payment.country,
			"is_primary":     i == 0,
			"verification":   "completed",
		})
		if err != nil {
			log.Printf("Failed to track payment method %d: %v", i, err)
		}
	}

	// 8. Industry-Specific Segmentation
	fmt.Println("\n8. Industry Segmentation:")
	
	industries := []struct {
		industry usercanal.Industry
		size     usercanal.CompanySize
		plan     usercanal.PlanType
	}{
		{usercanal.IndustryTechnology, usercanal.CompanyMedium, usercanal.PlanProfessional},
		{usercanal.IndustryFinance, usercanal.CompanyEnterprise, usercanal.PlanEnterprise},
		{usercanal.IndustryHealthcare, usercanal.CompanyLarge, usercanal.PlanPremium},
		{usercanal.IndustryEducation, usercanal.CompanySmall, usercanal.PlanBasic},
		{usercanal.IndustryNonProfit, usercanal.CompanySolopreneur, usercanal.PlanFree},
	}

	for i, seg := range industries {
		err = client.EventGroup(ctx, fmt.Sprintf("user_industry_%d", i), fmt.Sprintf("org_%s_%d", seg.industry, i), usercanal.Properties{
			"industry":     seg.industry,
			"company_size": seg.size,
			"plan_type":    seg.plan,
			"region":       "North America",
			"founded_year": 2020 - i,
		})
		if err != nil {
			log.Printf("Failed to track industry group %d: %v", i, err)
		}
	}

	// Flush all events
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush events: %v", err)
	}

	// Summary
	fmt.Println("\n✅ CDP Constants Showcase Complete!")
	fmt.Println("\n📊 Constants Demonstrated:")
	fmt.Printf("• %d Event Types (User, Revenue, Subscription, Trial, Shopping, Communication)\n", 44)
	fmt.Printf("• %d Currencies (including crypto: BTC, ETH, USDC, USDT)\n", 71)
	fmt.Printf("• %d Payment Methods (traditional + modern: Apple Pay, crypto, etc.)\n", 16)
	fmt.Printf("• %d Channel Types (organic, paid, social, email, etc.)\n", 16)
	fmt.Printf("• %d Traffic Sources (Google, Facebook, newsletters, etc.)\n", 23)
	fmt.Printf("• %d Device Types (desktop, mobile, IoT, VR, etc.)\n", 9)
	fmt.Printf("• %d Operating Systems (Windows, macOS, iOS, Android, etc.)\n", 14)
	fmt.Printf("• %d Browsers (Chrome, Safari, Firefox, etc.)\n", 10)
	fmt.Printf("• %d Auth Methods (Google, GitHub, SSO, etc.)\n", 5)
	fmt.Printf("• %d Subscription Intervals (daily to lifetime)\n", 8)
	fmt.Printf("• %d Plan Types (free to enterprise)\n", 10)
	fmt.Printf("• %d User Roles (owner to analyst)\n", 12)
	fmt.Printf("• %d Company Sizes (solopreneur to mega corp)\n", 7)
	fmt.Printf("• %d Industries (technology to non-profit)\n", 15)

	fmt.Println("\n🌍 Global CDP Benefits:")
	fmt.Println("• Complete international currency support (71 currencies)")
	fmt.Println("• Modern payment methods (crypto, digital wallets)")
	fmt.Println("• Cross-platform device tracking")
	fmt.Println("• Multi-channel attribution")
	fmt.Println("• B2B segmentation (roles, company sizes, industries)")
	fmt.Println("• SaaS business model support (trials, subscriptions)")
	fmt.Println("• Human-readable values for dashboard display")
	fmt.Println("• Consistent with major CDP platforms (Segment, Mixpanel, Amplitude)")
}