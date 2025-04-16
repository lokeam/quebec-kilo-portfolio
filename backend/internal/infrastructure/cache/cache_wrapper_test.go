package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Value     bool    `json:"value"`
}

func TestCacheWrapper(t *testing.T) {
	testLogger := testutils.NewTestLogger()

	t.Run(`NewCacheWrapper() correctly initializes itself`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheWrapper
			AND a properly configured CacheClient
			WHEN NewCacheWrapper() is called
			THEN NewCacheWrapper returns a properly initialized CachedWrapper
		*/

		testTTL := 10 * time.Minute
		testTimeout := 100 * time.Millisecond
		testCache := &mocks.MockCacheClient{
			GetFunc: func(ctx context.Context, key string) (string, error) {
				return "", nil
			},
			SetFunc: func(ctx context.Context, key string, value any, ttl time.Duration) error {
				return nil
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		cacheWrapper, testErr := NewCacheWrapper(testCache, testTTL, testTimeout, testLogger)

		assert.NoError(t, testErr)
		assert.NotNil(t, cacheWrapper)
		assert.Equal(t, testTTL, cacheWrapper.timeToLive)
		assert.Equal(t, testTimeout, cacheWrapper.redisTimeout)
		assert.Equal(t, testLogger, cacheWrapper.logger)
	})

	// ------ GetCachedResults() ------
	t.Run(`GetCachedResults() method returns a cache miss when the Redis connection fails or times out`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheWrapper
			AND a properly configured CacheClient that simulates a Redis connection failure
			AND a valid cache key that links to a cache value that exists
			WHEN GetCachedResults() is called with the key
			THEN GetCachedResults() should return an error from the Redis client indicating a cache miss
			AND no data should be returned along with the error
		*/

		expectedError := errors.New("cache not reachable")
		testCache := &mocks.MockCacheClient{
			GetFunc: func(ctx context.Context, key string) (string, error) {
				return "", expectedError
			},
			SetFunc: func(ctx context.Context, key string, value any, ttl time.Duration) error {
				return nil
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		testCacheWrapper, testErr := NewCacheWrapper(
			testCache,
			10 * time.Minute, // Test TTL
			100 * time.Millisecond, // Redis timeout
			testLogger,
		)
		assert.NoError(t, testErr)

		var result testData
		cacheHit, testErr := testCacheWrapper.GetCachedResults(
			context.Background(),
			"test:key",
			&result,
		)

		assert.Error(t, testErr)
		assert.False(t, cacheHit)
	})

	t.Run(`GetCachedResults() method returns an unmarshalling failure`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheWrapper
			AND a properly configured CacheClient that simulates an unmarshalling failure (such as invalid JSON)
			AND a valid cache key with a cache value that exists
			WHEN GetCachedResults() is called with the key
			THEN GetCachedResults() should return an error from the Redis client indicating a cache miss
			AND no data should be returned
			AND a warning should be logged
		*/

		testCache := &mocks.MockCacheClient{
			GetFunc: func(ctx context.Context, key string) (string, error) {
				return "invalid JSON", nil
			},
			SetFunc: func(ctx context.Context, key string, value any, ttl time.Duration) error {
				return nil
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		testCacheWrapper, testErr := NewCacheWrapper(
			testCache,
			10 * time.Minute, // Test TTL
			100 * time.Millisecond, // Redis timeout
			testLogger,
		)
		assert.NoError(t, testErr)

		var testResult testData
		cacheHit, testErr := testCacheWrapper.GetCachedResults(
			context.Background(),
			"test:key",
			&testResult,
		)

		assert.NoError(t, testErr)
		assert.False(t, cacheHit)
		assert.Greater(t, len(testLogger.WarnCalls), 0)
	})

	// Context cancellation test
	t.Run(`GetCachedResults() method handles context cancellation`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheWrapper
			AND a CacheClient that will delay long enough for the context to be cancelled
			WHEN GetCachedResults() is called
			THEN the operation should be aborted
			AND we should return a context cancellation error
		*/

		// Create a mock that delays until context is cancelled
		blockingCache := &mocks.MockCacheClient{
			GetFunc: func(ctx context.Context, key string) (string, error) {
				// Block until context is cancelled
				<-ctx.Done()
				return "", ctx.Err()
			},
			SetFunc: func(ctx context.Context, key string, value any, ttl time.Duration) error {
				return nil
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		testCacheWrapper, testErr := NewCacheWrapper(
			blockingCache,
			10 * time.Minute, // Test TTL
			100 * time.Millisecond, // Redis timeout
			testLogger,
		)
		assert.NoError(t, testErr)

		// Create cancellable context
		testCtx, cancel := context.WithCancel(context.Background())

		// Setup a channel to signal when goroutine has started
		goRoutineStarted := make(chan struct{})

		// Setup a channel to collect results
		resultChannel := make(chan error)

		// Call GetCachedResults() in a goroutine
		go func() {
			var result testData
			close(goRoutineStarted) // NOTE: signal that goroutine has started
			_, err := testCacheWrapper.GetCachedResults(testCtx, "test:key", &result)
			resultChannel <- err
		}()

		// For the routine to start
		<-goRoutineStarted

		// Cancel the context after a short delay
		time.Sleep(10 * time.Millisecond)
		cancel()

		// Wait for the result with a timeout
		select {
		case err := <- resultChannel:
			assert.Error(t, err)

			/*
				NOTE: context.Canceled is the error type
				(ie: we need to specifically listen for "context canceled" as opposed to "context cancelled")
			*/
			assert.Contains(t, err.Error(), "context canceled", "Error should indicate context cancellation")
		case <-time.After(200 *time.Millisecond):
			t.Fatal("Test timed out - GetCachedResults() did not respect context cancellation")
		}
	})

	t.Run(`Happy Path: GetCachedResults() method returns a cache hit`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheWrapper
			AND a properly configured CacheClient
			AND a properly configured testData struct
			AND a properly configured CacheKey
			WHEN GetCachedResults() is called
			THEN GetCachedResults() returns a cache hit and the data is returned successfully
		*/
		expectedResult := &testData{
			ID: 1,
			Name: "Dark Souls",
			Value: true,
		}
		jsonData, testErr := json.Marshal(expectedResult)
		assert.NoError(t, testErr)

		testCache := &mocks.MockCacheClient{
			GetFunc: func(ctx context.Context, key string) (string, error) {
				return string(jsonData), nil
			},
			SetFunc: func(ctx context.Context, key string, value any, ttl time.Duration) error {
				return nil
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		testCacheWrapper, testErr := NewCacheWrapper(
			testCache,
			10 * time.Minute, // Test TTL
			100 * time.Millisecond, // Redis timeout
			testLogger,
		)
		assert.NoError(t, testErr)

		var result testData
		cacheHit, testErr := testCacheWrapper.GetCachedResults(
			context.Background(),
			"test:key",
			&result,
		)

		assert.NoError(t, testErr)
		assert.True(t, cacheHit)
		assert.Equal(t, expectedResult.ID, result.ID)
		assert.Equal(t, expectedResult.Name, result.Name)
		assert.Equal(t, expectedResult.Value, result.Value)
	})

	// ------ SetCachedResults() ------
	t.Run(`SetCachedResults() method handles context cancellation`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheWrapper
			AND a CacheClient that will delay long enough for the context to be cancelled
			WHEN SetCachedResults() is called with a context that gets cancelled
			THEN the opeation should be aborted
			AND the context cancellation error should be returned
		*/

		blockingCache := &mocks.MockCacheClient{
			SetFunc: func(ctx context.Context, key string, result any, ttl time.Duration) error {
				// Block until context cancelled
				<-ctx.Done()
				return ctx.Err()
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		testCacheWrapper, testErr := NewCacheWrapper(
			blockingCache,
			10 * time.Minute, // Test TTL
			100 * time.Millisecond, // Redis timeout
			testLogger,
		)
		assert.NoError(t, testErr)

		// Create context with short timeout
		testCtx, cancel := context.WithTimeout(context.Background(), 20 * time.Millisecond)
		defer cancel()

		// Call SetCachedResults() and verify that it respects the context
		testData := &testData{
			ID: 1,
			Name: "Dark Souls",
			Value: true,
		}
		ctxErr := testCacheWrapper.SetCachedResults(
			testCtx,
			"test:key",
			testData,
		)

		assert.Error(t, ctxErr)
		assert.Contains(t, ctxErr.Error(), "context deadline exceeded", "Error should indicate context timeout")
	})

	t.Run(`SetCachedResults() method experienced unmarshalling failure`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheWrapper
			AND a properly configured CacheClient that simulates an unmarshalling failure (such as invalid JSON)
			AND a properly configured testData struct
			AND a properly configured CacheKey
			WHEN SetCachedResults() is called
			THEN SetCachedResults() should return an error from the Redis client indicating a cache miss
			AND no data should be cached
			AND a warning should be logged
		*/

		expectedError := errors.New("cannot connect to Redis cache")
		testCache := &mocks.MockCacheClient{
			SetFunc: func(ctx context.Context, key string, result any, ttl time.Duration) error {
				return expectedError
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		testCacheWrapper, testErr := NewCacheWrapper(
			testCache,
			10 * time.Minute, // Test TTL
			100 * time.Millisecond, // Redis timeout
			testLogger,
		)
		assert.NoError(t, testErr)

		testData := &testData{
			ID: 1,
			Name: "Dark Souls",
			Value: true,
		}
		testErr = testCacheWrapper.SetCachedResults(
			context.Background(),
			"test:key",
			testData,
		)

		assert.Error(t, testErr)
		assert.Greater(t, len(testLogger.ErrorCalls), 0)
	})

	t.Run(`Happy Path: SetCachedResults() successfully caches a result`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheWrapper
			AND a properly configured CacheClient
			AND a properly configured testData struct
			AND a properly configured CacheKey
			WHEN SetCachedResults() is called
			THEN SetCachedResults() successfully caches the result
			AND returns no error
		*/

		testCache := &mocks.MockCacheClient{
			SetFunc: func(ctx context.Context, key string, result any, ttl time.Duration) error {
				return nil
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		testCacheWrapper, testErr := NewCacheWrapper(
			testCache,
			10 * time.Minute, // Test TTL
			100 * time.Millisecond, // Redis timeout
			testLogger,
		)
		assert.NoError(t, testErr)

		testData := &testData{
			ID: 1,
			Name: "Dark Souls",
			Value: true,
		}
		testErr = testCacheWrapper.SetCachedResults(
			context.Background(),
			"test:key",
			testData,
		)

		assert.NoError(t, testErr)
	})

	// ------ Edge cases ------
	t.Run(`CacheWrapper handles concurrent access - is thread safe`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheWrapper
			AND a properly configured CacheClient
			WHEN the CacheWrapper is accessed by many goroutines
			THEN all operations complete without race conditions or panics
			AND all operations produce the expected results
		*/

		// NOTE: We need to cound calls to detect race conditions
		var getCalls, setCalls int32
		testCache := &mocks.MockCacheClient{
			GetFunc: func(ctx context.Context, key string) (string, error) {
				// Atomic increment to avoid race condition in the counter itself
				atomic.AddInt32(&getCalls, 1)
				time.Sleep(5 * time.Millisecond) // Force some concurrency

				// Extract ID from they key for JSON return
				keyParts := strings.Split(key, "-")
				id := "0"
				if len(keyParts) > 1 {
					id = keyParts[1]
				}

				return fmt.Sprintf(`{"id": %s, "name": "Test", "value": true}`, id), nil
			},
			SetFunc: func(ctx context.Context, key string, value any, ttl time.Duration) error {
				atomic.AddInt32(&setCalls, 1)
				time.Sleep(5 * time.Millisecond) // Force some concurrency
				return nil
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		testCacheWrapper, testErr := NewCacheWrapper(
			testCache,
			10 * time.Minute, // Test TTL
			100 * time.Millisecond, // Redis timeout
			testLogger,
		)
		assert.NoError(t, testErr)

		// Create wait group to wait for all goroutines to finish
		testWaitGroup := sync.WaitGroup{}
		concurrentOperations := 10
		testWaitGroup.Add(concurrentOperations * 2) // NOTE: we need to account for both GET and SET ops

		// Run the ops in parallel
		for i := 0; i < concurrentOperations; i++ {
			go func(index int) {
				defer testWaitGroup.Done()
				var result testData
				key := fmt.Sprintf("key-%d", index)
				cacheHit, err := testCacheWrapper.GetCachedResults(
					context.Background(),
					key,
					&result,
				)
				assert.NoError(t, err)
				assert.True(t, cacheHit)
				assert.Equal(t, key, fmt.Sprintf("key-%d", result.ID))
			}(i)
		}

		for i := 0; i < concurrentOperations; i++ {
			go func(index int) {
				defer testWaitGroup.Done()
				key := fmt.Sprintf("key-%d", index)
				data := testData{
					ID: index,
					Name: "Test",
					Value: true,
				}
				err := testCacheWrapper.SetCachedResults(context.Background(), key, data)
				assert.NoError(t, err)
			}(i)
		}

		// Wait for the ops to complete
		testWaitGroup.Wait()

		// Verify that we did the things with the codes
		assert.Equal(t, int32(concurrentOperations), getCalls)
		assert.Equal(t, int32(concurrentOperations), setCalls)
	})

	t.Run(`Edge case test - CacheWrapper handles empty keys`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheWrapper
			WHEN GetCachedResults() OR SetCachedResults() is called with an empty key
			THEN we should return an appropriate error
			AND no cache operation should be attempted
		*/

		// Track if we called mock method
		var getCalled, setCalled bool
		testCache := &mocks.MockCacheClient{
			GetFunc: func(ctx context.Context, key string) (string, error) {
				getCalled = true
				return "", nil
			},
			SetFunc: func(ctx context.Context, key string, value any, ttl time.Duration) error {
				setCalled = true
				return nil
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		testCacheWrapper, testErr := NewCacheWrapper(
			testCache,
			10 * time.Minute, // Test TTL
			100 * time.Millisecond, // Redis timeout
			testLogger,
		)

		assert.NoError(t, testErr)

		// Test GetCachedResults with empty key
		var result testData
		cacheHit, testErr := testCacheWrapper.GetCachedResults(
			context.Background(),
			"",
			&result,
		)

		assert.Error(t, testErr)
		assert.False(t, cacheHit)
		assert.False(t, getCalled, "Cache client should not be called with empty key")

		// Test SetCachedResults with empty key
		testData := &testData{
			ID: 1,
			Name: "Test",
			Value: true,
		}
		testErr = testCacheWrapper.SetCachedResults(
			context.Background(),
			"",
			testData,
		)

		assert.Error(t, testErr)
		assert.False(t, setCalled, "Cache client should not be called with empty key")
	})

	t.Run(`Edge case test - CacheWrapper handles zero or negative TTL values for SetCachedResults()`, func(t *testing.T) {
		/*
			GIVEN a properly configured CacheClient
			WHEN a NewCacheWrapper is called with a zero or negative time to live value
			THEN an appropriate error should be logged OR a reasonable default value should be used
			AND the resulting wrapper should still work as expected
		*/

		testCache := &mocks.MockCacheClient{
			SetFunc: func(ctx context.Context, key string, value any, ttl time.Duration) error {
				// Make sure that a reasonable default is always used even if value is zero or negative
				assert.True(t, ttl > 0, "TTL should be a positive value")
				return nil
			},
			DeleteFunc: func(ctx context.Context, key string) error {
				return nil
			},
		}

		zeroTTLWrapper, testErr := NewCacheWrapper(
			testCache,
			0, // Zero TTL
			100 * time.Millisecond, // Redis timeout
			testLogger,
		)

		// Either we return an error OR use the default TTL value
		if testErr != nil {
			assert.Error(t, testErr)
			assert.Nil(t, zeroTTLWrapper)
		} else {
			assert.NotZero(t, zeroTTLWrapper.timeToLive, "Should use default TTL when zero provided")

			// Test that the wrapper still works as expected
			testData := &testData{
				ID: 1,
				Name: "Dark Souls",
				Value: true,
			}
			testErr = zeroTTLWrapper.SetCachedResults(
				context.Background(),
				"test:key",
				testData,
			)
			assert.NoError(t, testErr)
		}

		// Test using a negative TTL value
		negativeTTLWrapper, testErr := NewCacheWrapper(
			testCache,
			-5 * time.Minute, // Negative TTL
			100 * time.Millisecond, // Redis timeout,
			testLogger,
		)

		// Either we return an error OR a default value is used
		if testErr != nil {
			assert.Error(t, testErr)
			assert.Nil(t, negativeTTLWrapper)
		} else {
			assert.True(t, negativeTTLWrapper.timeToLive > 0, "Should use positive TTL value")
			// Test that the wrapper still works
			testData := &testData{
				ID: 1,
				Name: "Dark Souls",
				Value: true,
			}
			testErr = negativeTTLWrapper.SetCachedResults(
				context.Background(),
				"test:key",
				testData,
			)
			assert.NoError(t, testErr)
		}
	})
}
