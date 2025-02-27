package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
)

type CacheWrapper struct {
	cacheClient     interfaces.CacheClient
	timeToLive      time.Duration
	redisTimeout    time.Duration
	logger          interfaces.Logger
}

func NewCacheWrapper(
	cacheClient interfaces.CacheClient,
	timeToLive time.Duration,
	redisTimeout time.Duration,
	logger interfaces.Logger,
) (*CacheWrapper, error) {
	// Validate time to live value
	if timeToLive <= 0 {
		const defaultTimeToLive = 10 * time.Minute
		logger.Warn("Cache Wrapper - Invalid TTL provided, using default", map[string]any{
			"providedTTL": timeToLive,
			"defaultTTL":  defaultTimeToLive,
		})
		timeToLive = defaultTimeToLive
	}

	return &CacheWrapper{
		cacheClient:   cacheClient,
		timeToLive:    timeToLive,
		redisTimeout:  redisTimeout,
		logger:        logger,
	}, nil
}

func (cw *CacheWrapper) GetCachedResults(
	ctx context.Context,
	cacheKey string,
	result any,
) (bool, error) {
	// Validate cache key
	if cacheKey == "" {
		cw.logger.Error("Cache Wrapper - Cache key cannot be empty", map[string]any{
			"cacheKey": cacheKey,
		})
		return false, errors.New("cache key cannot be empty")
	}

	start := time.Now()

	cw.logger.Info("Cache Wrapper - attempting cache lookup", map[string]any{
		"cacheKey": cacheKey,
		"start":    start,
	})

	redisContext, cancel := context.WithTimeout(ctx, cw.redisTimeout)
	defer cancel()

	cachedData, err := cw.cacheClient.Get(redisContext, cacheKey)
	if err != nil {
		cw.logger.Error("Cache Wrapper - Cache GET failed or timed out", map[string]any{
			"cacheKey": cacheKey,
			"error":    err,
		})
		return false, err
	}

	if err := json.Unmarshal([]byte(cachedData), result); err == nil {
		cw.logger.Debug("Cache Wrapper - Cache hit", map[string]any{
			"cacheKey": cacheKey,
			"duration": time.Since(start),
		})
		return true, nil
	}

	cw.logger.Warn("Cache Wrapper - Cache unmarshalling failed", map[string]any{
		"cacheKey": cacheKey,
		"duration": time.Since(start),
	})

	return false, nil
}

func (cw *CacheWrapper) SetCachedResults(
	ctx context.Context,
	cacheKey string,
	data any,
) error {
	// Validate cache key
	if cacheKey == "" {
		cw.logger.Error("Cache Wrapper - Cache key cannot be empty", map[string]any{
			"cacheKey": cacheKey,
		})
		return errors.New("cache key cannot be empty")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := cw.cacheClient.Set(ctx, cacheKey, string(jsonData), cw.timeToLive); err != nil {
		cw.logger.Error("Cache Wrapper - Cache SET failed", map[string]any{
			"cacheKey": cacheKey,
			"error":    err,
		})
		return err
	}

	cw.logger.Info("Cache Wrapper - Cache SET successful", map[string]any{
		"cacheKey": cacheKey,
	})

	return nil
}
