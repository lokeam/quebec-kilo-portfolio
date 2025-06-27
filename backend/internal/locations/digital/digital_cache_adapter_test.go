package digital

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/stretchr/testify/assert"
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
}

func (m *MockCacheWrapper) GetCachedResults(ctx context.Context, key string, result interface{}) (bool, error) {
	args := m.Called(ctx, key, result)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheWrapper) SetCachedResults(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
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
		expectedLocations := []models.DigitalLocation{testLocation}

		mockCache.On("GetCachedResults", mock.Anything, "digital:test-user-id", mock.Anything).
			Return(true, nil).
			Run(func(args mock.Arguments) {
				// Copy the test location data into the result parameter
				result := args.Get(2).(*[]models.DigitalLocation)
				*result = expectedLocations
			})

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
	t.Run(`GetCachedDigitalLocations - Cache error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "digital:test-user-id", mock.Anything).Return(false, testError)

		adapter := createAdapter(mockCache)

		// WHEN
		locations, err := adapter.GetCachedDigitalLocations(context.Background(), testUserID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
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
	t.Run(`SetCachedDigitalLocations - Success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		locations := []models.DigitalLocation{testLocation}
		mockCache.On("SetCachedResults", mock.Anything, "digital:test-user-id", locations).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedDigitalLocations(context.Background(), testUserID, locations)

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
	t.Run(`SetCachedDigitalLocations - Error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		locations := []models.DigitalLocation{testLocation}
		mockCache.On("SetCachedResults", mock.Anything, "digital:test-user-id", locations).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedDigitalLocations(context.Background(), testUserID, locations)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
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
		mockCache.On("GetCachedResults", mock.Anything, "digital:test-user-id:location:test-location-id", mock.Anything).
			Return(true, nil).
			Run(func(args mock.Arguments) {
				// Copy the test location data into the result parameter
				result := args.Get(2).(*models.DigitalLocation)
				*result = testLocation
			})

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
		locationKey := fmt.Sprintf("digital:%s:location:%s", testUserID, testLocationID)
		userKey := fmt.Sprintf("digital:%s", testUserID)
		mockCache.On("SetCachedResults", mock.Anything, locationKey, nil).Return(nil)
		mockCache.On("SetCachedResults", mock.Anything, userKey, nil).Return(nil)

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
		locationKey := fmt.Sprintf("digital:%s:location:%s", testUserID, testLocationID)
		mockCache.On("SetCachedResults", mock.Anything, locationKey, nil).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateDigitalLocationCache(context.Background(), testUserID, testLocationID)

		// THEN
		if err == nil {
			t.Error("Expected an error, got nil")
		}
		expectedErrMsg := fmt.Sprintf("failed to invalidate cache for location %s (user: %s, key: %s): %v",
			testLocationID, testUserID, locationKey, testError)
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message %q, got %q", expectedErrMsg, err.Error())
		}
		mockCache.AssertExpectations(t)
	})
}

func TestGetCachedSubscription(t *testing.T) {
	ctx := context.Background()
	locationID := "test-location"
	subscription := models.Subscription{
		ID:              1,
		LocationID:      locationID,
		BillingCycle:    "monthly",
		CostPerCycle:    9.99,
		NextPaymentDate: time.Now(),
		PaymentMethod:   "Visa",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockCache := &MockCacheWrapper{}
	mockCache.On("GetCachedResults", ctx, fmt.Sprintf("digital:subscription:%s", locationID), mock.Anything).
		Return(true, nil).
		Run(func(args mock.Arguments) {
			result := args.Get(2).(*models.Subscription)
			*result = subscription
		})

	adapter, err := NewDigitalCacheAdapter(mockCache)
	assert.NoError(t, err)

	// Test cache hit
	result, found, err := adapter.GetCachedSubscription(ctx, locationID)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, subscription.ID, result.ID)

	// Test cache miss
	mockCache = &MockCacheWrapper{}
	mockCache.On("GetCachedResults", ctx, fmt.Sprintf("digital:subscription:%s", locationID), mock.Anything).
		Return(false, nil)

	adapter, err = NewDigitalCacheAdapter(mockCache)
	assert.NoError(t, err)

	result, found, err = adapter.GetCachedSubscription(ctx, locationID)
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, result)

	mockCache.AssertExpectations(t)
}

func TestSetCachedSubscription(t *testing.T) {
	ctx := context.Background()
	locationID := "test-location"
	subscription := models.Subscription{
		ID:              1,
		LocationID:      locationID,
		BillingCycle:    "monthly",
		CostPerCycle:    9.99,
		NextPaymentDate: time.Now(),
		PaymentMethod:   "Visa",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockCache := &MockCacheWrapper{}
	mockCache.On("SetCachedResults", ctx, fmt.Sprintf("digital:subscription:%s", locationID), subscription).
		Return(nil)

	adapter, err := NewDigitalCacheAdapter(mockCache)
	assert.NoError(t, err)

	err = adapter.SetCachedSubscription(ctx, locationID, subscription)
	assert.NoError(t, err)

	mockCache.AssertExpectations(t)
}

func TestInvalidateSubscriptionCache(t *testing.T) {
	ctx := context.Background()
	locationID := "test-location"

	mockCache := &MockCacheWrapper{}
	mockCache.On("SetCachedResults", ctx, fmt.Sprintf("digital:subscription:%s", locationID), nil).
		Return(nil)

	adapter, err := NewDigitalCacheAdapter(mockCache)
	assert.NoError(t, err)

	err = adapter.InvalidateSubscriptionCache(ctx, locationID)
	assert.NoError(t, err)

	mockCache.AssertExpectations(t)
}

func TestGetCachedPayments(t *testing.T) {
	ctx := context.Background()
	locationID := "test-location"
	payments := []models.Payment{
		{
			ID:            1,
			LocationID:    locationID,
			Amount:        1000,
			PaymentDate:   time.Now(),
			PaymentMethod: "Visa",
			CreatedAt:     time.Now(),
		},
		{
			ID:            2,
			LocationID:    locationID,
			Amount:        2000,
			PaymentDate:   time.Now(),
			PaymentMethod: "Mastercard",
			CreatedAt:     time.Now(),
		},
	}

	mockCache := &MockCacheWrapper{}
	mockCache.On("GetCachedResults", ctx, fmt.Sprintf("digital:payments:%s", locationID), mock.Anything).
		Return(true, nil).
		Run(func(args mock.Arguments) {
			result := args.Get(2).(*[]models.Payment)
			*result = payments
		})

	adapter, err := NewDigitalCacheAdapter(mockCache)
	assert.NoError(t, err)

	// Test cache hit
	result, err := adapter.GetCachedPayments(ctx, locationID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, payments[0].ID, result[0].ID)
	assert.Equal(t, payments[1].ID, result[1].ID)

	// Test cache miss
	mockCache = &MockCacheWrapper{}
	mockCache.On("GetCachedResults", ctx, fmt.Sprintf("digital:payments:%s", locationID), mock.Anything).
		Return(false, nil)

	adapter, err = NewDigitalCacheAdapter(mockCache)
	assert.NoError(t, err)

	result, err = adapter.GetCachedPayments(ctx, locationID)
	assert.NoError(t, err)
	assert.Nil(t, result)

	mockCache.AssertExpectations(t)
}

func TestSetCachedPayments(t *testing.T) {
	ctx := context.Background()
	locationID := "test-location"
	payments := []models.Payment{
		{
			ID:            1,
			LocationID:    locationID,
			Amount:        1000,
			PaymentDate:   time.Now(),
			PaymentMethod: "Visa",
			CreatedAt:     time.Now(),
		},
	}

	mockCache := &MockCacheWrapper{}
	mockCache.On("SetCachedResults", ctx, fmt.Sprintf("digital:payments:%s", locationID), payments).
		Return(nil)

	adapter, err := NewDigitalCacheAdapter(mockCache)
	assert.NoError(t, err)

	err = adapter.SetCachedPayments(ctx, locationID, payments)
	assert.NoError(t, err)

	mockCache.AssertExpectations(t)
}

func TestInvalidatePaymentsCache(t *testing.T) {
	ctx := context.Background()
	locationID := "test-location"

	mockCache := &MockCacheWrapper{}
	mockCache.On("SetCachedResults", ctx, fmt.Sprintf("digital:payments:%s", locationID), nil).
		Return(nil)

	adapter, err := NewDigitalCacheAdapter(mockCache)
	assert.NoError(t, err)

	err = adapter.InvalidatePaymentsCache(ctx, locationID)
	assert.NoError(t, err)

	mockCache.AssertExpectations(t)
}

func TestInvalidateDigitalLocationCache(t *testing.T) {
	ctx := context.Background()
	userID := "test-user"
	locationID := "test-location"
	expectedLocationKey := fmt.Sprintf("digital:%s:location:%s", userID, locationID)
	expectedUserKey := fmt.Sprintf("digital:%s", userID)

	t.Run("successful cache invalidation", func(t *testing.T) {
		mockCache := &MockCacheWrapper{}
		mockCache.On("SetCachedResults", ctx, expectedLocationKey, nil).Return(nil)
		mockCache.On("SetCachedResults", ctx, expectedUserKey, nil).Return(nil)

		adapter, err := NewDigitalCacheAdapter(mockCache)
		assert.NoError(t, err)

		err = adapter.InvalidateDigitalLocationCache(ctx, userID, locationID)
		assert.NoError(t, err)

		mockCache.AssertExpectations(t)
	})

	t.Run("location cache invalidation failure", func(t *testing.T) {
		expectedErr := errors.New("cache error")
		mockCache := &MockCacheWrapper{}
		mockCache.On("SetCachedResults", ctx, expectedLocationKey, nil).Return(expectedErr)

		adapter, err := NewDigitalCacheAdapter(mockCache)
		assert.NoError(t, err)

		err = adapter.InvalidateDigitalLocationCache(ctx, userID, locationID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to invalidate cache for location")
		assert.Contains(t, err.Error(), userID)
		assert.Contains(t, err.Error(), expectedLocationKey)
		assert.Contains(t, err.Error(), expectedErr.Error())

		mockCache.AssertExpectations(t)
	})

	t.Run("user cache invalidation failure", func(t *testing.T) {
		expectedErr := errors.New("cache error")
		mockCache := &MockCacheWrapper{}
		mockCache.On("SetCachedResults", ctx, expectedLocationKey, nil).Return(nil)
		mockCache.On("SetCachedResults", ctx, expectedUserKey, nil).Return(expectedErr)

		adapter, err := NewDigitalCacheAdapter(mockCache)
		assert.NoError(t, err)

		err = adapter.InvalidateDigitalLocationCache(ctx, userID, locationID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to invalidate user cache")
		assert.Contains(t, err.Error(), userID)
		assert.Contains(t, err.Error(), expectedUserKey)
		assert.Contains(t, err.Error(), expectedErr.Error())

		mockCache.AssertExpectations(t)
	})
}

