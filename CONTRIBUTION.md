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

### Transport Layer (`internal/transport/`)
* TCP communication
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

### Event Serialization (`schema/` & `internal/event/`)
* FlatBuffers schema definitions
* Efficient binary serialization
* Type-safe event handling

### Protocol Conversion (`go-sdk/convert/`)
* FlatBuffers conversion
  * Go types to FlatBuffers
  * Efficient zero-copy conversion
  * Proper type handling and validation
  * Memory pooling for builders

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

### Environment Setup
1. Install required tools:
   * Go 1.23 or later
   * Protocol Buffers compiler
   * Make
