package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type LibraryDbAdapter interface {
	GetSingleLibraryGame(ctx context.Context, userID string, gameID int64) (models.LibraryGame, bool, error)
	GetAllLibraryGames(ctx context.Context, userID string) ([]models.LibraryGame, error)
	GetUserLibraryItems(ctx context.Context, userID string) ([]models.LibraryGame, error)
	UpdateLibraryGame(ctx context.Context, game models.LibraryGame) error
	CreateLibraryGame(ctx context.Context, userID string, game models.LibraryGame) error
	DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error
	IsGameInLibrary(ctx context.Context, userID string, gameID int64) (bool, error)
}