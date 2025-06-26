package spend_tracking

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/types"
)

type SpendTrackingService struct {
    dbAdapter    interfaces.SpendTrackingDbAdapter
    cacheWrapper interfaces.SpendTrackingCacheWrapper
    validator    interfaces.SpendTrackingValidator
    logger       interfaces.Logger
}

func NewSpendTrackingService(
    appContext *appcontext.AppContext,
    dbAdapter interfaces.SpendTrackingDbAdapter,
    cacheWrapper interfaces.SpendTrackingCacheWrapper,
) (*SpendTrackingService, error) {
    if dbAdapter == nil {
        return nil, fmt.Errorf("dbAdapter is required")
    }
    if cacheWrapper == nil {
        return nil, fmt.Errorf("cacheWrapper is required")
    }

    return &SpendTrackingService{
        dbAdapter:    dbAdapter,
        cacheWrapper: cacheWrapper,
        validator:    NewSpendTrackingValidator(),
        logger:       appContext.Logger,
    }, nil
}

// GetSpendTrackingBFFResponse retrieves the spend tracking BFF response for a user
func (sts *SpendTrackingService) GetSpendTrackingBFFResponse(
    ctx context.Context,
    userID string,
) (types.SpendTrackingBFFResponseFINAL, error) {
    // 1. Validate userID
    if err := sts.validator.ValidateUserID(userID); err != nil {
        return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("invalid user ID: %w", err)
    }

    // 2. Try to get from cache first
    cachedResponse, err := sts.cacheWrapper.GetCachedSpendTrackingBFF(ctx, userID)
    if err == nil && len(cachedResponse.CurrentTotalThisMonth) > 0 {
        sts.logger.Debug("Cache hit for spend tracking BFF response", map[string]any{
            "userID": userID,
        })
        return cachedResponse, nil
    }

    // 3. Cache miss, get from database
    sts.logger.Debug("Cache miss for spend tracking BFF response, fetching from database", map[string]any{
        "userID": userID,
    })
    response, err := sts.dbAdapter.GetSpendTrackingBFFResponse(ctx, userID)
    if err != nil {
        return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("failed to get spend tracking BFF response: %w", err)
    }

    // 4. Cache the response
    if err := sts.cacheWrapper.SetCachedSpendTrackingBFF(ctx, userID, response); err != nil {
        sts.logger.Error("Failed to cache spend tracking BFF response", map[string]any{
            "error":  err,
            "userID": userID,
        })
    }

    return response, nil
}
