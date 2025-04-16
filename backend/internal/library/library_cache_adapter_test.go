package library

import (
	"context"
	"errors"
	"testing"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/stretchr/testify/mock"
)

/*
	Behavior:
	- Getting cached library items for a user
	- Setting cached library items for a user
	- Getting a cached game for a user
	- Setting a cached game for a user
	- Invalidating cache for a user
	- Invalidating cache for a specific game

	Scenarios:
	- GetCachedLibraryItems with cache hit
	- GetCachedLibraryItems with cache miss
	- GetCachedLibraryItems with cache error
	- SetCachedLibraryItems success
	- SetCachedLibraryItems error
	- GetCachedGame with cache hit
	- GetCachedGame with cache miss
	- GetCachedGame with cache error
	- SetCachedGame success
	- SetCachedGame error
	- InvalidateUserCache success
	- InvalidateUserCache error
	- InvalidateGameCache success
	- InvalidateGameCache error
*/

type MockCacheWrapper struct {
	mock.Mock
	games []models.Game
}

// Mock implementations of CacheWrapper methods
func (m *MockCacheWrapper) GetCachedResults(ctx context.Context, key string, result any) (bool, error) {
	args := m.Called(ctx, key, result)

	// If there is a result and the third argument is a slice of games
	if args.Get(0).(bool) && result != nil {
		switch value := result.(type) {
		case *[]models.Game:
			// Copy mock games to result pointer
			*value = m.games
		case *models.Game:
			// Copy mock game to result pointer
			*value = m.games[0]
		}
	}

	return args.Get(0).(bool), args.Error(1)
}

func (m *MockCacheWrapper) SetCachedResults(ctx context.Context, key string, result any) error {
	args := m.Called(ctx, key, result)
	return args.Error(0)
}

func (m *MockCacheWrapper) DeleteCacheKey(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func TestLibraryCacheAdapter(t *testing.T) {
	// Setup test data
	testUserID := "test-user-id"
	testGameID := int64(123)
	testGame := models.Game{ID: testGameID, Name: "Test Game"}
	testGames := []models.Game{testGame}
	testError := errors.New("cache error")

	// Helper func to create a new adapter with a mock cache wrapper
	createAdapter := func(mockCache *MockCacheWrapper) *LibraryCacheAdapter {
		adapter, _ := NewLibraryCacheAdapter(mockCache)
		return adapter.(*LibraryCacheAdapter)
	}

	// ------ GetCachedLibraryItems() ------
	/*
		GIVEN a cache wrapper that returns cached games
		WHEN GetCachedLibraryItems is called
		THEN it should return the cached games without error
	*/
	t.Run("GetCachedLibraryItems with cache hit", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.games = testGames
		mockCache.On("GetCachedResults", mock.Anything, "library:test-user-id", mock.Anything).Return(true, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		games, err := adapter.GetCachedLibraryItems(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(games) != 1 {
			t.Errorf("Expected 1 game, got %d", len(games))
		}
		if games[0].ID != testGameID {
			t.Errorf("Expected game ID %d, got %d", testGameID, games[0].ID)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns a cache miss
		WHEN GetCachedLibraryItems is called
		THEN it should return nil games without error
	*/
	t.Run(`GetCachedLibraryItems with cache miss`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "library:test-user-id", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		games, err := adapter.GetCachedLibraryItems(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if games != nil {
			t.Errorf("Expected nil games, got %v", games)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error
		WHEN GetCachedLibraryItems is called
		THEN it should return the error
	*/
	t.Run(`GetCachedLibraryItems with cache error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "library:test-user-id", mock.Anything).Return(false, testError)

		adapter := createAdapter(mockCache)

		// WHEN
		games, err := adapter.GetCachedLibraryItems(context.Background(), testUserID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		if games != nil {
			t.Errorf("Expected nil games, got %v", games)
		}
		mockCache.AssertExpectations(t)
	})

	// ------ SetCachedLibraryItems() ------
	/*
		GIVEN a cache wrapper that successfully sets cached results
		WHEN SetCachedLibraryItems is called
		THEN it should not return an error
	*/
	t.Run(`SetCachedLibraryItems success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id", testGames).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedLibraryItems(context.Background(), testUserID, testGames)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)

	})

	/*
		GIVEN a cache wrapper that returns an error when setting cached results
		WHEN SetCachedLibraryItems is called
		THEN it should return the error
	*/
	t.Run(`SetCachedLibraryItems error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id", testGames).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedLibraryItems(context.Background(),testUserID, testGames)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v but instead got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})

	// ------ GetCachedGame() ------\
	/*
		GIVEN a cache wrapper that returns a cached game
		WHEN GetCachedGame is called
		THEN it should return the cached game without error
	*/
	t.Run(`GetCachedGame with cache hit`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.games = testGames
		mockCache.On("GetCachedResults", mock.Anything, "library:test-user-id:game:123", mock.Anything).Return(true, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		game, found, err := adapter.GetCachedGame(context.Background(), testUserID, testGameID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !found {
			t.Errorf("Expected game to be found in cache")
		}
		if game == nil {
			t.Errorf("Expected non-nil game")
		} else if game.ID != testGameID {
			t.Errorf("Expected game ID %d, got %d", testGameID, game.ID)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns a cache miss
		WHEN GetCachedGame is called
		THEN it should indicate the game was not found without error
	*/
	t.Run(`GetCachedGame with cache miss`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "library:test-user-id:game:123", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		game, found, err := adapter.GetCachedGame(context.Background(), testUserID, testGameID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if found {
			t.Errorf("Expected found to be false")
		}
		if game != nil {
			t.Errorf("Expected nil game, got %v", game)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error
		WHEN GetCachedGame is called
		THEN it should return an error
	*/
	t.Run(`GetCachedGame with cache error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "library:test-user-id:game:123", mock.Anything).Return(false, testError)

		adapter := createAdapter(mockCache)

		// WHEN
		game, found, err := adapter.GetCachedGame(context.Background(), testUserID, testGameID)

		// THEN
		if err == nil || err == testError {
			t.Errorf("Expected ErrDatabaseConnection error but instead got %v", err)
		}
		if found {
			t.Errorf("Expected game not to be found")
		}
		if game != nil {
			t.Errorf("Expected nil game but instead got %v", game)
		}
		mockCache.AssertExpectations(t)
	})

	// ------ SetCachedGame() ------
	/*
		GIVEN a cache wrapper that successfully sets a cached game
		WHEN SetCachedGame is called
		THEN it should not return an error
	*/
	t.Run(`SetCachedGame success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id:game:123", testGame).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedGame(context.Background(), testUserID, testGame)

		// THEN
		if err != nil {
			t.Errorf("Expected no error but instead got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error when setting a cached game
		WHEN SetCachedGame is called
		THEN it should return the error
	*/
	t.Run(`SetCachedGame error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id:game:123", testGame).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedGame(context.Background(), testUserID, testGame)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})

	// ------ InvalidateUserCache() ------
	/*
		GIVEN a cache wrapper that successfully invalidates user cache
		WHEN InvalidateUserCache is called
		THEN it should not return an error
	*/
	t.Run(`InvalidateUserCache success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id:game", nil).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateUserCache(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error but instead got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error when invalidating user cache
		WHEN InvalidateUserCache is called
		THEN it should return the error
	*/
	t.Run(`InvalidateUserCache error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id:game", nil).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateUserCache(context.Background(), testUserID)

		// THEN
		// THEN it should return the error
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})

	// ------ InvalidateGameCache() ------
	/*
		GIVEN a cache wrapper that successfully invalidates game cache
		WHEN InvalidateGameCache is called
		THEN it should not return an error
	*/
	t.Run(`InvalidateGameCache success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id:game:123", nil).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateGameCache(context.Background(), testUserID, testGameID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error when invalidating game cache
		WHEN InvalidateGameCache is called
		THEN it should return the error
	*/
	t.Run(`InvalidateGameCache error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id:game:123", nil).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateGameCache(context.Background(), testUserID, testGameID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})
}
