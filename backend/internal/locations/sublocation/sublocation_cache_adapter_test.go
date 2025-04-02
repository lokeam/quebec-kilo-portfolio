package sublocation

import (
	"context"
	"errors"
	"testing"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/stretchr/testify/mock"
)

/*
Behavior:
- Getting all cached sublocations for a user
- Setting all cached sublocations for a user
- Getting a single cached sublocation for a user
- Setting a single cached sublocation for a user
- Invalidating all cached sublocations for a user
- Invalidating a single cached sublocation for a user

Scenarios:
- GetCachedSublocations:
  - Cache hit
  - Cache miss
  - Cache error
- SetCachedSublocations:
  - Success
  - Error
- GetSingleCachedSublocation:
  - Cache hit
  - Cache miss
  - Cache error
- SetSingleCachedSublocation:
  - Success
  - Error
- InvalidateUserCache:
  - Success
  - Error
- InvalidateSublocationCache:
  - Success
  - Error
*/

type MockCacheWrapper struct {
	mock.Mock
	sublocations []models.Sublocation
}

func (m *MockCacheWrapper) GetCachedResults(
	ctx context.Context,
	key string,
	result any,
) (bool, error) {
	args := m.Called(ctx, key, result)

	if args.Get(0).(bool) && result != nil && len(m.sublocations) > 0 {
		switch value := result.(type) {
		case *[]models.Sublocation:
			// Copy mock sublocations to result pointer
			*value = m.sublocations
		case *models.Sublocation:
			// Copy mock sublocation to result pointer
			*value = m.sublocations[0]
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

func TestSublocationCacheAdapter(t *testing.T) {
	// Setup test data
	testUserID := "test-user-id"
	testSublocationID := "test-sublocation-id"
	testSublocation := models.Sublocation{
		ID:              testSublocationID,
		UserID:          testUserID,
		Name:            "Test Sublocation",
		LocationType:    "shelf",
		BgColor:         "blue",
		Capacity:        20,
	}
	testSublocations := []models.Sublocation{testSublocation}
	testError := errors.New("cache error")

	// Helper fn to create a new adapter w/ mock cache wrapper
	createAdapter := func(mockCache *MockCacheWrapper) *SublocationCacheAdapter {
		adapter, _ := NewSublocationCacheAdapter(mockCache)
		return adapter.(*SublocationCacheAdapter)
	}

	// ------------ GetCachedSublocations
	/*
		GIVEN a cache wrapper that returns cached sublocations
	  WHEN GetCachedSublocations is called
	  THEN it should return the cached sublocations without error
	*/
	t.Run(`GetCachedSublocations - Cache hit`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.sublocations = testSublocations
		mockCache.On("GetCachedResults", mock.Anything, "sublocation:test-user-id", mock.Anything).Return(true, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		sublocations, err := adapter.GetCachedSublocations(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(sublocations) != 1 {
			t.Errorf("Expected 1 sublocation, got %d", len(sublocations))
		}
		if sublocations[0].ID != testSublocationID {
			t.Errorf("Expected sublocation ID %s, got %s", testSublocationID, sublocations[0].ID)
		}
		mockCache.AssertExpectations(t)
	})


	/*
	   GIVEN a cache wrapper that returns a cache miss
	   WHEN GetCachedSublocations is called
	   THEN it should return nil sublocations without error
	*/
	t.Run(`GetCachedSublocations - Cache Miss`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "sublocation:test-user-id", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		sublocations, err := adapter.GetCachedSublocations(context.Background(), testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if sublocations != nil {
			t.Errorf("Expected nil sublocations, got %v", sublocations)
		}

		mockCache.AssertExpectations(t)
	})

	/*
	   GIVEN a cache wrapper that returns an error
	   WHEN GetCachedSublocations is called
	   THEN it should return an error
	*/
	t.Run(`GetCachedSublocations - Cache Error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "sublocation:test-user-id", mock.Anything).Return(false, testError)

		adapter := createAdapter(mockCache)

		// WHEN
		sublocations, err := adapter.GetCachedSublocations(context.Background(), testUserID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		if sublocations != nil {
			t.Errorf("Expected nil sublocations, got %v", sublocations)
		}

		mockCache.AssertExpectations(t)
	})

	/*
	   GIVEN a cache wrapper that successfully sets cached results
	   WHEN SetCachedSublocations is called
	   THEN it should not return an error
	*/
	t.Run(`SetCachedSublocations - Success`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "sublocation:test-user-id", testSublocations).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedSublocations(context.Background(), testUserID, testSublocations)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	/*
	   GIVEN a cache wrapper that returns an error when setting cached results
	   WHEN SetCachedSublocations is called
	   THEN it should return the error
	*/
	t.Run(`SetCachedSublocations - error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "sublocation:test-user-id", testSublocations).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetCachedSublocations(context.Background(), testUserID, testSublocations)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v but instead got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})


	// ----------- GetSingleCachedSublocation -----------
	/*
	   GIVEN a cache wrapper that returns a cached sublocation
	   WHEN GetSingleCachedSublocation is called
	   THEN it should return the cached sublocation without error
	*/
	t.Run(`GetSingleCachedSublocation - Cache Hit`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.sublocations = testSublocations
		mockCache.On("GetCachedResults", mock.Anything, "sublocation:test-user-id:sublocation:test-sublocation-id", mock.Anything).Return(true, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		sublocation, found, err := adapter.GetSingleCachedSublocation(context.Background(), testUserID, testSublocationID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !found {
			t.Errorf("Expected sublocation to be found in cache")
		}
		if sublocation == nil {
			t.Errorf("Expected non-nil sublocation")
		} else if sublocation.ID != testSublocationID {
			t.Errorf("Expected sublocation ID %s, got %s", testSublocationID, sublocation.ID)
		}

		mockCache.AssertExpectations(t)
	})

	/*
	   GIVEN a cache wrapper that returns a cache miss
	   WHEN GetSingleCachedSublocation is called
	   THEN it should indicate that the sublocation was not found without error
	*/
	t.Run(`GetSingleCachedSublocatin - Cache Miss`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "sublocation:test-user-id:sublocation:test-sublocation-id", mock.Anything).Return(false, nil)

		adapter := createAdapter(mockCache)

		// WHEN
		sublocation, found, err := adapter.GetSingleCachedSublocation(context.Background(), testUserID, testSublocationID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if found {
			t.Errorf("Expected found to be false")
		}
		if sublocation != nil {
			t.Errorf("Expected nil sublocation, got %v", sublocation)
		}
		mockCache.AssertExpectations(t)
	})


	/*
	   GIVEN a cache wrapper that returns an error
	   WHEN GetSingleCachedSublocation is called
	   THEN it should return an error
	*/
	t.Run(`GetSingleCachedSublocation - Cache Error`, func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("GetCachedResults", mock.Anything, "sublocation:test-user-id:sublocation:test-sublocation-id", mock.Anything).Return(false, testError)

		adapter := createAdapter(mockCache)

		// WHEN
		sublocation, found, err := adapter.GetSingleCachedSublocation(context.Background(), testUserID, testSublocationID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		if found {
			t.Errorf("Expected found to be false")
		}
		if sublocation != nil {
			t.Errorf("Expected nil sublocation, got %v", sublocation)
		}
		mockCache.AssertExpectations(t)
	})


	/*
	   GIVEN a cache wrapper that successfully sets a cached sublocation
	   WHEN SetSingleCachedSublocation is called
	   THEN it should not return an error
	*/
	t.Run("SetSingleCachedSublocation - success", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "sublocation:test-user-id:sublocation:test-sublocation-id", testSublocation).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetSingleCachedSublocation(context.Background(), testUserID, testSublocation)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	/*
	   GIVEN a cache wrapper that returns an error when setting a cached sublocation
	   WHEN SetSingleCachedSublocation is called
	   THEN it should return the error
	*/
	t.Run("SetSingleCachedSublocation - error", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "sublocation:test-user-id:sublocation:test-sublocation-id", testSublocation).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.SetSingleCachedSublocation(context.Background(), testUserID, testSublocation)

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
	t.Run("InvalidateUserCache - success", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "sublocation:test-user-id", nil).Return(nil)

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
	t.Run("InvalidateUserCache - error", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "sublocation:test-user-id", nil).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateUserCache(context.Background(), testUserID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})

	// ----------- InvalidateSublocationCache -----------
	/*
	   GIVEN a cache wrapper that successfully invalidates the sublocation cache
	   WHEN InvalidateSublocationCache is called
	   THEN it should not return an error
	*/
	t.Run("InvalidateSublocationCache - success", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "sublocation:test-user-id:sublocation:test-sublocation-id", nil).Return(nil)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateSublocationCache(context.Background(), testUserID, testSublocationID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		mockCache.AssertExpectations(t)
	})

	/*
	   GIVEN a cache wrapper that returns an error when invalidating sublocation cache
	   WHEN InvalidateSublocationCache is called
	   THEN it should return the error
	*/
	t.Run("InvalidateSublocationCache - error", func(t *testing.T) {
		// GIVEN
		mockCache := new(MockCacheWrapper)
		mockCache.On("SetCachedResults", mock.Anything, "sublocation:test-user-id:sublocation:test-sublocation-id", nil).Return(testError)

		adapter := createAdapter(mockCache)

		// WHEN
		err := adapter.InvalidateSublocationCache(context.Background(), testUserID, testSublocationID)

		// THEN
		if err != testError {
			t.Errorf("Expected error %v, got %v", testError, err)
		}
		mockCache.AssertExpectations(t)
	})


}