package physical

import (
	"context"
	"database/sql"
	"encoding/json"
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

	getPhysicalLocationSublocationsQuery = `
		SELECT json_agg(
			json_build_object(
				'id', sl.id,
				'user_id', sl.user_id,
				'physical_location_id', sl.physical_location_id,
				'name', sl.name,
				'location_type', sl.location_type,
				'bg_color', sl.bg_color,
				'stored_items', sl.stored_items,
				'created_at', sl.created_at,
				'updated_at', sl.updated_at,
				'items', COALESCE(
					(
						SELECT json_agg(
							json_build_object(
								'id', g.id,
								'name', g.name,
								'summary', g.summary,
								'cover_id', g.cover_id,
								'cover_url', g.cover_url,
								'first_release_date', g.first_release_date,
								'rating', g.rating
							)
						)
						FROM games g
						JOIN user_games ug ON ug.game_id = g.id
						JOIN physical_game_locations pgl ON pgl.user_game_id = ug.id
						WHERE pgl.sublocation_id = sl.id
					),
					'[]'::json
				)
			)
		)
		FROM sublocations sl
		WHERE sl.physical_location_id = $1
	`

	getUserPhysicalLocationsQuery = `
		SELECT id, user_id, name, label, location_type, map_coordinates, created_at, updated_at
		FROM physical_locations
		WHERE user_id = $1
		ORDER BY created_at
	`

	getUserPhysicalLocationSublocationsQuery = `
		SELECT json_agg(
			json_build_object(
				'id', sl.id,
				'user_id', sl.user_id,
				'physical_location_id', sl.physical_location_id,
				'name', sl.name,
				'location_type', sl.location_type,
				'bg_color', sl.bg_color,
				'stored_items', sl.stored_items,
				'created_at', sl.created_at,
				'updated_at', sl.updated_at,
				'items', COALESCE(
					(
						SELECT json_agg(
							json_build_object(
								'id', g.id,
								'name', g.name,
								'summary', g.summary,
								'cover_id', g.cover_id,
								'cover_url', g.cover_url,
								'first_release_date', g.first_release_date,
								'rating', g.rating
							)
						)
						FROM games g
						JOIN user_games ug ON ug.game_id = g.id
						JOIN physical_game_locations pgl ON pgl.user_game_id = ug.id
						WHERE pgl.sublocation_id = sl.id
					),
					'[]'::json
				)
			)
		)
		FROM sublocations sl
		WHERE sl.physical_location_id = $1
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

	// Start transaction
	tx, err := pa.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	var location models.PhysicalLocation
	err = tx.QueryRowContext(ctx, getPhysicalLocationQuery, locationID, userID).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.PhysicalLocation{}, ErrLocationNotFound
		}
		return models.PhysicalLocation{}, fmt.Errorf("failed to get physical location: %w", err)
	}

	// Get sublocations
	var subLocationsJSON []byte
	err = tx.QueryRowContext(ctx, getPhysicalLocationSublocationsQuery, locationID).Scan(&subLocationsJSON)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.PhysicalLocation{}, fmt.Errorf("failed to get sublocations: %w", err)
	}

	if len(subLocationsJSON) > 0 {
		var subLocations []models.Sublocation
		if err := json.Unmarshal(subLocationsJSON, &subLocations); err != nil {
			return models.PhysicalLocation{}, fmt.Errorf("failed to unmarshal sublocations: %w", err)
		}
		location.SubLocations = &subLocations
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("failed to commit transaction: %w", err)
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

	// Start transaction
	tx, err := pa.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Get base locations
	rows, err := tx.QueryContext(ctx, getUserPhysicalLocationsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user physical locations: %w", err)
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
			return nil, fmt.Errorf("failed to scan physical location: %w", err)
		}
		locations = append(locations, location)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating physical locations: %w", err)
	}

	// Get sublocations for each location
	for i := range locations {
		var subLocationsJSON []byte
		err := tx.QueryRowContext(ctx, getUserPhysicalLocationSublocationsQuery, locations[i].ID).Scan(&subLocationsJSON)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get sublocations for location %s: %w", locations[i].ID, err)
		}

		if len(subLocationsJSON) > 0 {
			var subLocations []models.Sublocation
			if err := json.Unmarshal(subLocationsJSON, &subLocations); err != nil {
				return nil, fmt.Errorf("failed to unmarshal sublocations for location %s: %w", locations[i].ID, err)
			}
			locations[i].SubLocations = &subLocations
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	pa.logger.Debug("GetUserPhysicalLocations success", map[string]any{
		"locations": locations,
	})

	return locations, nil
}

// POST
func (pa *PhysicalDbAdapter) ensureUserExists(ctx context.Context, userID string) error {
	pa.logger.Debug("Ensuring user exists", map[string]any{"userID": userID})

	var exists bool
	err := pa.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)
	`, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}

	if !exists {
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

	// Start transaction
	tx, err := pa.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Ensure user exists
	if err := pa.ensureUserExists(ctx, userID); err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("error ensuring user exists: %w", err)
	}

	// Check for duplicate location name
	var existingID string
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM physical_locations
		WHERE user_id = $1 AND LOWER(name) = LOWER($2)
	`, userID, location.Name).Scan(&existingID)

	if err == nil {
		return models.PhysicalLocation{}, ErrDuplicateLocation
	} else if err != sql.ErrNoRows {
		return models.PhysicalLocation{}, fmt.Errorf("error checking for existing location: %w", err)
	}

	// Generate a new UUID if ID is not provided
	if location.ID == "" {
		location.ID = uuid.New().String()
	}

	// Set the user ID and timestamps
	location.UserID = userID
	now := time.Now()
	location.CreatedAt = now
	location.UpdatedAt = now

	var newLocation models.PhysicalLocation
	err = tx.QueryRowContext(ctx, createPhysicalLocationQuery,
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

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	pa.logger.Debug("AddPhysicalLocation success", map[string]any{
		"location": newLocation,
	})

	return newLocation, nil
}

// PUT
func (pa *PhysicalDbAdapter) UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	pa.logger.Debug("Updating physical location", map[string]any{
		"userID": userID,
		"location": location,
	})

	// Start transaction
	tx, err := pa.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	var updatedLocation models.PhysicalLocation
	err = tx.QueryRowContext(ctx, updatePhysicalLocationQuery,
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.PhysicalLocation{}, ErrLocationNotFound
		}
		return models.PhysicalLocation{}, fmt.Errorf("failed to update physical location: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	pa.logger.Debug("UpdatePhysicalLocation success", map[string]any{
		"location": updatedLocation,
	})

	return updatedLocation, nil
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
			// First try to find the location by ID
			var id string
			err := tx.QueryRowxContext(ctx, `
				SELECT id FROM physical_locations
				WHERE id = $1 AND user_id = $2
			`, locationID, userID).Scan(&id)

			if err != nil {
				if err == sql.ErrNoRows {
					return ErrLocationNotFound
				}
				return fmt.Errorf("error checking physical location: %w", err)
			}

			// Delete the location
			result, err := tx.ExecContext(ctx, `
				DELETE FROM physical_locations
				WHERE id = $1 AND user_id = $2
			`, id, userID)
			if err != nil {
				return fmt.Errorf("error deleting physical location: %w", err)
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("error getting rows affected: %w", err)
			}

			if rowsAffected == 0 {
				return ErrLocationNotFound
			}

			pa.logger.Debug("RemovePhysicalLocation success", map[string]any{
				"rowsAffected": rowsAffected,
				"locationID":   id,
			})

			return nil
		})
}

// Helpers