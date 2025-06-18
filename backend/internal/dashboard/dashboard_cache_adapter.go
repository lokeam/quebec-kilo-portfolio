package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/types"
)

type DashboardCacheAdapter struct {
	cacheWrapper interfaces.CacheWrapper
}

// Constants for cache keys
const (
	dashboardBFFCacheKey = "dashboard:bff:%s"
	dashboardLastUpdateKey = "dashboard:last_update:%s"
	cacheTTL = 5 * time.Minute
)

func NewDashboardCacheAdapter(cacheWrapper interfaces.CacheWrapper) (interfaces.DashboardCacheWrapper, error) {
	if cacheWrapper == nil {
		return nil, fmt.Errorf("cacheWrapper is required")
	}

	return &DashboardCacheAdapter{
		cacheWrapper: cacheWrapper,
	}, nil
}

// Helper fns:
func (dca *DashboardCacheAdapter) getLastUpdateTimestamp(
	ctx context.Context,
	userID string,
) (time.Time, error) {
	cacheKey := fmt.Sprintf(dashboardLastUpdateKey, userID)

	var timestamp time.Time
	cacheHit, err := dca.cacheWrapper.GetCachedResults(ctx, cacheKey, &timestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get last update timestamp: %w", err)
	}

	if !cacheHit {
		return time.Time{}, nil
	}

	return timestamp, nil
}

func (dca *DashboardCacheAdapter) updateLastUpdateTimestamp(
	ctx context.Context,
	userID string,
) error {
	cacheKey := fmt.Sprintf(dashboardLastUpdateKey, userID)
	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, time.Now())
}


func (dca *DashboardCacheAdapter) GetCachedDashboardBFF(
	ctx context.Context,
	userID string) (types.DashboardBFFResponse, error) {
    // 1. Check last update timestamp
    lastUpdate, err := dca.getLastUpdateTimestamp(ctx, userID)
    if err != nil {
        return types.DashboardBFFResponse{}, err
    }

    // 2. If timestamp is recent, try to get cached data
    if !lastUpdate.IsZero() && time.Since(lastUpdate) < cacheTTL {
        cacheKey := fmt.Sprintf(dashboardBFFCacheKey, userID)
        var response types.DashboardBFFResponse

        cacheHit, err := dca.cacheWrapper.GetCachedResults(ctx, cacheKey, &response)
        if err != nil {
            return types.DashboardBFFResponse{}, err
        }

        if cacheHit {
            return response, nil
        }
    }

    // 3. If no cache or cache is stale, return empty response to trigger fresh fetch
    return types.DashboardBFFResponse{}, nil
}

func (dca *DashboardCacheAdapter) SetCachedDashboardBFF(
	ctx context.Context,
	userID string,
	response types.DashboardBFFResponse,
) error {
	// Update the last update timestamp
	if err := dca.updateLastUpdateTimestamp(ctx, userID); err != nil {
		return fmt.Errorf("failed to update last update timestamp: %w", err)
	}

	// Cache the response
	cacheKey := fmt.Sprintf(dashboardBFFCacheKey, userID)
	return dca.cacheWrapper.SetCachedResults(ctx, cacheKey, response)
}

func (dca *DashboardCacheAdapter) InvalidateUserCache(ctx context.Context, userID string) error {
	// Delete the dashboard BFF response cache
	bffCacheKey := fmt.Sprintf(dashboardBFFCacheKey, userID)
	if err := dca.cacheWrapper.DeleteCacheKey(ctx, bffCacheKey); err != nil {
		return fmt.Errorf("failed to delete dashboard bff cache: %w", err)
	}

	// Delete the last update timestamp cache
	lastUpdateKey := fmt.Sprintf(dashboardLastUpdateKey, userID)
	if err := dca.cacheWrapper.DeleteCacheKey(ctx, lastUpdateKey); err != nil {
		return fmt.Errorf("failed to delete last update timestamp cache: %w", err)
	}

	return nil
}