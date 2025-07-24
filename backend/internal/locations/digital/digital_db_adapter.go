package digital

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

// DeletionResult tracks the results of a bulk deletion operation
type DeletionResult struct {
	SuccessCount int64
	FailedIDs    []string
	Error        error
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
func (da *DigitalDbAdapter) GetSingleDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
) (models.DigitalLocation, error) {
	da.logger.Debug("GetSingleDigitalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
	})

	// Fetch digital location
	var digitalLocation models.DigitalLocation
	err := da.db.GetContext(
		ctx,
		&digitalLocation,
		GetSingleDigitalLocationQuery,
		locationID,
		userID,
	)
	if err != nil{
		if err == sql.ErrNoRows {
			return models.DigitalLocation{}, fmt.Errorf("digital location not found: %w", err)
		}
		return models.DigitalLocation{}, fmt.Errorf("error getting digital location: %w", err)
	}

	// Optionally fetch subscription info
	var subscription models.Subscription
	err = da.db.GetContext(
		ctx,
		&subscription,
		GetSingleDigitalLocationSubscriptionQuery,
		locationID,
	)
	switch {
	case err == nil:
		digitalLocation.Subscription = &subscription
	case err == sql.ErrNoRows:
		digitalLocation.Subscription = nil
	default:
		return models.DigitalLocation{}, fmt.Errorf("error getting subscription: %w", err)
	}

	da.logger.Debug("GetSingleDigitalLocation success", map[string]any{
		"location": digitalLocation,
	})

	return digitalLocation, nil
}

func (da *DigitalDbAdapter) GetAllDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocation, error) {
	// Define the DigitalLocationJoin struct
	type DigitalLocationJoin struct {
		models.DigitalLocation
		ItemsJSON     []byte    `db:"items"`
		SubID         *int64    `db:"sub_id"`
		BillingCycle  *string   `db:"billing_cycle"`
		CostPerCycle  *float64  `db:"cost_per_cycle"`
		AnchorDate    *time.Time `db:"anchor_date"`
		LastPaymentDate *time.Time `db:"last_payment_date"`
		NextPaymentDate *time.Time `db:"next_payment_date"`
		PaymentMethod *string   `db:"payment_method"`
		SubCreatedAt  *time.Time `db:"sub_created_at"`
		SubUpdatedAt  *time.Time `db:"sub_updated_at"`
	}

	// Execute the query to get raw data
	var locationsJoin []DigitalLocationJoin
	err := da.db.SelectContext(
		ctx,
		&locationsJoin,
		GetLocationsWithSubscriptionDataQuery,
		userID,
	)
	if err != nil {
			da.logger.Error("Failed to get user digital locations", map[string]any{"error": err})
			return nil, fmt.Errorf("error getting user digital locations: %w", err)
	}

	// Log the raw data from database
	if len(locationsJoin) > 0 {
			da.logger.Debug("Raw location data from DB", map[string]any{
					"count": len(locationsJoin),
					"first_location_id": locationsJoin[0].ID,
					"first_location_has_sub_id": locationsJoin[0].SubID != nil,
			})
	} else {
			da.logger.Debug("No locations found for user", map[string]any{"userID": userID})
			return []models.DigitalLocation{}, nil
	}

	// Create result array with exact size needed
	locations := make([]models.DigitalLocation, len(locationsJoin))

	// Process each location safely - using index-based for loop to avoid range variable copy issues
	for i := 0; i < len(locationsJoin); i++ {
			// Access the source data directly by index (no temporary variables)
			source := &locationsJoin[i]

			// Step 1: Build the base location from embedded struct
			baseLocation := source.DigitalLocation

			// Step 2: Unmarshal items
			var items []models.Game
			if err := json.Unmarshal(source.ItemsJSON, &items); err != nil {
					da.logger.Error("Failed to unmarshal items", map[string]any{
							"location_id": source.ID,
							"error": err,
					})
					return nil, fmt.Errorf("error unmarshaling items for location %s: %w", source.ID, err)
			}
			baseLocation.Items = items

			// Step 3: Add subscription if it exists
			if source.SubID != nil {
					da.logger.Debug("Adding subscription data to location", map[string]any{
							"locationID": source.ID,
							"subID": *source.SubID,
					})

					baseLocation.Subscription = &models.Subscription{
							ID:                *source.SubID,
							LocationID:        source.ID,
							BillingCycle:      *source.BillingCycle,
							CostPerCycle:      *source.CostPerCycle,
							AnchorDate:        *source.AnchorDate,
							NextPaymentDate:   *source.NextPaymentDate,
							LastPaymentDate:   source.LastPaymentDate,
							PaymentMethod:     *source.PaymentMethod,
							CreatedAt:         *source.SubCreatedAt,
							UpdatedAt:         *source.SubUpdatedAt,
					}
			}

			// Step 4: Add the fully constructed location to results
			locations[i] = baseLocation
	}

	da.logger.Debug("GetAllDigitalLocations success", map[string]any{
			"count": len(locations),
	})

	return locations, nil
}


// POST

func (a *DigitalDbAdapter) CreateDigitalLocation(
	ctx context.Context,
	userID string,
	location models.DigitalLocation,
) (models.DigitalLocation, error) {
	a.logger.Debug("Adding digital location", map[string]any{
		"userID": userID,
		"location": location,
		"is_active": location.IsActive,
	})

	// Validate userID is not empty
	if userID == "" {
		return models.DigitalLocation{}, fmt.Errorf("user ID cannot be empty")
	}

	// Check if a location with this name already exists for the user
	var existingID string
	err := a.db.QueryRowContext(
		ctx,
		CheckIfLocationExistsForUserQuery,
		userID,
		location.Name,
	).Scan(&existingID)

	if err == nil {
		return models.DigitalLocation{}, fmt.Errorf("a digital location with the name '%s' already exists", location.Name)
	} else if err != sql.ErrNoRows {
		return models.DigitalLocation{}, fmt.Errorf("error checking for existing location: %w", err)
	}

	// Set the user ID and timestamps
	location.UserID = userID
	now := time.Now()
	location.CreatedAt = now
	location.UpdatedAt = now

	// Use a transaction to ensure both location and subscription are saved
	err = postgres.WithTransaction(ctx, a.db, a.logger, func(tx *sqlx.Tx) error {
		// Insert the digital location
		locationQuery := `
			INSERT INTO digital_locations (id, user_id, name, is_subscription, is_active, url, payment_method, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, user_id, name, is_subscription, is_active, url, payment_method, created_at, updated_at
		`

		err = tx.QueryRowxContext(
			ctx,
			locationQuery,
			location.ID,
			userID,
			location.Name,
			location.IsSubscription,
			location.IsActive,
			location.URL,
			location.PaymentMethod,
			location.CreatedAt,
			location.UpdatedAt,
		).StructScan(&location)

		if err != nil {
			return fmt.Errorf("error adding digital location: %w", err)
		}

		// If subscription data exists, save it
		if location.Subscription != nil {
			a.logger.Debug("Saving subscription data", map[string]any{
				"locationID": location.ID,
				"subscription": location.Subscription,
			})

			var subID int64
			var digitalLocationID string
			var billingCycle string
			var costPerCycle float64
			var anchorDate time.Time
			var lastPaymentDate *time.Time
			var nextPaymentDate *time.Time
			var paymentMethod string
			var createdAt time.Time
			var updatedAt time.Time

			err = tx.QueryRowxContext(
				ctx,
				CreateSubscriptionWithAnchorDateQuery,
				location.ID,
				location.Subscription.BillingCycle,
				location.Subscription.CostPerCycle,
				location.Subscription.AnchorDate,
				location.Subscription.PaymentMethod,
				now,
				now,
			).Scan(&subID, &digitalLocationID, &billingCycle, &costPerCycle, &anchorDate, &lastPaymentDate, &nextPaymentDate, &paymentMethod, &createdAt, &updatedAt)

			if err != nil {
				return fmt.Errorf("error adding subscription: %w", err)
			}

			location.Subscription.ID = subID
			location.Subscription.CreatedAt = now
			location.Subscription.UpdatedAt = now

			a.logger.Debug("Subscription saved successfully", map[string]any{
				"subID": subID,
			})
		} else {
			a.logger.Debug("No subscription data to save", map[string]any{
				"locationID": location.ID,
			})
		}

		return nil
	})

	if err != nil {
		return models.DigitalLocation{}, err
	}

	a.logger.Debug("Digital location created successfully", map[string]any{
		"location": location,
		"is_active": location.IsActive,
	})

	return location, nil
}

// PUT
func (a *DigitalDbAdapter) UpdateDigitalLocation(
	ctx context.Context,
	userID string,
	location models.DigitalLocation,
) error {
	a.logger.Debug("Updating digital location", map[string]any{
		"userID": userID,
		"location": location,
		"is_active": location.IsActive, // Explicitly log the is_active value
	})

	// Validate ID is present
	if location.ID == "" {
		return fmt.Errorf("id is required")
	}

	// Check if the location belongs to the user
	if location.UserID != userID {
		return ErrUnauthorizedLocation
	}

	now := time.Now()
	query := `
			UPDATE digital_locations
			SET name = $1, is_active = $2, url = $3, updated_at = $4, is_subscription = $5, payment_method = $6
			WHERE id = $7 AND user_id = $8
	`

	a.logger.Debug("Executing SQL update", map[string]any{
		"query": query,
		"values": []interface{}{
			location.Name,
			location.IsActive,
			location.URL,
			now,
			location.IsSubscription,
			location.PaymentMethod,
			location.ID,
			userID,
		},
	})

	result, err := a.db.ExecContext(
			ctx,
			query,
			location.Name,
			location.IsActive,
			location.URL,
			now,
			location.IsSubscription,
			location.PaymentMethod,
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
		return fmt.Errorf("digital location not found")
	}

	a.logger.Debug("Update completed successfully", map[string]any{
		"rowsAffected": rowsAffected,
	})

	return nil
}

// DELETE
func (da *DigitalDbAdapter) DeleteDigitalLocation(
	ctx context.Context,
	userID string,
	locationIDs []string,
) (int64, error) {
	isBulk := len(locationIDs) > 1
	da.logger.Debug("DeleteDigitalLocation called", map[string]any{
		"userID":         userID,
		"locationIDs":    locationIDs,
		"isBulkOperation": isBulk,
	})

	// Validate input parameters
	if userID == "" {
		return 0, fmt.Errorf("user ID cannot be empty")
	}

	if len(locationIDs) == 0 {
		return 0, fmt.Errorf("no location IDs provided for deletion")
	}

	var totalDeleted int64

	// Use the transaction utility
	err := postgres.WithTransaction(ctx, da.db, da.logger, func(tx *sqlx.Tx) error {
		// First verify all locations exist and belong to the user
		var count int
		err := tx.QueryRowxContext(
			ctx,
			CheckIfAllLocationsExistForUserQuery,
			pq.Array(locationIDs),
			userID,
		).Scan(&count)
		if err != nil {
			return fmt.Errorf("error verifying locations: %w", err)
		}
		if count != len(locationIDs) {
			return fmt.Errorf("one or more locations not found or do not belong to user")
		}

		// Delete all related records and locations in one go
		result, err := tx.ExecContext(
			ctx,
			CascadingDeleteDigitalLocationQuery,
			pq.Array(locationIDs),
			userID,
		)
		if err != nil {
			return fmt.Errorf("error executing delete: %w", err)
		}

		// Get the number of deleted rows
		totalDeleted, err = result.RowsAffected()
		if err != nil {
			return fmt.Errorf("error getting rows affected: %w", err)
		}

		// If not all locations were deleted, return an error
		if totalDeleted < int64(len(locationIDs)) {
			return fmt.Errorf("partial deletion: %d of %d locations deleted", totalDeleted, len(locationIDs))
		}

		// Clean up games that are not associated with any digital location
		_, err = tx.ExecContext(
			ctx,
			DeleteOrphanedUserGamesQuery,
			userID,
		)
		if err != nil {
			return fmt.Errorf("error cleaning up orphaned user games: %w", err)
		}

		da.logger.Debug("DeleteDigitalLocation success", map[string]any{
			"totalDeleted":    totalDeleted,
			"isBulkOperation": isBulk,
		})

		return nil
	})

	if err != nil {
		da.logger.Error("DeleteDigitalLocation failed", map[string]any{
			"error":           err,
			"totalDeleted":    totalDeleted,
			"isBulkOperation": isBulk,
		})
		return totalDeleted, err
	}

	return totalDeleted, nil
}

