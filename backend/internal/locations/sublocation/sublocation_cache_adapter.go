package sublocation

import (
	"context"
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
	cacheKey := fmt.Sprintf("sublocation:%s:sublocation:%s", userID, sublocation.ID)
	return sca.cacheWrapper.SetCachedResults(ctx, cacheKey, sublocation)
}

func (sca *SublocationCacheAdapter) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	cacheKey := fmt.Sprintf("sublocation:%s", userID)
	return sca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil)
}

func (sca *SublocationCacheAdapter) InvalidateSublocationCache(
	ctx context.Context,
	userID string,
	sublocationID string,
) error {
	cacheKey := fmt.Sprintf("sublocation:%s:sublocation:%s", userID, sublocationID)
	return sca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil)
}
