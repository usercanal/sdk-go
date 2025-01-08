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
   * Go 1.21 or later
   * Protocol Buffers compiler
   * Make

### Build Commands
\`\`\`bash
make build      # Build for current platform
make install    # Install locally
make release    # Build for all platforms
\`\`\`

## Release Process

### Versioning
1. Follow semantic versioning (MAJOR.MINOR.PATCH)
2. Beta releases use -beta.N suffix
3. Release candidates use -rc.N suffix

### Release Steps
1. Update version in Makefile
2. Update CHANGELOG.md
3. Run full test suite
4. Tag release in git
5. Build release binaries
6. Update documentation

### Documentation
1. Update README.md
2. Update API documentation
3. Update examples
4. Check godoc formatting

## Best Practices

### SDK Usage
1. Always provide contexts for cancellation
2. Handle errors appropriately
3. Configure batch sizes based on usage
4. Monitor SDK health via stats
5. Implement proper shutdown

### Production Deployment
1. Set appropriate timeouts
2. Configure retry policies
3. Monitor event queue size
4. Handle backpressure
5. Implement proper logging

## Support

### Getting Help
1. Open an issue for bugs
2. Use discussions for questions
3. Check existing issues first
4. Provide minimal reproduction

### Contributing
1. Fork the repository
2. Create a feature branch
3. Add tests for changes
4. Update documentation
5. Submit pull request