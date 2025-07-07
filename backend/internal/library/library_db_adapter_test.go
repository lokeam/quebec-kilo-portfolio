package library

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
)

/*
	Behavior:
	- Retrieving games from a user's library
	- Adding games to a user's library
	- Removing games from a user's library
	- Checking if a game exists in a user's library
	- Updating game information in the database
	- Handling database transactions

	Scenarios:
	- GetSingleLibraryGame returns false when game not found
	- GetSingleLibraryGame handles database errors
	- GetLibraryItems handles database errors
	- CreateLibraryGame successfully adds a new game
	- CreateLibraryGame handles existing games
	- DeleteLibraryGame successfully removes a game
	- IsGameInLibrary correctly identifies if a game is in library
*/

// MockGameScanner implements the GameScanner interface for testing
type MockGameScanner struct {
	ScanGameFunc  func(row *sqlx.Row) (models.Game, error)
	ScanGamesFunc func(rows *sqlx.Rows) ([]models.Game, error)
}

func (mgs *MockGameScanner) ScanGame(row *sqlx.Row) (models.Game, error) {
	return mgs.ScanGameFunc(row)
}

func (mgs *MockGameScanner) ScanGames(rows *sqlx.Rows) ([]models.Game, error) {
	return mgs.ScanGamesFunc(rows)
}

func TestLibraryDbAdapter(t *testing.T) {
	// Setup test data
	userID := "test-user-id"
	gameID := int64(123)
	expectedGame := models.Game{ID: gameID, Name: "Test Game"}

	// Helper function to setup mock database
	setupMockDB := func() (*LibraryDbAdapter, sqlmock.Sqlmock, error) {
		db, mock, err := sqlmock.New()
		if err != nil {
			return nil, nil, err
		}

		sqlxDB := sqlx.NewDb(db, "postgres")
		appContext := &appcontext.AppContext{
			Logger: testutils.NewTestLogger(),
		}

		adapter := &LibraryDbAdapter{
			db:      sqlxDB,
			logger:  appContext.Logger,
			scanner: &MockGameScanner{},
		}

		return adapter, mock, nil
	}

	/*
		GIVEN a request to get a specific game from a user's library
		WHEN the game exists in the database
		THEN the adapter returns the game and true
	*/
	t.Run("GetSingleLibraryGame returns game when found", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Error setting up mock DB: %v", err)
		}
		defer adapter.db.Close()

		// Set up mock expectations
		mock.ExpectQuery("SELECT (.+) FROM user_library ul").
			WithArgs(userID, gameID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Mock the game data query
		mock.ExpectQuery("SELECT (.+) FROM user_games ug").
			WithArgs(userID, gameID).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "name", "cover_url", "game_type_display_text",
				"game_type_normalized_text", "is_favorite", "created_at",
			}).AddRow(
				expectedGame.ID, expectedGame.Name, "", "Physical", "physical", true, "2024-01-01",
			))

		// Mock the locations query
		mock.ExpectQuery("SELECT (.+) FROM user_games ug").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{
				"game_id", "platform_id", "platform_name", "category", "created_at",
				"parent_location_id", "parent_location_name", "parent_location_type",
				"parent_location_bg_color", "sublocation_id", "sublocation_name", "sublocation_type",
			}))

		// Execute the fn
		game, err := adapter.GetSingleLibraryGame(context.Background(), userID, gameID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if game.ID != expectedGame.ID {
			t.Errorf("Expected game ID %d, got %d", expectedGame.ID, game.ID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get a specific game from a user's library
		WHEN the game does not exist in the database
		THEN the adapter returns ErrGameNotFound
	*/
	t.Run("GetSingleLibraryGame returns error when game not found", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Error setting up mock DB: %v", err)
		}
		defer adapter.db.Close()

		// Set up mock expectations
		mock.ExpectQuery("SELECT (.+) FROM user_library ul").
			WithArgs(userID, gameID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Execute
		_, err = adapter.GetSingleLibraryGame(context.Background(), userID, gameID)

		// Verify
		if err == nil {
			t.Errorf("Expected ErrGameNotFound, got nil")
		}
		if err != ErrGameNotFound {
			t.Errorf("Expected ErrGameNotFound, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	/*
		GIVEN a request to add a game to a user's library
		WHEN the game doesn't already exist in the database
		THEN the adapter adds the game to the database and the user's library
	*/
	t.Run("CreateLibraryGame - Successfully adds a new game", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Error setting up mock DB: %v", err)
		}
		defer adapter.db.Close()

		gameToSave := models.GameToSave{
			GameID:               gameID,
			GameName:             "Test Game",
			GameCoverURL:         "https://example.com/cover.jpg",
			GameFirstReleaseDate: 1640995200,
			GameRating:           8.5,
			GameType: models.GameToSaveIGDBType{
				DisplayText:    "Physical",
				NormalizedText: "physical",
			},
			PlatformLocations: []models.GameToSaveLocation{
				{
					PlatformID:   1,
					PlatformName: "PlayStation 5",
					Type:         "physical",
					Location: models.GameToSaveLocationDetails{
						SublocationID: "shelf-1",
					},
				},
			},
		}

		// Set up mock expectations for transaction
		mock.ExpectBegin()

		// Insert game into games table
		mock.ExpectExec("INSERT INTO games").
			WithArgs(gameToSave.GameID, gameToSave.GameName, gameToSave.GameCoverURL,
				gameToSave.GameFirstReleaseDate, gameToSave.GameRating).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Insert platform
		mock.ExpectExec("INSERT INTO platforms").
			WithArgs(gameToSave.PlatformLocations[0].PlatformID,
				gameToSave.PlatformLocations[0].PlatformName, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Insert user game
		mock.ExpectQuery("INSERT INTO user_games").
			WithArgs(userID, gameToSave.GameID, gameToSave.PlatformLocations[0].PlatformID,
				gameToSave.PlatformLocations[0].Type).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Commit transaction
		mock.ExpectCommit()

		// Execute
		err = adapter.CreateLibraryGame(context.Background(), userID, gameToSave)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	/*
		GIVEN a request to add a game to a user's library
		WHEN the game already exists in the database
		THEN the adapter only adds the game to the user's library
	*/
	t.Run("CreateLibraryGame - Handles existing games", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Error setting up mock DB: %v", err)
		}
		defer adapter.db.Close()

		gameToSave := models.GameToSave{
			GameID:               gameID,
			GameName:             "Test Game",
			GameCoverURL:         "https://example.com/cover.jpg",
			GameFirstReleaseDate: 1640995200,
			GameRating:           8.5,
			GameType: models.GameToSaveIGDBType{
				DisplayText:    "Physical",
				NormalizedText: "physical",
			},
			PlatformLocations: []models.GameToSaveLocation{
				{
					PlatformID:   1,
					PlatformName: "PlayStation 5",
					Type:         "physical",
					Location: models.GameToSaveLocationDetails{
						SublocationID: "shelf-1",
					},
				},
			},
		}

		// Set up mock expectations for transaction
		mock.ExpectBegin()

		// Insert game into games table (will be ignored due to ON CONFLICT DO NOTHING)
		mock.ExpectExec("INSERT INTO games").
			WithArgs(gameToSave.GameID, gameToSave.GameName, gameToSave.GameCoverURL,
				gameToSave.GameFirstReleaseDate, gameToSave.GameRating).
			WillReturnResult(sqlmock.NewResult(1, 0))

		// Insert platform (will be ignored due to ON CONFLICT DO NOTHING)
		mock.ExpectExec("INSERT INTO platforms").
			WithArgs(gameToSave.PlatformLocations[0].PlatformID,
				gameToSave.PlatformLocations[0].PlatformName, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 0))

		// Insert user game (will fail with unique constraint, then select existing)
		mock.ExpectQuery("INSERT INTO user_games").
			WithArgs(userID, gameToSave.GameID, gameToSave.PlatformLocations[0].PlatformID,
				gameToSave.PlatformLocations[0].Type).
			WillReturnError(sql.ErrNoRows)

		// Select existing user game
		mock.ExpectQuery("SELECT id FROM user_games").
			WithArgs(userID, gameToSave.GameID, gameToSave.PlatformLocations[0].PlatformID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Commit transaction
		mock.ExpectCommit()

		// Execute
		err = adapter.CreateLibraryGame(context.Background(), userID, gameToSave)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	/*
		GIVEN a request to remove a game from a user's library
		WHEN the database operation is successful
		THEN the adapter removes the game from the user's library
	*/
	t.Run("DeleteLibraryGame - Successfully removes a game", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Error setting up mock DB: %v", err)
			defer adapter.db.Close()
		}

		// Set up mock expectations for transaction
		mock.ExpectBegin()

		// Delete from user library
		mock.ExpectExec("DELETE FROM user_library").
			WithArgs(userID, gameID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit transaction
		mock.ExpectCommit()

		// Execute
		err = adapter.DeleteLibraryGame(context.Background(), userID, gameID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	/*
		GIVEN a request to update a game in a user's library
		WHEN the database operation is successful
		THEN the adapter updates the game in the user's library
	*/
	t.Run("UpdateLibraryGame - Successfully updates a game", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Error setting up mock DB: %v", err)
		}
		defer adapter.db.Close()

		gameToSave := models.GameToSave{
			GameID:               gameID,
			GameName:             "Updated Test Game",
			GameCoverURL:         "https://example.com/updated-cover.jpg",
			GameFirstReleaseDate: 1640995200,
			GameRating:           9.0,
			GameType: models.GameToSaveIGDBType{
				DisplayText:    "Digital",
				NormalizedText: "digital",
			},
			PlatformLocations: []models.GameToSaveLocation{
				{
					PlatformID:   2,
					PlatformName: "Steam",
					Type:         "digital",
					Location: models.GameToSaveLocationDetails{
						DigitalLocationID: "steam-lib",
					},
				},
			},
		}

		// Set up mock expectations
		// Update game
		mock.ExpectExec("UPDATE games SET").
			WithArgs(gameToSave.GameID, gameToSave.GameName, gameToSave.GameCoverURL,
				gameToSave.GameType.DisplayText, gameToSave.GameType.NormalizedText).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Delete existing locations
		mock.ExpectExec("DELETE FROM user_games").
			WithArgs(gameToSave.GameID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Insert new location
		mock.ExpectExec("INSERT INTO user_games").
			WithArgs(gameToSave.GameID, gameToSave.PlatformLocations[0].PlatformID,
				gameToSave.PlatformLocations[0].PlatformName, gameToSave.PlatformLocations[0].Type,
				gameToSave.PlatformLocations[0].Location.SublocationID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Execute
		err = adapter.UpdateLibraryGame(context.Background(), gameToSave)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}