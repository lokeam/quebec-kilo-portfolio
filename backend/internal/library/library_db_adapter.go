package library

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // NOTE: this registers pgx with database/sql
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/postgres"
)

type LibraryDbAdapter struct {
	client   *postgres.PostgresClient
	db       *sqlx.DB
	logger   interfaces.Logger
	scanner  *postgres.GameScanner
}

func NewLibraryDbAdapter(appContext *appcontext.AppContext) (*LibraryDbAdapter, error) {
	appContext.Logger.Debug("Creating LibraryDbAdapter", map[string]any{"appContext": appContext})

	// Create a PostgresClient
	client, err := postgres.NewPostgresClient(appContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres client %w", err)
	}

	// Create sqlx db from px pool
	db, err := sqlx.Connect("pgx", appContext.Config.Postgres.ConnectionString)
	if err != nil {
			return nil, fmt.Errorf("failed to create sqlx connection: %w", err)
	}

	// Register custom types for PostgreSQL arrays so that sqlx can handle string array types
	db.MapperFunc(strings.ToLower)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &LibraryDbAdapter{
		client: client,
		db:     db,
		logger: appContext.Logger,
		scanner: &postgres.GameScanner{},
	}, nil
}

// GET
func (la *LibraryDbAdapter) GetUserGame(ctx context.Context, userID string, gameID int64) (models.Game, bool, error) {
	la.logger.Debug("LibraryDbAdapter - GetUserGame called", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})

	query := `
		SELECT g.id, g.name, g.summary, g.cover_url, g.first_release_date, g.rating,
			   g.platform_names, g.genre_names, g.theme_names
		FROM user_library ul
		JOIN games g ON ul.game_id = g.id
		WHERE ul.user_id = $1 AND g.id = $2
		LIMIT 1
	`

	// Scanner util handles manual scanning of array columns
	row := la.db.QueryRowxContext(ctx, query, userID, gameID)
	game, err := la.scanner.ScanGame(row)

	if err == pgx.ErrNoRows {
		return models.Game{}, false, nil
	}

	if err != nil {
		return models.Game{}, false, fmt.Errorf("error querying user game: %w", err)
	}

	return game, true, nil
}

func (la *LibraryDbAdapter) GetLibraryItems(ctx context.Context, userID string) ([]models.Game, error) {
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

	// Scanner util handles manual scanning of array columns
	rows, err := la.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user library: %w", err)
	}
	defer rows.Close()

	return la.scanner.ScanGames(rows)
}

func (la *LibraryDbAdapter) GetUserLibraryItems(ctx context.Context, userID string) ([]models.Game, error) {
	// Just delegate to the existing method
	return la.GetLibraryItems(ctx, userID)
}

// PUT
func (la *LibraryDbAdapter) UpdateGameInLibrary(ctx context.Context, game models.Game) error {
	la.logger.Debug("LibraryDbAdapter - UpdateGameInLibrary called", map[string]any{
		"gameID": game.ID,
	})

	return la.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		query := `
			UPDATE games
			SET name = :name,
					summary = :summary,
					cover_url = :cover_url,
					first_release_date = :first_release_date,
					rating = :rating,
					platform_names = :platform_names,
					genre_names = :genre_names,
					theme_names = :theme_names
			WHERE id = :id
		`

		// Wrap Go slices in special type that knows how to convert to PostgreSql array syntax
		_, err := tx.ExecContext(ctx, query,
			game.Name,
			game.Summary,
			game.CoverURL,
			game.FirstReleaseDate,
			game.Rating,
			pq.Array(game.PlatformNames),
			pq.Array(game.GenreNames),
			pq.Array(game.ThemeNames),
			game.ID,
		)
		if err != nil {
			return fmt.Errorf("error updating game: %w", err)
		}

		return nil
	})
}

// POST
func (la *LibraryDbAdapter) AddGameToLibrary(ctx context.Context, userID string, gameID int64) error {
	la.logger.Info("LibraryDbAdapter - AddGameToLibrary called", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})

	return la.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		// First check if game exists in table, if not add it
		var gameExists bool
		err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM games WHERE id = $1)", gameID).Scan(&gameExists)
		if err != nil {
			return fmt.Errorf("error checking if game exists: %w", err)
		}

		if !gameExists {
			// Insert game if its not there
			_, err := tx.ExecContext(ctx, `
				INSERT INTO games (id, name, summary, cover_url, first_release_date, rating, platform_names, genre_names, theme_names)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, gameID, "", "", "", "", "", "", "", "")
			if err != nil {
				return fmt.Errorf("error inserting game: %w", err)
			}
		}

		// Add to user's library if it doesn't already exist
		_, err = tx.ExecContext(ctx, `
			INSERT INTO user_library (user_id, game_id, added_at)
			VALUES ($1, $2, NOW())
			ON CONFLICT (user_id, game_id) DO NOTHING
		`, userID, gameID)
		if err != nil {
			return fmt.Errorf("error adding game to user library: %w", err)
		}

		return nil
	})
}

// DELETE
func (la *LibraryDbAdapter) RemoveGameFromLibrary(ctx context.Context, userID string, gameID int64) error {
	la.logger.Info("LibraryDbAdapter - RemoveGameFromLibrary called", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})

	return la.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `
			DELETE FROM user_library
			WHERE user_id = $1 AND game_id = $2
		`, userID, gameID)
		if err != nil {
			return fmt.Errorf("error removing game from library: %w", err)
		}

		return nil
	})
}

// Helper fn to support transactions
func (la *LibraryDbAdapter) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	// Start a transaction
	tx, err := la.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Defer a rollback in the event anything catches fire
	defer func() {
		if panicValue := recover(); panicValue != nil {
			tx.Rollback()
			panic(panicValue) // re-throw panic after rollback so that it may be handled further up callstack
		}
	}()

	// Do the thing
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// Helper fn to determine if game exists in user library
func (la *LibraryDbAdapter) IsGameInLibrary(ctx context.Context, userID string, gameID int64) (bool, error) {
	la.logger.Debug("LibraryDbAdapter - IsGameInLibrary called", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})

	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_library
			WHERE user_id = $1 AND game_id = $2
		)
	`

	err := la.db.GetContext(ctx, &exists, query, userID, gameID)
	if err != nil {
		return false, fmt.Errorf("error checking if game is in library: %w", err)
	}

	return exists, nil
}
