package physical

import (
	"context"
	"errors"
	"testing"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/stretchr/testify/mock"
)

/*
	Behavior:
	- Getting all cached physical locations for a user
	- Setting all cached physical locations for a user
	- Getting a single cached physical location for a user
	- Setting a single cached physical location for a user
	- Invalidating all cached physical locations for a user
	- Invalidating a single cached physical location for a user

	Scenarios:
	- GetCachedPyysicalLocations:
		- Cache hit
		- Cache miss
		- Cache error
	- SetCachedPhysicalLocations:
		- Success
		- Error
	- GetSingleCachedPhysicalLocation:
		- Cache hit
		- Cache miss
		- Cache error
	- SetSingleCachedPhysicalLocation:
		- Success
		- Error
	- InvalidateUserCache:
		- Success
		- Error
	- InvalidateLocationCache:
		- Success
		- Error
*/

type MockCacheWrapper struct {
	mock.Mock
	locations []models.PhysicalLocation
}

// Mock implementation of CacheWrapper methods

func (m *MockCacheWrapper) GetCachedResults(
	ctx context.Context,
	key string,
	result any,
) (bool, error) {
	args := m.Called(ctx, key, result)

	// Try and copy data if:
	// 1. There is a cache hit
	// 2. We have locations to copy
	if args.Get(0).(bool) && result != nil && len(m.locations) > 0 {
		switch value := result.(type) {
		case *[]models.PhysicalLocation:
			// Copy mock locations to result pointer
			*value = m.locations
		case *models.PhysicalLocation:
			// Copy mock location to result pointer
			*value = m.locations[0]
		}
	}

	return args.Get(0).(bool), args.Error(1)
}

func (m *MockCacheWrapper) SetCachedResults(
	ctx context.Context,
	key string,
	result any,
) error {
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

func TestPhysicalCacheAdapter(t *testing.T) {
	// Setup test data
	testUserID := "test-user-id"
	testLocationID := "test-location-id"
	testLocation := models.PhysicalLocation{
		ID:             testLocationID,
		UserID:         testUserID,
		Name:           "Test Location",
		Label:          "Home",
		LocationType:   "Residence",
		MapCoordinates: models.PhysicalMapCoordinates{
			Coords:         "123,456",
			GoogleMapsLink: "https://www.google.com/maps?q=123,456",
		},
	}
	testLocations := []models.PhysicalLocation{testLocation}
	testError := errors.New("cache error")

	// Helper function to create a new adapter with mock cache wrapper
	createAdapter := func(mockCache *MockCacheWrapper) *PhysicalCacheAdapter {
		adapter, _ := NewPhysicalCacheAdapter(mockCache)
		return adapter.(*PhysicalCacheAdapter)
	}

	// ----------- GetCachedPhysicalLocations -----------
	/*
		GIVEN a cache wrapper that returns cached locations
		WHEN GetCachedPhysicalLocations is called
		THEN it should return the cached locations without error
	*/
	t.Run(`GetCachedPhysicalLocations - Cache hit`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.locations = testLocations
		mockCache.On("GetCachedResults", mock.Anything, "physical:test-user-id", mock.Anything).Return(true, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		locations, err := adapter.GetCachedPhysicalLocations(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != 1 {
			t.Errorf("Expected 1 location, got %d", len(locations))
		}
		if locations[0].ID != testLocationID {
			t.Errorf("Expected location ID %s, got %s", testLocationID, locations[0].ID)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns a cache miss
		WHEN GetCachedPhysicalLocations is called
		THEN it should return nil locations without error
	*/
	t.Run(`GetCachedPhysicalLocations - cache miss`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "physical:test-user-id", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		locations, err := adapter.GetCachedPhysicalLocations(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if locations != nil {
			t.Errorf("Expected nil locations, got %v", locations)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error
		WHEN GetCachedPhysicalLocations is called
		THEN it should return an error
	*/
	t.Run(`GetCachedPhysicalLocations - cache error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "physical:test-user-id", mock.Anything).Return(false, testError)

		adapter := createAdapter(mockCache)

		// WHEN
		locations, err := adapter.GetCachedPhysicalLocations(context.Background(), testUserID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		if locations != nil {
			t.Errorf("Expected nil locations, got %v", locations)
		}

		mockCache.AssertExpectations(t)
	})

	// ----------- SetCachedPhysicalLocations -----------
	/*
		GIVEN a cache wrapper that succesfully sets cached results
		WHEN SetCachedPhysicalLocations is called
		THEN it should not return an error
	*/
	t.Run(`SetCachedPhysicalLocations - success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "physical:test-user-id", testLocations).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedPhysicalLocations(context.Background(), testUserID, testLocations)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error when setting cached results
		WHEN SetCachedPhysicalLocations is called
		THEN it should return the error
	*/
	t.Run(`SetCachedPhysicalLocations - error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "physical:test-user-id", testLocations).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedPhysicalLocations(context.Background(), testUserID, testLocations)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v but instead got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})

	// ----------- GetSingleCachedPhysicalLocation -----------
	/*
		GIVEN a cache wrapper that returns a cached location
		WHEN GetCachedPhysicalLocation is called
		THEN it should return the cached location without error
	*/
	t.Run(`GetSingleCachedPhysicalLocation - cache hit`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.locations = testLocations
		mockCache.On("GetCachedResults", mock.Anything, "physical:test-user-id:location:test-location-id", mock.Anything).Return(true, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		location, found, err := adapter.GetSingleCachedPhysicalLocation(context.Background(), testUserID, testLocationID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !found {
			t.Errorf("Expected location to be found in cache")
		}
		if location == nil {
			t.Errorf("Expected non-nil location")
		} else if location.ID != testLocationID {
			t.Errorf("Expected location ID %s, got %s", testLocationID, location.ID)
		}

		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns a cache miss
		WHEN GetCachedPhysicalLocation is called
		THEN it should indicate that the location was found without error
	*/
	t.Run("GetSingleCachedPhysicalLocation with cache miss", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "physical:test-user-id:location:test-location-id", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		location, found, err := adapter.GetSingleCachedPhysicalLocation(context.Background(), testUserID, testLocationID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if found {
			t.Errorf("Expected found to be false")
		}
		if location != nil {
			t.Errorf("Expected nil location, got %v", location)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error
		WHEN GetCachedPhysicalLocation is called
		THEN it should return an error
	*/
	t.Run("GetCachedPhysicalLocations with cache error", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "physical:test-user-id", mock.Anything).Return(false, testError)

		adapter := createAdapter(mockCache)

		// WHEN
		locations, err := adapter.GetCachedPhysicalLocations(context.Background(), testUserID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		if locations != nil {
			t.Errorf("Expected nil locations, got %v", locations)
		}
		mockCache.AssertExpectations(t)
	})

	// ----------- SetSingleCachedPhysicalLocation -----------
	/*
		GIVEN a cache wrapper that successfully sets a cached location
		WHEN SetSingleCachedPhysicalLocation is called
		THEN it should not return an error
	*/
	t.Run("SetCachedPhysicalLocations success", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "physical:test-user-id", testLocations).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedPhysicalLocations(context.Background(), testUserID, testLocations)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error when setting a cached location
		WHEN SetCachedPhysicalLocation is called
		THEN it should return the error
	*/
	t.Run("SetCachedPhysicalLocations error", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "physical:test-user-id", testLocations).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedPhysicalLocations(context.Background(), testUserID, testLocations)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v but instead got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})

	// ----------- InvalidateUserCache -----------
	/*
		GIVEN a cache wrapper that successfully invalidates the user cache
		WHEN InvalidateUserCache is called
		THEN it should not return an error
	*/
	t.Run("InvalidateUserCache success", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("DeleteCacheKey", mock.Anything, "physical:test-user-id").Return(nil)
		mockCache.On("DeleteCacheKey", mock.Anything, "physical:bff:test-user-id").Return(nil)

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
	t.Run("InvalidateUserCache error", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("DeleteCacheKey", mock.Anything, "physical:test-user-id").Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateUserCache(context.Background(), testUserID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})


	// ----------- InvalidateLocationCache -----------
	/*
		GIVEN a cache wrapper that successfully invalidates the location cache
		WHEN InvalidateLocationCache is called
		THEN it should not return an error
	*/
	t.Run("InvalidateLocationCache success", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("DeleteCacheKey", mock.Anything, "physical:test-user-id:location:test-location-id").Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateLocationCache(context.Background(), testUserID, testLocationID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)
	})


	/*
		GIVEN a cache wrapper that returns an error when invalidating location cache
		WHEN InvalidateLocationCache is called
		THEN it should return the error
	*/
	t.Run("InvaliateLocationCache error", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("DeleteCacheKey", mock.Anything, "physical:test-user-id:location:test-location-id").Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateLocationCache(context.Background(), testUserID, testLocationID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})
}