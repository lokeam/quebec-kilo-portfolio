package redisclient

import (
	"context"
	"time"
)

type RedisClient interface {
	Ping(ctx context.Context) error
	IsReady() bool
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}
