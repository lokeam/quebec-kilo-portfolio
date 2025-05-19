package library

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lokeam/qko-beta/internal/appcontext_test"
	"github.com/lokeam/qko-beta/internal/models"
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
	- GetUserGame successfully retrieves a game
	- GetUserGame returns false when game not found
	- GetUserGame handles database errors
	- GetAllLibraryGames successfully retrieves all games
	- GetLibraryItems handles database errors
	- CreateLibraryGame successfully adds a new game
	- CreateLibraryGame handles existing games
	- RemoveGameFromLibrary successfully removes a game
	- IsGameInLibrary correctly identifies if a game is in library
*/

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
	// Set up base app context for testing
	baseAppCtx := appcontext_test.NewTestingAppContext("test-token", nil)

	// Helper fn - Create mock DB + adapter
	setupMockDB := func() (*LibraryDbAdapter, sqlmock.Sqlmock, error) {
		// Create a sqlmock database
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			return nil, nil, err
		}

		// Create a sqlx wrapper around the mock DB
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		// Create the adapter with the mock DB
		adapter := &LibraryDbAdapter{
			db:      sqlxDB,
			logger:  baseAppCtx.Logger,
			scanner: &MockGameScanner{},
		}

		return adapter, mock, nil
	}


	/*
		GIVEN a request to get a specific game from a user's library
		WHEN the game exists in the database
		THEN the adapter returns the game and true
	*/
	t.Run(`GetUserGame - Successfully retrieves a game`, func(t *testing.T) {
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		gameID := int64(123)
		expectedGame := models.Game{
			ID:   gameID,
			Name: "Test Game",
		}

		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "name", "summary", "cover_url", "first_release_date", "rating", "platform_names", "genre_names", "theme_names"})
		mock.ExpectQuery("SELECT (.+) FROM user_library ul").
			WithArgs(userID, gameID).
			WillReturnRows(rows)

		// Set up mock scanner
		adapter.scanner = &MockGameScanner{
			ScanGameFunc: func(row *sqlx.Row) (models.Game, error) {
				return expectedGame, nil
			},
		}

		// Execute the fn
		game, exists, err := adapter.GetUserGame(context.Background(), userID, gameID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !exists {
			t.Errorf("Expected game to exist, but got false")
		}
		if game.ID != expectedGame.ID || game.Name != expectedGame.Name {
			t.Errorf("Expected game to be %+v, but got %+v", expectedGame, game)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})


	/*
		GIVEN a request to get a specific game from a user's library
		WHEN the game does not exist in the database
		THEN the adapter returns false and no error
	*/
	t.Run("GetUserGame returns false when game not found", func(t *testing.T) {
    // Setup
    adapter, mock, err := setupMockDB()
    if err != nil {
        t.Fatalf("Error setting up mock DB: %v", err)
    }
    defer adapter.db.Close()

    userID := "test-user-id"
    gameID := int64(123)

    // Set up mock expectations
    mock.ExpectQuery("SELECT (.+) FROM user_library ul").
        WithArgs(userID, gameID).
        WillReturnRows(sqlmock.NewRows([]string{}))

    // Set up mock scanner
    adapter.scanner = &MockGameScanner{
        ScanGameFunc: func(row *sqlx.Row) (models.Game, error) {
            return models.Game{}, pgx.ErrNoRows
        },
    }

    // Execute
    _, exists, err := adapter.GetUserGame(context.Background(), userID, gameID)

    // Verify
    if err == nil {
        t.Errorf("Expected ErrGameNotFound, got nil")
    }
    if err != ErrGameNotFound {
        t.Errorf("Expected ErrGameNotFound, got %v", err)
    }
    if exists {
        t.Errorf("Expected game to not exist, got true")
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
})


	/*
		GIVEN a request to get all games from a user's library
		WHEN the database query is successful
		THEN the adapter returns all games
	*/
	t.Run(`GetAllLibraryGames - Successfully retrieves all games`, func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		expectedGames := []models.Game{
			{ ID: 123, Name: "Game 1" },
			{ ID: 456, Name: "Game 2" },
		}

		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id", "name", "summary", "cover_url", "first_release_date", "rating", "platform_names", "genre_names", "theme_names"})
		mock.ExpectQuery("SELECT (.+) FROM user_library ul").
			WithArgs(userID).
			WillReturnRows(rows)

		// Set up mock scanner
		adapter.scanner = &MockGameScanner{
			ScanGamesFunc: func(rows *sqlx.Rows) ([]models.Game, error) {
				return expectedGames, nil
			},
		}

		// Execute
		games, err := adapter.GetAllLibraryGames(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(games) != len(expectedGames) {
			t.Errorf("Expected %d games, got %d", len(expectedGames), len(games))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get all games from a user's library
		WHEN the database returns an error
		THEN the adapter returns the error
	*/
	t.Run("GetAllLibraryGames handles database errors", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Error setting up mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		dbError := ErrDatabaseConnection

		// Setup mock expectations
		mock.ExpectQuery("SELECT (.+) FROM user_library ul").
			WithArgs(userID).
			WillReturnError(dbError)

		// Execute
		_, err = adapter.GetAllLibraryGames(context.Background(), userID)

		// Verify
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, dbError) && !strings.Contains(err.Error(), dbError.Error()) {
			t.Errorf("Expected error to contain %v, got %v", dbError, err)
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

		userID := "test-user-id"
		gameID := int64(123)

		// Set up mock expectations for transaction
		mock.ExpectBegin()

		// Check if game exists
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(gameID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Insert game
		mock.ExpectExec("INSERT INTO games").
			WithArgs(gameID, "", "", "", "", "", "", "", "").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Add to user library
		mock.ExpectExec("INSERT INTO user_library").
			WithArgs(userID, gameID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit transaction
		mock.ExpectCommit()

		// Execute
		err = adapter.CreateLibraryGame(context.Background(), userID, gameID)

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

		userID := "test-user-id"
		gameID := int64(123)

		// Set up mock expectations for transaction
		mock.ExpectBegin()

		// Check if game exists and return true
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(gameID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Add to user library
		mock.ExpectExec("INSERT INTO user_library").
			WithArgs(userID, gameID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit transaction
		mock.ExpectCommit()

		// Execute
		err = adapter.CreateLibraryGame(context.Background(), userID, gameID)

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
	t.Run("RemoveGameFromLibrary - Successfully removes a game", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Error setting up mock DB: %v", err)
			defer adapter.db.Close()
		}

		userID := "test-user-id"
		gameID := int64(123)

		// Set up mock expectations for transaction
		mock.ExpectBegin()

		// Delete from user library
		mock.ExpectExec("DELETE FROM user_library").
			WithArgs(userID, gameID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit transaction
		mock.ExpectCommit()

		// Execute
		err = adapter.RemoveGameFromLibrary(context.Background(), userID, gameID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})


	/*
		GIVEN a request to check if a game is in a user's library
		WHEN the game is in the library
		THEN the adapter returns true
	*/
	t.Run("IsGameInLibrary correctly identifies if a game is in library", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Error setting up mock DB: %v", err)
		}
		defer adapter.db.Close()

		userID := "test-user-id"
		gameID := int64(123)

		// Set up mock expectations
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(userID, gameID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Execute
		exists, err := adapter.IsGameInLibrary(context.Background(), userID, gameID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !exists {
			t.Errorf("Expected game to exist in library, got false")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})



	/*
		GIVEN a request to update a game in the database
		WHEN the database operation is successful
		THEN the adapter updates the game information
	*/
	t.Run("UpdateLibraryGame successfully updates a game", func(t *testing.T) {
		// Setup
		adapter, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Error setting up mock DB: %v", err)
		}
		defer adapter.db.Close()

		game := models.Game{
			ID:              123,
			Name:            "Updated Game",
			Summary:         "Updated Summary",
			CoverURL:        "http://example.com/cover.jpg",
			FirstReleaseDate: time.Now().Unix(),
			Rating:          85.5,
			PlatformNames:   []string{"Platform1", "Platform2"},
			GenreNames:      []string{"Genre1", "Genre2"},
			ThemeNames:      []string{"Theme1", "Theme2"},
		}

		// Set up mock expectations for transaction
		mock.ExpectBegin()

		// Update game
		mock.ExpectExec("UPDATE games").
			WithArgs(
				game.Name,
				game.Summary,
				game.CoverURL,
				game.FirstReleaseDate,
				game.Rating,
				pq.Array(game.PlatformNames),
				pq.Array(game.GenreNames),
				pq.Array(game.ThemeNames),
				game.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit transaction
		mock.ExpectCommit()

		// Execute
		err = adapter.UpdateLibraryGame(context.Background(), game)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}