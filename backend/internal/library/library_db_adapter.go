package library

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // NOTE: this registers pgx with database/sql
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/postgres"
	"github.com/lokeam/qko-beta/internal/types"
)

type LibraryDbAdapter struct {
	client   *postgres.PostgresClient
	db       *sqlx.DB
	logger   interfaces.Logger
	scanner  interfaces.GameScanner
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
// GetUserGame retrieves a game from a user's library.
// Returns:
// - The game if found
// - A boolean indicating if the game was found
// - An error, which will be:
//   - nil if the operation was successful or the game wasn't found
//   - ErrGameNotFound if the game doesn't exist in the user's library
//   - Another error if a database error occurred
func (la *LibraryDbAdapter) GetUserGame(
	ctx context.Context,
	userID string,
	gameID int64,
) (types.LibraryGameDBResult, []types.LibraryGamePhysicalLocationDBResponse, []types.LibraryGameDigitalLocationDBResponse, bool, error) {
	la.logger.Debug("LibraryDbAdapter - GetUserGame called", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})

	// 1. Get basic game info and library status
	gameQuery := `
	SELECT
			g.id,
			g.name,
			g.cover_url,
			g.first_release_date,
			g.rating,
			ug.favorite,
			w.id IS NOT NULL as is_in_wishlist,
			ug.game_type as game_type_display,
			LOWER(ug.game_type) as game_type_normalized
	FROM games g
	JOIN user_games ug ON g.id = ug.game_id
	LEFT JOIN wishlist w ON g.id = w.game_id AND w.user_id = $1
	WHERE ug.user_id = $1 AND g.id = $2
	LIMIT 1
`

	var game types.LibraryGameDBResult
	err := la.db.GetContext(ctx, &game, gameQuery, userID, gameID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.LibraryGameDBResult{}, nil, nil, false, ErrGameNotFound
		}
		return types.LibraryGameDBResult{}, nil, nil, false, fmt.Errorf("error querying user game: %w", err)
	}

	// 2. Get physical locations
	physicalQuery := `
		SELECT
			ug.game_id,
			pl.name as location_name,
			pl.location_type,
			sl.name as sublocation_name,
			sl.location_type as sublocation_type,
			sl.bg_color as sublocation_bg_color,
			p.name as platform_name
		FROM user_games ug
		JOIN physical_game_locations pgl ON ug.id = pgl.user_game_id
		JOIN sublocations sl ON pgl.sublocation_id = sl.id
		JOIN physical_locations pl ON sl.physical_location_id = pl.id
		JOIN platforms p ON ug.platform_id = p.id
		WHERE ug.user_id = $1 AND ug.game_id = $2
	`

	var physicalLocations []types.LibraryGamePhysicalLocationDBResponse
	if err := la.db.SelectContext(ctx, &physicalLocations, physicalQuery, userID, gameID); err != nil {
		return types.LibraryGameDBResult{}, nil, nil, false, fmt.Errorf("error querying physical locations: %w", err)
	}

	// 3. Get digital locations
	digitalQuery := `
		SELECT
			ug.game_id,
			dl.name as location_name,
			dl.name as normalized_name,
			p.name as platform_name
		FROM user_games ug
		JOIN digital_game_locations dgl ON ug.id = dgl.user_game_id
		JOIN digital_locations dl ON dgl.digital_location_id = dl.id
		JOIN platforms p ON ug.platform_id = p.id
		WHERE ug.user_id = $1 AND ug.game_id = $2
	`

	var digitalLocations []types.LibraryGameDigitalLocationDBResponse
	if err := la.db.SelectContext(ctx, &digitalLocations, digitalQuery, userID, gameID); err != nil {
		return types.LibraryGameDBResult{}, nil, nil, false, fmt.Errorf("error querying digital locations: %w", err)
	}

	return game, physicalLocations, digitalLocations, true, nil
}

// GetAllLibraryGames retrieves all games in a user's library with their locations
func (la *LibraryDbAdapter) GetAllLibraryGames(
	ctx context.Context,
	userID string,
) ([]types.LibraryGameDBResult, []types.LibraryGamePhysicalLocationDBResponse, []types.LibraryGameDigitalLocationDBResponse, error) {
	la.logger.Debug("LibraryDbAdapter - GetAllLibraryGames called", map[string]any{
		"userID": userID,
	})

	// 1. Get basic game info and library status
	gameQuery := `
	SELECT
			g.id,
			g.name,
			g.cover_url,
			g.first_release_date,
			g.rating,
			ug.favorite,
			w.id IS NOT NULL as is_in_wishlist,
			ug.game_type as game_type_display,
			LOWER(ug.game_type) as game_type_normalized
	FROM games g
	JOIN user_games ug ON g.id = ug.game_id
	LEFT JOIN wishlist w ON g.id = w.game_id AND w.user_id = $1
	WHERE ug.user_id = $1
`

	var games []types.LibraryGameDBResult
	if err := la.db.SelectContext(ctx, &games, gameQuery, userID); err != nil {
		return nil, nil, nil, fmt.Errorf("error querying user library games: %w", err)
	}

	// 2. Get physical locations
	physicalQuery := `
		SELECT
			ug.game_id,
			pl.name as location_name,
			pl.location_type,
			sl.name as sublocation_name,
			sl.location_type as sublocation_type,
			sl.bg_color as sublocation_bg_color,
			p.name as platform_name
		FROM user_games ug
		JOIN physical_game_locations pgl ON ug.id = pgl.user_game_id
		JOIN sublocations sl ON pgl.sublocation_id = sl.id
		JOIN physical_locations pl ON sl.physical_location_id = pl.id
		JOIN platforms p ON ug.platform_id = p.id
		WHERE ug.user_id = $1
	`

	var physicalLocations []types.LibraryGamePhysicalLocationDBResponse
	if err := la.db.SelectContext(ctx, &physicalLocations, physicalQuery, userID); err != nil {
		return nil, nil, nil, fmt.Errorf("error querying physical locations: %w", err)
	}

	// 3. Get digital locations
	digitalQuery := `
		SELECT
			ug.game_id,
			dl.name as location_name,
			dl.name as normalized_name,
			p.name as platform_name
		FROM user_games ug
		JOIN digital_game_locations dgl ON ug.id = dgl.user_game_id
		JOIN digital_locations dl ON dgl.digital_location_id = dl.id
		JOIN platforms p ON ug.platform_id = p.id
		WHERE ug.user_id = $1
	`

	var digitalLocations []types.LibraryGameDigitalLocationDBResponse
	if err := la.db.SelectContext(ctx, &digitalLocations, digitalQuery, userID); err != nil {
		return nil, nil, nil, fmt.Errorf("error querying digital locations: %w", err)
	}

	return games, physicalLocations, digitalLocations, nil
}

// GetUserLibraryItems is an alias for GetAllLibraryGames to maintain backward compatibility
func (la *LibraryDbAdapter) GetUserLibraryItems(
	ctx context.Context,
	userID string,
) ([]types.LibraryGameDBResult, []types.LibraryGamePhysicalLocationDBResponse, []types.LibraryGameDigitalLocationDBResponse, error) {
	return la.GetAllLibraryGames(ctx, userID)
}

// PUT
func (la *LibraryDbAdapter) UpdateLibraryGame(ctx context.Context, game models.LibraryGame) error {
	la.logger.Debug("LibraryDbAdapter - UpdateLibraryGame called", map[string]any{
		"gameID": game.GameID,
	})

	return postgres.WithTransaction(ctx, la.db, la.logger, func(tx *sqlx.Tx) error {
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
			game.GameName,
			game.GameCoverURL,
			game.GameFirstReleaseDate,
			game.GameRating,
			pq.Array(game.PlatformLocations),
			pq.Array(game.GameThemeNames),
			game.GameID,
		)
		if err != nil {
			return fmt.Errorf("error updating game: %w", err)
		}

		return nil
	})
}

// POST
func (la *LibraryDbAdapter) CreateLibraryGame(ctx context.Context, userID string, game models.LibraryGame) error {
	la.logger.Info("LibraryDbAdapter - CreateLibraryGame called", map[string]any{
		"userID": userID,
		"gameID": game.GameID,
	})

	return postgres.WithTransaction(ctx, la.db, la.logger, func(tx *sqlx.Tx) error {
		// First check if game exists in table, if not add it
		var gameExists bool
		err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM games WHERE id = $1)", game.GameID).Scan(&gameExists)
		if err != nil {
			return fmt.Errorf("error checking if game exists: %w", err)
		}

		// Insert game if it doesn't exist
		if !gameExists {
			_, err := tx.ExecContext(ctx, `
					INSERT INTO games (id, name, cover_url, first_release_date, rating)
					VALUES ($1, $2, $3, $4, $5)
			`, game.GameID, game.GameName, game.GameCoverURL, game.GameFirstReleaseDate, game.GameRating)
			if err != nil {
					return fmt.Errorf("error inserting game: %w", err)
			}
	}

		// For each platform location:
		for i := 0; i < len(game.PlatformLocations); i++ {
			location := game.PlatformLocations[i]

			// Ensure platform exists
			_, err = tx.ExecContext(ctx, `
				INSERT INTO platforms (id, name, category, model)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (id) DO NOTHING
			`, location.PlatformID, location.PlatformName,
			getPlatformCategory(location.PlatformName),
			getPlatformModel(location.PlatformName))
			if err != nil {
			return fmt.Errorf("error ensuring platform exists at index %d: %w", i, err)
			}

			// Ensure game-platform combination exists
			_, err = tx.ExecContext(ctx, `
				INSERT INTO game_platforms (game_id, platform_id)
				VALUES ($1, $2)
				ON CONFLICT (game_id, platform_id) DO NOTHING
			`, game.GameID, location.PlatformID)
			if err != nil {
			return fmt.Errorf("error ensuring game-platform combination at index %d: %w", i, err)
			}

			// Insert user_game entry
			var userGameID int
			err = tx.QueryRowContext(ctx, `
					INSERT INTO user_games (user_id, game_id, platform_id, game_type)
					VALUES ($1, $2, $3, $4)
					RETURNING id
			`, userID, game.GameID, location.PlatformID, location.Type).Scan(&userGameID)
			if err != nil {
					if strings.Contains(err.Error(), "unique constraint") {
							// If this game+platform already exists, get its ID
							err = tx.QueryRowContext(ctx, `
									SELECT id FROM user_games
									WHERE user_id = $1 AND game_id = $2 AND platform_id = $3
							`, userID, game.GameID, location.PlatformID).Scan(&userGameID)
							if err != nil {
									return fmt.Errorf("error finding existing user game at index %d: %w", i, err)
							}
					} else {
							return fmt.Errorf("error inserting user game at index %d: %w", i, err)
					}
			}

			// Insert location mapping
			if location.Type == "physical" {
				// First get the sublocation UUID by name
				var sublocationID uuid.UUID
				err = tx.QueryRowContext(ctx, `
						SELECT id FROM sublocations
						WHERE id = $1 AND user_id = $2
				`, location.Location.SublocationID, userID).Scan(&sublocationID)
				if err != nil {
						return fmt.Errorf("error finding sublocation at index %d: %w", i, err)
				}

				_, err = tx.ExecContext(ctx, `
						INSERT INTO physical_game_locations (user_game_id, sublocation_id)
						VALUES ($1, $2)
				`, userGameID, sublocationID)
			} else {
					// First get the digital location UUID by name
					var digitalLocationID uuid.UUID
					err = tx.QueryRowContext(ctx, `
							SELECT id FROM digital_locations
							WHERE id = $1 AND user_id = $2
					`, location.Location.DigitalLocationID, userID).Scan(&digitalLocationID)
					if err != nil {
							return fmt.Errorf("error finding digital location at index %d: %w", i, err)
					}

					_, err = tx.ExecContext(ctx, `
							INSERT INTO digital_game_locations (user_game_id, digital_location_id)
							VALUES ($1, $2)
					`, userGameID, digitalLocationID)
			}
			if err != nil {
					return fmt.Errorf("error inserting game location at index %d: %w", i, err)
			}
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

	return postgres.WithTransaction(ctx, la.db, la.logger, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `
			DELETE FROM user_games
			WHERE user_id = $1 AND game_id = $2
		`, userID, gameID)
		if err != nil {
			return fmt.Errorf("error removing game from library: %w", err)
		}

		return nil
	})
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
			SELECT 1 FROM user_games
			WHERE user_id = $1 AND game_id = $2
		)
	`

	err := la.db.GetContext(ctx, &exists, query, userID, gameID)
	if err != nil {
		return false, fmt.Errorf("error checking if game is in library: %w", err)
	}

	return exists, nil
}

// Helper function to determine platform category
func getPlatformCategory(platformName string) string {
	switch {
	case strings.Contains(strings.ToLower(platformName), "pc"):
			return "pc"
	case strings.Contains(strings.ToLower(platformName), "mobile"):
			return "mobile"
	default:
			return "console"
	}
}

// Helper function to determine platform model
func getPlatformModel(platformName string) string {
	return platformName // For now, use the platform name as the model
}