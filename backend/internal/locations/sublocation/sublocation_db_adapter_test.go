package sublocation

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext_test"
	"github.com/lokeam/qko-beta/internal/models"
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

func TestSublocationDbAdapter(t *testing.T) {
	baseAppCtx := appcontext_test.NewTestingAppContext("test-token", nil)

	setupMockDB := func() (*SublocationDbAdapter, sqlmock.Sqlmock, error) {
		// Create mock sqldatabase
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			return nil, nil, err
		}

		// Create a sqlx wrapper around mock data
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		// Create the adapter with the mock DB
		adapter := &SublocationDbAdapter{
			db:        sqlxDB,
			logger:    baseAppCtx.Logger,
		}

		return adapter, mock, nil
	}

	// -------- GetSingleSublocation -------
	/*
		GIVEN a valid sublocation ID
		WHEN GetSingleSublocation is called
		THEN it returns the sublocation
	*/
	t.Run("GetSingleSublocation", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		sublocationID := "test-sublocation-id"
		physicalLocationID := "test-physical-location-id"
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "physical_location_id", "name",
			"location_type", "bg_color", "stored_items",
			"created_at", "updated_at",
		}).
			AddRow(
				sublocationID, userID, physicalLocationID,
				"Test Sublocation", "shelf", "blue", 20,
				now, now,
			)

		mock.ExpectQuery("SELECT (.+) FROM sublocations").
			WithArgs(sublocationID, userID).
			WillReturnRows(rows)

		gameRows := sqlmock.NewRows([]string{
			"id", "name", "summary", "cover_id", "cover_url",
			"first_release_date", "rating",
		}).AddRow(
			int64(123), "Test Game", "A test game",
			int64(456), "http://example.com/cover.jpg",
			int64(1609459200), 4.5,
		)

		mock.ExpectQuery("SELECT g.\\* FROM games g JOIN game_sub_locations gsl ON g.id = gsl.game_id WHERE gsl.sub_location_id = \\$1 AND g.user_id = \\$2").
			WithArgs(sublocationID, userID).
			WillReturnRows(gameRows)

		sublocation, err := adapter.GetSingleSublocation(context.Background(), userID, sublocationID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if sublocation.ID != sublocationID {
			t.Errorf("Expected sublocation ID %s, got %s", sublocationID, sublocation.ID)
		}
		if sublocation.PhysicalLocationID != physicalLocationID {
			t.Errorf("Expected physical location ID %s, got %s", physicalLocationID, sublocation.PhysicalLocationID)
		}
		if len(sublocation.Items) != 1 {
			t.Errorf("Expected 1 game, got %d", len(sublocation.Items))
		} else if sublocation.Items[0].ID != 123 {
			t.Errorf("Expected game ID 123, got %d", sublocation.Items[0].ID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a valid sublocation ID
		WHEN GetGamesBySublocationID is called
		THEN it returns the games
	*/
	t.Run("GetGamesBySublocationID_ValidSublocation", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		sublocationID := "test-subloc-id"

		gameRows := sqlmock.NewRows([]string{
			"id", "name", "summary", "cover_id", "cover_url",
			"first_release_date", "rating",
		}).
			AddRow(
				int64(123), "Test Game 1", "A test game",
				int64(456), "http://example.com/cover1.jpg",
				int64(1609459200), 4.5,
			).
			AddRow(
				int64(456), "Test Game 2", "Another test game",
				int64(789), "http://example.com/cover2.jpg",
				int64(1609459200), 4.7,
			)

		mock.ExpectQuery("SELECT g.\\* FROM games g JOIN game_sub_locations gsl ON g.id = gsl.game_id WHERE gsl.sub_location_id = \\$1 AND g.user_id = \\$2").
			WithArgs(sublocationID, userID).
			WillReturnRows(gameRows)

		games, err := adapter.GetGamesBySublocationID(context.Background(), userID, sublocationID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(games) != 2 {
			t.Errorf("Expected 2 games, got %d", len(games))
			return
		}
		if games[0].ID != 123 {
			t.Errorf("Expected first game ID 123, got %d", games[0].ID)
		}
		if games[1].ID != 456 {
			t.Errorf("Expected second game ID 456, got %d", games[1].ID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a non-existent sublocation ID
		WHEN GetSingleSublocation is called
		THEN it returns an error
	*/
	t.Run("GetSublocation_NotFound", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		sublocationID := "nonexistent-id"

		mock.ExpectQuery("SELECT (.+) FROM sublocations").
			WithArgs(sublocationID, userID).
			WillReturnError(sql.ErrNoRows)

		_, err = adapter.GetSingleSublocation(context.Background(), userID, sublocationID)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a valid user ID
		WHEN GetAllSublocations is called
		THEN it returns all user sublocations
	*/
	t.Run("GetAllSublocations", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		now := time.Now()

		// Set up mock expectations for sublocations query
		rows := sqlmock.NewRows([]string{
			"id", "user_id", "physical_location_id", "name",
			"location_type", "bg_color", "stored_items",
			"created_at", "updated_at",
		}).
			AddRow("subloc-1", userID, "physical-loc-1", "Sublocation 1",
			"shelf", "red", 20, now, now).
			AddRow("subloc-2", userID, "physical-loc-1", "Sublocation 2",
			"cabinet", "blue", 30, now, now)

		mock.ExpectQuery("SELECT (.+) FROM sublocations").
			WithArgs(userID).
			WillReturnRows(rows)

		// Set up mock expectations for games query for first sublocation
		gameRows1 := sqlmock.NewRows([]string{
			"id", "name", "summary", "cover_id", "cover_url",
			"first_release_date", "rating",
		}).AddRow(
			int64(123), "Test Game 1", "A test game", int64(456), "http://example.com/cover1.jpg",
			int64(1609459200), 4.5,
		)

		mock.ExpectQuery("SELECT g.\\* FROM games g JOIN game_sub_locations gsl ON g.id = gsl.game_id WHERE gsl.sub_location_id = \\$1 AND g.user_id = \\$2").
			WithArgs("subloc-1", userID).
			WillReturnRows(gameRows1)

		// Set up mock expectations for games query for second sublocation
		gameRows2 := sqlmock.NewRows([]string{
			"id", "name", "summary", "cover_id", "cover_url",
			"first_release_date", "rating",
		}).AddRow(
			int64(456), "Test Game 2", "Another test game", int64(789), "http://example.com/cover2.jpg",
			int64(1609459200), 4.7,
		)

		mock.ExpectQuery("SELECT g.\\* FROM games g JOIN game_sub_locations gsl ON g.id = gsl.game_id WHERE gsl.sub_location_id = \\$1 AND g.user_id = \\$2").
			WithArgs("subloc-2", userID).
			WillReturnRows(gameRows2)

		// Execute
		sublocations, err := adapter.GetAllSublocations(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(sublocations) != 2 {
			t.Errorf("Expected 2 sublocations, got %d", len(sublocations))
		} else {
			if sublocations[0].ID != "subloc-1" {
				t.Errorf("Expected sublocation ID subloc-1, got %s", sublocations[0].ID)
			}
			if len(sublocations[0].Items) != 1 {
				t.Errorf("Expected 1 game in first sublocation, got %d", len(sublocations[0].Items))
			}
			if sublocations[1].ID != "subloc-2" {
				t.Errorf("Expected sublocation ID subloc-2, got %s", sublocations[1].ID)
			}
			if len(sublocations[1].Items) != 1 {
				t.Errorf("Expected 1 game in second sublocation, got %d", len(sublocations[1].Items))
			}
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a valid sublocation
		WHEN UpdateSublocation is called
		THEN it updates the sublocation
	*/
	t.Run("UpdateSublocation", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		physicalLocationID := "test-physical-location-id"
		now := time.Now()

		sublocation := models.Sublocation{
			ID:                 "existing-sublocation-id",
			UserID:            userID,
			PhysicalLocationID: physicalLocationID,
			Name:              "Updated Sublocation",
			LocationType:      "shelf",
			StoredItems:       5,
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		mock.ExpectExec("UPDATE sublocations").
			WithArgs(
				sublocation.Name,
				sublocation.LocationType,
				sublocation.StoredItems,
				sqlmock.AnyArg(),
				sublocation.ID,
				sublocation.UserID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = adapter.UpdateSublocation(context.Background(), userID, sublocation)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a sublocation with mismatched user ID
		WHEN UpdateSublocation is called
		THEN it returns an error
	*/
	t.Run("UpdateSublocation_NotFound", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		sublocationID := "nonexistent-id"

		sublocation := models.Sublocation{
			ID:           sublocationID,
			UserID:       userID,
			Name:         "Updated Sublocation",
			LocationType: "cabinet",
			StoredItems:     40,
		}

		// Set up mock expectations - no rows affected
		mock.ExpectExec("UPDATE sublocations").
			WithArgs("Updated Sublocation", "cabinet", "purple", 40, sqlmock.AnyArg(), sublocationID, userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Execute
		err = adapter.UpdateSublocation(context.Background(), userID, sublocation)

		// Verify
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a valid sublocation
		WHEN CreateSublocation is called
		THEN it creates the sublocation
	*/
	t.Run("CreateSublocation", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		// Use a valid UUID format for physicalLocationID to pass uuid.Parse validation
		physicalLocationID := "123e4567-e89b-12d3-a456-426614174000"
		now := time.Now()

		sublocation := models.Sublocation{
			ID:                 "new-sublocation-id",
			UserID:            userID,
			PhysicalLocationID: physicalLocationID,
			Name:              "New Sublocation",
			LocationType:      "shelf",
			StoredItems:       0,
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		mock.ExpectQuery("INSERT INTO sublocations").
			WithArgs(
				sublocation.ID,
				sublocation.UserID,
				sublocation.PhysicalLocationID,
				sublocation.Name,
				sublocation.LocationType,
				sublocation.StoredItems,
				sqlmock.AnyArg(), // Use AnyArg() for timestamps to avoid precision issues
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "physical_location_id", "name", "location_type", "bg_color", "stored_items", "created_at", "updated_at"}).
				AddRow(sublocation.ID, sublocation.UserID, sublocation.PhysicalLocationID, sublocation.Name, sublocation.LocationType, sublocation.StoredItems, sublocation.CreatedAt, sublocation.UpdatedAt))

		result, err := adapter.CreateSublocation(context.Background(), userID, sublocation)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result.ID != sublocation.ID {
			t.Errorf("Expected sublocation ID %s, got %s", sublocation.ID, result.ID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a valid sublocation ID
		WHEN DeleteSublocation is called
		THEN it deletes the sublocation
	*/
	t.Run("DeleteSublocation", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		sublocationID := "test-subloc-id"

		// Set up mock expectations
		mock.ExpectBegin()

		// Check if sublocation exists
		rows := sqlmock.NewRows([]string{"id"}).AddRow(sublocationID)
		mock.ExpectQuery("SELECT id FROM sublocations").
			WithArgs(sublocationID, userID).
			WillReturnRows(rows)

		// Delete the sublocation
		mock.ExpectExec("DELETE FROM sublocations").
			WithArgs(sublocationID, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		// Execute
		err = adapter.DeleteSublocation(context.Background(), userID, sublocationID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a non-existent sublocation ID
		WHEN DeleteSublocation is called
		THEN it returns an error
	*/
	t.Run("RemoveSublocation_NotFound", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		sublocationID := "nonexistent-id"

		// Set up mock expectations
		mock.ExpectBegin()

		// Check if sublocation exists - return no rows
		mock.ExpectQuery("SELECT id FROM sublocations").
			WithArgs(sublocationID, userID).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectRollback()

		// Execute
		err = adapter.DeleteSublocation(context.Background(), userID, sublocationID)

		// Verify
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN valid game and sublocation IDs
		WHEN AddGameToSublocation is called
		THEN it creates the association
	*/
	t.Run("AddGameToSublocation", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		gameID := "test-game-id"
		sublocationID := "test-subloc-id"

		// Set up mock expectations
		mock.ExpectBegin()

		// Check if game exists
		gameRows := sqlmock.NewRows([]string{"id"}).AddRow(gameID)
		mock.ExpectQuery("SELECT id FROM games").
			WithArgs(gameID, userID).
			WillReturnRows(gameRows)

		// Check if sublocation exists
		sublocRows := sqlmock.NewRows([]string{"id"}).AddRow(sublocationID)
		mock.ExpectQuery("SELECT id FROM sublocations").
			WithArgs(sublocationID, userID).
			WillReturnRows(sublocRows)

		// Add the relationship
		mock.ExpectExec("INSERT INTO game_sub_locations").
			WithArgs(gameID, sublocationID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Execute
		err = adapter.AddGameToSublocation(context.Background(), userID, gameID, sublocationID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a game that doesn't belong to the user
		WHEN AddGameToSublocation is called
		THEN it returns an error
	*/
	t.Run("AddGameToSublocation_GameNotBelongToUser", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		gameID := "other-user-game-id"
		sublocationID := "test-subloc-id"

		// Set up mock expectations
		mock.ExpectBegin()

		// Check if game exists - return no rows since it belongs to another user
		mock.ExpectQuery("SELECT id FROM games").
			WithArgs(gameID, userID).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectRollback()

		// Execute
		err = adapter.AddGameToSublocation(context.Background(), userID, gameID, sublocationID)

		// Verify
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN valid game and sublocation IDs
		WHEN RemoveGameFromSublocation is called
		THEN it removes the association
	*/
	t.Run("RemoveGameFromSublocation", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		gameID := "test-game-id"
		sublocationID := "test-subloc-id"

		// Set up mock expectations
		mock.ExpectBegin()

		// Check if game exists
		gameRows := sqlmock.NewRows([]string{"id"}).AddRow(gameID)
		mock.ExpectQuery("SELECT id FROM games").
			WithArgs(gameID, userID).
			WillReturnRows(gameRows)

		// Check if sublocation exists
		sublocRows := sqlmock.NewRows([]string{"id"}).AddRow(sublocationID)
		mock.ExpectQuery("SELECT id FROM sublocations").
			WithArgs(sublocationID, userID).
			WillReturnRows(sublocRows)

		// Remove the relationship
		mock.ExpectExec("DELETE FROM game_sub_locations").
			WithArgs(gameID, sublocationID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		// Execute
		err = adapter.RemoveGameFromSublocation(context.Background(), userID, gameID, sublocationID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a non-existent association
		WHEN RemoveGameFromSublocation is called
		THEN it returns an error
	*/
	t.Run("RemoveGameFromSublocation_NonExistentAssociation", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		gameID := "test-game-id"
		sublocationID := "test-subloc-id"

		// Set up mock expectations
		mock.ExpectBegin()

		// Check if game exists
		gameRows := sqlmock.NewRows([]string{"id"}).AddRow(gameID)
		mock.ExpectQuery("SELECT id FROM games").
			WithArgs(gameID, userID).
			WillReturnRows(gameRows)

		// Check if sublocation exists
		sublocRows := sqlmock.NewRows([]string{"id"}).AddRow(sublocationID)
		mock.ExpectQuery("SELECT id FROM sublocations").
			WithArgs(sublocationID, userID).
			WillReturnRows(sublocRows)

		// Remove the relationship - no rows affected
		mock.ExpectExec("DELETE FROM game_sub_locations").
			WithArgs(gameID, sublocationID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		// Execute
		err = adapter.RemoveGameFromSublocation(context.Background(), userID, gameID, sublocationID)

		// Verify
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a sublocation with no games
		WHEN GetGamesBySublocationID is called
		THEN it returns an empty slice
	*/
	t.Run("GetGamesBySublocationID_NoGames", func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		sublocationID := "test-subloc-id"

		// Create empty rows
		gameRows := sqlmock.NewRows([]string{
			"id", "name", "summary", "cover_id", "cover_url",
			"first_release_date", "rating",
		})

		mock.ExpectQuery("SELECT g.\\* FROM games g JOIN game_sub_locations gsl ON g.id = gsl.game_id WHERE gsl.sub_location_id = \\$1 AND g.user_id = \\$2").
			WithArgs(sublocationID, userID).
			WillReturnRows(gameRows)

		// Execute
		games, err := adapter.GetGamesBySublocationID(context.Background(), userID, sublocationID)

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
}