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

func DefaultPhysicalValidator() *MockPhysicalValidator {
	return &MockPhysicalValidator{
		ValidatePhysicalLocationFunc: func(location models.PhysicalLocation) (models.PhysicalLocation, error) {
			return location, nil
		},
	}
}

func DefaultPhysicalDbAdapter() *MockPhysicalDbAdapter {
	return &MockPhysicalDbAdapter{
		GetPhysicalLocationsFunc: func(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
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
