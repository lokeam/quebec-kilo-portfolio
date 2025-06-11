package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/types"
)

type SearchAddGameFormDbAdapter interface {
	GetAllGameStorageLocationsBFF(ctx context.Context, userID string) (types.AddGameFormStorageLocationsResponse, error)
}