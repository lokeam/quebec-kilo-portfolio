package token

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cache "github.com/lokeam/qko-beta/internal/infrastructure/cache/rueidis"
	"github.com/lokeam/qko-beta/internal/shared/logger"
)

type TokenInfo struct {
	AccessToken string
	ExpiresAt   time.Time
}

// SaveTokenInRedis saves the token info in Redis.
func SaveTokenInRedis(
	ctx context.Context,
	key string,
	redisClient *cache.RueidisClient,
	tokenInfo TokenInfo,
	ttl time.Duration,
	log logger.Logger,
) error {
	data, err := json.Marshal(tokenInfo)
	if err != nil {
		log.Error("Failed to marshal token info", map[string]any{"error": err.Error()})
		return err
	}
	err = redisClient.Set(ctx, key, string(data), ttl)
	if err != nil {
		log.Error("Failed to save token in Redis", map[string]any{"error": err.Error()})
		return fmt.Errorf("failed to save token in Redis: %w", err)
	}
	log.Info("Token saved in Redis", map[string]any{"key": key, "ttl": ttl.String()})
	return nil
}
