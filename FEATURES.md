# Usercanal SDK Features

## Core Architecture Features

### Unified Protocol Support
- **Dual protocol handling**: Analytics events and structured logging in a single SDK
- **Shared transport layer**: Both protocols use the same optimized TCP connection
- **Independent batching**: Separate batchers prevent cross-protocol blocking
- **Unified configuration**: Single config applies to both events and logs

### High-Performance Protocol
- **Binary FlatBuffers**: Zero-copy serialization for maximum throughput (20M+ events/second)
- **Smart batching**: Configurable size (default: 100 items) and time-based (default: 10s) batching
- **Connection pooling**: Single persistent TCP connection with keepalive
- **Schema validation**: Type-safe binary schemas prevent data corruption

### Enterprise-Grade Reliability
- **Built-in authentication**: API key embedded in every batch header
- **Batch tracing**: Unique batch IDs for delivery tracking and debugging
- **Zero data loss**: Automatic retry with exponential backoff
- **DNS failover**: Multiple endpoint support for high availability
- **Graceful degradation**: Continues operation during network issues

## Event Analytics Features

### User Behavior Tracking
- **Standard event types**: Pre-defined constants for common user actions
- **Custom events**: Support for application-specific event tracking
- **User identification**: Associate events with user identities
- **Group analytics**: Track organizational and cohort behavior
- **Session tracking**: Correlate events within user sessions

### Revenue Analytics
- **Multiple revenue types**: Subscription and one-time payment tracking
- **Multi-currency support**: USD, EUR, GBP with extensible currency system
- **Product tracking**: Detailed product information in revenue events
- **Order management**: Complete transaction lifecycle tracking
- **Conversion funnels**: Track user progression through purchase flows

### Type Safety
- **Strongly typed events**: EventName enum prevents typos
- **Validated properties**: Type-safe property maps with validation
- **Currency constraints**: Enum-based currency codes
- **Revenue validation**: Amount and currency validation

## Advanced Logging Features

### Structured Logging
- **Binary format**: Optimized binary logging vs traditional text-based syslog
- **Structured data**: Rich property maps with mixed data types
- **Multiple severity levels**: Standard syslog levels (0-8) from EMERGENCY to TRACE
- **Service isolation**: Clear service and source identification
- **Context correlation**: Distributed tracing via context IDs

### High-Throughput Processing
- **Batch processing**: Log batching not available in traditional syslog
- **Non-blocking I/O**: Asynchronous log delivery
- **Bulk operations**: LogBatch() for high-volume scenarios
- **Memory efficient**: Streaming batch processing

### Enterprise Integration
- **Workspace routing**: Multi-tenant isolation via API keys
- **Real-time processing**: Sub-millisecond routing and processing
- **Event type routing**: Configurable log routing (collect/enrich/auth)
- **Metadata enhancement**: Pipeline for log enrichment

## Transport Layer Features

### Network Reliability
- **Connection management**: Automatic reconnection with exponential backoff
- **Health monitoring**: Continuous connection health checks
- **Timeout handling**: Configurable timeouts with context cancellation
- **Error recovery**: Automatic retry for transient failures
- **Resource cleanup**: Proper connection and resource disposal

### Performance Optimization
- **TCP optimization**: Single persistent connection with keepalive
- **Message framing**: Efficient length-prefixed message protocol
- **Zero-copy processing**: FlatBuffers eliminate serialization overhead
- **Batch efficiency**: Optimal network utilization through batching
- **Memory management**: Bounded memory usage with configurable limits

### Security
- **TLS support**: Production-grade encryption for data in transit
- **API key authentication**: Secure workspace isolation
- **Schema validation**: Prevent malformed data injection
- **Input sanitization**: Comprehensive validation and bounds checking

## Observability Features

### Built-in Metrics
- **Throughput tracking**: Events and logs sent per second
- **Queue monitoring**: Real-time queue depth and processing stats
- **Success/failure rates**: Detailed delivery statistics
- **Timing metrics**: Batch flush intervals and network latency
- **Resource usage**: Memory footprint and connection status

### Health Monitoring
- **Connection status**: TCP connection health indicators
- **Batch delivery**: Recent success/failure tracking
- **Error categorization**: Validation, network, and timeout error tracking
- **Performance indicators**: Latency and throughput measurements

### Debug Support
- **Verbose logging**: Configurable debug output
- **Batch tracing**: Track individual batch delivery
- **Error details**: Comprehensive error messages with context
- **Status dumping**: Runtime state inspection for troubleshooting

## Configuration Features

### Flexible Configuration
- **Environment-specific**: Different configs for dev/staging/production
- **Runtime adjustment**: Hot configuration updates where applicable
- **Validation**: Configuration validation with sensible defaults
- **Option patterns**: Clean, composable configuration API

### Performance Tuning
- **Batch size tuning**: Optimize for throughput vs latency
- **Flush interval adjustment**: Balance real-time vs efficiency
- **Retry configuration**: Customize retry behavior for different environments
- **Connection tuning**: Timeout and keepalive customization

### Development Support
- **Local development**: Easy local collector configuration
- **Debug mode**: Enhanced logging for development
- **Test mode**: Special configuration for testing scenarios
- **Mock support**: Interface design enables easy mocking

## Quality Assurance Features

### Thread Safety
- **Concurrent access**: Thread-safe public API
- **Background processing**: Safe concurrent batch processing
- **Resource synchronization**: Proper locking around shared resources
- **Graceful shutdown**: Coordinated cleanup across all threads

### Error Handling
- **Comprehensive errors**: Detailed error types for different failure modes
- **Error propagation**: Clean error handling through the call stack
- **Recovery mechanisms**: Automatic recovery from transient failures
- **User feedback**: Clear error messages for debugging

### Testing Support
- **Interface design**: Clean interfaces enable comprehensive testing
- **Dependency injection**: Configurable dependencies for testing
- **Mock-friendly**: Easy to mock for unit testing
- **Integration testing**: Support for end-to-end testing scenarios