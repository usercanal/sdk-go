# Contributing to UserCanal Go SDK

This document contains information for SDK developers. For SDK users, please refer to the main README.md.

## Project Structure

### Core SDK Package (`go-sdk/usercanal/`)
* Main package that users import
* Re-exports public types and interfaces
* Provides high-level configuration
* Version information and SDK metadata

### Public API (`go-sdk/api/`)
* `client.go`: Main client implementation
  * Core event tracking methods (Track, Identify, Group, Revenue)
  * Client lifecycle management
  * Stats and monitoring
* `stats.go`: Statistics and monitoring
  * Event tracking metrics
  * Connection status
  * Performance metrics
* Configuration via functional options pattern

### Types and Validation (`go-sdk/types/`)
* `types.go`: Core data structures
  * Event, Identity, Group types
  * Properties and metadata
* `event_types.go`: Event type definitions
  * Standard event names
  * Type-safe event constants
* `validation.go`: Data validation
  * Struct validation rules
  * Property type checking
* `errors.go`: Error definitions
  * Custom error types
  * Error wrapping rules

### Event Batching (`go-sdk/batch/`)
* Efficient event batching
  * Configurable batch size
  * Time-based auto-flushing
  * Memory-efficient queue
* Thread-safe operations
  * Concurrent event adding
  * Safe batch flushing
* Error handling
  * Failed batch retry
  * Event requeuing

### Transport Layer (`go-sdk/transport/`)
* gRPC communication
  * Connection management
  * State monitoring
  * Keep-alive handling
* Retry mechanism
  * Exponential backoff with jitter
  * Configurable retry limits
* Metrics tracking
  * Success/failure counts
  * Latency tracking
  * Batch statistics

### Protocol Conversion (`go-sdk/convert/`)
* Protocol buffer conversion
  * Go types to protobuf
  * Proper type handling
  * Validation during conversion

## Development Guidelines

### Error Handling
1. Use custom error types from `types.ErrorX`
2. Implement proper error wrapping
3. Provide context in error messages
4. Handle context cancellation properly
5. Use appropriate error types for different scenarios:
   * ValidationError for input validation
   * NetworkError for transport issues
   * TimeoutError for deadlines

### Concurrency
1. Use RWMutex for better read performance
2. Implement clean shutdown procedures
3. Proper goroutine management
4. Channel-based communication where appropriate
5. Document concurrency guarantees

### Testing
1. Table-driven tests with good coverage
2. Integration tests for full flow
3. Benchmark tests for performance
4. Concurrency tests with -race
5. Mock external dependencies using interfaces

### Performance Considerations
1. Efficient batching strategies
2. Memory allocation optimization
3. Lock contention minimization
4. Connection pooling
5. Buffer reuse where possible

### Code Style
1. Follow Go standard formatting
2. Use meaningful variable names
3. Add comprehensive comments
4. Document exported symbols
5. Use consistent error handling patterns

## Build Process

### Version Management
\`\`\`bash
# Build with version information
make build VERSION=1.0.0-beta.1
\`\`\`

### Environment Setup
1. Install required tools:
   * Go 1.23 or later
   * Protocol Buffers compiler
   * Make
