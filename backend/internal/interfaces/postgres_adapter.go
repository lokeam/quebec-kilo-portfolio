package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type PostgresAdapter interface {
	GetUserLibraryItems(ctx context.Context, userID string) ([]models.Game, error)
	CreateLibraryGame(ctx context.Context, userID string, gameID int64) error
	DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error
}
