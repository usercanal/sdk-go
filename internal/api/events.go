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

	// Use minimal enrichment for server-side (device_id only, no auto session generation)
	transportEvent = c.identityMgr.EnrichEventMinimal(transportEvent)

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

	// Use minimal enrichment for server-side (device_id only, no auto session generation)
	transportEvent = c.identityMgr.EnrichEventMinimal(transportEvent)

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

	// Use minimal enrichment for server-side (device_id only, no auto session generation)
	transportEvent = c.identityMgr.EnrichEventMinimal(transportEvent)

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

	// Use minimal enrichment for server-side (device_id only, no auto session generation)
	transportEvent = c.identityMgr.EnrichEventMinimal(transportEvent)

	if err := c.eventBatcher.Add(ctx, transportEvent); err != nil {
		return fmt.Errorf("failed to add revenue event: %w", err)
	}

	return nil
}

// TrackAdvanced sends an analytics event with advanced options for device/session override
func (c *Client) TrackAdvanced(ctx context.Context, event types.EventAdvanced) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	// Set timestamp if not provided
	timestamp := time.Now()
	if event.Timestamp != nil {
		timestamp = *event.Timestamp
	}

	// Convert to regular Event for transport conversion
	regularEvent := types.Event{
		UserId:     event.UserId,
		Name:       event.Name,
		Properties: event.Properties,
		Timestamp:  timestamp,
	}

	transportEvent, err := convert.EventToInternal(&regularEvent)
	if err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidInput, err)
	}

	// Handle advanced overrides - use minimal enrichment for server-side scenarios
	if event.DeviceID != nil || event.SessionID != nil {
		// Apply manual overrides, use minimal enrichment to avoid auto-session generation
		transportEvent = c.identityMgr.EnrichEventMinimal(transportEvent)

		if event.DeviceID != nil {
			transportEvent.DeviceID = *event.DeviceID
		}
		if event.SessionID != nil {
			transportEvent.SessionID = *event.SessionID
		}
	} else {
		// Use minimal enrichment for server-side (no auto session generation)
		transportEvent = c.identityMgr.EnrichEventMinimal(transportEvent)
	}

	if err := c.eventBatcher.Add(ctx, transportEvent); err != nil {
		return fmt.Errorf("failed to add advanced event: %w", err)
	}

	return nil
}
