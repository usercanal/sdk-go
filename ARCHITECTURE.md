# UserCanal Go SDK Architecture

High-performance, unified client for analytics events and structured logging using binary FlatBuffers over TCP.

```
Application ───▶ SDK Client ───▶ UserCanal Collector
```

## Package Structure

### Public API (`usercanal.go`)
- Main client interface with simplified methods
- Type re-exports and constants

### Internal API (`internal/api/`)
- Core business logic and coordination
- Input validation and error handling
- Separate processing for events and logs

### Types (`types/`)
- Data structures and validation rules
- Comprehensive CDP constants (71 currencies, 44+ events)
- Error types and constructors

### Batching (`internal/batch/`)
- Event/log aggregation for performance
- Time and size-based flushing
- Thread-safe queue operations

### Transport (`internal/transport/`)
- TCP communication with auto-retry
- Connection management and health monitoring
- DNS failover support

### Schema (`internal/schema/`, `internal/convert/`)
- FlatBuffers binary protocol
- Zero-copy serialization
- Type-safe conversion between Go types and binary format

## Data Flow

```
Application
    │
    ▼
API Layer - validation & conversion
    │
    ▼
Batch Queue - time/size-based aggregation
    │
    ▼
Transport - TCP transmission with retry
    │
    ▼
UserCanal Collector
```

## Concurrency

- **Thread-safe public API** with channel-based coordination
- **Background goroutines** for batch flushing and connection health
- **Context-aware shutdown** for graceful cleanup

## Protocol

- **FlatBuffers binary format** - zero-copy, type-safe, compact
- **Length-prefixed framing** over TCP
- **Batched payloads** with API key authentication
- **Schema versioning** for compatibility

## Error Handling

- **Typed errors**: ValidationError, NetworkError, TimeoutError
- **Exponential backoff**: 1s → 1.5s → 2.25s → ... (max 30s)
- **Circuit breaker** pattern for sustained failures

## Performance

- **20M+ events/second** throughput with batching
- **Sub-millisecond latency** for API calls
- **<1% CPU overhead**, <10MB memory footprint
- **Zero-copy processing** with object pooling

## Configuration

```go
DefaultEndpoint      = "collect.usercanal.com:50000"
DefaultBatchSize     = 100
DefaultFlushInterval = 10 * time.Second
DefaultMaxRetries    = 3
```

Environment overrides supported for dev/staging/production.

## Extension Points

- **Custom events** via string-based EventName
- **Log routing** via EventType field (collect/enrich)
- **Observability** via GetStats() and debug mode
- **Distributed tracing** via context IDs