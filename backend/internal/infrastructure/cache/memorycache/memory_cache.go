package cache

import (
	"context"
	"sync"
	"time"
)

// Structs / Interfaces
type MemoryCache struct {
	mu      sync.RWMutex
	items   map[string]cacheItem
}

type cacheItem struct {
	value       string
	expiration  time.Time
}


// Constructor
func NewMemoryCache() (*MemoryCache, error) {
	return &MemoryCache{
		items: make(map[string]cacheItem),
	}, nil
}


// Retrieves an item if it exists + is not expired
func (mc *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, found := mc.items[key]

	if !found || time.Now().After(item.expiration) {
		return "", nil
	}

	return item.value, nil
}

// Stores an item with a given TTL
func (mc *MemoryCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.items[key] = cacheItem{
		value:       value,
		expiration:  time.Now().Add(ttl),
	}

	return nil
}
