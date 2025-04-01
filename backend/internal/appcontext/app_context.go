package appcontext

import (
	"github.com/lokeam/qko-beta/config"
	memcache "github.com/lokeam/qko-beta/internal/infrastructure/cache/memorycache"
	cache "github.com/lokeam/qko-beta/internal/infrastructure/cache/rueidis"
	interfaces "github.com/lokeam/qko-beta/internal/interfaces"
	twitch "github.com/lokeam/qko-beta/internal/shared/twitch"
)

type AppContext struct {
	Config                 *config.Config
	Logger                 interfaces.Logger
	MemCache               *memcache.MemoryCache
	RedisClient            *cache.RueidisClient
	TwitchTokenRetriever   interfaces.TokenRetriever
}

func NewAppContext(
	config *config.Config,
	logger interfaces.Logger,
	memCache *memcache.MemoryCache,
	redisClient *cache.RueidisClient,
) *AppContext {

	twitchTokenRetriever := twitch.NewTwitchTokenRetriever(
		memCache,
		redisClient,
		config.IGDB.ClientID,
		config.IGDB.ClientSecret,
		config.IGDB.AuthURL,
	)

	return &AppContext{
		Config: config,
		Logger: logger,
		MemCache: memCache,
		RedisClient: redisClient,
		TwitchTokenRetriever: twitchTokenRetriever,
	}
}
