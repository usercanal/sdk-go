# SDK Maturity TODOs - Organized by Domain

## 🏗️ **API Domain**

### #2 **Missing Key Public Methods** [interface, log]
**Problem**: Several important methods mentioned in checklist are missing

**Add to `usercanal.go`**:
```go
// Missing logging convenience methods for different levels
func (c *Client) LogWarning(ctx context.Context, service, source, message string, data map[string]interface{}) error {
    return c.internal.LogWarning(ctx, service, source, message, data)
}
// Add LogCritical, LogAlert, LogEmergency, LogNotice, LogTrace
```

### #8 **Missing API Methods in Internal Client** [log]
**Fix in `internal/api/logs.go`**:
```go
func (c *Client) LogWarning(ctx context.Context, service, source, message string, data map[string]interface{}) error {
    return c.Log(ctx, types.LogEntry{
        Level:     types.LogWarning,
        EventType: types.LogCollect,
        Service:   service,
        Source:    source,
        Message:   message,
        Data:      data,
    })
}
```

### #3 **Incomplete Statistics Implementation** [api]
**Problem**: `GetStats()` has commented out functionality
```go
// In transport/sender.go - add missing methods
func (s *Sender) GetConnectionInfo() ConnectionInfo {
    return ConnectionInfo{
        State: s.connMgr.GetState(),
        Uptime: s.Uptime(),
        ReconnectCount: s.connMgr.GetReconnectCount(),
    }
}
```

### #17 **Incomplete Close() Chain** [api]
**Problem**: Close() doesn't properly cascade through all components
```go
func (c *Client) Close() error {
    // Missing: c.eventBatcher.Close()
    // Missing: c.logBatcher.Close()
}
```

### **TODO-ORIGINAL: Interface Simplification** [interface, event, log]
**Problem**: Current interface could be more user-friendly
```go
// Desired interface:
err := client.Track("test_user", EventFeatureUsed).
    WithProperty("test", true).
    Send()

err := client.Track("test_user", "video.viewed").
    WithProperties("duration", 120, "quality", "hd").
    Send()

// Better method names:
client.TrackEvent(ctx, event)   // Analytics event
client.SendLog(ctx, entry)      // Structured log
client.LogInfo(ctx, message)    // Quick logging helpers
```

## 🔄 **Batch Domain**

### #6 **Resource Management Issues** [batch]
**Problem**: Potential goroutine leaks and incomplete cleanup

**Fix in `batch/batch.go`**:
```go
func (m *Manager) Close() error {
    m.ticker.Stop()
    close(m.done)
    m.wg.Wait() // Add WaitGroup tracking

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    return m.Flush(ctx)
}
```

### #40 **Missing Graceful Degradation** [batch, transport]
**Problem**: No fallback behavior when collector is completely unavailable
```go
// Should have configurable options:
// - Drop data silently (default)
// - Store to local file
// - Circuit breaker pattern
// Current implementation keeps retrying forever
```

## 🔄 **Convert Domain**

### #5 **Missing Context Handling** [convert, log]
**Problem**: Context IDs for logging are not properly implemented
```go
func LogToInternal(l *types.LogEntry) (*transport.Log, error) {
    contextID := l.ContextID
    if contextID == 0 {
        contextID = generateContextID() // Generate if not provided
    }
    return &transport.Log{ContextID: contextID, ...}, nil
}
```

### #23 **Convert Package Error Handling** [convert]
**Problem**: Inconsistent error wrapping in convert functions
```go
// Should be consistent about error types
// Sometimes returns validation errors, sometimes generic errors
```

### #38 **Transport Event/Log Type Conversion** [convert, event]
**Problem**: Missing validation for EventType mapping
```go
eventType, ok := eventTypeMap[e.Name]
if !ok {
    return nil, fmt.Errorf("unmapped event type: %s", e.Name)
    // Should this fail or use a default?
}
```

## 👤 **Identity Domain**

### #12 **Identity Manager Not Integrated** [identity]
**Problem**: Complete identity manager exists but is never used
```go
// In internal/api/client.go
type Client struct {
    // ... existing fields
    identityMgr *identity.Manager  // Add this
}
// Use it to enrich events with session info
```

## 📝 **Logger Domain**

### #18 **Debug Logging Inconsistencies** [logger]
**Problem**: Some components use logger, others use fmt/log directly
```go
// Standardize on internal logger package everywhere
```

## 📋 **Schema Domain**

### #9 **Schema Path Issues** [schema]
**Problem**: Schema imports might cause issues across SDKs
```go
// Make sure all SDKs use identical .fbs files
// Consider putting schemas in a shared repository
```

### #15 **Flatbuffers Builder Size Issues** [schema]
**Problem**: Fixed builder sizes might be too small for large batches
```go
// Dynamic size calculation
estimatedSize := estimateBatchSize(events)
builder := flatbuffers.NewBuilder(estimatedSize)
```

### #20 **Schema Version Handling Missing** [schema]
**Problem**: FlatBuffers schemas have version fields but they're not used
```go
// Should set/check schema_version for compatibility
```

### #33 **Schema Consistency Check** [schema]
**Problem**: No verification that generated Go code matches .fbs files
```go
// Build step that regenerates and diffs
```

## 🚀 **Transport Domain**

### #10 **Connection Manager Signal Method** [transport]
**Problem**: `signalRetry()` is not exported but called from other packages
```go
// Either export it or handle internally in sender
func (cm *ConnManager) SignalRetry() { ... }
```

### #13 **Incomplete Metrics Tracking** [transport]
**Problem**: Transport metrics don't track all the data that Stats expects
```go
type Stats struct {
    ResolvedEndpoints []string      // Not implemented
    LastDNSResolution time.Time     // Not implemented
    DNSFailures       int64         // Not implemented
}
```

### #19 **Missing Timeout Configurations** [transport]
**Problem**: Hardcoded timeouts that should be configurable
```go
// Make timeouts configurable instead of hardcoded
```

### #21 **API Key Format Validation** [transport]
**Problem**: API key validation assumes hex format but doesn't validate length
```go
// Should validate it's exactly 16 bytes (32 hex chars)
```

### #22 **Missing Batch ID Usage** [transport]
**Problem**: Batch IDs are generated but never exposed for tracking
```go
// Users need access to batch IDs for debugging
```

## 📅 **Version Domain**

### #25 **Version Info Incomplete** [version]
**Problem**: Version info has build-time variables that default to "unknown"
```go
// version/version.go ldflags variables need better defaults
```

## 📚 **Examples Domain**

### #24 **Missing Examples Integration** [examples]
**Problem**: Examples directory exists but examples might not reflect actual API
```go
// Verify examples work with current API
// Add error case examples
```

## 🎯 **Types Domain**

### #4 **Error Handling Inconsistencies** [types]
**Problem**: Mixed error wrapping patterns
```go
// Create error constructors in types/
func NewNetworkError(operation, message string) error {
    return &NetworkError{Operation: operation, Message: message}
}
```

### #7 **Missing Validation Methods** [types]
**Problem**: Some validation is incomplete
```go
func validateEventName(name EventName) error {
    if len(string(name)) > 128 {
        return NewValidationError("EventName", "exceeds maximum length of 128")
    }
    return nil
}
```

### #11 **Type System Inconsistencies** [types]
**Problem**: Missing string conversion methods
```go
// Add String() methods for LogLevel and LogEventType
func (l LogLevel) String() string { ... }
func (t LogEventType) String() string { ... }
```

### #14 **Missing Constants Export** [types]
**Problem**: Important constants are not exported in the public API
```go
// AuthMethodGoogle, AuthMethodEmail, PaymentMethodCard
// Should be accessible to users
```

### #35 **Error Message Consistency** [types]
**Problem**: Error messages use different formatting patterns
```go
// Standardize error message format across all validation
```

### #39 **Timestamp Precision Inconsistency** [types]
**Problem**: Mixed timestamp precision handling
```go
// Standardize on milliseconds everywhere
return uint64(time.Now().UnixMilli())
```

### **TODO-ORIGINAL: Event Names as Strings** [types, event]
**Problem**: Make event names more flexible
```go
// Allow both string and EventName constants
// EventName should be string type, not custom type
```

## 🏗️ **Interface Domain**

### #30 **Missing Package Documentation** [interface]
**Problem**: No package-level documentation comments
```go
// Package usercanal provides a unified SDK for analytics events and structured logging.
package usercanal
```

### #31 **Internal Package Exposure** [interface]
**Problem**: Some internal types might be exposed unintentionally
```go
// Audit all exported types in internal/ packages
```

## ⚙️ **Configuration Domain**

### #1 **Configuration Inconsistencies** [config]
**Problem**: Multiple default endpoints and inconsistent config handling
```go
// Centralize in internal/config/defaults.go
const DefaultEndpoint = "collect.usercanal.com:50000"
```

## 🏗️ **Build System Domain**

### #26 **Makefile and Build System** [build]
**Problem**: Missing verification for flatbuffers code generation
```go
// Add linting/formatting checks in build pipeline
```

### #27 **Module Dependencies** [build]
**Problem**: Should verify minimal required dependencies
```go
// Review external dependencies for necessity
```

### #32 **Missing Benchmark Tests** [build]
**Problem**: No performance benchmarks to validate claims
```go
// BenchmarkTrackEvent, BenchmarkLogEntry, BenchmarkBatchSerialization
```

### #34 **Missing Race Condition Tests** [build]
**Problem**: Claims thread safety but no race tests
```go
// Tests with -race flag and concurrent usage tests
```

### #37 **Missing go.mod Version Requirements** [build]
**Problem**: Ensure reproducible builds across environments

## 📄 **Documentation Domain**

### #16 **Missing Thread Safety Documentation** [docs]
**Problem**: Some components claim thread safety but don't clearly document which methods are safe

### #28 **File Extension Issue** [docs]
**Problem**: `SDK-implementation.sdk` should be `.md`

### #29 **CONTRIBUTION.md vs Contributing** [docs]
**Problem**: Should be `CONTRIBUTING.md` (standard)

---

## 🔥 **Critical Priority (Fix First)**
- #1 Configuration Inconsistencies
- #2 Missing Key Public Methods
- #6 Resource Management Issues
- #17 Incomplete Close() Chain
- #40 Missing Graceful Degradation

## 📋 **High Priority (Cross-Language Impact)**
- #4 Error Handling Inconsistencies
- #9 Schema Path Issues
- #11 Type System Inconsistencies
- #20 Schema Version Handling

## 🔧 **Medium Priority (Quality & Completeness)**
- #3 Incomplete Statistics Implementation
- #5 Missing Context Handling
- #12 Identity Manager Integration
- #13 Incomplete Metrics Tracking
