// sdk-go/internal/convert/common.go
package convert

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/usercanal/sdk-go/types"
)

// Common payload marshaling utilities
func marshalPayload(data map[string]interface{}) ([]byte, error) {
	if data == nil {
		data = make(map[string]interface{})
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return payload, nil
}

// Common timestamp handling
func resolveTimestamp(t time.Time) uint64 {
	if t.IsZero() {
		return uint64(time.Now().UnixMilli())
	}
	return uint64(t.UnixMilli())
}

// Common validation helpers
func validateRequired(field, value string) error {
	if value == "" {
		return types.NewValidationError(field, "is required")
	}
	return nil
}

// Validate numeric values
func validatePositive(field string, value float64) error {
	if value <= 0 {
		return types.NewValidationError(field, "must be positive")
	}
	return nil
}

// Validate non-negative values
func validateNonNegative(field string, value float64) error {
	if value < 0 {
		return types.NewValidationError(field, "cannot be negative")
	}
	return nil
}
