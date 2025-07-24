package digital

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

// GetSubscription retrieves a subscription for a digital location
func (da *DigitalDbAdapter) GetSubscription(
	ctx context.Context,
	locationID string,
) (*models.Subscription, error) {
	da.logger.Debug("GetSubscription called", map[string]any{
		"locationID": locationID,
	})

	var subscription models.Subscription
	err := da.db.GetContext(
		ctx,
		&subscription,
		GetSubscriptionByLocationIDQuery,
		locationID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription not found: %w", err)
		}
		return nil, fmt.Errorf("error getting subscription: %w", err)
	}

	da.logger.Debug("GetSubscription success", map[string]any{
		"subscription": subscription,
	})

	return &subscription, nil
}


// CreateSubscription creates a new subscription for a digital location
func (da *DigitalDbAdapter) CreateSubscription(
	ctx context.Context,
	subscription models.Subscription,
) (*models.Subscription, error) {
	da.logger.Debug("CreateSubscription called", map[string]any{
		"subscription": subscription,
	})

	// Validate billing cycle format
	switch subscription.BillingCycle {
	case "1 month", "3 month", "6 month", "12 month":
		// Valid billing cycles
	default:
		return nil, fmt.Errorf("invalid billing cycle: %s. Must be one of: 1 month, 3 month, 6 month, 12 month", subscription.BillingCycle)
	}

	// Validate anchor date is provided
	if subscription.AnchorDate.IsZero() {
		return nil, fmt.Errorf("anchor_date is required for subscription creation")
	}

	now := time.Now()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	err := da.db.QueryRowxContext(
		ctx,
		CreateSubscriptionWithAnchorDateQuery,
		subscription.LocationID,
		subscription.BillingCycle,
		subscription.CostPerCycle,
		subscription.AnchorDate,
		subscription.PaymentMethod,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	).StructScan(&subscription)

	if err != nil {
		return nil, fmt.Errorf("error adding subscription: %w", err)
	}

	da.logger.Debug("CreateSubscription success", map[string]any{
		"subscription": subscription,
	})

	return &subscription, nil
}


// UpdateSubscription updates an existing subscription
func (da *DigitalDbAdapter) UpdateSubscription(
	ctx context.Context,
	subscription models.Subscription,
) error {
	da.logger.Debug("UpdateSubscription called", map[string]any{
		"subscription": subscription,
	})

	// Validate billing cycle format
	switch subscription.BillingCycle {
	case "1 month", "3 month", "6 month", "12 month":
		// Valid billing cycles
	default:
		return fmt.Errorf("invalid billing cycle: %s. Must be one of: 1 month, 3 month, 6 month, 12 month", subscription.BillingCycle)
	}

	// Validate anchor date is provided
	if subscription.AnchorDate.IsZero() {
		return fmt.Errorf("anchor_date is required for subscription updates")
	}

	subscription.UpdatedAt = time.Now()
	result, err := da.db.ExecContext(
		ctx,
		UpdateSubscriptionQuery,
		subscription.BillingCycle,
		subscription.CostPerCycle,
		subscription.AnchorDate,
		subscription.PaymentMethod,
		subscription.UpdatedAt,
		subscription.LocationID,
	)

	if err != nil {
		return fmt.Errorf("error updating subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	da.logger.Debug("UpdateSubscription success", map[string]any{
		"rowsAffected": rowsAffected,
	})

	return nil
}


// DeleteSubscription deletes a subscription for a digital location
func (da *DigitalDbAdapter) DeleteSubscription(
	ctx context.Context,
	locationID string,
) error {
	da.logger.Debug("DeleteSubscription called", map[string]any{
		"locationID": locationID,
	})

	result, err := da.db.ExecContext(
		ctx,
		DeleteSubscriptionQuery,
		locationID,
	)
	if err != nil {
		return fmt.Errorf("error removing subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	da.logger.Debug("DeleteSubscription success", map[string]any{
		"rowsAffected": rowsAffected,
	})

	return nil
}
