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

	const getDigitalLocationQuery = `
		WITH location_games AS (
			SELECT
				dl.id as location_id,
				json_agg(
					json_build_object(
						'id', g.id,
						'name', g.name,
						'summary', g.summary,
						'cover_id', g.cover_id,
						'cover_url', g.cover_url,
						'first_release_date', g.first_release_date,
						'rating', g.rating,
						'platform_names', COALESCE(
							(
								SELECT array_agg(p.name ORDER BY p.name)
								FROM game_platforms gp
								JOIN platforms p ON p.id = gp.platform_id
								WHERE gp.game_id = g.id
							),
							'{}'
						),
						'genre_names', COALESCE(
							(
								SELECT array_agg(gn.name ORDER BY gn.name)
								FROM game_genres gg
								JOIN genres gn ON gn.id = gg.genre_id
								WHERE gg.game_id = g.id
							),
							'{}'
						),
						'theme_names', COALESCE(
							(
								SELECT array_agg(t.name ORDER BY t.name)
								FROM game_themes gt
								JOIN themes t ON t.id = gt.theme_id
								WHERE gt.game_id = g.id
							),
							'{}'
						)
					)
				) as items
			FROM digital_locations dl
			LEFT JOIN digital_game_locations dgl ON dgl.digital_location_id = dl.id
			LEFT JOIN user_games ug ON ug.id = dgl.user_game_id
			LEFT JOIN games g ON g.id = ug.game_id
			WHERE dl.id = $1 AND dl.user_id = $2
			GROUP BY dl.id
		)
		SELECT
			dl.*,
			COALESCE(lg.items, '[]'::json) as items,
			dls.id as sub_id,
			dls.billing_cycle,
			dls.cost_per_cycle,
			dls.next_payment_date,
			dls.payment_method,
			dls.created_at as sub_created_at,
			dls.updated_at as sub_updated_at
		FROM digital_locations dl
		LEFT JOIN location_games lg ON lg.location_id = dl.id
		LEFT JOIN digital_location_subscriptions dls ON dls.digital_location_id = dl.id
		WHERE dl.id = $1 AND dl.user_id = $2
	`

	type DigitalLocationJoin struct {
		models.DigitalLocation
		ItemsJSON     []byte    `db:"items"`
		SubID         *int64    `db:"sub_id"`
		BillingCycle  *string   `db:"billing_cycle"`
		CostPerCycle  *float64  `db:"cost_per_cycle"`
		NextPaymentDate *time.Time `db:"next_payment_date"`
		PaymentMethod *string   `db:"payment_method"`
		SubCreatedAt  *time.Time `db:"sub_created_at"`
		SubUpdatedAt  *time.Time `db:"sub_updated_at"`
	}

	var locationJoin DigitalLocationJoin
	err := da.db.GetContext(ctx, &locationJoin, getDigitalLocationQuery, locationID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.DigitalLocation{}, fmt.Errorf("digital location not found: %w", err)
		}
		return models.DigitalLocation{}, fmt.Errorf("error getting digital location: %w", err)
	}

	// Unmarshal the JSON array into the Items field
	if err := json.Unmarshal(locationJoin.ItemsJSON, &locationJoin.Items); err != nil {
		return models.DigitalLocation{}, fmt.Errorf("error unmarshaling items: %w", err)
	}

	// If subscription data exists, add it to the location
	if locationJoin.SubID != nil {
		locationJoin.Subscription = &models.Subscription{
			ID:              *locationJoin.SubID,
			LocationID:      locationJoin.ID,
			BillingCycle:    *locationJoin.BillingCycle,
			CostPerCycle:    *locationJoin.CostPerCycle,
			NextPaymentDate: *locationJoin.NextPaymentDate,
			PaymentMethod:   *locationJoin.PaymentMethod,
			CreatedAt:       *locationJoin.SubCreatedAt,
			UpdatedAt:       *locationJoin.SubUpdatedAt,
		}
	}

	da.logger.Debug("GetDigitalLocation success", map[string]any{
		"location": locationJoin.DigitalLocation,
	})

	return locationJoin.DigitalLocation, nil
}

func (da *DigitalDbAdapter) GetUserDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
	// Define the DigitalLocationJoin struct
	type DigitalLocationJoin struct {
			models.DigitalLocation
			ItemsJSON     []byte    `db:"items"`
			SubID         *int64    `db:"sub_id"`
			BillingCycle  *string   `db:"billing_cycle"`
			CostPerCycle  *float64  `db:"cost_per_cycle"`
			NextPaymentDate *time.Time `db:"next_payment_date"`
			PaymentMethod *string   `db:"payment_method"`
			SubCreatedAt  *time.Time `db:"sub_created_at"`
			SubUpdatedAt  *time.Time `db:"sub_updated_at"`
	}

	// Simple query to get locations with their subscription data
	// FIXED: Removed the location_games JOIN which doesn't exist
	query := `
			SELECT
					dl.*,
					'[]'::json as items,
					dls.id as sub_id,
					dls.billing_cycle,
					dls.cost_per_cycle,
					dls.next_payment_date,
					dls.payment_method,
					dls.created_at as sub_created_at,
					dls.updated_at as sub_updated_at
			FROM digital_locations dl
			LEFT JOIN digital_location_subscriptions dls ON dls.digital_location_id = dl.id
			WHERE dl.user_id = $1
			ORDER BY dl.created_at
	`

	var locationsJoin []DigitalLocationJoin
	err := da.db.SelectContext(ctx, &locationsJoin, query, userID)
	if err != nil {
			da.logger.Error("Failed to get user digital locations", map[string]any{"error": err})
			return nil, fmt.Errorf("error getting user digital locations: %w", err)
	}

	// Log the raw data from database
	var firstLocationID string
	var hasSubID bool
	if len(locationsJoin) > 0 {
			firstLocationID = locationsJoin[0].ID
			hasSubID = locationsJoin[0].SubID != nil
	}
	da.logger.Debug("Raw location data from DB", map[string]any{
			"count": len(locationsJoin),
			"first_location_id": firstLocationID,
			"first_location_has_sub_id": hasSubID,
	})

	// Convert to DigitalLocation array and unmarshal items
	locations := make([]models.DigitalLocation, len(locationsJoin))
	for i, loc := range locationsJoin {
			if err := json.Unmarshal(loc.ItemsJSON, &loc.Items); err != nil {
					da.logger.Error("Failed to unmarshal items", map[string]any{
							"location_id": loc.ID,
							"error": err,
					})
					return nil, fmt.Errorf("error unmarshaling items for location %s: %w", loc.ID, err)
			}

			// If subscription data exists, add it to the location
			if loc.SubID != nil {
					da.logger.Debug("Found subscription data for location", map[string]any{
							"locationID": loc.ID,
							"subID": *loc.SubID,
							"billingCycle": *loc.BillingCycle,
							"costPerCycle": *loc.CostPerCycle,
					})

					loc.Subscription = &models.Subscription{
							ID:              *loc.SubID,
							LocationID:      loc.ID,
							BillingCycle:    *loc.BillingCycle,
							CostPerCycle:    *loc.CostPerCycle,
							NextPaymentDate: *loc.NextPaymentDate,
							PaymentMethod:   *loc.PaymentMethod,
							CreatedAt:       *loc.SubCreatedAt,
							UpdatedAt:       *loc.SubUpdatedAt,
					}
			}

			locations[i] = loc.DigitalLocation
	}

	da.logger.Debug("GetUserDigitalLocations success", map[string]any{
			"locations": locations,
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

func (a *DigitalDbAdapter) AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
	a.logger.Debug("Adding digital location", map[string]any{
		"userID": userID,
		"location": location,
		"is_active": location.IsActive, // Explicitly log the is_active value
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
		// Location with this name already exists
		return models.DigitalLocation{}, fmt.Errorf("a digital location with the name '%s' already exists", location.Name)
	} else if err != sql.ErrNoRows {
		// Some other database error occurred
		return models.DigitalLocation{}, fmt.Errorf("error checking for existing location: %w", err)
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

	// Set default service_type if not provided
	if location.ServiceType == "" {
		location.ServiceType = "basic"
	}

	a.logger.Debug("Executing SQL insert", map[string]any{
		"id": location.ID,
		"userID": userID,
		"name": location.Name,
		"service_type": location.ServiceType,
		"is_active": location.IsActive,
		"url": location.URL,
	})

	query := `
		INSERT INTO digital_locations (id, user_id, name, service_type, is_active, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, name, service_type, is_active, url, created_at, updated_at
	`

	err = a.db.QueryRowxContext(
		ctx,
		query,
		location.ID,
		userID,
		location.Name,
		location.ServiceType,
		location.IsActive,
		location.URL,
		location.CreatedAt,
		location.UpdatedAt,
	).StructScan(&location)

	if err != nil {
		// Check if the error is due to a unique constraint violation
		if strings.Contains(err.Error(), "digital_locations_user_id_name_key") {
			return models.DigitalLocation{}, fmt.Errorf("a digital location with the name '%s' already exists", location.Name)
		}
		return models.DigitalLocation{}, fmt.Errorf("error adding digital location: %w", err)
	}

	// Log the returned location to verify is_active value
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
			SET name = $1, is_active = $2, url = $3, updated_at = $4, service_type = $5
			WHERE id = $6 AND user_id = $7
	`

	a.logger.Debug("Executing SQL update", map[string]any{
		"query": query,
		"values": []interface{}{
			location.Name,
			location.IsActive,
			location.URL,
			now,
			location.ServiceType,
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
			location.ServiceType,
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
			// First try to find the location by ID or name
			var id string
			checkQuery := `
				SELECT id FROM digital_locations
				WHERE (
					(id::text = $1) OR
					(name = $1) OR
					(LOWER(name) = LOWER($1))
				) AND user_id = $2
			`

			err := tx.QueryRowxContext(ctx, checkQuery, locationID, userID).Scan(&id)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("digital location not found or does not belong to user")
				}
				return fmt.Errorf("error checking digital location: %w", err)
			}

			// Delete the location using the found ID
			deleteQuery := `
				DELETE FROM digital_locations
				WHERE id = $1 AND user_id = $2
			`

			result, err := tx.ExecContext(ctx, deleteQuery, id, userID)
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

	query := `
		SELECT id, digital_location_id, billing_cycle, cost_per_cycle,
		       next_payment_date, payment_method, created_at, updated_at
		FROM digital_location_subscriptions
		WHERE digital_location_id = $1
	`

	var subscription models.Subscription
	err := da.db.GetContext(ctx, &subscription, query, locationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No subscription found is not an error
		}
		return nil, fmt.Errorf("error getting subscription: %w", err)
	}

	return &subscription, nil
}

// AddSubscription creates a new subscription for a digital location
func (da *DigitalDbAdapter) AddSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
	da.logger.Debug("AddSubscription called", map[string]any{
		"subscription": subscription,
	})

	query := `
		INSERT INTO digital_location_subscriptions
			(digital_location_id, billing_cycle, cost_per_cycle,
			 next_payment_date, payment_method, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, digital_location_id, billing_cycle, cost_per_cycle,
		          next_payment_date, payment_method, created_at, updated_at
	`

	now := time.Now()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	err := da.db.QueryRowxContext(
		ctx,
		query,
		subscription.LocationID,
		subscription.BillingCycle,
		subscription.CostPerCycle,
		subscription.NextPaymentDate,
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

	query := `
		UPDATE digital_location_subscriptions
		SET billing_cycle = $1,
			cost_per_cycle = $2,
			next_payment_date = $3,
			payment_method = $4,
			updated_at = $5
		WHERE digital_location_id = $6
	`

	subscription.UpdatedAt = time.Now()
	result, err := da.db.ExecContext(
		ctx,
		query,
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

// RemoveSubscription deletes a subscription for a digital location
func (da *DigitalDbAdapter) RemoveSubscription(ctx context.Context, locationID string) error {
	da.logger.Debug("RemoveSubscription called", map[string]any{
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

// GetPayments retrieves all payments for a digital location
func (da *DigitalDbAdapter) GetPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	da.logger.Debug("GetPayments called", map[string]any{
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

// AddPayment records a new payment for a digital location
func (da *DigitalDbAdapter) AddPayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	da.logger.Debug("AddPayment called", map[string]any{
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

// GetPayment retrieves a specific payment by ID
func (da *DigitalDbAdapter) GetPayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
	da.logger.Debug("GetPayment called", map[string]any{
		"paymentID": paymentID,
	})

	query := `
		SELECT id, digital_location_id, amount, payment_date,
		       payment_method, transaction_id, created_at
		FROM digital_location_payments
		WHERE id = $1
	`

	var payment models.Payment
	err := da.db.GetContext(ctx, &payment, query, paymentID)
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

	query := `
		SELECT g.*
		FROM games g
		JOIN user_games ug ON ug.game_id = g.id
		JOIN digital_game_locations dgl ON dgl.user_game_id = ug.id
		WHERE dgl.digital_location_id = $1 AND ug.user_id = $2
	`

	var games []models.Game
	err := da.db.SelectContext(ctx, &games, query, locationID, userID)
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