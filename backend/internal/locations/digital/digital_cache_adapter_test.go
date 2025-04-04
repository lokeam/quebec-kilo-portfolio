package digital

import (
	"context"
	"errors"
	"testing"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/stretchr/testify/mock"
)

/*
Behavior:
- Getting all cached digital locations for a user
- Setting all cached digital locations for a user
- Getting a single cached digital location for a user
- Setting a single cached digital location for a user
- Invalidating all cached digital locations for a user
- Invalidating a single cached digital location for a user

Scenarios:
- GetCachedDigitalLocations:
  - Cache hit
  - Cache miss
  - Cache error
- SetCachedDigitalLocations:
  - Success
  - Error
- GetSingleCachedDigitalLocation:
  - Cache hit
  - Cache miss
  - Cache error
- SetSingleCachedDigitalLocation:
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
	locations []models.DigitalLocation
}

func (m *MockCacheWrapper) GetCachedResults(
	ctx context.Context,
	key string,
	result any,
) (bool, error) {
	args := m.Called(ctx, key, result)

	if args.Get(0).(bool) && result != nil && len(m.locations) > 0 {
		switch value := result.(type) {
		case *[]models.DigitalLocation:
			// Copy mock locations to result pointer
			*value = m.locations
		case *models.DigitalLocation:
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

func TestDigitalCacheAdapter(t *testing.T) {
	// Setup test data
	testUserID := "test-user-id"
	testLocationID := "test-location-id"
	testLocation := models.DigitalLocation{
		ID:       testLocationID,
		UserID:   testUserID,
		Name:     "Test Location",
		IsActive: true,
		URL:      "https://example.com",
	}
	testLocations := []models.DigitalLocation{testLocation}
	testError := errors.New("cache error")

	// Helper function to create a new adapter with mock cache wrapper
	createAdapter := func(mockCache *MockCacheWrapper) *DigitalCacheAdapter {
		adapter, _ := NewDigitalCacheAdapter(mockCache)
		return adapter.(*DigitalCacheAdapter)
	}

	// ----------- GetCachedDigitalLocations -----------
	/*
		GIVEN a cache wrapper that returns cached locations
		WHEN GetCachedDigitalLocations is called
		THEN it should return the cached locations without error
	*/
	t.Run(`GetCachedDigitalLocations - Cache hit`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.locations = testLocations
		mockCache.On("GetCachedResults", mock.Anything, "digital:test-user-id", mock.Anything).Return(true, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		locations, err := adapter.GetCachedDigitalLocations(context.Background(), testUserID)

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
		WHEN GetCachedDigitalLocations is called
		THEN it should return nil locations without error
	*/
	t.Run(`GetCachedDigitalLocations - Cache miss`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "digital:test-user-id", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		locations, err := adapter.GetCachedDigitalLocations(context.Background(), testUserID)

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
		WHEN GetCachedDigitalLocations is called
		THEN it should return an error
	*/
	t.Run(`GetCachedDigitalLocations - cache miss`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "digital:test-user-id", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		locations, err := adapter.GetCachedDigitalLocations(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if locations != nil {
			t.Errorf("Expected nil locations, got %v", locations)
		}
		mockCache.AssertExpectations(t)
	})

	// ----------- SetCachedDigitalLocations -----------
	/*
		GIVEN a cache wrapper that successfully sets cached results
		WHEN SetCachedDigitalLocations is called
		THEN it should not return an error
	*/
	t.Run(`SetSingleCachedDigitalLocation - Success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "digital:test-user-id:location:test-location-id", testLocation).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetSingleCachedDigitalLocation(context.Background(), testUserID, testLocation)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	/*
		GIVEN a cache wrapper that returns an error when setting cached results
		WHEN SetCachedDigitalLocations is called
		THEN it should return the error
	*/
	t.Run(`InvalidateUserCache - success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "digital:test-user-id", nil).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateUserCache(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error but instead got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	// ----------- GetSingleCachedDigitalLocation -----------
	/*
		GIVEN a cache wrapper that returns a cached location
		WHEN GetSingleCachedDigitalLocation is called
		THEN it should return the cached location without error
	*/
	t.Run(`GetSingleCachedDigitalLocation - cache hit`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.locations = testLocations
		mockCache.On("GetCachedResults", mock.Anything, "digital:test-user-id:location:test-location-id", mock.Anything).Return(true, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		location, found, err := adapter.GetSingleCachedDigitalLocation(context.Background(), testUserID, testLocationID)

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
		WHEN GetSingleCachedDigitalLocation is called
		THEN it should indicate that the location was not found without error
	*/
	t.Run(`GetSingleCachedDigitalLocation - cache miss`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "digital:test-user-id:location:test-location-id", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		location, found, err := adapter.GetSingleCachedDigitalLocation(context.Background(), testUserID, testLocationID)

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
		WHEN GetSingleCachedDigitalLocation is called
		THEN it should return an error
	*/
	t.Run(`GetSingleCachedDigitalLocation - cache error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "digital:test-user-id:location:test-location-id", mock.Anything).Return(false, testError)

		adapter := createAdapter(mockCache)

		// WHEN
		location, found, err := adapter.GetSingleCachedDigitalLocation(context.Background(), testUserID, testLocationID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		if found {
			t.Errorf("Expected found to be false")
		}
		if location != nil {
			t.Errorf("Expected nil location, got %v", location)
		}
		mockCache.AssertExpectations(t)
	})


	// ----------- SetSingleCachedDigitalLocation -----------
	/*
		GIVEN a cache wrapper that successfully sets a cached location
		WHEN SetSingleCachedDigitalLocation is called
		THEN it should not return an error
	*/
	t.Run(`SetSingleCachedDigitalLocation - success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "digital:test-user-id:location:test-location-id", testLocation).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetSingleCachedDigitalLocation(context.Background(), testUserID, testLocation)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)
	})


	/*
		GIVEN a cache wrapper that returns an error when setting a cached location
		WHEN SetSingleCachedDigitalLocation is called
		THEN it should return the error
	*/
	t.Run(`SetSingleCachedDigitalLocation - error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "digital:test-user-id:location:test-location-id", testLocation).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetSingleCachedDigitalLocation(context.Background(), testUserID, testLocation)

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
	t.Run(`InvalidateUserCache - success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "digital:test-user-id", nil).Return(nil)

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
	t.Run(`InvalidateUserCache - error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "digital:test-user-id", nil).Return(testError)

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
	t.Run(`InvalidateLocationCache - success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "digital:test-user-id:location:test-location-id", nil).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateDigitalLocationCache(context.Background(), testUserID, testLocationID)

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
	t.Run(`InvalidateLocationCache - error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "digital:test-user-id:location:test-location-id", nil).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateDigitalLocationCache(context.Background(), testUserID, testLocationID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})
}

