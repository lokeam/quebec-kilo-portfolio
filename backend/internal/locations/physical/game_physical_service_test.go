package physical

import (
	"context"
	"reflect"
	"testing"

	"github.com/lokeam/qko-beta/internal/appcontext"
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
	mockDb := new(mocks.MockPhysicalDbAdapter)
	mockCache := new(mocks.MockPhysicalCacheWrapper)
	mockLogger := testutils.NewTestLogger()

	appCtx := &appcontext.AppContext{
		Logger: mockLogger,
	}

	service, err := NewGamePhysicalService(appCtx)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	service.dbAdapter = mockDb
	service.cacheWrapper = mockCache

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

		mockCache.On("GetSingleCachedPhysicalLocation", ctx, "user1", "loc1").
			Return(nil, false, nil)
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

		mockCache.On("GetCachedPhysicalLocations", ctx, "user1").
			Return(nil, nil)
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
		mockCache.On("InvalidateUserCache", ctx, "user1").
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

	// Test DeletePhysicalLocation
	t.Run("DeletePhysicalLocation", func(t *testing.T) {
		mockDb.On("RemovePhysicalLocation", ctx, "user1", "loc1").
			Return(nil)
		mockCache.On("InvalidateUserCache", ctx, "user1").
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
}

func TestUpdatePhysicalLocation(t *testing.T) {
	tests := []struct {
		name             string
		userID          string
		location        models.PhysicalLocation
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
				UserID:         "user1",
				Name:           "Updated Location",
				Label:          "Updated Label",
				LocationType:   "type1",
				MapCoordinates: "1.0,2.0",
			},
			expectedError:     ErrUnauthorizedLocation,
			expectedLocation: models.PhysicalLocation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockDb := new(mocks.MockPhysicalDbAdapter)
			mockCache := new(mocks.MockPhysicalCacheWrapper)
			mockLogger := testutils.NewTestLogger()

			appCtx := &appcontext.AppContext{
				Logger: mockLogger,
			}

			service, err := NewGamePhysicalService(appCtx)
			if err != nil {
				t.Fatalf("Failed to create service: %v", err)
			}
			service.dbAdapter = mockDb
			service.cacheWrapper = mockCache

			// Setup expectations
			if tt.expectedError == nil {
				mockDb.On("UpdatePhysicalLocation", ctx, tt.userID, tt.location).
					Return(tt.expectedLocation, nil)
				mockCache.On("SetSingleCachedPhysicalLocation", ctx, tt.userID, tt.expectedLocation).
					Return(nil)
			}

			location, err := service.UpdatePhysicalLocation(ctx, tt.userID, tt.location)
			if err != tt.expectedError {
				t.Errorf("UpdatePhysicalLocation() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if err == nil && !reflect.DeepEqual(location, tt.expectedLocation) {
				t.Errorf("UpdatePhysicalLocation() location = %v, expectedLocation %v", location, tt.expectedLocation)
			}

			mockDb.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}