package interfaces

import (
	"context"

	igdb "github.com/Henry-Sarabia/igdb"
)

type IGDBAdapter interface {
	SearchGames(
		ctx context.Context,
		query string,
		limit int,
	) ([]*igdb.Game, error)
}
