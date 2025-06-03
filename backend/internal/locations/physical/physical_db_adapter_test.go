package physical

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lokeam/qko-beta/internal/testutils"
)

/*
	Behavior:
	- Retrieving a single physical location for a user
	- Retrieving all physical locations for a user
	- Adding a new physical location
	- Updating an existing physical location
	- Removing a physical location
	- Handling db errors
	- Ensuring user may only access their own locations

	Scenarios:
	- GetSinglePhysicalLocation:
		- Successfully retrieves a valid location
		- Returns error when location not found
		- Handles db errors
	- GetAllPhysicalLocations:
		- Successfully retrieves all locations for a user
		- Returns empty slice when no locations exist
		- Handles db errors
	- CreatePhysicalLocation:
		- Successfully adds new location
		- Handles db errors
	- UpdatePhysicalLocation:
		- Successfully updates a location
		- Returns error when location not found
		- Handles db errors
	- DeletePhysicalLocation:
		- Successfully removes a location
		- Returns errors when location not found
		- Handles db errors
*/

func TestPhysicalDbAdapter(t *testing.T) {
	// Set up base app context for testing
	baseAppCtx := appcontext_test.NewTestingAppContext("test-token", nil)

	// Create mock DB + adapter
	setupMockDB := func() (*PhysicalDbAdapter, sqlmock.Sqlmock, error) {
		// Create mock sqldatabase
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			return nil, nil, err
		}

		// Create a sqlx wrapper around mock data
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		// Create the adapter with the mock DB
		adapter := &PhysicalDbAdapter{
			db:     sqlxDB,
			logger: baseAppCtx.Logger,
		}

		return adapter, mock, nil
	}

	/*
		GIVEN a request to get a specific physical location
		WHEN the location exists in the database
		THEN the adapter returns the location
	*/
	t.Run(`GetSinglePhysicalLocation - Successfully retrieves a valid location`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		now := time.Now()

		expectedLocation := models.PhysicalLocation{
			ID:               locationID,
			UserID:           userID,
			Name:             "Test Location",
			Label:            "Home",
			LocationType:     "Residence",
			MapCoordinates:   "123.456,789.012",
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "label", "location_type", "map_coordinates", "created_at", "updated_at"}).
			AddRow(locationID, userID, "Test Location", "Home", "Residence", "123,456", now, now)

		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnRows(rows)

		// Execute
		location, err := adapter.GetSinglePhysicalLocation(context.Background(), userID, locationID)

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
		GIVEN a request to get a specific physical location
		WHEN the location does not exist in the db
		THEN the adapter returns an error
	*/
	t.Run(``, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"

		// Setup mock expectations
		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnError(sql.ErrNoRows)

		// Execute
		_, err = adapter.GetSinglePhysicalLocation(context.Background(), userID, locationID)

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
		GIVEN a request to get A SPECIFIC physical location
		WHEN the db returns an error
		THEN the adapter also returns an error
	*/
	t.Run(`GetSinglePhysicalLocation - Handles database errors`, func(t *testing.T) {
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
		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnError(dbError)

		// Execute the fucntion
		_, err = adapter.GetSinglePhysicalLocation(context.Background(), userID, locationID)

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
		GIVEN a request to get ALL physical locations for a user
		WHEN locations exist in the db
		THEN the adapter returns all locations
	*/
	t.Run("GetAllPhysicalLocations - Successfully retrieves all locations", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		now := time.Now()

		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "label", "location_type", "map_coordinates", "created_at", "updated_at"}).
			AddRow("loc1", userID, "Location 1", "Home", "Residence", "123,456", now, now).
			AddRow("loc2", userID, "Location 2", "Office", "Work", "789,012", now, now)

		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(userID).
			WillReturnRows(rows)

		// Execute the function
		locations, err := adapter.GetAllPhysicalLocations(context.Background(), userID)

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
		GIVEN a request to get ALL physical locations for a user
		WHEN NO locations exist in the db
		THEN the adapter returns an empty slice
	*/
	t.Run(`GetAllPhysicalLocations - Returns empty slice when locations exist`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"

		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "label", "location_type", "map_coordinates", "created_at", "updated_at"})

		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(userID).
			WillReturnRows(rows)

		// Execute the function
		locations, err := adapter.GetAllPhysicalLocations(context.Background(), userID)

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
		GIVEN a request to get ALL physical locations for a user
		WHEN the database returns an error
		THEN the adapter returns the error
	*/
	t.Run(`GetAllPhysicalLocations - Handles database errors`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		dbError := ErrDatabaseError

		// Set up mock expectations
		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(userID).
			WillReturnError(dbError)

		// Execute
		_, err = adapter.GetAllPhysicalLocations(context.Background(), userID)

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
		GIVEN a request to ADD a new physical location
		WHEN the db operation is successful
		THEN the adapter adds the location and returns it with generated fields
	*/
	t.Run(`CreatePhysicalLocation - Successfully adds a new location`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		now := time.Now()

		location := models.PhysicalLocation{
			ID:             locationID,
			Name:           "New Location",
			Label:          "Home",
			LocationType:   "Residence",
			MapCoordinates: "123,456",
		}

		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "label", "location_type", "map_coordinates", "created_at", "updated_at"}).
			AddRow(locationID, userID, "New Location", "Home", "Residence", "123,456", now, now)

		mock.ExpectQuery("INSERT INTO physical_locations").
			WithArgs(locationID, userID, "New Location", "Home", "Residence", "123,456", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Execute
		result, err := adapter.CreatePhysicalLocation(context.Background(), userID, location)

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
		GIVEN a request to ADD a new physical location
		WHEN the db returns an error
		THEN then the adapter also returns an error
	*/
	t.Run("CreatePhysicalLocation - Handles database errors", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		dbError := ErrDatabaseError

		location := models.PhysicalLocation{
			ID:             locationID,
			Name:           "New Location",
			Label:          "Home",
			LocationType:   "Residence",
			MapCoordinates: "123,456",
		}

		// Set up mock expectations
		mock.ExpectQuery("INSERT INTO physical_locations").
			WithArgs(locationID, userID, "New Location", "Home", "Residence", "123,456", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(dbError)

		// Execute the function
		_, err = adapter.CreatePhysicalLocation(context.Background(), userID, location)

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
		GIVEN a request to UPDATE a physical location
		WHEN the location exists AND belongs to the user
		THEN the adapter updates the location
	*/
	t.Run(`UpdatePhysicalLocation - Successfully updates a location`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		now := time.Now()

		location := models.PhysicalLocation{
			ID:             locationID,
			UserID:         userID,
			Name:           "Updated Location",
			Label:          "Updated Label",
			LocationType:   "Updated Type",
			MapCoordinates: "789,012",
		}

		// Setup mock expectations
		mock.ExpectQuery("SELECT id FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(locationID))

		mock.ExpectQuery("UPDATE physical_locations").
			WithArgs(
				locationID,
				userID,
				"Updated Location",
				"Updated Label",
				"Updated Type",
				"789,012",
			).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "user_id", "name", "label", "location_type", "map_coordinates", "created_at", "updated_at",
			}).AddRow(
				locationID, userID, "Updated Location", "Updated Label", "Updated Type", "789,012", now, now,
			))

		// Execute
		updatedLocation, err := adapter.UpdatePhysicalLocation(context.Background(), userID, location)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, location.ID, updatedLocation.ID)
		assert.Equal(t, location.UserID, updatedLocation.UserID)
		assert.Equal(t, location.Name, updatedLocation.Name)
		assert.Equal(t, location.Label, updatedLocation.Label)
		assert.Equal(t, location.LocationType, updatedLocation.LocationType)
		assert.Equal(t, location.MapCoordinates, updatedLocation.MapCoordinates)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run(`UpdatePhysicalLocation - Returns error when location not found`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"

		location := models.PhysicalLocation{
			ID:             locationID,
			UserID:         userID,
			Name:           "Updated Location",
			Label:          "Updated Label",
			LocationType:   "Updated Type",
			MapCoordinates: "789,012",
		}

		// Set up mock expectations
		mock.ExpectQuery("SELECT id FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnError(sql.ErrNoRows)

		// Execute
		_, err = adapter.UpdatePhysicalLocation(context.Background(), userID, location)

		// Verify
		assert.Error(t, err)
		assert.Equal(t, ErrUnauthorizedLocation, err)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run(`UpdatePhysicalLocation - Handles database errors`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		dbError := ErrDatabaseError

		location := models.PhysicalLocation{
			ID:             locationID,
			UserID:         userID,
			Name:           "Updated Location",
			Label:          "Updated Label",
			LocationType:   "Updated Type",
			MapCoordinates: "789,012",
		}

		// Set up mock expectations
		mock.ExpectQuery("SELECT id FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(locationID))

		mock.ExpectQuery("UPDATE physical_locations").
			WithArgs(
				locationID,
				userID,
				"Updated Location",
				"Updated Label",
				"Updated Type",
				"789,012",
			).
			WillReturnError(dbError)

		// Execute the function
		_, err = adapter.UpdatePhysicalLocation(context.Background(), userID, location)

		// Verify
		assert.Error(t, err)
		assert.True(t, errors.Is(err, dbError))
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	t.Run("UpdatePhysicalLocation - Returns error when user is not authorized", func(t *testing.T) {
		// Setup
		adapter, _, err := setupMockDB()
		require.NoError(t, err)
		defer adapter.db.Close()

		userID := "test-user-id"
		wrongUserID := "wrong-user-id"
		locationID := "test-location-id"

		location := models.PhysicalLocation{
			ID:             locationID,
			UserID:         wrongUserID,
			Name:           "Updated Location",
			Label:          "Updated Label",
			LocationType:   "Updated Type",
			MapCoordinates: "789,012",
		}

		// Execute the function
		_, err = adapter.UpdatePhysicalLocation(context.Background(), userID, location)

		// Verify
		assert.Error(t, err)
		assert.Equal(t, ErrUnauthorizedLocation, err)
	})

	/*
		GIVEN a request to REMOVE a physical location
		WHEN the location exists AND belongs to the user
		THEN the adapter removes the location
	*/
	t.Run(`DeletePhysicalLocation - Successfully removes a location`, func(t *testing.T) {
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
		mock.ExpectQuery("SELECT id FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnRows(rows)

		// Delete location
		mock.ExpectExec("DELETE FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Commit transaction
		mock.ExpectCommit()

		// Execute
		err = adapter.DeletePhysicalLocation(context.Background(), userID, locationID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to REMOVE a physical location
		WHEN the location doesn't exist or doesn't belong to the user
		THEN the adapter returns an error
	*/
	t.Run(`DeletePhysicalLocation - Returns error when location not found`, func(t *testing.T) {
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
		mock.ExpectQuery("SELECT id FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnError(sql.ErrNoRows)

		// Rollback transaction
		mock.ExpectRollback()

		// Execute the function
		err = adapter.DeletePhysicalLocation(context.Background(), userID, locationID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to REMOVE a physical location
		WHEN the database returns an error during deletion
		THEN the adapter returns the error
	*/
	t.Run(`DeletePhysicalLocation - Handles database errors during deletion`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"
		dbError := ErrDatabaseError

		// Setup mock expectations
		mock.ExpectBegin()

		// Check if the location exists
		rows := sqlmock.NewRows([]string{"id"}).AddRow(locationID)
		mock.ExpectQuery("SELECT id FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnRows(rows)

		// Delete the location and get an error
		mock.ExpectExec("DELETE FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnError(dbError)

		// Rollback the transaction
		mock.ExpectRollback()

		// Execute the function
		err = adapter.DeletePhysicalLocation(context.Background(), userID, locationID)

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

// func setupTestAdapter(t *testing.T) (*PhysicalDbAdapter, func()) {
// 	appCtx := &appcontext.AppContext{
// 		Config: &config.Config{
// 			Postgres: &config.PostgresConfig{
// 				ConnectionString: "postgres://postgres:postgres@localhost:5432/qko_test?sslmode=disable",
// 			},
// 		},
// 		Logger: testutils.NewTestLogger(),
// 	}

// 	adapter, err := NewPhysicalDbAdapter(appCtx)
// 	require.NoError(t, err)

// 	cleanup := func() {
// 		adapter.db.Close()
// 	}

// 	return adapter, cleanup
// }

func setupMockDB() (*PhysicalDbAdapter, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	sqlxDB := sqlx.NewDb(db, "postgres")
	adapter := &PhysicalDbAdapter{
		db:     sqlxDB,
		logger: testutils.NewTestLogger(),
	}

	return adapter, mock, nil
}

func TestUpdatePhysicalLocation_Success(t *testing.T) {
	adapter, mock, err := setupMockDB()
	require.NoError(t, err)
	defer adapter.db.Close()

	ctx := context.Background()
	userID := "user1"
	locationID := "loc1"
	now := time.Now()

	// Mock the authorization check
	mock.ExpectQuery("SELECT id FROM physical_locations").
		WithArgs(locationID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(locationID))

	// Mock the update query
	mock.ExpectQuery("UPDATE physical_locations").
		WithArgs(
			locationID,
			userID,
			"Test Location",
			"Test Label",
			"test",
			"45.5017,-73.5673",
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "name", "label", "location_type", "map_coordinates", "created_at", "updated_at",
		}).AddRow(
			locationID, userID, "Test Location", "Test Label", "test", "45.5017,-73.5673", now, now,
		))

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          userID,
		Name:            "Test Location",
		Label:           "Test Label",
		LocationType:    "test",
		MapCoordinates:  "45.5017,-73.5673",
	}

	updatedLocation, err := adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.NoError(t, err)
	assert.Equal(t, locationID, updatedLocation.ID)
	assert.Equal(t, userID, updatedLocation.UserID)
	assert.Equal(t, "Test Location", updatedLocation.Name)
	assert.Equal(t, "Test Label", updatedLocation.Label)
	assert.Equal(t, "test", updatedLocation.LocationType)
	assert.Equal(t, "45.5017,-73.5673", updatedLocation.MapCoordinates)
}

func TestUpdatePhysicalLocation_Unauthorized(t *testing.T) {
	adapter, mock, err := setupMockDB()
	require.NoError(t, err)
	defer adapter.db.Close()

	ctx := context.Background()
	userID := "user1"
	locationID := "loc1"

	// Mock the database query to return no rows
	mock.ExpectQuery("SELECT id FROM physical_locations").
		WithArgs(locationID, userID).
		WillReturnError(sql.ErrNoRows)

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          "user2", // Different user ID
		Name:            "Test Location",
		Label:           "Test Label",
		LocationType:    "test",
		MapCoordinates:  "45.5017,-73.5673",
	}

	_, err = adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorizedLocation, err)
}

func TestUpdatePhysicalLocation_DatabaseError(t *testing.T) {
	adapter, mock, err := setupMockDB()
	require.NoError(t, err)
	defer adapter.db.Close()

	ctx := context.Background()
	userID := "user1"
	locationID := "loc1"

	// Mock the database query to return an error
	mock.ExpectQuery("SELECT id FROM physical_locations").
		WithArgs(locationID, userID).
		WillReturnError(sql.ErrConnDone)

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          userID,
		Name:            "Test Location",
		Label:           "Test Label",
		LocationType:    "test",
		MapCoordinates:  "45.5017,-73.5673",
	}

	_, err = adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.Error(t, err)
}

func TestUpdatePhysicalLocation_NotFound(t *testing.T) {
	adapter, mock, err := setupMockDB()
	require.NoError(t, err)
	defer adapter.db.Close()

	ctx := context.Background()
	userID := "user1"
	locationID := "loc1"

	// Mock the database query to return no rows
	mock.ExpectQuery("SELECT id FROM physical_locations").
		WithArgs(locationID, userID).
		WillReturnError(sql.ErrNoRows)

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          userID,
		Name:            "Test Location",
		Label:           "Test Label",
		LocationType:    "test",
		MapCoordinates:  "45.5017,-73.5673",
	}

	_, err = adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorizedLocation, err)
}

func TestUpdatePhysicalLocation_InvalidCoordinates(t *testing.T) {
	adapter, mock, err := setupMockDB()
	require.NoError(t, err)
	defer adapter.db.Close()

	ctx := context.Background()
	userID := "user1"
	locationID := "loc1"

	// Mock the authorization check
	mock.ExpectQuery("SELECT id FROM physical_locations").
		WithArgs(locationID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(locationID))

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          userID,
		Name:            "Test Location",
		Label:           "Test Label",
		LocationType:    "test",
		MapCoordinates:  "invalid",
	}

	_, err = adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.Error(t, err)
}

func TestUpdatePhysicalLocation_InvalidBgColor(t *testing.T) {
	adapter, mock, err := setupMockDB()
	require.NoError(t, err)
	defer adapter.db.Close()

	ctx := context.Background()
	userID := "user1"
	locationID := "loc1"

	// Mock the authorization check
	mock.ExpectQuery("SELECT id FROM physical_locations").
		WithArgs(locationID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(locationID))

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          userID,
		Name:            "Test Location",
		Label:           "Test Label",
		LocationType:    "test",
		MapCoordinates:  "45.5017,-73.5673",
	}

	_, err = adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.Error(t, err)
}

func TestUpdatePhysicalLocation_EmptyName(t *testing.T) {
	adapter, mock, err := setupMockDB()
	require.NoError(t, err)
	defer adapter.db.Close()

	ctx := context.Background()
	userID := "user1"
	locationID := "loc1"

	// Mock the authorization check
	mock.ExpectQuery("SELECT id FROM physical_locations").
		WithArgs(locationID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(locationID))

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          userID,
		Name:            "",
		Label:           "Test Label",
		LocationType:    "test",
		MapCoordinates:  "45.5017,-73.5673",
	}

	_, err = adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.Error(t, err)
}

func TestUpdatePhysicalLocation_EmptyLocationType(t *testing.T) {
	adapter, mock, err := setupMockDB()
	require.NoError(t, err)
	defer adapter.db.Close()

	ctx := context.Background()
	userID := "user1"
	locationID := "loc1"

	// Mock the authorization check
	mock.ExpectQuery("SELECT id FROM physical_locations").
		WithArgs(locationID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(locationID))

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          userID,
		Name:            "Test Location",
		Label:           "Test Label",
		LocationType:    "",
		MapCoordinates:  "45.5017,-73.5673",
	}

	_, err = adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.Error(t, err)
}