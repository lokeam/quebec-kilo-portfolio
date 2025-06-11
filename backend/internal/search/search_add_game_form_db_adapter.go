package search

import (
	"context"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/postgres"
	"github.com/lokeam/qko-beta/internal/types"
)

type SearchAddGameFormDbAdapter struct {
	client   *postgres.PostgresClient
	db       *sqlx.DB
	logger   interfaces.Logger
}

const (
	getAllPhysicalLocationsQuery = `
			SELECT
					pl.id as parent_location_id,
					pl.name as parent_location_name,
					pl.location_type as parent_location_type,
					pl.bg_color as parent_location_bg_color,
					sl.id as sublocation_id,
					sl.name as sublocation_name,
					sl.location_type as sublocation_type
			FROM physical_locations pl
			LEFT JOIN sublocations sl ON sl.physical_location_id = pl.id
			WHERE pl.user_id = $1
			ORDER BY pl.name, sl.name
	`

	getAllDigitalLocationsQuery = `
			SELECT
					id as digital_location_id,
					name as digital_location_name,
					is_subscription,
					is_active
			FROM digital_locations
			WHERE user_id = $1
			ORDER BY name
	`
)

func NewSearchAddGameFormDBAdapter(appContext *appcontext.AppContext) (*SearchAddGameFormDbAdapter, error) {
	appContext.Logger.Debug("Creating SearchAddGameFormDbAdapter", map[string]any{"appContext": appContext})

	// Create a PostgresClient
	client, err := postgres.NewPostgresClient(appContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres client in search db adapter: %w", err)
	}


	// Create sql from px pool
	db, err := sqlx.Connect("pgx", appContext.Config.Postgres.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlx connection in search db adapter: %w", err)
	}

	// Register custom types for PostgreSQL arrays so sqlx can handle string array types
	db.MapperFunc(strings.ToLower)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &SearchAddGameFormDbAdapter{
		client: client,
		db:     db,
		logger: appContext.Logger,
	}, nil
}

func (a *SearchAddGameFormDbAdapter) GetAllGameStorageLocationsBFF(ctx context.Context, userID string) (types.AddGameFormStorageLocationsResponse, error) {
	a.logger.Debug("GetAllGameStorageLocationsBFF called", map[string]any{
			"userID": userID,
	})

	// Start transaction
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
			return types.AddGameFormStorageLocationsResponse{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Initialize with empty slices instead of nil
	response := types.AddGameFormStorageLocationsResponse{
			PhysicalLocations: make([]types.AddGameFormPhysicalLocationsResponse, 0),
			DigitalLocations:  make([]types.AddGameFormDigitalLocationsResponse, 0),
	}

	// Get physical locations
	physicalRows, err := tx.QueryContext(ctx, getAllPhysicalLocationsQuery, userID)
	if err != nil {
			return types.AddGameFormStorageLocationsResponse{}, fmt.Errorf("failed to query physical locations: %w", err)
	}
	defer physicalRows.Close()

	for physicalRows.Next() {
			var loc types.AddGameFormPhysicalLocationsResponse
			err := physicalRows.Scan(
					&loc.ParentLocationID,
					&loc.ParentLocationName,
					&loc.ParentLocationType,
					&loc.ParentLocationBgColor,
					&loc.SublocationID,
					&loc.SublocationName,
					&loc.SublocationType,
			)
			if err != nil {
					return types.AddGameFormStorageLocationsResponse{}, fmt.Errorf("failed to scan physical location: %w", err)
			}

			// Unescape HTML entities in the names
			loc.ParentLocationName = html.UnescapeString(loc.ParentLocationName)
			loc.SublocationName = html.UnescapeString(loc.SublocationName)

			response.PhysicalLocations = append(response.PhysicalLocations, loc)
	}

	if err = physicalRows.Err(); err != nil {
			return types.AddGameFormStorageLocationsResponse{}, fmt.Errorf("error iterating physical location rows: %w", err)
	}

	// Get digital locations
	digitalRows, err := tx.QueryContext(ctx, getAllDigitalLocationsQuery, userID)
	if err != nil {
			return types.AddGameFormStorageLocationsResponse{}, fmt.Errorf("failed to query digital locations: %w", err)
	}
	defer digitalRows.Close()

	for digitalRows.Next() {
			var loc types.AddGameFormDigitalLocationsResponse
			err := digitalRows.Scan(
					&loc.DigitalLocationID,
					&loc.DigitalLocationName,
					&loc.IsSubscription,
					&loc.IsActive,
			)
			if err != nil {
					return types.AddGameFormStorageLocationsResponse{}, fmt.Errorf("failed to scan digital location: %w", err)
			}
			response.DigitalLocations = append(response.DigitalLocations, loc)
	}

	if err = digitalRows.Err(); err != nil {
			return types.AddGameFormStorageLocationsResponse{}, fmt.Errorf("error iterating digital location rows: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
			return types.AddGameFormStorageLocationsResponse{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return response, nil
}