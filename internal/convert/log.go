// sdk-go/internal/convert/log.go
package convert

import (
	"fmt"

	schema_log "github.com/usercanal/sdk-go/internal/schema/log"
	"github.com/usercanal/sdk-go/internal/transport"
	"github.com/usercanal/sdk-go/types"
)

// Map SDK log levels to FlatBuffer log levels
var logLevelMap = map[types.LogLevel]schema_log.LogLevel{
	types.LogEmergency: schema_log.LogLevelEMERGENCY,
	types.LogAlert:     schema_log.LogLevelALERT,
	types.LogCritical:  schema_log.LogLevelCRITICAL,
	types.LogError:     schema_log.LogLevelERROR,
	types.LogWarning:   schema_log.LogLevelWARNING,
	types.LogNotice:    schema_log.LogLevelNOTICE,
	types.LogInfo:      schema_log.LogLevelINFO,
	types.LogDebug:     schema_log.LogLevelDEBUG,
	types.LogTrace:     schema_log.LogLevelTRACE,
}

// Map SDK log event types to FlatBuffer log event types
var logEventTypeMap = map[types.LogEventType]schema_log.LogEventType{
	types.LogUnknown: schema_log.LogEventTypeUNKNOWN,
	types.LogCollect: schema_log.LogEventTypeCOLLECT,
	types.LogEnrich:  schema_log.LogEventTypeENRICH,
	types.LogAuth:    schema_log.LogEventTypeAUTH,
}

// LogToInternal converts a types.LogEntry to an internal transport.LogEntry
func LogToInternal(l *types.LogEntry) (*transport.Log, error) {
	// Validate required fields
	if err := validateRequired("Service", l.Service); err != nil {
		return nil, err
	}
	if err := validateRequired("Source", l.Source); err != nil {
		return nil, err
	}

	// Map log level with validation
	fbLogLevel, ok := logLevelMap[l.Level]
	if !ok {
		return nil, fmt.Errorf("invalid log level: %d", l.Level)
	}

	// Map log event type with validation
	fbEventType, ok := logEventTypeMap[l.EventType]
	if !ok {
		return nil, fmt.Errorf("invalid log event type: %d", l.EventType)
	}

	// Prepare payload - combine message and data
	payload := make(map[string]interface{})
	if l.Message != "" {
		payload["message"] = l.Message
	}
	if l.Data != nil {
		for k, v := range l.Data {
			payload[k] = v
		}
	}

	payloadBytes, err := marshalPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal log payload: %w", err)
	}

	return &transport.Log{
		EventType: fbEventType,
		ContextID: l.ContextID,
		Level:     fbLogLevel,
		Timestamp: resolveTimestamp(l.Timestamp),
		Source:    l.Source,
		Service:   l.Service,
		Payload:   payloadBytes,
	}, nil
}
