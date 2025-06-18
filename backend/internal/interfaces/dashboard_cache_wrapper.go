package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/types"
)

type DashboardCacheWrapper interface {
	GetCachedDashboardBFF(ctx context.Context, userID string) (types.DashboardBFFResponse, error)

	SetCachedDashboardBFF(ctx context.Context, userID string, response types.DashboardBFFResponse) error

	InvalidateUserCache(ctx context.Context, userID string) error
}