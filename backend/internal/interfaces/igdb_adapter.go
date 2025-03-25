package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type IGDBAdapter interface {
	SearchGames(
		ctx context.Context,
		query string,
		limit int,
	) ([]*models.Game, error)
}
