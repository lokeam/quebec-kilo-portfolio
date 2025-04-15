package physical

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
)

type PhysicalCacheAdapter struct {
	cacheWrapper interfaces.CacheWrapper
}

func NewPhysicalCacheAdapter(
	cacheWrapper interfaces.CacheWrapper,
) (interfaces.PhysicalCacheWrapper, error) {
	return &PhysicalCacheAdapter{
		cacheWrapper: cacheWrapper,
	}, nil
}

func (pca *PhysicalCacheAdapter) GetCachedPhysicalLocations(
	ctx context.Context,
	userID string,
) ([]models.PhysicalLocation, error) {
	cacheKey := fmt.Sprintf("physical:%s", userID)

	var locations []models.PhysicalLocation
	cacheHit, err := pca.cacheWrapper.GetCachedResults(ctx, cacheKey, &locations)
	if err != nil {
		return nil, err
	}

	if cacheHit {
		return locations, nil
	}

	return nil, nil
}

func (pca *PhysicalCacheAdapter) SetCachedPhysicalLocations(
	ctx context.Context,
	userID string,
	locations []models.PhysicalLocation,
	) error {
		cacheKey := fmt.Sprintf("physical:%s", userID)
		return pca.cacheWrapper.SetCachedResults(ctx, cacheKey, locations)
	}

func (pca *PhysicalCacheAdapter) GetSingleCachedPhysicalLocation(
	ctx context.Context,
	userID string,
	locationID string,
) (*models.PhysicalLocation, bool, error) {
	cacheKey := fmt.Sprintf("physical:%s:location:%s", userID, locationID)

	var location models.PhysicalLocation
	cacheHit, err := pca.cacheWrapper.GetCachedResults(ctx, cacheKey, &location)
	if err != nil {
		return nil, false, err
	}

	if cacheHit {
		return &location, true, nil
	}

	return nil, false, nil
}

func (pca *PhysicalCacheAdapter) SetSingleCachedPhysicalLocation(
	ctx context.Context,
	userID string,
	location models.PhysicalLocation,
) error {
	cacheKey := fmt.Sprintf("physical:%s:location:%s", userID, location.ID)
	return pca.cacheWrapper.SetCachedResults(ctx, cacheKey, location)
}

func (pca *PhysicalCacheAdapter) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	cacheKey := fmt.Sprintf("physical:%s", userID)
	return pca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil)
}

func (pca *PhysicalCacheAdapter) InvalidateLocationCache(
	ctx context.Context,
	userID string,
	locationID string,
) error {
	cacheKey := fmt.Sprintf("physical:%s:location:%s", userID, locationID)
	return pca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil)
}
