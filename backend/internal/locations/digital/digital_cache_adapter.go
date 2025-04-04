package digital

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
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
	cacheKey := fmt.Sprintf("digital:%s:location:%s", userID, locationID)
	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil)
}
