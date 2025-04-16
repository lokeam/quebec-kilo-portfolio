package sublocation

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
	- GetSublocations retrieves sublocations for a user
		- Attempts to retrieve from cache first
		- Falls back to database on cache miss
		- Caches database results for future requests

	- AddSublocation adds a new sublocation
		- Validates the sublocation data
		- Adds the sublocation to the database
		- Invalidates the user's cache

	- UpdateSublocation updates an existing sublocation
		- Validates the sublocation data
		- Updates the sublocation in the database
		- Invalidates both the user's cache and the specific sublocation cache

	- DeleteSublocation deletes a sublocation
		- Deletes the sublocation from the database
		- Invalidates both the user's cache and the specific sublocation cache

	Scenarios:
		- GetSublocations:
			- Cache hit
			- Cache miss
			- Database error
			- No locations found

		- AddSublocation:
			- validation success + db success
			- validation failure
			- validation success + db failure
			- cache invalidation failure (should not block success)

		- UpdateSublocation:
			- validation successs + db success
			- validation failure
			- validation success + db failure
			- cache invalidation failure (should not block success)

		- DeleteSublocation:
			- db success
			- db failure
			- cache invalidation failure (should not block success)
*/

type MockSublocationDbAdapter struct {
	GetSublocationsFunc func(ctx context.Context, userID string) ([]models.Sublocation, error)
	AddSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) error
	UpdateSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) error
	DeleteSublocationFunc func(ctx context.Context, userID string, sublocationID string) error
}

func newMockGameSublocationServiceWithDefaults(logger *testutils.TestLogger) *GameSublocationService {
	mockConfig := mocks.NewMockConfig()

	return &GameSublocationService{
		dbAdapter:      mocks.DefaultSublocationDbAdapter(),
		config:         mockConfig,
		logger:         logger,
		validator:      mocks.DefaultSublocationValidator(),
		sanitizer:      mocks.DefaultSanitizer(),
		cacheWrapper:   mocks.DefaultSublocationCacheWrapper(),
	}
}

func TestGameSublocationService(t *testing.T) {
	ctx := context.Background()
	testUserID := "test-user-id"
	testSublocationID := "test-sublocation-id"

	testSublocation := models.Sublocation{
		ID:           testSublocationID,
		UserID:       testUserID,
		Name:         "Test Sublocation",
		LocationType: "shelf",
		BgColor:      "blue",
		StoredItems:  50,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// ---------- GetSublocations - Cache Hit ----------
	/*
		GIVEN a user ID and cached sublocations exist
		WHEN the GetSublocations method is called
		THEN the service should return the cached sublocations without querying the database
	*/
	t.Run("GetSublocations() - Cache Hit", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameSublocationServiceWithDefaults(testLogger)

		// GIVEN
		mockCache := mocks.DefaultSublocationCacheWrapper()
		mockCache.GetCachedSublocationsFunc = func(
			ctx context.Context,
			userID string,
		) ([]models.Sublocation, error) {
			return []models.Sublocation{testSublocation}, nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		sublocations, err := service.GetSublocations(ctx, testUserID)

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
	})


	// --------- GetSublocations: Cache Miss ---------
	/*
		GIVEN a user ID and no cached sublocations exist
		WHEN the GetSublocations method is called
		THEN the service should query the database and cache the results
	*/
	t.Run("GetSublocation - Cache Miss", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameSublocationServiceWithDefaults(testLogger)

		// GIVEN
		mockCache := mocks.DefaultSublocationCacheWrapper()
		mockCache.GetSingleCachedSublocationFunc = func(ctx context.Context, userID, sublocationID string) (*models.Sublocation, bool, error) {
			return nil, false, nil
		}
		service.cacheWrapper = mockCache

		// Override DB adapter to return test data
		mockDb := mocks.DefaultSublocationDbAdapter()
		mockDb.GetSublocationFunc = func(ctx context.Context, userID, sublocationID string) (models.Sublocation, error) {
			return testSublocation, nil
		}
		service.dbAdapter = mockDb

		// WHEN
		sublocation, err := service.GetSublocation(ctx, testUserID, testSublocationID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if sublocation.ID != testSublocationID {
			t.Errorf("Expected sublocation ID %s, got %s", testSublocationID, sublocation.ID)
		}
	})


	// --------- AddSublocation: Success ---------
	/*
		GIVEN a user ID and a valid sublocation
		WHEN the AddSublocation method is called
		THEN the service should add the sublocation to the database and invalidate the cache
	*/
	t.Run(`AddSublocation - Success`, func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameSublocationServiceWithDefaults(testLogger)

		// Track if AddSublocation func was called
		dbAddCalled := false
		mockDB := mocks.DefaultSublocationDbAdapter()
		mockDB.AddSublocationFunc = func(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error) {
			dbAddCalled = true
			return sublocation, nil
		}
		service.dbAdapter = mockDB

		// Track if InvalidateUserCache was called
		cacheInvalidatedCalled := false
		mockCache := mocks.DefaultSublocationCacheWrapper()
		mockCache.InvalidateUserCacheFunc = func(ctx context.Context, userID string) error {
			cacheInvalidatedCalled = true
			return nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		createdSublocation, err := service.AddSublocation(ctx, testUserID, testSublocation)

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
		if createdSublocation.ID != testSublocationID {
			t.Errorf("Expected created sublocation ID %s, got %s", testSublocationID, createdSublocation.ID)
		}
	})


	// --------- AddSublocation: Validation Failure ---------
	/*
		GIVEN a user ID and an invalid sublocation
		WHEN the AddSublocation method is called
		THEN the service should return a validation error without calling the database
	*/
	t.Run("AddSublocation - Validation Failure", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameSublocationServiceWithDefaults(testLogger)

		// Override validator to simulate validation failure
		mockValidator := mocks.DefaultSublocationValidator()
		mockValidator.ValidateSublocationFunc = func(sublocation models.Sublocation) (models.Sublocation, error) {
			return models.Sublocation{}, errors.New("validation error")
		}
		service.validator = mockValidator

		// Track if database is called
		dbCalled := false
		mockDb := mocks.DefaultSublocationDbAdapter()
		mockDb.AddSublocationFunc = func(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error) {
			dbCalled = true
			return sublocation, nil
		}
		service.dbAdapter = mockDb

		// WHEN
		createdSublocation, err := service.AddSublocation(ctx, testUserID, testSublocation)

		// THEN
		if err == nil {
			t.Error("Expected validation error, got nil")
		}
		if dbCalled {
			t.Error("Expected database not to be called")
		}
		if createdSublocation.ID != "" {
			t.Error("Expected empty sublocation on validation error")
		}
	})


	// --------- UpdateSublocation: Success ---------
	/*
		GIVEN a user ID and a valid sublocation
		WHEN the UpdateSublocation method is called
		THEN the service should update the sublocation in the database and invalidate both caches
	*/
	t.Run("UpdateSublocation - Success", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameSublocationServiceWithDefaults(testLogger)

		// Track if UpdateSublocation was called
		dbUpdateCalled := false
		mockDb := mocks.DefaultSublocationDbAdapter()
		mockDb.UpdateSublocationFunc = func(ctx context.Context, userID string, sublocation models.Sublocation) error {
			dbUpdateCalled = true
			return nil
		}
		service.dbAdapter = mockDb

		// Track if cache invalidation was called
		userCacheInvalidated := false
		sublocationCacheInvalidated := false
		mockCache := mocks.DefaultSublocationCacheWrapper()
		mockCache.InvalidateUserCacheFunc = func(ctx context.Context, userID string) error {
			userCacheInvalidated = true
			return nil
		}
		mockCache.InvalidateSublocationCacheFunc = func(ctx context.Context, userID, sublocationID string) error {
			sublocationCacheInvalidated = true
			return nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		err := service.UpdateSublocation(ctx, testUserID, testSublocation)

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
		if !sublocationCacheInvalidated {
			t.Error("Expected sublocation cache invalidation to be called")
		}
	})


	// --------- DeleteSublocation: Success ---------
	/*
		GIVEN a user ID and a sublocation ID
		WHEN the DeleteSublocation method is called
		THEN the service should delete the sublocation from the database and invalidate both caches
	*/
	t.Run("DeleteSublocation - Success", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service := newMockGameSublocationServiceWithDefaults(testLogger)

		// Track if RemoveSublocation was called
		dbDeleteCalled := false
		mockDb := mocks.DefaultSublocationDbAdapter()
		mockDb.DeleteSublocationFunc = func(ctx context.Context, userID, sublocationID string) error {
			dbDeleteCalled = true
			return nil
		}
		service.dbAdapter = mockDb

		// Track if cache invalidation was called
		userCacheInvalidated := false
		sublocationCacheInvalidated := false
		mockCache := mocks.DefaultSublocationCacheWrapper()
		mockCache.InvalidateUserCacheFunc = func(ctx context.Context, userID string) error {
			userCacheInvalidated = true
			return nil
		}
		mockCache.InvalidateSublocationCacheFunc = func(ctx context.Context, userID, sublocationID string) error {
			sublocationCacheInvalidated = true
			return nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		err := service.DeleteSublocation(ctx, testUserID, testSublocationID)

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
		if !sublocationCacheInvalidated {
			t.Error("Expected sublocation cache invalidation to be called")
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
		service := newMockGameSublocationServiceWithDefaults(testLogger)

		// Override DB adapter to return an error
		mockDb := mocks.DefaultSublocationDbAdapter()
		mockDb.GetSublocationsFunc = func(ctx context.Context, userID string) ([]models.Sublocation, error) {
			return nil, errors.New("database error")
		}
		service.dbAdapter = mockDb

		// WHEN
		_, err := service.GetSublocations(ctx, testUserID)

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
		service := newMockGameSublocationServiceWithDefaults(testLogger)

		// Override cache to simulate error during set
		mockCache := mocks.DefaultSublocationCacheWrapper()
		mockCache.GetCachedSublocationsFunc = func(ctx context.Context, userID string) ([]models.Sublocation, error) {
			return nil, errors.New("cache miss")
		}
		mockCache.SetCachedSublocationsFunc = func(ctx context.Context, userID string, sublocations []models.Sublocation) error {
			return errors.New("cache error")
		}
		service.cacheWrapper = mockCache

		// Override DB adapter to return test data
		mockDb := mocks.DefaultSublocationDbAdapter()
		mockDb.GetSublocationsFunc = func(ctx context.Context, userID string) ([]models.Sublocation, error) {
			return []models.Sublocation{testSublocation}, nil
		}
		service.dbAdapter = mockDb

		// WHEN
		sublocations, err := service.GetSublocations(ctx, testUserID)

		// THEN
		if err != nil {
			t.Errorf("Expected no error despite cache failure, got %v", err)
		}
		if len(sublocations) != 1 {
			t.Errorf("Expected 1 sublocation despite cache failure, got %d", len(sublocations))
		}
	})

}