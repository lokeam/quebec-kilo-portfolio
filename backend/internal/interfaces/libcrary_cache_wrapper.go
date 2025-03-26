package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type LibraryCacheWrapper interface {
	GetCachedLibraryItems(ctx context.Context, userID string) ([]models.Game, error)
	SetCachedLibraryItems(ctx context.Context, userID string, games []models.Game) error
	GetCachedGame(ctx context.Context, userID string, gameID int64) (*models.Game, bool, error)
	SetCachedGame(ctx context.Context, userID string, game models.Game) error
	InvalidateUserCache(ctx context.Context, userID string) error
	InvalidateGameCache(ctx context.Context, userID string, gameID int64) error
}