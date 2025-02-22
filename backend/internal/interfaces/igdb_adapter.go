package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/types"
)

type IGDBAdapter interface {
	SearchGames(
		ctx context.Context,
		query string,
		limit int,
	) ([]*types.Game, error)
}
