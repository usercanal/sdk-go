# Component Behavior Specifications

## Client Facade Component

### Purpose
Provides language-idiomatic public API while abstracting internal complexity.

### State Management
- **Initialization**: Validates API key, applies configuration, initializes internal client
- **Active**: Accepts user calls, delegates to internal client, handles errors gracefully
- **Closing**: Rejects new calls, flushes pending data, coordinates shutdown
- **Closed**: Rejects all calls with appropriate error messages

### Key Behaviors

#### Constructor Behavior
```
Input: api_key (required), config (optional)
Validation: 
  - api_key must be non-empty string
  - config values must be positive where applicable
Output: Client instance or error
Side Effects: Creates internal client, starts background timers
```

#### Method Delegation Pattern
```
1. Check client state (not closed/closing)
2. Validate input parameters
3. Convert to internal types if needed
4. Delegate to internal client method
5. Handle errors and convert to public error types
6. Return result to user
```

#### Error Conversion Strategy
- Internal validation errors → Public validation errors
- Internal network errors → Public network errors with retry hints
- Internal timeout errors → Public timeout errors
- Unknown errors → Generic client errors with debug information

#### Graceful Shutdown
```
1. Mark client as "closing" (reject new calls)
2. Flush all pending data with timeout
3. Close internal client and transport
4. Mark client as "closed"
5. Return flush errors if any occurred
```

## Internal Client Component

### Purpose
Coordinates between batching, transport, and user requests while managing lifecycle.

### State Transitions
```
Created → Initialized → Active → Closing → Closed
```

### Key Behaviors

#### Dual Protocol Coordination
- Maintains separate batch managers for events and logs
- Routes method calls to appropriate batcher
- Aggregates statistics from both batchers
- Coordinates shutdown sequence

#### Configuration Application
```
1. Apply defaults for unspecified values
2. Validate configuration constraints
3. Pass configuration to child components
4. Enable debug logging if requested
```

#### Statistics Aggregation
```
Collect from event batcher:
  - Queue size, success count, failure count
  - Last flush time, last failure time
Collect from log batcher:
  - Same metrics as events
Aggregate into unified stats structure
```

#### Flush Coordination
```
1. Flush event batcher (wait for completion)
2. Flush log batcher (wait for completion)  
3. Return combined error if either failed
4. Maintain operation atomicity per batcher
```

## Batch Manager Component

### Purpose
Aggregates items into batches and coordinates delivery with configurable timing.

### State Management
- **Active**: Accepting items, managing queue, periodic flushing
- **Flushing**: Temporarily blocking new items during send operation
- **Closing**: Final flush attempt, no new items accepted

### Core Algorithms

#### Item Addition Logic
```
function add_item(item):
    if context_cancelled:
        return timeout_error
    
    acquire_write_lock()
    queue.append(item)
    needs_flush = queue.size >= batch_size
    release_write_lock()
    
    if needs_flush:
        return flush()
    return success
```

#### Periodic Flush Logic
```
function periodic_flush():
    while not_closing:
        wait_for_timer_or_shutdown()
        if has_items():
            flush()
```

#### Flush Implementation
```
function flush():
    acquire_write_lock()
    if queue.empty():
        release_write_lock()
        return success
    
    items = queue.copy()
    queue.clear()
    release_write_lock()
    
    result = send_function(items)
    
    if result.failed() and not_context_cancelled():
        acquire_write_lock()
        queue.prepend(items)  # Re-queue failed items
        release_write_lock()
        increment_failure_count(items.length)
        return network_error
    
    increment_success_count(items.length)
    update_last_flush_time()
    return success
```

#### Graceful Shutdown
```
function close():
    stop_periodic_timer()
    flush_with_timeout(5_seconds)
    log_remaining_items_if_any()
```

### Thread Safety Requirements
- Queue modifications must be atomic
- Statistics updates must be atomic
- Flush operations must not interfere with item addition
- Multiple flush calls should be serialized

## Transport Layer Component

### Purpose
Manages TCP connection lifecycle, implements retry logic, and handles binary protocol.

### Connection State Machine
```
Disconnected → Connecting → Connected → Reconnecting → Disconnected
                    ↓           ↓           ↑
                 Failed ←────── Error ──────┘
```

### Key Behaviors

#### Connection Establishment
```
function establish_connection():
    for each_endpoint in failover_list:
        try:
            socket = create_tcp_socket()
            socket.connect(endpoint, timeout=5s)
            if tls_enabled:
                socket = wrap_with_tls(socket)
            return socket
        catch connection_error:
            continue
    return connection_failed_error
```

#### Send Implementation
```
function send_batch(items):
    for attempt in 1..max_retries:
        try:
            if not_connected():
                establish_connection()
            
            serialized = serialize_to_flatbuffer(items)
            socket.write_uint32_be(serialized.length)
            socket.write_bytes(serialized.data)
            socket.flush()
            
            return success
            
        catch network_error as e:
            increment_failure_count()
            if attempt < max_retries:
                delay = calculate_backoff_delay(attempt)
                sleep(delay)
                continue
            return e
```

#### Retry Backoff Calculation
```
function calculate_backoff_delay(attempt):
    base_delay = 1000ms  
    max_delay = 60000ms
    multiplier = 2.0
    jitter_percent = 0.1
    
    delay = min(base_delay * (multiplier ^ (attempt - 1)), max_delay)
    jitter = delay * jitter_percent * random(-1, 1)
    return delay + jitter
```

#### Health Check Implementation
```
function periodic_health_check():
    while connected:
        wait(30_seconds)
        if connection_idle_for(60_seconds):
            send_heartbeat_or_reconnect()
```

#### Connection Recovery
```
function handle_connection_error():
    close_current_connection()
    mark_as_reconnecting()
    apply_exponential_backoff()
    attempt_reconnection()
```

## Binary Protocol Handler

### Purpose
Serializes data structures to FlatBuffers format and manages protocol compliance.

### Serialization Behaviors

#### Event Batch Serialization
```
1. Create FlatBuffer builder
2. Serialize batch header:
   - Generate unique batch_id
   - Include api_key for authentication  
   - Set current timestamp
   - Set batch_type = 1 (events)
3. For each event:
   - Serialize properties map
   - Create event with all fields
   - Add to events vector
4. Create batch with header and events
5. Finalize and return bytes
```

#### Log Batch Serialization
```
1. Create FlatBuffer builder
2. Serialize batch header:
   - Generate unique batch_id
   - Include api_key for authentication
   - Set current timestamp  
   - Set batch_type = 2 (logs)
3. For each log entry:
   - Serialize data map
   - Create log with all fields
   - Add to logs vector
4. Create batch with header and logs
5. Finalize and return bytes
```

#### Property Map Serialization
```
function serialize_properties(properties):
    keys = []
    values = []
    
    for key, value in properties:
        keys.append(serialize_string(key))
        values.append(serialize_value(value))
    
    return create_property_vector(keys, values)
```

#### Value Type Handling
```
function serialize_value(value):
    switch value.type:
        case string: return StringValue(value)
        case int64:  return IntValue(value)  
        case float64: return FloatValue(value)
        case bool:   return BoolValue(value)
        case null:   return NullValue()
        default:     return StringValue(value.toString())
```

### Validation Behaviors

#### Pre-serialization Validation
```
function validate_event(event):
    require_non_empty(event.user_id, max_length=255)
    require_non_empty(event.event_name, max_length=128)
    validate_properties(event.properties, max_count=64)
    validate_timestamp_range(event.timestamp)
```

#### Property Validation
```
function validate_properties(properties, max_count):
    require(properties.size <= max_count)
    
    for key, value in properties:
        require_valid_key(key, max_length=64)
        require_valid_value(value, max_size=1024)
```

#### Batch Size Validation
```
function validate_batch_size(items):
    require(items.length <= 1000)
    
    estimated_size = calculate_serialized_size(items)
    require(estimated_size <= 10MB)
```

These component specifications define the exact behaviors expected from each major component, enabling consistent implementation across different programming languages while maintaining the performance and reliability characteristics of the original Go implementation.