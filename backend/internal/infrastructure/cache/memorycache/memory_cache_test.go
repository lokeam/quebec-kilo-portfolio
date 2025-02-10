package cache

import (
	"context"
	"sync"
	"testing"
	"time"
)

/*
	Behaviors:
		- Constructor:
			* MemCache constructor should creates a map of items. This map should not be nil and the map should be empty

		- Set() method:
			* The Set() method takes in a key, a value in and a TimeToLive. It stores an item where the expiration is set to current time + TTL.

		- Get() method:
			* The Get() method takes in a key, returns a value a value and an error
			* The Get() method first requires a read lock, then checks if the items exists AND whether or not it has expired
			* If the item EXISTS AND has not exired, Get() returns the value
			* If the item DOES NOT EXIST OR has expired, Get() returns an empty string

	Scenarios:
		- Constructor should return a cache that is not nil and an empty map
		- Both Set() and Get() use sync.Mutex and both methods should be thread safe (ie: cache should not behave unexpectedly or cause data races)
		- Calling Get() on a key that does not exist should return an empty string
		- After a key is set, calling Get() should return the stored value
		- If the expiration time of a key is in the past, the key should not be retreived
*/

func TestMemoryCache(t *testing.T) {}

func TestMemoryCacheBehaviors(t *testing.T) {
	ctx := context.Background()

	// New MemoryCache
	t.Run(
		`GIVEN we create a new memory cache`,
		func(t *testing.T) {
			testMemCache, err := NewMemoryCache()
			if err != nil {
				t.Fatalf("we expected an empty cache, but we got an error: %v", err)
			}
			if testMemCache == nil {
				t.Fatalf("we expected a non nil memory cache, but memory cache is nil")
			}

			if testMemCache == nil {
				t.Fatalf("we expected the cache's items map to initialized")
			}

			// --------- Get() Method, TTL not set ---------
			// Store an item with proper value and expiration time
			t.Run(
				`WHEN a memcache key isn't set, THEN Get() returns an empty string`,
				func(t *testing.T) {
					value, err := testMemCache.Get(ctx, "key")
					if err != nil {
						t.Fatalf("we expected no error, but we got an error: %v", err)
					}
					if value != "" {
						t.Errorf("we expected to receive an empty string for a missing key, but we got: '%s'", value)
					}
				},
			)

			// --------- Get() Method, Happy Path ---------
			t.Run(
				`WHEN a memcache key is saved with a valid TTL, THEN the Get() method returns the value`,
				func(t *testing.T) {
					memCacheKey := "test-key-valid-ttl"
					expectedValue := "test-value-valid-ttl"
					testTTL := 500 * time.Millisecond

					if err := testMemCache.Set(ctx, memCacheKey, expectedValue, testTTL); err != nil {
						t.Fatalf("calling the Set() method returned this error: %v", err)
					}

					// Immediately retrieve the value
					retrievedValue, err := testMemCache.Get(ctx, memCacheKey)
					if err != nil {
						t.Fatalf("calling the Get() method returned this error: %v", err)
					}
					if retrievedValue != expectedValue {
						t.Errorf("we expected this key: %q to have expired and return an empty string but got this value: %q", expectedValue, retrievedValue)
					}
				},
			)

			// --------- Get() Method, Expired Key ---------
			t.Run(
				`WHEN a memCache key is set with a short TTL and then expires, THEN the Get() method returns an empty string`,
				func(t *testing.T) {
					memCacheKey := "test-key-expired-ttl"
					expectedValue := "test-value-expired-ttl"
					testTTL := 50 * time.Millisecond

					if err := testMemCache.Set(ctx, memCacheKey, expectedValue, testTTL); err != nil {
						t.Fatalf("calling the Set() method returns this error: %v", err)
					}

					// Wait for the TTL to expire
					time.Sleep(100 * time.Millisecond)

					expectedValue, err := testMemCache.Get(ctx, memCacheKey)
					if err != nil {
						t.Fatalf("calling the Get() method returns hits error: %v", err)
					}
					if expectedValue != "" {
						t.Errorf("we expected this key: %q to have expired and return an empty string, but we got this value:%q", memCacheKey, expectedValue)
					}
				},
			)

			// --------- Thread safety ---------
			// When multiple goroutines call Get or Set at the same time, nothing catches on fire
			t.Run(
				`WHEN multiple goroutines access the cache concurrently, THEN the cache should not behave unexpectedly or cause data races`,
				func(t *testing.T) {
					memCacheKey := "test-key-thread-safety"
					expectedValue := "test-value-thread-safety"
					testTTL := 500 * time.Millisecond

					var testWaitGroup sync.WaitGroup

					// Fire off many concurrent goroutines to set and get the same key
					for i := 0; i < 100; i++ {
						testWaitGroup.Add(2)
						go func() {
							defer testWaitGroup.Done()
							_ = testMemCache.Set(ctx, memCacheKey, expectedValue, testTTL)
						}()
						go func() {
							defer testWaitGroup.Done()
							_, _ = testMemCache.Get(ctx, memCacheKey)
						}()
					}

					testWaitGroup.Wait()

					actualValue, err := testMemCache.Get(ctx, memCacheKey)
					if err != nil {
						t.Fatalf("calling the Get() method returned this error: %v", err)
					}

					// TTL is set long enough, the key should still exist
					if expectedValue == "" {
						t.Errorf("we expected this value: %q to still exist, but we got this value instead: %q", expectedValue, actualValue)
					}
				},
			)
		},
	)
}
