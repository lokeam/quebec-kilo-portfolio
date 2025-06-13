package search

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lokeam/qko-beta/internal/appcontext_test"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/shared/constants"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
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

// mockSearchService implements the SearchService interface for testing
type mockSearchService struct {
	searchServiceResult   *searchdef.SearchResult
	searchServiceError    error
	mostRecentRequest     *searchdef.SearchRequest
	storageLocations      types.AddGameFormStorageLocationsResponse
	storageLocationsError error
}

func (mss *mockSearchService) Search(
	ctx context.Context,
	req searchdef.SearchRequest,
) (*searchdef.SearchResult, error) {
	mss.mostRecentRequest = &req
	return mss.searchServiceResult, mss.searchServiceError
}

func (mss *mockSearchService) GetAllGameStorageLocationsBFF(
	ctx context.Context,
	userID string,
) (types.AddGameFormStorageLocationsResponse, error) {
	return mss.storageLocations, mss.storageLocationsError
}

// mockLibraryService implements the LibraryService interface for testing
type mockLibraryService struct {
	isInLibrary bool
	err         error
}

func (mls *mockLibraryService) CreateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error {
	return nil
}

func (mls *mockLibraryService) GetAllLibraryItemsBFF(ctx context.Context, userID string) (types.LibraryBFFResponseFINAL, error) {
	return types.LibraryBFFResponseFINAL{}, nil
}

func (mls *mockLibraryService) GetSingleLibraryGame(ctx context.Context, userID string, gameID int64) (types.LibraryGameItemBFFResponseFINAL, error) {
	return types.LibraryGameItemBFFResponseFINAL{}, nil
}

func (mls *mockLibraryService) UpdateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error {
	return nil
}

func (mls *mockLibraryService) DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error {
	return nil
}

func (mls *mockLibraryService) InvalidateUserCache(ctx context.Context, userID string) error {
	return nil
}

func (mls *mockLibraryService) IsGameInLibraryBFF(ctx context.Context, userID string, gameID int64) (bool, error) {
	return mls.isInLibrary, mls.err
}

// mockWishlistService implements the WishlistService interface for testing
type mockWishlistService struct {
	wishlistItems []models.GameToSave
	err           error
}

func (mws *mockWishlistService) GetWishlistItems(ctx context.Context, userID string) ([]models.GameToSave, error) {
	return mws.wishlistItems, mws.err
}

// Helper function to create a SearchResult
func mockSearchResultWithGames(games []models.Game) *searchdef.SearchResult {
	return &searchdef.SearchResult{
		Games: games,
	}
}

// Helper function to create a request with userID in context
func createRequestWithUserID(method, path, body string) (*http.Request, *httptest.ResponseRecorder) {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set(httputils.XRequestIDHeader, "test-request-id")

	// Add userID to context
	ctx := context.WithValue(req.Context(), constants.UserIDKey, "test-user-id")
	req = req.WithContext(ctx)

	return req, httptest.NewRecorder()
}

func TestSearchHandler(t *testing.T) {
	testIGDBToken := "valid-token"
	baseAppCtx := appcontext_test.NewTestingAppContext(testIGDBToken, nil)

	// Create mock services
	mockSearchService := &mockSearchService{}
	mockLibraryService := &mockLibraryService{}
	mockWishlistService := &mockWishlistService{}

	// Create handler with mock services
	handler := NewSearchHandler(
		baseAppCtx,
		mockSearchService,
		mockLibraryService,
		mockWishlistService,
	)

	t.Run("Missing Query Parameter", func(t *testing.T) {
		req, recorder := createRequestWithUserID(http.MethodPost, "/search", `{}`)

		handler.Search(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing query, got %d", recorder.Code)
		}

		var errorResponse map[string]string
		if err := json.Unmarshal(recorder.Body.Bytes(), &errorResponse); err != nil {
			t.Fatalf("Failed to unmarshal error response: %v", err)
		}
		if errorResponse["error"] != "search query is required" {
			t.Errorf("Unexpected error message: %v", errorResponse["error"])
		}
	})

	t.Run("IGDB Token Retrieval Failure", func(t *testing.T) {
		// Create app context with invalid token
		invalidAppCtx := appcontext_test.NewTestingAppContext("", errors.New("token error"))
		handler := NewSearchHandler(
			invalidAppCtx,
			mockSearchService,
			mockLibraryService,
			mockWishlistService,
		)

		req, recorder := createRequestWithUserID(http.MethodPost, "/search", `{"query": "dark souls"}`)

		handler.Search(recorder, req)

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500 for IGDB token retrieval failure, got %d", recorder.Code)
		}
	})

	t.Run("Happy Path - Search Request", func(t *testing.T) {
		// Setup mock search result
		games := []models.Game{{ID: 1, Name: "Dark Souls"}}
		mockSearchService.searchServiceResult = mockSearchResultWithGames(games)
		mockSearchService.searchServiceError = nil

		// Setup mock library service
		mockLibraryService.isInLibrary = true
		mockLibraryService.err = nil

		// Setup mock wishlist service
		mockWishlistService.wishlistItems = []models.GameToSave{{GameID: 1}}
		mockWishlistService.err = nil

		req, recorder := createRequestWithUserID(http.MethodPost, "/search", `{"query": "dark souls"}`)

		handler.Search(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", recorder.Code)
		}

		// Verify the request was made with correct parameters
		if mockSearchService.mostRecentRequest.Query != "dark souls" {
			t.Errorf("Expected query 'dark souls', got %s", mockSearchService.mostRecentRequest.Query)
		}

		// Verify the response
		var response searchdef.SearchResponse
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response.Games) != 1 {
			t.Errorf("Expected 1 game in response, got %d", len(response.Games))
		}

		if !response.Games[0].IsInLibrary {
			t.Error("Expected game to be in library")
		}

		if !response.Games[0].IsInWishlist {
			t.Error("Expected game to be in wishlist")
		}
	})

	t.Run("GetGameStorageLocationsBFF", func(t *testing.T) {
		// Setup mock storage locations
		locations := types.AddGameFormStorageLocationsResponse{
			PhysicalLocations: []types.AddGameFormPhysicalLocationsResponse{
				{
					ParentLocationID:   "1",
					ParentLocationName: "Shelf 1",
					SublocationID:      "1",
					SublocationName:    "Shelf 1",
				},
				{
					ParentLocationID:   "2",
					ParentLocationName: "Shelf 2",
					SublocationID:      "2",
					SublocationName:    "Shelf 2",
				},
			},
		}
		mockSearchService.storageLocations = locations
		mockSearchService.storageLocationsError = nil

		req, recorder := createRequestWithUserID(http.MethodGet, "/search/bff", "")

		handler.GetGameStorageLocationsBFF(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", recorder.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		storageLocations, ok := response["storage_locations"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected storage_locations in response")
		}

		physicalLocations, ok := storageLocations["physical_locations"].([]interface{})
		if !ok {
			t.Fatal("Expected physical_locations in response")
		}

		if len(physicalLocations) != 2 {
			t.Errorf("Expected 2 storage locations, got %d", len(physicalLocations))
		}
	})

	t.Run("GetGameStorageLocationsBFF Error", func(t *testing.T) {
		// Setup mock error
		mockSearchService.storageLocationsError = errors.New("storage error")

		req, recorder := createRequestWithUserID(http.MethodGet, "/search/bff", "")

		handler.GetGameStorageLocationsBFF(recorder, req)

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", recorder.Code)
		}
	})
}
