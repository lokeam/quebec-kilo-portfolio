package analytics

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/interfaces"
)

// AnalyticsCacheAdapter wraps the CacheWrapper to provide domain-specific caching for analytics
type AnalyticsCacheAdapter struct {
	cacheWrapper interfaces.CacheWrapper
}

// NewAnalyticsCacheAdapter creates a new analytics cache adapter
func NewAnalyticsCacheAdapter(cacheWrapper interfaces.CacheWrapper) (*AnalyticsCacheAdapter, error) {
	if cacheWrapper == nil {
		return nil, fmt.Errorf("cache wrapper cannot be nil")
	}

	return &AnalyticsCacheAdapter{
		cacheWrapper: cacheWrapper,
	}, nil
}

// Implement CacheWrapper interface methods
func (aca *AnalyticsCacheAdapter) GetCachedResults(ctx context.Context, key string, result any) (bool, error) {
	return aca.cacheWrapper.GetCachedResults(ctx, key, result)
}

func (aca *AnalyticsCacheAdapter) SetCachedResults(ctx context.Context, key string, data any) error {
	return aca.cacheWrapper.SetCachedResults(ctx, key, data)
}

func (aca *AnalyticsCacheAdapter) DeleteCacheKey(ctx context.Context, key string) error {
	return aca.cacheWrapper.DeleteCacheKey(ctx, key)
}

// Domain-specific methods
// GetCachedAnalytics attempts to retrieve analytics data from cache
func (aca *AnalyticsCacheAdapter) GetCachedAnalytics(ctx context.Context, userID string, domain string) (any, bool, error) {
	cacheKey := fmt.Sprintf(CacheKeyFormat, userID, domain)
	var result any
	found, err := aca.cacheWrapper.GetCachedResults(ctx, cacheKey, &result)
	if err != nil {
		return nil, false, err
	}
	return result, found, nil
}

// SetCachedAnalytics stores analytics data in cache
func (aca *AnalyticsCacheAdapter) SetCachedAnalytics(ctx context.Context, userID string, domain string, data any) error {
	cacheKey := fmt.Sprintf(CacheKeyFormat, userID, domain)
	return aca.cacheWrapper.SetCachedResults(ctx, cacheKey, data)
}

// InvalidateDomain removes a specific domain from cache for a user
func (aca *AnalyticsCacheAdapter) InvalidateDomain(ctx context.Context, userID string, domain string) error {
	cacheKey := fmt.Sprintf(CacheKeyFormat, userID, domain)
	return aca.cacheWrapper.DeleteCacheKey(ctx, cacheKey)
}

// InvalidateDomains removes multiple domains from cache for a user
func (aca *AnalyticsCacheAdapter) InvalidateDomains(ctx context.Context, userID string, domains []string) error {
	var errs []error

	// Use index-based loop for better safety
	for i := 0; i < len(domains); i++ {
		domain := domains[i]
		err := aca.InvalidateDomain(ctx, userID, domain)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to invalidate some domains: %v", errs)
	}
	return nil
}

// InvalidateCache invalidates a specific cache key
func (aca *AnalyticsCacheAdapter) InvalidateCache(ctx context.Context, cacheKey string) error {
	return aca.cacheWrapper.DeleteCacheKey(ctx, cacheKey)
}