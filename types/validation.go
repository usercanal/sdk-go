// types/validation.go
package types

import (
	"fmt"
	"time"

	"github.com/usercanal/sdk-go/internal/logger"
)

// Validate validates an Event
func (e *Event) Validate() error {
	if e.UserId == "" {
		return NewValidationError("UserId", "is required")
	}
	if e.Name == "" {
		return NewValidationError("Name", "is required")
	}
	if !e.Name.IsStandardEvent() {
		// Allow custom events, but log a warning
		logger.Warn("Non-standard event name used: %s", e.Name)
	}
	if err := validateProperties(e.Properties); err != nil {
		return fmt.Errorf("properties validation failed: %w", err)
	}
	if e.Timestamp.IsZero() {
		return NewValidationError("Timestamp", "is required")
	}
	return nil
}

// Validate validates an Identity
func (i *Identity) Validate() error {
	if i.UserId == "" {
		return NewValidationError("UserId", "is required")
	}
	if err := validateProperties(i.Properties); err != nil {
		return fmt.Errorf("properties validation failed: %w", err)
	}
	return nil
}

// Validate validates a GroupInfo
func (g *GroupInfo) Validate() error {
	if g.UserId == "" {
		return NewValidationError("UserId", "is required")
	}
	if g.GroupId == "" {
		return NewValidationError("GroupId", "is required")
	}
	if err := validateProperties(g.Properties); err != nil {
		return fmt.Errorf("properties validation failed: %w", err)
	}
	return nil
}

// Validate validates a Revenue
func (r *Revenue) Validate() error {
	if r.OrderID == "" {
		return NewValidationError("OrderID", "is required")
	}
	if r.Amount <= 0 {
		return NewValidationError("Amount", "must be positive")
	}
	if string(r.Currency) == "" {
		return NewValidationError("Currency", "is required")
	}
	if string(r.Type) == "" {
		return NewValidationError("Type", "is required")
	}

	for i, p := range r.Products {
		if err := p.Validate(); err != nil {
			return fmt.Errorf("product[%d] validation failed: %w", i, err)
		}
	}

	if err := validateProperties(r.Properties); err != nil {
		return fmt.Errorf("properties validation failed: %w", err)
	}

	return nil
}

// Validate validates a Product
func (p *Product) Validate() error {
	if p.ID == "" {
		return NewValidationError("ID", "is required")
	}
	if p.Price < 0 {
		return NewValidationError("Price", "cannot be negative")
	}
	if p.Quantity <= 0 {
		return NewValidationError("Quantity", "must be positive")
	}
	return nil
}

// validateProperties checks if properties contain valid values
func validateProperties(props Properties) error {
	if props == nil {
		return nil
	}

	for key, value := range props {
		if key == "" {
			return NewValidationError("PropertyKey", "cannot be empty")
		}

		if err := validatePropertyValue(value); err != nil {
			return fmt.Errorf("property '%s' validation failed: %w", key, err)
		}
	}

	return nil
}

// validatePropertyValue checks if a property value is of a supported type
func validatePropertyValue(value interface{}) error {
	switch v := value.(type) {
	case nil:
		return nil
	case string:
		return nil
	case int, int32, int64, float32, float64:
		return nil
	case bool:
		return nil
	case time.Time:
		return nil
	case EventName:
		return nil
	case Currency:
		return nil
	case RevenueType:
		return nil
	case AuthMethod:
		return nil
	case PaymentMethod:
		return nil
	case []interface{}:
		for i, item := range v {
			if err := validatePropertyValue(item); err != nil {
				return fmt.Errorf("array item[%d] validation failed: %w", i, err)
			}
		}
		return nil
	case map[string]interface{}:
		return validateProperties(Properties(v))
	default:
		return NewValidationError("PropertyValue", fmt.Sprintf("unsupported type: %T", value))
	}
}
