package physical

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	baseAppCtx := &appcontext.AppContext{
		Logger: testutils.NewTestLogger(),
	}

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
			MapCoordinates:   models.PhysicalMapCoordinates{Coords: "123.456,789.012"},
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		// Expect transaction begin
		mock.ExpectBegin()
		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "label", "location_type", "bg_color", "map_coordinates", "created_at", "updated_at"}).
			AddRow(locationID, userID, "Test Location", "Home", "Residence", "", "123.456,789.012", now, now)
		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnRows(rows)
		// Mock sublocations query
		mock.ExpectQuery("SELECT json_agg").
			WithArgs(locationID).
			WillReturnRows(sqlmock.NewRows([]string{"json_agg"}).AddRow("[]"))
		// Expect commit
		mock.ExpectCommit()

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
	t.Run(`GetSinglePhysicalLocation - Returns error when location not found`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		locationID := "test-location-id"

		// Expect transaction begin
		mock.ExpectBegin()
		// Setup mock expectations
		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnError(sql.ErrNoRows)
		// Expect rollback
		mock.ExpectRollback()

		// Execute
		_, err = adapter.GetSinglePhysicalLocation(context.Background(), userID, locationID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error but got nil")
		}
		if !errors.Is(err, ErrLocationNotFound) {
			t.Errorf("Expected error to contain %v, but got %v", ErrLocationNotFound, err)
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

		// Expect transaction begin
		mock.ExpectBegin()
		// Set up mock expectations
		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(locationID, userID).
			WillReturnError(dbError)
		// Expect rollback
		mock.ExpectRollback()

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

		// Expect transaction begin
		mock.ExpectBegin()
		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "label", "location_type", "bg_color", "map_coordinates", "created_at", "updated_at"}).
			AddRow("loc1", userID, "Location 1", "Home", "Residence", "", "123,456", now, now).
			AddRow("loc2", userID, "Location 2", "Office", "Work", "", "789,012", now, now)
		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(userID).
			WillReturnRows(rows)
		// Mock sublocations queries for each location
		mock.ExpectQuery("SELECT json_agg").
			WithArgs("loc1").
			WillReturnRows(sqlmock.NewRows([]string{"json_agg"}).AddRow("[]"))
		mock.ExpectQuery("SELECT json_agg").
			WithArgs("loc2").
			WillReturnRows(sqlmock.NewRows([]string{"json_agg"}).AddRow("[]"))
		// Expect commit
		mock.ExpectCommit()

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
	t.Run(`GetAllPhysicalLocations - Returns empty slice when no locations exist`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"

		// Expect transaction begin
		mock.ExpectBegin()
		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "label", "location_type", "bg_color", "map_coordinates", "created_at", "updated_at"})
		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(userID).
			WillReturnRows(rows)
		// Expect commit
		mock.ExpectCommit()

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

		// Expect transaction begin
		mock.ExpectBegin()
		// Set up mock expectations
		mock.ExpectQuery("SELECT (.+) FROM physical_locations").
			WithArgs(userID).
			WillReturnError(dbError)
		// Expect rollback
		mock.ExpectRollback()

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
			MapCoordinates: models.PhysicalMapCoordinates{Coords: "123,456"},
		}

		// Expect transaction begin
		mock.ExpectBegin()
		// Set up mock expectations
		// First, expect the user existence check
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id = \\$1\\)").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		// Then expect the duplicate check
		mock.ExpectQuery("SELECT id FROM physical_locations").
			WithArgs(userID, "New Location").
			WillReturnError(sql.ErrNoRows)
		// Finally expect the insert
		rows := sqlmock.NewRows([]string{"id", "user_id", "name", "label", "location_type", "bg_color", "map_coordinates", "created_at", "updated_at"}).
			AddRow(locationID, userID, "New Location", "Home", "Residence", "", "123,456", now, now)
		mock.ExpectQuery("INSERT INTO physical_locations").
			WithArgs(locationID, userID, "New Location", "Home", "Residence", "123,456", "").
			WillReturnRows(rows)
		// Expect commit
		mock.ExpectCommit()

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
			MapCoordinates: models.PhysicalMapCoordinates{Coords: "123,456"},
		}

		// Expect transaction begin
		mock.ExpectBegin()
		// Set up mock expectations
		// First, expect the user existence check
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id = \\$1\\)").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		// Then expect the duplicate check
		mock.ExpectQuery("SELECT id FROM physical_locations").
			WithArgs(userID, "New Location").
			WillReturnError(sql.ErrNoRows)
		// Finally expect the insert to fail
		mock.ExpectQuery("INSERT INTO physical_locations").
			WithArgs(locationID, userID, "New Location", "Home", "Residence", "123,456", "").
			WillReturnError(dbError)
		// Expect rollback
		mock.ExpectRollback()

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
			MapCoordinates: models.PhysicalMapCoordinates{Coords: "789,012"},
		}

		// Expect transaction begin
		mock.ExpectBegin()
		// Setup mock expectations
		mock.ExpectQuery("UPDATE physical_locations").
			WithArgs(
				locationID,
				userID,
				"Updated Location",
				"Updated Label",
				"Updated Type",
				"789,012",
				"",
			).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "user_id", "name", "label", "location_type", "map_coordinates", "bg_color", "created_at", "updated_at",
			}).AddRow(
				locationID, userID, "Updated Location", "Updated Label", "Updated Type", "789,012", "", now, now,
			))
		// Expect commit
		mock.ExpectCommit()

		// Execute
		updatedLocation, err := adapter.UpdatePhysicalLocation(context.Background(), userID, location)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, location.ID, updatedLocation.ID)
		assert.Equal(t, location.UserID, updatedLocation.UserID)
		assert.Equal(t, location.Name, updatedLocation.Name)
		assert.Equal(t, location.Label, updatedLocation.Label)
		assert.Equal(t, location.LocationType, updatedLocation.LocationType)
		assert.Equal(t, location.MapCoordinates.Coords, updatedLocation.MapCoordinates.Coords)
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
			MapCoordinates: models.PhysicalMapCoordinates{Coords: "789,012"},
		}

		// Expect transaction begin
		mock.ExpectBegin()
		// Set up mock expectations
		mock.ExpectQuery("UPDATE physical_locations").
			WithArgs(
				locationID,
				userID,
				"Updated Location",
				"Updated Label",
				"Updated Type",
				"789,012",
				"",
			).
			WillReturnError(sql.ErrNoRows)
		// Expect rollback
		mock.ExpectRollback()

		// Execute
		_, err = adapter.UpdatePhysicalLocation(context.Background(), userID, location)

		// Verify
		assert.Error(t, err)
		assert.Equal(t, ErrLocationNotFound, err)
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
			MapCoordinates: models.PhysicalMapCoordinates{Coords: "789,012"},
		}

		// Expect transaction begin
		mock.ExpectBegin()
		// Set up mock expectations
		mock.ExpectQuery("UPDATE physical_locations").
			WithArgs(
				locationID,
				userID,
				"Updated Location",
				"Updated Label",
				"Updated Type",
				"789,012",
				"",
			).
			WillReturnError(dbError)
		// Expect rollback
		mock.ExpectRollback()

		// Execute the function
		_, err = adapter.UpdatePhysicalLocation(context.Background(), userID, location)

		// Verify
		assert.Error(t, err)
		assert.True(t, errors.Is(err, dbError))
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
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

		// Expect transaction begin
		mock.ExpectBegin()
		// Check if location exists - use the actual COUNT query from the implementation
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM physical_locations").
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		// Get sublocations (empty in this case)
		mock.ExpectQuery("SELECT id FROM sublocations").
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		// Delete location
		mock.ExpectExec("DELETE FROM physical_locations").
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		// Expect commit
		mock.ExpectCommit()

		// Execute
		rowsAffected, err := adapter.DeletePhysicalLocation(context.Background(), userID, []string{locationID})

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if rowsAffected != 1 {
			t.Errorf("Expected 1 row affected, got %d", rowsAffected)
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

		// Expect transaction begin
		mock.ExpectBegin()
		// Check if location exists - not found
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM physical_locations").
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// Expect rollback
		mock.ExpectRollback()

		// Execute the function
		_, err = adapter.DeletePhysicalLocation(context.Background(), userID, []string{locationID})

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

		// Expect transaction begin
		mock.ExpectBegin()
		// Check if the location exists
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM physical_locations").
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		// Get sublocations (empty in this case)
		mock.ExpectQuery("SELECT id FROM sublocations").
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		// Delete the location and get an error
		mock.ExpectExec("DELETE FROM physical_locations").
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnError(dbError)
		// Expect rollback
		mock.ExpectRollback()

		// Execute the function
		_, err = adapter.DeletePhysicalLocation(context.Background(), userID, []string{locationID})

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

	// Mock the update query
	mock.ExpectBegin()
	mock.ExpectQuery("UPDATE physical_locations").
		WithArgs(
			locationID,
			userID,
			"Test Location",
			"Test Label",
			"test",
			"45.5017,-73.5673",
			"",
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "name", "label", "location_type", "map_coordinates", "bg_color", "created_at", "updated_at",
		}).AddRow(
			locationID, userID, "Test Location", "Test Label", "test", "45.5017,-73.5673", "", now, now,
		))
	mock.ExpectCommit()

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          userID,
		Name:            "Test Location",
		Label:           "Test Label",
		LocationType:    "test",
		MapCoordinates:  models.PhysicalMapCoordinates{Coords: "45.5017,-73.5673"},
	}

	updatedLocation, err := adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.NoError(t, err)
	assert.Equal(t, locationID, updatedLocation.ID)
	assert.Equal(t, userID, updatedLocation.UserID)
	assert.Equal(t, "Test Location", updatedLocation.Name)
	assert.Equal(t, "Test Label", updatedLocation.Label)
	assert.Equal(t, "test", updatedLocation.LocationType)
	assert.Equal(t, "45.5017,-73.5673", updatedLocation.MapCoordinates.Coords)
}

func TestUpdatePhysicalLocation_NotFound(t *testing.T) {
	adapter, mock, err := setupMockDB()
	require.NoError(t, err)
	defer adapter.db.Close()

	ctx := context.Background()
	userID := "user1"
	locationID := "loc1"

	// Mock the database query to return no rows
	mock.ExpectBegin()
	mock.ExpectQuery("UPDATE physical_locations").
		WithArgs(
			locationID,
			userID,
			"Test Location",
			"Test Label",
			"test",
			"45.5017,-73.5673",
			"",
		).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          userID,
		Name:            "Test Location",
		Label:           "Test Label",
		LocationType:    "test",
		MapCoordinates:  models.PhysicalMapCoordinates{Coords: "45.5017,-73.5673"},
	}

	_, err = adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.Error(t, err)
	assert.Equal(t, ErrLocationNotFound, err)
}

func TestUpdatePhysicalLocation_DatabaseError(t *testing.T) {
	adapter, mock, err := setupMockDB()
	require.NoError(t, err)
	defer adapter.db.Close()

	ctx := context.Background()
	userID := "user1"
	locationID := "loc1"

	// Mock the database query to return an error
	mock.ExpectBegin()
	mock.ExpectQuery("UPDATE physical_locations").
		WithArgs(
			locationID,
			userID,
			"Test Location",
			"Test Label",
			"test",
			"45.5017,-73.5673",
			"",
		).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	location := models.PhysicalLocation{
		ID:              locationID,
		UserID:          userID,
		Name:            "Test Location",
		Label:           "Test Label",
		LocationType:    "test",
		MapCoordinates:  models.PhysicalMapCoordinates{Coords: "45.5017,-73.5673"},
	}

	_, err = adapter.UpdatePhysicalLocation(ctx, userID, location)
	assert.Error(t, err)
}