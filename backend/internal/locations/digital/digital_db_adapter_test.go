package digital

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
)

/*
Behavior:
- Retrieving a single digital location for a user
- Retrieving all digital locations for a user
- Adding a new digital location
- Updating an existing digital location
- Removing a digital location
- Handling db errors
- Ensuring user may only access their own locations

Scenarios:
- GetDigitalLocation:
  - Successfully retrieves a valid location
  - Returns error when location not found
  - Handles db errors
- GetUserDigitalLocations:
  - Successfully retrieves all locations for a user
  - Returns empty slice when no locations exist
  - Handles db errors
- AddDigitalLocation:
  - Successfully adds new location
  - Handles db errors
- UpdateDigitalLocation:
  - Successfully updates a location
  - Returns error when location not found
  - Handles db errors
  - Returns error when user is not authorized
- RemoveDigitalLocation:
  - Successfully removes a location
  - Returns errors when location not found
  - Handles db errors
*/

func TestDigitalDbAdapter(t *testing.T) {
	// Set up base app context for testing
	testLogger := testutils.NewTestLogger()

	// Create mock DB + adapter
	setupMockDB := func() (*DigitalDbAdapter, sqlmock.Sqlmock, error) {
		// Create mock sqldatabase
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			return nil, nil, err
		}

		// Create a sqlx wrapper around mock data
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		// Create the adapter with the mock DB
		adapter := &DigitalDbAdapter{
			db:     sqlxDB,
			logger: testLogger,
		}

		return adapter, mock, nil
	}

	/*
		GIVEN a request to get a specific digital location
		WHEN the location exists in the database
		THEN the adapter returns the location
	*/
	t.Run(`GetDigitalLocation - Successfully retrieves a valid location`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		now := time.Now()

		// Add a valid JSON array for items
		itemsJSON := `[]` // Empty JSON array

		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at", "items"}).
			AddRow(locationID, userID, "Test Location", true, "https://example.com", now, now, itemsJSON)

		mock.ExpectQuery("SELECT (.+) FROM digital_locations").
			WithArgs(locationID, userID).
			WillReturnRows(rows)

		// Execute
		location, err := adapter.GetDigitalLocation(context.Background(), userID, locationID)

		// Manually set the features for verification
		expectedLocation := models.DigitalLocation{
			ID:        locationID,
			UserID:    userID,
			Name:      "Test Location",
			IsActive:  true,
			URL:       "https://example.com",
			CreatedAt: now,
			UpdatedAt: now,
			Items:     []models.Game{}, // Empty slice of games
		}

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if location.ID != expectedLocation.ID || location.Name != expectedLocation.Name {
			t.Errorf("Expected location to be %+v, but got %+v", expectedLocation, location)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get a specific digital location
		WHEN the location does not exist in the db
		THEN the adapter returns an error
	*/
	t.Run(`GetDigitalLocation - Returns error when location not found`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"

		// Setup mock expectations
		mock.ExpectQuery("SELECT (.+) FROM digital_locations").
			WithArgs(locationID, userID).
			WillReturnError(sql.ErrNoRows)

		// Execute
		_, err = adapter.GetDigitalLocation(context.Background(), userID, locationID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error but got nil")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Expected error to contain %v, but got %v", sql.ErrNoRows, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get A SPECIFIC digital location
		WHEN the db returns an error
		THEN the adapter also returns an error
	*/
	t.Run(`GetDigitalLocation - Handles database errors`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock db: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		dbError := errors.New("database error")

		// Set up mock expectations
		mock.ExpectQuery("SELECT (.+) FROM digital_locations").
			WithArgs(locationID, userID).
			WillReturnError(dbError)

		// Execute the function
		_, err = adapter.GetDigitalLocation(context.Background(), userID, locationID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error but got nil")
		}
		if !errors.Is(err, dbError) {
			t.Errorf("Expected error to contain %v, but instead got %v", dbError, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get ALL digital locations for a user
		WHEN locations exist in the db
		THEN the adapter returns all locations
	*/
	t.Run("GetUserDigitalLocations - Successfully retrieves all locations", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		now := time.Now()

		// Add valid JSON arrays for items
		itemsJSON := `[]` // Empty JSON array

		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at", "items"}).
			AddRow("loc1", userID, "Location 1", true, "https://example1.com", now, now, itemsJSON).
			AddRow("loc2", userID, "Location 2", true, "https://example2.com", now, now, itemsJSON)

		mock.ExpectQuery("SELECT (.+) FROM digital_locations").
			WithArgs(userID).
			WillReturnRows(rows)

		// Execute the function
		locations, err := adapter.GetUserDigitalLocations(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != 2 {
			t.Errorf("Expected 2 locations, got %d", len(locations))
		}
		if locations[0].Name != "Location 1" || locations[1].Name != "Location 2" {
			t.Errorf("Locations not returned correctly: %+v", locations)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get ALL digital locations for a user
		WHEN NO locations exist in the db
		THEN the adapter returns an empty slice
	*/
	t.Run(`GetUserDigitalLocations - Returns empty slice when no locations exist`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"

		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at", "items"})

		mock.ExpectQuery("SELECT (.+) FROM digital_locations").
			WithArgs(userID).
			WillReturnRows(rows)

		// Execute the function
		locations, err := adapter.GetUserDigitalLocations(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != 0 {
			t.Errorf("Expected 0 locations, got %d", len(locations))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get ALL digital locations for a user
		WHEN the database returns an error
		THEN the adapter returns the error
	*/
	t.Run(`GetUserDigitalLocations - Handles database errors`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		dbError := errors.New("database error")

		// Set up mock expectations
		mock.ExpectQuery("SELECT (.+) FROM digital_locations").
			WithArgs(userID).
			WillReturnError(dbError)

		// Execute
		_, err = adapter.GetUserDigitalLocations(context.Background(), userID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error but got nil")
		}
		if !errors.Is(err, dbError) {
			t.Errorf("Expected error to contain %v, but instead got %v", dbError, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to ADD a new digital location
		WHEN the db operation is successful
		THEN the adapter adds the location and returns it with generated fields
	*/
	t.Run(`AddDigitalLocation - Successfully adds a new location`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		now := time.Now()

		location := models.DigitalLocation{
			ID:       locationID,
			Name:     "New Location",
			IsActive: true,
			URL:      "https://example.com",
		}

		// Set up mock expectations
		// For ensureUserExists - user exists
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Check if a location with the same name already exists
		mock.ExpectQuery("SELECT id FROM digital_locations").
			WithArgs(userID, "New Location").
			WillReturnError(sql.ErrNoRows)

		// Add the location
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at"}).
			AddRow(locationID, userID, "New Location", true, "https://example.com", now, now)

		mock.ExpectQuery("INSERT INTO digital_locations").
			WithArgs(
				sqlmock.AnyArg(), // ID could be generated if empty
				userID,
				"New Location",
				true,
				"https://example.com",
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnRows(rows)

		// Execute
		result, err := adapter.AddDigitalLocation(context.Background(), userID, location)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result.ID != locationID || result.UserID != userID || result.Name != "New Location" {
			t.Errorf("Expected location with ID %s and name 'New Location', got %+v", locationID, result)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to ADD a new digital location
		WHEN the db returns an error
		THEN then the adapter also returns an error
	*/
	t.Run("AddDigitalLocation - Handles database errors", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		dbError := errors.New("database error")

		location := models.DigitalLocation{
			ID:       locationID,
			Name:     "New Location",
			IsActive: true,
			URL:      "https://example.com",
		}

		// Set up mock expectations
		// For ensureUserExists - user exists
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Check if a location with the same name already exists
		mock.ExpectQuery("SELECT id FROM digital_locations").
			WithArgs(userID, "New Location").
			WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery("INSERT INTO digital_locations").
			WithArgs(
				sqlmock.AnyArg(), // ID could be generated if empty
				userID,
				"New Location",
				true,
				"https://example.com",
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnError(dbError)

		// Execute
		_, err = adapter.AddDigitalLocation(context.Background(), userID, location)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, dbError) {
			t.Errorf("Expected error to contain %v, got %v", dbError, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to UPDATE a digital location
		WHEN the location exists AND belongs to the user
		THEN the adapter updates the location
	*/
	t.Run(`UpdateDigitalLocation - Successfully updates a location`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"

		location := models.DigitalLocation{
			ID:       locationID,
			UserID:   userID,
			Name:     "Updated Location",
			IsActive: true,
			URL:      "https://updated-example.com",
		}

		// Set up mock expectations
		mock.ExpectExec("UPDATE digital_locations").
			WithArgs(
				"Updated Location",
				true,
				"https://updated-example.com",
				sqlmock.AnyArg(), // updated_at
				locationID,
				userID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Execute the function
		err = adapter.UpdateDigitalLocation(context.Background(), userID, location)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to UPDATE a digital location
		WHEN the location doesn't exist or doesn't belong to the user
		THEN the adapter returns an error
	*/
	t.Run(`UpdateDigitalLocation - Returns error when location not found`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"

		location := models.DigitalLocation{
			ID:       locationID,
			UserID:   userID,
			Name:     "Updated Location",
			IsActive: true,
			URL:      "https://updated-example.com",
		}

		// Set up mock expectations
		mock.ExpectExec("UPDATE digital_locations").
			WithArgs(
				"Updated Location",
				true,
				"https://updated-example.com",
				sqlmock.AnyArg(), // updated_at
				locationID,
				userID,
			).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Execute the function
		err = adapter.UpdateDigitalLocation(context.Background(), userID, location)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to UPDATE a digital location
		WHEN the database returns an error
		THEN the adapter returns the error
	*/
	t.Run(`UpdateDigitalLocation - Handles database errors`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		dbError := errors.New("database error")

		location := models.DigitalLocation{
			ID:       locationID,
			UserID:   userID,
			Name:     "Updated Location",
			IsActive: true,
			URL:      "https://updated-example.com",
		}

		// Set up mock expectations
		mock.ExpectExec("UPDATE digital_locations").
			WithArgs(
				"Updated Location",
				true,
				"https://updated-example.com",
				sqlmock.AnyArg(), // updated_at
				locationID,
				userID,
			).
			WillReturnError(dbError)

		// Execute the function
		err = adapter.UpdateDigitalLocation(context.Background(), userID, location)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, dbError) {
			t.Errorf("Expected error to contain %v, got %v", dbError, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to UPDATE a digital location
		WHEN the userID doesn't match the location's userID
		THEN the adapter returns an UNAUTHORIZED error
	*/
	t.Run("UpdateDigitalLocation - Returns error when user is not authorized", func(t *testing.T) {
		// Setup
		adapter, _, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		wrongUserID := "wrong-user-id"
		locationID := "test-location-id"

		location := models.DigitalLocation{
			ID:       locationID,
			UserID:   wrongUserID,
			Name:     "Updated Location",
			IsActive: true,
			URL:      "https://updated-example.com",
		}

		// Execute the function
		err = adapter.UpdateDigitalLocation(context.Background(), userID, location)

		// Verify
		if err == nil {
			t.Errorf("Expected an unauthorized error, got nil")
		}
		if !errors.Is(err, ErrUnauthorizedLocation) {
			t.Errorf("Expected unauthorized error, got %v", err)
		}
	})

	/*
		GIVEN a request to REMOVE a digital location
		WHEN the location exists AND belongs to the user
		THEN the adapter removes the location
	*/
	t.Run(`RemoveDigitalLocation - Successfully removes a location`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"

		// Setup mock expectations for transaction
		mock.ExpectBegin()

		// Check if location exists
		rows := sqlmock.NewRows([]string{"id"}).AddRow(locationID)
		mock.ExpectQuery("SELECT id FROM digital_locations").
			WithArgs(locationID, userID).
			WillReturnRows(rows)

		// Delete location
		mock.ExpectExec("DELETE FROM digital_locations").
			WithArgs(locationID, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Commit transaction
		mock.ExpectCommit()

		// Execute
		err = adapter.RemoveDigitalLocation(context.Background(), userID, locationID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to REMOVE a digital location
		WHEN the location doesn't exist or doesn't belong to the user
		THEN the adapter returns an error
	*/
	t.Run(`RemoveDigitalLocation - Returns error when location not found`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"

		// Setup mock expectations
		mock.ExpectBegin()

		// Check if location exists - not found
		mock.ExpectQuery("SELECT id FROM digital_locations").
			WithArgs(locationID, userID).
			WillReturnError(sql.ErrNoRows)

		// Rollback transaction
		mock.ExpectRollback()

		// Execute the function
		err = adapter.RemoveDigitalLocation(context.Background(), userID, locationID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to REMOVE a digital location
		WHEN the database returns an error during deletion
		THEN the adapter returns the error
	*/
	t.Run(`RemoveDigitalLocation - Handles database errors during deletion`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		dbError := errors.New("database error")

		// Setup mock expectations
		mock.ExpectBegin()

		// Check if the location exists
		rows := sqlmock.NewRows([]string{"id"}).AddRow(locationID)
		mock.ExpectQuery("SELECT id FROM digital_locations").
			WithArgs(locationID, userID).
			WillReturnRows(rows)

		// Delete the location and get an error
		mock.ExpectExec("DELETE FROM digital_locations").
			WithArgs(locationID, userID).
			WillReturnError(dbError)

		// Rollback the transaction
		mock.ExpectRollback()

		// Execute the function
		err = adapter.RemoveDigitalLocation(context.Background(), userID, locationID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, dbError) {
			t.Errorf("Expected error to contain %v, got %v", dbError, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})
}

func TestGetSubscription(t *testing.T) {
	tests := []struct {
		name           string
		locationID     string
		mockSetup      func(sqlmock.Sqlmock)
		expectedResult *models.Subscription
		expectedError  error
	}{
		{
			name:       "Successfully get subscription",
			locationID: "test-location-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "location_id", "billing_cycle", "cost_per_cycle", "next_payment_date", "payment_method", "created_at", "updated_at"}).
					AddRow(1, "test-location-id", "monthly", 9.99, time.Now(), "Visa", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, location_id, billing_cycle, cost_per_cycle, next_payment_date, payment_method, created_at, updated_at FROM digital_location_subscriptions WHERE location_id = \\$1").
					WithArgs("test-location-id").
					WillReturnRows(rows)
			},
			expectedResult: &models.Subscription{
				ID:             1,
				LocationID:     "test-location-id",
				BillingCycle:   "monthly",
				CostPerCycle:   9.99,
				NextPaymentDate: time.Now(),
				PaymentMethod:  "Visa",
			},
			expectedError: nil,
		},
		{
			name:       "Subscription not found",
			locationID: "non-existent-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, location_id, billing_cycle, cost_per_cycle, next_payment_date, payment_method, created_at, updated_at FROM digital_location_subscriptions WHERE location_id = \\$1").
					WithArgs("non-existent-id").
					WillReturnError(sql.ErrNoRows)
			},
			expectedResult: nil,
			expectedError:  sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			// Create a mock app context
			appCtx := &appcontext.AppContext{
				Logger: testutils.NewTestLogger(),
			}

			// Create the adapter with the mock DB
			sqlxDB := sqlx.NewDb(db, "sqlmock")
			if sqlxDB == nil {
				t.Fatal("Failed to create sqlx.DB")
			}

			adapter := &DigitalDbAdapter{
				db:     sqlxDB,
				logger: appCtx.Logger,
			}

			result, err := adapter.GetSubscription(context.Background(), tt.locationID)

			if err != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expectedResult != nil {
				if result == nil {
					t.Error("expected non-nil result, got nil")
				} else {
					// Compare fields individually with nil checks
					if result.ID != tt.expectedResult.ID {
						t.Errorf("expected ID %v, got %v", tt.expectedResult.ID, result.ID)
					}
					if result.LocationID != tt.expectedResult.LocationID {
						t.Errorf("expected LocationID %v, got %v", tt.expectedResult.LocationID, result.LocationID)
					}
					if result.BillingCycle != tt.expectedResult.BillingCycle {
						t.Errorf("expected BillingCycle %v, got %v", tt.expectedResult.BillingCycle, result.BillingCycle)
					}
					if result.CostPerCycle != tt.expectedResult.CostPerCycle {
						t.Errorf("expected CostPerCycle %v, got %v", tt.expectedResult.CostPerCycle, result.CostPerCycle)
					}
					if result.PaymentMethod != tt.expectedResult.PaymentMethod {
						t.Errorf("expected PaymentMethod %v, got %v", tt.expectedResult.PaymentMethod, result.PaymentMethod)
					}
				}
			} else if result != nil {
				t.Error("expected nil result, got non-nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAddSubscription(t *testing.T) {
	tests := []struct {
		name           string
		subscription   models.Subscription
		mockSetup      func(sqlmock.Sqlmock)
		expectedResult *models.Subscription
		expectedError  error
	}{
		{
			name: "Successfully add subscription",
			subscription: models.Subscription{
				LocationID:     "test-location-id",
				BillingCycle:   "monthly",
				CostPerCycle:   9.99,
				NextPaymentDate: time.Now(),
				PaymentMethod:  "Visa",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO digital_location_subscriptions").
					WithArgs("test-location-id", "monthly", 9.99, sqlmock.AnyArg(), "Visa").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedResult: &models.Subscription{
				ID:             1,
				LocationID:     "test-location-id",
				BillingCycle:   "monthly",
				CostPerCycle:   9.99,
				NextPaymentDate: time.Now(),
				PaymentMethod:  "Visa",
			},
			expectedError: nil,
		},
		{
			name: "Failed to add subscription",
			subscription: models.Subscription{
				LocationID:     "test-location-id",
				BillingCycle:   "monthly",
				CostPerCycle:   9.99,
				NextPaymentDate: time.Now(),
				PaymentMethod:  "Visa",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO digital_location_subscriptions").
					WithArgs("test-location-id", "monthly", 9.99, sqlmock.AnyArg(), "Visa").
					WillReturnError(sql.ErrConnDone)
			},
			expectedResult: nil,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			// Create a mock app context
			appCtx := &appcontext.AppContext{
				Logger: testutils.NewTestLogger(),
			}

			// Create the adapter with the mock DB
			sqlxDB := sqlx.NewDb(db, "sqlmock")
			if sqlxDB == nil {
				t.Fatal("Failed to create sqlx.DB")
			}

			adapter := &DigitalDbAdapter{
				db:     sqlxDB,
				logger: appCtx.Logger,
			}

			result, err := adapter.AddSubscription(context.Background(), tt.subscription)

			if err != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expectedResult != nil {
				if result == nil {
					t.Error("expected non-nil result, got nil")
				} else {
					// Compare fields individually with nil checks
					if result.ID != tt.expectedResult.ID {
						t.Errorf("expected ID %v, got %v", tt.expectedResult.ID, result.ID)
					}
					if result.LocationID != tt.expectedResult.LocationID {
						t.Errorf("expected LocationID %v, got %v", tt.expectedResult.LocationID, result.LocationID)
					}
					if result.BillingCycle != tt.expectedResult.BillingCycle {
						t.Errorf("expected BillingCycle %v, got %v", tt.expectedResult.BillingCycle, result.BillingCycle)
					}
					if result.CostPerCycle != tt.expectedResult.CostPerCycle {
						t.Errorf("expected CostPerCycle %v, got %v", tt.expectedResult.CostPerCycle, result.CostPerCycle)
					}
					if result.PaymentMethod != tt.expectedResult.PaymentMethod {
						t.Errorf("expected PaymentMethod %v, got %v", tt.expectedResult.PaymentMethod, result.PaymentMethod)
					}
				}
			} else if result != nil {
				t.Error("expected nil result, got non-nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetPayments(t *testing.T) {
	tests := []struct {
		name           string
		locationID     string
		mockSetup      func(sqlmock.Sqlmock)
		expectedResult []models.Payment
		expectedError  error
	}{
		{
			name:       "Successfully get payments",
			locationID: "test-location-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "location_id", "amount", "payment_date", "payment_method", "transaction_id", "created_at"}).
					AddRow(1, "test-location-id", 9.99, time.Now(), "Visa", "tx123", time.Now()).
					AddRow(2, "test-location-id", 9.99, time.Now(), "Visa", "tx124", time.Now())
				mock.ExpectQuery("SELECT id, location_id, amount, payment_date, payment_method, transaction_id, created_at FROM digital_location_payments WHERE location_id = \\$1 ORDER BY payment_date DESC").
					WithArgs("test-location-id").
					WillReturnRows(rows)
			},
			expectedResult: []models.Payment{
				{
					ID:            1,
					LocationID:    "test-location-id",
					Amount:        9.99,
					PaymentDate:   time.Now(),
					PaymentMethod: "Visa",
					TransactionID: "tx123",
				},
				{
					ID:            2,
					LocationID:    "test-location-id",
					Amount:        9.99,
					PaymentDate:   time.Now(),
					PaymentMethod: "Visa",
					TransactionID: "tx124",
				},
			},
			expectedError: nil,
		},
		{
			name:       "No payments found",
			locationID: "non-existent-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, location_id, amount, payment_date, payment_method, transaction_id, created_at FROM digital_location_payments WHERE location_id = \\$1 ORDER BY payment_date DESC").
					WithArgs("non-existent-id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "location_id", "amount", "payment_date", "payment_method", "transaction_id", "created_at"}))
			},
			expectedResult: []models.Payment{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			// Create a mock app context
			appCtx := &appcontext.AppContext{
				Logger: testutils.NewTestLogger(),
			}

			// Create the adapter with the mock DB
			adapter := &DigitalDbAdapter{
				db:     sqlx.NewDb(db, "sqlmock"),
				logger: appCtx.Logger,
			}

			result, err := adapter.GetPayments(context.Background(), tt.locationID)

			if err != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if len(result) != len(tt.expectedResult) {
				t.Errorf("expected %d payments, got %d", len(tt.expectedResult), len(result))
			}

			for i, payment := range result {
				if payment.ID != tt.expectedResult[i].ID ||
					payment.LocationID != tt.expectedResult[i].LocationID ||
					payment.Amount != tt.expectedResult[i].Amount ||
					payment.PaymentMethod != tt.expectedResult[i].PaymentMethod ||
					payment.TransactionID != tt.expectedResult[i].TransactionID {
					t.Errorf("expected payment %+v, got %+v", tt.expectedResult[i], payment)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestNewDigitalDbAdapter(t *testing.T) {
	logger := testutils.NewTestLogger()

	appCtx := &appcontext.AppContext{
		Logger: logger,
	}

	adapter, err := NewDigitalDbAdapter(appCtx)
	if err != nil {
		t.Fatalf("Failed to create DigitalDbAdapter: %v", err)
	}

	// Verify the adapter was created with the correct components
	if adapter == nil {
		t.Error("Expected adapter to be non-nil")
	}
}

func TestGetUserDigitalLocations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up mock expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at"}).
		AddRow("test-location-id", "test-user", "Test Location", true, "https://test.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT (.+) FROM digital_locations").
		WithArgs("test-user").
		WillReturnRows(rows)

	// Create a mock app context
	appCtx := &appcontext.AppContext{
		Logger: testutils.NewTestLogger(),
	}

	// Create the adapter with the mock DB
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if sqlxDB == nil {
		t.Fatal("Failed to create sqlx.DB")
	}

	adapter := &DigitalDbAdapter{
		db:     sqlxDB,
		logger: appCtx.Logger,
	}

	// Test cases
	t.Run("success", func(t *testing.T) {
		locations, err := adapter.GetUserDigitalLocations(context.Background(), "test-user")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if locations == nil {
			t.Error("Expected non-nil locations")
		}
		if len(locations) == 0 {
			t.Error("Expected at least one location")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetDigitalLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up mock expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at"}).
		AddRow("test-location-id", "test-user", "Test Location", true, "https://test.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT (.+) FROM digital_locations").
		WithArgs("test-location-id", "test-user").
		WillReturnRows(rows)

	// Create a mock app context
	appCtx := &appcontext.AppContext{
		Logger: testutils.NewTestLogger(),
	}

	// Create the adapter with the mock DB
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if sqlxDB == nil {
		t.Fatal("Failed to create sqlx.DB")
	}

	adapter := &DigitalDbAdapter{
		db:     sqlxDB,
		logger: appCtx.Logger,
	}

	// Test cases
	t.Run("success", func(t *testing.T) {
		location, err := adapter.GetDigitalLocation(context.Background(), "test-user", "test-location-id")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if location.ID == "" {
			t.Error("Expected non-empty location ID")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddDigitalLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up mock expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at"}).
		AddRow("test-location-id", "test-user", "Test Location", true, "https://test.com", time.Now(), time.Now())
	mock.ExpectQuery("INSERT INTO digital_locations").
		WithArgs("test-location-id", "test-user", "Test Location", true, "https://test.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	// Create a mock app context
	appCtx := &appcontext.AppContext{
		Logger: testutils.NewTestLogger(),
	}

	// Create the adapter with the mock DB
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if sqlxDB == nil {
		t.Fatal("Failed to create sqlx.DB")
	}

	adapter := &DigitalDbAdapter{
		db:     sqlxDB,
		logger: appCtx.Logger,
	}

	// Test cases
	t.Run("success", func(t *testing.T) {
		testLocation := models.DigitalLocation{
			Name:     "Test Location",
			IsActive: true,
			URL:      "https://test.com",
		}

		location, err := adapter.AddDigitalLocation(context.Background(), "test-user", testLocation)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if location.ID == "" {
			t.Error("Expected non-empty location ID")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddGameToDigitalLocation(t *testing.T) {
	// Setup
	testLogger := testutils.NewTestLogger()

	// Create mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	adapter := &DigitalDbAdapter{
		db:     sqlxDB,
		logger: testLogger,
	}

	userID := "test-user-id"
	locationID := "test-location-id"
	gameID := int64(123)
	userGameID := 456

	t.Run("successfully adds game to digital location", func(t *testing.T) {
		// Set up expectations
		mock.ExpectQuery("SELECT id FROM user_games").
			WithArgs(userID, gameID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userGameID))

		mock.ExpectExec("INSERT INTO digital_game_locations").
			WithArgs(userGameID, locationID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Execute
		err := adapter.AddGameToDigitalLocation(context.Background(), userID, locationID, gameID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("game not found in user's library", func(t *testing.T) {
		// Set up expectations
		mock.ExpectQuery("SELECT id FROM user_games").
			WithArgs(userID, gameID).
			WillReturnError(sql.ErrNoRows)

		// Execute
		err := adapter.AddGameToDigitalLocation(context.Background(), userID, locationID, gameID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "game not found in user's library") {
			t.Errorf("Expected error about game not found, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("game already exists in location", func(t *testing.T) {
		// Set up expectations
		mock.ExpectQuery("SELECT id FROM user_games").
			WithArgs(userID, gameID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userGameID))

		duplicateErr := errors.New("ERROR: duplicate key value violates unique constraint \"digital_game_locations_user_game_id_digital_location_id_key\"")
		mock.ExpectExec("INSERT INTO digital_game_locations").
			WithArgs(userGameID, locationID).
			WillReturnError(duplicateErr)

		// Execute
		err := adapter.AddGameToDigitalLocation(context.Background(), userID, locationID, gameID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "game already exists in this digital location") {
			t.Errorf("Expected error about game already existing, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})
}

func TestRemoveGameFromDigitalLocation(t *testing.T) {
	// Setup
	testLogger := testutils.NewTestLogger()

	// Create mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	adapter := &DigitalDbAdapter{
		db:     sqlxDB,
		logger: testLogger,
	}

	userID := "test-user-id"
	locationID := "test-location-id"
	gameID := int64(123)
	userGameID := 456

	t.Run("successfully removes game from digital location", func(t *testing.T) {
		// Set up expectations
		mock.ExpectQuery("SELECT id FROM user_games").
			WithArgs(userID, gameID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userGameID))

		mock.ExpectExec("DELETE FROM digital_game_locations").
			WithArgs(userGameID, locationID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Execute
		err := adapter.RemoveGameFromDigitalLocation(context.Background(), userID, locationID, gameID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("game not found in user's library", func(t *testing.T) {
		// Set up expectations
		mock.ExpectQuery("SELECT id FROM user_games").
			WithArgs(userID, gameID).
			WillReturnError(sql.ErrNoRows)

		// Execute
		err := adapter.RemoveGameFromDigitalLocation(context.Background(), userID, locationID, gameID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "game not found in user's library") {
			t.Errorf("Expected error about game not found, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("game not found in digital location", func(t *testing.T) {
		// Set up expectations
		mock.ExpectQuery("SELECT id FROM user_games").
			WithArgs(userID, gameID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userGameID))

		mock.ExpectExec("DELETE FROM digital_game_locations").
			WithArgs(userGameID, locationID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Execute
		err := adapter.RemoveGameFromDigitalLocation(context.Background(), userID, locationID, gameID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "game not found in digital location") {
			t.Errorf("Expected error about game not found in location, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})
}

func TestGetGamesByDigitalLocationID(t *testing.T) {
	// Setup
	testLogger := testutils.NewTestLogger()

	// Create mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	adapter := &DigitalDbAdapter{
		db:     sqlxDB,
		logger: testLogger,
	}

	userID := "test-user-id"
	locationID := "test-location-id"

	t.Run("successfully gets games for a digital location", func(t *testing.T) {
		// Set up expectations
		rows := sqlmock.NewRows([]string{"id", "name", "summary", "cover_id", "cover_url", "first_release_date", "rating"}).
			AddRow(1, "Game 1", "Summary 1", "cover1", "http://cover1.jpg", time.Now(), 4.5).
			AddRow(2, "Game 2", "Summary 2", "cover2", "http://cover2.jpg", time.Now(), 4.8)

		mock.ExpectQuery("SELECT g.* FROM games g").
			WithArgs(locationID, userID).
			WillReturnRows(rows)

		// Execute
		games, err := adapter.GetGamesByDigitalLocationID(context.Background(), userID, locationID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(games) != 2 {
			t.Errorf("Expected 2 games, got %d", len(games))
		}
		if games[0].Name != "Game 1" || games[1].Name != "Game 2" {
			t.Errorf("Games not returned correctly: %+v", games)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("no games found for digital location", func(t *testing.T) {
		// Set up expectations
		rows := sqlmock.NewRows([]string{"id", "name", "summary", "cover_id", "cover_url", "first_release_date", "rating"})

		mock.ExpectQuery("SELECT g.* FROM games g").
			WithArgs(locationID, userID).
			WillReturnRows(rows)

		// Execute
		games, err := adapter.GetGamesByDigitalLocationID(context.Background(), userID, locationID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(games) != 0 {
			t.Errorf("Expected 0 games, got %d", len(games))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		// Set up expectations
		dbError := errors.New("database error")

		mock.ExpectQuery("SELECT g.* FROM games g").
			WithArgs(locationID, userID).
			WillReturnError(dbError)

		// Execute
		_, err := adapter.GetGamesByDigitalLocationID(context.Background(), userID, locationID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, dbError) {
			t.Errorf("Expected error %v, got %v", dbError, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})
}

func TestUpdateSubscription(t *testing.T) {
	// Setup
	testLogger := testutils.NewTestLogger()

	// Create mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	adapter := &DigitalDbAdapter{
		db:     sqlxDB,
		logger: testLogger,
	}

	locationID := "test-location-id"
	subscription := models.Subscription{
		ID:             1,
		LocationID:     locationID,
		BillingCycle:   "monthly",
		CostPerCycle:   9.99,
		NextPaymentDate: time.Now(),
		PaymentMethod:  "Credit Card",
	}

	t.Run("successfully updates subscription", func(t *testing.T) {
		// Set up expectations
		mock.ExpectExec("UPDATE digital_location_subscriptions").
			WithArgs(
				subscription.BillingCycle,
				subscription.CostPerCycle,
				sqlmock.AnyArg(), // NextPaymentDate can be any time
				subscription.PaymentMethod,
				sqlmock.AnyArg(), // UpdatedAt can be any time
				subscription.LocationID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Execute
		err := adapter.UpdateSubscription(context.Background(), subscription)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("subscription not found", func(t *testing.T) {
		// Set up expectations
		mock.ExpectExec("UPDATE digital_location_subscriptions").
			WithArgs(
				subscription.BillingCycle,
				subscription.CostPerCycle,
				sqlmock.AnyArg(), // NextPaymentDate can be any time
				subscription.PaymentMethod,
				sqlmock.AnyArg(), // UpdatedAt can be any time
				subscription.LocationID,
			).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Execute
		err := adapter.UpdateSubscription(context.Background(), subscription)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "subscription not found") {
			t.Errorf("Expected error about subscription not found, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		// Set up expectations
		dbError := errors.New("database error")
		mock.ExpectExec("UPDATE digital_location_subscriptions").
			WithArgs(
				subscription.BillingCycle,
				subscription.CostPerCycle,
				sqlmock.AnyArg(), // NextPaymentDate can be any time
				subscription.PaymentMethod,
				sqlmock.AnyArg(), // UpdatedAt can be any time
				subscription.LocationID,
			).
			WillReturnError(dbError)

		// Execute
		err := adapter.UpdateSubscription(context.Background(), subscription)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, dbError) {
			t.Errorf("Expected error %v, got %v", dbError, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})
}

func TestRemoveSubscription(t *testing.T) {
	// Setup
	testLogger := testutils.NewTestLogger()

	// Create mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	adapter := &DigitalDbAdapter{
		db:     sqlxDB,
		logger: testLogger,
	}

	locationID := "test-location-id"

	t.Run("successfully removes subscription", func(t *testing.T) {
		// Set up expectations
		mock.ExpectExec("DELETE FROM digital_location_subscriptions").
			WithArgs(locationID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Execute
		err := adapter.RemoveSubscription(context.Background(), locationID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("subscription not found", func(t *testing.T) {
		// Set up expectations
		mock.ExpectExec("DELETE FROM digital_location_subscriptions").
			WithArgs(locationID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Execute
		err := adapter.RemoveSubscription(context.Background(), locationID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "subscription not found") {
			t.Errorf("Expected error about subscription not found, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		// Set up expectations
		dbError := errors.New("database error")
		mock.ExpectExec("DELETE FROM digital_location_subscriptions").
			WithArgs(locationID).
			WillReturnError(dbError)

		// Execute
		err := adapter.RemoveSubscription(context.Background(), locationID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, dbError) {
			t.Errorf("Expected error %v, got %v", dbError, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})
}
