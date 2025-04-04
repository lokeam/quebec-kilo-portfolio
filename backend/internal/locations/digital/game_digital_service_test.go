package digital

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
)

/*
Behavior:
- GetDigitalLocations retrieves digital locations for a user
  - Attempts to retrieve from cache first
  - Falls back to database on cache miss
  - Caches database results for future requests

- AddDigitalLocation adds a new digital location
  - Validates the location data
  - Adds the location to the database
  - Invalidates the user's cache

- UpdateDigitalLocation updates an existing digital location
  - Validates the location data
  - Updates the location in the database
  - Invalidates both the user's cache and the specific location cache

- DeleteDigitalLocation deletes a digital location
  - Deletes the location from the database
  - Invalidates both the user's cache and the specific location cache

Scenarios:
- GetDigitalLocations:
  - Cache hit
  - Cache miss
  - Database error
  - No locations found

- AddDigitalLocation:
  - Validation success + db success
  - Validation failure
  - Validation success + db failure
  - Cache invalidation failure (should not block success)

- UpdateDigitalLocation:
  - Validation success + db success
  - Validation failure
  - Validation success + db failure
  - Cache invalidation failure (should not block success)

- DeleteDigitalLocation:
  - DB success
  - DB failure
  - Cache invalidation failure (should not block success)
*/

type MockDigitalDbAdapter struct {
	GetDigitalLocationsFunc    func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetDigitalLocationFunc     func(ctx context.Context, userID, locationID string) (*models.DigitalLocation, error)
	CreateDigitalLocationFunc  func(ctx context.Context, userID string, location models.DigitalLocation) error
	UpdateDigitalLocationFunc  func(ctx context.Context, userID string, location models.DigitalLocation) error
	DeleteDigitalLocationFunc  func(ctx context.Context, userID, locationID string) error
}

func newMockGameDigitalServiceWithDefaults(logger *testutils.TestLogger) *GameDigitalService {
	mockConfig := mocks.NewMockConfig()

	return &GameDigitalService{
		dbAdapter:    mocks.DefaultDigitalDbAdapter(),
		config:       mockConfig,
		logger:       logger,
		validator:    mocks.DefaultDigitalValidator(),
		sanitizer:    mocks.DefaultSanitizer(),
		cacheWrapper: mocks.DefaultDigitalCacheWrapper(),
	}
}

func TestGameDigitalService(t *testing.T) {
	ctx := context.Background()
	testUserID := "test-user-id"
	testLocationID := "test-location-id"

	testLocation := models.DigitalLocation{
		ID:        testLocationID,
		UserID:    testUserID,
		Name:      "Test Digital Location",
		URL:       "https://example.com",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// --------- GetDigitalLocations: Cache Hit ---------
	/*
		GIVEN a user ID and cached locations exist
		WHEN the GetDigitalLocations method is called
		THEN the service should return the cached locations without querying the database
	*/
	t.Run("GetDigitalLocations - Cache Hit", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameDigitalServiceWithDefaults(testLogger)

		// GIVEN
		mockCache := mocks.DefaultDigitalCacheWrapper()
		mockCache.GetCachedDigitalLocationsFunc = func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{testLocation}, nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		locations, err := service.GetDigitalLocations(ctx, testUserID)

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
	})


	// --------- GetDigitalLocations Cache Miss ---------
	/*
		GIVEN a user ID and no cached locations exist
		WHEN the GetDigitalLocations method is called
		THEN the service should query the database and cache the results
	*/
	t.Run("GetDigitalLocations - Cache Miss", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameDigitalServiceWithDefaults(testLogger)

		// GIVEN
		mockCache := mocks.DefaultDigitalCacheWrapper()
		mockCache.GetCachedDigitalLocationsFunc = func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return nil, errors.New("cache miss")
		}

		// Track if SetCachedDigitalLocationsFunc was called
		cacheSetCalled := false
		mockCache.SetCachedDigitalLocationsFunc = func(ctx context.Context, userID string, locations []models.DigitalLocation) error {
			cacheSetCalled = true
			return nil
		}
		service.cacheWrapper = mockCache

		// Override DB adapter to return test data
		mockDb := mocks.DefaultDigitalDbAdapter()
		mockDb.GetDigitalLocationsFunc = func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{testLocation}, nil
		}
		service.dbAdapter = mockDb

		// WHEN
		locations, err := service.GetDigitalLocations(ctx, testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != 1 {
			t.Errorf("Expected 1 location, got %d", len(locations))
		}
		if !cacheSetCalled {
			t.Error("Expected cache set to be called")
		}
	})


	// --------- GetDigitalLocation Cache Hit ---------
	/*
		GIVEN a user ID, location ID, and cached location exists
		WHEN the GetDigitalLocation method is called
		THEN the service should return the cached location without querying the database
	*/
	t.Run("GetDigitalLocation - Cache Hit", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameDigitalServiceWithDefaults(testLogger)

		// GIVEN
		mockCache := mocks.DefaultDigitalCacheWrapper()
		mockCache.GetSingleCachedDigitalLocationFunc = func(ctx context.Context, userID, locationID string) (*models.DigitalLocation, bool, error) {
			return &testLocation, true, nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		location, err := service.GetDigitalLocation(ctx, testUserID, testLocationID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if location == nil {
			t.Error("Expected location, got nil")
		} else if location.ID != testLocationID {
			t.Errorf("Expected location ID %s, got %s", testLocationID, location.ID)
		}
	})


	// --------- GetDigitalLocation Cache Miss ---------
	/*
		GIVEN a user ID, location ID, and no cached location exists
		WHEN the GetDigitalLocation method is called
		THEN the service should query the database and cache the result
	*/
	t.Run("GetDigitalLocation - Cache Miss", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameDigitalServiceWithDefaults(testLogger)

		// GIVEN
		mockCache := mocks.DefaultDigitalCacheWrapper()
		mockCache.GetSingleCachedDigitalLocationFunc = func(ctx context.Context, userID, locationID string) (*models.DigitalLocation, bool, error) {
			return nil, false, nil
		}

		// Track if SetSingleCachedDigitalLocationFunc was called
		cacheSetCalled := false
		mockCache.SetSingleCachedDigitalLocationFunc = func(ctx context.Context, userID string, location models.DigitalLocation) error {
			cacheSetCalled = true
			return nil
		}
		service.cacheWrapper = mockCache

		// Override DB adapter to return test data
		mockDb := mocks.DefaultDigitalDbAdapter()
		mockDb.GetDigitalLocationFunc = func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return testLocation, nil
		}
		service.dbAdapter = mockDb

		// WHEN
		location, err := service.GetDigitalLocation(ctx, testUserID, testLocationID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if location == nil {
			t.Error("Expected location, got nil")
		}
		if !cacheSetCalled {
			t.Error("Expected cache set to be called")
		}
	})


	// --------- AddDigitalLocation Success ---------
	/*
		GIVEN a user ID and a valid location
		WHEN the AddDigitalLocation method is called
		THEN the service should add the location to the database and invalidate the cache
	*/
	t.Run("AddDigitalLocation - Success", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameDigitalServiceWithDefaults(testLogger)

		// GIVEN
		// Track if AddDigitalLocation func was called
		dbAddCalled := false
		mockDB := mocks.DefaultDigitalDbAdapter()
		mockDB.AddDigitalLocationFunc = func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
			dbAddCalled = true
			return location, nil
		}
		service.dbAdapter = mockDB

		// Track if InvalidateUserCache was called
		cacheInvalidatedCalled := false
		mockCache := mocks.DefaultDigitalCacheWrapper()
		mockCache.InvalidateUserCacheFunc = func(ctx context.Context, userID string) error {
			cacheInvalidatedCalled = true
			return nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		err := service.AddDigitalLocation(ctx, testUserID, testLocation)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !dbAddCalled {
			t.Error("Expected database add to be called")
		}
		if !cacheInvalidatedCalled {
			t.Error("Expected cache invalidation to be called")
		}
	})


	// --------- AddDigitalLocation Validation Failure ---------
	/*
		GIVEN a user ID and an invalid location
		WHEN the AddDigitalLocation method is called
		THEN the service should return a validation error without calling the database
	*/
	t.Run("AddDigitalLocation - Validation Failure", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameDigitalServiceWithDefaults(testLogger)

		// GIVEN
		// Override validator to simulate validation failure
		mockValidator := mocks.DefaultDigitalValidator()
		mockValidator.ValidateDigitalLocationFunc = func(location models.DigitalLocation) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, errors.New("validation error")
		}
		service.validator = mockValidator

		// Track if database is called
		dbCalled := false
		mockDb := mocks.DefaultDigitalDbAdapter()
		mockDb.AddDigitalLocationFunc = func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
			dbCalled = true
			return location, nil
		}
		service.dbAdapter = mockDb

		// WHEN
		err := service.AddDigitalLocation(ctx, testUserID, testLocation)

		// THEN
		if err == nil {
			t.Error("Expected validation error, got nil")
		}
		if dbCalled {
			t.Error("Expected database not to be called")
		}
	})


	// --------- UpdateDigitalLocation Success ---------
	/*
		GIVEN a user ID and a valid location
		WHEN the UpdateDigitalLocation method is called
		THEN the service should update the location in the database and invalidate both caches
	*/
	t.Run("UpdateDigitalLocation - Success", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameDigitalServiceWithDefaults(testLogger)

		// GIVEN
		// Track if UpdateDigitalLocation was called
		dbUpdateCalled := false
		mockDb := mocks.DefaultDigitalDbAdapter()
		mockDb.UpdateDigitalLocationFunc = func(ctx context.Context, userID string, location models.DigitalLocation) error {
			dbUpdateCalled = true
			return nil
		}
		service.dbAdapter = mockDb

		// Track if cache invalidation was called
		userCacheInvalidated := false
		locationCacheInvalidated := false
		mockCache := mocks.DefaultDigitalCacheWrapper()
		mockCache.InvalidateUserCacheFunc = func(ctx context.Context, userID string) error {
			userCacheInvalidated = true
			return nil
		}
		mockCache.InvalidateDigitalLocationCacheFunc = func(ctx context.Context, userID, locationID string) error {
			locationCacheInvalidated = true
			return nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		err := service.UpdateDigitalLocation(ctx, testUserID, testLocation)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !dbUpdateCalled {
			t.Error("Expected database update to be called")
		}
		if !userCacheInvalidated {
			t.Error("Expected user cache invalidation to be called")
		}
		if !locationCacheInvalidated {
			t.Error("Expected location cache invalidation to be called")
		}
	})


	// --------- DeleteDigitalLocation Success ---------
	/*
		GIVEN a user ID and a location ID
		WHEN the DeleteDigitalLocation method is called
		THEN the service should delete the location from the database and invalidate both caches
	*/
	t.Run("DeleteDigitalLocation - Success", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameDigitalServiceWithDefaults(testLogger)

		// GIVEN
		// Track if RemoveDigitalLocation was called
		dbDeleteCalled := false
		mockDb := mocks.DefaultDigitalDbAdapter()
		mockDb.DeleteDigitalLocationFunc = func(ctx context.Context, userID, locationID string) error {
			dbDeleteCalled = true
			return nil
		}
		service.dbAdapter = mockDb

		// Track if cache invalidation was called
		userCacheInvalidated := false
		locationCacheInvalidated := false
		mockCache := mocks.DefaultDigitalCacheWrapper()
		mockCache.InvalidateUserCacheFunc = func(ctx context.Context, userID string) error {
			userCacheInvalidated = true
			return nil
		}
		mockCache.InvalidateDigitalLocationCacheFunc = func(ctx context.Context, userID, locationID string) error {
			locationCacheInvalidated = true
			return nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		err := service.DeleteDigitalLocation(ctx, testUserID, testLocationID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !dbDeleteCalled {
			t.Error("Expected database delete to be called")
		}
		if !userCacheInvalidated {
			t.Error("Expected user cache invalidation to be called")
		}
		if !locationCacheInvalidated {
			t.Error("Expected location cache invalidation to be called")
		}
	})


	// --------- DB Error Handling ---------
	/*
		GIVEN a user ID
		WHEN the database returns an error
		THEN the service should propagate that error
	*/
	t.Run("DB error is properly propagated", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameDigitalServiceWithDefaults(testLogger)

		// GIVEN
		// Override DB adapter to return an error
		mockDb := mocks.DefaultDigitalDbAdapter()
		mockDb.GetDigitalLocationsFunc = func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return nil, errors.New("database error")
		}
		service.dbAdapter = mockDb

		// WHEN
		_, err := service.GetDigitalLocations(ctx, testUserID)

		// THEN
		if err == nil {
			t.Error("Expected database error, got nil")
		}
	})


	// --------- Cache Error Handling ---------
	/*
		GIVEN a user ID and successful DB operation
		WHEN the cache returns an error during set
		THEN the service should still return success
	*/
	t.Run("Cache error during set doesn't block operation", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameDigitalServiceWithDefaults(testLogger)

		// GIVEN
		// Override cache to simulate error during set
		mockCache := mocks.DefaultDigitalCacheWrapper()
		mockCache.GetCachedDigitalLocationsFunc = func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return nil, errors.New("cache miss")
		}
		mockCache.SetCachedDigitalLocationsFunc = func(ctx context.Context, userID string, locations []models.DigitalLocation) error {
			return errors.New("cache error")
		}
		service.cacheWrapper = mockCache

		// Override DB adapter to return test data
		mockDb := mocks.DefaultDigitalDbAdapter()
		mockDb.GetDigitalLocationsFunc = func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{testLocation}, nil
		}
		service.dbAdapter = mockDb

		// WHEN
		locations, err := service.GetDigitalLocations(ctx, testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error despite cache failure, got %v", err)
		}
		if len(locations) != 1 {
			t.Errorf("Expected 1 location despite cache failure, got %d", len(locations))
		}
	})
}
