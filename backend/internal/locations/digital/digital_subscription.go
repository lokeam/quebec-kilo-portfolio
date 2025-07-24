package digital

import (
	"context"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)


func HandleCreateSubscription(
	ctx context.Context,
	dbAdapter interfaces.DigitalDbAdapter,
	location models.DigitalLocation,
	request types.DigitalLocationRequest,
) error {
	if request.IsSubscription && location.Subscription != nil {
		// The subscription is already transformed and attached to the location
		// Just set the LocationID and timestamps
		location.Subscription.LocationID = location.ID
		location.Subscription.CreatedAt = location.CreatedAt
		location.Subscription.UpdatedAt = location.UpdatedAt

		// Create the subscription in the database
		_, err := dbAdapter.CreateSubscription(ctx, *location.Subscription)
		if err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}
	}

	return nil
}


func HandleUpdateSubscription(
	ctx context.Context,
	dbAdapter interfaces.DigitalDbAdapter,
	existingLocation models.DigitalLocation,
	location models.DigitalLocation,
) error {
	locationID := location.ID

	if location.IsSubscription && location.Subscription != nil {
		// Update existing subscription with new data
		subscription := *location.Subscription
		subscription.LocationID = locationID
		subscription.UpdatedAt = time.Now()

		if err := dbAdapter.UpdateSubscription(ctx, subscription); err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}
	} else if location.IsSubscription && existingLocation.Subscription != nil && location.Subscription == nil {
		// Remove subscription if service type changed to non-subscription
		if err := dbAdapter.DeleteSubscription(ctx, locationID); err != nil {
			return fmt.Errorf("failed to remove subscription: %w", err)
		}
	} else if !location.IsSubscription && existingLocation.Subscription != nil {
		// Remove subscription if service type changed
		if err := dbAdapter.DeleteSubscription(ctx, locationID); err != nil {
			return fmt.Errorf("failed to remove subscription: %w", err)
		}
	}

	return nil
}