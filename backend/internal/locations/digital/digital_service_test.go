package digital

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/models"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/stretchr/testify/assert"
)

/*
Behavior:
- GetDigitalLocations retrieves digital locations for a user
  - Attempts to retrieve from cache first
  - Falls back to database on cache miss
  - Caches database results for future requests

- CreateDigitalLocation adds a new digital location
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

- CreateDigitalLocation:
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
	RemoveDigitalLocationFunc func(ctx context.Context, userID string, locationIDs []string) (int64, error)

	// Subscription Operations
	GetSubscriptionFunc func(ctx context.Context, locationID string) (*models.Subscription, error)
	AddSubscriptionFunc func(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscriptionFunc func(ctx context.Context, subscription models.Subscription) error
	RemoveSubscriptionFunc func(ctx context.Context, locationID string) error
	ValidateSubscriptionExistsFunc func(ctx context.Context, locationID string) (*models.Subscription, error)

	// Payment Operations
	GetPaymentsFunc func(ctx context.Context, locationID string) ([]models.Payment, error)
	AddPaymentFunc func(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPaymentFunc func(ctx context.Context, paymentID int64) (*models.Payment, error)

	// Game Operations
	AddGameToDigitalLocationFunc func(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocationFunc func(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationIDFunc func(ctx context.Context, userID string, locationID string) ([]models.Game, error)

	// BFF Response
	GetAllDigitalLocationsBFFFunc func(ctx context.Context, userID string) (types.DigitalLocationsBFFResponse, error)
}

func (m *MockDigitalDbAdapter) GetAllDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
	return m.GetUserDigitalLocationsFunc(ctx, userID)
}

func (m *MockDigitalDbAdapter) GetSingleDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
	return m.GetDigitalLocationFunc(ctx, userID, locationID)
}

func (m *MockDigitalDbAdapter) FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
	return m.FindDigitalLocationByNameFunc(ctx, userID, name)
}

func (m *MockDigitalDbAdapter) CreateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
	return m.AddDigitalLocationFunc(ctx, userID, location)
}

func (m *MockDigitalDbAdapter) UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error {
	return m.UpdateDigitalLocationFunc(ctx, userID, location)
}

func (m *MockDigitalDbAdapter) DeleteDigitalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error) {
	return m.RemoveDigitalLocationFunc(ctx, userID, locationIDs)
}

// Subscription Operations
func (m *MockDigitalDbAdapter) GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error) {
	return m.GetSubscriptionFunc(ctx, locationID)
}

func (m *MockDigitalDbAdapter) CreateSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
	return m.AddSubscriptionFunc(ctx, subscription)
}

func (m *MockDigitalDbAdapter) UpdateSubscription(ctx context.Context, subscription models.Subscription) error {
	return m.UpdateSubscriptionFunc(ctx, subscription)
}

func (m *MockDigitalDbAdapter) DeleteSubscription(ctx context.Context, locationID string) error {
	return m.RemoveSubscriptionFunc(ctx, locationID)
}

// Payment Operations
func (m *MockDigitalDbAdapter) GetAllPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	return m.GetPaymentsFunc(ctx, locationID)
}

func (m *MockDigitalDbAdapter) CreatePayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	return m.AddPaymentFunc(ctx, payment)
}

func (m *MockDigitalDbAdapter) GetSinglePayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
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

// BFF Response
func (m *MockDigitalDbAdapter) GetAllDigitalLocationsBFF(ctx context.Context, userID string) (types.DigitalLocationsBFFResponse, error) {
	if m.GetAllDigitalLocationsBFFFunc != nil {
		return m.GetAllDigitalLocationsBFFFunc(ctx, userID)
	}
	return types.DigitalLocationsBFFResponse{}, nil
}

// MockDashboardCacheWrapper is a mock implementation of interfaces.DashboardCacheWrapper
type MockDashboardCacheWrapper struct {
	GetCachedDashboardBFFFunc func(ctx context.Context, userID string) (types.DashboardBFFResponse, error)
	SetCachedDashboardBFFFunc func(ctx context.Context, userID string, response types.DashboardBFFResponse) error
	InvalidateUserCacheFunc   func(ctx context.Context, userID string) error
}

func (m *MockDashboardCacheWrapper) GetCachedDashboardBFF(ctx context.Context, userID string) (types.DashboardBFFResponse, error) {
	if m.GetCachedDashboardBFFFunc != nil {
		return m.GetCachedDashboardBFFFunc(ctx, userID)
	}
	return types.DashboardBFFResponse{}, nil
}

func (m *MockDashboardCacheWrapper) SetCachedDashboardBFF(ctx context.Context, userID string, response types.DashboardBFFResponse) error {
	if m.SetCachedDashboardBFFFunc != nil {
		return m.SetCachedDashboardBFFFunc(ctx, userID, response)
	}
	return nil
}

func (m *MockDashboardCacheWrapper) InvalidateUserCache(ctx context.Context, userID string) error {
	if m.InvalidateUserCacheFunc != nil {
		return m.InvalidateUserCacheFunc(ctx, userID)
	}
	return nil
}

// MockSpendTrackingCacheWrapper is a mock implementation of interfaces.SpendTrackingCacheWrapper
type MockSpendTrackingCacheWrapper struct {
	GetCachedSpendTrackingBFFFunc func(ctx context.Context, userID string) (types.SpendTrackingBFFResponseFINAL, error)
	SetCachedSpendTrackingBFFFunc func(ctx context.Context, userID string, response types.SpendTrackingBFFResponseFINAL) error
	InvalidateUserCacheFunc       func(ctx context.Context, userID string) error
}

func (m *MockSpendTrackingCacheWrapper) GetCachedSpendTrackingBFF(ctx context.Context, userID string) (types.SpendTrackingBFFResponseFINAL, error) {
	if m.GetCachedSpendTrackingBFFFunc != nil {
		return m.GetCachedSpendTrackingBFFFunc(ctx, userID)
	}
	return types.SpendTrackingBFFResponseFINAL{}, nil
}

func (m *MockSpendTrackingCacheWrapper) SetCachedSpendTrackingBFF(ctx context.Context, userID string, response types.SpendTrackingBFFResponseFINAL) error {
	if m.SetCachedSpendTrackingBFFFunc != nil {
		return m.SetCachedSpendTrackingBFFFunc(ctx, userID, response)
	}
	return nil
}

func (m *MockSpendTrackingCacheWrapper) InvalidateUserCache(ctx context.Context, userID string) error {
	if m.InvalidateUserCacheFunc != nil {
		return m.InvalidateUserCacheFunc(ctx, userID)
	}
	return nil
}

// MockDigitalCacheWrapper is a mock implementation of interfaces.DigitalCacheWrapper
type MockDigitalCacheWrapper struct {
	GetCachedDigitalLocationsFunc      func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	SetCachedDigitalLocationsFunc      func(ctx context.Context, userID string, locations []models.DigitalLocation) error
	GetSingleCachedDigitalLocationFunc func(ctx context.Context, userID, locationID string) (*models.DigitalLocation, bool, error)
	SetSingleCachedDigitalLocationFunc func(ctx context.Context, userID string, location models.DigitalLocation) error
	InvalidateUserCacheFunc            func(ctx context.Context, userID string) error
	InvalidateDigitalLocationCacheFunc func(ctx context.Context, userID, locationID string) error
	InvalidateDigitalLocationsBulkFunc func(ctx context.Context, userID string, locationIDs []string) error

	// Subscription caching
	GetCachedSubscriptionFunc          func(ctx context.Context, locationID string) (*models.Subscription, bool, error)
	SetCachedSubscriptionFunc          func(ctx context.Context, locationID string, subscription models.Subscription) error
	InvalidateSubscriptionCacheFunc    func(ctx context.Context, locationID string) error

	// Payment caching
	GetCachedPaymentsFunc              func(ctx context.Context, locationID string) ([]models.Payment, error)
	SetCachedPaymentsFunc              func(ctx context.Context, locationID string, payments []models.Payment) error
	InvalidatePaymentsCacheFunc        func(ctx context.Context, locationID string) error

	// BFF Response caching
	GetCachedDigitalLocationsBFFFunc       func(ctx context.Context, userID string) (types.DigitalLocationsBFFResponse, error)
	SetCachedDigitalLocationsBFFFunc       func(ctx context.Context, userID string, response types.DigitalLocationsBFFResponse) error
	InvalidateDigitalLocationsBFFCacheFunc func(ctx context.Context, userID string) error
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

func (m *MockDigitalCacheWrapper) InvalidateDigitalLocationsBulk(ctx context.Context, userID string, locationIDs []string) error {
	if m.InvalidateDigitalLocationsBulkFunc != nil {
		return m.InvalidateDigitalLocationsBulkFunc(ctx, userID, locationIDs)
	}
	return nil
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

// BFF Response caching
func (m *MockDigitalCacheWrapper) GetCachedDigitalLocationsBFF(ctx context.Context, userID string) (types.DigitalLocationsBFFResponse, error) {
	if m.GetCachedDigitalLocationsBFFFunc != nil {
		return m.GetCachedDigitalLocationsBFFFunc(ctx, userID)
	}
	return types.DigitalLocationsBFFResponse{}, nil
}

func (m *MockDigitalCacheWrapper) SetCachedDigitalLocationsBFF(ctx context.Context, userID string, response types.DigitalLocationsBFFResponse) error {
	if m.SetCachedDigitalLocationsBFFFunc != nil {
		return m.SetCachedDigitalLocationsBFFFunc(ctx, userID, response)
	}
	return nil
}

func (m *MockDigitalCacheWrapper) InvalidateDigitalLocationsBFFCache(ctx context.Context, userID string) error {
	if m.InvalidateDigitalLocationsBFFCacheFunc != nil {
		return m.InvalidateDigitalLocationsBFFCacheFunc(ctx, userID)
	}
	return nil
}

func newMockGameDigitalServiceWithDefaults(logger *testutils.TestLogger) *GameDigitalService {
	mockDbAdapter := &MockDigitalDbAdapter{
		GetUserDigitalLocationsFunc: func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{
				{ID: "1", Name: "Location 1", URL: "http://example.com/1", IsSubscription: false, IsActive: true},
				{ID: "2", Name: "Location 2", URL: "http://example.com/2", IsSubscription: true, IsActive: true},
			}, nil
		},
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{ID: locationID, Name: "Test Location", URL: "http://example.com", IsSubscription: false, IsActive: true}, nil
		},
		FindDigitalLocationByNameFunc: func(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
			return models.DigitalLocation{ID: "test-id", Name: name, URL: "http://example.com", IsSubscription: false, IsActive: true}, nil
		},
		AddDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
			location.ID = "new-id"
			return location, nil
		},
		UpdateDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) error {
			return nil
		},
		RemoveDigitalLocationFunc: func(ctx context.Context, userID string, locationIDs []string) (int64, error) {
			return int64(len(locationIDs)), nil
		},
		GetSubscriptionFunc: func(ctx context.Context, locationID string) (*models.Subscription, error) {
			return &models.Subscription{ID: 1, LocationID: locationID}, nil
		},
		AddSubscriptionFunc: func(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
			return &models.Subscription{ID: 1, LocationID: subscription.LocationID}, nil
		},
		UpdateSubscriptionFunc: func(ctx context.Context, subscription models.Subscription) error {
			return nil
		},
		RemoveSubscriptionFunc: func(ctx context.Context, locationID string) error {
			return nil
		},
		GetPaymentsFunc: func(ctx context.Context, locationID string) ([]models.Payment, error) {
			return []models.Payment{{ID: 1, LocationID: locationID}}, nil
		},
		AddPaymentFunc: func(ctx context.Context, payment models.Payment) (*models.Payment, error) {
			return &models.Payment{ID: 1, LocationID: payment.LocationID}, nil
		},
		GetPaymentFunc: func(ctx context.Context, paymentID int64) (*models.Payment, error) {
			return &models.Payment{ID: paymentID}, nil
		},
		AddGameToDigitalLocationFunc: func(ctx context.Context, userID string, locationID string, gameID int64) error {
			return nil
		},
		RemoveGameFromDigitalLocationFunc: func(ctx context.Context, userID string, locationID string, gameID int64) error {
			return nil
		},
		GetGamesByDigitalLocationIDFunc: func(ctx context.Context, userID string, locationID string) ([]models.Game, error) {
			return []models.Game{{ID: 1, Name: "Test Game"}}, nil
		},
		ValidateSubscriptionExistsFunc: func(ctx context.Context, locationID string) (*models.Subscription, error) {
			return &models.Subscription{ID: 1, LocationID: locationID}, nil
		},
	}

	mockCacheWrapper := &MockDigitalCacheWrapper{
		GetCachedDigitalLocationsFunc: func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{
				{ID: "1", Name: "Location 1", URL: "http://example.com/1", IsSubscription: false, IsActive: true},
				{ID: "2", Name: "Location 2", URL: "http://example.com/2", IsSubscription: true, IsActive: true},
			}, nil
		},
		SetCachedDigitalLocationsFunc: func(ctx context.Context, userID string, locations []models.DigitalLocation) error {
			return nil
		},
		GetSingleCachedDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (*models.DigitalLocation, bool, error) {
			return &models.DigitalLocation{ID: locationID, Name: "Test Location", URL: "http://example.com", IsSubscription: false, IsActive: true}, true, nil
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
		GetCachedSubscriptionFunc: func(ctx context.Context, locationID string) (*models.Subscription, bool, error) {
			return &models.Subscription{ID: 1, LocationID: locationID}, true, nil
		},
		SetCachedSubscriptionFunc: func(ctx context.Context, locationID string, subscription models.Subscription) error {
			return nil
		},
		InvalidateSubscriptionCacheFunc: func(ctx context.Context, locationID string) error {
			return nil
		},
		GetCachedPaymentsFunc: func(ctx context.Context, locationID string) ([]models.Payment, error) {
			return []models.Payment{{ID: 1, LocationID: locationID}}, nil
		},
		SetCachedPaymentsFunc: func(ctx context.Context, locationID string, payments []models.Payment) error {
			return nil
		},
		InvalidatePaymentsCacheFunc: func(ctx context.Context, locationID string) error {
			return nil
		},
	}

	// Create a mock config
	mockConfig := &config.Config{
		Redis: config.RedisConfig{
			RedisTTL:     60,
			RedisTimeout: 5,
		},
	}

	// Create mock sanitizer and validator
	mockSanitizer, _ := security.NewSanitizer()
	mockValidator, _ := NewDigitalValidator(mockSanitizer)

	// Create mock dashboard and spend tracking cache wrappers
	mockDashboardCacheWrapper := &MockDashboardCacheWrapper{
		InvalidateUserCacheFunc: func(ctx context.Context, userID string) error {
			return nil
		},
	}

	mockSpendTrackingCacheWrapper := &MockSpendTrackingCacheWrapper{
		InvalidateUserCacheFunc: func(ctx context.Context, userID string) error {
			return nil
		},
	}

	// Directly create the service with mocks instead of using NewGameDigitalService
	return &GameDigitalService{
		dbAdapter:                 mockDbAdapter,
		cacheWrapper:              mockCacheWrapper,
		dashboardCacheWrapper:     mockDashboardCacheWrapper,
		spendTrackingCacheWrapper: mockSpendTrackingCacheWrapper,
		logger:                    logger,
		config:                    mockConfig,
		sanitizer:                 mockSanitizer,
		validator:                 mockValidator,
	}
}

func TestGameDigitalService(t *testing.T) {
	// Set up test logger
	logger := testutils.NewTestLogger()

	// Test cases
	t.Run("GetAllDigitalLocations - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)

		// Execute
		locations, err := service.GetAllDigitalLocations(context.Background(), "test-user")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != 2 {
			t.Errorf("Expected 2 locations, got %d", len(locations))
		}
	})

	t.Run("GetAllDigitalLocations - Error", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := errors.New("test error")

		// Override the default mock
		mockDb := &MockDigitalDbAdapter{}
		// Copy all the defaults from the original mock
		*mockDb = *service.dbAdapter.(*MockDigitalDbAdapter)
		// Override just the function we want to test
		mockDb.GetUserDigitalLocationsFunc = func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return nil, expectedErr
		}
		service.dbAdapter = mockDb

		// Set cache to return an error to force DB lookup
		mockCache := &MockDigitalCacheWrapper{}
		*mockCache = *service.cacheWrapper.(*MockDigitalCacheWrapper)
		mockCache.GetCachedDigitalLocationsFunc = func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return nil, errors.New("cache miss")
		}
		service.cacheWrapper = mockCache

		// Execute
		locations, err := service.GetAllDigitalLocations(context.Background(), "test-user")

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if len(locations) != 0 {
			t.Errorf("Expected empty locations slice, got %v", locations)
		}
	})

	t.Run("GetSingleDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)

		// Execute
		location, err := service.GetSingleDigitalLocation(context.Background(), "test-user", "test-location")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if location.ID == "" {
			t.Error("Expected location to be returned")
		}
	})

	t.Run("GetSingleDigitalLocation - Not Found", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := sql.ErrNoRows

		// Override the default mock
		mockDb := &MockDigitalDbAdapter{}
		*mockDb = *service.dbAdapter.(*MockDigitalDbAdapter)
		mockDb.GetDigitalLocationFunc = func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, expectedErr
		}
		service.dbAdapter = mockDb

		// Set cache to return an error to force DB lookup
		mockCache := &MockDigitalCacheWrapper{}
		*mockCache = *service.cacheWrapper.(*MockDigitalCacheWrapper)
		mockCache.GetSingleCachedDigitalLocationFunc = func(ctx context.Context, userID, locationID string) (*models.DigitalLocation, bool, error) {
			return nil, false, errors.New("cache miss")
		}
		service.cacheWrapper = mockCache

		// Execute
		location, err := service.GetSingleDigitalLocation(context.Background(), "test-user", "test-location")

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

	t.Run("CreateDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		newLocation := types.DigitalLocationRequest{
			Name:        "Test Location",
			URL:         "http://example.com",
			IsSubscription: false,
			IsActive:    true,
			PaymentMethod: "visa",
		}

		// Execute
		createdLocation, err := service.CreateDigitalLocation(context.Background(), "test-user", newLocation)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if createdLocation.Name != newLocation.Name {
			t.Errorf("Expected location name %s, got %s", newLocation.Name, createdLocation.Name)
		}
	})

	t.Run("CreateDigitalLocation - Error", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := errors.New("test error")
		newLocation := types.DigitalLocationRequest{
			Name:        "Test Location",
			URL:         "http://example.com",
			IsSubscription: false,
			IsActive:    true,
			PaymentMethod: "visa",
		}

		// Override the default mock
		mockDb := &MockDigitalDbAdapter{}
		*mockDb = *service.dbAdapter.(*MockDigitalDbAdapter)
		mockDb.AddDigitalLocationFunc = func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, expectedErr
		}
		service.dbAdapter = mockDb

		// Execute
		createdLocation, err := service.CreateDigitalLocation(context.Background(), "test-user", newLocation)

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) && err.Error() != "validation failed: "+expectedErr.Error() {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if createdLocation.ID != "" || createdLocation.Name != "" || createdLocation.IsActive != false {
			t.Errorf("Expected empty location, got %v", createdLocation)
		}
	})

	t.Run("UpdateDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		location := types.DigitalLocationRequest{
			ID:          "test-location",
			Name:        "Updated Location",
			URL:         "http://example.com",
			IsSubscription: false,
			IsActive:    true,
			PaymentMethod: "visa",
		}

		// Execute
		err := service.UpdateDigitalLocation(context.Background(), "test-user", location)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("DeleteDigitalLocation - Success", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)

		// Execute
		count, err := service.DeleteDigitalLocation(context.Background(), "test-user", []string{"123e4567-e89b-12d3-a456-426614174000"})

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		assert.Equal(t, int64(1), count)
	})

	t.Run("DeleteDigitalLocation - Not Found", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedErr := sql.ErrNoRows

		// Override the default mock
		mockDb := &MockDigitalDbAdapter{}
		*mockDb = *service.dbAdapter.(*MockDigitalDbAdapter)
		mockDb.RemoveDigitalLocationFunc = func(ctx context.Context, userID string, locationIDs []string) (int64, error) {
			return 0, expectedErr
		}
		service.dbAdapter = mockDb

		// Execute
		count, err := service.DeleteDigitalLocation(context.Background(), "test-user", []string{"123e4567-e89b-12d3-a456-426614174000"})

		// Verify
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		assert.Equal(t, int64(0), count)
	})

	t.Run("GetAllDigitalLocations - Cache Hit", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedLocations := []models.DigitalLocation{
			{ID: "1", Name: "Location 1", URL: "http://example.com/1", IsSubscription: false, IsActive: true},
			{ID: "2", Name: "Location 2", URL: "http://example.com/2", IsSubscription: true, IsActive: true},
		}

		// Execute
		locations, err := service.GetAllDigitalLocations(context.Background(), "test-user")

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(locations) != len(expectedLocations) {
			t.Errorf("Expected %d locations, got %d", len(expectedLocations), len(locations))
		}
	})

	t.Run("GetAllDigitalLocations - Cache Miss", func(t *testing.T) {
		// Setup
		service := newMockGameDigitalServiceWithDefaults(logger)
		expectedLocations := []models.DigitalLocation{
			{ID: "1", Name: "Location 1", URL: "http://example.com/1", IsSubscription: false, IsActive: true},
			{ID: "2", Name: "Location 2", URL: "http://example.com/2", IsSubscription: true, IsActive: true},
		}

		// Override cache to simulate a miss
		mockCache := &MockDigitalCacheWrapper{}
		*mockCache = *service.cacheWrapper.(*MockDigitalCacheWrapper)
		mockCache.GetCachedDigitalLocationsFunc = func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return nil, errors.New("cache miss")
		}
		service.cacheWrapper = mockCache

		// Override DB to return expected locations
		mockDb := &MockDigitalDbAdapter{}
		*mockDb = *service.dbAdapter.(*MockDigitalDbAdapter)
		mockDb.GetUserDigitalLocationsFunc = func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return expectedLocations, nil
		}
		service.dbAdapter = mockDb

		// Execute
		locations, err := service.GetAllDigitalLocations(context.Background(), "test-user")

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
	expectedLocation := models.DigitalLocation{
		ID:          "test-id",
		Name:        "Test Location",
		URL:         "http://example.com",
		IsSubscription: false,
		IsActive:    true,
	}

	// Set up for error test
	mockDb1 := &MockDigitalDbAdapter{}
	*mockDb1 = *service.dbAdapter.(*MockDigitalDbAdapter)
	mockDb1.GetDigitalLocationFunc = func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
		return models.DigitalLocation{}, expectedErr
	}

	// Set cache to miss
	mockCache1 := &MockDigitalCacheWrapper{}
	*mockCache1 = *service.cacheWrapper.(*MockDigitalCacheWrapper)
	mockCache1.GetSingleCachedDigitalLocationFunc = func(ctx context.Context, userID, locationID string) (*models.DigitalLocation, bool, error) {
		return nil, false, errors.New("cache miss")
	}

	service.dbAdapter = mockDb1
	service.cacheWrapper = mockCache1

	// Execute error test
	location, err := service.GetSingleDigitalLocation(context.Background(), "test-user", "test-location")

	// Verify error case
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if location.ID != "" || location.Name != "" || location.IsActive != false {
		t.Errorf("Expected empty location, got %v", location)
	}

	// Set up success test
	mockDb2 := &MockDigitalDbAdapter{}
	*mockDb2 = *service.dbAdapter.(*MockDigitalDbAdapter)
	mockDb2.GetDigitalLocationFunc = func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
		return expectedLocation, nil
	}

	// Set cache to miss
	mockCache2 := &MockDigitalCacheWrapper{}
	*mockCache2 = *service.cacheWrapper.(*MockDigitalCacheWrapper)
	mockCache2.GetSingleCachedDigitalLocationFunc = func(ctx context.Context, userID, locationID string) (*models.DigitalLocation, bool, error) {
		return nil, false, errors.New("cache miss")
	}

	service.dbAdapter = mockDb2
	service.cacheWrapper = mockCache2

	// Execute success test
	location, err = service.GetSingleDigitalLocation(context.Background(), "test-user", "test-location")

	// Verify success case
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if location.ID != expectedLocation.ID {
		t.Errorf("Expected location ID %s, got %s", expectedLocation.ID, location.ID)
	}
}

func TestGameDigitalService_AddDigitalLocation(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())
	expectedErr := errors.New("test error")
	location := types.DigitalLocationRequest{
		Name:        "Test Location",
		URL:         "http://example.com",
		IsSubscription: false,
		IsActive:    true,
		PaymentMethod: "visa",
	}

	// Set up error test
	mockDb1 := &MockDigitalDbAdapter{}
	*mockDb1 = *service.dbAdapter.(*MockDigitalDbAdapter)
	mockDb1.AddDigitalLocationFunc = func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
		return models.DigitalLocation{}, expectedErr
	}
	service.dbAdapter = mockDb1

	// Execute error test
	createdLocation, err := service.CreateDigitalLocation(context.Background(), "test-user", location)

	// Verify error case
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) && err.Error() != "validation failed: "+expectedErr.Error() {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if createdLocation.ID != "" || createdLocation.Name != "" || createdLocation.IsActive != false {
		t.Errorf("Expected empty location, got %v", createdLocation)
	}

	// Set up success test
	mockDb2 := &MockDigitalDbAdapter{}
	*mockDb2 = *service.dbAdapter.(*MockDigitalDbAdapter)
	mockDb2.AddDigitalLocationFunc = func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
		location.ID = "new-id"
		return location, nil
	}
	service.dbAdapter = mockDb2

	// Execute success test
	createdLocation, err = service.CreateDigitalLocation(context.Background(), "test-user", location)

	// Verify success case
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if createdLocation.Name != location.Name {
		t.Errorf("Expected location name %s, got %s", location.Name, createdLocation.Name)
	}
	if createdLocation.ID != "new-id" {
		t.Errorf("Expected location ID 'new-id', got %s", createdLocation.ID)
	}
}

func TestGameDigitalService_UpdateDigitalLocation(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())
	expectedErr := sql.ErrNoRows
	location := types.DigitalLocationRequest{
		ID:          "test-location",
		Name:        "Updated Location",
		URL:         "http://example.com",
		IsSubscription: false,
		IsActive:    true,
		PaymentMethod: "visa",
	}

	// Set up error test
	mockDb1 := &MockDigitalDbAdapter{}
	*mockDb1 = *service.dbAdapter.(*MockDigitalDbAdapter)
	mockDb1.UpdateDigitalLocationFunc = func(ctx context.Context, userID string, location models.DigitalLocation) error {
		return expectedErr
	}
	service.dbAdapter = mockDb1

	// Execute error test
	err := service.UpdateDigitalLocation(context.Background(), "test-user", location)

	// Verify error case
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) && !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected error containing %v, got %v", expectedErr, err)
	}

	// Set up success test
	mockDb2 := &MockDigitalDbAdapter{}
	*mockDb2 = *service.dbAdapter.(*MockDigitalDbAdapter)
	mockDb2.UpdateDigitalLocationFunc = func(ctx context.Context, userID string, location models.DigitalLocation) error {
		return nil
	}
	service.dbAdapter = mockDb2

	// Execute success test
	err = service.UpdateDigitalLocation(context.Background(), "test-user", location)

	// Verify success case
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGameDigitalService_RemoveDigitalLocation(t *testing.T) {
	ctx := context.Background()
	log := testutils.NewTestLogger()

	t.Run("single location success", func(t *testing.T) {
		mockDb := &MockDigitalDbAdapter{}
		mockDb.RemoveDigitalLocationFunc = func(ctx context.Context, userID string, locationIDs []string) (int64, error) {
			return 1, nil
		}

		mockCache := &MockDigitalCacheWrapper{}
		mockCache.InvalidateDigitalLocationCacheFunc = func(ctx context.Context, userID, locationID string) error {
			return nil
		}

		service := newMockGameDigitalServiceWithDefaults(log)
		service.dbAdapter = mockDb
		service.cacheWrapper = mockCache

		locationIDs := []string{"123e4567-e89b-12d3-a456-426614174000"}
		count, err := service.DeleteDigitalLocation(ctx, "test-user", locationIDs)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("bulk locations success", func(t *testing.T) {
		mockDb := &MockDigitalDbAdapter{}
		mockDb.RemoveDigitalLocationFunc = func(ctx context.Context, userID string, locationIDs []string) (int64, error) {
			return int64(len(locationIDs)), nil
		}

		mockCache := &MockDigitalCacheWrapper{}
		mockCache.InvalidateDigitalLocationCacheFunc = func(ctx context.Context, userID, locationID string) error {
			return nil
		}

		service := newMockGameDigitalServiceWithDefaults(log)
		service.dbAdapter = mockDb
		service.cacheWrapper = mockCache

		locationIDs := []string{
			"123e4567-e89b-12d3-a456-426614174000",
			"123e4567-e89b-12d3-a456-426614174001",
			"123e4567-e89b-12d3-a456-426614174002",
		}
		count, err := service.DeleteDigitalLocation(ctx, "test-user", locationIDs)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("database error", func(t *testing.T) {
		mockDb := &MockDigitalDbAdapter{}
		mockDb.RemoveDigitalLocationFunc = func(ctx context.Context, userID string, locationIDs []string) (int64, error) {
			return 0, errors.New("database error")
		}

		service := newMockGameDigitalServiceWithDefaults(log)
		service.dbAdapter = mockDb

		locationIDs := []string{"123e4567-e89b-12d3-a456-426614174000"}
		count, err := service.DeleteDigitalLocation(ctx, "test-user", locationIDs)
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("cache invalidation error", func(t *testing.T) {
		mockDb := &MockDigitalDbAdapter{}
		mockDb.RemoveDigitalLocationFunc = func(ctx context.Context, userID string, locationIDs []string) (int64, error) {
			return 1, nil
		}

		mockCache := &MockDigitalCacheWrapper{}
		mockCache.InvalidateDigitalLocationCacheFunc = func(ctx context.Context, userID, locationID string) error {
			return errors.New("cache error")
		}

		service := newMockGameDigitalServiceWithDefaults(log)
		service.dbAdapter = mockDb
		service.cacheWrapper = mockCache

		locationIDs := []string{"123e4567-e89b-12d3-a456-426614174000"}
		count, err := service.DeleteDigitalLocation(ctx, "test-user", locationIDs)
		assert.NoError(t, err) // Cache errors should not fail the operation
		assert.Equal(t, int64(1), count)
	})

	t.Run("validation error", func(t *testing.T) {
		service := newMockGameDigitalServiceWithDefaults(log)

		locationIDs := []string{} // Empty location IDs should fail validation
		count, err := service.DeleteDigitalLocation(ctx, "test-user", locationIDs)
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestRemoveDigitalLocation_ErrorHandling(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())

	t.Run("empty location IDs", func(t *testing.T) {
		// Execute
		count, err := service.DeleteDigitalLocation(context.Background(), "test-user", []string{})

		// Verify
		if err == nil {
			t.Error("expected error but got nil")
		}
		assert.Equal(t, int64(0), count)
	})

	t.Run("empty user ID", func(t *testing.T) {
		// Execute
		count, err := service.DeleteDigitalLocation(context.Background(), "", []string{"test-location"})

		// Verify
		if err == nil {
			t.Error("expected error but got nil")
		}
		assert.Equal(t, int64(0), count)
	})
}

func TestRemoveDigitalLocation_Transaction(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())

	t.Run("transaction success", func(t *testing.T) {
		// Setup
		mockDb := &MockDigitalDbAdapter{}
		*mockDb = *service.dbAdapter.(*MockDigitalDbAdapter)
		mockDb.RemoveDigitalLocationFunc = func(ctx context.Context, userID string, locationIDs []string) (int64, error) {
			return int64(len(locationIDs)), nil
		}
		service.dbAdapter = mockDb

		// Execute
		count, err := service.DeleteDigitalLocation(context.Background(), "test-user", []string{
			"123e4567-e89b-12d3-a456-426614174000",
			"123e4567-e89b-12d3-a456-426614174001",
		})

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assert.Equal(t, int64(2), count)
	})

	t.Run("transaction rollback", func(t *testing.T) {
		// Setup
		mockDb := &MockDigitalDbAdapter{}
		*mockDb = *service.dbAdapter.(*MockDigitalDbAdapter)
		mockDb.RemoveDigitalLocationFunc = func(ctx context.Context, userID string, locationIDs []string) (int64, error) {
			return 0, fmt.Errorf("transaction failed")
		}
		service.dbAdapter = mockDb

		// Execute
		count, err := service.DeleteDigitalLocation(context.Background(), "test-user", []string{"123e4567-e89b-12d3-a456-426614174000"})

		// Verify
		if err == nil {
			t.Error("expected error but got nil")
		}
		assert.Equal(t, int64(0), count)
	})
}

func TestRemoveDigitalLocation_Performance(t *testing.T) {
	// Setup
	service := newMockGameDigitalServiceWithDefaults(testutils.NewTestLogger())

	t.Run("large number of locations", func(t *testing.T) {
		// Setup
		locationIDs := make([]string, 100)
		for i := range locationIDs {
			locationIDs[i] = fmt.Sprintf("123e4567-e89b-12d3-a456-426614174%03d", i)
		}

		mockDb := &MockDigitalDbAdapter{}
		*mockDb = *service.dbAdapter.(*MockDigitalDbAdapter)
		mockDb.RemoveDigitalLocationFunc = func(ctx context.Context, userID string, locationIDs []string) (int64, error) {
			return int64(len(locationIDs)), nil
		}
		service.dbAdapter = mockDb

		// Execute
		count, err := service.DeleteDigitalLocation(context.Background(), "test-user", locationIDs)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assert.Equal(t, int64(100), count)
	})
}
