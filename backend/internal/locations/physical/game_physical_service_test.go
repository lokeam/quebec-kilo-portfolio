package physical

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	"github.com/stretchr/testify/mock"
)

/*
	Behavior:
	- GetPhysicalLocations retrieves physical locations for a user
		- Attempts to retrieve from cache first
		- Falls back to database on cache miss
		- Caches database results for future requests

	- AddPhysicalLocation adds a new physical location
		- Validates the location data
		- Adds the location to the database
		- Invalidates the user's cache

	- UpdatePhysicalLocation updates an existing physical location
		- Validates the location data
		- Updates the location in the database
		- Invalidates both the user's cache and the specific location cache

	- DeletePhysicalLocation deletes a physical location
		- Deletes the location from the database
		- Invalidates both the user's cache and the specific location cache

	Scenarios:
		- GetPhysicalLocations:
			- Cache hit
			- Cache miss
			- Database error
			- No locations found

		- AddPhysicalLocation:
			- validation success + db success
			- validation failure
			- validation success + db failure
			- cache invalidation failure (should not block success)

		- UpdatePhysicalLocation:
			- validation successs + db success
			- validation failure
			- validation success + db failure
			- cache invalidation failure (should not block success)

		- DeletePhysicalLocation:
			- db success
			- db failure
			- cache invalidation failure (should not block success)
*/

type MockPhysicalDbAdapter struct {
	mock.Mock
}

func (m *MockPhysicalDbAdapter) GetPhysicalLocation(ctx context.Context, userID string, locationID string) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, locationID)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) GetUserPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) AddPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, location)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, location)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) RemovePhysicalLocation(ctx context.Context, userID string, locationID string) error {
	args := m.Called(ctx, userID, locationID)
	return args.Error(0)
}

// func newMockGamePhysicalServiceWithDefaults(logger *testutils.TestLogger) *GamePhysicalService {
// 	mockConfig := mocks.NewMockConfig()
// 	mockDbAdapter := &MockPhysicalDbAdapter{}

// 	// Set up default mock responses
// 	mockDbAdapter.On("GetPhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
// 		Return(models.PhysicalLocation{}, nil)
// 	mockDbAdapter.On("GetUserPhysicalLocations", mock.Anything, mock.Anything).
// 		Return([]models.PhysicalLocation{}, nil)
// 	mockDbAdapter.On("AddPhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
// 		Return(models.PhysicalLocation{}, nil)
// 	mockDbAdapter.On("UpdatePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
// 		Return(models.PhysicalLocation{}, nil)
// 	mockDbAdapter.On("RemovePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
// 		Return(nil)

// 	return &GamePhysicalService{
// 		dbAdapter:      mockDbAdapter,
// 		config:         mockConfig,
// 		logger:         logger,
// 		validator:      mocks.DefaultPhysicalValidator(),
// 		sanitizer:      mocks.DefaultSanitizer(),
// 		cacheWrapper:   mocks.DefaultPhysicalCacheWrapper(),
// 	}
// }

func TestGamePhysicalService(t *testing.T) {
	ctx := context.Background()

	// Create a simple function to create the service with mocks for testing
	createTestService := func(mockDb *mocks.MockPhysicalDbAdapter, mockCache *mocks.MockPhysicalCacheWrapper) *GamePhysicalService {
		mockLogger := testutils.NewTestLogger()

		// Create mock validator and sanitizer
		mockValidator := &mocks.MockPhysicalValidator{
			ValidatePhysicalLocationFunc: func(location models.PhysicalLocation) (models.PhysicalLocation, error) {
				// Simply return the location as valid by default
				return location, nil
			},
		}
		mockSanitizer := &mocks.MockSanitizer{
			SanitizeFunc: func(s string) (string, error) {
				return s, nil // Default implementation just returns the input
			},
		}

		// Create a minimal service without calling the constructor
		return &GamePhysicalService{
			dbAdapter:    mockDb,
			cacheWrapper: mockCache,
			logger:       mockLogger,
			validator:    mockValidator,
			sanitizer:    mockSanitizer,
		}
	}

	// Create fresh mocks for this test
	mockDb := new(mocks.MockPhysicalDbAdapter)
	mockCache := new(mocks.MockPhysicalCacheWrapper)

	// Create the service with mocks
	service := createTestService(mockDb, mockCache)

	// Test GetPhysicalLocation
	t.Run("GetPhysicalLocation", func(t *testing.T) {
		expectedLocation := models.PhysicalLocation{
			ID:             "loc1",
			UserID:         "user1",
			Name:           "Test Location",
			Label:          "Test Label",
			LocationType:   "type1",
			MapCoordinates: "1.0,2.0",
		}

		// Be explicit about the nil pointer return to avoid type conversion issues
		var nilLocation *models.PhysicalLocation = nil
		mockCache.On("GetSingleCachedPhysicalLocation", ctx, "user1", "loc1").
			Return(nilLocation, false, fmt.Errorf("cache miss"))
		mockDb.On("GetPhysicalLocation", ctx, "user1", "loc1").
			Return(expectedLocation, nil)
		mockCache.On("SetSingleCachedPhysicalLocation", ctx, "user1", expectedLocation).
			Return(nil)

		location, err := service.GetPhysicalLocation(ctx, "user1", "loc1")
		if err != nil {
			t.Errorf("GetPhysicalLocation() error = %v", err)
			return
		}
		if !reflect.DeepEqual(location, expectedLocation) {
			t.Errorf("GetPhysicalLocation() location = %v, expectedLocation %v", location, expectedLocation)
		}

		mockDb.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	// Test GetUserPhysicalLocations
	t.Run("GetUserPhysicalLocations", func(t *testing.T) {
		expectedLocations := []models.PhysicalLocation{
			{
				ID:             "loc1",
				UserID:         "user1",
				Name:           "Test Location 1",
				Label:          "Test Label 1",
				LocationType:   "type1",
				MapCoordinates: "1.0,2.0",
			},
			{
				ID:             "loc2",
				UserID:         "user1",
				Name:           "Test Location 2",
				Label:          "Test Label 2",
				LocationType:   "type2",
				MapCoordinates: "3.0,4.0",
			},
		}

		// Return an empty slice instead of nil to avoid type casting issues
		mockCache.On("GetCachedPhysicalLocations", ctx, "user1").
			Return([]models.PhysicalLocation{}, fmt.Errorf("cache miss"))
		mockDb.On("GetUserPhysicalLocations", ctx, "user1").
			Return(expectedLocations, nil)
		mockCache.On("SetCachedPhysicalLocations", ctx, "user1", expectedLocations).
			Return(nil)

		locations, err := service.GetUserPhysicalLocations(ctx, "user1")
		if err != nil {
			t.Errorf("GetUserPhysicalLocations() error = %v", err)
			return
		}
		if !reflect.DeepEqual(locations, expectedLocations) {
			t.Errorf("GetUserPhysicalLocations() locations = %v, expectedLocations %v", locations, expectedLocations)
		}

		mockDb.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	// Test AddPhysicalLocation
	t.Run("AddPhysicalLocation", func(t *testing.T) {
		// Create fresh mocks for this test
		mockDb := new(mocks.MockPhysicalDbAdapter)
		mockCache := new(mocks.MockPhysicalCacheWrapper)
		service := createTestService(mockDb, mockCache)

		location := models.PhysicalLocation{
			ID:             "loc1",
			UserID:         "user1",
			Name:           "Test Location",
			Label:          "Test Label",
			LocationType:   "type1",
			MapCoordinates: "1.0,2.0",
		}

		mockDb.On("AddPhysicalLocation", ctx, "user1", location).
			Return(location, nil)

		// After adding, it fetches all locations to update cache
		allLocations := []models.PhysicalLocation{location}
		mockDb.On("GetUserPhysicalLocations", ctx, "user1").
			Return(allLocations, nil)
		mockCache.On("SetCachedPhysicalLocations", ctx, "user1", allLocations).
			Return(nil)
		mockCache.On("SetSingleCachedPhysicalLocation", ctx, "user1", location).
			Return(nil)

		createdLocation, err := service.AddPhysicalLocation(ctx, "user1", location)
		if err != nil {
			t.Errorf("AddPhysicalLocation() error = %v", err)
			return
		}
		if !reflect.DeepEqual(createdLocation, location) {
			t.Errorf("AddPhysicalLocation() createdLocation = %v, location %v", createdLocation, location)
		}

		mockDb.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	// Test UpdatePhysicalLocation
	t.Run("UpdatePhysicalLocation", func(t *testing.T) {
		userID := "9a4aeee6-fb31-4839-a921-f61b0525046d"
		locationID := "loc1" // This should be the ID of the existing "John's condo" location

		// existingLocation := models.PhysicalLocation{
		// 	ID:             locationID,
		// 	UserID:         userID,
		// 	Name:           "John's condo",
		// 	Label:          "Home",
		// 	LocationType:   "apartment",
		// 	MapCoordinates: "40.69041162815012, -74.04432918344848",
		// }

		updatedLocation := models.PhysicalLocation{
			ID:             locationID,
			UserID:         userID,
			Name:           "John's condo",
			Label:          "Home",
			LocationType:   "apartment",
			MapCoordinates: "40.69041162815012, -74.04432918344848",
		}

		mockDb.On("UpdatePhysicalLocation", ctx, userID, updatedLocation).
			Return(updatedLocation, nil)
		mockCache.On("SetSingleCachedPhysicalLocation", ctx, userID, updatedLocation).
			Return(nil)
		mockCache.On("InvalidateUserCache", ctx, userID).
			Return(nil)
		mockDb.On("GetUserPhysicalLocations", ctx, userID).
			Return([]models.PhysicalLocation{updatedLocation}, nil)
		mockCache.On("SetCachedPhysicalLocations", ctx, userID, []models.PhysicalLocation{updatedLocation}).
			Return(nil)

		result, err := service.UpdatePhysicalLocation(ctx, userID, updatedLocation)
		if err != nil {
			t.Errorf("UpdatePhysicalLocation() error = %v", err)
			return
		}
		if !reflect.DeepEqual(result, updatedLocation) {
			t.Errorf("UpdatePhysicalLocation() result = %v, expected %v", result, updatedLocation)
		}

		mockDb.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	// Test DeletePhysicalLocation - Success Path
	t.Run("DeletePhysicalLocation - Success Path", func(t *testing.T) {
		// Create fresh mocks for this test
		mockDb := new(mocks.MockPhysicalDbAdapter)
		mockCache := new(mocks.MockPhysicalCacheWrapper)
		service := createTestService(mockDb, mockCache)

		mockDb.On("RemovePhysicalLocation", ctx, "user1", "loc1").
			Return(nil)
		mockDb.On("GetUserPhysicalLocations", ctx, "user1").
			Return([]models.PhysicalLocation{}, nil)
		mockCache.On("SetCachedPhysicalLocations", ctx, "user1", []models.PhysicalLocation{}).
			Return(nil)
		mockCache.On("InvalidateLocationCache", ctx, "user1", "loc1").
			Return(nil)

		err := service.DeletePhysicalLocation(ctx, "user1", "loc1")
		if err != nil {
			t.Errorf("DeletePhysicalLocation() error = %v", err)
			return
		}

		mockDb.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	// Test DeletePhysicalLocation with cache update error
	t.Run("DeletePhysicalLocation - Cache Error Handling", func(t *testing.T) {
		// Create fresh mocks for this test
		mockDb := new(mocks.MockPhysicalDbAdapter)
		mockCache := new(mocks.MockPhysicalCacheWrapper)
		service := createTestService(mockDb, mockCache)

		mockDb.On("RemovePhysicalLocation", ctx, "user1", "loc1").
			Return(nil)
		mockDb.On("GetUserPhysicalLocations", ctx, "user1").
			Return([]models.PhysicalLocation{}, nil)

		// Simulate cache update error - NOTE: should not block deletion
		cacheErr := fmt.Errorf("cache error")
		mockCache.On("SetCachedPhysicalLocations", ctx, "user1", []models.PhysicalLocation{}).
			Return(cacheErr)
		mockCache.On("InvalidateLocationCache", ctx, "user1", "loc1").
			Return(nil)

		err := service.DeletePhysicalLocation(ctx, "user1", "loc1")
		if err != nil {
			t.Errorf("DeletePhysicalLocation() should succeed despite cache error, but got: %v", err)
			return
		}

		mockDb.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})

	// Test DeletePhysicalLocation with db error
	t.Run("DeletePhysicalLocation - DB Error", func(t *testing.T) {
		// Create fresh mocks for this test
		mockDb := new(mocks.MockPhysicalDbAdapter)
		mockCache := new(mocks.MockPhysicalCacheWrapper)
		service := createTestService(mockDb, mockCache)

		dbErr := fmt.Errorf("database error")
		mockDb.On("RemovePhysicalLocation", ctx, "user1", "loc1").
			Return(dbErr)

		err := service.DeletePhysicalLocation(ctx, "user1", "loc1")
		if err == nil {
			t.Errorf("DeletePhysicalLocation() should fail on DB error")
			return
		}
		if err.Error() != "database error" {
			t.Errorf("Expected error 'database error', got: %v", err)
		}

		mockDb.AssertExpectations(t)
		// Cache functions should not be called if DB operation fails
	})
}

func TestUpdatePhysicalLocation(t *testing.T) {
	ctx := context.Background()

	// Create a simple function to create the service with mocks for testing
	createTestService := func(mockDb *mocks.MockPhysicalDbAdapter, mockCache *mocks.MockPhysicalCacheWrapper) *GamePhysicalService {
		mockLogger := testutils.NewTestLogger()

		// Create mock validator and sanitizer
		mockValidator := &mocks.MockPhysicalValidator{
			ValidatePhysicalLocationFunc: func(location models.PhysicalLocation) (models.PhysicalLocation, error) {
				// Simply return the location as valid by default
				return location, nil
			},
		}
		mockSanitizer := &mocks.MockSanitizer{
			SanitizeFunc: func(s string) (string, error) {
				return s, nil // Default implementation just returns the input
			},
		}

		// Create a minimal service without calling the constructor
		return &GamePhysicalService{
			dbAdapter:    mockDb,
			cacheWrapper: mockCache,
			logger:       mockLogger,
			validator:    mockValidator,
			sanitizer:    mockSanitizer,
		}
	}

	tests := []struct {
		name            string
		userID          string
		location        models.PhysicalLocation
		mockSetup       func(*mocks.MockPhysicalDbAdapter, *mocks.MockPhysicalCacheWrapper)
		expectedError   error
		expectedLocation models.PhysicalLocation
	}{
		{
			name:   "successful update",
			userID: "user1",
			location: models.PhysicalLocation{
				ID:             "loc1",
				UserID:         "user1",
				Name:           "Updated Location",
				Label:          "Updated Label",
				LocationType:   "type1",
				MapCoordinates: "1.0,2.0",
			},
			mockSetup: func(mockDb *mocks.MockPhysicalDbAdapter, mockCache *mocks.MockPhysicalCacheWrapper) {
				updatedLocation := models.PhysicalLocation{
					ID:             "loc1",
					UserID:         "user1",
					Name:           "Updated Location",
					Label:          "Updated Label",
					LocationType:   "type1",
					MapCoordinates: "1.0,2.0",
				}

				mockDb.On("UpdatePhysicalLocation", ctx, "user1", mock.Anything).
					Return(updatedLocation, nil)
				mockCache.On("SetSingleCachedPhysicalLocation", ctx, "user1", updatedLocation).
					Return(nil)
				mockCache.On("InvalidateUserCache", ctx, "user1").
					Return(nil)

				updatedLocations := []models.PhysicalLocation{updatedLocation}
				mockDb.On("GetUserPhysicalLocations", ctx, "user1").
					Return(updatedLocations, nil)
				mockCache.On("SetCachedPhysicalLocations", ctx, "user1", updatedLocations).
					Return(nil)
			},
			expectedError: nil,
			expectedLocation: models.PhysicalLocation{
				ID:             "loc1",
				UserID:         "user1",
				Name:           "Updated Location",
				Label:          "Updated Label",
				LocationType:   "type1",
				MapCoordinates: "1.0,2.0",
			},
		},
		{
			name:   "unauthorized update",
			userID: "user2",
			location: models.PhysicalLocation{
				ID:             "loc1",
				UserID:         "user1", // Different user ID than the one making the request
				Name:           "Updated Location",
				Label:          "Updated Label",
				LocationType:   "type1",
				MapCoordinates: "1.0,2.0",
			},
			mockSetup: func(mockDb *mocks.MockPhysicalDbAdapter, mockCache *mocks.MockPhysicalCacheWrapper) {
				// NOTE: no mocks needed since we expect an early return due to user ID mismatch
			},
			expectedError:     ErrUnauthorizedLocation,
			expectedLocation: models.PhysicalLocation{},
		},
		{
			name:   "database error",
			userID: "user1",
			location: models.PhysicalLocation{
				ID:             "loc1",
				UserID:         "user1",
				Name:           "Updated Location",
				Label:          "Updated Label",
				LocationType:   "type1",
				MapCoordinates: "1.0,2.0",
			},
			mockSetup: func(mockDb *mocks.MockPhysicalDbAdapter, mockCache *mocks.MockPhysicalCacheWrapper) {
				dbErr := fmt.Errorf("database error")
				mockDb.On("UpdatePhysicalLocation", ctx, "user1", mock.Anything).
					Return(models.PhysicalLocation{}, dbErr)
			},
			expectedError:     fmt.Errorf("failed to update physical location: database error"),
			expectedLocation: models.PhysicalLocation{},
		},
		{
			name:   "cache error handling",
			userID: "user1",
			location: models.PhysicalLocation{
				ID:             "loc1",
				UserID:         "user1",
				Name:           "Updated Location",
				Label:          "Updated Label",
				LocationType:   "type1",
				MapCoordinates: "1.0,2.0",
			},
			mockSetup: func(mockDb *mocks.MockPhysicalDbAdapter, mockCache *mocks.MockPhysicalCacheWrapper) {
				updatedLocation := models.PhysicalLocation{
					ID:             "loc1",
					UserID:         "user1",
					Name:           "Updated Location",
					Label:          "Updated Label",
					LocationType:   "type1",
					MapCoordinates: "1.0,2.0",
				}

				mockDb.On("UpdatePhysicalLocation", ctx, "user1", mock.Anything).
					Return(updatedLocation, nil)

				// Simulate cache errors - these should not block the update
				cacheErr := fmt.Errorf("cache error")
				mockCache.On("SetSingleCachedPhysicalLocation", ctx, "user1", updatedLocation).
					Return(cacheErr)
				mockCache.On("InvalidateUserCache", ctx, "user1").
					Return(cacheErr)

				updatedLocations := []models.PhysicalLocation{updatedLocation}
				mockDb.On("GetUserPhysicalLocations", ctx, "user1").
					Return(updatedLocations, nil)
				mockCache.On("SetCachedPhysicalLocations", ctx, "user1", updatedLocations).
					Return(cacheErr)
			},
			expectedError: nil, // Should still succeed despite cache errors
			expectedLocation: models.PhysicalLocation{
				ID:             "loc1",
				UserID:         "user1",
				Name:           "Updated Location",
				Label:          "Updated Label",
				LocationType:   "type1",
				MapCoordinates: "1.0,2.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(mocks.MockPhysicalDbAdapter)
			mockCache := new(mocks.MockPhysicalCacheWrapper)

			// Create the service using our helper function
			service := createTestService(mockDb, mockCache)

			// Setup expectations
			if tt.mockSetup != nil {
				tt.mockSetup(mockDb, mockCache)
			}

			location, err := service.UpdatePhysicalLocation(ctx, tt.userID, tt.location)

			// For the error case with a specific message
			if tt.expectedError != nil && err != nil {
				if tt.expectedError.Error() == "failed to update physical location: database error" {
					if err.Error() != tt.expectedError.Error() {
						t.Errorf("UpdatePhysicalLocation() error = %v, expectedError %v", err, tt.expectedError)
					}
				} else if err != tt.expectedError && err.Error() != tt.expectedError.Error() {
					t.Errorf("UpdatePhysicalLocation() error = %v, expectedError %v", err, tt.expectedError)
				}
			} else if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("UpdatePhysicalLocation() error = %v, expectedError %v", err, tt.expectedError)
			}

			if err == nil && !reflect.DeepEqual(location, tt.expectedLocation) {
				t.Errorf("UpdatePhysicalLocation() location = %v, expectedLocation %v", location, tt.expectedLocation)
			}

			mockDb.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}