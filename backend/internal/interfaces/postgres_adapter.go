package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/types"
)

type PostgresAdapter interface {
	GetUserLibraryItems(ctx context.Context, userID string) ([]types.Game, error)
	AddGameToLibrary(ctx context.Context, userID string, gameID int64) error
	RemoveGameFromLibrary(ctx context.Context, userID string, gameID int64) error
}
