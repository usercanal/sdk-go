// sdk-go/internal/api/events.go
package api

import (
	"context"
	"fmt"
	"time"

	"github.com/usercanal/sdk-go/internal/convert"
	"github.com/usercanal/sdk-go/types"
)

// Track sends an analytics event
func (c *Client) Track(ctx context.Context, event types.Event) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	// Set timestamp if not set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	transportEvent, err := convert.EventToInternal(&event)
	if err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	if err := c.eventBatcher.Add(ctx, transportEvent); err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}

	return nil
}

// Identify associates a user with their traits
func (c *Client) Identify(ctx context.Context, identity types.Identity) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := identity.Validate(); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	transportEvent, err := convert.IdentityToInternal(&identity)
	if err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	if err := c.eventBatcher.Add(ctx, transportEvent); err != nil {
		return fmt.Errorf("failed to add identity event: %w", err)
	}

	return nil
}

// Group associates a user with a group
func (c *Client) Group(ctx context.Context, groupInfo types.GroupInfo) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := groupInfo.Validate(); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	transportEvent, err := convert.GroupToInternal(&groupInfo)
	if err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	if err := c.eventBatcher.Add(ctx, transportEvent); err != nil {
		return fmt.Errorf("failed to add group event: %w", err)
	}

	return nil
}

// Revenue tracks a revenue event
func (c *Client) Revenue(ctx context.Context, rev types.Revenue) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := rev.Validate(); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	transportEvent, err := convert.RevenueToInternal(&rev)
	if err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	if err := c.eventBatcher.Add(ctx, transportEvent); err != nil {
		return fmt.Errorf("failed to add revenue event: %w", err)
	}

	return nil
}
