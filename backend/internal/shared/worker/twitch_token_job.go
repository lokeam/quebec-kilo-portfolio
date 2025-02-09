package worker

import (
	"context"
	"time"

	memorycache "github.com/lokeam/qko-beta/internal/infrastructure/cache/memorycache"
	cache "github.com/lokeam/qko-beta/internal/infrastructure/cache/rueidis"
	"github.com/lokeam/qko-beta/internal/shared/logger"
	"github.com/lokeam/qko-beta/internal/shared/token"
)

func UpdateTwitchTokenJob(
	ctx context.Context,
	redisKey string,
	redisClient *cache.RueidisClient,
	memCache *memorycache.MemoryCache,
	tokenInfo token.TokenInfo,
	logger *logger.Logger,
) error {
	// Save token in memory cache
	logger.Info("Saving token in memory cache", nil)
	if err := memCache.Set(
		ctx,
		redisKey,
		tokenInfo.AccessToken,
		time.Until(tokenInfo.ExpiresAt),
		); err != nil {
			logger.Error("Failed to save token in memory cache", map[string]any{"error": err})
			// Continue even if memory cache fails
	}
	logger.Debug("Token saved in memory cache", map[string]any{"key": redisKey, "ttl": time.Until(
		tokenInfo.ExpiresAt).String(),
	})
	// Save token in Redis
	if err := token.SaveTokenInRedis(
		ctx,
		redisKey,
		redisClient,
		tokenInfo,
		time.Until(tokenInfo.ExpiresAt),
			*logger,
	); err != nil {
		logger.Error("Failed to save token in Redis", map[string]any{"error": err})
		return err
	}
	logger.Debug("Token saved in Redis", nil)
	return nil
}