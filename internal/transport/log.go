// sdk-go/internal/transport/log.go
package transport

import (
	"context"
	"fmt"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	schema_common "github.com/usercanal/sdk-go/internal/schema/common"
	schema_log "github.com/usercanal/sdk-go/internal/schema/log"
	"github.com/usercanal/sdk-go/types"
)

func (s *Sender) SendLogs(ctx context.Context, logs []*Log) error {
	if len(logs) == 0 {
		return nil
	}

	// Add default timeout if none exists
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	// Size and count validation for critical environments
	if len(logs) > MaxBatchItems {
		return types.NewValidationError("logs", fmt.Sprintf("batch too large (max %d logs)", MaxBatchItems))
	}

	totalSize := 0
	for i, log := range logs {
		// Validate required fields
		if log.Timestamp == 0 {
			return types.NewValidationError("Timestamp", fmt.Sprintf("log[%d] timestamp is required", i))
		}
		if log.Source == "" {
			return types.NewValidationError("Source", fmt.Sprintf("log[%d] source is required", i))
		}
		if log.Service == "" {
			return types.NewValidationError("Service", fmt.Sprintf("log[%d] service is required", i))
		}
		if len(log.Payload) == 0 {
			return types.NewValidationError("Payload", fmt.Sprintf("log[%d] payload is required", i))
		}

		// Size validation
		if len(log.Payload) > MaxLogSize {
			return types.NewValidationError("payload", fmt.Sprintf("log[%d] payload too large (max %d bytes)", i, MaxLogSize))
		}
		totalSize += len(log.Payload)
	}

	if totalSize > MaxBatchSize {
		return types.NewValidationError("batch", fmt.Sprintf("total payload size %d exceeds limit %d", totalSize, MaxBatchSize))
	}

	select {
	case <-s.ctx.Done():
		return types.NewValidationError("sender", "is shutting down")
	default:
	}

	builder := flatbuffers.NewBuilder(1024 * len(logs))

	// Create logs vector
	logOffsets := make([]flatbuffers.UOffsetT, len(logs))
	for i := len(logs) - 1; i >= 0; i-- {
		log := logs[i]

		payloadOffset := builder.CreateByteVector(log.Payload)
		sourceOffset := builder.CreateString(log.Source)
		serviceOffset := builder.CreateString(log.Service)

		schema_log.LogEntryStart(builder)
		schema_log.LogEntryAddEventType(builder, log.EventType)
		schema_log.LogEntryAddContextId(builder, log.ContextID)
		schema_log.LogEntryAddLevel(builder, log.Level)
		schema_log.LogEntryAddTimestamp(builder, log.Timestamp)
		schema_log.LogEntryAddSource(builder, sourceOffset)
		schema_log.LogEntryAddService(builder, serviceOffset)
		schema_log.LogEntryAddPayload(builder, payloadOffset)
		logOffsets[i] = schema_log.LogEntryEnd(builder)
	}

	logsVec := builder.CreateVectorOfTables(logOffsets)

	// Create LogData
	schema_log.LogDataStart(builder)
	schema_log.LogDataAddLogs(builder, logsVec)
	logDataEnd := schema_log.LogDataEnd(builder)

	builder.Finish(logDataEnd)
	logDataBytes := builder.FinishedBytes()

	// logger.Debug("About to send batch with SchemaTypeLOG = %d", int(schema_common.SchemaTypeLOG))

	// Send as batch
	err := s.sendBatch(ctx, schema_common.SchemaTypeLOG, logDataBytes)
	if err == nil {
		s.recordLogSuccess(len(logs))
	}
	return err
}
