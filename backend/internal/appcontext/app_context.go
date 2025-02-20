package appcontext

import (
	"github.com/lokeam/qko-beta/config"
	memcache "github.com/lokeam/qko-beta/internal/infrastructure/cache/memorycache"
	cache "github.com/lokeam/qko-beta/internal/infrastructure/cache/rueidis"
	"github.com/lokeam/qko-beta/internal/interfaces"
	twitch "github.com/lokeam/qko-beta/internal/shared/twitch"
)

type AppContext struct {
	Config                 *config.Config
	Logger                 interfaces.Logger
	MemCache               *memcache.MemoryCache
	RedisClient            *cache.RueidisClient
	TwitchTokenRetriever   *twitch.TwitchTokenRetriever
}

func NewAppContext(
	config *config.Config,
	logger interfaces.Logger,
	memCache *memcache.MemoryCache,
	redisClient *cache.RueidisClient,
) *AppContext {
	return &AppContext{
		Config: config,
		Logger: logger,
		MemCache: memCache,
		RedisClient: redisClient,
		TwitchTokenRetriever: twitch.NewTwitchTokenRetriever(
			memCache,
			redisClient,
			config.IGDB.ClientID,
			config.IGDB.ClientSecret,
			config.IGDB.AuthURL,
		),
	}
}
