package sublocation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/postgres"
	"github.com/lokeam/qko-beta/internal/types"
)

var ErrUnauthorizedLocation = errors.New("unauthorized: sublocation does not belong to user")

const (
	// DeleteSublocation queries
	queryFindOrphanedGames = `
		WITH games_in_sublocation AS (
			SELECT DISTINCT ug.id
			FROM user_games ug
			JOIN physical_game_locations pgl ON ug.id = pgl.user_game_id
			WHERE pgl.sublocation_id = ANY($1)
			AND ug.user_id = $2
		)
		SELECT
			ug.id as user_game_id,
			g.id as game_id,
			g.name as game_name,
			p.name as platform_name
		FROM games_in_sublocation gis
		JOIN user_games ug ON gis.id = ug.id
		JOIN games g ON ug.game_id = g.id
		JOIN platforms p ON ug.platform_id = p.id
		LEFT JOIN physical_game_locations pgl ON ug.id = pgl.user_game_id
		LEFT JOIN digital_game_locations dgl ON ug.id = dgl.user_game_id
		WHERE pgl.id IS NULL AND dgl.id IS NULL
	`

	queryDeleteSublocations = `
		DELETE FROM sublocations
		WHERE id = ANY($1) AND user_id = $2
		RETURNING id
	`

	queryDeleteOrphanedGames = `
		DELETE FROM user_games
		WHERE id = ANY($1)
	`
)

type SublocationDbAdapter struct {
	client  *postgres.PostgresClient
	db      *sqlx.DB
	logger  interfaces.Logger
}

func NewSublocationDbAdapter(appContext *appcontext.AppContext) (*SublocationDbAdapter, error) {
	appContext.Logger.Debug("Creating SublocationDbAdapter", map[string]any{"appContext": appContext})

	// Create a PostgresClient
	client, err := postgres.NewPostgresClient(appContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres client %w", err)
	}

	// Create sqlx from px pool
	db, err := sqlx.Connect("pgx", appContext.Config.Postgres.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlx connection: %w", err)
	}

	// Register custom types for PostgreSQL array so sqlx can handle string array types
	db.MapperFunc(strings.ToLower)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &SublocationDbAdapter{
		client: client,
		db:     db,
		logger: appContext.Logger,
	}, nil
}

// GET
func (sa *SublocationDbAdapter) GetSingleSublocation(ctx context.Context, userID string, sublocationID string) (models.Sublocation, error) {
	sa.logger.Debug("GetSingleSublocation called", map[string]any{
		"userID":        userID,
		"sublocationID": sublocationID,
	})

	query := `
		SELECT id, user_id, physical_location_id, name, location_type, stored_items, created_at, updated_at
		FROM sublocations
		WHERE id = $1 AND user_id = $2
	`

	var sublocation models.Sublocation
	err := sa.db.GetContext(ctx, &sublocation, query, sublocationID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Sublocation{}, fmt.Errorf("sublocation not found: %w", err)
		}

		return models.Sublocation{}, fmt.Errorf("error getting sublocation: %w", err)
	}

	// Get games for this sublocation
	games, err := sa.GetGamesBySublocationID(ctx, userID, sublocationID)
	if err != nil {
		sa.logger.Error("Error getting games for sublocation", map[string]any{
			"error":        err,
			"sublocationID": sublocationID,
		})
		// Continue without games rather than failing the whole request
	} else {
		sublocation.Items = games
	}

	sa.logger.Debug("GetSingleSublocation success", map[string]any{
		"sublocation": sublocation,
	})

	return sublocation, nil
}

func (sa *SublocationDbAdapter) GetAllSublocations(ctx context.Context, userID string) ([]models.Sublocation, error) {
	sa.logger.Debug("GetAllSublocations called", map[string]any{
		"userID": userID,
	})

	query := `
		SELECT id, user_id, physical_location_id, name, location_type, stored_items, created_at, updated_at
		FROM sublocations
		WHERE user_id = $1
		ORDER BY name
	`

	var sublocations []models.Sublocation
	err := sa.db.SelectContext(ctx, &sublocations, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user sublocations: %w", err)
	}

	// For each sublocation, get its games
	for i := 0; i < len(sublocations); i++ {
		games, err := sa.GetGamesBySublocationID(ctx, userID, sublocations[i].ID)
		if err != nil {
			sa.logger.Error("Error getting games for sublocation", map[string]any{
				"error":        err,
				"sublocationID": sublocations[i].ID,
			})
			// Continue without games rather than failing the whole request
			continue
		}
		sublocations[i].Items = games
	}

	sa.logger.Debug("GetAllSublocations success", map[string]any{
		"sublocations": sublocations,
	})

	return sublocations, nil
}

// Get games for a specific sublocation
func (sa *SublocationDbAdapter) GetGamesBySublocationID(ctx context.Context, userID string, sublocationID string) ([]models.Game, error) {
	sa.logger.Debug("GetGamesBySublocationID called", map[string]any{
		"userID":        userID,
		"sublocationID": sublocationID,
	})

	query := `
		SELECT g.*
		FROM games g
		JOIN game_sub_locations gsl ON g.id = gsl.game_id
		WHERE gsl.sub_location_id = $1 AND g.user_id = $2
	`

	var games []models.Game
	err := sa.db.SelectContext(ctx, &games, query, sublocationID, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting games for sublocation: %w", err)
	}

	sa.logger.Debug("GetGamesBySublocationID success", map[string]any{
		"games": games,
	})

	return games, nil
}

// PUT
func (sa *SublocationDbAdapter) UpdateSublocation(ctx context.Context, userID string, sublocation models.Sublocation) error {
	sa.logger.Debug("UpdateSublocation called", map[string]any{
		"userID":        userID,
		"sublocationID": sublocation.ID,
	})

	// Ensure sublocation belongs to user
	if sublocation.UserID != userID {
		return ErrUnauthorizedLocation
	}

	query := `
		UPDATE sublocations
		SET name = $1, location_type = $2, stored_items = $3, updated_at = $4
		WHERE id = $5 AND user_id = $6
  `

	now := time.Now()
	result, err := sa.db.ExecContext(
		ctx,
		query,
		sublocation.Name,
		sublocation.LocationType,
		sublocation.StoredItems,
		now,
		sublocation.ID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("error updating sublocation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("sublocation not found or not updated")
	}

	sa.logger.Debug("UpdateSublocation success", map[string]any{
		"rowsAffected": rowsAffected,
	})

	return nil
}

// POST
func (sa *SublocationDbAdapter) CreateSublocation(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error) {
	sa.logger.Debug("CreateSublocation called", map[string]any{
		"userID": userID,
		"sublocation": sublocation,
	})

	// Generate a new UUID if ID is not provided
	if sublocation.ID == "" {
		sublocation.ID = uuid.New().String()
	}

	// Set the user ID
	sublocation.UserID = userID

	// Set timestamps
	now := time.Now()
	sublocation.CreatedAt = now
	sublocation.UpdatedAt = now

	// Validate the physical location ID is a valid UUID
	_, err := uuid.Parse(sublocation.PhysicalLocationID)
	if err != nil {
		return models.Sublocation{}, fmt.Errorf("invalid physical_location_id: %w", err)
	}

	sa.logger.Debug("Executing query", map[string]any{
		"id": sublocation.ID,
		"userID": userID,
		"physicalLocationID": sublocation.PhysicalLocationID,
		"name": sublocation.Name,
		"locationType": sublocation.LocationType,
		"storedItems": sublocation.StoredItems,
	})

	query := `
		INSERT INTO sublocations (id, user_id, physical_location_id, name, location_type, stored_items, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, physical_location_id, name, location_type, stored_items, created_at, updated_at
  `

	var createdSublocation models.Sublocation
	err = sa.db.QueryRowxContext(
		ctx,
		query,
		sublocation.ID,
		userID,
		sublocation.PhysicalLocationID,
		sublocation.Name,
		sublocation.LocationType,
		sublocation.StoredItems,
		sublocation.CreatedAt,
		sublocation.UpdatedAt,
	).StructScan(&createdSublocation)

	if err != nil {
		sa.logger.Error("Error executing query", map[string]any{
			"error": err,
			"sublocation": sublocation,
		})
		return models.Sublocation{}, fmt.Errorf("error adding sublocation: %w", err)
	}

	// Initialize empty items slice
	createdSublocation.Items = []models.Game{}

	sa.logger.Debug("CreateSublocation success", map[string]any{
		"createdSublocation": createdSublocation,
	})

	return createdSublocation, nil
}

// DELETE
func (sa *SublocationDbAdapter) DeleteSublocation(
	ctx context.Context,
	userID string,
	sublocationIDs []string,
) (types.DeleteSublocationResponse, error) {
	response := types.DeleteSublocationResponse{
		Success:        false,
		DeletedCount:   0,
		SublocationIDs: make([]string, 0),
		DeletedGames:   make([]types.DeletedGameDetails, 0),
	}

	err := postgres.WithTransaction(ctx, sa.db, sa.logger, func(tx *sqlx.Tx) error {
		// 1. Get games in sublocation that will be orphaned
		var orphanedGames []types.DeletedGameDetails
		if err := tx.SelectContext(
			ctx,
			&orphanedGames,
			queryFindOrphanedGames,
			pq.Array(sublocationIDs),
			userID,
		); err != nil {
			return fmt.Errorf("error finding orphaned games: %w", err)
		}

		// 2. Delete sublocations
		rows, err := tx.QueryContext(ctx, queryDeleteSublocations, pq.Array(sublocationIDs), userID)
		if err != nil {
			return fmt.Errorf("error deleting sublocations: %w", err)
		}
		defer rows.Close()

		// Collect deleted IDs
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return fmt.Errorf("error scanning deleted ID: %w", err)
			}
			response.SublocationIDs = append(response.SublocationIDs, id)
		}

		if err := rows.Err(); err != nil {
			return fmt.Errorf("error iterating delete results: %w", err)
		}

		// 3. Delete orphaned games
		if len(orphanedGames) > 0 {
			userGameIDs := make([]int, len(orphanedGames))
			for i, game := range orphanedGames {
				userGameIDs[i] = game.UserGameID
			}

			if _, err := tx.ExecContext(ctx, queryDeleteOrphanedGames, pq.Array(userGameIDs)); err != nil {
				return fmt.Errorf("error deleting orphaned games: %w", err)
			}
		}

		// Update response
		response.Success = true
		response.DeletedCount = len(response.SublocationIDs)
		response.DeletedGames = orphanedGames

		// Log the deletion
		sa.logger.Info("Successfully deleted sublocations and orphaned games", map[string]any{
			"user_id": userID,
			"deleted_count": response.DeletedCount,
			"deleted_ids": response.SublocationIDs,
			"deleted_games_count": len(response.DeletedGames),
			"deleted_games": response.DeletedGames,
		})

		return nil
	})

	if err != nil {
		return response, err
	}

	return response, nil
}

// Game-Sublocation relationship management
func (sa *SublocationDbAdapter) AddGameToSublocation(ctx context.Context, userID string, gameID string, sublocationID string) error {
	sa.logger.Debug("AddGameToSublocation called", map[string]any{
		"userID":        userID,
		"gameID":        gameID,
		"sublocationID": sublocationID,
	})

	return postgres.WithTransaction(
		ctx,
		sa.db,
		sa.logger,
		func(tx *sqlx.Tx) error {
			// Verify game belongs to user
			checkGameQuery := `SELECT id FROM games WHERE id = $1 AND user_id = $2`
			var gameIDResult string
			err := tx.QueryRowxContext(ctx, checkGameQuery, gameID, userID).Scan(&gameIDResult)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("game not found or does not belong to user")
				}
				return fmt.Errorf("error checking game: %w", err)
			}

			// Verify sublocation belongs to user
			checkSublocQuery := `SELECT id FROM sublocations WHERE id = $1 AND user_id = $2`
			var sublocIDResult string
			err = tx.QueryRowxContext(ctx, checkSublocQuery, sublocationID, userID).Scan(&sublocIDResult)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("sublocation not found or does not belong to user")
				}
				return fmt.Errorf("error checking sublocation: %w", err)
			}

			// Add the relationship
			insertQuery := `
				INSERT INTO game_sub_locations (game_id, sub_location_id)
				VALUES ($1, $2)
				ON CONFLICT (game_id, sub_location_id) DO NOTHING
			`
			_, err = tx.ExecContext(ctx, insertQuery, gameID, sublocationID)
			if err != nil {
				return fmt.Errorf("error adding game to sublocation: %w", err)
			}

			return nil
		})
}

func (sa *SublocationDbAdapter) RemoveGameFromSublocation(ctx context.Context, userID string, userGameID string) error {
	sa.logger.Debug("RemoveGameFromSublocation called", map[string]any{
		"userID":     userID,
		"userGameID": userGameID,
	})

	return postgres.WithTransaction(ctx, sa.db, sa.logger, func(tx *sqlx.Tx) error {
		// 1. Verify game belongs to user
		checkGameQuery := `SELECT id FROM user_games WHERE id = $1 AND user_id = $2`
		var gameIDResult string
		err := tx.QueryRowxContext(ctx, checkGameQuery, userGameID, userID).Scan(&gameIDResult)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("game not found or does not belong to user")
			}
			return fmt.Errorf("error checking game: %w", err)
		}

		// 2. Remove game from current sublocation
		deleteQuery := `DELETE FROM physical_game_locations WHERE user_game_id = $1`
		result, err := tx.ExecContext(ctx, deleteQuery, userGameID)
		if err != nil {
			return fmt.Errorf("error removing game from sublocation: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("error getting rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("game is not in any sublocation")
		}

		return nil
	})
}

// CheckGameInAnySublocation checks if a game is in any sublocation
func (sa *SublocationDbAdapter) CheckGameInAnySublocation(ctx context.Context, userGameID string) (bool, error) {
	sa.logger.Debug("CheckGameInAnySublocation called", map[string]any{
		"userGameID": userGameID,
	})

	query := `
		SELECT EXISTS (
			SELECT 1 FROM physical_game_locations
			WHERE user_game_id = $1
		)
	`

	var exists bool
	err := sa.db.GetContext(ctx, &exists, query, userGameID)
	if err != nil {
		return false, fmt.Errorf("error checking game location: %w", err)
	}

	return exists, nil
}

// CheckGameInSublocation checks if a game is in a specific sublocation
func (sa *SublocationDbAdapter) CheckGameInSublocation(ctx context.Context, userGameID string, sublocationID string) (bool, error) {
	sa.logger.Debug("CheckGameInSublocation called", map[string]any{
		"userGameID":    userGameID,
		"sublocationID": sublocationID,
	})

	query := `
		SELECT EXISTS (
			SELECT 1 FROM physical_game_locations
			WHERE user_game_id = $1 AND sublocation_id = $2
		)
	`

	var exists bool
	err := sa.db.GetContext(ctx, &exists, query, userGameID, sublocationID)
	if err != nil {
		return false, fmt.Errorf("error checking game location: %w", err)
	}

	return exists, nil
}

// CheckGameOwnership checks if a user owns a specific game
func (sa *SublocationDbAdapter) CheckGameOwnership(ctx context.Context, userID string, userGameID string) (bool, error) {
	sa.logger.Debug("CheckGameOwnership called", map[string]any{
		"userID":     userID,
		"userGameID": userGameID,
	})

	query := `
		SELECT EXISTS (
			SELECT 1 FROM user_games
			WHERE id = $1 AND user_id = $2
		)
	`

	var exists bool
	err := sa.db.GetContext(ctx, &exists, query, userGameID, userID)
	if err != nil {
		return false, fmt.Errorf("error checking game ownership: %w", err)
	}

	return exists, nil
}

// Check if sublocation with same name exists in the same physical location
func (sa *SublocationDbAdapter) CheckDuplicateSublocation(ctx context.Context, userID string, physicalLocationID string, name string) (bool, error) {
	sa.logger.Debug("CheckDuplicateSublocation called", map[string]any{
		"userID": userID,
		"physicalLocationID": physicalLocationID,
		"name": name,
	})

	query := `
		SELECT COUNT(*)
		FROM sublocations
		WHERE user_id = $1 AND physical_location_id = $2 AND LOWER(name) = LOWER($3)
	`

	var count int
	err := sa.db.GetContext(ctx, &count, query, userID, physicalLocationID, name)
	if err != nil {
		return false, fmt.Errorf("error checking for duplicate sublocation: %w", err)
	}

	return count > 0, nil
}

// MoveGameToSublocation moves a game to a new sublocation in a transaction
func (sa *SublocationDbAdapter) MoveGameToSublocation(ctx context.Context, userID string, userGameID string, targetSublocationID string) error {
	sa.logger.Debug("MoveGameToSublocation called", map[string]any{
		"userID":        userID,
		"userGameID":    userGameID,
		"targetSublocationID": targetSublocationID,
	})

	return postgres.WithTransaction(ctx, sa.db, sa.logger, func(tx *sqlx.Tx) error {
		// 1. Lock and get current sublocation
		var currentSublocationID string
		err := tx.QueryRowxContext(ctx,
			`SELECT sublocation_id FROM physical_game_locations
			 WHERE user_game_id = $1 FOR UPDATE`,
			userGameID,
		).Scan(&currentSublocationID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error getting current sublocation: %w", err)
		}

		// 2. Verify game ownership
		checkGameQuery := `SELECT id FROM user_games WHERE id = $1 AND user_id = $2`
		var gameIDResult string
		err = tx.QueryRowxContext(ctx, checkGameQuery, userGameID, userID).Scan(&gameIDResult)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("game not found or does not belong to user")
			}
			return fmt.Errorf("error checking game: %w", err)
		}

		// 3. Verify target sublocation ownership
		checkSublocQuery := `SELECT id FROM sublocations WHERE id = $1 AND user_id = $2`
		var sublocIDResult string
		err = tx.QueryRowxContext(ctx, checkSublocQuery, targetSublocationID, userID).Scan(&sublocIDResult)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("target sublocation not found or does not belong to user")
			}
			return fmt.Errorf("error checking sublocation: %w", err)
		}

		// 4. Move the game
		deleteQuery := `DELETE FROM physical_game_locations WHERE user_game_id = $1`
		_, err = tx.ExecContext(ctx, deleteQuery, userGameID)
		if err != nil {
			return fmt.Errorf("error removing game from current sublocation: %w", err)
		}

		insertQuery := `
			INSERT INTO physical_game_locations (user_game_id, sublocation_id)
			VALUES ($1, $2)
			ON CONFLICT (user_game_id) DO UPDATE
			SET sublocation_id = $2
		`
		_, err = tx.ExecContext(ctx, insertQuery, userGameID, targetSublocationID)
		if err != nil {
			return fmt.Errorf("error adding game to new sublocation: %w", err)
		}

		return nil
	})
}