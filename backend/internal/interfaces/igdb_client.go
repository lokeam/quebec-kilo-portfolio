package interfaces

import "github.com/lokeam/qko-beta/internal/models"

type IGDBClient interface {
	SearchGames(query string) ([]*models.Game, error)
}
