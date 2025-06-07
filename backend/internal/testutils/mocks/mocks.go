package mocks

import (
	"context"
	"errors"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
)

// FakeSanitizer is a simple implementation of interfaces.Sanitizer used across tests.
type MockSanitizer struct {
	// Allow tests to override behavior.
	SanitizeFunc func(query string) (string, error)
}

func (ms *MockSanitizer) SanitizeSearchQuery(query string) (string, error) {
	if query == "<script>alert('trigger xss sanitizer error');</script>" {
		return "", errors.New("sanitizer failure")
	}

	return query, nil
}

// FakeValidator implements interfaces.SearchValidator.
type MockValidator struct {
	ValidateFunc func(query searchdef.SearchQuery) error
}

func (mv *MockValidator) ValidateQuery(query searchdef.SearchQuery) error {
	if mv.ValidateFunc != nil {
		return mv.ValidateFunc(query)
	}
	return nil
}

// FakeIGDBAdapter implements interfaces.IGDBAdapter.
type MockIGDBAdapter struct {
	SearchGamesFunc func(ctx context.Context, query string, limit int) ([]*models.Game, error)
	UpdateTokenFunc func(token string) error
}

func (mv *MockIGDBAdapter) SearchGames(ctx context.Context, query string, limit int) ([]*models.Game, error) {
	if mv.SearchGamesFunc != nil {
		return mv.SearchGamesFunc(ctx, query, limit)
	}
	return nil, errors.New("SearchGamesFunc not defined")
}

func (mv *MockIGDBAdapter) UpdateToken(token string) error {
	if mv.UpdateTokenFunc != nil {
		return mv.UpdateTokenFunc(token)
	}
	return nil
}

// FakeCacheWrapper implements interfaces.IGDBCacheWrapper.
type MockCacheWrapper struct {
	GetCachedResultsFunc func(ctx context.Context, sq searchdef.SearchQuery) (*searchdef.SearchResult, error)
	SetCachedResultsFunc func(ctx context.Context, sq searchdef.SearchQuery, result *searchdef.SearchResult) error
	TimeToLive           int // or time.Duration
	InvalidateCacheFunc  func(ctx context.Context, cacheKey string) error
}

func (mcw *MockCacheWrapper) GetCachedResults(ctx context.Context, sq searchdef.SearchQuery) (*searchdef.SearchResult, error) {
	if mcw.GetCachedResultsFunc != nil {
		return mcw.GetCachedResultsFunc(ctx, sq)
	}
	return nil, nil
}

func (mcw *MockCacheWrapper) SetCachedResults(ctx context.Context, sq searchdef.SearchQuery, result *searchdef.SearchResult) error {
	if mcw.SetCachedResultsFunc != nil {
		return mcw.SetCachedResultsFunc(ctx, sq, result)
	}
	return nil
}

func (mcw *MockCacheWrapper) InvalidateCache(ctx context.Context, cacheKey string) error {
	if mcw.InvalidateCacheFunc != nil {
		return mcw.InvalidateCacheFunc(ctx, cacheKey)
	}
	return nil
}