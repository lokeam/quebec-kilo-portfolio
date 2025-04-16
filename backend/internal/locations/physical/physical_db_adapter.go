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

const (
	createPhysicalLocationQuery = `
		INSERT INTO physical_locations (
			id, user_id, name, label, location_type, map_coordinates
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, name, label, location_type, map_coordinates, created_at, updated_at
	`

	getPhysicalLocationQuery = `
		SELECT id, user_id, name, label, location_type, map_coordinates, created_at, updated_at
		FROM physical_locations
		WHERE id = $1 AND user_id = $2
	`

	getUserPhysicalLocationsQuery = `
		SELECT id, user_id, name, label, location_type, map_coordinates, created_at, updated_at
		FROM physical_locations
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	updatePhysicalLocationQuery = `
		UPDATE physical_locations
		SET name = $3, label = $4, location_type = $5, map_coordinates = $6, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, name, label, location_type, map_coordinates, created_at, updated_at
	`
)

// GET
func (pa *PhysicalDbAdapter) GetPhysicalLocation(ctx context.Context, userID string, locationID string) (models.PhysicalLocation, error) {
	pa.logger.Debug("GetPhysicalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
	})

	var location models.PhysicalLocation
	err := pa.db.QueryRowContext(ctx, getPhysicalLocationQuery, locationID, userID).Scan(
		&location.ID,
		&location.UserID,
		&location.Name,
		&location.Label,
		&location.LocationType,
		&location.MapCoordinates,
		&location.CreatedAt,
		&location.UpdatedAt,
	)
	if err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("failed to get physical location: %w", err)
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

	rows, err := pa.db.QueryContext(ctx, getUserPhysicalLocationsQuery, userID)
	if err != nil {
		return []models.PhysicalLocation{}, fmt.Errorf("failed to get user physical locations: %w", err)
	}
	defer rows.Close()

	var locations []models.PhysicalLocation
	for rows.Next() {
		var location models.PhysicalLocation
		err := rows.Scan(
			&location.ID,
			&location.UserID,
			&location.Name,
			&location.Label,
			&location.LocationType,
			&location.MapCoordinates,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		if err != nil {
			return []models.PhysicalLocation{}, fmt.Errorf("failed to scan physical location: %w", err)
		}
		locations = append(locations, location)
	}

	pa.logger.Debug("GetUserPhysicalLocations success", map[string]any{
		"locations": locations,
	})

	return locations, nil
}

// PUT
func (pa *PhysicalDbAdapter) UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	pa.logger.Debug("Updating physical location", map[string]any{
		"userID": userID,
		"location": location,
	})

	var updatedLocation models.PhysicalLocation
	err := pa.db.QueryRowContext(ctx, updatePhysicalLocationQuery,
		location.ID,
		userID,
		location.Name,
		location.Label,
		location.LocationType,
		location.MapCoordinates,
	).Scan(
		&updatedLocation.ID,
		&updatedLocation.UserID,
		&updatedLocation.Name,
		&updatedLocation.Label,
		&updatedLocation.LocationType,
		&updatedLocation.MapCoordinates,
		&updatedLocation.CreatedAt,
		&updatedLocation.UpdatedAt,
	)
	if err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("failed to update physical location: %w", err)
	}

	pa.logger.Debug("UpdatePhysicalLocation success", map[string]any{
		"location": updatedLocation,
	})

	return updatedLocation, nil
}

// POST
func (pa *PhysicalDbAdapter) ensureUserExists(ctx context.Context, userID string) error {
	pa.logger.Debug("Ensuring user exists", map[string]any{"userID": userID})

	// Check if user exists
	var exists bool
	err := pa.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)
	`, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}

	if !exists {
		// Create user with all required fields
		_, err = pa.db.ExecContext(ctx, `
			INSERT INTO users (id, email, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW())
		`, userID, fmt.Sprintf("%s@example.com", userID))
		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
	}

	return nil
}

func (pa *PhysicalDbAdapter) AddPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	pa.logger.Debug("Adding physical location", map[string]any{"userID": userID, "location": location})

	// Ensure user exists first
	if err := pa.ensureUserExists(ctx, userID); err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("error ensuring user exists: %w", err)
	}

	// Check if a location with this name already exists for the user
	var existingID string
	err := pa.db.QueryRowContext(ctx, `
		SELECT id FROM physical_locations
		WHERE user_id = $1 AND LOWER(name) = LOWER($2)
	`, userID, location.Name).Scan(&existingID)

	if err == nil {
		// Location with this name already exists
		return models.PhysicalLocation{}, fmt.Errorf("a physical location with the name '%s' already exists", location.Name)
	} else if err != sql.ErrNoRows {
		// Some other database error occurred
		return models.PhysicalLocation{}, fmt.Errorf("error checking for existing location: %w", err)
	}

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

	var newLocation models.PhysicalLocation
	err = pa.db.QueryRowContext(ctx, createPhysicalLocationQuery,
		location.ID,
		userID,
		location.Name,
		location.Label,
		location.LocationType,
		location.MapCoordinates,
	).Scan(
		&newLocation.ID,
		&newLocation.UserID,
		&newLocation.Name,
		&newLocation.Label,
		&newLocation.LocationType,
		&newLocation.MapCoordinates,
		&newLocation.CreatedAt,
		&newLocation.UpdatedAt,
	)
	if err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("failed to create physical location: %w", err)
	}

	pa.logger.Debug("AddPhysicalLocation success", map[string]any{
		"location": newLocation,
	})

	return newLocation, nil
}

// DELETE
func (pa *PhysicalDbAdapter) RemovePhysicalLocation(ctx context.Context, userID string, locationID string) error {
	pa.logger.Debug("Removing physical location", map[string]any{
		"userID":     userID,
		"locationID": locationID,
	})

	return postgres.WithTransaction(
		ctx,
		pa.db,
		pa.logger,
		func(tx *sqlx.Tx) error {
			// First try to find the location by ID or name
			var id string
			checkQuery := `
				SELECT id FROM physical_locations
				WHERE (
					(id::text = $1) OR
					(name = $1) OR
					(LOWER(name) = LOWER($1))
				) AND user_id = $2
			`

			err := tx.QueryRowxContext(ctx, checkQuery, locationID, userID).Scan(&id)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("physical location not found or does not belong to user")
				}
				return fmt.Errorf("error checking physical location: %w", err)
			}

			// Delete the location using the found ID
			deleteQuery := `
				DELETE FROM physical_locations
				WHERE id = $1 AND user_id = $2
			`

			result, err := tx.ExecContext(ctx, deleteQuery, id, userID)
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

			pa.logger.Debug("RemovePhysicalLocation success", map[string]any{
				"rowsAffected": rowsAffected,
				"locationID":   id,
			})

			return nil
		})
}

// Helpers