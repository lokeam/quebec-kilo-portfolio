package worker

import (
	memcache "github.com/lokeam/qko-beta/internal/infrastructure/cache/memorycache"
	cache "github.com/lokeam/qko-beta/internal/infrastructure/cache/rueidis"
)

type CacheClients struct {
	RedisClient *cache.RueidisClient
	MemCache    *memcache.MemoryCache
}