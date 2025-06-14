# Usercanal Analytics SDK for Go

<p align="center">
  <a href="https://pkg.go.dev/github.com/usercanal/sdk-go"><img src="https://pkg.go.dev/badge/github.com/usercanal/sdk-go.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/usercanal/sdk-go"><img src="https://goreportcard.com/badge/github.com/usercanal/sdk-go" alt="Go Report Card"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
</p>

## Overview

Usercanal is a unified, high-performance SDK for **analytics events** and **structured logging**. Built for applications where performance matters, it combines the best of product analytics (like PostHog/Mixpanel) with enterprise-grade logging in a single, lightweight package.

## Official SDKs

| SDK | Events | Logs |
|-----|---------|---------|
| [Go SDK](https://github.com/usercanal/sdk-go) | ✅ | ✅ |
| [TypeScript SDK](https://github.com/usercanal/sdk-ts) | ⏳ | ⏳ |
| [Swift SDK](https://github.com/usercanal/sdk-swift) | ⏳ | ⏳ |

## Key Features

### Core Infrastructure
- **Ultra-lightweight**: Minimal impact on your application's performance. No serialization/deserialization
- **Binary FlatBuffers protocol**: Zero-copy processing for maximum speed (20M+ events/second throughput)
- **Smart batching**: Automatic batching with configurable size and intervals for optimal performance
- **Built-in authentication**: API key-based auth in every batch header for enterprise security
- **Batch tracing**: Unique batch IDs for tracking data delivery and debugging lost data
- **Built-in high availability**: DNS failover, smart reconnection, backoff, zero data loss guarantees
- **Unified transport**: Events and logs share the same optimized TCP pipeline
- **Workspace isolation**: Multi-tenant architecture with API key-based workspace routing
- **Real-time processing**: Sub-millisecond event processing with immediate routing
- **Schema validation**: Type-safe binary schemas prevent data corruption at the protocol level
- **Production monitoring**: Built-in metrics and health checking for observability

### Event Analytics Features
- **User behavior tracking**: Product analytics with funnel and cohort support
- **Revenue tracking**: Built-in subscription and one-time payment tracking with currency support
- **Type-safe constants**: Pre-defined event types for common user actions
- **User & group analytics**: Individual and cohort behavior analysis

### Advanced Logging Features
- **Binary structured logging**: Unlike traditional syslog, all data sent in optimized binary format
- **High-throughput batch processing**: Log batching capabilities not available in traditional syslog
- **Context correlation**: Distributed tracing via context IDs across microservices
- **Multiple severity levels**: Standard levels from EMERGENCY to TRACE
- **Enrichment pipeline**: Experimental feature for metadata enhancement
- **Service isolation**: Clear service/source identification for multi-service architectures
- **Workspace routing**: Automatic tenant isolation based on API key authentication
- **Real-time processing**: Immediate routing and processing without buffering delays

## Installation

Install the SDK with a single command:
```bash
go get github.com/usercanal/sdk-go
```

## Quick Start: Analytics Events

Perfect for tracking user behavior, feature usage, and business metrics:

```go
import (
    "context"
    "github.com/usercanal/sdk-go"
)

func main() {
    // Initialize the client
    client, _ := usercanal.NewClient("YOUR_API_KEY")
    defer client.Close(context.Background())

    ctx := context.Background()

    // Track user events - simplified API
    client.Event(ctx, "user_123", usercanal.FeatureUsed, usercanal.Properties{
        "feature_name": "export",
        "duration_ms":  1500,
    })

    // Track revenue - simplified parameters
    client.EventRevenue(ctx, "user_123", "order_456", 29.99, usercanal.CurrencyUSD, usercanal.Properties{
        "type": usercanal.RevenueTypeSubscription,
    })

    // Identify users
    client.EventIdentify(ctx, "user_123", usercanal.Properties{
        "name":  "John Doe",
        "email": "john@example.com",
    })

    client.Flush(ctx) // Ensure delivery
}
```

## Quick Start: Structured Logging

Perfect for application monitoring, debugging, and observability:

```go
import (
    "context"
    "github.com/usercanal/sdk-go"
)

func main() {
    client, _ := usercanal.NewClient("YOUR_API_KEY")
    defer client.Close(context.Background())

    ctx := context.Background()

    // Simple logging - hostname auto-set by SDK
    client.LogInfo(ctx, "api-server", "User login successful", map[string]interface{}{
        "user_id": "123",
        "method":  "oauth",
    })

    client.LogError(ctx, "payment-service", "Payment failed", map[string]interface{}{
        "user_id": "123",
        "amount":  99.99,
        "reason":  "insufficient_funds",
    })

    client.Flush(ctx) // Ensure delivery
}
```

## Configuration

```go
// Simple - uses defaults
client, _ := usercanal.NewClient("YOUR_API_KEY")

// Advanced configuration
client, _ := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
    Endpoint:      "collect.usercanal.com:50000", // Production endpoint
    BatchSize:     100,                          // Events/logs per batch
    FlushInterval: 5 * time.Second,              // Max time between sends
    MaxRetries:    3,                            // Retry attempts
    Debug:         true,                         // Enable debug logging
})
```

## Protocol Advantages

### vs Traditional Logging (syslog, etc.)
| Feature | Usercanal | Traditional Syslog |
|---------|-----------|-------------------|
| **Format** | Binary (FlatBuffers) | Text-based |
| **Batching** | ✅ Built-in | ❌ Single messages |
| **Authentication** | ✅ API keys | ❌ Network-based only |
| **Context IDs** | ✅ Distributed tracing | ❌ No correlation |
| **Performance** | Zero-copy processing | Text parsing overhead |
| **Batch Tracking** | ✅ Unique batch IDs | ❌ No delivery tracking |

## API Overview

```go
// Event Tracking
client.Event(ctx, userID, eventName, properties)
client.EventIdentify(ctx, userID, traits)
client.EventGroup(ctx, userID, groupID, properties)
client.EventRevenue(ctx, userID, orderID, amount, currency, properties)

// Structured Logging (hostname auto-set)
client.LogInfo(ctx, service, message, data)
client.LogError(ctx, service, message, data)
// + LogDebug, LogWarning, LogCritical, LogAlert, LogEmergency, LogNotice, LogTrace

// Management
client.Flush(ctx)        // Force send
client.Close(ctx)        // Graceful shutdown
client.GetStats()        // Connection & queue stats
```

## Why Usercanal?

- **Unified SDK**: One library for both analytics and logging
- **Performance First**: Binary protocol with zero-copy processing
- **Enterprise Ready**: Built-in auth, batching, and delivery tracking
- **Developer Friendly**: Type-safe APIs with built-in constants
- **Lightning Fast Collector**: Non-blocking pipeline with 20M+ events/second throughput
- **Resource Efficient**: Minimal CPU and memory footprint

## Local Development

For local development:

```go
client, _ := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
    Endpoint: "localhost:50000",
    Debug:    true,
})
```

## Documentation

### SDK Documentation
- [📋 Features Overview](FEATURES.md) - Complete feature list and capabilities

### External Resources
- [📚 Full Documentation](https://docs.usercanal.com/docs/SDKs/go) - Dive deeper into the SDK capabilities.
- [📊 Dashboard](https://app.usercanal.com) - Manage your analytics in real-time. (coming)
- [🚀 Collector](https://github.com/usercanal/cdp-collector) - Self-hosted data pipeline

## License

This SDK is distributed under the MIT License. See the [LICENSE](LICENSE) for more information.
