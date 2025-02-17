package search

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

/*

	Behaviors:
		- Wrapper is created correctly
		- When cache hit occurs (example: cache returns valid JSON) the cache returns valid JSON, the result is unmarshalled and the CacheHit flag is set
		- When cache miss occurs (example: cache client returns an error) an error is returned
		- When the cached value is NOT valid JSON, the method logs a warning and returns no result
		- If JSON marshalling OR SetCachedResults() fails, the proper error is returned and logged

	Scenarios:
		- Cache hit

*/

func TestIGDBCacheWrapper(t *testing.T) {

	testLogger := testutils.NewTestLogger()

	t.Run(
		`NewIGDBCacheWrapper correctly intializes`,
		func(t *testing.T) {
			/*
				GIVEN a valid Redis client, a time-to-live value a redis timeout value and a logger
				WHEN NewIGDBCacheWrapper is called
				THEN a new IGDBCacheWrapper is returned with its fields properly set
			*/

			testTTL := 10 * time.Minute
			testTimeout := 100 * time.Millisecond
			testCache := &mocks.MockCacheClient{}

			igdbCacheWrapper, testErr := NewIGDBCacheWrapper(testCache, testTTL, testTimeout, testLogger)

			// NOTE: experimenting with the assert package instead of writing my own error msgs
			assert.NoError(t, testErr)
			assert.NotNil(t, igdbCacheWrapper)

			assert.Equal(t, testCache, igdbCacheWrapper.cacheClient, "cache client wasn't set correctly")
			assert.Equal(t, testTTL, igdbCacheWrapper.timeToLive, "ttl wasn't set correctly")
			assert.Equal(t, testTimeout, igdbCacheWrapper.redisTimeout, "redis timeout wasn't set correctly")
			assert.Equal(t, testLogger, igdbCacheWrapper.logger, "the logger wasn't configured correctly")
		},
	)

	t.Run(
		`GetCacheResults() Cache Hit`,
		func(t *testing.T) {
			/*
				GIVEN a cache client that returns valid JSON for a cached search result
				WHEN GetCachedResults() is called with a valid search query
				THEN the cached result is properly unmarshalled and marked as a cache hit
			*/
			expectedResult := &searchdef.SearchResult{
				Games: []searchdef.Game{{ID: 1, Name: "Dark Souls"}},
				Meta: searchdef.SearchMeta{
					CacheHit: true,
					CacheTTL: time.Duration(1 * time.Hour),
				},
			}
			testData, testErr := json.Marshal(expectedResult)
			assert.NoError(t, testErr)

			testCache := &mocks.MockCacheClient{
				GetFunc: func(ctx context.Context, key string) (string, error) {
					return string(testData), nil
				},
			}
			testIGDBCacheWrapper, testErr := NewIGDBCacheWrapper(testCache, 10 * time.Minute, 100 * time.Millisecond, testLogger)
			assert.NoError(t, testErr)

			testQuery := searchdef.SearchQuery{Query: "sample query", Limit: 10}
			testResult, testErr := testIGDBCacheWrapper.GetCachedResults(context.Background(), testQuery)

			assert.NoError(t, testErr)
			assert.NotNil(t, testResult)
			assert.True(t, testResult.Meta.CacheHit)
		},
	)

	t.Run(
		`GetCacheResults() Cache GET experiences failure`,
		func(t *testing.T) {
			/*
				GIVEN a cache client that returns an error
				WHEN GetCachedResults() is called with a valid search query
				THEN an error is returned and the result is nil
			*/
			expectedError := errors.New("cache not reachable")
			testCache := &mocks.MockCacheClient{
				GetFunc: func(ctx context.Context, key string) (string, error) {
					return "", expectedError
				},
			}
			testIGDBCacheWrapper, testErr := NewIGDBCacheWrapper(testCache, 10 * time.Minute, 100 * time.Millisecond, testLogger)
			assert.NoError(t, testErr)

			testQuery := searchdef.SearchQuery{Query: "sample query", Limit: 10}
			testResult, testErr := testIGDBCacheWrapper.GetCachedResults(context.Background(), testQuery)

			assert.Error(t, testErr)
			assert.Nil(t, testResult)
		},
	)

	t.Run(
		`GetCachedResults() experiences unmarshalling failure`,
		func(t *testing.T) {
			/*
				GIVEN a cache client that returns an invalid JSON string
				WHEN GetCachedResults() is called
				AND invalid JSON is returned
				THEN the method logs a warning and returns no result
			*/
			testCache := &mocks.MockCacheClient{
				GetFunc: func(ctx context.Context, key string) (string, error) {
					return "invalid JSON", nil
				},
			}

			testLogger := &testutils.TestLogger{}
			testIGDBCacheWrapper, testErr := NewIGDBCacheWrapper(testCache, 10 * time.Minute, 100 * time.Millisecond, testLogger)
			assert.NoError(t, testErr)

			testQuery := searchdef.SearchQuery{Query: "sample query", Limit: 10}
			expectedResult, testErr := testIGDBCacheWrapper.GetCachedResults(context.Background(), testQuery)

			assert.NoError(t, testErr)
			assert.Nil(t, expectedResult)
			assert.Greater(t, len(testLogger.WarnCalls), 0, "expected a warning on json unmarshalling failure")
		},
	)

	t.Run(
		`SetCachedResults() experiences Happy Path, successfully caches results`,
		func(t *testing.T) {
			/*
				GIVEN a cache client that successfully accepts Set calls
				WHEN SetCachedResults() is called with a valid search query
				THEN the results are properly marshalled and cached
			*/
			testCache := &mocks.MockCacheClient{
				SetFunc: func(ctx context.Context, key string, value any, ttl time.Duration) error {
					return nil
				},
			}
			testIGDBCacheWrapper, testErr := NewIGDBCacheWrapper(testCache, 10 * time.Minute, 100 * time.Millisecond, testLogger)
			assert.NoError(t, testErr)

			query := searchdef.SearchQuery{Query: "sample query", Limit: 10}
			expectedResult := &searchdef.SearchResult{
				Games: []searchdef.Game{{ID: 1, Name: "Dark Souls"}},
				Meta: searchdef.SearchMeta{
					CacheHit: false,
					CacheTTL: time.Duration(1 * time.Hour),
				},
			}

			testErr = testIGDBCacheWrapper.SetCachedResults(context.Background(), query, expectedResult)
			assert.NoError(t, testErr)
		},
	)

	t.Run(
		`SetCachedResults() experiences marshalling failure or Set failure`,
		func(t *testing.T) {
			/*
				GIVEN a cache client that returns an error on Set calls
				WHEN SetCachedResults() is called with a valid search result
				THEN an error is returned and the result is nil
			*/
			expectedError := errors.New("cannot connect to redis cache")
			testCache := &mocks.MockCacheClient{
				SetFunc: func(ctx context.Context, key string, value any, ttl time.Duration) error {
					return expectedError
				},
			}

			testLogger := &testutils.TestLogger{}
			testIGDBCacheWrapper, testErr := NewIGDBCacheWrapper(testCache, 10 * time.Minute, 100 * time.Millisecond, testLogger)
			assert.NoError(t, testErr)

			query := searchdef.SearchQuery{Query: "sample query", Limit: 10}
			expectedResult := &searchdef.SearchResult{
				Games: []searchdef.Game{{ID: 1, Name: "Dark Souls"}},
				Meta: searchdef.SearchMeta{
					CacheHit: false,
					CacheTTL: time.Duration(1 * time.Hour),
				},
			}

			testErr = testIGDBCacheWrapper.SetCachedResults(context.Background(), query, expectedResult)
			assert.Error(t, testErr)
			assert.Greater(t, len(testLogger.ErrorCalls), 0, "expected an error on SetCachedResults() failure")
		},
	)
}
