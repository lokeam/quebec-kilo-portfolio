package interfaces

import (
	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/models"
)

type GameScanner interface {
	ScanGame(row *sqlx.Row) (models.Game, error)
	ScanGames(rows *sqlx.Rows) ([]models.Game, error)
}
