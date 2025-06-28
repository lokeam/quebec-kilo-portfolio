package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type SpendTrackingDbAdapter interface {
	GetSpendTrackingBFFResponse(ctx context.Context, userID string) (types.SpendTrackingBFFResponseFINAL, error)
	CreateOneTimePurchase(ctx context.Context, userID string, request models.SpendTrackingOneTimePurchaseDB) (models.SpendTrackingOneTimePurchaseDB, error)
	UpdateOneTimePurchase(ctx context.Context, userID string, request models.SpendTrackingOneTimePurchaseDB) (models.SpendTrackingOneTimePurchaseDB, error)
	GetSingleSpendTrackingItem(ctx context.Context, userID string, itemID string) (models.SpendTrackingOneTimePurchaseDB, error)
	DeleteSpendTrackingItems(ctx context.Context, userID string, itemIDs []string) (int64, error)
}