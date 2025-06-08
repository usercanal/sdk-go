# UserCanal Go SDK - TODO

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