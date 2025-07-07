package sublocation

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lokeam/qko-beta/internal/appcontext_test"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/postgres"
)

/*
	Behavior:
	- Retrieving a single sublocation for a user
	- Retrieving all sublocations for a user
	- Adding a new sublocation
	- Updating an existing sublocation
	- Removing a sublocation
	- Handling db errors
	- Ensuring user may only access their own sublocations

	Scenarios:
	- GetSingleSublocation:
		- Successfully retrieves a valid sublocation
		- Returns error when sublocation not found
		- Handles db errors
	- GetAllSublocations:
		- Successfully retrieves all sublocations for a user
		- Returns empty slice when no sublocations exist
		- Handles db errors
	- CreateSublocation:
		- Successfully adds new sublocation
		- Handles db errors
	- UpdateSublocation:
		- Successfully updates a sublocation
		- Returns error when sublocation not found
		- Handles db errors
	- DeleteSublocation:
		- Successfully removes a sublocation
		- Returns errors when sublocation not found
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

		// Create a mock PostgresClient
		mockClient := &postgres.PostgresClient{}

		// Create the adapter with the mock DB
		adapter := &SublocationDbAdapter{
			client: mockClient,
			db:     sqlxDB,
			logger: baseAppCtx.Logger,
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
			"location_type", "stored_items",
			"created_at", "updated_at",
		}).
			AddRow(
				sublocationID, userID, physicalLocationID,
				"Test Sublocation", "shelf", 20,
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
			"location_type", "stored_items",
			"created_at", "updated_at",
		}).
			AddRow("subloc-1", userID, "physical-loc-1", "Sublocation 1",
			"shelf", 20, now, now).
			AddRow("subloc-2", userID, "physical-loc-1", "Sublocation 2",
			"cabinet", 30, now, now)

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
			StoredItems:  40,
		}

		// Set up mock expectations - no rows affected
		mock.ExpectExec("UPDATE sublocations").
			WithArgs("Updated Sublocation", "cabinet", 40, sqlmock.AnyArg(), sublocationID, userID).
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
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "physical_location_id", "name", "location_type", "stored_items", "created_at", "updated_at"}).
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

		// Set up mock expectations for transaction
		mock.ExpectBegin()

		// Get user_game_ids in this sublocation
		userGameRows := sqlmock.NewRows([]string{"user_game_id"}).AddRow(123)
		mock.ExpectQuery("SELECT pgl.user_game_id FROM physical_game_locations pgl WHERE pgl.sublocation_id = \\$1").
			WithArgs(sublocationID).
			WillReturnRows(userGameRows)

		// Check if game exists in other locations
		existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(123, sublocationID).
			WillReturnRows(existsRows)

		// Get game details before deletion
		gameDetailsRows := sqlmock.NewRows([]string{"user_game_id", "game_id", "game_name", "platform_name"}).
			AddRow(123, 123, "Test Game", "PS5")
		mock.ExpectQuery("SELECT ug.id as user_game_id, g.id as game_id, g.name as game_name, p.name as platform_name FROM user_games ug JOIN games g ON ug.game_id = g.id JOIN platforms p ON ug.platform_id = p.id WHERE ug.id = \\$1").
			WithArgs(123).
			WillReturnRows(gameDetailsRows)

		// Delete the orphaned game
		mock.ExpectExec("DELETE FROM user_games WHERE id = ANY\\(\\$1\\)").
			WithArgs(123).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Delete the sublocation
		mock.ExpectExec("DELETE FROM sublocations WHERE id = ANY\\(\\$1\\) AND user_id = \\$2").
			WithArgs(pq.Array([]string{sublocationID}), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		// Execute
		_, err = adapter.DeleteSublocation(context.Background(), userID, []string{sublocationID})

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

		// Get user_game_ids in this sublocation - return no rows
		userGameRows := sqlmock.NewRows([]string{"user_game_id"})
		mock.ExpectQuery("SELECT pgl.user_game_id FROM physical_game_locations pgl WHERE pgl.sublocation_id = \\$1").
			WithArgs(sublocationID).
			WillReturnRows(userGameRows)

		// No need to expect DELETE FROM sublocations if no user_game_ids are found and code does not execute it

		mock.ExpectRollback()

		// Execute
		_, err = adapter.DeleteSublocation(context.Background(), userID, []string{sublocationID})

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
		userGameID := "test-user-game-id"

		// Set up mock expectations
		mock.ExpectBegin()

		// Check if game exists
		gameRows := sqlmock.NewRows([]string{"id"}).AddRow(userGameID)
		mock.ExpectQuery("SELECT id FROM user_games").
			WithArgs(userGameID, userID).
			WillReturnRows(gameRows)

		// Remove the relationship
		mock.ExpectExec("DELETE FROM physical_game_locations").
			WithArgs(userGameID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		// Execute
		err = adapter.RemoveGameFromSublocation(context.Background(), userID, userGameID)

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
		userGameID := "test-user-game-id"

		// Set up mock expectations
		mock.ExpectBegin()

		// Check if game exists
		gameRows := sqlmock.NewRows([]string{"id"}).AddRow(userGameID)
		mock.ExpectQuery("SELECT id FROM user_games").
			WithArgs(userGameID, userID).
			WillReturnRows(gameRows)

		// Remove the relationship - no rows affected
		mock.ExpectExec("DELETE FROM physical_game_locations").
			WithArgs(userGameID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		// Execute
		err = adapter.RemoveGameFromSublocation(context.Background(), userID, userGameID)

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