package analytics

import (
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/infrastructure/cache"
)

// NewAnalyticsService creates and initializes the analytics service with all required dependencies
func NewAnalyticsService(appCtx *appcontext.AppContext) (Service, error) {
	// Create the analytics DB adapter
	analyticsDbAdapter, err := NewAnalyticsDbAdapter(appCtx)
	if err != nil {
		appCtx.Logger.Error("Failed to create analytics DB adapter", map[string]any{
			"error": err,
		})
		return nil, err
	}

	// Create the cache wrapper with configuration from appContext
	cacheWrapper, err := cache.NewCacheWrapper(
		appCtx.RedisClient,
		appCtx.Config.Redis.RedisTTL,
		appCtx.Config.Redis.RedisTimeout,
		appCtx.Logger)
	if err != nil {
		appCtx.Logger.Error("Failed to create cache wrapper", map[string]any{
			"error": err,
		})
		return nil, err
	}

	// Create analytics cache adapter which implements CacheWrapper
	analyticsCacheAdapter, err := NewAnalyticsCacheAdapter(cacheWrapper)
	if err != nil {
		appCtx.Logger.Error("Failed to create analytics cache adapter", map[string]any{
			"error": err,
		})
		return nil, err
	}

	// Initialize the analytics service with adapter and cache
	return NewService(analyticsDbAdapter, analyticsCacheAdapter, appCtx.Logger), nil
}