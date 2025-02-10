1. Simplify interface.
2. Make event names string for both strings and event constants usage.
3. Remove all grpc and protobuf imports
--

// Using predefined event
err := client.Track("test_user", EventFeatureUsed).
    WithProperty("test", true).
    Send()

// Using custom event
err := client.Track("test_user", "video.viewed").
    WithProperties(
        "duration", 120,
        "quality", "hd",
    ).
    Send()

// With context and timestamp
err := client.Track("test_user", EventOrderCompleted, ctx).
    WithProperties(
        "revenue", 99.99,
        "currency", "USD",
    ).
    WithTimestamp(time.Now()).
    Send()

3. Context ..
For the context properties, we could add clear documentation:
```go
// WithContext adds additional context to the event
// Available context properties:
//   - "ip": string          - Source IP address
//   - "user_agent": string  - User agent string
//   - "locale": string      - User's locale
//   - "app_version": string - Application version
//   - "platform": string    - Platform identifier
func (t *TrackBuilder) WithContext(keyValues ...interface{}) *TrackBuilder {
    // Similar implementation to WithProperties but with validation
    // for allowed context properties
    return t
}
```
