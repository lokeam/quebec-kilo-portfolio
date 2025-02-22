package search

import (
	"context"
	"errors"
	"testing"

	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	"github.com/lokeam/qko-beta/internal/types"
)

/*
	Behavior:
		Search()
			- Game search service calls sanitizer to sanitize the query
			- Game search service calls validator to validate the query
			- Game search service builds the query from the request
			- Game search service attempts to retrieve cached results
				- Game search service returns cached results if
				- Game search service calls the igdb adapter to search for games on cache miss
				- Game search service converts the igdb games result to a searchdef.SearchResult
				- Game search service caches this recently retrieved and converted result

	Scenarios:
		- Sanitation Failure:
			* Search request fails sanitization, logs sanitizer error, returns and logs error
		- Validation Failure:
			* Search request passes sanitization but fails sanitation, returns and logs error
		- Cache Hit:
			* Search request passes both sanitization and validation AND cache exists, returns cached result
		- Cache Miss:
			* Search request passes both sanitization and validation BUT no cached result exists
			  service then calls IGDB adapter to fetch data, convert results and cache data and return result
		- Adapter Failure (Search Error):
			* Search request passes both sanitization and validation AND cache exists
			  IGDB adapter returns an error
			  Game search service logs the error and returns the error
		- Caching Fresh Data Failure Does NOT Block Result Return:
			* Search request passes both sanitization and validation AND cache exists
			  Caching wrapper returns an error (Redis failure)
			  Game search service logs the error BUT still returns the result
*/

/* ------ Helper to create a GameSearchService with Mocks ------ */

func newMockGameSearchServiceWithDefaults(logger *testutils.TestLogger) *GameSearchService {
	mockConfig := mocks.NewMockConfig()

	return &GameSearchService{
		adapter:        mocks.DefaultIGDBAdapter(),
		config:         mockConfig,
		logger:         logger,
		validator:      mocks.DefaultValidator(),
		sanitizer:      mocks.DefaultSanitizer(),
		cacheWrapper:   mocks.DefaultCacheWrapper(),
	}
}

func TestGameSearchService(t *testing.T) {
	ctx := context.Background()

	// TBD
	// --------- Sanitization Failure ---------
	t.Run(
		`Search service fails sanitization`,
		func(t *testing.T) {
			/*
				GIVEN a search request with a query that causes the sanitizer to fail
				WHEN the Search() method is called
				THEN the service should log the sanitizer error and return an error WITHOUT proceeding to validation or cache lookup
			*/
			testLogger := testutils.NewTestLogger()
			testSearchService := newMockGameSearchServiceWithDefaults(testLogger)

			// Override the sanitizer to simulate a failure
			testSearchService.sanitizer = &mocks.MockSanitizer{
				SanitizeFunc: func(query string) (string, error) {
					return "", errors.New("sanitization failure")
				},
			}

			_, testErr := testSearchService.Search(ctx, searchdef.SearchRequest{
				Query: "<script>alert('trigger xss sanitizer error');</script>",
			})
			if testErr == nil || testErr.Error() != "sanitizer failure" {
				t.Errorf("expected sanitization failure error, but instead got: %v", testErr)
			}
		},
	)

	// --------- Validation Fa1ilure ---------
	t.Run(
		`Search service fails validation`,
		func(t *testing.T) {
			/*
				GIVEN a search request where the query is successfully sanitized but fails validation
				WHEN the Search() method is called
				THEN the service should log the validation error and return the validation error without attempting to cehck the cache or search via the adapter
			*/
			testLogger := testutils.NewTestLogger()
			testSearchService := newMockGameSearchServiceWithDefaults(testLogger)

			testSearchService.validator = &mocks.MockValidator{
				ValidateFunc: func(query searchdef.SearchQuery) error {
					return errors.New("validation failure")
				},
			}

			_, testErr := testSearchService.Search(ctx, searchdef.SearchRequest{
				Query: "some query",
			})
			if testErr == nil || testErr.Error() != "validation failure" {
				t.Errorf("expected validation failure error, but instead got: %v", testErr)
			}
		},
	)

	// --------- Cache Hit ---------
	t.Run(
		`Search service returns cached result on cache hit`,
		func(t *testing.T) {
			/*
				GIVEN a fully valid search request (successful sanitization and validation) and a cache entry exists
				WHEN the Search() method is called
				THEN the service should return the cached result
			*/
			testLogger := testutils.NewTestLogger()
			testSearchService := newMockGameSearchServiceWithDefaults(testLogger)

			expectedResult := searchdef.SearchResult{
				Games: []searchdef.Game{
					{
						ID:   1,
						Name: "Dark Souls",
					},
				},
			}

			testSearchService.cacheWrapper = &mocks.MockCacheWrapper{
				GetCachedResultsFunc: func(ctx context.Context, sq searchdef.SearchQuery) (*searchdef.SearchResult, error) {
					return &expectedResult, nil
				},
				SetCachedResultsFunc: mocks.DefaultCacheWrapper().SetCachedResultsFunc,
				TimeToLive:           60,
			}

			actualResult, testErr := testSearchService.Search(ctx, searchdef.SearchRequest{
				Query: "some query",
			})
			if testErr != nil {
				t.Errorf("expected no error on cache hit, but instead got: %v", testErr)
			}
			if !searchResultEqual(*actualResult, expectedResult) {
				t.Errorf("expected result to be %v, but instead got: %v", expectedResult, *actualResult)
			}
		},
	)

	// --------- Cache Miss with Successful Adapter Search ---------
	t.Run(
		`Search service performs adapter search on cache miss`,
		func(t *testing.T) {
			/*
				GIVEN a fully valid search request (successful sanitization and validation) and no cached result exists
				WHEN the Search() method is called
				THEN the service should call the IGDB adapter to fetch data, cache the fresh data and then return a valid Seawrch result
			*/
			testLogger := testutils.NewTestLogger()
			testSearchService := newMockGameSearchServiceWithDefaults(testLogger)

			// Override the adapter to simulate a successful search
			fetchedGameResponse := []*types.Game{
				{
					ID:   1,
					Name: "Dark Souls",
				},
				{
					ID:   2,
					Name: "Dark Souls 2",
				},
				{
					ID:   3,
					Name: "Dark Souls 3",
				},
			}

			testSearchService.adapter = &mocks.MockIGDBAdapter{
				SearchGamesFunc: func(ctx context.Context, query string, limit int) ([]*types.Game, error) {
					return fetchedGameResponse, nil
				},
			}

			searchServiceResult, err := testSearchService.Search(ctx, searchdef.SearchRequest{
				Query: "Dark Souls",
			})
			if err != nil {
				t.Errorf("expected no error on successful adapter search, got: %v", err)
			}
			if searchServiceResult == nil {
				t.Errorf("expected valid search result on cache miss, got nil")
			}
		},
	)

	// --------- Adapter Failure (Search Error) ---------
	t.Run(
		`Search service fails on adapter error`,
		func(t *testing.T) {
			/*
				GIVEN a fully valid search request with a cache miss
				WHEN the IGDB adapter fails (example: returning an error during the search call to IGDB)
				THEN the service should log the adapter error and return that error without further processing
			*/
			testLogger := testutils.NewTestLogger()
			testSearchService := newMockGameSearchServiceWithDefaults(testLogger)

			testSearchService.adapter = &mocks.MockIGDBAdapter{
				SearchGamesFunc: func(ctx context.Context, query string, limit int) ([]*types.Game, error) {
					return nil, errors.New("adapter failure")
				},
			}

			_, testErr := testSearchService.Search(ctx, searchdef.SearchRequest{
				Query: "Dark Souls",
			})
			if testErr == nil || testErr.Error() != "adapter failure" {
				t.Errorf("expected adapter failure error, but instead got: %v", testErr)
			}
		},
	)

	// --------- Caching Fresh Data Failure Does NOT Block Result ---------
	t.Run(
		`Search service fails to cache fresh data but still returns result`,
		func(t *testing.T) {
			/*
				GIVEN a valid search request with a cache miss and successful adapter search
				WHEN the service attempts to cache the fresh data but encounters and error (example: Redis fails)
				THEN the service should log the cache error but still return the valid SearchResult (since failing to cache should not block a successful search response)
			*/
			testLogger := testutils.NewTestLogger()
			testSearchService := newMockGameSearchServiceWithDefaults(testLogger)

			fetchedGameResponse := []*types.Game{
				{
					ID:   1,
					Name: "Dark Souls",
				},
			}

			testSearchService.adapter = &mocks.MockIGDBAdapter{
				SearchGamesFunc: func(ctx context.Context, query string, limit int) ([]*types.Game, error) {
					return fetchedGameResponse, nil
				},
			}

			testSearchService.cacheWrapper = &mocks.MockCacheWrapper{
				SetCachedResultsFunc: func(ctx context.Context, sq searchdef.SearchQuery, sr *searchdef.SearchResult) error {
					return errors.New("cache failure")
				},
				TimeToLive: 60,
			}

			searchServiceResult, err := testSearchService.Search(ctx, searchdef.SearchRequest{
				Query: "Dark Souls",
			})
			if err != nil {
				t.Errorf("expected no error on cache miss, but instead got: %v", err)
			}
			if searchServiceResult == nil {
				t.Errorf("expected valid search result on cache miss, got nil")
			}
		},
	)
}


// searchResultEqual compares two SearchResult values.
func searchResultEqual(a, b searchdef.SearchResult) bool {
	// Compare number of games.
	if len(a.Games) != len(b.Games) {
			return false
	}
	// Compare each game (adjust the fields as needed).
	for i := range a.Games {
			if a.Games[i].ID != b.Games[i].ID || a.Games[i].Name != b.Games[i].Name {
					return false
			}
	}
	// Compare meta fields.
	if a.Meta.CacheHit != b.Meta.CacheHit ||
	a.Meta.CacheTTL != b.Meta.CacheTTL ||
	a.Meta.TimestampUTC != b.Meta.TimestampUTC {
			return false
	}
	return true
}