package search

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/appcontext_test"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

/*
	Behaviour:
	- Reading the query and parameters from request (example: domain, limit)
	- Retrieving the IGDB token from cache for use in the search (POST) request
	- Send the search request to approriate domains-specific search services
	- Respond with JSON on success or error on failure (both flows)

	Scenarios:
	- Missing query parameter
	- Missing domain parameter
	- IGDB token failure
	- Valid search request with default (game) domain
	- Valid search request with unsupported domain
	- Valid search request with a custom limit parameter
	- Valid search request with a search service failure
*/

// --------- Mock Search Service ---------
type mockSearchService struct {
	searchServiceResult   *searchdef.SearchResult
	searchServiceError    error
	mostRecentRequest     *searchdef.SearchRequest
}

func (mss *mockSearchService) Search(
	ctx context.Context,
	req searchdef.SearchRequest,
) (*searchdef.SearchResult, error) {
	mss.mostRecentRequest = &req

	return mss.searchServiceResult, mss.searchServiceError
}

// Helper fn to create a SearchResult
func mockSearchResultWithGames(games []searchdef.Game) *searchdef.SearchResult {
	return &searchdef.SearchResult{
		Games: games,
	}
}

func TestSearchHandler(t *testing.T) {
	testIGDBToken := "valid-token"
	baseAppCtx := appcontext_test.NewTestingAppContext(testIGDBToken, nil)

	// --------- Missing Query Parameter ---------
	t.Run(
		`Missing Query Parameter`,
		func(t *testing.T) {
			/*
				GIVEN an HTTP request that doesn't contain a query parameter
				WHEN the SearchHandler is called
				THEN the error response is returned with httputils.RespondWithError containing the message "search query is required"
			*/
			mockSearchService := &mockSearchService{}
			mockSearchHandler := &SearchHandler{
				appContext:            baseAppCtx,
				gameSearchService:     mockSearchService,
			}

			// NOTE: IGDB requires every search request to use POST instead of GET
			testRequest := httptest.NewRequest(http.MethodPost, "/search", nil)
			// testRequest := httptest.NewRequest(http.MethodGet, "/search", nil)

			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-1")
			testResponseRecorder := httptest.NewRecorder()

			testError := mockSearchHandler.HandleSearch(testResponseRecorder, testRequest)
			if testError == nil || !strings.Contains(testError.Error(), "search query is required") {
				t.Errorf("expected an error for a missing query parameter, but got: %v", testError)
			}
		},
	)

	// --------- No token or error from IGDB.GetAccessTokenKey ---------
	t.Run(
		`IGDB Token Retrieval Failure`,
		func(t *testing.T) {
			/*
				GIVEN an HTTP request with a valid query parameter + any necessary headers
				AND an app context whose IGDB configuration is set to return an error
				WHEN we call HandleSearch()
				THEN httputils.RespondWithError produces an error response indicating that the IGDB token could not be retrieved
			*/
			tokenError := errors.New("failed to retrieve token")
			brokenAppCtx := &appcontext.AppContext{
				Logger: baseAppCtx.Logger,
				Config: &config.Config{
					IGDB: &config.IGDBConfig{
						ClientID:       "dummyID",
						ClientSecret:   "dummySecret",
						AuthURL:        "dummyAuthURL",
						BaseURL:        "dummyBaseURL",
						TokenTTL:       24 * time.Hour,
						// Simulate an error by leaving AccessTokenKey empty:
						AccessTokenKey: "",
					},
				},
			}
			mockSearchService := &mockSearchService{}
			mockSearchHandler := &SearchHandler{
				appContext:         brokenAppCtx,
				gameSearchService:  mockSearchService,
			}

			// NOTE: IGDB requires every search request to use POST instead of GET
			testRequest := httptest.NewRequest(http.MethodPost, "/search?query=darksouls", nil)
			// testRequest := httptest.NewRequest(http.MethodGet, "/search?query=darksouls", nil)

			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-2")
			testResponseRecorder := httptest.NewRecorder()

			testError := mockSearchHandler.HandleSearch(testResponseRecorder, testRequest)
			if testError == nil || !strings.Contains(testError.Error(), tokenError.Error()) {
				t.Fatalf("expected an error for a missing IGDB token, but got: %v", testError)
			}
		},
	)

	// --------- Valid Search Request with Default Domain ---------
	t.Run(
		`Happy Path - Search Request + Default Domain all ok`,
		func(t *testing.T) {
			/*
				GIVEN an HTTP request with a valid query parameter, no domain and no limit parameters
				AND an app context with a valid IGDB token
				AND a mock GameSearchService set up to return a valid SearchResult (such as a list of games)
				WHEN the SearchHandler() is called
				THEN the handler defaults the domain to "games" and limit to 50
				AND calls the GameSearchService.Search() method containing the query, domain and default limit
				AND returns a JSON response containing search results
			*/

			// Simulate valid search result with an example game
			games := []searchdef.Game{{ ID: 1, Name: "Dark Souls"}}
			mockSearchService := &mockSearchService{
				searchServiceResult: mockSearchResultWithGames(games),
			}
			mockSearchHandler := &SearchHandler{
				appContext:            baseAppCtx,
				gameSearchService:     mockSearchService,
			}

			// NOTE: IGDB requires every search request to use POST instead of GET
			testRequest := httptest.NewRequest(http.MethodPost, "/search?query=darksouls", nil)
			// testRequest := httptest.NewRequest(http.MethodGet, "/search?query=darksouls", nil)

			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-3")
			testResponseRecorder := httptest.NewRecorder()

			testError := mockSearchHandler.HandleSearch(testResponseRecorder, testRequest)
			if testError != nil {
				t.Fatalf("ran HandleSearch() and didn't expect to get an error, but instead we got: %v", testError)
			}

			// Verify the default limit is set to 50
			if mockSearchService.mostRecentRequest.Limit != 50 {
				t.Fatalf("expected the default limit to be 50 but instead we got: %d", mockSearchService.mostRecentRequest.Limit)
			}

			// Check the HTTP response
			if testResponseRecorder.Code != http.StatusOK {
				t.Fatalf("expected the HTTP status code to be 200 but instead we got: %d", testResponseRecorder.Code)
			}

			// Convert the raw response body data into something we can work with
			var testSearchResponse searchdef.SearchResponse
			rawRRBodyData, _ := io.ReadAll(testResponseRecorder.Body)
			if testError := json.Unmarshal(rawRRBodyData, &testSearchResponse); testError != nil {
				t.Fatalf("failed to unmarshal the response body data: %v", testError)
			}

			// Check the search result length
			if testSearchResponse.Total != len(games) {
				t.Fatalf("expected the total number of search results to be: %d but instead we got: %d", len(games), testSearchResponse.Total)
			}
		},
	)

	// --------- Valid Search Request with Unsupported Domain ---------
	t.Run(
		`Valid Search Request but Unsupported Domain error`,
		func(t *testing.T) {
			/*
				GIVEN an HTTP request with a valid query paramter and domain parameter that is not supported (example: "music")
				WHEN the SearchHandler() is called
				THEN httputils.RespondWithError produces an error with message "unsupported search domain"
			*/
			mockSearchService := &mockSearchService{}
			mockSearchHandler := &SearchHandler{
				appContext:            baseAppCtx,
				gameSearchService:     mockSearchService,
			}

			// NOTE: IGDB requires every search request to use POST instead of GET
			testRequest := httptest.NewRequest(http.MethodPost, "/search?query=darksouls&domain=music", nil)
			// testRequest := httptest.NewRequest(http.MethodGet, "/search?query=darksouls&domain=music", nil)

			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-4")
			testResponseRecorder := httptest.NewRecorder()

			testError := mockSearchHandler.HandleSearch(testResponseRecorder, testRequest)
			if testError == nil || !strings.Contains(testError.Error(), "unsupported search domain") {
				t.Fatalf("expected an error for an unsupported domain, but instead we got: %v", testError)
			}
		},
	)

	// --------- Valid Search Request with Custom Limit Parameter ---------
	t.Run(
		`Happy Path - Search Request + Custom Limit Parameter`,
		func(t *testing.T) {
			/*
				GIVEN an HTTP request with a valid query parameter, and a limit parameter that is a valid integer (example: limit=100)
				WHEN the SearchHandler() is called
				THEN the limit used in the constructed Search request is the the provided integer instead of the default (example: 50)
				AND the JSON success response reflects this limit with the correct number of search results
			*/
			games := []searchdef.Game{{ ID: 1, Name: "Dark Souls" }}
			mockSearchService := &mockSearchService{
				searchServiceResult: mockSearchResultWithGames(games),
			}
			mockSearchHandler := &SearchHandler{
				appContext:            baseAppCtx,
				gameSearchService:     mockSearchService,
			}

			// NOTE: IGDB requires every search request to use POST instead of GET
			testRequest := httptest.NewRequest(http.MethodPost, "/search?query=darksouls&limit=100", nil)
			// testRequest := httptest.NewRequest(http.MethodGet, "/search?query=darksouls&limit=100", nil)

			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-5")
			testResponseRecorder := httptest.NewRecorder()

			testError := mockSearchHandler.HandleSearch(testResponseRecorder, testRequest)
			if testError != nil {
				t.Fatalf("ran HandleSearch() and didn't expect to get an error, but instead we got: %v", testError)
			}
			if mockSearchService.mostRecentRequest == nil {
				t.Fatalf("expected the most recent request to be set, but it was not")
			}
			if mockSearchService.mostRecentRequest.Limit != 100 {
				t.Fatalf("expected the limit to be 100 but instead we got: %d", mockSearchService.mostRecentRequest.Limit)
			}
			if testResponseRecorder.Code != http.StatusOK {
				t.Fatalf("expected the HTTP status code to be 200 but instead we got: %d", testResponseRecorder.Code)
			}
		},
	)

	// --------- Valid Search Request with Search Service Failure ---------
	t.Run(
		`Valid Search Request but Search Service Fails`,
		func(t *testing.T) {
			/*
				GIVEN an HTTP request with a valid query parameters,
				AND a mock SearchService set up to return a non-nil error (simulating failure in the search service logic)
				WHEN the SearchHandler() is called
				THEN the error is passed to httputils.RespondWithError
				AND the resulting response contains the error message returned by the mock SearchService
			*/
			mockSearchServiceError := errors.New("search service failure")
			mockSearchService := &mockSearchService{
				searchServiceError: mockSearchServiceError,
			}
			mockSearchHandler := &SearchHandler{
				appContext:            baseAppCtx,
				gameSearchService:     mockSearchService,
			}

			// NOTE: IGDB requires every search request to use POST instead of GET
			testRequest := httptest.NewRequest(http.MethodPost, "/search?query=darksouls", nil)
			// testRequest := httptest.NewRequest(http.MethodGet, "/search?query=darksouls", nil)

			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-6")
			testResponseRecorder := httptest.NewRecorder()

			testError := mockSearchHandler.HandleSearch(testResponseRecorder, testRequest)
			if testError == nil || !strings.Contains(testError.Error(), mockSearchServiceError.Error()) {
				t.Fatalf("expected a search service error, but instead we got: %v", testError)
			}
		},
	)
}
