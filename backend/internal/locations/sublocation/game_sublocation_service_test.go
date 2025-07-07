package sublocation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/stretchr/testify/assert"
)

/*
	Behavior:
	- GetSublocations retrieves sublocations for a user
		- Attempts to retrieve from cache first
		- Falls back to database on cache miss
		- Caches database results for future requests

	- CreateSublocation adds a new sublocation
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

		- CreateSublocation:
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
	GetSublocationFunc func(ctx context.Context, userID string, sublocationID string) (models.Sublocation, error)
	GetUserSublocationsFunc func(ctx context.Context, userID string) ([]models.Sublocation, error)
	AddSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error)
	UpdateSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) error
	RemoveSublocationFunc func(ctx context.Context, userID string, sublocationID string) error
	DeleteSublocationFunc func(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error)
}

func newMockGameSublocationServiceWithDefaults(logger *testutils.TestLogger) (*GameSublocationService, error) {
	mockConfig := mocks.NewMockConfig()

	// Use nil for cache clients since they're not used in these tests
	appCtx := appcontext.NewAppContext(mockConfig, logger, nil, nil)

	// Create a stub physical service
	physicalService := &mocks.MockPhysicalService{}

	service, err := NewGameSublocationService(appCtx, physicalService)
	if err != nil {
		return nil, err
	}

	return service, nil
}

type mockDbAdapter struct {
	GetSublocationFunc func(ctx context.Context, userID string, sublocationID string) (models.Sublocation, error)
	GetUserSublocationsFunc func(ctx context.Context, userID string) ([]models.Sublocation, error)
	AddSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error)
	UpdateSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) error
	DeleteSublocationFunc func(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error)
	CheckGameInAnySublocationFunc func(ctx context.Context, userGameID string) (bool, error)
	CheckGameInSublocationFunc func(ctx context.Context, userGameID string, sublocationID string) (bool, error)
	CheckGameOwnershipFunc func(ctx context.Context, userID string, userGameID string) (bool, error)
	MoveGameToSublocationFunc func(ctx context.Context, userID string, userGameID string, targetSublocationID string) error
	RemoveGameFromSublocationFunc func(ctx context.Context, userID string, userGameID string) error
	CheckDuplicateSublocationFunc func(ctx context.Context, userID string, physicalLocationID string, name string) (bool, error)
}

func (m *mockDbAdapter) GetAllSublocations(ctx context.Context, userID string) ([]models.Sublocation, error) {
	return m.GetUserSublocationsFunc(ctx, userID)
}

func (m *mockDbAdapter) GetSingleSublocation(ctx context.Context, userID, sublocationID string) (models.Sublocation, error) {
	return m.GetSublocationFunc(ctx, userID, sublocationID)
}

func (m *mockDbAdapter) CreateSublocation(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error) {
	return m.AddSublocationFunc(ctx, userID, sublocation)
}

func (m *mockDbAdapter) UpdateSublocation(ctx context.Context, userID string, sublocation models.Sublocation) error {
	return m.UpdateSublocationFunc(ctx, userID, sublocation)
}

func (m *mockDbAdapter) DeleteSublocation(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error) {
	if m.DeleteSublocationFunc != nil {
		return m.DeleteSublocationFunc(ctx, userID, sublocationIDs)
	}
	return types.DeleteSublocationResponse{
		Success: true,
		DeletedCount: len(sublocationIDs),
		SublocationIDs: sublocationIDs,
	}, nil
}

func (m *mockDbAdapter) CheckGameInAnySublocation(ctx context.Context, userGameID string) (bool, error) {
	if m.CheckGameInAnySublocationFunc != nil {
		return m.CheckGameInAnySublocationFunc(ctx, userGameID)
	}
	return false, nil
}

func (m *mockDbAdapter) CheckGameInSublocation(ctx context.Context, userGameID string, sublocationID string) (bool, error) {
	if m.CheckGameInSublocationFunc != nil {
		return m.CheckGameInSublocationFunc(ctx, userGameID, sublocationID)
	}
	return false, nil
}

func (m *mockDbAdapter) CheckGameOwnership(ctx context.Context, userID string, userGameID string) (bool, error) {
	if m.CheckGameOwnershipFunc != nil {
		return m.CheckGameOwnershipFunc(ctx, userID, userGameID)
	}
	return true, nil
}

func (m *mockDbAdapter) MoveGameToSublocation(ctx context.Context, userID string, userGameID string, targetSublocationID string) error {
	if m.MoveGameToSublocationFunc != nil {
		return m.MoveGameToSublocationFunc(ctx, userID, userGameID, targetSublocationID)
	}
	return nil
}

func (m *mockDbAdapter) RemoveGameFromSublocation(ctx context.Context, userID string, userGameID string) error {
	if m.RemoveGameFromSublocationFunc != nil {
		return m.RemoveGameFromSublocationFunc(ctx, userID, userGameID)
	}
	return nil
}

func (m *mockDbAdapter) CheckDuplicateSublocation(ctx context.Context, userID string, physicalLocationID string, name string) (bool, error) {
	if m.CheckDuplicateSublocationFunc != nil {
		return m.CheckDuplicateSublocationFunc(ctx, userID, physicalLocationID, name)
	}
	return false, nil
}

func TestGameSublocationService(t *testing.T) {
	ctx := context.Background()
	testUserID := "test-user-id"
	testSublocationID := "test-sublocation-id"

	testSublocation := models.Sublocation{
		ID:                 testSublocationID,
		UserID:            testUserID,
		PhysicalLocationID: "physical-location-1",
		Name:              "Test Sublocation",
		LocationType:      "shelf",
		StoredItems:       50,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// ---------- GetSublocations - Cache Hit ----------
	t.Run("GetSublocations() - Cache Hit", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service, err := newMockGameSublocationServiceWithDefaults(testLogger)
		assert.NoError(t, err)

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
		assert.NoError(t, err)
		assert.Len(t, sublocations, 1)
		assert.Equal(t, testSublocationID, sublocations[0].ID)
	})

	// --------- GetSublocations: Cache Miss ---------
	t.Run("GetSingleSublocation - Cache Miss", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service, err := newMockGameSublocationServiceWithDefaults(testLogger)
		assert.NoError(t, err)

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
		sublocation, err := service.GetSingleSublocation(ctx, testUserID, testSublocationID)

		// THEN
		assert.NoError(t, err)
		assert.Equal(t, testSublocationID, sublocation.ID)
	})

	// --------- CreateSublocation: Success ---------
	t.Run("CreateSublocation - Success", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service, err := newMockGameSublocationServiceWithDefaults(testLogger)
		assert.NoError(t, err)

		// Track if CreateSublocation func was called
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
		req := types.CreateSublocationRequest{
			Name: "Test Sublocation",
			LocationType: "shelf",
			PhysicalLocationID: "physical-location-1",
		}
		createdSublocation, err := service.CreateSublocation(ctx, testUserID, req)

		// THEN
		assert.NoError(t, err)
		assert.True(t, dbAddCalled)
		assert.True(t, cacheInvalidatedCalled)
		assert.Equal(t, req.Name, createdSublocation.Name)
	})

	// --------- CreateSublocation: Validation Failure ---------
	t.Run("CreateSublocation - Validation Failure", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service, err := newMockGameSublocationServiceWithDefaults(testLogger)
		assert.NoError(t, err)

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
		req := types.CreateSublocationRequest{
			Name: "Test Sublocation",
			LocationType: "shelf",
			PhysicalLocationID: "physical-location-1",
		}
		createdSublocation, err := service.CreateSublocation(ctx, testUserID, req)

		// THEN
		assert.Error(t, err)
		assert.False(t, dbCalled)
		assert.Empty(t, createdSublocation.ID)
	})

	// --------- UpdateSublocation: Success ---------
	t.Run("UpdateSublocation - Success", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service, err := newMockGameSublocationServiceWithDefaults(testLogger)
		assert.NoError(t, err)

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
		req := types.UpdateSublocationRequest{
			Name: "Updated Sublocation",
			LocationType: "console",
		}
		err = service.UpdateSublocation(ctx, testUserID, testSublocationID, req)

		// THEN
		assert.NoError(t, err)
		assert.True(t, dbUpdateCalled)
		assert.True(t, userCacheInvalidated)
		assert.True(t, sublocationCacheInvalidated)
	})

	// --------- DeleteSublocation: Success ---------
	t.Run("DeleteSublocation - Success", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service, err := newMockGameSublocationServiceWithDefaults(testLogger)
		assert.NoError(t, err)

		// Setup mock to return a sublocation for GetSingleSublocation call
		mockDb := mocks.DefaultSublocationDbAdapter()
		mockDb.GetSublocationFunc = func(ctx context.Context, userID, sublocationID string) (models.Sublocation, error) {
			return testSublocation, nil
		}

		// Track if DeleteSublocation was called via DeleteSublocationFunc
		dbDeleteCalled := false
		mockDb.DeleteSublocationFunc = func(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error) {
			dbDeleteCalled = true
			return types.DeleteSublocationResponse{
				Success: true,
				DeletedCount: len(sublocationIDs),
				SublocationIDs: sublocationIDs,
			}, nil
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
		mockCache.InvalidateLocationCacheFunc = func(ctx context.Context, userID, locationID string) error {
			return nil
		}
		service.cacheWrapper = mockCache

		// WHEN
		response, err := service.DeleteSublocation(ctx, testUserID, []string{testSublocationID})

		// THEN
		assert.NoError(t, err)
		assert.True(t, dbDeleteCalled)
		assert.True(t, userCacheInvalidated)
		assert.True(t, sublocationCacheInvalidated)
		assert.True(t, response.Success)
		assert.Equal(t, 1, response.DeletedCount)
		assert.Equal(t, []string{testSublocationID}, response.SublocationIDs)
	})

	// --------- DB Error Handling ---------
	t.Run("DB error is properly propagated", func(t *testing.T) {
		testLogger := testutils.NewTestLogger()
		service, err := newMockGameSublocationServiceWithDefaults(testLogger)
		assert.NoError(t, err)

		// Create a custom mock DB adapter
		mock := &mockDbAdapter{
			GetUserSublocationsFunc: func(ctx context.Context, userID string) ([]models.Sublocation, error) {
				return nil, errors.New("database error")
			},
			GetSublocationFunc: func(ctx context.Context, userID, sublocationID string) (models.Sublocation, error) {
				return models.Sublocation{}, nil
			},
			AddSublocationFunc: func(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error) {
				return sublocation, nil
			},
			UpdateSublocationFunc: func(ctx context.Context, userID string, sublocation models.Sublocation) error {
				return nil
			},
			DeleteSublocationFunc: func(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error) {
				return types.DeleteSublocationResponse{
					Success: true,
					DeletedCount: len(sublocationIDs),
					SublocationIDs: sublocationIDs,
				}, nil
			},
			CheckGameInAnySublocationFunc: func(ctx context.Context, userGameID string) (bool, error) {
				return false, nil
			},
			CheckGameInSublocationFunc: func(ctx context.Context, userGameID string, sublocationID string) (bool, error) {
				return false, nil
			},
			CheckGameOwnershipFunc: func(ctx context.Context, userID string, userGameID string) (bool, error) {
				return true, nil
			},
			MoveGameToSublocationFunc: func(ctx context.Context, userID string, userGameID string, targetSublocationID string) error {
				return nil
			},
			RemoveGameFromSublocationFunc: func(ctx context.Context, userID string, userGameID string) error {
				return nil
			},
			CheckDuplicateSublocationFunc: func(ctx context.Context, userID string, physicalLocationID string, name string) (bool, error) {
				return false, nil
			},
		}
		service.dbAdapter = mock

		// WHEN
		_, err = service.GetSublocations(ctx, testUserID)

		// THEN
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestCreateSublocation(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	testUserID := "test-user"

	testLogger := testutils.NewTestLogger()
	service, err := newMockGameSublocationServiceWithDefaults(testLogger)
	assert.NoError(t, err)

	// WHEN
	req := types.CreateSublocationRequest{
		Name: "Test Sublocation",
		LocationType: "shelf",
		PhysicalLocationID: "physical-location-1",
	}
	sublocation, err := service.CreateSublocation(ctx, testUserID, req)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, req.Name, sublocation.Name)
	assert.Equal(t, req.LocationType, sublocation.LocationType)
	assert.Equal(t, req.PhysicalLocationID, sublocation.PhysicalLocationID)
}

func TestUpdateSublocation(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	testUserID := "test-user"
	testSublocationID := "test-sublocation"

	testLogger := testutils.NewTestLogger()
	service, err := newMockGameSublocationServiceWithDefaults(testLogger)
	assert.NoError(t, err)

	// WHEN
	req := types.UpdateSublocationRequest{
		Name: "Updated Sublocation",
		LocationType: "console",
	}
	err = service.UpdateSublocation(ctx, testUserID, testSublocationID, req)

	// THEN
	assert.NoError(t, err)
}

func TestDeleteSublocation(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	testUserID := "test-user"
	testSublocationID := "test-sublocation"

	testLogger := testutils.NewTestLogger()
	service, err := newMockGameSublocationServiceWithDefaults(testLogger)
	assert.NoError(t, err)

	// WHEN
	response, err := service.DeleteSublocation(ctx, testUserID, []string{testSublocationID})

	// THEN
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, 1, response.DeletedCount)
	assert.Equal(t, []string{testSublocationID}, response.SublocationIDs)
}

func TestDeleteSublocationWithOrphanedGames(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	testUserID := "test-user"
	testSublocationID := "test-sublocation"

	testLogger := testutils.NewTestLogger()
	service, err := newMockGameSublocationServiceWithDefaults(testLogger)
	assert.NoError(t, err)

	// Set up mock to return orphaned games
	mockDb := mocks.DefaultSublocationDbAdapter()
	mockDb.DeleteSublocationFunc = func(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error) {
		return types.DeleteSublocationResponse{
			Success: true,
			DeletedCount: len(sublocationIDs),
			SublocationIDs: sublocationIDs,
			DeletedGames: []types.DeletedGameDetails{
				{
					UserGameID: 1,
					GameID: 123,
					GameName: "Test Game",
					PlatformName: "PS5",
				},
			},
		}, nil
	}
	service.dbAdapter = mockDb

	// WHEN
	response, err := service.DeleteSublocation(ctx, testUserID, []string{testSublocationID})

	// THEN
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, 1, response.DeletedCount)
	assert.Equal(t, []string{testSublocationID}, response.SublocationIDs)
	assert.Len(t, response.DeletedGames, 1)
	assert.Equal(t, "Test Game", response.DeletedGames[0].GameName)
}