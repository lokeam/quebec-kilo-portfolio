//go:build !ide
// +build !ide

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
	"github.com/lib/pq"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
- GetSingleDigitalLocation:
  - Successfully retrieves a valid location
  - Returns error when location not found
  - Handles db errors
- GetAllDigitalLocations:
  - Successfully retrieves all locations for a user
  - Returns empty slice when no locations exist
  - Handles db errors
- CreateDigitalLocation:
  - Successfully adds new location
  - Handles db errors
- UpdateDigitalLocation:
  - Successfully updates a location
  - Returns error when location not found
  - Handles db errors
  - Returns error when user is not authorized
- DeleteDigitalLocation:
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
	t.Run(`GetSingleDigitalLocation - Successfully retrieves a valid location`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		now := time.Now()

		// Only include real DB columns - remove 'items' since it's not in the actual query
		rows := sqlmock.NewRows([]string{
			"id", "user_id", "name", "is_subscription", "is_active", "url", "payment_method", "created_at", "updated_at",
			"disk_size_value", "disk_size_unit",
		}).
			AddRow(locationID, userID, "Test Location", false, true, "https://example.com", "paypal", now, now, nil, nil)

		mock.ExpectQuery("SELECT (.+) FROM digital_locations").
			WithArgs(locationID, userID).
			WillReturnRows(rows)

		// Also expect the subscription query, which returns no rows
		mock.ExpectQuery("SELECT (.+) FROM digital_location_subscriptions").
			WithArgs(locationID).
			WillReturnError(sql.ErrNoRows)

		// Execute
		location, err := adapter.GetSingleDigitalLocation(context.Background(), userID, locationID)

		// Manually set the features for verification
		expectedLocation := models.DigitalLocation{
			ID:        locationID,
			UserID:    userID,
			Name:      "Test Location",
			IsSubscription: false,
			IsActive:  true,
			URL:       "https://example.com",
			PaymentMethod: "paypal",
			CreatedAt: now,
			UpdatedAt: now,
			DiskSizeValue: nil,
			DiskSizeUnit: nil,
			Items:     []models.Game{}, // Will be empty
			Subscription: nil,
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
	t.Run(`GetSingleDigitalLocation - Returns error when location not found`, func(t *testing.T) {
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
		_, err = adapter.GetSingleDigitalLocation(context.Background(), userID, locationID)

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
	t.Run(`GetSingleDigitalLocation - Handles database errors`, func(t *testing.T) {
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
		_, err = adapter.GetSingleDigitalLocation(context.Background(), userID, locationID)

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
	t.Run("GetAllDigitalLocations - Successfully retrieves all locations", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		now := time.Now()

		// Use the actual columns returned by GetLocationsWithSubscriptionDataQuery - include 'items'
		rows := sqlmock.NewRows([]string{
			"id", "user_id", "name", "is_subscription", "is_active", "url", "payment_method", "created_at", "updated_at",
			"items", "sub_id", "billing_cycle", "cost_per_cycle", "anchor_date", "last_payment_date", "next_payment_date", "payment_method", "sub_created_at", "sub_updated_at",
		}).
			AddRow(
				"test-location-id",
				userID,
				"Test Location",
				false,
				true,
				"https://test.com",
				"paypal",
				now,
				now,
				[]byte("[]"), // items as JSON
				nil, // sub_id
				nil, // billing_cycle
				nil, // cost_per_cycle
				nil, // anchor_date
				nil, // last_payment_date
				nil, // next_payment_date
				nil, // payment_method (subscription)
				nil, // sub_created_at
				nil, // sub_updated_at
			)

		mock.ExpectQuery("SELECT (.+) FROM digital_locations").
			WithArgs(userID).
			WillReturnRows(rows)

		// Execute the function
		locations, err := adapter.GetAllDigitalLocations(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != 1 {
			t.Errorf("Expected 1 location, got %d", len(locations))
		}
		if locations[0].Name != "Test Location" {
			t.Errorf("Expected location name 'Test Location', got %s", locations[0].Name)
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
	t.Run(`GetAllDigitalLocations - Returns empty slice when no locations exist`, func(t *testing.T) {
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
		locations, err := adapter.GetAllDigitalLocations(context.Background(), userID)

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
	t.Run(`GetAllDigitalLocations - Handles database errors`, func(t *testing.T) {
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
		_, err = adapter.GetAllDigitalLocations(context.Background(), userID)

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
	t.Run(`CreateDigitalLocation - Successfully adds a new location`, func(t *testing.T) {
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
		// Check if a location with the same name already exists - returns no rows
		mock.ExpectQuery("SELECT id FROM digital_locations").
			WithArgs(userID, "New Location").
			WillReturnError(sql.ErrNoRows)

		// Begin transaction
		mock.ExpectBegin()

		// Insert the location
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_subscription", "is_active", "url", "payment_method", "created_at", "updated_at"}).
			AddRow(locationID, userID, "New Location", false, true, "https://example.com", "", now, now)

		mock.ExpectQuery("INSERT INTO digital_locations").
			WithArgs(
				locationID,
				userID,
				"New Location",
				false,
				true,
				"https://example.com",
				"",
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnRows(rows)

		// Commit transaction
		mock.ExpectCommit()

		// Execute
		result, err := adapter.CreateDigitalLocation(context.Background(), userID, location)

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
	t.Run("CreateDigitalLocation - Handles database errors", func(t *testing.T) {
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
		// Check if a location with the same name already exists - returns no rows
		mock.ExpectQuery("SELECT id FROM digital_locations").
			WithArgs(userID, "New Location").
			WillReturnError(sql.ErrNoRows)

		// Begin transaction
		mock.ExpectBegin()

		// Insert fails
		mock.ExpectQuery("INSERT INTO digital_locations").
			WithArgs(
				locationID,
				userID,
				"New Location",
				false,
				true,
				"https://example.com",
				"",
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnError(dbError)

		// Rollback transaction
		mock.ExpectRollback()

		// Execute
		_, err = adapter.CreateDigitalLocation(context.Background(), userID, location)

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
			URL:      "https://example.com",
		}

		// Set up mock expectations
		mock.ExpectExec("UPDATE digital_locations").
			WithArgs(
				"Updated Location",
				true,
				"https://example.com",
				sqlmock.AnyArg(), // updated_at
				false, // is_subscription
				"", // payment_method
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
			URL:      "https://example.com",
		}

		// Set up mock expectations
		mock.ExpectExec("UPDATE digital_locations").
			WithArgs(
				"Updated Location",
				true,
				"https://example.com",
				sqlmock.AnyArg(), // updated_at
				false, // is_subscription
				"", // payment_method
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
			URL:      "https://example.com",
		}

		// Set up mock expectations
		mock.ExpectExec("UPDATE digital_locations").
			WithArgs(
				"Updated Location",
				true,
				"https://example.com",
				sqlmock.AnyArg(), // updated_at
				false, // is_subscription
				"", // payment_method
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
			URL:      "https://example.com",
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
	t.Run(`DeleteDigitalLocation - Successfully removes a location`, func(t *testing.T) {
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
		locationIDs := []string{locationID}

		// Set up expectations
		mock.ExpectBegin()

		// Check if location exists and belongs to user
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM digital_locations").
			WithArgs(pq.Array(locationIDs), userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Execute delete
		mock.ExpectExec("WITH deleted_related AS").
			WithArgs(pq.Array(locationIDs), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Clean up orphaned user games
		mock.ExpectExec("DELETE FROM user_games").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit()

		// Execute
		count, err := adapter.DeleteDigitalLocation(context.Background(), userID, locationIDs)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 location to be deleted, got %d", count)
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
	t.Run(`DeleteDigitalLocation - Returns error when location not found`, func(t *testing.T) {
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
		locationIDs := []string{locationID}

		// Set up expectations
		mock.ExpectBegin()

		// Check if location exists - not found
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM digital_locations").
			WithArgs(pq.Array(locationIDs), userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectRollback()

		// Execute
		count, err := adapter.DeleteDigitalLocation(context.Background(), userID, locationIDs)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "one or more locations not found") {
			t.Errorf("Expected error about location not found, got %v", err)
		}
		if count != 0 {
			t.Errorf("Expected 0 locations to be deleted, got %d", count)
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
	t.Run(`DeleteDigitalLocation - Handles database errors during deletion`, func(t *testing.T) {
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
		locationIDs := []string{locationID}

		// Set up expectations
		mock.ExpectBegin()

		// Check if location exists
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM digital_locations").
			WithArgs(pq.Array(locationIDs), userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Delete the location and get an error
		dbError := errors.New("database error")
		mock.ExpectExec("WITH deleted_related AS").
			WithArgs(pq.Array(locationIDs), userID).
			WillReturnError(dbError)

		mock.ExpectRollback()

		// Execute
		count, err := adapter.DeleteDigitalLocation(context.Background(), userID, locationIDs)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, dbError) {
			t.Errorf("Expected error %v, got %v", dbError, err)
		}
		if count != 0 {
			t.Errorf("Expected 0 locations to be deleted, got %d", count)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("returns error for empty location IDs", func(t *testing.T) {
		// Setup
		testLogger := testutils.NewTestLogger()
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		adapter := &DigitalDbAdapter{
			db:     sqlxDB,
			logger: testLogger,
		}

		// Execute
		count, err := adapter.DeleteDigitalLocation(context.Background(), "test-user-id", []string{})

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "no location IDs provided") {
			t.Errorf("Expected error about no location IDs, got %v", err)
		}
		if count != 0 {
			t.Errorf("Expected 0 locations to be deleted, got %d", count)
		}
	})

	t.Run("returns error for empty user ID", func(t *testing.T) {
		// Setup
		testLogger := testutils.NewTestLogger()
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		adapter := &DigitalDbAdapter{
			db:     sqlxDB,
			logger: testLogger,
		}

		// Execute
		count, err := adapter.DeleteDigitalLocation(context.Background(), "", []string{"test-location-id"})

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "user ID cannot be empty") {
			t.Errorf("Expected error about empty user ID, got %v", err)
		}
		if count != 0 {
			t.Errorf("Expected 0 locations to be deleted, got %d", count)
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
				now := time.Now()
				rows := sqlmock.NewRows([]string{"id", "digital_location_id", "billing_cycle", "cost_per_cycle", "anchor_date", "last_payment_date", "next_payment_date", "payment_method", "created_at", "updated_at"}).
					AddRow(1, "test-location-id", "1 month", 9.99, now, now, now, "Visa", now, now)
				mock.ExpectQuery("SELECT id, digital_location_id, billing_cycle, cost_per_cycle, anchor_date, last_payment_date, next_payment_date, payment_method, created_at, updated_at FROM digital_location_subscriptions WHERE digital_location_id = \\$1").
					WithArgs("test-location-id").
					WillReturnRows(rows)
			},
			expectedResult: &models.Subscription{
				ID:             1,
				LocationID:     "test-location-id",
				BillingCycle:   "1 month",
				CostPerCycle:   9.99,
				PaymentMethod:  "Visa",
			},
			expectedError: nil,
		},
		{
			name:       "Subscription not found",
			locationID: "non-existent-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, digital_location_id, billing_cycle, cost_per_cycle, anchor_date, last_payment_date, next_payment_date, payment_method, created_at, updated_at FROM digital_location_subscriptions WHERE digital_location_id = \\$1").
					WithArgs("non-existent-id").
					WillReturnError(sql.ErrNoRows)
			},
			expectedResult: nil,
			expectedError: sql.ErrNoRows, // Expect error when not found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			adapter := &DigitalDbAdapter{db: sqlxDB, logger: testutils.NewTestLogger()}

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			sub, err := adapter.GetSubscription(context.Background(), tt.locationID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, sub)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, sub)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddSubscription(t *testing.T) {
	now := time.Now()
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
				BillingCycle:   "1 month",
				CostPerCycle:   9.99,
				AnchorDate:     now,
				PaymentMethod:  "Visa",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				otherTime := time.Now()
				rows := sqlmock.NewRows([]string{"id", "digital_location_id", "billing_cycle", "cost_per_cycle", "anchor_date", "last_payment_date", "next_payment_date", "payment_method", "created_at", "updated_at"}).
					AddRow(1, "test-location-id", "1 month", 9.99, now, otherTime, otherTime, "Visa", otherTime, otherTime)
				mock.ExpectQuery("INSERT INTO digital_location_subscriptions").
					WithArgs(
						"test-location-id",
						"1 month",
						9.99,
						now, // AnchorDate
						"Visa",
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
					).
					WillReturnRows(rows)
			},
			expectedResult: &models.Subscription{
				ID:             1,
				LocationID:     "test-location-id",
				BillingCycle:   "1 month",
				CostPerCycle:   9.99,
				AnchorDate:     now,
				PaymentMethod:  "Visa",
			},
			expectedError: nil,
		},
		{
			name: "Failed to add subscription",
			subscription: models.Subscription{
				LocationID:     "test-location-id",
				BillingCycle:   "1 month",
				CostPerCycle:   9.99,
				AnchorDate:     now,
				PaymentMethod:  "Visa",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO digital_location_subscriptions").
					WithArgs(
						"test-location-id",
						"1 month",
						9.99,
						now, // AnchorDate
						"Visa",
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
					).
					WillReturnError(sql.ErrConnDone)
			},
			expectedResult: nil,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			adapter := &DigitalDbAdapter{db: sqlxDB, logger: testutils.NewTestLogger()}

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			result, err := adapter.CreateSubscription(context.Background(), tt.subscription)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
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
				rows := sqlmock.NewRows([]string{"id", "digital_location_id", "amount", "payment_date", "payment_method", "transaction_id", "created_at"}).
					AddRow(1, "test-location-id", 9.99, time.Now(), "Visa", "tx123", time.Now()).
					AddRow(2, "test-location-id", 9.99, time.Now(), "Visa", "tx124", time.Now())
				mock.ExpectQuery("SELECT id, digital_location_id, amount, payment_date, payment_method, transaction_id, created_at FROM digital_location_payments WHERE digital_location_id = \\$1 ORDER BY payment_date DESC").
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
				mock.ExpectQuery("SELECT id, digital_location_id, amount, payment_date, payment_method, transaction_id, created_at FROM digital_location_payments WHERE digital_location_id = \\$1 ORDER BY payment_date DESC").
					WithArgs("non-existent-id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "digital_location_id", "amount", "payment_date", "payment_method", "transaction_id", "created_at"}))
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

			result, err := adapter.GetAllPayments(context.Background(), tt.locationID)

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
	// Create a mock database
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a mock logger
	log := testutils.NewTestLogger()

	// Create sqlx wrapper around mock database
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Create the adapter directly with the mock database
	adapter := &DigitalDbAdapter{
		db:     sqlxDB,
		logger: log,
	}

	// Verify the adapter was created correctly
	if adapter.db == nil {
		t.Error("Expected adapter.db to be non-nil")
	}
	if adapter.logger == nil {
		t.Error("Expected adapter.logger to be non-nil")
	}
}

func TestGetUserDigitalLocations(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

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

	t.Run("success", func(t *testing.T) {
		userID := "test-user-id"
		now := time.Now()

		// Use the actual columns from GetLocationsWithSubscriptionDataQuery - include 'items'
		rows := sqlmock.NewRows([]string{
			"id", "user_id", "name", "is_subscription", "is_active", "url", "payment_method", "created_at", "updated_at",
			"items", "sub_id", "billing_cycle", "cost_per_cycle", "anchor_date", "last_payment_date", "next_payment_date", "payment_method", "sub_created_at", "sub_updated_at",
		}).
			AddRow(
				"test-location-id",
				userID,
				"Test Location",
				false,
				true,
				"https://test.com",
				"paypal",
				now,
				now,
				[]byte("[]"), // items as JSON
				nil, // sub_id
				nil, // billing_cycle
				nil, // cost_per_cycle
				nil, // anchor_date
				nil, // last_payment_date
				nil, // next_payment_date
				nil, // payment_method (subscription)
				nil, // sub_created_at
				nil, // sub_updated_at
			)

		mock.ExpectQuery("SELECT (.+) FROM digital_locations").
			WithArgs(userID).
			WillReturnRows(rows)

		// Execute
		locations, err := adapter.GetAllDigitalLocations(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if locations == nil {
			t.Errorf("Expected non-nil locations")
		}
		if len(locations) == 0 {
			t.Errorf("Expected at least one location")
		}
		if locations[0].Name != "Test Location" {
			t.Errorf("Expected location name 'Test Location', got %s", locations[0].Name)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})
}

func TestGetDigitalLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up mock expectations - use the actual columns from the database
	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "name", "is_subscription", "is_active", "url", "payment_method",
		"created_at", "updated_at", "disk_size_value", "disk_size_unit",
	}).
		AddRow(
			"test-location-id",
			"test-user",
			"Test Location",
			false,
			true,
			"https://test.com",
			"paypal",
			now,
			now,
			nil,
			nil,
		)

	mock.ExpectQuery("SELECT (.+) FROM digital_locations").
		WithArgs("test-location-id", "test-user").
		WillReturnRows(rows)

	// Also expect the subscription query, which returns no rows
	mock.ExpectQuery("SELECT (.+) FROM digital_location_subscriptions").
		WithArgs("test-location-id").
		WillReturnError(sql.ErrNoRows)

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
		location, err := adapter.GetSingleDigitalLocation(context.Background(), "test-user", "test-location-id")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if location.ID == "" {
			t.Error("Expected non-empty location ID")
		}
		if location.Name != "Test Location" {
			t.Errorf("Expected location name 'Test Location', got %s", location.Name)
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

	// Set up mock expectations for checking existing location name
	mock.ExpectQuery("SELECT id FROM digital_locations").
		WithArgs("test-user", "Test Location").
		WillReturnError(sql.ErrNoRows)

	// Begin transaction
	mock.ExpectBegin()

	// Set up mock expectations for location insertion
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_subscription", "is_active", "url", "payment_method", "created_at", "updated_at"}).
		AddRow("test-location-id", "test-user", "Test Location", false, true, "https://test.com", "", now, now)
	mock.ExpectQuery("INSERT INTO digital_locations").
		WithArgs(
			"test-location-id",
			"test-user",
			"Test Location",
			false,
			true,
			"https://test.com",
			"",
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnRows(rows)

	// Commit transaction
	mock.ExpectCommit()

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
			ID:       "test-location-id",
			Name:     "Test Location",
			IsActive: true,
			URL:      "https://test.com",
		}

		location, err := adapter.CreateDigitalLocation(context.Background(), "test-user", testLocation)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if location.ID == "" {
			t.Error("Expected non-empty location ID")
		}
		if location.Name != "Test Location" {
			t.Errorf("Expected location name 'Test Location', got %s", location.Name)
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
			AddRow(int64(1), "Game 1", "Summary 1", int64(123), "http://cover1.jpg", int64(1735142400), 4.5).
			AddRow(int64(2), "Game 2", "Summary 2", int64(456), "http://cover2.jpg", int64(1735142400), 4.8)

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
	now := time.Now()
	tests := []struct {
		name           string
		subscription   models.Subscription
		mockSetup      func(sqlmock.Sqlmock)
		expectedError  error
	}{
		{
			name: "successfully updates subscription",
			subscription: models.Subscription{
				LocationID:     "test-location-id",
				BillingCycle:   "1 month",
				CostPerCycle:   9.99,
				AnchorDate:     now,
				PaymentMethod:  "Credit Card",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE digital_location_subscriptions").
					WithArgs(
						"1 month",
						9.99,
						sqlmock.AnyArg(), // anchor_date
						"Credit Card",
						sqlmock.AnyArg(), // updated_at
						"test-location-id",
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			name: "subscription not found",
			subscription: models.Subscription{
				LocationID:     "test-location-id",
				BillingCycle:   "1 month",
				CostPerCycle:   9.99,
				AnchorDate:     now,
				PaymentMethod:  "Credit Card",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE digital_location_subscriptions").
					WithArgs(
						"1 month",
						9.99,
						sqlmock.AnyArg(), // anchor_date
						"Credit Card",
						sqlmock.AnyArg(), // updated_at
						"test-location-id",
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: sql.ErrNoRows,
		},
		{
			name: "database error",
			subscription: models.Subscription{
				LocationID:     "test-location-id",
				BillingCycle:   "1 month",
				CostPerCycle:   9.99,
				AnchorDate:     now,
				PaymentMethod:  "Credit Card",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE digital_location_subscriptions").
					WithArgs(
						"1 month",
						9.99,
						sqlmock.AnyArg(), // anchor_date
						"Credit Card",
						sqlmock.AnyArg(), // updated_at
						"test-location-id",
					).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			adapter := &DigitalDbAdapter{db: sqlxDB, logger: testutils.NewTestLogger()}

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			err = adapter.UpdateSubscription(context.Background(), tt.subscription)
			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
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
		err := adapter.DeleteSubscription(context.Background(), locationID)

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
		err := adapter.DeleteSubscription(context.Background(), locationID)

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
		err := adapter.DeleteSubscription(context.Background(), locationID)

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
