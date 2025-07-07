package library

import (
	"context"
	"errors"
	"testing"

	"github.com/lokeam/qko-beta/internal/types"
	"github.com/stretchr/testify/mock"
)

/*
	Behavior:
	- Getting cached library items for a user
	- Setting cached library items for a user
	- Getting a cached game for a user
	- Setting a cached game for a user
	- Getting cached BFF library items for a user
	- Setting cached BFF library items for a user
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
	- GetCachedLibraryItemsBFF with cache hit
	- GetCachedLibraryItemsBFF with cache miss
	- GetCachedLibraryItemsBFF with cache error
	- SetCachedLibraryItemsBFF success
	- SetCachedLibraryItemsBFF error
	- InvalidateUserCache success
	- InvalidateUserCache error
	- InvalidateGameCache success
	- InvalidateGameCache error
*/

type MockCacheWrapper struct {
	mock.Mock
}

// Mock implementations of CacheWrapper methods
func (m *MockCacheWrapper) GetCachedResults(ctx context.Context, key string, result any) (bool, error) {
	args := m.Called(ctx, key, result)
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

func (m *MockCacheWrapper) InvalidateCache(ctx context.Context, cacheKey string) error {
	args := m.Called(ctx, cacheKey)
	return args.Error(0)
}

func TestLibraryCacheAdapter(t *testing.T) {
	// Setup test data
	testUserID := "test-user-id"
	testGameID := int64(123)
	testError := errors.New("cache error")

	// Test data for library items
	testGames := []types.LibraryGameDBResult{
		{
			ID:                    testGameID,
			Name:                  "Test Game",
			CoverURL:              "https://example.com/cover.jpg",
			FirstReleaseDate:      1640995200,
			Rating:                8.5,
			ThemeNames:            []string{"Action", "Adventure"},
			Favorite:              true,
			IsInWishlist:          false,
			GameTypeDisplay:       "Physical",
			GameTypeNormalized:    "physical",
			PlatformID:            1,
			PlatformName:          "PlayStation 5",
		},
	}

	testPhysicalLocations := []types.LibraryGamePhysicalLocationDBResponse{
		{
			ID:               1,
			PlatformID:       1,
			PlatformName:     "PlayStation 5",
			LocationID:       "loc-1",
			LocationName:     "Living Room",
			LocationType:     "physical",
			SublocationID:    "sub-1",
			SublocationName:  "Shelf A",
			SublocationType:  "shelf",
			SublocationBgColor: "#FF0000",
		},
	}

	testDigitalLocations := []types.LibraryGameDigitalLocationDBResponse{
		{
			ID:           1,
			PlatformID:   2,
			PlatformName: "Steam",
			LocationID:   "dig-loc-1",
			LocationName: "Steam Library",
			IsActive:     true,
		},
	}

	testGame := types.LibraryGameItemBFFResponseFINAL{
		ID:                    testGameID,
		Name:                  "Test Game",
		CoverURL:              "https://example.com/cover.jpg",
		GameTypeDisplayText:   "Physical",
		GameTypeNormalizedText: "physical",
		IsFavorite:            true,
		GamesByPlatformAndLocation: []types.LibraryGamesByPlatformAndLocationItemFINAL{
			{
				ID:                 1,
				PlatformID:         1,
				PlatformName:       "PlayStation 5",
				IsPC:               false,
				IsMobile:           false,
				DateAdded:          1640995200,
				ParentLocationID:   "loc-1",
				ParentLocationName: "Living Room",
				ParentLocationType: "physical",
				ParentLocationBgColor: "#FF0000",
				SublocationID:      "sub-1",
				SublocationName:    "Shelf A",
				SublocationType:    "shelf",
			},
		},
	}

	testBFFResponse := types.LibraryBFFResponseFINAL{
		LibraryItems:  []types.LibraryGameItemBFFResponseFINAL{testGame},
		RecentlyAdded: []types.LibraryGameItemBFFResponseFINAL{testGame},
	}

	// Helper func to create a new adapter with a mock cache wrapper
	createAdapter := func(mockCache *MockCacheWrapper) *LibraryCacheAdapter {
		adapter, _ := NewLibraryCacheAdapter(mockCache)
		return adapter.(*LibraryCacheAdapter)
	}

	// ------ GetCachedLibraryItems() ------
	/*
		GIVEN a cache wrapper that returns cached library items
		WHEN GetCachedLibraryItems is called
		THEN it should return the cached items without error
	*/
	t.Run("GetCachedLibraryItems with cache hit", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "library:test-user-id", mock.Anything).Return(true, nil).Run(func(args mock.Arguments) {
			// Set the result in the passed pointer
			result := args.Get(2).(*cachedLibraryItems)
			*result = cachedLibraryItems{
				Games:             testGames,
				PhysicalLocations: testPhysicalLocations,
				DigitalLocations:  testDigitalLocations,
			}
		})

		adapter := createAdapter(mockCache)

		// WHEN
		games, physicalLocations, digitalLocations, err := adapter.GetCachedLibraryItems(context.Background(), testUserID)

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
		if len(physicalLocations) != 1 {
			t.Errorf("Expected 1 physical location, got %d", len(physicalLocations))
		}
		if len(digitalLocations) != 1 {
			t.Errorf("Expected 1 digital location, got %d", len(digitalLocations))
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns a cache miss
		WHEN GetCachedLibraryItems is called
		THEN it should return nil items without error
	*/
	t.Run(`GetCachedLibraryItems with cache miss`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "library:test-user-id", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		games, physicalLocations, digitalLocations, err := adapter.GetCachedLibraryItems(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if games != nil {
			t.Errorf("Expected nil games, got %v", games)
		}
		if physicalLocations != nil {
			t.Errorf("Expected nil physical locations, got %v", physicalLocations)
		}
		if digitalLocations != nil {
			t.Errorf("Expected nil digital locations, got %v", digitalLocations)
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
		games, physicalLocations, digitalLocations, err := adapter.GetCachedLibraryItems(context.Background(), testUserID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		if games != nil {
			t.Errorf("Expected nil games, got %v", games)
		}
		if physicalLocations != nil {
			t.Errorf("Expected nil physical locations, got %v", physicalLocations)
		}
		if digitalLocations != nil {
			t.Errorf("Expected nil digital locations, got %v", digitalLocations)
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
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id", mock.Anything).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedLibraryItems(context.Background(), testUserID, testGames, testPhysicalLocations, testDigitalLocations)

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
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id", mock.Anything).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedLibraryItems(context.Background(), testUserID, testGames, testPhysicalLocations, testDigitalLocations)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v but instead got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})

	// ------ GetCachedGame() ------
	/*
		GIVEN a cache wrapper that returns a cached game
		WHEN GetCachedGame is called
		THEN it should return the cached game without error
	*/
	t.Run(`GetCachedGame with cache hit`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "library:test-user-id:game:123", mock.Anything).Return(true, nil).Run(func(args mock.Arguments) {
			// Set the result in the passed pointer
			result := args.Get(2).(*cachedGame)
			*result = cachedGame{Game: testGame}
		})

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
		if game.ID != testGameID {
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
		if game.ID != 0 {
			t.Errorf("Expected empty game struct, got %v", game)
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
		if err != testError {
			t.Errorf("Expected error %v but instead got %v", testError, err)
		}
		if found {
			t.Errorf("Expected game not to be found")
		}
		if game.ID != 0 {
			t.Errorf("Expected empty game struct but instead got %v", game)
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
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id:game:123", mock.Anything).Return(nil)

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
		mockCache.On("SetCachedResults", mock.Anything, "library:test-user-id:game:123", mock.Anything).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedGame(context.Background(), testUserID, testGame)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})

	// ------ GetCachedLibraryItemsBFF() ------
	/*
		GIVEN a cache wrapper that returns cached BFF library items
		WHEN GetCachedLibraryItemsBFF is called
		THEN it should return the cached BFF response without error
	*/
	t.Run(`GetCachedLibraryItemsBFF with cache hit`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "library:bff:test-user-id", mock.Anything).Return(true, nil).Run(func(args mock.Arguments) {
			// Set the result in the passed pointer
			result := args.Get(2).(*types.LibraryBFFResponseFINAL)
			*result = testBFFResponse
		})

		adapter := createAdapter(mockCache)

		// WHEN
		response, err := adapter.GetCachedLibraryItemsBFF(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(response.LibraryItems) != 1 {
			t.Errorf("Expected 1 library item, got %d", len(response.LibraryItems))
		}
		if response.LibraryItems[0].ID != testGameID {
			t.Errorf("Expected game ID %d, got %d", testGameID, response.LibraryItems[0].ID)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns a cache miss
		WHEN GetCachedLibraryItemsBFF is called
		THEN it should return empty response without error
	*/
	t.Run(`GetCachedLibraryItemsBFF with cache miss`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "library:bff:test-user-id", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		response, err := adapter.GetCachedLibraryItemsBFF(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(response.LibraryItems) != 0 {
			t.Errorf("Expected empty library items, got %d", len(response.LibraryItems))
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error
		WHEN GetCachedLibraryItemsBFF is called
		THEN it should return the error
	*/
	t.Run(`GetCachedLibraryItemsBFF with cache error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "library:bff:test-user-id", mock.Anything).Return(false, testError)

		adapter := createAdapter(mockCache)

		// WHEN
		response, err := adapter.GetCachedLibraryItemsBFF(context.Background(), testUserID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		if len(response.LibraryItems) != 0 {
			t.Errorf("Expected empty library items, got %d", len(response.LibraryItems))
		}
		mockCache.AssertExpectations(t)
	})

	// ------ SetCachedLibraryItemsBFF() ------
	/*
		GIVEN a cache wrapper that successfully sets cached BFF results
		WHEN SetCachedLibraryItemsBFF is called
		THEN it should not return an error
	*/
	t.Run(`SetCachedLibraryItemsBFF success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "library:bff:test-user-id", testBFFResponse).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedLibraryItemsBFF(context.Background(), testUserID, testBFFResponse)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error when setting cached BFF results
		WHEN SetCachedLibraryItemsBFF is called
		THEN it should return the error
	*/
	t.Run(`SetCachedLibraryItemsBFF error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "library:bff:test-user-id", testBFFResponse).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedLibraryItemsBFF(context.Background(), testUserID, testBFFResponse)

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
		mockCache.On("DeleteCacheKey", mock.Anything, "library:test-user-id").Return(nil)
		mockCache.On("DeleteCacheKey", mock.Anything, "library:bff:test-user-id").Return(nil)

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
		mockCache.On("DeleteCacheKey", mock.Anything, "library:test-user-id").Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateUserCache(context.Background(), testUserID)

		// THEN
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
		mockCache.On("DeleteCacheKey", mock.Anything, "library:test-user-id:game:123").Return(nil)

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
		mockCache.On("DeleteCacheKey", mock.Anything, "library:test-user-id:game:123").Return(testError)

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
