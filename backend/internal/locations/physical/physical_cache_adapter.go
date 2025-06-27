package physical

import (
	"context"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type PhysicalCacheAdapter struct {
	cacheWrapper interfaces.CacheWrapper
}

// Constants for cache keys
const (
	physicalLocationsCacheKey = "physical:%s"
	physicalLocationsBFFCacheKey = "physical:bff:%s"
	physicalLocationSingleLocationCacheKey = "physical:%s:location:%s"
	cacheTTL = 1 * time.Hour
)

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
	cacheKey := fmt.Sprintf(physicalLocationsCacheKey, userID)

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
	cacheKey := fmt.Sprintf(physicalLocationsCacheKey, userID)
	return pca.cacheWrapper.SetCachedResults(ctx, cacheKey, locations)
}

// --- BFF ---
func (pca *PhysicalCacheAdapter) GetCachedPhysicalLocationsBFF(
	ctx context.Context,
	userID string,
) (types.LocationsBFFResponse, error) {
	cacheKey := fmt.Sprintf(physicalLocationsBFFCacheKey, userID)

	var response types.LocationsBFFResponse
	cacheHit, err := pca.cacheWrapper.GetCachedResults(ctx, cacheKey, &response)
	if err != nil {
			return types.LocationsBFFResponse{}, err
	}

	if cacheHit {
			return response, nil
	}

	return types.LocationsBFFResponse{}, nil
}

func (pca *PhysicalCacheAdapter) SetCachedPhysicalLocationsBFF(
	ctx context.Context,
	userID string,
	response types.LocationsBFFResponse,
) error {
	cacheKey := fmt.Sprintf(physicalLocationsBFFCacheKey, userID)
	return pca.cacheWrapper.SetCachedResults(ctx, cacheKey, response)
}

// -- BFF --

func (pca *PhysicalCacheAdapter) GetSingleCachedPhysicalLocation(
	ctx context.Context,
	userID string,
	locationID string,
) (*models.PhysicalLocation, bool, error) {
	cacheKey := fmt.Sprintf(physicalLocationSingleLocationCacheKey, userID, locationID)

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
	cacheKey := fmt.Sprintf(physicalLocationSingleLocationCacheKey, userID, location.ID)
	return pca.cacheWrapper.SetCachedResults(ctx, cacheKey, location)
}

func (pca *PhysicalCacheAdapter) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	// Invalidate regular physical locations cache
	cacheKey := fmt.Sprintf(physicalLocationsCacheKey, userID)
	if err := pca.cacheWrapper.DeleteCacheKey(ctx, cacheKey); err != nil {
		return err
	}

	// Also invalidate BFF cache
	bffCacheKey := fmt.Sprintf(physicalLocationsBFFCacheKey, userID)
	return pca.cacheWrapper.DeleteCacheKey(ctx, bffCacheKey)
}

func (pca *PhysicalCacheAdapter) InvalidateLocationCache(
	ctx context.Context,
	userID string,
	locationID string,
) error {
	cacheKey := fmt.Sprintf(physicalLocationSingleLocationCacheKey, userID, locationID)
	return pca.cacheWrapper.DeleteCacheKey(ctx, cacheKey)
}

// InvalidateCache invalidates a specific cache key
func (pca *PhysicalCacheAdapter) InvalidateCache(ctx context.Context, cacheKey string) error {
	return pca.cacheWrapper.InvalidateCache(ctx, cacheKey)
}
