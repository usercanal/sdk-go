# Usercanal Analytics SDK for Go

<p align="center">
  <a href="https://pkg.go.dev/github.com/usercanal/sdk-go"><img src="https://pkg.go.dev/badge/github.com/usercanal/sdk-go.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/usercanal/sdk-go"><img src="https://goreportcard.com/badge/github.com/usercanal/sdk-go" alt="Go Report Card"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
</p>

## Overview

Usercanal helps you track user behavior and business metrics across your applications with ease. Built on gRPC, our Go SDK offers efficient event tracking, automatic batching, and type-safe constants, all while providing flexibility for custom analytics.

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

## Key Features

- **Efficient gRPC-based communication** for low-latency event tracking
- **Automatic batching and retry** mechanisms for reliable data delivery
- **Built-in high availability**
  - Automatic DNS-based failover
  - Smart reconnection with exponential backoff
  - No data loss during collector upgrades or outages
- **Built-in support for revenue and subscription tracking** with type-safe constants
- **User & Group Analytics** to understand user behavior and group dynamics

## Explore More

- [📚 Full Documentation](https://usercanal.com/docs/sdks/go) - Dive deeper into the SDK capabilities.
- [🔬 API Reference](https://pkg.go.dev/github.com/usercanal/sdk-go) - Comprehensive API details.
- [💡 Example Apps](https://github.com/usercanal/examples) - Learn from practical examples.
- [📊 Dashboard](https://app.usercanal.com) - Manage your analytics in real-time.

## Join Our Community

- [Join our Discord](https://discord.gg/usercanal) to connect with other developers and get support.
- [Contact Support](mailto:support@usercanal.com) for any questions or issues.

## License

This SDK is distributed under the MIT License. See the [LICENSE](LICENSE) for more information.
