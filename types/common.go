// sdk-go/types/common.go
package types

import "fmt"

// Properties represents a map of property values
type Properties map[string]interface{}

// Common error types
var (
	ErrInvalidInput   = fmt.Errorf("invalid input")
	ErrNetworkFailure = fmt.Errorf("network failure")
	ErrTimeout        = fmt.Errorf("operation timed out")
	ErrNotConnected   = fmt.Errorf("not connected")
)

// Error constructors for consistent error handling patterns

// NewNetworkError creates a standardized network error
func NewNetworkError(operation, message string) error {
	return &NetworkError{
		Operation: operation,
		Message:   message,
		Retries:   0,
	}
}

// NewTimeoutError creates a standardized timeout error
func NewTimeoutError(operation, duration string) error {
	return &TimeoutError{
		Operation: operation,
		Duration:  duration,
	}
}

// WrapError wraps an error with context for consistent error chains
func WrapError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", operation, err)
}

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
