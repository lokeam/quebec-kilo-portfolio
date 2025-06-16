package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/types"
)

type SpendTrackingCacheWrapper interface {
	GetCachedSpendTrackingBFF(ctx context.Context, userID string) (types.SpendTrackingBFFResponseFINAL, error)

	SetCachedSpendTrackingBFF(ctx context.Context, userID string, response types.SpendTrackingBFFResponseFINAL) error

	InvalidateUserCache(ctx context.Context, userID string) error
}