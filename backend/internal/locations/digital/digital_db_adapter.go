package digital

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

// --- QUERIES ---
const (
	getLocationsWithSubscriptionDataQuery = `
		SELECT
			dl.*,
			'[]'::json as items,
			dls.id as sub_id,
			dls.billing_cycle,
			dls.cost_per_cycle,
			dls.anchor_date,
			dls.last_payment_date,
			dls.next_payment_date,
			dls.payment_method,
			dls.created_at as sub_created_at,
			dls.updated_at as sub_updated_at
		FROM digital_locations dl
		LEFT JOIN digital_location_subscriptions dls ON dls.digital_location_id = dl.id
		WHERE dl.user_id = $1
		ORDER BY dl.created_at
	`

	getSubscriptionByLocationIDQuery =  `
		SELECT id, digital_location_id, billing_cycle, cost_per_cycle,
		  anchor_date, last_payment_date, next_payment_date, payment_method, created_at, updated_at
		FROM digital_location_subscriptions
		WHERE digital_location_id = $1
	`

	getAllGamesInDigitalLocationQuery = `
		SELECT g.*
		FROM games g
		JOIN user_games ug ON ug.game_id = g.id
		JOIN digital_game_locations dgl ON dgl.user_game_id = ug.id
		WHERE dgl.digital_location_id = $1 AND ug.user_id = $2
	`

	//retrieves a specific payment by ID
	getSinglePaymentQuery = `
		SELECT id, digital_location_id, amount, payment_date,
			payment_method, transaction_id, created_at
		FROM digital_location_payments
		WHERE id = $1
	`

	recordPaymentQuery = `
		INSERT INTO digital_location_payments
			(digital_location_id, amount, payment_date, payment_method, transaction_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	updateSubscriptionQuery = `
		UPDATE digital_location_subscriptions
		SET billing_cycle = $1,
			cost_per_cycle = $2,
			anchor_date = $3,
			payment_method = $4,
			updated_at = $5
		WHERE digital_location_id = $6
	`

	updateSubscriptionLastPaymentDateQuery = `
		UPDATE digital_location_subscriptions
      SET last_payment_date = $1, updated_at = $2
    WHERE digital_location_id = $3
	`

	subscriptionAnchorDateQuery = `
		INSERT INTO digital_location_subscriptions
			(digital_location_id, billing_cycle, cost_per_cycle,
			 anchor_date, payment_method, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, digital_location_id, billing_cycle, cost_per_cycle,
		  anchor_date, last_payment_date, next_payment_date, payment_method, created_at, updated_at
	`
)


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
func (da *DigitalDbAdapter) GetSingleDigitalLocation(ctx context.Context, userID string, locationID string) (models.DigitalLocation, error) {
	da.logger.Debug("GetSingleDigitalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
	})

	// Fetch digital location
	var digitalLocation models.DigitalLocation
	err := da.db.GetContext(
		ctx,
		&digitalLocation,
		`SELECT * FROM digital_locations WHERE id = $1 AND user_id = $2`,
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
		`SELECT * FROM digital_location_subscriptions WHERE digital_location_id = $1`,
		locationID,
	)
	if err == nil {
		digitalLocation.Subscription = &subscription
	} else if err == sql.ErrNoRows {
		digitalLocation.Subscription = nil
	} else {
		return models.DigitalLocation{}, fmt.Errorf("error getting subscription: %w", err)
	}

	da.logger.Debug("GetSingleDigitalLocation success", map[string]any{
		"location": digitalLocation,
	})

	return digitalLocation, nil
}

func (da *DigitalDbAdapter) GetAllDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
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
		getLocationsWithSubscriptionDataQuery,
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
func (a *DigitalDbAdapter) ensureUserExists(ctx context.Context, userID string) error {
	a.logger.Debug("Ensuring user exists", map[string]any{"userID": userID})

	// Check if user exists
	var exists bool
	err := a.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)
	`, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}

	if !exists {
		// Create user with all required fields
		_, err = a.db.ExecContext(ctx, `
			INSERT INTO users (id, email, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW())
		`, userID, fmt.Sprintf("%s@example.com", userID))
		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
	}

	return nil
}

func (a *DigitalDbAdapter) CreateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
	a.logger.Debug("Adding digital location", map[string]any{
		"userID": userID,
		"location": location,
		"is_active": location.IsActive,
	})

	// Ensure user exists first
	if err := a.ensureUserExists(ctx, userID); err != nil {
		return models.DigitalLocation{}, fmt.Errorf("error ensuring user exists: %w", err)
	}

	// Check if a location with this name already exists for the user
	var existingID string
	err := a.db.QueryRowContext(ctx, `
		SELECT id FROM digital_locations
		WHERE user_id = $1 AND LOWER(name) = LOWER($2)
	`, userID, location.Name).Scan(&existingID)

	if err == nil {
		return models.DigitalLocation{}, fmt.Errorf("a digital location with the name '%s' already exists", location.Name)
	} else if err != sql.ErrNoRows {
		return models.DigitalLocation{}, fmt.Errorf("error checking for existing location: %w", err)
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
			subscriptionQuery := `
				INSERT INTO digital_location_subscriptions
					(digital_location_id, billing_cycle, cost_per_cycle, next_payment_date, payment_method, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id
			`

			var subID int64
			err = tx.QueryRowxContext(
				ctx,
				subscriptionQuery,
				location.ID,
				location.Subscription.BillingCycle,
				location.Subscription.CostPerCycle,
				location.Subscription.NextPaymentDate,
				location.Subscription.PaymentMethod,
				now,
				now,
			).Scan(&subID)

			if err != nil {
				return fmt.Errorf("error adding subscription: %w", err)
			}

			location.Subscription.ID = subID
			location.Subscription.CreatedAt = now
			location.Subscription.UpdatedAt = now
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
func (a *DigitalDbAdapter) UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error {
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
func (da *DigitalDbAdapter) DeleteDigitalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error) {
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
		checkQuery := `
			SELECT COUNT(*) FROM digital_locations
			WHERE id = ANY($1) AND user_id = $2
		`
		var count int
		err := tx.QueryRowxContext(ctx, checkQuery, pq.Array(locationIDs), userID).Scan(&count)
		if err != nil {
			return fmt.Errorf("error verifying locations: %w", err)
		}
		if count != len(locationIDs) {
			return fmt.Errorf("one or more locations not found or do not belong to user")
		}

		// Delete all related records and locations in one go
		deleteQuery := `
			WITH deleted_related AS (
				DELETE FROM digital_location_subscriptions
				WHERE digital_location_id = ANY($1)
			),
			deleted_games AS (
				DELETE FROM digital_game_locations
				WHERE digital_location_id = ANY($1)
			),
			deleted_payments AS (
				DELETE FROM digital_location_payments
				WHERE digital_location_id = ANY($1)
			)
			DELETE FROM digital_locations
			WHERE id = ANY($1) AND user_id = $2
		`

		// Execute the delete query
		result, err := tx.ExecContext(ctx, deleteQuery, pq.Array(locationIDs), userID)
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
		cleanupOrphanedUserGamesQuery := `
			DELETE FROM user_games
			WHERE user_id = $1
				AND id IN (
					SELECT ug.id
					FROM user_games ug
					LEFT JOIN digital_game_locations dgl ON ug.id = dgl.user_game_id
					WHERE ug.user_id = $1 AND dgl.user_game_id IS NULL
				)
		`
		_, err = tx.ExecContext(ctx, cleanupOrphanedUserGamesQuery, userID)
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

// FindDigitalLocationByName finds a digital location by name and user ID
func (da *DigitalDbAdapter) FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
	da.logger.Debug("FindDigitalLocationByName called", map[string]any{
		"userID": userID,
		"name":   name,
	})

	query := `
		SELECT id, user_id, name, service_type, is_active, url, created_at, updated_at
		FROM digital_locations
		WHERE user_id = $1 AND LOWER(name) = LOWER($2)
		`

	var location models.DigitalLocation
	err := da.db.GetContext(ctx, &location, query, userID, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.DigitalLocation{}, fmt.Errorf("digital location not found: %w", err)
		}
		return models.DigitalLocation{}, fmt.Errorf("error finding digital location: %w", err)
	}

	da.logger.Debug("FindDigitalLocationByName success", map[string]any{
		"location": location,
	})

	return location, nil
}

// GetSubscription retrieves a subscription for a digital location
func (da *DigitalDbAdapter) GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error) {
	da.logger.Debug("GetSubscription called", map[string]any{
		"locationID": locationID,
	})

	var subscription models.Subscription
	err := da.db.GetContext(
		ctx,
		&subscription,
		getSubscriptionByLocationIDQuery,
		locationID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No subscription found is not an error
		}
		return nil, fmt.Errorf("error getting subscription: %w", err)
	}

	return &subscription, nil
}

// CreateSubscription creates a new subscription for a digital location
func (da *DigitalDbAdapter) CreateSubscription(
	ctx context.Context,
	subscription models.Subscription,
) (*models.Subscription, error) {
	da.logger.Debug("CreateSubscription called", map[string]any{
		"subscription": subscription,
	})

	// Validate billing cycle format
	switch subscription.BillingCycle {
	case "1 month", "3 month", "6 month", "12 month":
		// Valid billing cycles
	default:
		return nil, fmt.Errorf("invalid billing cycle: %s. Must be one of: 1 month, 3 month, 6 month, 12 month", subscription.BillingCycle)
	}

	// Validate anchor date is provided
	if subscription.AnchorDate.IsZero() {
		return nil, fmt.Errorf("anchor_date is required for subscription creation")
	}

	now := time.Now()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	err := da.db.QueryRowxContext(
		ctx,
		subscriptionAnchorDateQuery,
		subscription.LocationID,
		subscription.BillingCycle,
		subscription.CostPerCycle,
		subscription.AnchorDate,
		subscription.PaymentMethod,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	).StructScan(&subscription)

	if err != nil {
		return nil, fmt.Errorf("error adding subscription: %w", err)
	}

	return &subscription, nil
}

// UpdateSubscription updates an existing subscription
func (da *DigitalDbAdapter) UpdateSubscription(ctx context.Context, subscription models.Subscription) error {
	da.logger.Debug("UpdateSubscription called", map[string]any{
		"subscription": subscription,
	})

	// Validate billing cycle format
	switch subscription.BillingCycle {
	case "1 month", "3 month", "6 month", "12 month":
		// Valid billing cycles
	default:
		return fmt.Errorf("invalid billing cycle: %s. Must be one of: 1 month, 3 month, 6 month, 12 month", subscription.BillingCycle)
	}

	// Validate anchor date is provided
	if subscription.AnchorDate.IsZero() {
		return fmt.Errorf("anchor_date is required for subscription updates")
	}

	subscription.UpdatedAt = time.Now()
	result, err := da.db.ExecContext(
		ctx,
		updateSubscriptionQuery,
		subscription.BillingCycle,
		subscription.CostPerCycle,
		subscription.NextPaymentDate,
		subscription.PaymentMethod,
		subscription.UpdatedAt,
		subscription.LocationID,
	)

	if err != nil {
		return fmt.Errorf("error updating subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// DeleteSubscription deletes a subscription for a digital location
func (da *DigitalDbAdapter) DeleteSubscription(ctx context.Context, locationID string) error {
	da.logger.Debug("DeleteSubscription called", map[string]any{
		"locationID": locationID,
	})

	query := `
		DELETE FROM digital_location_subscriptions
		WHERE digital_location_id = $1
	`

	result, err := da.db.ExecContext(ctx, query, locationID)
	if err != nil {
		return fmt.Errorf("error removing subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// GetAllPayments retrieves all payments for a digital location
func (da *DigitalDbAdapter) GetAllPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	da.logger.Debug("GetAllPayments called", map[string]any{
		"locationID": locationID,
	})

	query := `
		SELECT id, digital_location_id, amount, payment_date,
		       payment_method, transaction_id, created_at
		FROM digital_location_payments
		WHERE digital_location_id = $1
		ORDER BY payment_date DESC
	`

	var payments []models.Payment
	err := da.db.SelectContext(ctx, &payments, query, locationID)
	if err != nil {
		return nil, fmt.Errorf("error getting payments: %w", err)
	}

	return payments, nil
}

// CreatePayment records a new payment for a digital location
func (da *DigitalDbAdapter) CreatePayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	da.logger.Debug("CreatePayment called", map[string]any{
		"payment": payment,
	})

	query := `
		INSERT INTO digital_location_payments
			(digital_location_id, amount, payment_date,
			 payment_method, transaction_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, digital_location_id, amount, payment_date,
		          payment_method, transaction_id, created_at
	`

	payment.CreatedAt = time.Now()

	err := da.db.QueryRowxContext(
		ctx,
		query,
		payment.LocationID,
		payment.Amount,
		payment.PaymentDate,
		payment.PaymentMethod,
		payment.TransactionID,
		payment.CreatedAt,
	).StructScan(&payment)

	if err != nil {
		return nil, fmt.Errorf("error adding payment: %w", err)
	}

	return &payment, nil
}

// GetSinglePayment retrieves a specific payment by ID
func (da *DigitalDbAdapter) GetSinglePayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
	da.logger.Debug("GetSinglePayment called", map[string]any{
		"paymentID": paymentID,
	})

	var payment models.Payment
	err := da.db.GetContext(
		ctx,
		&payment,
		getSinglePaymentQuery,
		paymentID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("error getting payment: %w", err)
	}

	return &payment, nil
}

// AddGameToDigitalLocation adds a game to a digital location
func (da *DigitalDbAdapter) AddGameToDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error {
	da.logger.Debug("AddGameToDigitalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
		"gameID":     gameID,
	})

	// First, get the user_game_id for this user and game
	var userGameID int
	err := da.db.GetContext(ctx, &userGameID, `
		SELECT id FROM user_games
		WHERE user_id = $1 AND game_id = $2
	`, userID, gameID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("game not found in user's library")
		}
		return fmt.Errorf("error getting user game: %w", err)
	}

	// Then, add the game to the digital location
	_, err = da.db.ExecContext(ctx, `
		INSERT INTO digital_game_locations (user_game_id, digital_location_id)
		VALUES ($1, $2)
	`, userGameID, locationID)
	if err != nil {
		if strings.Contains(err.Error(), "digital_game_locations_user_game_id_digital_location_id_key") {
			return fmt.Errorf("game already exists in this digital location")
		}
		return fmt.Errorf("error adding game to digital location: %w", err)
	}

	return nil
}

// RemoveGameFromDigitalLocation removes a game from a digital location
func (da *DigitalDbAdapter) RemoveGameFromDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error {
	da.logger.Debug("RemoveGameFromDigitalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
		"gameID":     gameID,
	})

	// First, get the user_game_id for this user and game
	var userGameID int
	err := da.db.GetContext(ctx, &userGameID, `
		SELECT id FROM user_games
		WHERE user_id = $1 AND game_id = $2
	`, userID, gameID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("game not found in user's library")
		}
		return fmt.Errorf("error getting user game: %w", err)
	}

	// Then, remove the game from the digital location
	result, err := da.db.ExecContext(ctx, `
		DELETE FROM digital_game_locations
		WHERE user_game_id = $1 AND digital_location_id = $2
	`, userGameID, locationID)
	if err != nil {
		return fmt.Errorf("error removing game from digital location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("game not found in digital location")
	}

	return nil
}

// GetGamesByDigitalLocationID gets all games in a digital location
func (da *DigitalDbAdapter) GetGamesByDigitalLocationID(ctx context.Context, userID string, locationID string) ([]models.Game, error) {
	da.logger.Debug("GetGamesByDigitalLocationID called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
	})

	var games []models.Game
	err := da.db.SelectContext(
		ctx,
		&games,
		getAllGamesInDigitalLocationQuery,
		locationID,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting games for digital location: %w", err)
	}

	return games, nil
}

func (da *DigitalDbAdapter) ValidateSubscriptionExists(ctx context.Context, locationID string) (*models.Subscription, error) {
	da.logger.Debug("ValidateSubscriptionExists called", map[string]any{
			"locationID": locationID,
	})

	// Just check if subscription exists
	existingSub, err := da.GetSubscription(ctx, locationID)
	if err != nil {
			return nil, fmt.Errorf("failed to validate subscription: %w", err)
	}

	return existingSub, nil
}

func (da *DigitalDbAdapter) RecordPayment(
	ctx context.Context,
	payment models.Payment,
) error {
		return postgres.WithTransaction(
			ctx,
			da.db,
			da.logger,
			func(tx *sqlx.Tx) error {
				// 1. Record the payment
				_, err := tx.ExecContext(
					ctx,
					recordPaymentQuery,
					payment.LocationID,
					payment.Amount,
					payment.PaymentDate,
					payment.PaymentMethod,
					payment.TransactionID,
					time.Now(),
				)
				if err != nil {
					return fmt.Errorf("error recording payment: %w", err)
				}

				// 2. Update the subscription's last payment date
				result, err := tx.ExecContext(
					ctx,
					updateSubscriptionLastPaymentDateQuery,
					payment.PaymentDate,
					time.Now(),
					payment.LocationID,
				)
				if err != nil {
					return fmt.Errorf("error updating subscription last payment date: %w", err)
				}

				rowsAffected, err := result.RowsAffected()
				if err != nil {
					return fmt.Errorf("error getting rows affected: %w", err)
				}

				if rowsAffected == 0 {
					return fmt.Errorf("subscription not found for location: %s", payment.LocationID)
				}

				return nil
			})
}