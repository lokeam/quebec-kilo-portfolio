package mocks

import (
	"context"
	"errors"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
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
	}
}

func DefaultPhysicalDbAdapter() *MockPhysicalDbAdapter {
	return &MockPhysicalDbAdapter{
		GetPhysicalLocationsFunc: func(
			ctx context.Context,
			userID string,
		) ([]models.PhysicalLocation, error) {
			return []models.PhysicalLocation{
				{
					ID:             "location-1",
					UserID:         userID,
					Name:           "Home",
					Label:          "Primary",
					LocationType:   "Home",
					MapCoordinates: "40.7128,-74.0060",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
			}, nil
		},
		GetPhysicalLocationFunc: func(
			ctx context.Context,
			userID,
			locationID string,
		) (*models.PhysicalLocation, error) {
			return &models.PhysicalLocation{
				ID:              locationID,
				UserID:          userID,
				Name:            "Home",
				Label:           "Primary",
				LocationType:    "Home",
				MapCoordinates:  "40.7128,-74.0060",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}, nil
		},
		CreatePhysicalLocationFunc: func(
			ctx context.Context,
			userID string,
			location models.PhysicalLocation,
		) error {
			return nil
		},
		UpdatePhysicalLocationFunc: func(
			ctx context.Context,
			userID string,
			location models.PhysicalLocation,
		) error {
			return nil
		},
		DeletePhysicalLocationFunc: func(
			ctx context.Context,
			userID,
			locationID string,
		) error {
			return nil
		},
	}
}

func DefaultPhysicalCacheWrapper() *MockPhysicalCacheWrapper {
	return &MockPhysicalCacheWrapper{
		GetCachedPhysicalLocationsFunc: func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
			return nil, errors.New("cache miss")
		},
		SetCachedPhysicalLocationsFunc: func(ctx context.Context, userID string, locations []models.PhysicalLocation) error {
			return nil
		},
		InvalidateUserCacheFunc: func(ctx context.Context, userID string) error {
			return nil
		},
		InvalidateLocationCacheFunc: func(ctx context.Context, userID, locationID string) error {
			return nil
		},
	}
}

// ---------- Sublocation ----------
func DefaultSublocationValidator() *MockSublocationValidator {
	return &MockSublocationValidator{
		ValidateSublocationFunc: func(sublocation models.Sublocation) (models.Sublocation, error) {
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
				ID:           "sublocation-1",
				UserID:       userID,
				Name:         "Sublocation 1",
				LocationType: "house",
				BgColor:      "red",
				Capacity:     10,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		},
		GetSublocationsFunc: func(
			ctx context.Context,
			userID string,
		) ([]models.Sublocation, error) {
			return []models.Sublocation{
				{
					ID:           "sublocation-1",
					UserID:       userID,
					Name:         "Sublocation 1",
					LocationType: "house",
					BgColor:      "red",
					Capacity:     10,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
				{
					ID:           "sublocation-2",
					UserID:       userID,
					Name:         "Sublocation 2",
					LocationType: "apartment",
					BgColor:      "blue",
					Capacity:     10,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
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
			sublocationID string,
		) error {
			return nil
		},
	}
}

func DefaultSublocationCacheWrapper() *MockSublocationCacheWrapper {
	return &MockSublocationCacheWrapper{
		GetCachedSublocationsFunc: func(
			ctx context.Context,
			userID string,
		) ([]models.Sublocation, error) {
			return nil, errors.New("cache miss")
		},
		SetCachedSublocationsFunc: func(
			ctx context.Context,
			userID string,
			locations []models.Sublocation,
		) error {
			return nil
		},
		GetSingleCachedSublocationFunc: func(
			ctx context.Context,
			userID,
			sublocationID string,
		) (*models.Sublocation, bool, error) {
			return nil, false, errors.New("cache miss")
		},
		SetSingleCachedSublocationFunc: func(
			ctx context.Context,
			userID string,
			sublocation models.Sublocation,
		) error {
			return nil
		},
		InvalidateUserCacheFunc: func(
			ctx context.Context,
			userID string,
		) error {
			return nil
		},
		InvalidateSublocationCacheFunc: func(
			ctx context.Context,
			userID,
			sublocationID string,
		) error {
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
		GetDigitalLocationFunc: func(
			ctx context.Context,
			userID,
			digitalLocationID string,
		) (models.DigitalLocation, error) {
			return models.DigitalLocation{
				ID: digitalLocationID,
				UserID: userID,
				Name: "Digital Location 1",
				IsActive: true,
				URL: "https://example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
		GetDigitalLocationsFunc: func(
			ctx context.Context,
			userID string,
		) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{
				{
					ID: "digital-location-1",
					UserID: userID,
					Name: "Digital Location 1",
					IsActive: true,
					URL: "https://example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID: "digital-location-2",
					UserID: userID,
					Name: "Digital Location 2",
					IsActive: true,
					URL: "https://example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}, nil
		},
		AddDigitalLocationFunc: func(
			ctx context.Context,
			userID string,
			digitalLocation models.DigitalLocation,
		) (models.DigitalLocation, error) {
			return digitalLocation, nil
		},
		UpdateDigitalLocationFunc: func(
			ctx context.Context,
			userID string,
			digitalLocation models.DigitalLocation,
		) error {
			return nil
		},
		DeleteDigitalLocationFunc: func(
			ctx context.Context,
			userID string,
			digitalLocationID string,
		) error {
			return nil
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
	}
}
