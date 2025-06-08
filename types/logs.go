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
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogEmergency:
		return "emergency"
	case LogAlert:
		return "alert"
	case LogCritical:
		return "critical"
	case LogError:
		return "error"
	case LogWarning:
		return "warning"
	case LogNotice:
		return "notice"
	case LogInfo:
		return "info"
	case LogDebug:
		return "debug"
	case LogTrace:
		return "trace"
	default:
		return "unknown"
	}
}

// String returns the string representation of LogEventType
func (t LogEventType) String() string {
	switch t {
	case LogCollect:
		return "collect"
	case LogEnrich:
		return "enrich"
	default:
		return "unknown"
	}
}

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
