package library

import (
	"context"
	"fmt"
	"strings"
	"time"

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
) (types.LibraryGameItemBFFResponseFINAL, error) {
	la.logger.Debug("LibraryDbAdapter - GetSingleLibraryGame called", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})

	// First check if game exists within user's library
	var exists bool
	err := la.db.QueryRowContext(
		ctx,
		CheckIfGameIsInLibraryQuery,
		userID,
		gameID,
	).Scan(&exists)
	if err != nil {
		return types.LibraryGameItemBFFResponseFINAL{}, fmt.Errorf("error checking if game exists in library: %w", err)
	}

	if !exists {
		return types.LibraryGameItemBFFResponseFINAL{}, ErrGameNotFound
	}

	// Then get game data
	var game models.LibraryGameDB
	err = la.db.GetContext(
		ctx,
		&game,
		GetSingleLibraryGameQuery,
		userID,
		gameID,
	)
	if err != nil {
		return types.LibraryGameItemBFFResponseFINAL{}, fmt.Errorf("error getting game data for single game: %w", err)
	}

	// Get locations data
	var locations []models.GameLocationDatabaseEntry
	err = la.db.SelectContext(
		ctx,
		&locations,
		GetLibraryLocationsQuery,
		userID,
	)
	if err != nil {
		return types.LibraryGameItemBFFResponseFINAL{}, fmt.Errorf("error getting locations data for single game: %w", err)
	}

	// Transform locations to BFF response format
	var bffLocations []types.LibraryGamesByPlatformAndLocationItemFINAL
	for _, loc := range locations {
		if loc.GameID == gameID {
			bffLocations = append(bffLocations, TransformLocationDBToResponse(loc))
		}
	}

	// Transform game to BFF response format
	return TransformGameDBToResponse(game, bffLocations), nil
}

// PUT
// UpdateLibraryGame updates a game in the user's library.
// It updates the game details in the games table and any associated platform locations.
// Returns:
// - nil if the operation was successful
// - ErrGameNotFound if the game doesn't exist in the user's library
// - Another error if a database error occurred
func (la *LibraryDbAdapter) UpdateLibraryGame(ctx context.Context, game models.GameToSave) error {
	la.logger.Info("LibraryDbAdapter - UpdateLibraryGame called", map[string]any{
		"gameID": game.GameID,
	})

	// Update game in database
	_, err := la.db.ExecContext(
		ctx,
		UpdateLibraryGameQuery,
		game.GameID,
		game.GameName,
		game.GameCoverURL,
		game.GameType.DisplayText,
		game.GameType.NormalizedText,
	)
	if err != nil {
		return fmt.Errorf("error updating game in library: %w", err)
	}

	// Delete existing locations
	_, err = la.db.ExecContext(
		ctx,
		DeleteLibraryLocationsQuery,
		game.GameID,
	)
	if err != nil {
		return fmt.Errorf("error deleting existing locations: %w", err)
	}

	// Insert new locations
	for _, location := range game.PlatformLocations {
		_, err = la.db.ExecContext(
			ctx,
			InsertLibraryLocationQuery,
			game.GameID,
			location.PlatformID,
			location.PlatformName,
			location.Type,
			location.Location.SublocationID,
		)
		if err != nil {
			return fmt.Errorf("error inserting new location: %w", err)
		}
	}

	return nil
}

// POST
func (la *LibraryDbAdapter) CreateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error {
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
		err := tx.QueryRowContext(
			ctx,
			CheckIfGameIsInLibraryQuery,
			userID,
			gameID,
		).Scan(&exists)
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
		_, err = tx.ExecContext(
			ctx,
			CascadingDeleteLibraryGameQuery,
			userID,
			gameID,
		)
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

	err := la.db.GetContext(ctx, &exists, CheckIfGameIsInLibraryQuery, userID, gameID)
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

func (la *LibraryDbAdapter) GetLibraryBFFResponse(
	ctx context.Context,
	userID string,
) (types.LibraryBFFResponseFINAL, error) {
	la.logger.Debug("GetLibraryBFFResponse called", map[string]any{
		"userID": userID,
	})

	// Start transaction
	tx, err := la.db.BeginTxx(ctx, nil)
	if err != nil {
		return types.LibraryBFFResponseFINAL{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Get games data
	var games []models.LibraryGameDB
	if err := tx.SelectContext(ctx, &games, GetLibraryGamesBFFQuery, userID); err != nil {
		return types.LibraryBFFResponseFINAL{}, fmt.Errorf("error querying games: %w", err)
	}

	// Get locations data
	var locations []models.GameLocationDatabaseEntry
	if err := tx.SelectContext(ctx, &locations, GetLibraryLocationsQuery, userID); err != nil {
		return types.LibraryBFFResponseFINAL{}, fmt.Errorf("error querying locations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return types.LibraryBFFResponseFINAL{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Transform locations
	transformedLocations := make([]types.LibraryGamesByPlatformAndLocationItemFINAL, len(locations))
	for i, loc := range locations {
		transformedLocations[i] = TransformLocationDBToResponse(loc)
	}

	// Group locations by game ID
	locationsByGame := make(map[int64][]types.LibraryGamesByPlatformAndLocationItemFINAL)
	for _, loc := range transformedLocations {
		locationsByGame[loc.ID] = append(locationsByGame[loc.ID], loc)
	}

	// Transform games and build response
	response := types.LibraryBFFResponseFINAL{
		LibraryItems:  make([]types.LibraryGameItemBFFResponseFINAL, 0, len(games)),
		RecentlyAdded: make([]types.LibraryGameItemBFFResponseFINAL, 0),
	}

	// Calculate 6-month cutoff
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)

	// Process games
	for _, game := range games {
		gameItem := TransformGameDBToResponse(game, locationsByGame[game.ID])
		response.LibraryItems = append(response.LibraryItems, gameItem)

		// Check if game is recently added
		isRecentlyAdded := false
		for _, location := range locationsByGame[game.ID] {
			if location.DateAdded > sixMonthsAgo.Unix() {
				isRecentlyAdded = true
				break
			}
		}
		if isRecentlyAdded {
			response.RecentlyAdded = append(response.RecentlyAdded, gameItem)
		}
	}

	la.logger.Debug("GetLibraryBFFResponse success", map[string]any{
		"libraryItemsCount":  len(response.LibraryItems),
		"recentlyAddedCount": len(response.RecentlyAdded),
	})

	return response, nil
}

// GetUserLibraryItems retrieves all games in a user's library
func (la *LibraryDbAdapter) GetUserLibraryItems(ctx context.Context, userID string) ([]models.GameToSave, error) {
	la.logger.Debug("LibraryDbAdapter - GetUserLibraryItems called", map[string]any{
		"userID": userID,
	})

	// Get all games in user's library
	var games []models.LibraryGameDB
	err := la.db.SelectContext(
		ctx,
		&games,
		GetLibraryGamesBFFQuery,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting games from library: %w", err)
	}

	// Get locations data
	var locations []models.GameLocationDatabaseEntry
	err = la.db.SelectContext(
		ctx,
		&locations,
		GetLibraryLocationsQuery,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting locations data: %w", err)
	}

	// Group locations by game ID
	locationsByGame := make(map[int64][]models.GameLocationDatabaseEntry)
	for _, loc := range locations {
		locationsByGame[loc.GameID] = append(locationsByGame[loc.GameID], loc)
	}

	// Transform games to GameToSave format
	result := make([]models.GameToSave, len(games))
	for i, game := range games {
		gameLocations := locationsByGame[game.ID]
		platformLocations := make([]models.GameToSaveLocation, len(gameLocations))

		for j, loc := range gameLocations {
			platformLocations[j] = models.GameToSaveLocation{
				PlatformID:   loc.PlatformID,
				PlatformName: loc.PlatformName,
				Type:         "physical", // Default to physical since we don't have digital locations in this query
				Location: models.GameToSaveLocationDetails{
					SublocationID: loc.SublocationID.String,
				},
			}
		}

		result[i] = models.GameToSave{
			GameID:               game.ID,
			GameName:             game.Name,
			GameCoverURL:         game.CoverURL,
			GameFirstReleaseDate: 0, // Not available in this query
			GameType: models.GameToSaveIGDBType{
				DisplayText:    game.GameTypeDisplayText,
				NormalizedText: game.GameTypeNormalizedText,
			},
			PlatformLocations: platformLocations,
		}
	}

	return result, nil
}

// -- REFACTORED LIBRARY BFF RESPONSE METHODS DELETE LEGACY RESPONSE METHODS WHEN COMPLETE --
func (la *LibraryDbAdapter) GetPhysicalLocations(
	ctx context.Context,
	userID string,
) ([]models.PhysicalLocationDB, error) {
	la.logger.Debug("GetPhysicalLocations called", map[string]any{
			"userID": userID,
	})

	var locations []models.PhysicalLocationDB
	err := la.db.SelectContext(
		ctx,
		&locations,
		GetPhysicalLocationsRefactoredQuery,
		userID,
	)
	if err != nil {
			return nil, fmt.Errorf("error querying physical locations: %w", err)
	}

	la.logger.Debug("GetPhysicalLocations success", map[string]any{
			"userID": userID,
			"locationCount": len(locations),
	})

	return locations, nil
}


func (la *LibraryDbAdapter) GetDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocationDB, error) {
	la.logger.Debug("GetDigitalLocations called", map[string]any{
			"userID": userID,
	})

	var locations []models.DigitalLocationDB
	err := la.db.SelectContext(
		ctx,
		&locations,
		GetDigitalLocationsRefactoredQuery,
		userID,
	)
	if err != nil {
			return nil, fmt.Errorf("error querying digital locations: %w", err)
	}

	la.logger.Debug("GetDigitalLocations success", map[string]any{
			"userID": userID,
			"locationCount": len(locations),
	})

	return locations, nil
}


func (la *LibraryDbAdapter) GetGamesMetadata(
	ctx context.Context,
	userID string,
) ([]models.LibraryGameRefactoredDB, error) {
	la.logger.Debug("GetGamesMetadata called", map[string]any{
			"userID": userID,
	})

	// Use raw query to manually scan the array
	rows, err := la.db.QueryxContext(ctx, GetGamesMetadataRefactoredQuery, userID)
	if err != nil {
			return nil, fmt.Errorf("error querying games metadata: %w", err)
	}
	defer rows.Close()

	var games []models.LibraryGameRefactoredDB
	for rows.Next() {
			var game models.LibraryGameRefactoredDB
			var genreNames []string

			err := rows.Scan(
					&game.ID,
					&game.Name,
					&game.CoverURL,
					&game.FirstReleaseDate,
					&game.Rating,
					&game.GameTypeDisplayText,
					&game.GameTypeNormalizedText,
					&game.Favorite,
					&game.CreatedAt,
					&game.IsInWishlist,
					pq.Array(&genreNames),
			)
			if err != nil {
					return nil, fmt.Errorf("error scanning game metadata: %w", err)
			}

			// Manually assign the scanned array to the struct field
			game.GenreNames = genreNames
			games = append(games, game)
	}

	if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating games metadata: %w", err)
	}

	la.logger.Debug("GetGamesMetadata success", map[string]any{
			"userID": userID,
			"gameCount": len(games),
	})

	return games, nil
}

func (la *LibraryDbAdapter) GetLibraryRefactoredBFFResponse(
	ctx context.Context,
	userID string,
) (types.LibraryBFFRefactoredResponse, error) {
	la.logger.Debug("GetLibraryRefactoredBFFResponse called", map[string]any{
			"userID": userID,
	})

	// Execute queries sequentially
	games, err := la.GetGamesMetadata(ctx, userID)
	if err != nil {
			return types.LibraryBFFRefactoredResponse{}, fmt.Errorf("error getting games metadata: %w", err)
	}

	physicalLocations, err := la.GetPhysicalLocations(ctx, userID)
	if err != nil {
			return types.LibraryBFFRefactoredResponse{}, fmt.Errorf("error getting physical locations: %w", err)
	}

	digitalLocations, err := la.GetDigitalLocations(ctx, userID)
	if err != nil {
			return types.LibraryBFFRefactoredResponse{}, fmt.Errorf("error getting digital locations: %w", err)
	}

	// Transform to final response
	response := la.TransformToRefactoredResponse(games, physicalLocations, digitalLocations)

	la.logger.Debug("GetLibraryRefactoredResponse success", map[string]any{
			"userID": userID,
			"libraryItemsCount": len(response.LibraryItems),
			"recentlyAddedCount": len(response.RecentlyAdded),
	})

	return response, nil
}