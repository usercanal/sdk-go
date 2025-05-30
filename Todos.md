
## Integrate new log protocl
- Fix bugs
- Better interface naming and add Convenience Logging Methods
    client.TrackEvent(ctx, event)   // Analytics event
    client.SendLog(ctx, entry)      // Structured log
    client.LogInfo(ctx, message)    // Quick logging helpers
    client.LogError(ctx, err)       // Quick logging helpers

1. Simplify interface.
2. Make event names string for both strings and event constants usage.
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
