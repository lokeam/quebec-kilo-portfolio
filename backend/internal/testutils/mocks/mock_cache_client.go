package mocks

import (
	"context"
	"time"
)

// DummyCacheClient is a stub implementation of the cache client.
type MockCacheClient struct {
	GetFunc    func(ctx context.Context, key string) (string, error)
	SetFunc    func(ctx context.Context, key string, value any, ttl time.Duration) error
	DeleteFunc func(ctx context.Context, key string) error
}

func (d *MockCacheClient) Get(ctx context.Context, key string) (string, error) {
	return d.GetFunc(ctx, key)
}

func (d *MockCacheClient) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return d.SetFunc(ctx, key, value, ttl)
}

func (d *MockCacheClient) Delete(ctx context.Context, key string) error {
	return d.DeleteFunc(ctx, key)
}
