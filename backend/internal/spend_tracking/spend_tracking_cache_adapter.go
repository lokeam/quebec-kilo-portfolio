package spend_tracking

import (
	"context"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/types"
)

type SpendTrackingCacheAdapter struct {
    cacheWrapper interfaces.CacheWrapper
}

// Constants for cache keys
const (
    // spendTrackingBFFCacheKey is the pattern for spend tracking BFF cache keys
    spendTrackingBFFCacheKey = "spend_tracking:bff:%s"

    // spendTrackingLastUpdateKey is the pattern for tracking last update timestamps
    spendTrackingLastUpdateKey = "spend_tracking:last_update:%s"

    // cacheTTL is the time-to-live for cached responses
    cacheTTL = 1 * time.Hour
)

func NewSpendTrackingCacheAdapter(cacheWrapper interfaces.CacheWrapper) (interfaces.SpendTrackingCacheWrapper, error) {
    if cacheWrapper == nil {
        return nil, fmt.Errorf("cacheWrapper is required")
    }

    return &SpendTrackingCacheAdapter{
        cacheWrapper: cacheWrapper,
    }, nil
}

// getLastUpdateTimestamp retrieves the last update timestamp for a user's spend tracking data
func (sta *SpendTrackingCacheAdapter) getLastUpdateTimestamp(
    ctx context.Context,
    userID string,
) (time.Time, error) {
    cacheKey := fmt.Sprintf(spendTrackingLastUpdateKey, userID)

    var timestamp time.Time
    cacheHit, err := sta.cacheWrapper.GetCachedResults(ctx, cacheKey, &timestamp)
    if err != nil {
        return time.Time{}, fmt.Errorf("failed to get last update timestamp: %w", err)
    }

    if !cacheHit {
        return time.Time{}, nil
    }

    return timestamp, nil
}

// updateLastUpdateTimestamp updates the last update timestamp for a user's spend tracking data
func (sta *SpendTrackingCacheAdapter) updateLastUpdateTimestamp(
    ctx context.Context,
    userID string,
) error {
    cacheKey := fmt.Sprintf(spendTrackingLastUpdateKey, userID)
    return sta.cacheWrapper.SetCachedResults(ctx, cacheKey, time.Now())
}

func (sta *SpendTrackingCacheAdapter) GetCachedSpendTrackingBFF(
    ctx context.Context,
    userID string,
) (types.SpendTrackingBFFResponseFINAL, error) {
    // 1. Check last update timestamp
    lastUpdate, err := sta.getLastUpdateTimestamp(ctx, userID)
    if err != nil {
        return types.SpendTrackingBFFResponseFINAL{}, err
    }

    // 2. If timestamp is recent, try to get cached data
    if !lastUpdate.IsZero() && time.Since(lastUpdate) < cacheTTL {
        cacheKey := fmt.Sprintf(spendTrackingBFFCacheKey, userID)
        var response types.SpendTrackingBFFResponseFINAL

        cacheHit, err := sta.cacheWrapper.GetCachedResults(ctx, cacheKey, &response)
        if err != nil {
            return types.SpendTrackingBFFResponseFINAL{}, err
        }

        if cacheHit {
            return response, nil
        }
    }

    // 3. If no cache or cache is stale, return empty response to trigger fresh fetch
    return types.SpendTrackingBFFResponseFINAL{}, nil
}

func (sta *SpendTrackingCacheAdapter) SetCachedSpendTrackingBFF(
    ctx context.Context,
    userID string,
    response types.SpendTrackingBFFResponseFINAL,
) error {
    // 1. Update the last update timestamp
    if err := sta.updateLastUpdateTimestamp(ctx, userID); err != nil {
        return fmt.Errorf("failed to update last update timestamp: %w", err)
    }

    // 2. Cache the response
    cacheKey := fmt.Sprintf(spendTrackingBFFCacheKey, userID)
    return sta.cacheWrapper.SetCachedResults(ctx, cacheKey, response)
}

func (sta *SpendTrackingCacheAdapter) InvalidateUserCache(
    ctx context.Context,
    userID string,
) error {
    // Invalidate BFF cache
    bffCacheKey := fmt.Sprintf(spendTrackingBFFCacheKey, userID)
    if err := sta.cacheWrapper.DeleteCacheKey(ctx, bffCacheKey); err != nil {
        return fmt.Errorf("failed to invalidate BFF cache: %w", err)
    }

    // Also invalidate last update timestamp to ensure consistency
    lastUpdateKey := fmt.Sprintf(spendTrackingLastUpdateKey, userID)
    if err := sta.cacheWrapper.DeleteCacheKey(ctx, lastUpdateKey); err != nil {
        return fmt.Errorf("failed to invalidate last update timestamp: %w", err)
    }

    return nil
}

func (sta *SpendTrackingCacheAdapter) InvalidateSpendTrackingBFFCache(
    ctx context.Context,
    userID string,
) error {
    // Invalidate ONLY the BFF cache
    bffCacheKey := fmt.Sprintf(spendTrackingBFFCacheKey, userID)
    if err := sta.cacheWrapper.DeleteCacheKey(ctx, bffCacheKey); err != nil {
        return fmt.Errorf("failed to invalidate BFF cache: %w", err)
    }

    return nil
}