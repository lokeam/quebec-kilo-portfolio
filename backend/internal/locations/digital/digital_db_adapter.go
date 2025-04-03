package digital

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/postgres"
)

var ErrUnauthorizedLocation = errors.New("unauthorized: location does not belong to user")

type DigitalDbAdapter struct {
	client   *postgres.PostgresClient
	db       *sqlx.DB
	logger   interfaces.Logger
}

func NewDigitalDbAdapter(appContext *appcontext.AppContext) (*DigitalDbAdapter, error) {
	appContext.Logger.Debug("Creating DigitalDbAdapter", map[string]any{"appContext": appContext})

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

	return &DigitalDbAdapter{
		client: client,
		db:     db,
		logger: appContext.Logger,
	}, nil
}

// GET
func (da *DigitalDbAdapter) GetDigitalLocation(ctx context.Context, userID string, locationID string) (models.DigitalLocation, error) {
	da.logger.Debug("GetDigitalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
	})

	query := `
		SELECT id, user_id, name, is_active, url, features, created_at, updated_at
		FROM digital_locations
		WHERE id = $1 AND user_id = $2
		`

	var location models.DigitalLocation
	err := da.db.GetContext(ctx, &location, query, locationID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.DigitalLocation{}, fmt.Errorf("digital location not found: %w", err)
		}

		return models.DigitalLocation{}, fmt.Errorf("error getting digital location: %w", err)
	}

	da.logger.Debug("GetDigitalLocation success", map[string]any{
		"location": location,
	})

	return location, nil
}

func (da *DigitalDbAdapter) GetUserDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
	da.logger.Debug("GetUserDigitalLocations called", map[string]any{
		"userID": userID,
	})

	query := `
		SELECT id, user_id, name, is_active, url, features, created_at, updated_at
		FROM digital_locations
		WHERE user_id = $1
		ORDER BY name
		`

	var locations []models.DigitalLocation
	err := da.db.SelectContext(ctx, &locations, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user digital locations: %w", err)
	}

	da.logger.Debug("GetUserDigitalLocations success", map[string]any{
		"locations": locations,
	})

	return locations, nil
}

func (da *DigitalDbAdapter) GetGamesByDigitalLocationID(
	ctx context.Context,
	userID string,
	sublocationID string,
) ([]models.Game, error) {
	da.logger.Debug("GetGamesByDigitalLocationID called", map[string]any{
		"userID":         userID,
		"sublocationID":  sublocationID,
	})

	query := `
		SELECT g.*
		FROM games g
		JOIN game_digital_locations gdl ON g.id = gdl.game_id
		WHERE gdl.digital_location_id = $1 AND g.user_id = $2
	`

	var games []models.Game
	err := da.db.SelectContext(
		ctx,
		&games,
		query,
		sublocationID,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting games for digital location: %w", err)
	}

	da.logger.Debug("GetGamesByDigitalLocationID success", map[string]any{
		"games":  games,
	})

	return games, nil
}

// POST
func (a *DigitalDbAdapter) AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
	a.logger.Debug("Adding digital location", map[string]any{"userID": userID, "location": location})

	now := time.Now()
	query := `
			INSERT INTO digital_locations (id, user_id, name, is_active, url, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, user_id, name, is_active, url, created_at, updated_at
	`

	var result models.DigitalLocation
	err := a.db.QueryRowxContext(
			ctx,
			query,
			location.ID,
			userID,
			location.Name,
			location.IsActive,
			location.URL,
			now,
			now,
	).StructScan(&result)

	if err != nil {
			return models.DigitalLocation{}, fmt.Errorf("error adding digital location: %w", err)
	}

	return result, nil
}

// PUT
func (a *DigitalDbAdapter) UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error {
	a.logger.Debug("Updating digital location", map[string]any{"userID": userID, "location": location})

	// Check if the location belongs to the user
	if location.UserID != userID {
			return ErrUnauthorizedLocation
	}

	now := time.Now()
	query := `
			UPDATE digital_locations
			SET name = $1, is_active = $2, url = $3, updated_at = $4
			WHERE id = $5 AND user_id = $6
	`

	result, err := a.db.ExecContext(
			ctx,
			query,
			location.Name,
			location.IsActive,
			location.URL,
			now,
			location.ID,
			userID,
	)

	if err != nil {
			return fmt.Errorf("error updating digital location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
			return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
			return fmt.Errorf("digital location not found or not updated")
	}

	return nil
}

// DELETE
func (da *DigitalDbAdapter) RemoveDigitalLocation(ctx context.Context, userID string, locationID string) error {
	da.logger.Debug("RemoveDigitalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
	})

	return postgres.WithTransaction(
		ctx,
		da.db,
		da.logger,
		func(tx *sqlx.Tx) error {
			// Check if location exists AND belongs to user
			checkQuery := `
				SELECT id FROM digital_locations
				WHERE id = $1 AND user_id = $2
			`

			var id string
			err := tx.QueryRowxContext(ctx, checkQuery, locationID, userID).Scan(&id)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("digital location not found or does not belong to user")
				}
				return fmt.Errorf("error checking digital location: %w", err)
			}

			// Delete the location
			deleteQuery := `
				DELETE FROM digital_locations
				WHERE id = $1 AND user_id = $2
			`

			result, err := tx.ExecContext(ctx, deleteQuery, locationID, userID)
			if err != nil {
				return fmt.Errorf("error deleting digital location: %w", err)
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("error getting rows affected: %w", err)
			}

			if rowsAffected == 0 {
				return fmt.Errorf("digital location not found or not deleted")
			}

			return nil
		})
}
