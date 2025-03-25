package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/models"
)

// GameScanner handles scanning a db row into a Game model with proper array handling
type GameScanner struct{}

// Scans a single row into a Game model
func (gs *GameScanner) ScanGame(row *sqlx.Row) (models.Game, error) {
	var game models.Game
	var platforms, genres, themes []string

	err := row.Scan(
		&game.ID,
		&game.Name,
		&game.Summary,
		&game.CoverURL,
		&game.FirstReleaseDate,
		&game.Rating,
		&platforms,
		&genres,
		&themes,
	)
	if err != nil {
		return models.Game{}, fmt.Errorf("error scanning game row: %w", err)
	}

	game.PlatformNames = platforms
	game.GenreNames = genres
	game.ThemeNames = themes

	return game, nil
}

// Scans multiple rows into a slice of Game models
func (gs *GameScanner) ScanGames(rows *sqlx.Rows) ([]models.Game, error) {
	var games []models.Game

	for rows.Next() {
		var game models.Game
		var platforms, genres, themes []string

		err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.Summary,
			&game.CoverURL,
			&game.FirstReleaseDate,
			&game.Rating,
			&platforms,
			&genres,
      &themes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning game row: %w", err)
		}

		game.PlatformNames = platforms
		game.GenreNames = genres
		game.ThemeNames = themes

		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over game rows: %w", err)
	}

	return games, nil
}
