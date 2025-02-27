package search

import (
	"context"
	"errors"
	"testing"

	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCacheWrapper struct {
	mock.Mock
}


// Mock implementations of CacheWrapper methods
func (m *MockCacheWrapper) GetCachedResults(ctx context.Context, key string, result any) (bool, error) {
	args := m.Called(ctx, key, result)

    // If there's a result setup via the Run function, it will already be populated
    // in the result parameter (a pointer) by a test Run() function

    return args.Bool(0), args.Error(1)
}

func (m *MockCacheWrapper) SetCachedResults(ctx context.Context, key string, result any) error {
	args := m.Called(ctx, key, result)
	return args.Error(0)
}



func TestIGDBCacheAdapter(t *testing.T) {

	// ------ GetCachedResults() ------
	t.Run(`GetCachedResults() handles a cache miss`, func(t *testing.T) {
		/*
			GIVEN a properly configured IGDBCacheAdapter
			AND a MockCacheWrapper set up to return a cache miss
			AND a properly configured SearchQuery
			WHEN GetCachedResults() is called
			THEN GetCachedResults() returns a cache miss and an error
		*/
		mockCacheWrapper := &MockCacheWrapper{}
		testAdapter, err := NewIGDBCacheAdapter(mockCacheWrapper)
		assert.NoError(t, err)

		testQuery := searchdef.SearchQuery{Query: "Dark Souls", Limit: 5}

		// Setup mock to return cache miss
		mockCacheWrapper.On("GetCachedResults",
    mock.Anything,
    mock.AnythingOfType("string"),
    mock.AnythingOfType("*searchdef.SearchResult")).
    Return(false, nil, false)

			// Call method under test
			result, err := testAdapter.GetCachedResults(context.Background(), testQuery)

			assert.Nil(t, err, "Should not return on cache miss")
			assert.Nil(t, result, "Result should be nil on cache miss")
			mockCacheWrapper.AssertExpectations(t)
		})


	// ------ SetCachedResults() ------
	t.Run(`SetCachedResults() properly propagates errors from the underlying cache wrapper`, func(t *testing.T) {
		/*
			GIVEN a propelry configured IGDBAdapter
			AND a MockCacheWrapper set up to return an error
			AND a properly configured SearchQUery and SearchResult
			WHEN SetCachedResults() is called
			THEN the error from the cache wrapper is propagated
		*/
		mockCache := &MockCacheWrapper{}
		testAdapter, err := NewIGDBCacheAdapter(mockCache)
		assert.NoError(t, err)

		testQuery := searchdef.SearchQuery{Query: "Dark Souls", Limit: 5}
		testResult := &searchdef.SearchResult{
			Games: []types.Game{{ID: 1, Name: "Dark Souls"}},
		}

		expectedError := errors.New("cache error")
		mockCache.On("SetCachedResults",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.Anything).
			Return(expectedError)

			testAdapterErr := testAdapter.SetCachedResults(context.Background(), testQuery, testResult)

			assert.Error(t, testAdapterErr)
			assert.Equal(t, expectedError, testAdapterErr, "Should propagate the original error")
			mockCache.AssertExpectations(t)
	})

	t.Run(`SetCachedResults() delegates to the cache wrapper`, func(t *testing.T) {
		/*
			GIVEN a properly configured IGDBCacheAdapter
			AND a properly configured MockCacheWrapper
			AND a properly configured SearchQuery
			WHEN SetCachedResults() is called
			THEN SetCachedResults passes the call to the IGDB Cache Adapter
		*/
		mockCache := &MockCacheWrapper{}
		testAdapter, err := NewIGDBCacheAdapter(mockCache)
		assert.NoError(t, err)

		testQuery := searchdef.SearchQuery{Query: "Dark Souls", Limit: 5}
		testResult := &searchdef.SearchResult{
			Games: []types.Game{{ID: 1, Name: "Dark Souls"}},
		}

		mockCache.On("SetCachedResults",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.Anything).
			Return(nil)

		testAdapterErr := testAdapter.SetCachedResults(context.Background(), testQuery, testResult)
		assert.NoError(t, testAdapterErr)
		mockCache.AssertExpectations(t)
		mockCache.AssertNumberOfCalls(t, "SetCachedResults", 1)
	})

	t.Run(`Happy Path: GetCachedResults() returns a cache hit and gets data from wrapper`, func(t *testing.T) {
		/*
			GIVEN a properly configured IGDBCacheAdapter
			AND a properly configured MockCacheWrapper
			AND a properly configured SearchQuery
			WHEN GetCachedResults() is called
			THEN GetCachedResults returns a cache hit and the data is returned
		*/
		mockCache := &MockCacheWrapper{}
		testAdapter, err := NewIGDBCacheAdapter(mockCache)
		assert.NoError(t, err)

		testQuery := searchdef.SearchQuery{Query: "Dark Souls", Limit: 5}
		expectedResult := &searchdef.SearchResult{
			Games: []types.Game{{ID: 1, Name: "Dark Souls"}},
			Meta: searchdef.SearchMeta{
        CacheHit: true, // NOTE: this needs to be set in order to pass the full equality check
    	},
		}

		mockCache.On("GetCachedResults",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*searchdef.SearchResult")).
			Run(func(args mock.Arguments) {
				// Populate the result argument
				result := args.Get(2).(*searchdef.SearchResult)
				result.Games = expectedResult.Games
			}).
			Return(true, nil, true)

		result, testAdapterErr := testAdapter.GetCachedResults(context.Background(), testQuery)
		assert.NoError(t, testAdapterErr)
		assert.Equal(t, expectedResult, result)
		mockCache.AssertExpectations(t)
	})

	// ------ Edge cases ------
	t.Run(`Edge case test - IGDBCacheAdapter generates consistent cache keys for identical search queries`, func(t *testing.T) {
		/*
			GIVEN a properly configured IGDBCacheAdapter
			AND two identical SearchQueries
			WHEN generating cache keys for both
			THEN the generated keys should be the same
		*/
		mockCache := &MockCacheWrapper{}
		testAdapter, err := NewIGDBCacheAdapter(mockCache)
		assert.NoError(t, err)

		testQuery1 := searchdef.SearchQuery{Query: "Dark Souls", Limit: 5}
		testQuery2 := searchdef.SearchQuery{Query: "Dark Souls", Limit: 5}

		// Grab the cache keys by intercepting the calls to the mock
		var cacheKey1, cacheKey2 string
		mockCache.On("GetCachedResults",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*searchdef.SearchResult")).
			Run(func(args mock.Arguments) {
				if cacheKey1 == "" {
					cacheKey1 = args.String(1)
				} else {
					cacheKey2 = args.String(1)
				}
			}).
			Return(true, nil, true)

		testAdapter.GetCachedResults(context.Background(), testQuery1)
		testAdapter.GetCachedResults(context.Background(), testQuery2)

		assert.NotEmpty(t, cacheKey1, "First cache key should not be empty")
		assert.NotEmpty(t, cacheKey2, "Second cache key should not be empty")
		assert.Equal(t, cacheKey1, cacheKey2, "Cache keys should be identical")
	})

	t.Run(`Edge case test - IGDBCacheAdapter propagates context cancellation to cache wrapper`, func(t *testing.T) {
		/*
			GIVEN a properly configured IGDBCacheAdapter
			AND a cancelled context
			WHEN GetCachedResults() is called
			THEN the context cancellation should be propagated to the cache wrapper
		*/
		mockCache := &MockCacheWrapper{}
		testAdapter, err := NewIGDBCacheAdapter(mockCache)
		assert.NoError(t, err)

		testQuery := searchdef.SearchQuery{Query: "Dark Souls", Limit: 5}
		canceledContext, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Set up mock to verify the cancelled context is passed
		expectedError := context.Canceled
		mockCache.On("GetCachedResults",
			mock.MatchedBy(func(ctx context.Context) bool {
				return ctx.Err() != nil
			}),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*searchdef.SearchResult")).
			Return(false, expectedError, false)

		result, err := testAdapter.GetCachedResults(canceledContext, testQuery)
		assert.Error(t, err, "Should return error for cancelled context")
		assert.Equal(t, expectedError, err, "Should propagate context cancellation error")
		assert.Nil(t, result, "Result should be nil for cancelled context")
		mockCache.AssertExpectations(t)
	})
}
