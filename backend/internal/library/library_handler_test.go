package library

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/shared/constants"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/types"
)

/*
	Behaviour:
	- Reading the request body for POST requests
	- Reading URL parameters for DELETE requests
	- Retrieving the user ID from the request context
	- Dispatching to the appropriate library service based on domain
	- Responding with JSON on success or error on failure

	Scenarios:
	- Missing user ID in request context
	- Invalid request body for POST
	- GET library items successfully
	- POST add game to library successfully
	- POST add game to library with service error
	- DELETE game from library successfully
	- DELETE game from library with invalid ID
	- DELETE game from library with service error
	- Unsupported domain
*/

// MockLibraryService implements the LibraryService interface for testing
type MockLibraryService struct {
	CreateLibraryGameError        error
	MostRecentlyAddedGame        models.GameToSave

	DeleteGameError              error
	MostRecentlyDeletedGameID    int64

	GetLibraryRefactoredBFFResponseResult types.LibraryBFFRefactoredResponse
	GetLibraryRefactoredBFFResponseError   error

	InvalidateUserCacheError error
}

func (m *MockLibraryService) GetLibraryRefactoredBFFResponse(ctx context.Context, userID string) (types.LibraryBFFRefactoredResponse, error) {
	return m.GetLibraryRefactoredBFFResponseResult, m.GetLibraryRefactoredBFFResponseError
}

func (m *MockLibraryService) CreateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error {
	m.MostRecentlyAddedGame = game
	return m.CreateLibraryGameError
}

func (m *MockLibraryService) DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error {
	m.MostRecentlyDeletedGameID = gameID
	return m.DeleteGameError
}

func (m *MockLibraryService) InvalidateUserCache(ctx context.Context, userID string) error {
	return m.InvalidateUserCacheError
}

func (m *MockLibraryService) GetSingleLibraryGame(ctx context.Context, userID string, gameID int64) (types.LibraryGameItemBFFResponseFINAL, error) {
	return types.LibraryGameItemBFFResponseFINAL{}, nil
}

func (m *MockLibraryService) GetAllLibraryItemsBFF(ctx context.Context, userID string) (types.LibraryBFFResponseFINAL, error) {
	return types.LibraryBFFResponseFINAL{}, nil
}

func (m *MockLibraryService) UpdateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error {
	return nil
}

func (m *MockLibraryService) DeleteGameVersions(ctx context.Context, userID string, gameID int64, request types.BatchDeleteLibraryGameRequest) (types.BatchDeleteLibraryGameResponse, error) {
	return types.BatchDeleteLibraryGameResponse{}, nil
}

func (m *MockLibraryService) IsGameInLibraryBFF(ctx context.Context, userID string, gameID int64) (bool, error) {
	return false, nil
}

// MockAnalyticsService implements the analytics.Service interface for testing
type MockAnalyticsService struct {
	InvalidateDomainError error
}

func (m *MockAnalyticsService) InvalidateDomain(ctx context.Context, userID string, domain string) error {
	return m.InvalidateDomainError
}

func (m *MockAnalyticsService) GetAnalytics(ctx context.Context, userID string, domains []string) (map[string]any, error) {
	return make(map[string]any), nil
}

func (m *MockAnalyticsService) GetGeneralStats(ctx context.Context, userID string) (*analytics.GeneralStats, error) {
	return nil, nil
}

func (m *MockAnalyticsService) GetFinancialStats(ctx context.Context, userID string) (*analytics.FinancialStats, error) {
	return nil, nil
}

func (m *MockAnalyticsService) GetStorageStats(ctx context.Context, userID string) (*analytics.StorageStats, error) {
	return nil, nil
}

func (m *MockAnalyticsService) GetInventoryStats(ctx context.Context, userID string) (*analytics.InventoryStats, error) {
	return nil, nil
}

func (m *MockAnalyticsService) GetWishlistStats(ctx context.Context, userID string) (*analytics.WishlistStats, error) {
	return nil, nil
}

func (m *MockAnalyticsService) InvalidateDomains(ctx context.Context, userID string, domains []string) error {
	return nil
}

func TestLibraryHandler(t *testing.T) {
	// Setup base app context
	baseAppCtx := &appcontext.AppContext{
		Logger: testutils.NewTestLogger(),
	}

	// Helper fn - Create handler with given mock service
	createHandler := func(mockService *MockLibraryService, mockAnalytics *MockAnalyticsService) http.HandlerFunc {
		return CreateLibraryGame(baseAppCtx, mockService, mockAnalytics)
	}

	// Helper fn - create a request with userID in context
	createRequestWithUserID := func(method, path, body string) (*http.Request, *httptest.ResponseRecorder) {
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

	// Test cases
	/*
		GIVEN an HTTP request without a user ID in the context
		WHEN the LibraryHandler is called
		THEN an error response is returned with status 401 Unauthorized
	*/
	t.Run(`Missing user ID in request context`, func(t *testing.T) {
		// Create mock services
		mockService := &MockLibraryService{}
		mockAnalytics := &MockAnalyticsService{}

		// Create handler
		testLibraryHandler := createHandler(mockService, mockAnalytics)

		// Create request w/o userID in context
		req := httptest.NewRequest(http.MethodPost, "/library", nil)
		req.Header.Set(httputils.XRequestIDHeader, "test-request-id")
		testRecorder := httptest.NewRecorder()

		// Call handler
		testLibraryHandler.ServeHTTP(testRecorder, req)

		// Validate response
		if testRecorder.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 Unauthorized, got %d", testRecorder.Code)
		}

		// Validate error msg
		t.Logf("Response body: %s", testRecorder.Body.String())

		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(testRecorder.Body.Bytes(), &errorResponse); err != nil {
			t.Fatalf("Failed to unmarshal error message: %v", err)
		}
		if errorResponse.Error != "userID not found in request context" {
			t.Errorf("Expected error message 'userID not found in request context', got: %s", errorResponse.Error)
		}
	})

	/*
		GIVEN an HTTP POST request with a valid game in the body
		WHEN the LibraryHandler is called
		THEN the library service's CreateLibraryGame method is called with the game
		AND a JSON response with success is returned
	*/
	t.Run(`POST - Create Library Game Success`, func(t *testing.T) {
		// Create mock services
		mockService := &MockLibraryService{}
		mockAnalytics := &MockAnalyticsService{}

		// Create handler
		libraryHandler := createHandler(mockService, mockAnalytics)

		// Create request
		gameJSON := `{
			"id": 123,
			"name": "Test Game",
			"summary": "Test Summary",
			"cover_url": "https://example.com/cover.jpg",
			"release_date": 1640995200,
			"platform_locations": [
				{
					"platform_id": 1,
					"platform_name": "PlayStation 5",
					"type": "physical",
					"location": {
						"sublocation_id": "shelf-1"
					}
				}
			]
		}`
		testRequest, testRecorder := createRequestWithUserID(http.MethodPost, "/library", gameJSON)
		testRequest.Header.Set("Content-Type", "application/json")

		// Call handler
		libraryHandler.ServeHTTP(testRecorder, testRequest)

		// Validate response
		if testRecorder.Code != http.StatusCreated {
			t.Errorf("Expected status 201 for successful add, got %d", testRecorder.Code)
		}

		if mockService.MostRecentlyAddedGame.GameID != 123 {
			t.Errorf("Expected game ID 123, got %d", mockService.MostRecentlyAddedGame.GameID)
		}

		// Validate that the service was called with the correct game
		if mockService.MostRecentlyAddedGame.GameName != "Test Game" {
			t.Errorf("Expected game name 'Test Game', got %s", mockService.MostRecentlyAddedGame.GameName)
		}

		// Validate response body
		var response struct {
			Library struct {
				ID      int64  `json:"id"`
				Message string `json:"message"`
			} `json:"library"`
		}

		if err := json.Unmarshal(testRecorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response %v", err)
		}

		if response.Library.ID != 123 {
			t.Errorf("Expected game ID 123, got %d", response.Library.ID)
		}

		if response.Library.Message != "Game added to library successfully" {
			t.Errorf("Expected message 'Game added to library successfully' but instead got %s", response.Library.Message)
		}
	})

	/*
		GIVEN an HTTP POST request with an invalid body
		WHEN the LibraryHandler is called
		THEN an error response is returned
	*/
	t.Run(`POST - Add Game to Library with Invalid Request Body`, func(t *testing.T) {
		// Create mock services
		mockService := &MockLibraryService{}
		mockAnalytics := &MockAnalyticsService{}

		// Create handler
		libraryHandler := createHandler(mockService, mockAnalytics)

		// Create request with invalid JSON
		testRequest, testRecorder := createRequestWithUserID(http.MethodPost, "/library", `{invalid json}`)
		testRequest.Header.Set("Content-Type", "application/json")

		// Call handler
		libraryHandler.ServeHTTP(testRecorder, testRequest)

		// Validate response
		if testRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid body, got %d", testRecorder.Code)
		}
	})

	/*
		GIVEN an HTTP POST request with a valid game
		WHEN the library service returns an error
		THEN an error response is returned
	*/
	t.Run(`POST - Add Game to Library Service Error`, func(t *testing.T) {
		// Create mock service with error
		mockService := &MockLibraryService{
			CreateLibraryGameError: errors.New("service error"),
		}
		mockAnalytics := &MockAnalyticsService{}

		// Create handler
		libraryHandler := createHandler(mockService, mockAnalytics)

		// Create request with game in body
		gameJSON := `{
			"id": 123,
			"name": "Test Game",
			"summary": "Test Summary",
			"cover_url": "https://example.com/cover.jpg",
			"release_date": 1640995200,
			"platform_locations": [
				{
					"platform_id": 1,
					"platform_name": "PlayStation 5",
					"type": "physical",
					"location": {
						"sublocation_id": "shelf-1"
					}
				}
			]
		}`
		testRequest, testRecorder := createRequestWithUserID(http.MethodPost, "/library", gameJSON)
		testRequest.Header.Set("Content-Type", "application/json")

		// Call handler
		libraryHandler.ServeHTTP(testRecorder, testRequest)

		// Validate response
		if testRecorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500 for service error, got %d", testRecorder.Code)
		}
	})
}

