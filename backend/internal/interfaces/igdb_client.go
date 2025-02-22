package interfaces

import "github.com/lokeam/qko-beta/internal/types"

type IGDBClient interface {
	SearchGames(query string) ([]*types.Game, error)
}
