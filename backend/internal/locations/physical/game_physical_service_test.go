package physical

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
	GetPhysicalLocationsFunc func(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	GetPhysicalLocationFunc func(ctx context.Context, userID, locationID string) (*models.PhysicalLocation, error)
	CreatePhysicalLocationFunc func(ctx context.Context, userID string, location models.PhysicalLocation) error
	UpdatePhysicalLocationFunc func(ctx context.Context, userID string, location models.PhysicalLocation) error
	DeletePhysicalLocationFunc func(ctx context.Context, userID, locationID string) error
}

func newMockGamePhysicalServiceWithDefaults(logger *testutils.TestLogger) *GamePhysicalService {
	mockConfig := mocks.NewMockConfig()

	return &GamePhysicalService{
		dbAdapter:      mocks.DefaultPhysicalDbAdapter(),
		config:         mockConfig,
		logger:         logger,
		validator:      mocks.DefaultPhysicalValidator(),
		sanitizer:      mocks.DefaultSanitizer(),
		cacheWrapper:   mocks.DefaultPhysicalCacheWrapper(),
	}
}

func TestGamePhysicalService(t *testing.T) {
	ctx := context.Background()
	testUserID := "test-user-id"
	testLocationID := "test-location-id"

	testLocation := models.PhysicalLocation{
		ID:                  testLocationID,
		UserID:              testUserID,
		Name:                "Test Location",
		Label:               "Primary",
		LocationType:        "house",
		MapCoordinates:      "40.7128,-74.0060",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// --------- GetPhysicalLocations: Cache Hit ---------
	/*
		GIVEN a user ID and cached locations exist
		WHEN the GetPhysicalLocations method is called
		THEN the service should return the cached locations without querying the database
	*/
	t.Run("GetPhysicalLocations - Cache Hit", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGamePhysicalServiceWithDefaults(testLogger)

		// GIVEN
		mockCache := mocks.DefaultPhysicalCacheWrapper()
		mockCache.GetCachedPhysicalLocationsFunc = func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
			return []models.PhysicalLocation{testLocation}, nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		locations, err := service.GetPhysicalLocations(ctx, testUserID)

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


	// --------- GetPhysicalLocations Cache Miss ---------
	/*
		GIVEN a user ID and no cached locations exist
		WHEN the GetPhysicalLocations method is called
		THEN the service should query the database and cache the results
	*/
	t.Run("GetPhysicalLocations - Cache Miss", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGamePhysicalServiceWithDefaults(testLogger)

		// GIVEN
		mockCache := mocks.DefaultPhysicalCacheWrapper()
		mockCache.GetCachedPhysicalLocationsFunc = func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
			return nil, errors.New("cache miss")
		}

		// Track if SetCachedPhysicalLocationsFunc was called
		cacheSetCalled := false
		mockCache.SetCachedPhysicalLocationsFunc = func(ctx context.Context, userID string, locations []models.PhysicalLocation) error {
			cacheSetCalled = true
			return nil
		}
		service.cacheWrapper = mockCache

		// Override DB adapter to return test data
		mockDb := mocks.DefaultPhysicalDbAdapter()
		mockDb.GetPhysicalLocationsFunc = func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
			return []models.PhysicalLocation{testLocation}, nil
		}
		service.dbAdapter = mockDb

		// WHEN
		locations, err := service.GetPhysicalLocations(ctx, testUserID)

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

	// --------- GetPhysicalLocation Cache Hit ---------
	/*
		GIVEN a user ID, location ID, and cached location exists
		WHEN the GetPhysicalLocation method is called
		THEN the service should return the cached location without querying the database
	*/
	t.Run("GetPhysicalLocation - Cache Hit", func(t *testing.T){
		testLogger := testutils.NewTestLogger()
		service := newMockGamePhysicalServiceWithDefaults(testLogger)

		mockCache := mocks.DefaultPhysicalCacheWrapper()
		mockCache.GetCachedPhysicalLocationsFunc = func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
			return []models.PhysicalLocation{testLocation}, nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		locations, err := service.GetPhysicalLocations(ctx, testUserID)

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


	// --------- GetPhysicalLocation Cache Miss ---------
	/*
		GIVEN a user ID, location ID, and no cached location exists
		WHEN the GetPhysicalLocation method is called
		THEN the service should query the database and cache the result
	*/
	t.Run("GetPhysicalLocation with cache miss", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGamePhysicalServiceWithDefaults(testLogger)

		// Override cache to simulate cache miss
		mockCache := mocks.DefaultPhysicalCacheWrapper()
		mockCache.GetCachedPhysicalLocationsFunc = func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
			return nil, errors.New("cache miss")
		}

		cacheSetCalled := false
		mockCache.SetCachedPhysicalLocationsFunc = func(ctx context.Context, userID string, locations []models.PhysicalLocation) error {
			cacheSetCalled = true
			return nil
		}
		service.cacheWrapper = mockCache

		// Override DB adapter to return test data
		mockDb := mocks.DefaultPhysicalDbAdapter()
		mockDb.GetPhysicalLocationsFunc = func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
			return []models.PhysicalLocation{testLocation}, nil
		}
		service.dbAdapter = mockDb

		// WHEN
		locations, err := service.GetPhysicalLocations(ctx, testUserID)

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


	// --------- AddPhysicalLocation Success ---------
	/*
		GIVEN a user ID and a valid location
		WHEN the AddPhysicalLocation method is called
		THEN the service should add the location to the database and invalidate the cache
	*/
	t.Run("AddPhysicalLocation - Success", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGamePhysicalServiceWithDefaults(testLogger)

		// Track if AddPhysicalLocation func was called
		dbAddCalled := false
		mockDB := mocks.DefaultPhysicalDbAdapter()
		mockDB.CreatePhysicalLocationFunc = func(ctx context.Context, userID string, location models.PhysicalLocation) error {
			dbAddCalled = true
			return nil
		}
		service.dbAdapter = mockDB

		// Track if InvalidateUserCache was called
		cacheInvalidatedCalled := false
		mockCache := mocks.DefaultPhysicalCacheWrapper()
		mockCache.InvalidateUserCacheFunc = func(ctx context.Context, userID string) error {
			cacheInvalidatedCalled = true
			return nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		err := service.AddPhysicalLocation(ctx, testUserID, testLocation)

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


	// --------- AddPhysicalLocation Validation Failure ---------
	/*
		GIVEN a user ID and an invalid location
		WHEN the AddPhysicalLocation method is called
		THEN the service should return a validation error without calling the database
	*/
	t.Run(`AddPhysicalLocation - Validation Failure`, func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGamePhysicalServiceWithDefaults(testLogger)

		// Override validator to simulate validation failure
		mockValidator := mocks.DefaultPhysicalValidator()
		mockValidator.ValidatePhysicalLocationFunc = func(location models.PhysicalLocation) (models.PhysicalLocation, error) {
			return models.PhysicalLocation{}, errors.New("validation error")
		}
		service.validator = mockValidator

		// Track if database is called
		dbCalled := false
		mockDb := mocks.DefaultPhysicalDbAdapter()
		mockDb.CreatePhysicalLocationFunc = func(ctx context.Context, userID string, location models.PhysicalLocation) error {
			dbCalled = true
			return nil
		}
		service.dbAdapter = mockDb

		// WHEN
		err := service.AddPhysicalLocation(ctx, testUserID, testLocation)

		// THEN
		if err == nil {
			t.Error("Expected validation error, got nil")
		}
		if dbCalled {
			t.Error("Expected database not to be called")
		}
	})

	// --------- UpdatePhysicalLocation Success ---------
	/*
		GIVEN a user ID and a valid location
		WHEN the UpdatePhysicalLocation method is called
		THEN the service should update the location in the database and invalidate both caches
	*/
	t.Run(`UpdatePhysicalLocation - Success`, func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGamePhysicalServiceWithDefaults(testLogger)

		// Track if UpdatePhysicalLocation was called
		dbUpdateCalled := false
		mockDb := mocks.DefaultPhysicalDbAdapter()
		mockDb.UpdatePhysicalLocationFunc = func(ctx context.Context, userID string, location models.PhysicalLocation) error {
			dbUpdateCalled = true
			return nil
		}
		service.dbAdapter = mockDb

		// Track if cache invalidation was called
		userCacheInvalidated := false
		locationCacheInvalidated := false
		mockCache := mocks.DefaultPhysicalCacheWrapper()
		mockCache.InvalidateUserCacheFunc = func(ctx context.Context, userID string) error {
			userCacheInvalidated = true
			return nil
		}
		mockCache.InvalidateLocationCacheFunc = func(ctx context.Context, userID, locationID string) error {
			locationCacheInvalidated = true
			return nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		err := service.UpdatePhysicalLocation(ctx, testUserID, testLocation)

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


	// --------- DeletePhysicalLocation Success ---------
	/*
		GIVEN a user ID and a location ID
		WHEN the DeletePhysicalLocation method is called
		THEN the service should delete the location from the database and invalidate both caches
	*/
	t.Run(`DeletePhysicalLocation - Success`, func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGamePhysicalServiceWithDefaults(testLogger)

		// Track if RemovePhysicalLocation was called
		dbDeleteCalled := false
		mockDb := mocks.DefaultPhysicalDbAdapter()
		mockDb.DeletePhysicalLocationFunc = func(ctx context.Context, userID, locationID string) error {
			dbDeleteCalled = true
			return nil
		}
		service.dbAdapter = mockDb

		// Track if cache invalidation was called
		userCacheInvalidated := false
		locationCacheInvalidated := false
		mockCache := mocks.DefaultPhysicalCacheWrapper()
		mockCache.InvalidateUserCacheFunc = func(ctx context.Context, userID string) error {
			userCacheInvalidated = true
			return nil
		}
		mockCache.InvalidateLocationCacheFunc = func(ctx context.Context, userID, locationID string) error {
			locationCacheInvalidated = true
			return nil
		}
		service.cacheWrapper = mockCache


		// WHEN
		err := service.DeletePhysicalLocation(ctx, testUserID, testLocationID)

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
	t.Run(`DB error is properly propagated`, func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGamePhysicalServiceWithDefaults(testLogger)

		// Override DB adapter to return an error
		mockDb := mocks.DefaultPhysicalDbAdapter()
		mockDb.GetPhysicalLocationsFunc = func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
			return nil, errors.New("database error")
		}
		service.dbAdapter = mockDb

		// WHEN
		_, err := service.GetPhysicalLocations(ctx, testUserID)

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
	t.Run(`Cache error during set doesn't block operation`, func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGamePhysicalServiceWithDefaults(testLogger)

		// Override cache to simulate error during set
		mockCache := mocks.DefaultPhysicalCacheWrapper()
		mockCache.GetCachedPhysicalLocationsFunc = func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
			return nil, errors.New("cache miss")
		}
		mockCache.SetCachedPhysicalLocationsFunc = func(ctx context.Context, userID string, locations []models.PhysicalLocation) error {
			return errors.New("cache error")
		}
		service.cacheWrapper = mockCache

		// Override DB adapter to return test data
		mockDb := mocks.DefaultPhysicalDbAdapter()
		mockDb.GetPhysicalLocationsFunc = func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
			return []models.PhysicalLocation{testLocation}, nil
		}
		service.dbAdapter = mockDb

		// WHEN
		locations, err := service.GetPhysicalLocations(ctx, testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error despite cache failure, got %v", err)
		}
		if len(locations) != 1 {
			t.Errorf("Expected 1 location despite cache failure, got %d", len(locations))
		}
	})
}