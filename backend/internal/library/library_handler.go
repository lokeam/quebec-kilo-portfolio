package library

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
)

type LibraryRequestBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

type AddGameRequest struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Summary    string `json:"summary"`
	CoverURL   string `json:"cover_url"`
	ReleaseDate int64  `json:"release_date"`
}

func RegisterLibraryRoutes(
	r chi.Router,
	appCtx *appcontext.AppContext,
	libraryService services.LibraryService,
	analyticsService analytics.Service,
) {
	// Base routes
	r.Post("/", CreateLibraryGame(appCtx, libraryService, analyticsService))

	// Nested routes with ID
	r.Route("/games/{gameID}", func(r chi.Router) {
		r.Put("/", UpdateLibraryGame(appCtx, libraryService, analyticsService))
		r.Delete("/", DeleteGameFromLibrary(appCtx, libraryService, analyticsService))
	})

	// BFF route
	r.Get("/bff", GetAllLibraryItemsBFF(appCtx, libraryService))
}

// helper fn to standardize error handling
func handleError(
	w http.ResponseWriter,
	logger interfaces.Logger,
	requestID string,
	err error,
) {
	statusCode := GetStatusCodeForError(err)
	httputils.RespondWithError(
		httputils.NewResponseWriterAdapter(w),
		logger,
		requestID,
		err,
		statusCode,
	)
}

// GetAllLibraryItemsBFF handles GET requests for the /library/bff page
func GetAllLibraryItemsBFF(
	appCtx *appcontext.AppContext,
	libraryService services.LibraryService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)

		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			handleError(w, appCtx.Logger, requestID, errors.New("userID not found in request context"))
			return
		}

		appCtx.Logger.Info("Logging requestID and userID for library BFF request", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		libraryItemsPayload, err := libraryService.GetAllLibraryItemsBFF(r.Context(), userID)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Use standard response format
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"library": libraryItemsPayload,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// CreateLibraryGame handles POST requests for adding a game to the library
func CreateLibraryGame(
	appCtx *appcontext.AppContext,
	libraryService services.LibraryService,
	analyticsService analytics.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)
		userID := httputils.GetUserID(r)

		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			handleError(w, appCtx.Logger, requestID, errors.New("userID not found in request context"))
			return
		}

		appCtx.Logger.Info("Adding game to library", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		// Log the raw request body for debugging pissy JSON errors
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Note: NEED TO RESET body AFTER READING

		appCtx.Logger.Info("Request body received", map[string]any{
			"requestID": requestID,
			"body":      string(bodyBytes),
		})

		var tempAddGameRequest types.CreateLibraryGameRequest
		if err := json.NewDecoder(r.Body).Decode(&tempAddGameRequest); err != nil {
			handleError(w, appCtx.Logger, requestID, errors.New("invalid request body"))
			return
		}

		requestAdapter := NewLibraryRequestAdapter()
		libraryGame := requestAdapter.AdaptCreateRequestToLibraryGameModel(tempAddGameRequest)

		if err := libraryService.CreateLibraryGame(r.Context(), userID, libraryGame); err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Invalidate BFF cache
		if err := libraryService.InvalidateUserCache(r.Context(), userID); err != nil {
			appCtx.Logger.Warn("Failed to invalidate library BFF cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		// Invalidate analytics cache for inventory domain
		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"library": map[string]any{
				"id": libraryGame.GameID,
				"message": "Game added to library successfully",
			},
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusCreated,
			response,
		)
	}
}

// UpdateLibraryGame handles PUT requests for updating a library game
func UpdateLibraryGame(
	appCtx *appcontext.AppContext,
	libraryService services.LibraryService,
	analyticsService analytics.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)
		userID := httputils.GetUserID(r)

		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			handleError(w, appCtx.Logger, requestID, errors.New("userID not found in request context"))
			return
		}

		// Log the raw request body for debugging pissy JSON errors
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body after reading

		urlParts := strings.Split(r.URL.Path, "/")
		if len(urlParts) < 3 {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("invalid path"),
				http.StatusBadRequest,
			)
			return
		}

		gameIDStr := urlParts[len(urlParts)-1]
		gameIDint64, err := strconv.ParseInt(gameIDStr, 10, 64)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, errors.New("invalid ID: must be a number"))
			return
		}

		appCtx.Logger.Info("Updating library game", map[string]any{
			"requestID": requestID,
			"userID":    userID,
			"gameID":    gameIDint64,
		})

		// Parse request body
		var tempPutGameRequest types.UpdateLibraryGameRequest
		if err := json.NewDecoder(r.Body).Decode(&tempPutGameRequest); err != nil {
			handleError(w, appCtx.Logger, requestID, errors.New("invalid request body"))
			return
		}

		requestAdapter := NewLibraryRequestAdapter()
		libraryGame := requestAdapter.AdaptUpdateRequestToLibraryGameModel(tempPutGameRequest)
		libraryGame.GameID = gameIDint64

		// Get the current game state
		game, err := libraryService.GetSingleLibraryGame(
			r.Context(),
			userID,
			gameIDint64,
		)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Update the game
		if err := libraryService.UpdateLibraryGame(r.Context(), userID, libraryGame); err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Invalidate analytics cache for inventory domain
		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"library": map[string]any{
				"game": game,
			},
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// DeleteGameFromLibrary handles DELETE requests for deleting a game from the library
func DeleteGameFromLibrary(
	appCtx *appcontext.AppContext,
	libraryService services.LibraryService,
	analyticsService analytics.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)
		userID := httputils.GetUserID(r)

		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			handleError(w, appCtx.Logger, requestID, errors.New("userID not found in request context"))
			return
		}

		// Extract game ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("invalid path"),
				http.StatusBadRequest,
			)
			return
		}

		// Grab last part of the path as the gameID
		gameIDStr := parts[len(parts)-1]
		gameIDint64, err := strconv.ParseInt(gameIDStr, 10, 64)
		if err != nil {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("invalid ID: must be a number"),
				http.StatusBadRequest,
			)
			return
		}

		appCtx.Logger.Info("Deleting game from library", map[string]any{
			"requestID": requestID,
			"userID":    userID,
			"gameID":    gameIDint64,
		})

		if err := libraryService.DeleteLibraryGame(
			r.Context(),
			userID,
			gameIDint64,
		); err != nil {
			if errors.Is(err, ErrGameNotFound) {
				handleError(w, appCtx.Logger, requestID, err)
				return
			}
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Invalidate analytics cache for inventory domain
		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"message": "Game removed from library successfully",
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}
