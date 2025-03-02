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

	"github.com/lokeam/qko-beta/internal/appcontext_test"
	"github.com/lokeam/qko-beta/internal/library"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/lokeam/qko-beta/internal/wishlist"
	authMiddleware "github.com/lokeam/qko-beta/server/middleware"
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
func mockSearchResultWithGames(games []types.Game) *searchdef.SearchResult {
	return &searchdef.SearchResult{
		Games: games,
	}
}

func TestSearchHandler(t *testing.T) {
	testIGDBToken := "valid-token"
	baseAppCtx := appcontext_test.NewTestingAppContext(testIGDBToken, nil)

	libraryService, err := library.NewGameLibraryService(baseAppCtx)
	if err != nil {
		t.Fatalf("failed to create library service: %v", err)
	}
	wishlistService, err := wishlist.NewGameWishlistService(baseAppCtx)
	if err != nil {
		t.Fatalf("failed to create wishlist service: %v", err)
	}

	createHandler := func(mockService *mockSearchService) http.HandlerFunc {
    // Create a domain services map instead of using the factory
    searchServices := make(DomainSearchServices)
    searchServices["games"] = mockService
    return NewSearchHandler(baseAppCtx, searchServices, libraryService, wishlistService)
	}
			/*
				GIVEN an HTTP request that doesn't contain a query parameter
				WHEN the SearchHandler is called
				THEN the error response is returned with httputils.RespondWithError containing the message "search query is required"
			*/
	// --------- Missing Query Parameter ---------
	t.Run("Missing Query Parameter", func(t *testing.T) {
		// Create mock search service
		mockSearchService := &mockSearchService{}

		// Create handler
		searchHandler := createHandler(mockSearchService)

		// Create test request with empty JSON body
		reqBody := strings.NewReader(`{}`)
		req := httptest.NewRequest(http.MethodPost, "/search", reqBody)
		req.Header.Set(httputils.XRequestIDHeader, "test-request-id")
		recorder := httptest.NewRecorder()

		// Call handler
		searchHandler.ServeHTTP(recorder, req)

		// Validate response
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing query, got %d", recorder.Code)
		}

		// Optional: Validate error message
		var errorResponse map[string]string
		if err := json.Unmarshal(recorder.Body.Bytes(), &errorResponse); err != nil {
			t.Fatalf("Failed to unmarshal error response: %v", err)
		}
		if errorResponse["error"] != "search query is required" {
			t.Errorf("Unexpected error message: %v", errorResponse["error"])
		}
	})

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
			// tokenError := errors.New("failed to retrieve token")


			mockSearchService := &mockSearchService{}
			mockSearchHandler := createHandler(mockSearchService)

			// Create test request via JSON body
			reqBody := strings.NewReader(`{"query": "dark souls"}`)
			testRequest := httptest.NewRequest(http.MethodPost, "/search", reqBody)
			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-2")
			testResponseRecorder := httptest.NewRecorder()

			// Add in the userID to the request context
			ctx := context.WithValue(testRequest.Context(), authMiddleware.UserIDKey, "test-user-id")
			testRequest = testRequest.WithContext(ctx)

			mockSearchHandler.ServeHTTP(testResponseRecorder, testRequest)
			if testResponseRecorder.Code != http.StatusInternalServerError {
				t.Errorf("expected status 500 for IGDB token retrieval failure, got: %d", testResponseRecorder.Code)
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
			games := []types.Game{{ ID: 1, Name: "Dark Souls"}}
			mockSearchService := &mockSearchService{
				searchServiceResult: mockSearchResultWithGames(games),
			}
			mockSearchHandler := createHandler(mockSearchService)

			// Create test request via JSON body
			reqBody := strings.NewReader(`{"query": "dark souls"}`)
			testRequest := httptest.NewRequest(http.MethodPost, "/search", reqBody)
			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-3")
			testResponseRecorder := httptest.NewRecorder()

			// Add in the userID to the request context
			ctx := context.WithValue(testRequest.Context(), authMiddleware.UserIDKey, "test-user-id")
			testRequest = testRequest.WithContext(ctx)

			mockSearchHandler.ServeHTTP(testResponseRecorder, testRequest)
			if testResponseRecorder.Code != http.StatusOK {
				t.Fatalf("expected the HTTP status code to be 200 but instead we got: %d", testResponseRecorder.Code)
			}

			// Verify the default limit is set to 50
			if mockSearchService.mostRecentRequest.Limit != 5 {
				t.Fatalf("expected the default limit to be 5 but instead we got: %d", mockSearchService.mostRecentRequest.Limit)
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
			mockSearchHandler := createHandler(mockSearchService)

			// NOTE: IGDB requires every search request to use POST instead of GET
			testRequest := httptest.NewRequest(http.MethodPost, "/search?query=darksouls&domain=music", nil)
			// testRequest := httptest.NewRequest(http.MethodGet, "/search?query=darksouls&domain=music", nil)

			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-4")
			testResponseRecorder := httptest.NewRecorder()

			mockSearchHandler.ServeHTTP(testResponseRecorder, testRequest)
			if testResponseRecorder.Code != http.StatusBadRequest {
				t.Errorf("expected status 400 for unsupported domain, got: %d", testResponseRecorder.Code)
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
			games := []types.Game{{ ID: 1, Name: "Dark Souls" }}
			mockSearchService := &mockSearchService{
				searchServiceResult: mockSearchResultWithGames(games),
			}
			mockSearchHandler := createHandler(mockSearchService)

			// Create test request via JSON body
			reqBody := strings.NewReader(`{"query": "dark souls", "limit": 100}`)
			testRequest := httptest.NewRequest(http.MethodPost, "/search", reqBody)
			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-5")
			testResponseRecorder := httptest.NewRecorder()

			// Add in the userID to the request context
			ctx := context.WithValue(testRequest.Context(), authMiddleware.UserIDKey, "test-user-id")
			testRequest = testRequest.WithContext(ctx)

			mockSearchHandler.ServeHTTP(testResponseRecorder, testRequest)
			if testResponseRecorder.Code != http.StatusOK {
				t.Fatalf("expected the HTTP status code to be 200 but instead we got: %d", testResponseRecorder.Code)
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
			mockSearchHandler := createHandler(mockSearchService)

			// Create test request via JSON body
			reqBody := strings.NewReader(`{"query": "dark souls"}`)
			testRequest := httptest.NewRequest(http.MethodPost, "/search", reqBody)
			testRequest.Header.Set(httputils.XRequestIDHeader, "test-request-id-6")
			testResponseRecorder := httptest.NewRecorder()

			// Add in the userID to the request context
			ctx := context.WithValue(testRequest.Context(), authMiddleware.UserIDKey, "test-user-id")
			testRequest = testRequest.WithContext(ctx)

			mockSearchHandler.ServeHTTP(testResponseRecorder, testRequest)
			if testResponseRecorder.Code != http.StatusInternalServerError {
				t.Errorf("expected status 500 for search service failure, got: %d", testResponseRecorder.Code)
			}
		},
	)
}
