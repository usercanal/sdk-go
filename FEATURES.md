# SDK Features Overview

## Core Differentiators

### Unified Protocol
- **Single SDK** for both analytics events and structured logging
- **Shared transport** - events and logs use the same optimized connection
- **Independent batching** - separate queues prevent blocking between protocols

### High-Performance Binary Protocol
- **FlatBuffers format** - zero-copy serialization vs JSON/text parsing
- **Batch processing** - configurable batching (size + time-based)
- **20M+ events/second** throughput capability
- **TCP connection pooling** with automatic reconnection

### Enterprise Features
- **API key authentication** embedded in every batch
- **Batch tracking** with unique IDs for delivery verification
- **Zero data loss** with automatic retry and exponential backoff
- **Multi-tenant workspace** isolation

## Event Analytics

### Business Intelligence
- **Revenue tracking** - subscriptions, one-time payments, multi-currency
- **User behavior** - funnels, cohorts, feature usage
- **Type-safe constants** - 44+ predefined events (signup, purchase, etc.)
- **Custom events** - full string flexibility for domain-specific tracking

### User Management
- **Identity tracking** - associate events with users
- **Group analytics** - organizational/team behavior
- **Session correlation** - track user journeys

## Structured Logging

### vs Traditional Logging
| Feature | UserCanal | Syslog/JSON |
|---------|-----------|-------------|
| Format | Binary (efficient) | Text (overhead) |
| Batching | ✅ Built-in | ❌ Message-by-message |
| Authentication | ✅ Per-batch | ❌ Network-only |
| Context IDs | ✅ Distributed tracing | ❌ No correlation |
| Delivery tracking | ✅ Batch IDs | ❌ Fire-and-forget |

### Advanced Capabilities
- **Context correlation** - distributed tracing across microservices
- **Service isolation** - clear service/source identification
- **9 severity levels** - EMERGENCY to TRACE
- **Real-time processing** - sub-millisecond routing

## Developer Experience

### Simple API
```go
// Analytics
client.Event(ctx, userID, eventName, properties)
client.EventRevenue(ctx, userID, orderID, amount, currency, properties)

// Logging (hostname auto-set)
client.LogInfo(ctx, service, message, data)
client.LogError(ctx, service, message, data)
```

### Built-in Reliability
- **Smart reconnection** with DNS failover
- **Graceful degradation** during network issues
- **Thread-safe** concurrent access
- **Context-aware shutdown** with proper cleanup

### Observability
- **Built-in metrics** - queue depth, throughput, error rates
- **Health monitoring** - connection status, delivery tracking
- **Debug support** - configurable verbose logging

## Performance Characteristics

### Resource Efficiency
- **Minimal CPU** overhead from zero-copy processing
- **Bounded memory** usage with configurable limits
- **Single connection** reduces network overhead
- **Async processing** prevents blocking application threads

### Scalability
- **Horizontal scaling** through workspace isolation
- **High throughput** optimized for enterprise workloads
- **Real-time processing** with immediate routing
- **Configurable batching** to balance latency vs efficiency