package search

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/types"
)

// igdbCacheItem holds cached IGDB game data along with its expiration time.
type igdbCacheItem struct {
	games       []*types.Game
	expiration  time.Time
}

// IGDBCacheWrapper manages caching of IGDB search results using Redis.
// It first tries to return a cached result from Redis and falls back to fetching
// fresh data when needed. Redis calls use a configurable timeout.
type IGDBCacheWrapper struct {
	cacheClient  interfaces.CacheClient // Redis client.
	timeToLive   time.Duration        // How long items should stay in Redis.
	redisTimeout time.Duration        // Configurable timeout for Redis calls.
	logger       interfaces.Logger    // Logs for debugging and errors.
}

// NewIGDBCacheWrapper creates a new IGDBCacheWrapper.
// redisTimeout can be configured (e.g., via environment variables) so that Redis operations can have a shorter timeout.
func NewIGDBCacheWrapper(
	cacheClient interfaces.CacheClient,
	timeToLive time.Duration,
	redisTimeout time.Duration,
	logger interfaces.Logger,
) (*IGDBCacheWrapper, error) {
	return &IGDBCacheWrapper{
		cacheClient:  cacheClient,
		timeToLive:   timeToLive,
		redisTimeout: redisTimeout,
		logger:       logger,
	}, nil
}

// GetSearchResults first attempts to return a cached search result from Redis.
// It uses a short Redis-specific timeout. If no cached data is found or if there
// is an error, it fetches fresh data from IGDB via the dataRetriever.
func (d *IGDBCacheWrapper) GetCachedResults(
	ctx context.Context,
	query searchdef.SearchQuery,
) (*searchdef.SearchResult, error) {
	start := time.Now()
	sq := searchdef.SearchQuery{
		Query: query.Query,
		Limit: query.Limit,
	}
	cacheKey := sq.ToCacheKey()

	d.logger.Debug("attempting cache lookup", map[string]any{
		"key":   cacheKey,
		"query": query.Query,
	})

	redisCtx, cancel := context.WithTimeout(ctx, d.redisTimeout)
	defer cancel()

	cached, err := d.cacheClient.Get(redisCtx, cacheKey)
	if err != nil {
		d.logger.Warn("Cache GET failed or timed out", map[string]any{
			"key":   cacheKey,
			"error": err.Error(),
		})
		return nil, err
	}

	var result searchdef.SearchResult
	if err := json.Unmarshal([]byte(cached), &result); err == nil {
		result.Meta.CacheHit = true
		d.logger.Debug("Cache hit", map[string]any{
			"key":      cacheKey,
			"duration": time.Since(start).String(),
		})
		return &result, nil
	}

	d.logger.Warn("Cache unmarshal failed", map[string]any{
		"key":   cacheKey,
		"error": err,
	})

	return nil, nil
}

// SetCachedResults stores search results in Redis.
func (d *IGDBCacheWrapper) SetCachedResults(
	ctx context.Context,
	query searchdef.SearchQuery,
	result *searchdef.SearchResult,
) error {
	cacheKey := query.ToCacheKey()
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	if err := d.cacheClient.Set(ctx, cacheKey, string(data), d.timeToLive); err != nil {
		d.logger.Error("failed to set cache", map[string]any{
			"key":   cacheKey,
			"error": err.Error(),
		})
		return err
	}
	return nil
}