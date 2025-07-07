package mocks

import (
	"context"
	"errors"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/stretchr/testify/mock"
)

// DefaultSanitizer returns a MockSanitizer with a passing default.
func DefaultSanitizer() *MockSanitizer {
		return &MockSanitizer{
				// Override the SanitizeFunc to the default behavior.
				SanitizeFunc: func(query string) (string, error) {
						// Default: leave the query unchanged.
						return query, nil
				},
		}
}

// DefaultValidator returns a MockValidator with a passing default.
func DefaultValidator() *MockValidator {
		return &MockValidator{
				ValidateFunc: func(query searchdef.SearchQuery) error {
						// Default: always valid.
						return nil
				},
		}
}

// DefaultIGDBAdapter returns a MockIGDBAdapter with default (happy path) behavior.
func DefaultIGDBAdapter() *MockIGDBAdapter {
		return &MockIGDBAdapter{
				SearchGamesFunc: func(ctx context.Context, query string, limit int) ([]*models.Game, error) {
						// Default: return an empty slice (or a minimal dummy value).
						return []*models.Game{}, nil
				},
				UpdateTokenFunc: func(token string) error {
						// Default: do nothing, return success.
						return nil
				},
		}
}

// DefaultCacheWrapper returns a MockCacheWrapper with default (no cache hit) behavior.
func DefaultCacheWrapper() *MockCacheWrapper {
		return &MockCacheWrapper{
				GetCachedResultsFunc: func(ctx context.Context, sq searchdef.SearchQuery) (*searchdef.SearchResult, error) {
						// Default: simulate cache miss.
						return nil, nil
				},
				SetCachedResultsFunc: func(ctx context.Context, sq searchdef.SearchQuery, result *searchdef.SearchResult) error {
						// Default: do nothing.
						return nil
				},
				TimeToLive: 60, // example TTL in seconds.
		}
}

// ---------- Physical ----------
func DefaultPhysicalValidator() *MockPhysicalValidator {
	return &MockPhysicalValidator{
		ValidatePhysicalLocationFunc: func(location models.PhysicalLocation) (models.PhysicalLocation, error) {
			return location, nil
		},
		ValidatePhysicalLocationCreationFunc: func(location models.PhysicalLocation) (models.PhysicalLocation, error) {
			return location, nil
		},
		ValidatePhysicalLocationUpdateFunc: func(update, existing models.PhysicalLocation) (models.PhysicalLocation, error) {
			return update, nil
		},
		ValidateRemovePhysicalLocationFunc: func(userID string, locationIDs []string) ([]string, error) {
			return locationIDs, nil
		},
	}
}

func DefaultPhysicalDbAdapter() *MockPhysicalDbAdapter {
	m := &MockPhysicalDbAdapter{}

	defaultLocation := models.PhysicalLocation{
		ID:             "location-1",
		Name:           "Home",
		Label:          "Primary",
		LocationType:   "Home",
		MapCoordinates: models.PhysicalMapCoordinates{
			Coords:         "40.7128,-74.0060",
			GoogleMapsLink: "https://www.google.com/maps/search/?api=1&query=40.7128,-74.0060",
		},
		BgColor:        "red",
	}

	m.On("GetSinglePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(defaultLocation, nil)
	m.On("GetAllPhysicalLocations", mock.Anything, mock.Anything).
		Return([]models.PhysicalLocation{defaultLocation}, nil)
	m.On("CreatePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(defaultLocation, nil)
	m.On("UpdatePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(defaultLocation, nil)
	m.On("DeletePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(int64(1), nil)

	return m
}

func DefaultPhysicalCacheWrapper() *MockPhysicalCacheWrapper {
	m := &MockPhysicalCacheWrapper{}

	m.On("GetCachedPhysicalLocations", mock.Anything, mock.Anything).
		Return(nil, errors.New("cache miss"))
	m.On("SetCachedPhysicalLocations", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	m.On("GetSingleCachedPhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, false, errors.New("cache miss"))
	m.On("SetSingleCachedPhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	m.On("InvalidateUserCache", mock.Anything, mock.Anything).
		Return(nil)
	m.On("InvalidateLocationCache", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	return m
}

// DefaultPhysicalService returns a MockPhysicalService with default behavior
func DefaultPhysicalService() *MockPhysicalService {
	return &MockPhysicalService{}
}

// ---------- Sublocation ----------
func DefaultSublocationValidator() *MockSublocationValidator {
	return &MockSublocationValidator{
		ValidateSublocationFunc: func(sublocation models.Sublocation) (models.Sublocation, error) {
			return sublocation, nil
		},
		ValidateSublocationUpdateFunc: func(update, existing models.Sublocation) (models.Sublocation, error) {
			return update, nil
		},
		ValidateSublocationCreationFunc: func(sublocation models.Sublocation) (models.Sublocation, error) {
			return sublocation, nil
		},
	}
}

func DefaultSublocationDbAdapter() *MockSublocationDbAdapter {
	return &MockSublocationDbAdapter{
		GetSublocationFunc: func(
			ctx context.Context,
			userID string,
			sublocationID string,
		) (models.Sublocation, error) {
			return models.Sublocation{
				ID:                 "sublocation-1",
				UserID:            userID,
				PhysicalLocationID: "physical-location-1",
				Name:              "Sublocation 1",
				LocationType:      "shelf",
				StoredItems:       0,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}, nil
		},
		GetUserSublocationsFunc: func(
			ctx context.Context,
			userID string,
		) ([]models.Sublocation, error) {
			return []models.Sublocation{
				{
					ID:                 "sublocation-1",
					UserID:            userID,
					PhysicalLocationID: "physical-location-1",
					Name:              "Sublocation 1",
					LocationType:      "shelf",
					StoredItems:       0,
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				},
				{
					ID:                 "sublocation-2",
					UserID:            userID,
					PhysicalLocationID: "physical-location-1",
					Name:              "Sublocation 2",
					LocationType:      "console",
					StoredItems:       0,
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				},
			}, nil
		},
		AddSublocationFunc: func(
			ctx context.Context,
			userID string,
			sublocation models.Sublocation,
		) (models.Sublocation, error) {
			return sublocation, nil
		},
		UpdateSublocationFunc: func(
			ctx context.Context,
			userID string,
			sublocation models.Sublocation,
		) error {
			return nil
		},
		DeleteSublocationFunc: func(
			ctx context.Context,
			userID string,
			sublocationIDs []string,
		) (types.DeleteSublocationResponse, error) {
			return types.DeleteSublocationResponse{
				Success:      true,
				DeletedCount: len(sublocationIDs),
				SublocationIDs: sublocationIDs,
			}, nil
		},
		CheckGameInAnySublocationFunc: func(
			ctx context.Context,
			userGameID string,
		) (bool, error) {
			return false, nil
		},
		CheckGameInSublocationFunc: func(
			ctx context.Context,
			userGameID string,
			sublocationID string,
		) (bool, error) {
			return false, nil
		},
		CheckGameOwnershipFunc: func(
			ctx context.Context,
			userID string,
			userGameID string,
		) (bool, error) {
			return true, nil
		},
		MoveGameToSublocationFunc: func(
			ctx context.Context,
			userID string,
			userGameID string,
			targetSublocationID string,
		) error {
			return nil
		},
		RemoveGameFromSublocationFunc: func(
			ctx context.Context,
			userID string,
			userGameID string,
		) error {
			return nil
		},
	}
}

func DefaultSublocationCacheWrapper() *MockSublocationCacheWrapper {
	return &MockSublocationCacheWrapper{
		GetCachedSublocationsFunc: func(ctx context.Context, userID string) ([]models.Sublocation, error) {
			return nil, nil
		},
		SetCachedSublocationsFunc: func(ctx context.Context, userID string, sublocations []models.Sublocation) error {
			return nil
		},
		GetSingleCachedSublocationFunc: func(ctx context.Context, userID string, sublocationID string) (*models.Sublocation, bool, error) {
			return nil, false, nil
		},
		SetSingleCachedSublocationFunc: func(ctx context.Context, userID string, sublocation models.Sublocation) error {
			return nil
		},
		InvalidateUserCacheFunc: func(ctx context.Context, userID string) error {
			return nil
		},
		InvalidateSublocationCacheFunc: func(ctx context.Context, userID string, locationID string) error {
			return nil
		},
		InvalidateLocationCacheFunc: func(ctx context.Context, userID string, locationID string) error {
			return nil
		},
	}
}


// ---------- Digital ----------
func DefaultDigitalValidator() *MockDigitalValidator {
	return &MockDigitalValidator{
		ValidateDigitalLocationFunc: func(digitalLocation models.DigitalLocation) (models.DigitalLocation, error) {
			return digitalLocation, nil
		},
	}
}

func DefaultDigitalDbAdapter() *MockDigitalDbAdapter {
	return &MockDigitalDbAdapter{
		GetDigitalLocationFunc: func(ctx context.Context, userID, digitalLocationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		GetDigitalLocationsFunc: func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{}, nil
		},
		AddDigitalLocationFunc: func(ctx context.Context, userID string, digitalLocation models.DigitalLocation) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		UpdateDigitalLocationFunc: func(ctx context.Context, userID string, digitalLocation models.DigitalLocation) error {
			return nil
		},
		RemoveDigitalLocationFunc: func(ctx context.Context, userID string, locationIDs []string) (int64, error) {
			return 0, nil
		},
		FindDigitalLocationByNameFunc: func(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		// Subscription Operations
		GetSubscriptionFunc: func(ctx context.Context, locationID string) (*models.Subscription, error) {
			return &models.Subscription{
				ID:              1,
				LocationID:      locationID,
				BillingCycle:    "monthly",
				CostPerCycle:    9.99,
				NextPaymentDate: time.Now().AddDate(0, 1, 0),
				PaymentMethod:   "credit_card",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}, nil
		},
		AddSubscriptionFunc: func(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
			subscription.ID = 1
			subscription.CreatedAt = time.Now()
			subscription.UpdatedAt = time.Now()
			return &subscription, nil
		},
		UpdateSubscriptionFunc: func(ctx context.Context, subscription models.Subscription) error {
			return nil
		},
		RemoveSubscriptionFunc: func(ctx context.Context, locationID string) error {
			return nil
		},
		// Payment Operations
		GetPaymentsFunc: func(ctx context.Context, locationID string) ([]models.Payment, error) {
			return []models.Payment{
				{
					ID:            1,
					LocationID:    locationID,
					Amount:        9.99,
					PaymentDate:   time.Now(),
					PaymentMethod: "credit_card",
					TransactionID: "txn_123",
					CreatedAt:     time.Now(),
				},
				{
					ID:            2,
					LocationID:    locationID,
					Amount:        9.99,
					PaymentDate:   time.Now().AddDate(0, -1, 0),
					PaymentMethod: "credit_card",
					TransactionID: "txn_456",
					CreatedAt:     time.Now().AddDate(0, -1, 0),
				},
			}, nil
		},
		AddPaymentFunc: func(ctx context.Context, payment models.Payment) (*models.Payment, error) {
			payment.ID = 1
			payment.CreatedAt = time.Now()
			return &payment, nil
		},
		GetPaymentFunc: func(ctx context.Context, paymentID int64) (*models.Payment, error) {
			return &models.Payment{
				ID:            paymentID,
				LocationID:    "default-digital-location",
				Amount:        9.99,
				PaymentDate:   time.Now(),
				PaymentMethod: "credit_card",
				TransactionID: "txn_123",
				CreatedAt:     time.Now(),
			}, nil
		},
	}
}

func DefaultDigitalCacheWrapper() *MockDigitalCacheWrapper {
	return &MockDigitalCacheWrapper{
		GetCachedDigitalLocationsFunc: func(
			ctx context.Context,
			userID string,
		) ([]models.DigitalLocation, error) {
			return nil, errors.New("cache miss")
		},
		SetCachedDigitalLocationsFunc: func(
			ctx context.Context,
			userID string,
			locations []models.DigitalLocation,
		) error {
			return nil
		},
		GetSingleCachedDigitalLocationFunc: func(
			ctx context.Context,
			userID,
			digitalLocationID string,
		) (*models.DigitalLocation, bool, error) {
			return nil, false, errors.New("cache miss")
		},
		SetSingleCachedDigitalLocationFunc: func(
			ctx context.Context,
			userID string,
			location models.DigitalLocation,
		) error {
			return nil
		},
		InvalidateUserCacheFunc: func(
			ctx context.Context,
			userID string,
		) error {
			return nil
		},
		InvalidateDigitalLocationCacheFunc: func(
			ctx context.Context,
			userID,
			digitalLocationID string,
		) error {
			return nil
		},

		// Subscription caching
		GetCachedSubscriptionFunc: func(
			ctx context.Context,
			locationID string,
		) (*models.Subscription, bool, error) {
			return nil, false, errors.New("cache miss")
		},
		SetCachedSubscriptionFunc: func(
			ctx context.Context,
			locationID string,
			subscription models.Subscription,
		) error {
			return nil
		},
		InvalidateSubscriptionCacheFunc: func(
			ctx context.Context,
			locationID string,
		) error {
			return nil
		},

		// Payment caching
		GetCachedPaymentsFunc: func(
			ctx context.Context,
			locationID string,
		) ([]models.Payment, error) {
			return nil, errors.New("cache miss")
		},
		SetCachedPaymentsFunc: func(
			ctx context.Context,
			locationID string,
			payments []models.Payment,
		) error {
			return nil
		},
		InvalidatePaymentsCacheFunc: func(
			ctx context.Context,
			locationID string,
		) error {
			return nil
		},
	}
}
