package search

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

type SearchHandler struct {
	appContext *appcontext.AppContext
	searchService services.SearchService
	libraryService services.LibraryService
	wishlistService services.WishlistService
}

type SearchRequestBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

func NewSearchHandler(
	appCtx *appcontext.AppContext,
	searchService services.SearchService,
	libraryService services.LibraryService,
	wishlistService services.WishlistService,
) *SearchHandler {
	return &SearchHandler{
		appContext: appCtx,
		searchService: searchService,
		libraryService: libraryService,
		wishlistService: wishlistService,
	}
}

// RegisterSearchRoutes registers all search routes
func RegisterSearchRoutes(
	r chi.Router,
	appCtx *appcontext.AppContext,
	searchService services.SearchService,
	libraryService services.LibraryService,
	wishlistService services.WishlistService,
) {
	handler := NewSearchHandler(appCtx, searchService, libraryService, wishlistService)
	r.Post("/", handler.Search)
	r.Get("/bff", handler.GetGameStorageLocationsBFF)
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	h.appContext.Logger.Info("SearchHandler ServeHTTP called", map[string]any{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	// 1. Get common request information using utility fns
	requestID := httputils.GetRequestID(r)

	// 2. Get the userID from the request context
	userID := httputils.GetUserID(r)
	if userID == "" {
		h.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}
	h.appContext.Logger.Info("userID found in request context", map[string]any{
		"user_id": userID,
		"request_id": requestID,
	})

	// 3. Get the request body
	bodyBytes, _ := io.ReadAll(r.Body)
	h.appContext.Logger.Debug("Request body", map[string]any{
		"request_id": requestID,
		"body": string(bodyBytes),
	})
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Reset the body for parsing

	var body SearchRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.handleError(w, requestID, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 4. Get query from request body
	query := body.Query
	h.appContext.Logger.Debug("SearchHandler ServeHTTP called", map[string]any{
			"request_id": requestID,
			"query":      query,
	})
	// 5. Simple validation - NOTE: move this to middleware when implementing Auth0
	if query == "" {
		h.handleError(w, requestID, errors.New("search query is required"), http.StatusBadRequest)
		return
	}

	// 6. Retrieve the IGDB access token key.
	twitchAccessTokenKey, err := h.appContext.Config.IGDB.GetAccessTokenKey()
	if err != nil || twitchAccessTokenKey == "" {
		h.handleError(w, requestID, errors.New("failed to retrieve token"), http.StatusInternalServerError)
		return
	}

	h.appContext.Logger.Info("Successfully retrieved Twitch token", map[string]any{
		"request_id": requestID,
		"token_key":  twitchAccessTokenKey,
	})

	// 7. Optional limit parameter. Max default to 50.
	limit := 5 // DEBUG: cut this down to 5 for now
	if body.Limit > 0 {
		limit = body.Limit
	} else if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	h.appContext.Logger.Info("Handling search", map[string]any{
		"query":      query,
		"limit":      limit,
		"request_id": requestID,
	})

	// 8. Build the search request.
	req := searchdef.SearchRequest{Query: query, Limit: limit}
	var result *searchdef.SearchResult
	result, err = h.searchService.Search(r.Context(), req)
	if err != nil {
		h.handleError(w, requestID, err, http.StatusInternalServerError)
		return
	}

	// 9. Construct a unified response.
	response := searchdef.SearchResponse{
		Games: result.Games,
		Total: len(result.Games),
	}

	// 10. Check if the current search response contains items in a user's library or wishlist
	for i := 0; i < len(response.Games); i++ {
		game := response.Games[i]
		// Check if game is in library
		isInLibrary, err := h.libraryService.IsGameInLibraryBFF(r.Context(), userID, game.ID)
		if err != nil {
			h.appContext.Logger.Error("Failed to check if game is in library", map[string]any{
				"request_id": requestID,
				"game_id":    game.ID,
				"error":      err,
			})
		}
		response.Games[i].IsInLibrary = isInLibrary
	}

	wishlist, err := h.wishlistService.GetWishlistItems(r.Context(), userID)
	if err != nil {
		h.appContext.Logger.Error("Failed to fetch user wishlist", map[string]any{
			"request_id": requestID,
			"error":      err,
		})
	}

	for i, game := range response.Games {
		response.Games[i].IsInWishlist = containsGame(wishlist, game.ID)
	}

	// If no games found, return an empty response
	if response.Games == nil {
		response.Games = []models.Game{}
	}

	h.appContext.Logger.Info("Sending response to frontend", map[string]any{
		"response": response,
	})

	// 11. Return the search response as JSON with the expected structure
	apiResponse := httputils.NewAPIResponse(r, userID, response)
	httputils.RespondWithJSON(w, h.appContext.Logger, http.StatusOK, apiResponse)
}

func (h *SearchHandler) GetGameStorageLocationsBFF(w http.ResponseWriter, r *http.Request) {
	requestID := httputils.GetRequestID(r)

	userID := httputils.GetUserID(r)
	if userID == "" {
		h.appContext.Logger.Error("userID not found in request context", map[string]any{
			"request_id": requestID,
		})
		h.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	h.appContext.Logger.Info("Getting game storage locations", map[string]any{
		"requestID": requestID,
		"userID":    userID,
		})

	locations, err := h.searchService.GetAllGameStorageLocationsBFF(r.Context(), userID)
	if err != nil {
		h.handleError(w, requestID, err, http.StatusInternalServerError)
		return
	}

	// Use standard response format
	// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "physical" key, DO NOT use a struct{}
	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"storage_locations": locations,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		h.appContext.Logger,
		http.StatusOK,
		response,
	)
}

//  handleError standardizes error handling by formatting and sending HTTP error responses with appropriate status codes and logging
func (h *SearchHandler) handleError(w http.ResponseWriter, requestID string, err error, statusCode int) {
	httputils.RespondWithError(
		httputils.NewResponseWriterAdapter(w),
		h.appContext.Logger,
		requestID,
		err,
		statusCode,
	)
}

// Helper function to check if a game is in a list
func containsGame(games []models.GameToSave, gameID int64) bool {
	for _, game := range games {
		if game.GameID == gameID {
			return true
		}
	}
	return false
}