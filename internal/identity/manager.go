// identity/manager.go
package identity

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/usercanal/sdk-go/logger"
	pb "github.com/usercanal/sdk-go/proto"
)

type Manager struct {
	distinctID string
	contextID  string
	userID     string
	startTime  time.Time
	mu         sync.RWMutex
}

func NewManager() (*Manager, error) {
	distinctID := uuid.New().String()

	mgr := &Manager{
		distinctID: distinctID,
		contextID:  uuid.New().String(),
		startTime:  time.Now(),
	}

	logger.Debug("Identity manager initialized with distinctID: %s", distinctID)
	return mgr, nil
}

func (m *Manager) EnrichEvent(event *pb.Event) *pb.Event {
	if event == nil {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var base *pb.MessageBase
	switch e := event.Event.(type) {
	case *pb.Event_Track:
		base = e.Track.Base
	case *pb.Event_Identify:
		base = e.Identify.Base
	case *pb.Event_Group:
		base = e.Group.Base
	case *pb.Event_Alias:
		base = e.Alias.Base
	}

	if base == nil {
		base = &pb.MessageBase{}
	}

	// Add any missing required fields
	if base.DistinctId == "" {
		base.DistinctId = m.distinctID
	}
	if base.UserId == "" {
		base.UserId = m.userID
	}
	if base.ContextId == "" {
		base.ContextId = m.contextID
	}
	if base.MessageId == "" {
		base.MessageId = uuid.New().String()
	}
	if base.Timestamp == nil {
		base.Timestamp = timestamppb.Now()
	}

	// Update the base in the appropriate event type
	switch e := event.Event.(type) {
	case *pb.Event_Track:
		e.Track.Base = base
	case *pb.Event_Identify:
		e.Identify.Base = base
	case *pb.Event_Group:
		e.Group.Base = base
	case *pb.Event_Alias:
		e.Alias.Base = base
	}

	return event
}

func (m *Manager) EnrichIdentity(event *pb.Event) *pb.Event {
	if event == nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if e, ok := event.Event.(*pb.Event_Identify); ok {
		if e.Identify != nil && e.Identify.Base != nil && e.Identify.Base.UserId != "" {
			m.userID = e.Identify.Base.UserId
			logger.Debug("Updated user ID to: %s", m.userID)
		}
	}

	return m.EnrichEvent(event)
}

func (m *Manager) EnrichGroup(event *pb.Event) *pb.Event {
	if event == nil {
		return nil
	}

	if e, ok := event.Event.(*pb.Event_Group); ok {
		if e.Group != nil && e.Group.Base == nil {
			e.Group.Base = &pb.MessageBase{}
		}
		if e.Group != nil && e.Group.Base != nil && e.Group.Base.UserId == "" {
			m.mu.RLock()
			e.Group.Base.UserId = m.userID
			m.mu.RUnlock()
		}
	}

	return m.EnrichEvent(event)
}

func (m *Manager) GetIdentity() (string, string, string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.distinctID, m.userID, m.contextID
}

func (m *Manager) GetSessionDuration() time.Duration {
	return time.Since(m.startTime)
}
