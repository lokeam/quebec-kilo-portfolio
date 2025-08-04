package sublocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
)

type SublocationCacheAdapter struct {
	cacheWrapper interfaces.CacheWrapper
}

func NewSublocationCacheAdapter(
	cacheWrapper interfaces.CacheWrapper,
) (interfaces.SublocationCacheWrapper, error) {
	return &SublocationCacheAdapter{
		cacheWrapper: cacheWrapper,
	}, nil
}

func (sca *SublocationCacheAdapter) GetCachedSublocations(
	ctx context.Context,
	userID string,
) ([]models.Sublocation, error) {
	cacheKey := fmt.Sprintf("sublocation:%s", userID)

	var sublocations []models.Sublocation
	cacheHit, err := sca.cacheWrapper.GetCachedResults(ctx, cacheKey, &sublocations)
	if err != nil {
		return nil, err
	}

	if cacheHit {
		return sublocations, nil
	}

	return nil, nil
}

func (sca *SublocationCacheAdapter) SetCachedSublocations(
	ctx context.Context,
	userID string,
	sublocations []models.Sublocation,
) error {
	cacheKey := fmt.Sprintf("sublocation:%s", userID)
	return sca.cacheWrapper.SetCachedResults(ctx, cacheKey, sublocations)
}

func (sca *SublocationCacheAdapter) GetSingleCachedSublocation(
	ctx context.Context,
	userID string,
	sublocationID string,
) (*models.Sublocation, bool, error) {
	cacheKey := fmt.Sprintf("sublocation:%s:sublocation:%s", userID, sublocationID)

	var sublocation models.Sublocation
	cacheHit, err := sca.cacheWrapper.GetCachedResults(ctx, cacheKey, &sublocation)
	if err != nil {
		return nil, false, err
	}

	if cacheHit {
		return &sublocation, true, nil
	}

	return nil, false, nil
}

func (sca *SublocationCacheAdapter) SetSingleCachedSublocation(
	ctx context.Context,
	userID string,
	sublocation models.Sublocation,
) error {
	// Validate that the sublocation belongs to the user
	if sublocation.UserID != userID {
		return errors.New("sublocation does not belong to user")
	}

	cacheKey := fmt.Sprintf("sublocation:%s:sublocation:%s", userID, sublocation.ID)
	return sca.cacheWrapper.SetCachedResults(ctx, cacheKey, sublocation)
}

func (sca *SublocationCacheAdapter) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	cacheKey := fmt.Sprintf("sublocation:%s", userID)
	return sca.cacheWrapper.DeleteCacheKey(ctx, cacheKey)
}

func (sca *SublocationCacheAdapter) InvalidateSublocationCache(
	ctx context.Context,
	userID string,
	sublocationID string,
) error {
	cacheKey := fmt.Sprintf("sublocation:%s:sublocation:%s", userID, sublocationID)
	return sca.cacheWrapper.DeleteCacheKey(ctx, cacheKey)
}

func (sca *SublocationCacheAdapter) InvalidateLocationCache(
	ctx context.Context,
	userID string,
	locationID string,
) error {
	// Delete both the specific location key and the user's locations collection
	physicalLocationKey := fmt.Sprintf("physical:%s:location:%s", userID, locationID)
	physicalLocationsKey := fmt.Sprintf("physical:%s", userID)
	physicalBFFKey := fmt.Sprintf("physical:bff:%s", userID)

	// Delete specific location cache
	if err := sca.cacheWrapper.DeleteCacheKey(ctx, physicalLocationKey); err != nil {
		return err
	}

	// Also delete the collection of physical locations for this user
	// This ensures that when GetAllPhysicalLocations is called, it will fetch fresh data
	if err := sca.cacheWrapper.DeleteCacheKey(ctx, physicalLocationsKey); err != nil {
		return err
	}

	// Delete the BFF cache to ensure fresh data for the frontend
	return sca.cacheWrapper.DeleteCacheKey(ctx, physicalBFFKey)
}
