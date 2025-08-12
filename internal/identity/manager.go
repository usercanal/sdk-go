// identity/manager.go
package identity

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/usercanal/sdk-go/internal/logger"
	"github.com/usercanal/sdk-go/internal/transport"
)

type Manager struct {
	distinctID []byte // 16-byte UUID
	contextID  []byte // 16-byte UUID for session tracking
	userID     []byte // 16-byte UUID or custom ID
	startTime  time.Time
	mu         sync.RWMutex
}

// uuidToBytes converts a UUID to a byte slice
func uuidToBytes(u uuid.UUID) []byte {
	b := make([]byte, 16)
	copy(b, u[:])
	return b
}

func NewManager() (*Manager, error) {
	distinctID := uuidToBytes(uuid.New())
	contextID := uuidToBytes(uuid.New())

	mgr := &Manager{
		distinctID: distinctID,
		contextID:  contextID,
		startTime:  time.Now(),
	}

	logger.Debug("Identity manager initialized with distinctID: %x", distinctID)
	return mgr, nil
}

// EnrichEventMinimal does NOT auto-generate any IDs for server-side events
// Both device_id and session_id remain nil unless explicitly set by developer
func (m *Manager) EnrichEventMinimal(event *transport.Event) *transport.Event {
	if event == nil {
		return nil
	}

	// Do NOT auto-generate device_id or session_id for server SDKs
	// Server events should have nil IDs unless explicitly provided

	return event
}

// GetIdentity returns the current identity state
func (m *Manager) GetIdentity() (distinctID, userID, contextID []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	distinctID = make([]byte, len(m.distinctID))
	copy(distinctID, m.distinctID)

	if len(m.userID) > 0 {
		userID = make([]byte, len(m.userID))
		copy(userID, m.userID)
	}

	contextID = make([]byte, len(m.contextID))
	copy(contextID, m.contextID)

	return
}

// GetSessionDuration returns the current session duration
func (m *Manager) GetSessionDuration() time.Duration {
	return time.Since(m.startTime)
}

// GenerateEventID creates a new UUID for event tracking
func (m *Manager) GenerateEventID() []byte {
	return uuidToBytes(uuid.New())
}

// SetUserID allows manual setting of user ID
func (m *Manager) SetUserID(id []byte) {
	if len(id) == 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.userID = make([]byte, len(id))
	copy(m.userID, id)
}

// Reset clears the user ID but maintains the distinct ID
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.userID = nil
	m.contextID = uuidToBytes(uuid.New())
	m.startTime = time.Now()
}
