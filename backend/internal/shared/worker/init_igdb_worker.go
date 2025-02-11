package worker

import (
	"context"

	"github.com/lokeam/qko-beta/internal/shared/logger"
)

// StartInitIGDBJob launches the INIT_IGDB job as a goroutine.
func StartInitIGDBJob(
	ctx context.Context,
	redisKey string,
	caches *CacheClients,
	clientID,
	clientSecret,
	authURL string,
	log logger.LoggerInterface,
) {
	go func() {
		if err := InitIGDBJob(
			ctx,
			redisKey,
			caches.RedisClient,
			caches.MemCache,
			clientID,
			clientSecret,
			authURL,
			log,
		); err != nil {
			log.Error("INIT_IGDB job terminated with error", map[string]any{"error": err.Error()})
		}
	}()
}
