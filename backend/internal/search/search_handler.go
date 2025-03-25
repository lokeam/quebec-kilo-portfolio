package search

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/library"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/lokeam/qko-beta/internal/wishlist"
)

type DomainSearchServices map[string]SearchService

type SearchRequestBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

// NewSearchHandler returns an http.HandlerFunc which handles search requests.
func NewSearchHandler(
	appCtx *appcontext.AppContext,
	//searchServiceFactory SearchServiceFactory,
	searchServices DomainSearchServices,
	libraryService library.LibraryService,
	wishlistService wishlist.WishlistService,
) http.HandlerFunc {
	// Instantiate the concrete search service.
	appCtx.Logger.Info("NewSearchHandler created, initializing game search service", map[string]any{
		"appContext": appCtx,
		"availableDomains": getKeysFromMap(searchServices),
	})

	return func(w http.ResponseWriter, r *http.Request) {
		appCtx.Logger.Info("SearchHandler ServeHTTP called", map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
		})

		// 1. Get common request information using utility fns
		requestID := httputils.GetRequestID(r)

		// 2. Get the domain parameter. Default to "games" if not provided.
		domain := httputils.GetDomainFromRequest(r, "games")

		// 3. Get the userID from the request context
		userID := httputils.GetUserID(r)
		if userID == "" {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("userID not found in request context"),
				http.StatusUnauthorized,
			)
			return
		}
		appCtx.Logger.Info("userID found in request context", map[string]any{
			"user_id": userID,
		})


		// 4. Get the request body
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

		// 5. Get query from request body
		query := body.Query
    appCtx.Logger.Debug("SearchHandler ServeHTTP called", map[string]any{
        "request_id": r.Header.Get(httputils.XRequestIDHeader),
        "query":      query,
    })
		// Simple validation - NOTE: move this to middleware when implementing Auth0
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


		// 6. Retrieve the IGDB access token key.
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


		// 7. Optional limit parameter. Max default to 50.
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

		// 8. Build the search request.
		req := searchdef.SearchRequest{Query: query, Limit: limit}
		var result *searchdef.SearchResult

		// 9. Dispatch to the appropriate service.
		//service, err := searchServiceFactory.GetService(domain)
		service, exists := searchServices[domain]
		if !exists {
			domainErr := &types.DomainError{
				Domain: domain,
				Err:    fmt.Errorf("unsupported domain: %s", domain),
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

		// 10. Construct a unified response.
		response := searchdef.SearchResponse{
			Games: result.Games,
			Total: len(result.Games),
		}

		// 11. Check if the current search response contains items in a user's library or wishlist
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
			response.Games = []models.Game{}
		}

		appCtx.Logger.Info("Sending response to frontend", map[string]any{
			"response": response,
		})

		// 12. Return the search response as JSON.
		httputils.RespondWithJSON(w, appCtx.Logger, http.StatusOK, response)
	}
}

// Helper function to check if a game is in a list
func containsGame(games []models.Game, gameID int64) bool {
	for _, game := range games {
		if game.ID == gameID {
			return true
		}
	}
	return false
}

// Helper function to grab keys from map
func getKeysFromMap(m map[string]SearchService) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}