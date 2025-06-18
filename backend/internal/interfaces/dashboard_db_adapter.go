package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/types"
)

type DashboardDbAdapter interface {
	GetDashboardBFFResponse(ctx context.Context, userID string) (types.DashboardBFFResponse,error)
}