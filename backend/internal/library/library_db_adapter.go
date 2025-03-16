package library

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/postgres"
	"github.com/lokeam/qko-beta/internal/types"
)

type LibraryDbAdapter struct {
	client *postgres.PostgresClient
	logger interfaces.Logger
}

func NewLibraryDbAdapter(appContext *appcontext.AppContext) (*LibraryDbAdapter, error) {
	appContext.Logger.Debug("Creating LibraryDbAdapter", map[string]any{"appContext": appContext})

	// Create a PostgresClient
	client, err := postgres.NewPostgresClient(appContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres client %w", err)
	}

	return &LibraryDbAdapter{
		client: client,
		logger: appContext.Logger,
	}, nil
}

// GET
func (la *LibraryDbAdapter) GetLibraryItems(ctx context.Context, userID string) ([]types.Game, error) {
	la.logger.Debug("LibraryDbAdapter - GetUserLibraryItems called", map[string]any{
		"userID": userID,
	})

	// Todo: change user_library to whatever the table is called
	query := `
		SELECT g.id, g.name, g.summary, g.cover_url, g.first_release_date, g.rating,
		       g.platform_names, g.genre_names, g.theme_names
		FROM user_library ul
		JOIN games g ON ul.game_id = g.id
		WHERE ul.user_id = $1
	`

	rows, err := la.client.GetPool().Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user library %w", err)
	}
	defer rows.Close()

	var games []types.Game
	for rows.Next() {
		var game types.Game
		var platforms, genres, themes []string

		// Todo: Check these columns to make sure they are correct with what exists in db
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
			return nil, fmt.Errorf("error scanning game row %w", err)
		}

		game.PlatformNames = platforms
		game.GenreNames = genres
		game.ThemeNames = themes

		games = append(games, game)
	}

	if err := rows.Err(); err != nil{
		return nil, fmt.Errorf("error itereating rows: %w", err)
	}

	return games, nil
}

// POST
func (la *LibraryDbAdapter) AddGameToLibrary(ctx context.Context, userID string, gameID int64) error {
	la.logger.Info("LibraryDbAdapter - AddGameToLibrary called", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})

	// What does tx mean?
	tx, err := la.client.GetPool().BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// First check if game exists in games table, if not insert it
	var gameExists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM games WHERE id = $1)", gameID).Scan(&gameExists)
	if err != nil {
		return fmt.Errorf("error checking if game exists: %w", err)
	}

	if !gameExists {
		// Insert game if it doesn't exist
		_, err = tx.Exec(ctx, `
			INSERT INTO games (id, name, summary, cover_url, first_release_date, rating, platform_names, genre_names, theme_names)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, gameID, "", "", "", "", "", "", "", "")
		if err != nil {
			return fmt.Errorf("error inserting game: %w", err)
		}
	}

	// Now add to user library if not already there
	_, err = tx.Exec(ctx, `
		INSERT INTO user_library (user_id, game_id, added_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id, game_id) DO NOTHING
	`, userID, gameID)
	if err != nil {
		return fmt.Errorf("error adding game to library: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// DELETE
func (la *LibraryDbAdapter) RemoveGameFromLibrary(ctx context.Context, userID string, gameID int64) error {
	la.logger.Info("LibraryDbAdapter - RemoveGameFromLibrary called", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})

	_, err := la.client.GetPool().Exec(ctx, `
		DELETE FROM user_library
		WHERE user_id = $1 AND game_id = $2
		`, userID, gameID)
	if err != nil {
		return fmt.Errorf("error removing game from library: %w", err)
	}

	return nil
}

// Add this method to LibraryDbAdapter
func (la *LibraryDbAdapter) GetUserLibraryItems(ctx context.Context, userID string) ([]types.Game, error) {
	// Just delegate to the existing method
	return la.GetLibraryItems(ctx, userID)
}