// sdk-go/internal/transport/event.go
package transport

import (
	"context"
	"fmt"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	schema_common "github.com/usercanal/sdk-go/internal/schema/common"
	event_collector "github.com/usercanal/sdk-go/internal/schema/event"
	"github.com/usercanal/sdk-go/types"
)

func (s *Sender) SendEvents(ctx context.Context, events []*Event) error {
	if len(events) == 0 {
		return nil
	}

	// Add default timeout if none exists
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	// Size and count validation for critical environments
	if len(events) > MaxBatchItems {
		return types.NewValidationError("events", fmt.Sprintf("batch too large (max %d events)", MaxBatchItems))
	}

	totalSize := 0
	for i, evt := range events {
		// Validate required fields
		if evt.Timestamp == 0 {
			return types.NewValidationError("Timestamp", fmt.Sprintf("event[%d] timestamp is required", i))
		}
		if len(evt.UserID) == 0 {
			return types.NewValidationError("UserID", fmt.Sprintf("event[%d] userID is required", i))
		}
		if len(evt.Payload) == 0 {
			return types.NewValidationError("Payload", fmt.Sprintf("event[%d] payload is required", i))
		}

		// Size validation
		if len(evt.Payload) > MaxEventSize {
			return types.NewValidationError("payload", fmt.Sprintf("event[%d] payload too large (max %d bytes)", i, MaxEventSize))
		}
		totalSize += len(evt.Payload)
	}

	if totalSize > MaxBatchSize {
		return types.NewValidationError("batch", fmt.Sprintf("total payload size %d exceeds limit %d", totalSize, MaxBatchSize))
	}

	select {
	case <-s.ctx.Done():
		return types.NewValidationError("sender", "is shutting down")
	default:
	}

	builder := flatbuffers.NewBuilder(1024 * len(events))

	// Create events vector
	eventOffsets := make([]flatbuffers.UOffsetT, len(events))
	for i := len(events) - 1; i >= 0; i-- {
		evt := events[i]

		payloadOffset := builder.CreateByteVector(evt.Payload)
		userIDOffset := builder.CreateByteVector(evt.UserID)

		event_collector.EventStart(builder)
		event_collector.EventAddTimestamp(builder, evt.Timestamp)
		event_collector.EventAddEventType(builder, evt.EventType)
		event_collector.EventAddUserId(builder, userIDOffset)
		event_collector.EventAddPayload(builder, payloadOffset)
		eventOffsets[i] = event_collector.EventEnd(builder)
	}

	eventsVec := builder.CreateVectorOfTables(eventOffsets)

	// Create EventData
	event_collector.EventDataStart(builder)
	event_collector.EventDataAddEvents(builder, eventsVec)
	eventDataEnd := event_collector.EventDataEnd(builder)

	builder.Finish(eventDataEnd)
	eventDataBytes := builder.FinishedBytes()

	// Send as batch
	err := s.sendBatch(ctx, schema_common.SchemaTypeEVENT, eventDataBytes)
	if err == nil {
		s.recordEventSuccess(len(events))
	}
	return err
}
