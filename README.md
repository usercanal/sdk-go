# Usercanal Analytics SDK for Go

<p align="center">
  <a href="https://pkg.go.dev/github.com/usercanal/sdk-go"><img src="https://pkg.go.dev/badge/github.com/usercanal/sdk-go.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/usercanal/sdk-go"><img src="https://goreportcard.com/badge/github.com/usercanal/sdk-go" alt="Go Report Card"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
</p>

## Overview

Usercanal is an ultra-lightweight event tracking SDK designed for maximum performance and minimal overhead. Think PostHog, but faster and lighter. Perfect for applications where performance is crucial.

## Official SDKs
| SDK | Status | Version |
|-----|---------|---------|
| [Go SDK](https://github.com/usercanal/sdk-go) | ✅ | v1.0 Beta |
| [TypeScript SDK](https://github.com/usercanal/sdk-ts) | ⏳ | Soon |
| [Swift SDK](https://github.com/usercanal/sdk-swift) | ⏳ | Soon |

## Key Features

- **Ultra-lightweight**: Minimal impact on your application's performance
- **High-performance TCP communication**: Direct, efficient event delivery
- **Automatic batching**: Smart event grouping for optimal throughput
- **Built-in high availability**
  - Automatic DNS-based failover
  - Smart reconnection with exponential backoff
  - No data loss during collector upgrades or outages
- **Built-in support for revenue and subscription tracking** with type-safe constants
- **User & Group Analytics** to understand user behavior and group dynamics

## Installation

Install the SDK with a single command:
```bash
go get github.com/usercanal/sdk-go
```

## Quick Start

Start tracking events in just a few minutes!

```go
import (
    "context"
    "github.com/usercanal/sdk-go"
)

func main() {
    // Initialize the UserCanal client
    canal := usercanal.NewClient("YOUR_API_KEY")
    defer canal.Close()

    // Track a simple event
    ctx := context.Background()
    canal.Track(ctx, usercanal.Event{
        UserId: "123",
        Name:   usercanal.FeatureUsed,
        Properties: usercanal.Properties{
            "feature_name": "export",
            "duration_ms": 1500,
        },
    })
}
```

### Advanced Configuration
```go
// Initialize with custom configuration
canal := usercanal.NewClient("YOUR_API_KEY", usercanal.Config{
    Endpoint:      "collect.usercanal.com:9000",
    BatchSize:     100,
    FlushInterval: 5 * time.Second,
    MaxRetries:    3,
    Debug:         true,
})
```

## Why Usercanal?

- **Minimal Overhead**: Designed to be as lightweight as possible
- **Maximum Performance**: Direct TCP communication for faster event delivery
- **Simple Integration**: Start tracking events in minutes
- **Type-Safe**: Built-in constants for common events
- **Resource Efficient**: Optimized for minimal CPU and memory usage


## Explore More

- [📚 Full Documentation](https://usercanal.com/docs/sdks/go) - Dive deeper into the SDK capabilities.
- [📊 Dashboard](https://app.usercanal.com) - Manage your analytics in real-time. (coming)

## License

This SDK is distributed under the MIT License. See the [LICENSE](LICENSE) for more information.
