package digital

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type DigitalCacheAdapter struct {
	cacheWrapper interfaces.CacheWrapper
}

func NewDigitalCacheAdapter(
	cacheWrapper interfaces.CacheWrapper,
) (interfaces.DigitalCacheWrapper, error) {
	return &DigitalCacheAdapter{
		cacheWrapper: cacheWrapper,
	}, nil
}

func (dca *DigitalCacheAdapter) GetCachedDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocation, error) {
	cacheKey := fmt.Sprintf("digital:%s", userID)

	var locations []models.DigitalLocation
	cacheHit, err := dca.cacheWrapper.GetCachedResults(ctx, cacheKey, &locations)
	if err != nil {
		return nil, err
	}

	if cacheHit {
		return locations, nil
	}

	return nil, nil
}

func (dca *DigitalCacheAdapter) SetCachedDigitalLocations(
	ctx context.Context,
	userID string,
	locations []models.DigitalLocation,
) error {
	cacheKey := fmt.Sprintf("digital:%s", userID)
	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, locations)
}

func (dca *DigitalCacheAdapter) GetSingleCachedDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
) (*models.DigitalLocation, bool, error) {
	cacheKey := fmt.Sprintf("digital:%s:location:%s", userID, locationID)

	var location models.DigitalLocation
	cacheHit, err := dca.cacheWrapper.GetCachedResults(ctx, cacheKey, &location)
	if err != nil {
		return nil, false, err
	}

	if cacheHit {
		return &location, true, nil
	}

	return nil, false, nil
}

func (dca *DigitalCacheAdapter) SetSingleCachedDigitalLocation(
	ctx context.Context,
	userID string,
	location models.DigitalLocation,
) error {
	cacheKey := fmt.Sprintf("digital:%s:location:%s", userID, location.ID)

	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, location)
}

func (dca *DigitalCacheAdapter) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	cacheKey := fmt.Sprintf("digital:%s", userID)
	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil)
}

func (dca *DigitalCacheAdapter) InvalidateDigitalLocationCache(
	ctx context.Context,
	userID string,
	locationID string,
) error {
	// Invalidate the specific location cache
	cacheKey := fmt.Sprintf("digital:%s:location:%s", userID, locationID)
	if err := dca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil); err != nil {
		return fmt.Errorf("failed to invalidate cache for location %s (user: %s, key: %s): %w",
			locationID, userID, cacheKey, err)
	}

	// Also invalidate the user's locations list to ensure consistency
	userCacheKey := fmt.Sprintf("digital:%s", userID)
	if err := dca.cacheWrapper.SetCachedResults(ctx, userCacheKey, nil); err != nil {
		return fmt.Errorf("failed to invalidate user cache (user: %s, key: %s): %w",
			userID, userCacheKey, err)
	}

	// Log successful cache invalidation
	fmt.Printf("Successfully invalidated cache for location %s and user %s\n", locationID, userID)
	return nil
}

func (dca *DigitalCacheAdapter) GetCachedSubscription(
	ctx context.Context,
	locationID string,
) (*models.Subscription, bool, error) {
	cacheKey := fmt.Sprintf("digital:subscription:%s", locationID)

	var subscription models.Subscription
	cacheHit, err := dca.cacheWrapper.GetCachedResults(ctx, cacheKey, &subscription)
	if err != nil {
		return nil, false, err
	}

	if cacheHit {
		return &subscription, true, nil
	}

	return nil, false, nil
}

func (dca *DigitalCacheAdapter) SetCachedSubscription(
	ctx context.Context,
	locationID string,
	subscription models.Subscription,
) error {
	cacheKey := fmt.Sprintf("digital:subscription:%s", locationID)
	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, subscription)
}

func (dca *DigitalCacheAdapter) InvalidateSubscriptionCache(
	ctx context.Context,
	locationID string,
) error {
	cacheKey := fmt.Sprintf("digital:subscription:%s", locationID)
	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil)
}

func (dca *DigitalCacheAdapter) GetCachedPayments(
	ctx context.Context,
	locationID string,
) ([]models.Payment, error) {
	cacheKey := fmt.Sprintf("digital:payments:%s", locationID)

	var payments []models.Payment
	cacheHit, err := dca.cacheWrapper.GetCachedResults(ctx, cacheKey, &payments)
	if err != nil {
		return nil, err
	}

	if cacheHit {
		return payments, nil
	}

	return nil, nil
}

func (dca *DigitalCacheAdapter) SetCachedPayments(
	ctx context.Context,
	locationID string,
	payments []models.Payment,
) error {
	cacheKey := fmt.Sprintf("digital:payments:%s", locationID)
	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, payments)
}

func (dca *DigitalCacheAdapter) InvalidatePaymentsCache(
	ctx context.Context,
	locationID string,
) error {
	cacheKey := fmt.Sprintf("digital:payments:%s", locationID)
	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil)
}


// -- BFF --
func (dca *DigitalCacheAdapter) GetCachedDigitalLocationsBFF(
	ctx context.Context,
	userID string,
) (types.DigitalLocationsBFFResponse, error) {
	cacheKey := fmt.Sprintf("digital:bff:%s", userID)

	var response types.DigitalLocationsBFFResponse
	cacheHit, err := dca.cacheWrapper.GetCachedResults(ctx, cacheKey, &response)
	if err != nil {
		return types.DigitalLocationsBFFResponse{}, err
	}

	if cacheHit {
		return response, nil
	}

	return types.DigitalLocationsBFFResponse{}, nil
}

func (dca *DigitalCacheAdapter) SetCachedDigitalLocationsBFF(
	ctx context.Context,
	userID string,
	response types.DigitalLocationsBFFResponse,
) error {
	cacheKey := fmt.Sprintf("digital:bff:%s", userID)
	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, response)
}