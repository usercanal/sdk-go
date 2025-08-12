# UserCanal Go SDK - TODO

Consider builder pattern
```go
// Default: Auto-generate device_id, no session_id
client.Track("purchase_completed").
    WithProperties(props).
    Send(ctx)

// Override device_id (server knows user's device_id)
client.Track("subscription_renewed").
    WithDeviceID(userDeviceID).
    WithProperties(props).
    Send(ctx)

// Override both (rare, but needed for proxy scenarios)
client.Track("page_view").
    WithDeviceID(userDeviceID).
    WithSessionID(sessionID).
    WithProperties(props).
    Send(ctx)

// Explicit no session (server-side)
client.Track("api_call_made").
    WithDeviceID(userDeviceID).
    WithoutSession().
    WithProperties(props).
    Send(ctx)
```

## 📋 **BACKLOG**

### **Testing**
- Add basic smoke tests for Event/EventIdentify/EventRevenue/Logging APIs
- Add configuration pattern tests (default vs custom config)
- Add error handling validation tests
- Add race condition tests for thread safety claims

### **Documentation**
- Add Quick Start guide with examples
- Add API reference documentation
- Create SDK implementation guide for other languages

### **Performance & Quality**
- Add benchmark tests for Event/Log throughput claims
- Add graceful degradation when collector unavailable
- Add EventAdvanced for custom timestamps/event IDs (when customers need it)
- Add statistics implementation completion (ResolvedEndpoints, DNSFailures)
- Add identity manager integration for session enrichment
- Add convert package error handling consistency

### **Developer Experience**
- Add NewClient options pattern (WithDebug, WithBatchSize, etc.)
- Add batch ID exposure for debugging
- Add timeout configuration options
- Add validation for API key format (16 bytes hex)

### **Build & Maintenance**
- Add linting/formatting checks in build pipeline
- Add flatbuffers code generation verification
- Verify minimal external dependencies
- Add version info build-time variables
- Add package documentation comments
- Add missing constants export audit
- Add schema consistency validation
