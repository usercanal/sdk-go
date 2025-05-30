// sdk-go/types/logs.go
package types

import "time"

type LogLevel uint8
type LogEventType uint32

// Log levels (syslog standard)
const (
	LogEmergency LogLevel = 0
	LogAlert     LogLevel = 1
	LogCritical  LogLevel = 2
	LogError     LogLevel = 3
	LogWarning   LogLevel = 4
	LogNotice    LogLevel = 5
	LogInfo      LogLevel = 6
	LogDebug     LogLevel = 7
	LogTrace     LogLevel = 8
)

// Log event types for routing
const (
	LogUnknown LogEventType = 0
	LogCollect LogEventType = 1
	LogEnrich  LogEventType = 2
	LogAuth    LogEventType = 3
)

// LogEntry represents a log entry in the system
type LogEntry struct {
	EventType LogEventType
	ContextID uint64
	Level     LogLevel
	Timestamp time.Time
	Source    string
	Service   string
	Message   string                 // For simple string logs
	Data      map[string]interface{} // For structured data
}
