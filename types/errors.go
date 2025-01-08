// types/errors.go
package types

import "fmt"

// Common error types
var (
	ErrInvalidInput   = fmt.Errorf("invalid input")
	ErrNetworkFailure = fmt.Errorf("network failure")
	ErrTimeout        = fmt.Errorf("operation timed out")
	ErrNotConnected   = fmt.Errorf("not connected")
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// Is implements error matching for ValidationError
func (e *ValidationError) Is(target error) bool {
	return target == ErrInvalidInput
}

// NetworkError represents a network-related error
type NetworkError struct {
	Operation string
	Message   string
	Retries   int
}

func (e *NetworkError) Error() string {
	if e.Retries > 0 {
		return fmt.Sprintf("%s failed after %d retries: %s", e.Operation, e.Retries, e.Message)
	}
	return fmt.Sprintf("%s failed: %s", e.Operation, e.Message)
}

// Is implements error matching for NetworkError
func (e *NetworkError) Is(target error) bool {
	return target == ErrNetworkFailure
}

// TimeoutError represents a timeout error
type TimeoutError struct {
	Operation string
	Duration  string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("%s timed out after %s", e.Operation, e.Duration)
}

// Is implements error matching for TimeoutError
func (e *TimeoutError) Is(target error) bool {
	return target == ErrTimeout
}
