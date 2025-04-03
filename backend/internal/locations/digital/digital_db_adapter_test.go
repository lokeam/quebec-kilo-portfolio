package digital

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext_test"
	"github.com/lokeam/qko-beta/internal/models"
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
	baseAppCtx := appcontext_test.NewTestingAppContext("test-token", nil)

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
			logger: baseAppCtx.Logger,
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

		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at"}).
			AddRow(locationID, userID, "Test Location", true, "https://example.com", now, now)

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

		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at"}).
			AddRow("loc1", userID, "Location 1", true, "https://example1.com", now, now).
			AddRow("loc2", userID, "Location 2", true, "https://example2.com", now, now)

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
	t.Run(`GetUserDigitalLocations - Returns empty slice when locations exist`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"

		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at"})

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
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "is_active", "url", "created_at", "updated_at"}).
			AddRow(locationID, userID, "New Location", true, "https://example.com", now, now)

		mock.ExpectQuery("INSERT INTO digital_locations").
			WithArgs(
				locationID,
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
		mock.ExpectQuery("INSERT INTO digital_locations").
			WithArgs(
				locationID,
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
