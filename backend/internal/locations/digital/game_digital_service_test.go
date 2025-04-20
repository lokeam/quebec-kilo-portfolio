package digital

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/lokeam/qko-beta/internal/appcontext"
	mcache "github.com/lokeam/qko-beta/internal/infrastructure/cache/memorycache"
	rcache "github.com/lokeam/qko-beta/internal/infrastructure/cache/rueidis"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
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
	GetUserDigitalLocationsFunc func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetDigitalLocationFunc func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
	FindDigitalLocationByNameFunc func(ctx context.Context, userID string, name string) (models.DigitalLocation, error)
	AddDigitalLocationFunc func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocationFunc func(ctx context.Context, userID string, location models.DigitalLocation) error
	RemoveDigitalLocationFunc func(ctx context.Context, userID, locationID string) error

	// Subscription Operations
	GetSubscriptionFunc func(ctx context.Context, locationID string) (*models.Subscription, error)
	AddSubscriptionFunc func(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscriptionFunc func(ctx context.Context, subscription models.Subscription) error
	RemoveSubscriptionFunc func(ctx context.Context, locationID string) error

	// Payment Operations
	GetPaymentsFunc func(ctx context.Context, locationID string) ([]models.Payment, error)
	AddPaymentFunc func(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPaymentFunc func(ctx context.Context, paymentID int64) (*models.Payment, error)

	// Game Operations
	AddGameToDigitalLocationFunc func(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocationFunc func(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationIDFunc func(ctx context.Context, userID string, locationID string) ([]models.Game, error)
}

func (m *MockDigitalDbAdapter) GetUserDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
	return m.GetUserDigitalLocationsFunc(ctx, userID)
}

func (m *MockDigitalDbAdapter) GetDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
	return m.GetDigitalLocationFunc(ctx, userID, locationID)
}

func (m *MockDigitalDbAdapter) FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
	return m.FindDigitalLocationByNameFunc(ctx, userID, name)
}

func (m *MockDigitalDbAdapter) AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
	return m.AddDigitalLocationFunc(ctx, userID, location)
}

func (m *MockDigitalDbAdapter) UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error {
	return m.UpdateDigitalLocationFunc(ctx, userID, location)
}

func (m *MockDigitalDbAdapter) RemoveDigitalLocation(ctx context.Context, userID, locationID string) error {
	return m.RemoveDigitalLocationFunc(ctx, userID, locationID)
}

// Subscription Operations
func (m *MockDigitalDbAdapter) GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error) {
	return m.GetSubscriptionFunc(ctx, locationID)
}

func (m *MockDigitalDbAdapter) AddSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
	return m.AddSubscriptionFunc(ctx, subscription)
}

func (m *MockDigitalDbAdapter) UpdateSubscription(ctx context.Context, subscription models.Subscription) error {
	return m.UpdateSubscriptionFunc(ctx, subscription)
}

func (m *MockDigitalDbAdapter) RemoveSubscription(ctx context.Context, locationID string) error {
	return m.RemoveSubscriptionFunc(ctx, locationID)
}

// Payment Operations
func (m *MockDigitalDbAdapter) GetPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	return m.GetPaymentsFunc(ctx, locationID)
}

func (m *MockDigitalDbAdapter) AddPayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	return m.AddPaymentFunc(ctx, payment)
}

func (m *MockDigitalDbAdapter) GetPayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
	return m.GetPaymentFunc(ctx, paymentID)
}

// Game Operations
func (m *MockDigitalDbAdapter) AddGameToDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error {
	return m.AddGameToDigitalLocationFunc(ctx, userID, locationID, gameID)
}

func (m *MockDigitalDbAdapter) RemoveGameFromDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error {
	return m.RemoveGameFromDigitalLocationFunc(ctx, userID, locationID, gameID)
}

func (m *MockDigitalDbAdapter) GetGamesByDigitalLocationID(ctx context.Context, userID string, locationID string) ([]models.Game, error) {
	return m.GetGamesByDigitalLocationIDFunc(ctx, userID, locationID)
}

// MockDigitalCacheWrapper is a mock implementation of interfaces.DigitalCacheWrapper
type MockDigitalCacheWrapper struct {
	GetCachedDigitalLocationsFunc      func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	SetCachedDigitalLocationsFunc      func(ctx context.Context, userID string, locations []models.DigitalLocation) error
	GetSingleCachedDigitalLocationFunc func(ctx context.Context, userID, locationID string) (*models.DigitalLocation, bool, error)
	SetSingleCachedDigitalLocationFunc func(ctx context.Context, userID string, location models.DigitalLocation) error
	InvalidateUserCacheFunc            func(ctx context.Context, userID string) error
	InvalidateDigitalLocationCacheFunc func(ctx context.Context, userID, locationID string) error

	// Subscription caching
	GetCachedSubscriptionFunc          func(ctx context.Context, locationID string) (*models.Subscription, bool, error)
	SetCachedSubscriptionFunc          func(ctx context.Context, locationID string, subscription models.Subscription) error
	InvalidateSubscriptionCacheFunc    func(ctx context.Context, locationID string) error

	// Payment caching
	GetCachedPaymentsFunc              func(ctx context.Context, locationID string) ([]models.Payment, error)
	SetCachedPaymentsFunc              func(ctx context.Context, locationID string, payments []models.Payment) error
	InvalidatePaymentsCacheFunc        func(ctx context.Context, locationID string) error
}

func (m *MockDigitalCacheWrapper) GetCachedDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
	return m.GetCachedDigitalLocationsFunc(ctx, userID)
}

func (m *MockDigitalCacheWrapper) SetCachedDigitalLocations(ctx context.Context, userID string, locations []models.DigitalLocation) error {
	return m.SetCachedDigitalLocationsFunc(ctx, userID, locations)
}

func (m *MockDigitalCacheWrapper) GetSingleCachedDigitalLocation(ctx context.Context, userID, locationID string) (*models.DigitalLocation, bool, error) {
	return m.GetSingleCachedDigitalLocationFunc(ctx, userID, locationID)
}

func (m *MockDigitalCacheWrapper) SetSingleCachedDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error {
	return m.SetSingleCachedDigitalLocationFunc(ctx, userID, location)
}

func (m *MockDigitalCacheWrapper) InvalidateUserCache(ctx context.Context, userID string) error {
	return m.InvalidateUserCacheFunc(ctx, userID)
}

func (m *MockDigitalCacheWrapper) InvalidateDigitalLocationCache(ctx context.Context, userID, locationID string) error {
	return m.InvalidateDigitalLocationCacheFunc(ctx, userID, locationID)
}

// Subscription caching
func (m *MockDigitalCacheWrapper) GetCachedSubscription(ctx context.Context, locationID string) (*models.Subscription, bool, error) {
	return m.GetCachedSubscriptionFunc(ctx, locationID)
}

func (m *MockDigitalCacheWrapper) SetCachedSubscription(ctx context.Context, locationID string, subscription models.Subscription) error {
	return m.SetCachedSubscriptionFunc(ctx, locationID, subscription)
}

func (m *MockDigitalCacheWrapper) InvalidateSubscriptionCache(ctx context.Context, locationID string) error {
	return m.InvalidateSubscriptionCacheFunc(ctx, locationID)
}

// Payment caching
func (m *MockDigitalCacheWrapper) GetCachedPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	return m.GetCachedPaymentsFunc(ctx, locationID)
}

func (m *MockDigitalCacheWrapper) SetCachedPayments(ctx context.Context, locationID string, payments []models.Payment) error {
	return m.SetCachedPaymentsFunc(ctx, locationID, payments)
}

func (m *MockDigitalCacheWrapper) InvalidatePaymentsCache(ctx context.Context, locationID string) error {
	return m.InvalidatePaymentsCacheFunc(ctx, locationID)
}

func newMockGameDigitalServiceWithDefaults(logger *testutils.TestLogger) *GameDigitalService {
	mockDbAdapter := &MockDigitalDbAdapter{
		GetUserDigitalLocationsFunc: func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{}, nil
		},
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		FindDigitalLocationByNameFunc: func(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		AddDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		UpdateDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) error {
			return nil
		},
		RemoveDigitalLocationFunc: func(ctx context.Context, userID, locationID string) error {
			return nil
		},
		GetSubscriptionFunc: func(ctx context.Context, locationID string) (*models.Subscription, error) {
			return &models.Subscription{}, nil
		},
		AddSubscriptionFunc: func(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
			return &models.Subscription{}, nil
		},
		UpdateSubscriptionFunc: func(ctx context.Context, subscription models.Subscription) error {
			return nil
		},
		RemoveSubscriptionFunc: func(ctx context.Context, locationID string) error {
			return nil
		},
		GetPaymentsFunc: func(ctx context.Context, locationID string) ([]models.Payment, error) {
			return []models.Payment{}, nil
		},
		AddPaymentFunc: func(ctx context.Context, payment models.Payment) (*models.Payment, error) {
			return &models.Payment{}, nil
		},
		GetPaymentFunc: func(ctx context.Context, paymentID int64) (*models.Payment, error) {
			return &models.Payment{}, nil
		},
		AddGameToDigitalLocationFunc: func(ctx context.Context, userID string, locationID string, gameID int64) error {
			return nil
		},
		RemoveGameFromDigitalLocationFunc: func(ctx context.Context, userID string, locationID string, gameID int64) error {
			return nil
		},
		GetGamesByDigitalLocationIDFunc: func(ctx context.Context, userID string, locationID string) ([]models.Game, error) {
			return []models.Game{}, nil
		},
	}

	mockCacheWrapper := &MockDigitalCacheWrapper{
		GetCachedDigitalLocationsFunc: func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{}, nil
		},
		SetCachedDigitalLocationsFunc: func(ctx context.Context, userID string, locations []models.DigitalLocation) error {
			return nil
		},
		GetSingleCachedDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (*models.DigitalLocation, bool, error) {
			return &models.DigitalLocation{}, true, nil
		},
		SetSingleCachedDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) error {
			return nil
		},
		InvalidateUserCacheFunc: func(ctx context.Context, userID string) error {
			return nil
		},
		InvalidateDigitalLocationCacheFunc: func(ctx context.Context, userID, locationID string) error {
			return nil
		},

		// Subscription caching
		GetCachedSubscriptionFunc: func(ctx context.Context, locationID string) (*models.Subscription, bool, error) {
			return &models.Subscription{}, true, nil
		},
		SetCachedSubscriptionFunc: func(ctx context.Context, locationID string, subscription models.Subscription) error {
			return nil
		},
		InvalidateSubscriptionCacheFunc: func(ctx context.Context, locationID string) error {
			return nil
		},

		// Payment caching
		GetCachedPaymentsFunc: func(ctx context.Context, locationID string) ([]models.Payment, error) {
			return []models.Payment{}, nil
		},
		SetCachedPaymentsFunc: func(ctx context.Context, locationID string, payments []models.Payment) error {
			return nil
		},
		InvalidatePaymentsCacheFunc: func(ctx context.Context, locationID string) error {
			return nil
		},
	}

	// Create mock Redis client
	mockRedisClient := &rcache.RueidisClient{}

	// Create mock memory cache
	mockMemCache := &mcache.MemoryCache{}

	appCtx := &appcontext.AppContext{
		Logger:      logger,
		RedisClient: mockRedisClient,
		MemCache:    mockMemCache,
	}

	// Set the mock adapters on the service after creation
	gds, _ := NewGameDigitalService(appCtx)
	gds.dbAdapter = mockDbAdapter
	gds.cacheWrapper = mockCacheWrapper

	return gds
}

func TestGameDigitalService(t *testing.T) {
	// Set up test logger
	logger := testutils.NewTestLogger()

	// Test cases
	t.Run("GetUserDigitalLocations - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)

		// Execute
		locations, err := service.GetUserDigitalLocations(context.Background(), "test-user")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != 2 {
			t.Errorf("Expected 2 locations, got %d", len(locations))
		}
	})

	t.Run("GetUserDigitalLocations - Error", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := errors.New("test error")

		service.dbAdapter = &MockDigitalDbAdapter{
			GetUserDigitalLocationsFunc: func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
				return nil, expectedErr
			},
		}

		// Execute
		locations, err := service.GetUserDigitalLocations(context.Background(), "test-user")

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if len(locations) != 0 {
			t.Errorf("Expected empty locations slice, got %v", locations)
		}
	})

	t.Run("GetDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)

		// Execute
		location, err := service.GetDigitalLocation(context.Background(), "test-user", "test-location")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if location.ID == "" {
			t.Error("Expected location to be returned")
		}
	})

	t.Run("GetDigitalLocation - Not Found", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := sql.ErrNoRows

		service.dbAdapter = &MockDigitalDbAdapter{
			GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
				return models.DigitalLocation{}, expectedErr
			},
		}

		// Execute
		location, err := service.GetDigitalLocation(context.Background(), "test-user", "test-location")

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if location.ID != "" || location.Name != "" || location.IsActive != false {
			t.Errorf("Expected empty location, got %v", location)
		}
	})

	t.Run("AddDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		newLocation := models.DigitalLocation{
			Name:     "Test Location",
			IsActive: true,
		}

		// Execute
		createdLocation, err := service.AddDigitalLocation(context.Background(), "test-user", newLocation)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if createdLocation.Name != newLocation.Name {
			t.Errorf("Expected location name %s, got %s", newLocation.Name, createdLocation.Name)
		}
	})

	t.Run("AddDigitalLocation - Error", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := errors.New("test error")
		newLocation := models.DigitalLocation{
			Name:     "Test Location",
			IsActive: true,
		}

		service.dbAdapter = &MockDigitalDbAdapter{
			AddDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
				return models.DigitalLocation{}, expectedErr
			},
		}

		// Execute
		createdLocation, err := service.AddDigitalLocation(context.Background(), "test-user", newLocation)

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if createdLocation.ID != "" || createdLocation.Name != "" || createdLocation.IsActive != false {
			t.Errorf("Expected empty location, got %v", createdLocation)
		}
	})

	t.Run("UpdateDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		location := models.DigitalLocation{
			ID:       "test-location",
			Name:     "Updated Location",
			IsActive: true,
		}

		// Execute
		err := service.UpdateDigitalLocation(context.Background(), "test-user", location)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("RemoveDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)

		// Execute
		err := service.RemoveDigitalLocation(context.Background(), "test-user", "test-location")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("RemoveDigitalLocation - Not Found", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := sql.ErrNoRows

		service.dbAdapter = &MockDigitalDbAdapter{
			RemoveDigitalLocationFunc: func(ctx context.Context, userID, locationID string) error {
				return expectedErr
			},
		}

		// Execute
		err := service.RemoveDigitalLocation(context.Background(), "test-user", "test-location")

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("GetUserDigitalLocations - Cache Hit", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedLocations := []models.DigitalLocation{
			{ID: "1", Name: "Location 1"},
			{ID: "2", Name: "Location 2"},
		}

		// Execute
		locations, err := service.GetUserDigitalLocations(context.Background(), "test-user")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != len(expectedLocations) {
			t.Errorf("Expected %d locations, got %d", len(expectedLocations), len(locations))
		}
	})

	t.Run("GetUserDigitalLocations - Cache Miss", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedLocations := []models.DigitalLocation{
			{ID: "1", Name: "Location 1"},
			{ID: "2", Name: "Location 2"},
		}

		service.dbAdapter = &MockDigitalDbAdapter{
			GetUserDigitalLocationsFunc: func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
				return expectedLocations, nil
			},
		}

		// Execute
		locations, err := service.GetUserDigitalLocations(context.Background(), "test-user")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != len(expectedLocations) {
			t.Errorf("Expected %d locations, got %d", len(expectedLocations), len(locations))
		}
	})
}

func TestGameDigitalService_GetDigitalLocation(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())
	expectedErr := sql.ErrNoRows
	expectedLocation := models.DigitalLocation{}

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, expectedErr
		},
	}

	// Execute
	location, err := service.GetDigitalLocation(context.Background(), "test-user", "test-location")

	// Verify
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if location.ID != "" || location.Name != "" || location.IsActive != false {
		t.Errorf("Expected empty location, got %v", location)
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return expectedLocation, nil
		},
	}

	// Execute
	location, err = service.GetDigitalLocation(context.Background(), "test-user", "test-location")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if location.ID == "" {
		t.Error("Expected location to be returned")
	}
}

func TestGameDigitalService_AddDigitalLocation(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())
	expectedErr := errors.New("test error")
	location := models.DigitalLocation{
		Name:     "Test Location",
		IsActive: true,
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		FindDigitalLocationByNameFunc: func(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, expectedErr
		},
	}

	// Execute
	createdLocation, err := service.AddDigitalLocation(context.Background(), "test-user", location)

	// Verify
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if createdLocation.ID != "" || createdLocation.Name != "" || createdLocation.IsActive != false {
		t.Errorf("Expected empty location, got %v", createdLocation)
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		FindDigitalLocationByNameFunc: func(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		AddDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
			return location, nil
		},
	}

	// Execute
	createdLocation, err = service.AddDigitalLocation(context.Background(), "test-user", location)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if createdLocation.Name != location.Name {
		t.Errorf("Expected location name %s, got %s", location.Name, createdLocation.Name)
	}
}

func TestGameDigitalService_UpdateDigitalLocation(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())
	expectedErr := sql.ErrNoRows
	location := models.DigitalLocation{
		ID:       "test-location",
		Name:     "Updated Location",
		IsActive: true,
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, expectedErr
		},
	}

	// Execute
	err := service.UpdateDigitalLocation(context.Background(), "test-user", location)

	// Verify
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		UpdateDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) error {
			return nil
		},
	}

	// Execute
	err = service.UpdateDigitalLocation(context.Background(), "test-user", location)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGameDigitalService_RemoveDigitalLocation(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())
	expectedErr := sql.ErrNoRows

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, expectedErr
		},
	}

	// Execute
	err := service.RemoveDigitalLocation(context.Background(), "test-user", "test-location")

	// Verify
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	service.dbAdapter = &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		RemoveDigitalLocationFunc: func(ctx context.Context, userID, locationID string) error {
			return nil
		},
	}

	// Execute
	err = service.RemoveDigitalLocation(context.Background(), "test-user", "test-location")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
