package library

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/lokeam/qko-beta/internal/appcontext_test"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/shared/constants"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
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

type MockLibraryService struct {
	CreateLibraryGameError        error
	MostRecentlyAddedGame        *models.Game

	DeleteGameError              error
	MostRecentlyDeletedGameID    int64

	GetGameByIDResult            models.Game
	GetGameByIDError             error

	UpdateGameError              error
}

// Mock library service methods
func (m *MockLibraryService) GetGameByID(ctx context.Context, userID string, gameID int64) (models.Game, error) {
	return m.GetGameByIDResult, m.GetGameByIDError
}

func (m *MockLibraryService) CreateLibraryGame(ctx context.Context, userID string, game models.Game) error {
	m.MostRecentlyAddedGame = &game // NOTE: I don't understand why this works
	return m.CreateLibraryGameError
}

func (m *MockLibraryService) DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error {
	m.MostRecentlyDeletedGameID = gameID
	return m.DeleteGameError
}

func (m *MockLibraryService) UpdateLibraryGame(ctx context.Context, userID string, game models.Game) error {
	return m.UpdateGameError
}

// Tests
func TestLibraryHandler(t *testing.T) {
	// Set up base app context for testing
	baseAppCtx := appcontext_test.NewTestingAppContext("test-token", nil)

	// Helper fn - Create handler with given mock service
	createHandler := func(mockService *MockLibraryService) http.Handler {
		libraryServices := make(DomainLibraryServices)
		libraryServices["games"] = mockService
		return NewLibraryHandler(baseAppCtx, libraryServices)
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

	// Helper fn - create Chi router with the handler
	setupChiRouter := func(handler http.HandlerFunc) *chi.Mux {
		testRouter := chi.NewRouter()
		testRouter.Get("/library", handler)
		testRouter.Post("/library", handler)
		testRouter.Route("/library/games", func(r chi.Router) {
			r.Delete("/{id}", handler)
		})
		return testRouter
	}

	// Test cases
	/*
		GIVEN an HTTP request without a user ID in the context
		WHEN the LibraryHandler is called
		THEN an error response is returned with status 401 Unauthorized
	*/
	t.Run(`Missing user ID in request context`, func(t *testing.T) {
		// Create mock service
		mockService := &MockLibraryService{}

		// Create handler
		testLibraryHandler := createHandler(mockService)

		// Create request w/o userID in context
		req := httptest.NewRequest(http.MethodGet, "/library", nil)
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
		GIVEN an HTTP GET request with a valid user ID
		WHEN the LibraryHandler is called
		THEN the library service's GetAllLibraryGames method is called
		AND a JSON response with the library items is returned

		NOTE: REPLACE THIS WITH BFF TESTS
	*/

	/*
		GIVEN an HTTP POST request with a valid game in the body
		WHEN the LibraryHandler is called
		THEN the library service's CreateLibraryGame method is called with the game
		AND a JSON response with success is returned
	*/
	t.Run(`POST - Library Items Service Error`, func(t *testing.T) {
		// Create mock service with error
		mockService := &MockLibraryService{}

		// Create handler
		libraryHandler := createHandler(mockService)

		// Create request
		gameJSON := `{
			"id": 123,
			"name": "Test Game",
			"summary": "Test Summary",
			"cover_url": "https://example.com/cover.jpg",
			"first_release_date": "2024-01-01",
			"platform_names": ["Platform 1", "Platform 2"],
			"genre_names": ["Genre 1", "Genre 2"],
			"theme_names": ["Theme 1", "Theme 2"]
		}`
		testRequest, testRecorder := createRequestWithUserID(http.MethodPost, "/library", gameJSON)
		testRequest.Header.Set("Content-Type", "application/json")

		// Call handler
		libraryHandler.ServeHTTP(testRecorder, testRequest)

		// Validate response
		if testRecorder.Code != http.StatusOK {
			t.Errorf("Expected status 200 for successful add, got %d", testRecorder.Code)
		}

		if mockService.MostRecentlyAddedGame.ID != 123 {
			t.Errorf("Expected game ID 123, got %d", mockService.MostRecentlyAddedGame.ID)
		}

		// Validate that the service was called with the correct game
		if mockService.MostRecentlyAddedGame.Name != "Test Game" {
			t.Errorf("Expected game name 'Test Game', got %s", mockService.MostRecentlyAddedGame.Name)
		}

		// Validate response body
		var response struct {
			Success bool `json:"success"`
			Game    struct {
				ID    int64   `json:"id"`
				Name  string  `json:"name"`
			} `json:"game"`
		}

		if err := json.Unmarshal(testRecorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response %v", err)
		}

		if !response.Success {
			t.Errorf("Expected success to be true")
		}

		if response.Game.ID != 123 {
			t.Errorf("Expected game ID 123, got %d", response.Game.ID)
		}

		if response.Game.Name != "Test Game" {
			t.Errorf("Expected game name 'Test Game' but instead got %s", response.Game.Name)
		}
	})

	/*
		GIVEN an HTTP POST request with an invalid body
		WHEN the LibraryHandler is called
		THEN an error response is returned
	*/
	t.Run(`POST - Add Game to Library with Invalid Request Body`, func(t *testing.T) {
		// Create mock service
		mockService := &MockLibraryService{}

		// Create handler
		libraryHandler := createHandler(mockService)

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

		// Create handler
		libraryHandler := createHandler(mockService)

		// Create request with game IN body
		gameJSON := `{
			"id": 123,
			"name": "Test Game",
			"summary": "Test Summary",
			"cover_url": "https://example.com/cover.jpg",
			"first_release_date": "2024-01-01"
		}`
		testRequest, testRecorder := createRequestWithUserID(http.MethodPost, "/library", gameJSON)
		testRequest.Header.Set("Content-Type", "application/json")

		// Call handler
		libraryHandler.ServeHTTP(testRecorder, testRequest)

		// Validate response
		if testRecorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500 for service error but instead got %d", testRecorder.Code)
		}
	})

	/*
		GIVEN an HTTP DELETE request with a valid game ID
		WHEN the LibraryHandler is called
		THEN the library service's DeleteLibraryGame method is called with the ID
		AND a JSON response with success is returned
	*/
	t.Run(`DELETE - Delete Game from Library Successfully`, func(t *testing.T) {
		// Create mock service
		mockService := &MockLibraryService{}

		// Create handler
		libraryHandler := createHandler(mockService)

		// Create Chi router with handler
		testRouter := setupChiRouter(libraryHandler.(http.HandlerFunc)) // NOTE: Chi router expects http.HandlerFunc, not just any http.Handler

		// Create request
		testRequest, testRecorder := createRequestWithUserID(http.MethodDelete, "/library/games/123", "")

		// Call router
		testRouter.ServeHTTP(testRecorder, testRequest)

		// Validate response
		if testRecorder.Code != http.StatusOK {
			t.Errorf("Expected status 200 for successful delete, got %d", testRecorder.Code)
		}

		// Validate that service was called with correct ID
		if mockService.MostRecentlyDeletedGameID != 123 {
			t.Errorf("Expected DeleteLibraryGame to be called with ID 123, got %d", mockService.MostRecentlyDeletedGameID)
		}

		// Validate response body
		var response struct {
			Success bool  `json:"success"`
			ID      int64 `json:"id"`
		}

		if err := json.Unmarshal(testRecorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if !response.Success {
			t.Errorf("Expected success to be true")
		}

		if response.ID != 123 {
			t.Errorf("Expected ID 123, got %d", response.ID)
		}
	})

	/*
		GIVEN an HTTP DELETE request with an invalid game ID
		WHEN the LibraryHandler is called
		THEN an error response is returned
	*/
	t.Run(`DELETE - Game from Library Invalid ID`, func(t *testing.T) {
		// Create mock service
		mockService := &MockLibraryService{}

		// Create handler
		libraryHandler := createHandler(mockService)

		// Create Chi router with the handler
		testRouter := setupChiRouter(libraryHandler.(http.HandlerFunc))

		// Create reqeust with invalid ID
		testRequest, testRecorder := createRequestWithUserID(http.MethodDelete, "/library/games/invalid", "")

		// Call router
		testRouter.ServeHTTP(testRecorder, testRequest)

		// Validate response
		if testRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid ID, got %d", testRecorder.Code)
		}
	})

	/*
		GIVEN an HTTP DELETE request with a valid game ID
		WHEN the library service returns an error
		THEN an error response is returned
	*/
	t.Run(`DELETE - Game from Library Service Error`, func(t *testing.T) {
		// Create mock service with error
		mockService := &MockLibraryService{
			DeleteGameError: errors.New("service error"),
		}

		// Create Chi router with handler
		libraryHandler := createHandler(mockService)
		testRouter := setupChiRouter(libraryHandler.(http.HandlerFunc))

		// Create request
		testRequest, testRecorder := createRequestWithUserID(http.MethodDelete, "/library/games/123", "")

		// Call router
		testRouter.ServeHTTP(testRecorder, testRequest)

		// Validate response
		if testRecorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500 for service error but instead got %d", testRecorder.Code)
		}

	})

	/*
		GIVEN an HTTP DELETE request with a valid game ID
		WHEN the library service returns a "not found" error
		THEN a 404 Not Found response is returned
	*/
	t.Run(`DELETE - Game from Library Not Found`, func(t *testing.T) {
		// Create mock service with not found error
		mockService := &MockLibraryService{
			DeleteGameError: ErrGameNotFound,
		}

		// Create handler
		libraryHandler := createHandler(mockService)

		// Create Chi router with handler
		testRouter := setupChiRouter(libraryHandler.(http.HandlerFunc))

		// Create request
		testRequest, testRecorder := createRequestWithUserID(http.MethodDelete, "/library/games/123", "")

		// Call router
		testRouter.ServeHTTP(testRecorder, testRequest)

		// Validate response
		if testRecorder.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 for not found error but instead got %d", testRecorder.Code)
		}
	})

	/*
		GIVEN an HTTP request with an unsupported domain
		WHEN the LibraryHandler is called
		THEN an error response is returned
	*/
	t.Run(`Unsupported Domain`, func(t *testing.T) {
		// Create mock service
		mockService := &MockLibraryService{}

		// Create handler
		libraryHandler := createHandler(mockService)

		// Create request with unsupported domain
		testRequest, testRecorder := createRequestWithUserID(http.MethodGet, "/library?domain=unsupported", "")

		// Call handler
		libraryHandler.ServeHTTP(testRecorder, testRequest)

		// Validate response
		if testRecorder.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 for unsupported domain but instead got %d", testRecorder.Code)
		}
	})

	/*
		GIVEN an HTTP request with an unsupported method
		WHEN the LibraryHandler is called
		THEN a Method Not Allowed response is returned
	*/
	// Create mock service
	mockService := &MockLibraryService{}

	// Create handler
	libraryHandler := createHandler(mockService)

	// Create request with unsupported method
	testRequest, testRecorder := createRequestWithUserID(http.MethodPut, "/library", "")

	// Call handler
	libraryHandler.ServeHTTP(testRecorder, testRequest)

	// Validate response
	if testRecorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for unsupported method but instead got %d", testRecorder.Code)
	}
}

