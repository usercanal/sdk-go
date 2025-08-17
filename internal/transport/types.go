// sdk-go/internal/transport/types.go
package transport

import (
	event_schema "github.com/usercanal/sdk-go/internal/schema/event"
	log_schema "github.com/usercanal/sdk-go/internal/schema/log"
)

// Event represents an internal event structure for transport
type Event struct {
	Timestamp uint64
	EventType event_schema.EventType
	EventName string
	DeviceID  []byte
	SessionID []byte
	Payload   []byte
}

// Log represents an internal log structure for transport
type Log struct {
	EventType log_schema.LogEventType
	SessionID []byte // 16-byte session UUID for correlation
	Level     log_schema.LogLevel
	Timestamp uint64
	Source    string
	Service   string
	Payload   []byte
}
