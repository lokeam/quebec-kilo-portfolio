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
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/postgres"
)

var ErrUnauthorizedLocation = errors.New("unauthorized: sublocation does not belong to user")

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
func (sa *SublocationDbAdapter) GetSublocation(ctx context.Context, userID string, sublocationID string) (models.Sublocation, error) {
	sa.logger.Debug("GetSublocation called", map[string]any{
		"userID":        userID,
		"sublocationID": sublocationID,
	})

	query := `
		SELECT id, user_id, physical_location_id, name, location_type, bg_color, stored_items, created_at, updated_at
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

	sa.logger.Debug("GetSublocation success", map[string]any{
		"sublocation": sublocation,
	})

	return sublocation, nil
}

func (sa *SublocationDbAdapter) GetUserSublocations(ctx context.Context, userID string) ([]models.Sublocation, error) {
	sa.logger.Debug("GetUserSublocations called", map[string]any{
		"userID": userID,
	})

	query := `
		SELECT id, user_id, physical_location_id, name, location_type, bg_color, stored_items, created_at, updated_at
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
	for i := range sublocations {
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

	sa.logger.Debug("GetUserSublocations success", map[string]any{
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
		SET name = $1, location_type = $2, bg_color = $3, stored_items = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7
  `

	now := time.Now()
	result, err := sa.db.ExecContext(
		ctx,
		query,
		sublocation.Name,
		sublocation.LocationType,
		sublocation.BgColor,
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
func (sa *SublocationDbAdapter) AddSublocation(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error) {
	sa.logger.Debug("AddSublocation called", map[string]any{
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
		"bgColor": sublocation.BgColor,
		"storedItems": sublocation.StoredItems,
	})

	query := `
		INSERT INTO sublocations (id, user_id, physical_location_id, name, location_type, bg_color, stored_items, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, physical_location_id, name, location_type, bg_color, stored_items, created_at, updated_at
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
		sublocation.BgColor,
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

	sa.logger.Debug("AddSublocation success", map[string]any{
		"createdSublocation": createdSublocation,
	})

	return createdSublocation, nil
}

// DELETE
func (sa *SublocationDbAdapter) RemoveSublocation(ctx context.Context, userID string, sublocationID string) error {
	sa.logger.Debug("RemoveSublocation called", map[string]any{
		"userID":        userID,
		"sublocationID": sublocationID,
	})

	return postgres.WithTransaction(
		ctx,
		sa.db,
		sa.logger,
		func(tx *sqlx.Tx) error {
			// Check if sublocation exists AND belongs to user
			checkQuery := `
				SELECT id FROM sublocations
				WHERE id = $1 AND user_id = $2
			`

			var id string
			err := tx.QueryRowxContext(ctx, checkQuery, sublocationID, userID).Scan(&id)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("sublocation not found or does not belong to user")
				}
				return fmt.Errorf("error checking sublocation: %w", err)
			}

			// Delete the sublocation
			deleteQuery := `
				DELETE FROM sublocations
				WHERE id = $1 AND user_id = $2
			`

			result, err := tx.ExecContext(ctx, deleteQuery, sublocationID, userID)
			if err != nil {
				return fmt.Errorf("error deleting sublocation: %w", err)
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("error getting rows affected: %w", err)
			}

			if rowsAffected == 0 {
				return fmt.Errorf("sublocation not found or not deleted")
			}

			return nil
		})
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

func (sa *SublocationDbAdapter) RemoveGameFromSublocation(ctx context.Context, userID string, gameID string, sublocationID string) error {
	sa.logger.Debug("RemoveGameFromSublocation called", map[string]any{
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

			// Remove the relationship
			deleteQuery := `DELETE FROM game_sub_locations WHERE game_id = $1 AND sub_location_id = $2`
			result, err := tx.ExecContext(ctx, deleteQuery, gameID, sublocationID)
			if err != nil {
				return fmt.Errorf("error removing game from sublocation: %w", err)
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("error getting rows affected: %w", err)
			}

			if rowsAffected == 0 {
				return fmt.Errorf("game-sublocation relationship not found or not removed")
			}

			return nil
		})
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