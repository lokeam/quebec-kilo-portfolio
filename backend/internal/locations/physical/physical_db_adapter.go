package physical

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

var ErrUnauthorizedLocation = errors.New("unauthorized: location does not belong to user")

type PhysicalDbAdapter struct {
	client     *postgres.PostgresClient
	db         *sqlx.DB
	logger     interfaces.Logger
	scanner    interfaces.GameScanner
}

func NewPhysicalDbAdapter(appContext *appcontext.AppContext) (*PhysicalDbAdapter, error) {
	appContext.Logger.Debug("Creating PhysicalLibraryDbAdapter", map[string]any{"appContext": appContext})

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

	// Register custom types for PostgreSQL arrays so sqlx can handle string array types
	// NOTE: Do these values need to be repeated in every db adapter?
	// If so, can I separate these out into a shared fn?
	db.MapperFunc(strings.ToLower)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PhysicalDbAdapter{
		client: client,
		db:     db,
		logger: appContext.Logger,
	}, nil
}

// GET
func (pa *PhysicalDbAdapter) GetPhysicalLocation(ctx context.Context, userID string, locationID string) (models.PhysicalLocation, error) {
	pa.logger.Debug("GetPhysicalLocation called", map[string]any{
		"userID": userID,
		"locationID": locationID,
	})

	query := `
		SELECT id, user_id, name, label, location_type, map_coordinates, created_at, updated_at
		FROM physical_locations
		WHERE id = $1 AND user_id = $2
	`

	var location models.PhysicalLocation
	err := pa.db.GetContext(ctx, &location, query, locationID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.PhysicalLocation{}, fmt.Errorf("physical location not found: %w", err)
		}

		return models.PhysicalLocation{}, fmt.Errorf("error getting physical location: %w", err)
	}

	pa.logger.Debug("GetPhysicalLocation success", map[string]any{
		"location": location,
	})

	return location, nil
}

func (pa *PhysicalDbAdapter) GetUserPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
	pa.logger.Debug("GetUserPhysicalLocations called", map[string]any{
		"userID": userID,
	})

	query := `
		SELECT id, user_id, name, label, location_type, map_coordinates, created_at, updated_at
		FROM physical_locations
		WHERE user_id = $1
		ORDER BY name
	`

	var locations []models.PhysicalLocation
	err := pa.db.SelectContext(ctx, &locations, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user physical locations: %w", err)
	}

	pa.logger.Debug("GetUserPhysicalLocations success", map[string]any{
		"locations": locations,
	})

	return locations, nil
}

// PUT
func (pa *PhysicalDbAdapter) UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) error {
	pa.logger.Debug("UpdatePhysicalLocation called", map[string]any{
		"userID": userID,
		"locationID": location.ID,
	})

	// Ensure location belongs to user
	if location.UserID != userID {
    return ErrUnauthorizedLocation
	}

	query := `
		UPDATE physical_locations
		SET name = $1, label = $2, location_type = $3, map_coordinates = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7
	`

	now := time.Now()
	result, err := pa.db.ExecContext(
		ctx,
		query,
		location.Name,
		location.Label,
		location.LocationType,
		location.MapCoordinates,
		now,
		location.ID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("error updating physical location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("physical location not found or not updated")
	}

	pa.logger.Debug("UpdatePhysicalLocation success", map[string]any{
		"rowsAffected": rowsAffected,
	})

	return nil
}

// POST
func (pa *PhysicalDbAdapter) AddPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	pa.logger.Debug("AddPhysicalLocation called", map[string]any{
		"userID": userID,
	})

	// Generate a new UUID if ID is not provided
	if location.ID == "" {
		location.ID = uuid.New().String()
	}

	// Set the user ID
	location.UserID = userID

	// Set timestamps
	now := time.Now()
	location.CreatedAt = now
	location.UpdatedAt = now

	query := `
		INSERT INTO physical_locations (id, user_id, name, label, location_type, map_coordinates, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, name, label, location_type, map_coordinates, created_at, updated_at
	`

	err := pa.db.QueryRowxContext(
		ctx,
		query,
		location.ID,
		userID,
		location.Name,
		location.Label,
		location.LocationType,
		location.MapCoordinates,
		location.CreatedAt,
		location.UpdatedAt,
	).StructScan(&location)

	if err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("error adding physical location: %w", err)
	}

	return location, nil
}

// DELETE
func (pa *PhysicalDbAdapter) RemovePhysicalLocation(ctx context.Context, userID string, locationID string) error {
	pa.logger.Debug("RemovePhysicalLocation called", map[string]any{
		"userID": userID,
		"locationID": locationID,
	})

	return postgres.WithTransaction(
		ctx,
		pa.db,
		pa.logger,
		func(tx *sqlx.Tx) error {
			// Check if location exists AND belongs to user
			checkQuery := `
				SELECT id FROM physical_locations
				WHERE id = $1 AND user_id = $2
			`

			var id string
			err := tx.QueryRowxContext(ctx, checkQuery, locationID, userID).Scan(&id)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("physical location not found or does not belong to user")
				}
				return fmt.Errorf("error checking physical location: %w", err)
			}

			// Delete the location
			deleteQuery := `
				DELETE FROM physical_locations
				WHERE id = $1 AND user_id = $2
			`

			result, err := tx.ExecContext(ctx, deleteQuery, locationID, userID)
			if err != nil {
				return fmt.Errorf("error deleting physical location: %w", err)
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("error getting rows affected: %w", err)
			}

			if rowsAffected == 0 {
				return fmt.Errorf("physical location not found or not deleted")
			}

		return nil
	})
}

// Helpers