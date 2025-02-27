package interfaces

import (
	"context"
)

type CacheWrapper interface {
	GetCachedResults(ctx context.Context, key string, result any) (bool, error)
	SetCachedResults(ctx context.Context, key string, data any) error
}
