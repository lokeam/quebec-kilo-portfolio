package resourceinitializer

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/config"
	memcache "github.com/lokeam/qko-beta/internal/infrastructure/cache/memorycache"
	rueidis "github.com/lokeam/qko-beta/internal/infrastructure/cache/rueidis"
	logger "github.com/lokeam/qko-beta/internal/shared/logger"
)

type ResourceInitializer struct {
	// Cache
	RedisClient           *rueidis.RueidisClient

	// Memcahe
	MemCache              *memcache.MemoryCache

	// Postgres - TBD

	// Handlers - TBD
}

func NewResourceInitializer(
	ctx context.Context,
	config *config.Config,
	logger *logger.Logger,
	) (*ResourceInitializer, error) {
		// Initialize Redis
		cfg := rueidis.NewRueidisConfig()
		if err := cfg.LoadFromEnv(); err != nil {
			return nil, fmt.Errorf("redis config error: %w", err)
		}
		if err := cfg.Validate(); err != nil {
			return nil, fmt.Errorf("invalid redis config %w", err)
		}

		// Create Redis Client
		redisClient, err := rueidis.NewRueidisClient(cfg, *logger)
		if err != nil {
			logger.Error("redis client error", map[string]any{
				"error": err.Error(),
			})
			return nil, err
		}

		// Verify Redis connection with Ping
		if err := redisClient.Ping(ctx); err != nil {
			return nil, fmt.Errorf("initial redis ping failed: %w", err)
		}

		// Initialize Memcache
		memCache, err := memcache.NewMemoryCache()
		if err != nil {
			return nil, fmt.Errorf("memcache client error: %w", err)
		}


		return &ResourceInitializer{
			RedisClient: redisClient,
			MemCache:    memCache,
		}, nil
}

