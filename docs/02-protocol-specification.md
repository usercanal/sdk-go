# Protocol Specification: Usercanal Binary Protocol

## Protocol Overview

The Usercanal protocol is a high-performance binary protocol built on FlatBuffers for transmitting analytics events and structured logs over TCP connections. The protocol is designed for zero-copy processing and minimal serialization overhead.

## Connection Specification

### Transport Layer
- **Protocol**: TCP over IPv4/IPv6
- **Default Port**: 50000 (production), configurable for local development
- **Connection Model**: Persistent, long-lived connections with keepalive
- **Compression**: None (FlatBuffers already optimized)
- **Encryption**: TLS 1.2+ for production endpoints

### Connection Lifecycle
1. **Establishment**: TCP handshake + optional TLS negotiation
2. **Authentication**: First batch must include valid API key in header
3. **Keep-alive**: Periodic heartbeat messages during idle periods
4. **Graceful close**: Final flush before connection termination

## Message Format

### Batch Structure
Every message sent is a **batch** containing multiple items of the same type:

```
┌─────────────────────────────────────────────────────────┐
│                    Batch Header                         │
├─────────────────────────────────────────────────────────┤
│                   Item 1 (Event/Log)                   │
├─────────────────────────────────────────────────────────┤
│                   Item 2 (Event/Log)                   │
├─────────────────────────────────────────────────────────┤
│                        ...                              │
├─────────────────────────────────────────────────────────┤
│                   Item N (Event/Log)                   │
└─────────────────────────────────────────────────────────┘
```

### Batch Header Fields
- **batch_id**: uint64 - Unique identifier for delivery tracking
- **api_key**: string - Workspace authentication
- **timestamp**: uint64 - Batch creation time (Unix nanoseconds)
- **item_count**: uint32 - Number of items in batch
- **batch_type**: uint8 - Type discriminator (1=events, 2=logs)
- **schema_version**: uint8 - Protocol version for compatibility

### Wire Protocol
1. **Message Length**: 4-byte big-endian uint32 (message size in bytes)
2. **FlatBuffer Data**: Variable-length binary payload
3. **No framing**: Each message is self-contained

## Event Protocol Specification

### Event Message Structure
- **event_id**: string - Unique identifier (UUID recommended)
- **user_id**: string - User identifier (required)
- **event_name**: string - Event type (predefined constants or custom)
- **properties**: Map<string, Value> - Event metadata
- **timestamp**: uint64 - Event occurrence time (Unix nanoseconds)

### Standard Event Names
- **user_signed_up**: New user registration
- **user_logged_in**: Authentication success
- **feature_used**: Feature interaction
- **order_completed**: Purchase transaction
- **subscription_started**: Recurring payment initiation
- **subscription_changed**: Plan modification
- **subscription_canceled**: Cancellation event
- **cart_viewed**: Shopping cart interaction
- **checkout_started**: Purchase flow initiation
- **checkout_completed**: Purchase finalization

### Property Value Types
Properties support heterogeneous data types:
- **string**: UTF-8 encoded text
- **int64**: Signed 64-bit integer
- **float64**: IEEE 754 double precision
- **bool**: Boolean true/false
- **null**: Explicit null value

### Special Event Types

#### Identity Events
- **user_id**: string - Primary user identifier
- **properties**: Map<string, Value> - User attributes (name, email, etc.)

#### Group Events
- **user_id**: string - User being associated
- **group_id**: string - Group/organization identifier  
- **properties**: Map<string, Value> - Group metadata

#### Revenue Events
- **user_id**: string - Purchasing user
- **order_id**: string - Transaction identifier
- **amount**: float64 - Monetary value
- **currency**: string - ISO 4217 currency code (USD, EUR, GBP)
- **revenue_type**: string - "subscription" or "one_time"
- **products**: Array<Product> - Purchased items
- **properties**: Map<string, Value> - Additional metadata

#### Product Structure
- **product_id**: string - SKU or identifier
- **name**: string - Product name
- **price**: float64 - Unit price
- **quantity**: int32 - Number of units

## Log Protocol Specification

### Log Message Structure
- **level**: uint8 - Severity level (0-8, syslog standard)
- **timestamp**: uint64 - Log creation time (Unix nanoseconds)
- **service**: string - Originating service name
- **source**: string - Source file or component
- **message**: string - Human-readable log message
- **context_id**: uint64 - Distributed tracing correlation ID
- **event_type**: uint32 - Routing classification
- **data**: Map<string, Value> - Structured log data

### Log Levels (syslog RFC 5424)
- **0 (Emergency)**: System unusable
- **1 (Alert)**: Immediate action required
- **2 (Critical)**: Critical conditions
- **3 (Error)**: Error conditions
- **4 (Warning)**: Warning conditions
- **5 (Notice)**: Normal but significant
- **6 (Info)**: Informational messages
- **7 (Debug)**: Debug-level messages
- **8 (Trace)**: Fine-grained debug information

### Event Types (Routing)
- **0 (Unknown)**: Unclassified logs
- **1 (Collect)**: Standard log collection
- **2 (Enrich)**: Logs requiring metadata enhancement
- **3 (Auth)**: Authentication and authorization logs

## Validation Rules

### Event Validation
- **user_id**: Required, non-empty string, max 255 characters
- **event_name**: Required, non-empty string, max 128 characters
- **properties**: Optional, max 64 properties per event
- **property keys**: Max 64 characters, alphanumeric + underscore
- **property values**: Max 1KB per value when serialized

### Log Validation  
- **service**: Required, non-empty string, max 64 characters
- **source**: Required, non-empty string, max 128 characters
- **message**: Required, non-empty string, max 8KB
- **data**: Optional, max 32 properties per log entry
- **property keys**: Max 64 characters, alphanumeric + underscore

### Batch Validation
- **api_key**: Required, valid format (implementation-defined)
- **item_count**: Must match actual items in batch
- **batch_size**: Max 1000 items per batch
- **message_size**: Max 10MB per batch message

## Error Handling

### Protocol Errors
The server may respond with error codes for malformed requests:
- **400**: Invalid batch format or validation failure
- **401**: Authentication failure (invalid API key)
- **413**: Batch too large (exceeds size limits)
- **429**: Rate limiting (too many requests)
- **500**: Server internal error

### Client Behavior
- **Validation errors (400)**: Drop batch, log error, continue
- **Auth errors (401)**: Stop sending, report configuration error
- **Rate limiting (429)**: Exponential backoff, retry batch
- **Server errors (500)**: Exponential backoff, retry batch
- **Network errors**: Exponential backoff, retry batch

### Retry Strategy
- **Initial delay**: 1 second
- **Maximum delay**: 60 seconds  
- **Backoff multiplier**: 2.0
- **Maximum retries**: 3 attempts
- **Jitter**: ±10% random variation

## Schema Versioning

### Compatibility Rules
- **Forward compatible**: Newer clients can send to older servers
- **Backward compatible**: Older clients can send to newer servers
- **Field addition**: New optional fields may be added
- **Field removal**: Deprecated fields maintained for 2 major versions

### Version Negotiation
- Clients include schema_version in batch header
- Servers support multiple schema versions simultaneously
- Unsupported versions return 400 error with supported versions

This protocol specification ensures consistent implementation across all language SDKs while maintaining high performance and reliability.