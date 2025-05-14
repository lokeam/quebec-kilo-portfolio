package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type IGDBClient interface {
	SearchGames(query string) ([]*models.Game, error)
	ExecuteQuery(ctx context.Context, query *types.IGDBQuery) ([]*types.IGDBResponse, error)
}
