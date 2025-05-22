package library

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // NOTE: this registers pgx with database/sql
	"github.com/jmoiron/sqlx"
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
// GetSingleLibraryGame retrieves a game from a user's library.
// Returns:
// - The game if found
// - A boolean indicating if the game was found
// - An error, which will be:
//   - nil if the operation was successful or the game wasn't found
//   - ErrGameNotFound if the game doesn't exist in the user's library
//   - Another error if a database error occurred
func (la *LibraryDbAdapter) GetSingleLibraryGame(
	ctx context.Context,
	userID string,
	gameID int64,
) (
	types.LibraryGameDBResult,
	[]types.LibraryGamePhysicalLocationDBResponse,
	[]types.LibraryGameDigitalLocationDBResponse,
	error,
) {
	la.logger.Debug("LibraryDbAdapter - GetSingleLibraryGame called", map[string]any{
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
			LOWER(ug.game_type) as game_type_normalized,
			p.id as platform_id,
			p.name as platform_name
	FROM games g
	JOIN user_games ug ON g.id = ug.game_id
	LEFT JOIN wishlist w ON g.id = w.game_id AND w.user_id = $1
	LEFT JOIN platforms p ON ug.platform_id = p.id
	WHERE ug.user_id = $1 AND g.id = $2
	`

	var game types.LibraryGameDBResult
	err := la.db.GetContext(ctx, &game, gameQuery, userID, gameID)
	if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
					return types.LibraryGameDBResult{}, nil, nil, ErrGameNotFound
			}
			return types.LibraryGameDBResult{}, nil, nil, fmt.Errorf("error querying user game: %w", err)
	}

	// 2. Get locations
	locationQuery := `
    SELECT
        ug.game_id,
        p.id as platform_id,
        p.name as platform_name,
        ug.game_type as type,
        COALESCE(pl.id::text, dl.id::text) as location_id,
        COALESCE(pl.name, dl.name) as location_name,
        COALESCE(pl.location_type, 'digital') as location_type,
        sl.id::text as sublocation_id,
        sl.name as sublocation_name,
        sl.location_type as sublocation_type,
        sl.bg_color as sublocation_bg_color,
        dl.is_active
    FROM user_games ug
    JOIN platforms p ON ug.platform_id = p.id
    LEFT JOIN physical_game_locations pgl ON ug.id = pgl.user_game_id
    LEFT JOIN sublocations sl ON pgl.sublocation_id = sl.id
    LEFT JOIN physical_locations pl ON sl.physical_location_id = pl.id
    LEFT JOIN digital_game_locations dgl ON ug.id = dgl.user_game_id
    LEFT JOIN digital_locations dl ON dgl.digital_location_id = dl.id
    WHERE ug.user_id = $1 AND ug.game_id = $2
    `

		var locations []types.GameLocationDBResult
    if err := la.db.SelectContext(ctx, &locations, locationQuery, userID, gameID); err != nil {
        return types.LibraryGameDBResult{}, nil, nil, fmt.Errorf("error querying game locations: %w", err)
    }

		var physicalLocations []types.LibraryGamePhysicalLocationDBResponse
    var digitalLocations []types.LibraryGameDigitalLocationDBResponse

    for i := 0; i < len(locations); i++ {
			loc := locations[i]
			if loc.Type == "physical" {
				physicalLoc := types.LibraryGamePhysicalLocationDBResponse{
					GameID: loc.GameID,
					PlatformID: loc.PlatformID,
					PlatformName: loc.PlatformName,
					LocationID: loc.LocationID,
					LocationName: loc.LocationName,
					LocationType: loc.LocationType,
					SublocationID: loc.SublocationID,
					SublocationName: loc.SublocationName,
					SublocationType: loc.SublocationType,
					SublocationBgColor: loc.SublocationBgColor,
				}
				physicalLocations = append(physicalLocations, physicalLoc)
			} else {
					digitalLocations = append(digitalLocations, types.LibraryGameDigitalLocationDBResponse{
						GameID: loc.GameID,
						PlatformID: loc.PlatformID,
						PlatformName: loc.PlatformName,
						LocationID: loc.LocationID,
						LocationName: loc.LocationName,
						IsActive: loc.IsActive,
					})
			}
    }

	return game, physicalLocations, digitalLocations, nil
}

// GetAllLibraryGames retrieves all games in a user's library with their locations
func (la *LibraryDbAdapter) GetAllLibraryGames(
	ctx context.Context,
	userID string,
) (
	[]types.LibraryGameDBResult,
	[]types.LibraryGamePhysicalLocationDBResponse,
	[]types.LibraryGameDigitalLocationDBResponse,
	error,
	) {
	la.logger.Debug("LibraryDbAdapter - GetAllLibraryGames called", map[string]any{
		"userID": userID,
	})

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
			LOWER(ug.game_type) as game_type_normalized,
			p.id as platform_id,
			p.name as platform_name
	FROM games g
	JOIN user_games ug ON g.id = ug.game_id
	LEFT JOIN wishlist w ON g.id = w.game_id AND w.user_id = $1
	LEFT JOIN platforms p ON ug.platform_id = p.id
	WHERE ug.user_id = $1
	`

	var games []types.LibraryGameDBResult
	if err := la.db.SelectContext(ctx, &games, gameQuery, userID); err != nil {
		return nil, nil, nil, fmt.Errorf("error querying user library games: %w", err)
	}

	// 2. Get unified locations
	locationQuery := `
    SELECT
        ug.game_id,
        p.id as platform_id,
        p.name as platform_name,
        ug.game_type as type,
        COALESCE(pl.id::text, dl.id::text) as location_id,
        COALESCE(pl.name, dl.name) as location_name,
        COALESCE(pl.location_type, 'digital') as location_type,
        sl.id::text as sublocation_id,
        sl.name as sublocation_name,
        sl.location_type as sublocation_type,
        sl.bg_color as sublocation_bg_color,
        dl.is_active
    FROM user_games ug
    JOIN platforms p ON ug.platform_id = p.id
    LEFT JOIN physical_game_locations pgl ON ug.id = pgl.user_game_id
    LEFT JOIN sublocations sl ON pgl.sublocation_id = sl.id
    LEFT JOIN physical_locations pl ON sl.physical_location_id = pl.id
    LEFT JOIN digital_game_locations dgl ON ug.id = dgl.user_game_id
    LEFT JOIN digital_locations dl ON dgl.digital_location_id = dl.id
    WHERE ug.user_id = $1
    `

		var locations []types.GameLocationDBResult
    if err := la.db.SelectContext(ctx, &locations, locationQuery, userID); err != nil {
        return nil, nil, nil, fmt.Errorf("error querying game locations: %w", err)
    }

		var physicalLocations []types.LibraryGamePhysicalLocationDBResponse
    var digitalLocations []types.LibraryGameDigitalLocationDBResponse

		for i := 0; i < len(locations); i++ {
			loc := locations[i]
			if loc.Type == "physical" {
					physicalLoc := types.LibraryGamePhysicalLocationDBResponse{
							GameID: loc.GameID,
							PlatformID: loc.PlatformID,
							PlatformName: loc.PlatformName,
							LocationID: loc.LocationID,
							LocationName: loc.LocationName,
							LocationType: loc.LocationType,
							SublocationID: loc.SublocationID,
							SublocationName: loc.SublocationName,
							SublocationType: loc.SublocationType,
							SublocationBgColor: loc.SublocationBgColor,
					}
					physicalLocations = append(physicalLocations, physicalLoc)
			} else {
					digitalLocations = append(digitalLocations, types.LibraryGameDigitalLocationDBResponse{
							GameID: loc.GameID,
							PlatformID: loc.PlatformID,
							PlatformName: loc.PlatformName,
							LocationID: loc.LocationID,
							LocationName: loc.LocationName,
							IsActive: loc.IsActive,
					})
			}
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
// UpdateLibraryGame updates a game in the user's library.
// It updates the game details in the games table and any associated platform locations.
// Returns:
// - nil if the operation was successful
// - ErrGameNotFound if the game doesn't exist in the user's library
// - Another error if a database error occurred
func (la *LibraryDbAdapter) UpdateLibraryGame(ctx context.Context, userID string, game models.LibraryGame) error {
	la.logger.Info("LibraryDbAdapter - UpdateLibraryGame called", map[string]any{
		"userID": userID,
		"gameID": game.GameID,
	})

	return postgres.WithTransaction(ctx, la.db, la.logger, func(tx *sqlx.Tx) error {
		// First check if the game exists in the user's library
		var exists bool
		err := tx.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM user_games
				WHERE user_id = $1 AND game_id = $2
			)
		`, userID, game.GameID).Scan(&exists)
		if err != nil {
			la.logger.Error("Failed to check if game exists in library", map[string]any{
				"error": err,
				"userID": userID,
				"gameID": game.GameID,
			})
			return fmt.Errorf("error checking if game exists in library: %w", err)
		}

		if !exists {
			la.logger.Info("Game not found in user's library", map[string]any{
				"userID": userID,
				"gameID": game.GameID,
			})
			return ErrGameNotFound
		}

		// Update game details
		_, err = tx.ExecContext(ctx, `
			UPDATE games
			SET name = $1,
					cover_url = $2,
					first_release_date = $3,
					rating = $4
			WHERE id = $5
		`, game.GameName, game.GameCoverURL, game.GameFirstReleaseDate, game.GameRating, game.GameID)
		if err != nil {
			la.logger.Error("Failed to update game details", map[string]any{
				"error": err,
				"gameID": game.GameID,
			})
			return fmt.Errorf("error updating game details: %w", err)
		}

		// For each platform version, update or add location mapping
		for i, location := range game.PlatformLocations {
			la.logger.Info("Processing platform update", map[string]any{
				"index": i,
				"platformID": location.PlatformID,
				"platformName": location.PlatformName,
			})

			// Get the user_game_id for this platform
			var userGameID int
			err = tx.QueryRowContext(ctx, `
				SELECT id FROM user_games
				WHERE user_id = $1 AND game_id = $2 AND platform_id = $3
			`, userID, game.GameID, location.PlatformID).Scan(&userGameID)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					// Platform doesn't exist, add it
					err = tx.QueryRowContext(ctx, `
						INSERT INTO user_games (user_id, game_id, platform_id, game_type)
						VALUES ($1, $2, $3, $4)
						RETURNING id
					`, userID, game.GameID, location.PlatformID, location.Type).Scan(&userGameID)
					if err != nil {
						la.logger.Error("Failed to add new platform", map[string]any{
							"error": err,
							"platformID": location.PlatformID,
						})
						return fmt.Errorf("error adding new platform at index %d: %w", i, err)
					}
				} else {
					la.logger.Error("Failed to get user game ID", map[string]any{
						"error": err,
						"platformID": location.PlatformID,
					})
					return fmt.Errorf("error getting user game ID at index %d: %w", i, err)
				}
			}

			// Update location mapping based on type
			if location.Type == "physical" {
				// Delete existing physical location
				_, err = tx.ExecContext(ctx, `
					DELETE FROM physical_game_locations
					WHERE user_game_id = $1
				`, userGameID)
				if err != nil {
					la.logger.Error("Failed to delete existing physical location", map[string]any{
						"error": err,
						"userGameID": userGameID,
					})
					return fmt.Errorf("error deleting existing physical location at index %d: %w", i, err)
				}

				// Add new physical location
				_, err = tx.ExecContext(ctx, `
					INSERT INTO physical_game_locations (user_game_id, sublocation_id)
					VALUES ($1, $2)
				`, userGameID, location.Location.SublocationID)
			} else {
				// Delete existing digital location
				_, err = tx.ExecContext(ctx, `
					DELETE FROM digital_game_locations
					WHERE user_game_id = $1
				`, userGameID)
				if err != nil {
					la.logger.Error("Failed to delete existing digital location", map[string]any{
						"error": err,
						"userGameID": userGameID,
					})
					return fmt.Errorf("error deleting existing digital location at index %d: %w", i, err)
				}

				// Add new digital location
				_, err = tx.ExecContext(ctx, `
					INSERT INTO digital_game_locations (user_game_id, digital_location_id)
					VALUES ($1, $2)
				`, userGameID, location.Location.DigitalLocationID)
			}
			if err != nil {
				la.logger.Error("Failed to update location mapping", map[string]any{
					"error": err,
					"userGameID": userGameID,
					"type": location.Type,
				})
				return fmt.Errorf("error updating location mapping at index %d: %w", i, err)
			}
		}

		la.logger.Info("Successfully updated game in library", map[string]any{
			"userID": userID,
			"gameID": game.GameID,
		})
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
			// STEP 1: Ensure game exists in games table
			la.logger.Info("STEP 1: Ensuring game exists", map[string]any{
					"gameID": game.GameID,
			})
			_, err := tx.ExecContext(ctx, `
					INSERT INTO games (id, name, cover_url, first_release_date, rating)
					VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT (id) DO NOTHING
			`, game.GameID, game.GameName, game.GameCoverURL, game.GameFirstReleaseDate, game.GameRating)
			if err != nil {
					la.logger.Error("Failed to ensure game exists", map[string]any{
							"error": err,
							"gameID": game.GameID,
					})
					return fmt.Errorf("error ensuring game exists: %w", err)
			}

			// STEP 2: For each platform version the user wants to add
			for i, location := range game.PlatformLocations {
					la.logger.Info("Processing platform", map[string]any{
							"index": i,
							"platformID": location.PlatformID,
							"platformName": location.PlatformName,
					})

					// STEP 2a: Ensure platform exists
					_, err = tx.ExecContext(ctx, `
							INSERT INTO platforms (id, name, category, model)
							VALUES ($1, $2, $3, $4)
							ON CONFLICT (id) DO NOTHING
					`, location.PlatformID, location.PlatformName,
						 getPlatformCategory(location.PlatformName),
						 getPlatformModel(location.PlatformName))
					if err != nil {
							la.logger.Error("Failed to ensure platform exists", map[string]any{
									"error": err,
									"platformID": location.PlatformID,
							})
							return fmt.Errorf("error ensuring platform exists at index %d: %w", i, err)
					}

					// STEP 2b: Add game+platform to user's library
					var userGameID int
					err = tx.QueryRowContext(ctx, `
							INSERT INTO user_games (user_id, game_id, platform_id, game_type)
							VALUES ($1, $2, $3, $4)
							RETURNING id
					`, userID, game.GameID, location.PlatformID, location.Type).Scan(&userGameID)
					if err != nil {
							if strings.Contains(err.Error(), "unique constraint") {
									la.logger.Info("Game already exists, getting ID", map[string]any{
											"userID": userID,
											"gameID": game.GameID,
											"platformID": location.PlatformID,
									})
									err = tx.QueryRowContext(ctx, `
											SELECT id FROM user_games
											WHERE user_id = $1 AND game_id = $2 AND platform_id = $3
									`, userID, game.GameID, location.PlatformID).Scan(&userGameID)
									if err != nil {
											la.logger.Error("Failed to find existing user game", map[string]any{
													"error": err,
													"userID": userID,
													"gameID": game.GameID,
													"platformID": location.PlatformID,
											})
											return fmt.Errorf("error finding existing user game at index %d: %w", i, err)
									}
							} else {
									la.logger.Error("Failed to insert user game", map[string]any{
											"error": err,
											"userID": userID,
											"gameID": game.GameID,
											"platformID": location.PlatformID,
									})
									return fmt.Errorf("error inserting user game at index %d: %w", i, err)
							}
					}

					// STEP 3: Add location mapping based on type
					if location.Type == "physical" {
							la.logger.Info("Adding physical location", map[string]any{
									"userGameID": userGameID,
									"sublocationID": location.Location.SublocationID,
							})
							_, err = tx.ExecContext(ctx, `
									INSERT INTO physical_game_locations (user_game_id, sublocation_id)
									VALUES ($1, $2)
							`, userGameID, location.Location.SublocationID)
					} else {
							la.logger.Info("Adding digital location", map[string]any{
									"userGameID": userGameID,
									"digitalLocationID": location.Location.DigitalLocationID,
							})
							_, err = tx.ExecContext(ctx, `
									INSERT INTO digital_game_locations (user_game_id, digital_location_id)
									VALUES ($1, $2)
							`, userGameID, location.Location.DigitalLocationID)
					}
					if err != nil {
							la.logger.Error("Failed to insert game location", map[string]any{
									"error": err,
									"userGameID": userGameID,
									"type": location.Type,
							})
							return fmt.Errorf("error inserting game location at index %d: %w", i, err)
					}
			}

			return nil
	})
}

// DeleteLibraryGame removes a game from a user's library.
// It will automatically remove any associated physical or digital location mappings
// due to ON DELETE CASCADE constraints in the database.
// Returns:
// - nil if the operation was successful
// - ErrGameNotFound if the game doesn't exist in the user's library
// - Another error if a database error occurred
func (la *LibraryDbAdapter) DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error {
	la.logger.Info("LibraryDbAdapter - DeleteLibraryGame called", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})

	return postgres.WithTransaction(ctx, la.db, la.logger, func(tx *sqlx.Tx) error {
		// First check if the game exists in the user's library
		var exists bool
		err := tx.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM user_games
				WHERE user_id = $1 AND game_id = $2
			)
		`, userID, gameID).Scan(&exists)
		if err != nil {
			la.logger.Error("Failed to check if game exists in library", map[string]any{
				"error": err,
				"userID": userID,
				"gameID": gameID,
			})
			return fmt.Errorf("error checking if game exists in library: %w", err)
		}

		if !exists {
			la.logger.Info("Game not found in user's library", map[string]any{
				"userID": userID,
				"gameID": gameID,
			})
			return ErrGameNotFound
		}

		// Delete the game from user's library
		// Note: Due to ON DELETE CASCADE in our schema, this will automatically
		// remove related records from physical_game_locations and digital_game_locations
		_, err = tx.ExecContext(ctx, `
			DELETE FROM user_games
			WHERE user_id = $1 AND game_id = $2
		`, userID, gameID)
		if err != nil {
			la.logger.Error("Failed to delete game from library", map[string]any{
				"error": err,
				"userID": userID,
				"gameID": gameID,
			})
			return fmt.Errorf("error deleting game from library: %w", err)
		}

		la.logger.Info("Successfully deleted game from library", map[string]any{
			"userID": userID,
			"gameID": gameID,
		})
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