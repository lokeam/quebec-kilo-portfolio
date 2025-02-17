package mocks

import (
	"context"

	igdb "github.com/Henry-Sarabia/igdb"
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
				SearchGamesFunc: func(ctx context.Context, query string, limit int) ([]*igdb.Game, error) {
						// Default: return an empty slice (or a minimal dummy value).
						return []*igdb.Game{}, nil
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