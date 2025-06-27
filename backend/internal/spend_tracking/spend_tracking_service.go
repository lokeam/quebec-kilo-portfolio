package spend_tracking

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type SpendTrackingService struct {
    dbAdapter               interfaces.SpendTrackingDbAdapter
    cacheWrapper            interfaces.SpendTrackingCacheWrapper
    dashboardCacheWrapper   interfaces.DashboardCacheWrapper
    validator               interfaces.SpendTrackingValidator
    logger                  interfaces.Logger
}

func NewSpendTrackingService(
    appContext *appcontext.AppContext,
    dbAdapter interfaces.SpendTrackingDbAdapter,
    cacheWrapper interfaces.SpendTrackingCacheWrapper,
    dashboardCacheWrapper interfaces.DashboardCacheWrapper,
) (*SpendTrackingService, error) {
    if dbAdapter == nil {
        return nil, fmt.Errorf("dbAdapter is required")
    }
    if cacheWrapper == nil {
        return nil, fmt.Errorf("cacheWrapper is required")
    }
    if dashboardCacheWrapper == nil {
        return nil, fmt.Errorf("dashboardCacheWrapper is required")
    }

    return &SpendTrackingService{
        dbAdapter:    dbAdapter,
        cacheWrapper: cacheWrapper,
        dashboardCacheWrapper: dashboardCacheWrapper,
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


// --- WRITE OPERATIONS ---
func (sts *SpendTrackingService) CreateOneTimePurchase(
    ctx context.Context,
    userID string,
    request types.SpendTrackingRequest,
) (models.SpendTrackingOneTimePurchaseDB, error) {
    sts.logger.Debug("Creating one-time purchase", map[string]any{
        "userID":  userID,
        "request": request,
    })

    // Validate userID
    if err := sts.validator.ValidateUserID(userID); err != nil {
        return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("invalid user ID: %w", err)
    }

    // Validate the request
    if err := sts.validator.ValidateOneTimePurchase(request); err != nil {
        return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("validation failed: %w", err)
    }

    // Transform request to database model
    oneTimePurchase, err := TransformCreateRequestToModel(request, userID)
    if err != nil {
        return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("transformation failed: %w", err)
    }

    // Create in database
    createdPurchase, err := sts.dbAdapter.CreateOneTimePurchase(ctx, userID, oneTimePurchase)
    if err != nil {
        sts.logger.Error("Failed to create one-time purchase in DB", map[string]any{"error": err})
        return models.SpendTrackingOneTimePurchaseDB{}, err
    }

    // Invalidate the cache for this user
    if err := sts.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
        sts.logger.Error("Failed to invalidate user cache after creating purchase", map[string]any{
            "error":  err,
            "userID": userID,
        })
        // DB update successful, continue despite error
    }

    // Invalidate dashboard cache to refresh statistics
    if err := sts.dashboardCacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
        sts.logger.Error("Failed to invalidate dashboard cache after creating purchase", map[string]any{
            "error":  err,
            "userID": userID,
        })
        // DB update successful, continue despite error
    }

    // Invalidate spend tracking cache to refresh financial data
    if err := sts.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
        sts.logger.Error("Failed to invalidate spend tracking cache after creating purchase", map[string]any{
            "error":  err,
            "userID": userID,
        })
        // DB update successful, continue despite error
    }

    sts.logger.Debug("CreateOneTimePurchase success", map[string]any{
        "oneTimePurchase": createdPurchase,
    })

    return createdPurchase, nil
}