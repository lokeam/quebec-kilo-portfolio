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

// CreateOneTimePurchase creates a new one-time purchase
// func (sts *SpendTrackingService) CreateOneTimePurchase(
//     ctx context.Context,
//     userID string,
//     purchase models.OneTimePurchaseToSave,
// ) error {
//     sts.logger.Info("SpendTrackingService - CreateOneTimePurchase called", map[string]any{
//         "userID": userID,
//         "title":  purchase.Title,
//     })

//     // 1. Validate inputs
//     if err := sts.validator.ValidateUserID(userID); err != nil {
//         return fmt.Errorf("invalid user ID: %w", err)
//     }
//     if err := sts.validator.ValidateOneTimePurchase(purchase); err != nil {
//         return fmt.Errorf("invalid purchase: %w", err)
//     }

//     // 2. Add to database
//     if err := sts.dbAdapter.CreateOneTimePurchase(ctx, userID, purchase); err != nil {
//         return fmt.Errorf("failed to create one-time purchase: %w", err)
//     }

//     // 3. Invalidate cache
//     if err := sts.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
//         sts.logger.Error("Failed to invalidate user cache", map[string]any{
//             "error":  err,
//             "userID": userID,
//         })
//     }

//     return nil
// }

// UpdateOneTimePurchase updates an existing one-time purchase
// func (sts *SpendTrackingService) UpdateOneTimePurchase(
//     ctx context.Context,
//     userID string,
//     purchase models.OneTimePurchaseToSave,
// ) error {
//     sts.logger.Info("SpendTrackingService - UpdateOneTimePurchase called", map[string]any{
//         "userID": userID,
//         "title":  purchase.Title,
//     })

//     // 1. Validate inputs
//     if err := sts.validator.ValidateUserID(userID); err != nil {
//         return fmt.Errorf("invalid user ID: %w", err)
//     }
//     if err := sts.validator.ValidateOneTimePurchase(purchase); err != nil {
//         return fmt.Errorf("invalid purchase: %w", err)
//     }

//     // 2. Update in database
//     if err := sts.dbAdapter.UpdateOneTimePurchase(ctx, userID, purchase); err != nil {
//         return fmt.Errorf("failed to update one-time purchase: %w", err)
//     }

//     // 3. Invalidate cache
//     if err := sts.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
//         sts.logger.Error("Failed to invalidate user cache", map[string]any{
//             "error":  err,
//             "userID": userID,
//         })
//     }

//     return nil
// }

// DeleteOneTimePurchase deletes a one-time purchase
// func (sts *SpendTrackingService) DeleteOneTimePurchase(
//     ctx context.Context,
//     userID string,
//     purchaseID int64,
// ) error {
//     sts.logger.Info("SpendTrackingService - DeleteOneTimePurchase called", map[string]any{
//         "userID":     userID,
//         "purchaseID": purchaseID,
//     })

//     // 1. Validate inputs
//     if err := sts.validator.ValidateUserID(userID); err != nil {
//         return fmt.Errorf("invalid user ID: %w", err)
//     }
//     if err := sts.validator.ValidatePurchaseID(purchaseID); err != nil {
//         return fmt.Errorf("invalid purchase ID: %w", err)
//     }

//     // 2. Delete from database
//     if err := sts.dbAdapter.DeleteOneTimePurchase(ctx, userID, purchaseID); err != nil {
//         return fmt.Errorf("failed to delete one-time purchase: %w", err)
//     }

//     // 3. Invalidate cache
//     if err := sts.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
//         sts.logger.Error("Failed to invalidate user cache", map[string]any{
//             "error":  err,
//             "userID": userID,
//         })
//     }

//     return nil
// }

// UpdateSubscription updates a subscription
// func (sts *SpendTrackingService) UpdateSubscription(
//     ctx context.Context,
//     userID string,
//     subscription models.SubscriptionToSave,
// ) error {
//     sts.logger.Info("SpendTrackingService - UpdateSubscription called", map[string]any{
//         "userID": userID,
//         "name":   subscription.Name,
//     })

//     // 1. Validate inputs
//     if err := sts.validator.ValidateUserID(userID); err != nil {
//         return fmt.Errorf("invalid user ID: %w", err)
//     }
//     if err := sts.validator.ValidateSubscription(subscription); err != nil {
//         return fmt.Errorf("invalid subscription: %w", err)
//     }

//     // 2. Update in database
//     if err := sts.dbAdapter.UpdateSubscription(ctx, userID, subscription); err != nil {
//         return fmt.Errorf("failed to update subscription: %w", err)
//     }

//     // 3. Invalidate cache
//     if err := sts.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
//         sts.logger.Error("Failed to invalidate user cache", map[string]any{
//             "error":  err,
//             "userID": userID,
//         })
//     }

//     return nil
// }

// DeleteSubscription deletes a subscription
// func (sts *SpendTrackingService) DeleteSubscription(
//     ctx context.Context,
//     userID string,
//     subscriptionID int64,
// ) error {
//     sts.logger.Info("SpendTrackingService - DeleteSubscription called", map[string]any{
//         "userID":          userID,
//         "subscriptionID": subscriptionID,
//     })

//     // 1. Validate inputs
//     if err := sts.validator.ValidateUserID(userID); err != nil {
//         return fmt.Errorf("invalid user ID: %w", err)
//     }
//     if err := sts.validator.ValidateSubscriptionID(subscriptionID); err != nil {
//         return fmt.Errorf("invalid subscription ID: %w", err)
//     }

//     // 2. Delete from database
//     if err := sts.dbAdapter.DeleteSubscription(ctx, userID, subscriptionID); err != nil {
//         return fmt.Errorf("failed to delete subscription: %w", err)
//     }

//     // 3. Invalidate cache
//     if err := sts.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
//         sts.logger.Error("Failed to invalidate user cache", map[string]any{
//             "error":  err,
//             "userID": userID,
//         })
//     }

//     return nil
// }