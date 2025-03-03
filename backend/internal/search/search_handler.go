package search

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/library"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/lokeam/qko-beta/internal/wishlist"
	authMiddleware "github.com/lokeam/qko-beta/server/middleware"
)

type SearchRequestBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

// NewSearchHandler returns an http.HandlerFunc which handles search requests.
func NewSearchHandler(
	appCtx *appcontext.AppContext,
	searchServiceFactory SearchServiceFactory,
	libraryService library.LibraryService,
	wishlistService wishlist.WishlistService,
) http.HandlerFunc {
	// Instantiate the concrete search service.
	appCtx.Logger.Info("NewSearchHandler created, initializing game search service", map[string]any{
		"appContext": appCtx,
	})

	return func(w http.ResponseWriter, r *http.Request) {
		appCtx.Logger.Info("SearchHandler ServeHTTP called", map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
		})

		// 1. Get the search query parameter.
		requestID := r.Header.Get(httputils.XRequestIDHeader)

    bodyBytes, _ := io.ReadAll(r.Body)
    appCtx.Logger.Debug("Request body", map[string]any{
        "body": string(bodyBytes),
    })
    r.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Reset the body for parsing

		var body SearchRequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("invalid request body"),
				http.StatusBadRequest,
			)
			return
		}
		defer r.Body.Close()

		query := body.Query
    appCtx.Logger.Debug("SearchHandler ServeHTTP called", map[string]any{
        "request_id": r.Header.Get(httputils.XRequestIDHeader),
        "query":      query,
    })
		if query == "" {
			err := errors.New("search query is required")
			httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					r.Header.Get(httputils.XRequestIDHeader),
					err,
					http.StatusBadRequest,
			)
			return
		}


		// 2. Retrieve the userID from the request context
		userID, ok := r.Context().Value(authMiddleware.UserIDKey).(string)
		if !ok {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})

			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("userID not found in request context"),
				http.StatusUnauthorized,
			)
		}
		appCtx.Logger.Info("userID found in request context", map[string]any{
			"user_id": userID,
		})

		// 3. Retrieve the IGDB access token key.
		twitchAccessTokenKey, err := appCtx.Config.IGDB.GetAccessTokenKey()
		if err != nil || twitchAccessTokenKey == "" {
			err := errors.New("failed to retrieve token")
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				err,
				http.StatusInternalServerError,
			)
			return
		}

		appCtx.Logger.Info("Successfully retrieved Twitch token", map[string]any{
			"request_id": requestID,
			"token_key":  twitchAccessTokenKey,
		})

		// 4. Get the domain parameter. Default to "games" if not provided.
		domain := r.URL.Query().Get("domain")
		if domain == "" {
			domain = "games"
		}

		// 5. Optional limit parameter. Max default to 50.
		limit := 5 // DEBUG: cut this down to 5 for now
		if body.Limit > 0 {
			limit = body.Limit
		} else if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if parsed, err := strconv.Atoi(limitStr); err == nil {
				limit = parsed
			}
		}

		appCtx.Logger.Info("Handling search", map[string]any{
			"query":      query,
			"limit":      limit,
			"domain":     domain,
			"request_id": requestID,
		})

		// 6. Build the search request.
		req := searchdef.SearchRequest{Query: query, Limit: limit}
		var result *searchdef.SearchResult

		// 7. Dispatch to the appropriate service.
		service, err := searchServiceFactory.GetService(domain)
		if err != nil {

			domainErr := &types.DomainError{
				Domain: domain,
				Err:    err,
			}

			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				domainErr,
				http.StatusBadRequest,
			)
			return
		}

		result, err = service.Search(r.Context(), req)
		if err != nil {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				err,
			)
			return
		}

		// 8. Construct a unified response.
		response := searchdef.SearchResponse{
			Games: result.Games,
			Total: len(result.Games),
		}

		// 9. Check if the current search response contains items in a user's library or wishlist
		library, err := libraryService.GetLibraryItems(r.Context(), userID)
		if err != nil {
			appCtx.Logger.Error("Failed to fetch items in user's library", map[string]any{
				"request_id": requestID,
				"error":      err,
			})
		}

		wishlist, err := wishlistService.GetWishlistItems(r.Context(), userID)
		if err != nil {
			appCtx.Logger.Error("Failed to fetch user wishlist", map[string]any{
				"request_id": requestID,
				"error":      err,
			})
		}

		for i, game := range response.Games {
			response.Games[i].IsInLibrary = containsGame(library, game.ID)
			response.Games[i].IsInWishlist = containsGame(wishlist, game.ID)
		}


		// If no games found, return an empty response
		if response.Games == nil {
			response.Games = []types.Game{}
		}

		appCtx.Logger.Info("Sending response to frontend", map[string]any{
			"response": response,
		})

		// 10.Return the search response as JSON.
		httputils.RespondWithJSON(w, appCtx.Logger, http.StatusOK, response)
	}
}

// Helper function to check if a game is in a list
func containsGame(games []types.Game, gameID int64) bool {
	for _, game := range games {
		if game.ID == gameID {
			return true
		}
	}
	return false
}
