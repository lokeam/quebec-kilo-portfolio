package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/types"
)

type SpendTrackingDbAdapter interface {
	GetSpendTrackingBFFResponse(ctx context.Context, userID string) (types.SpendTrackingBFFResponseFINAL, error)
}