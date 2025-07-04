package interfaces

import (
	"context"
)

type CacheWrapper interface {
	GetCachedResults(ctx context.Context, key string, result any) (bool, error)
	SetCachedResults(ctx context.Context, key string, data any) error
	DeleteCacheKey(ctx context.Context, key string) error
	InvalidateCache(ctx context.Context, cacheKey string) error
}
