package library

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/constants"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
)

// Use the DomainLibraryServices from the services package
type DomainLibraryServices = services.DomainLibraryServices

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

func NewLibraryHandler(
	appCtx *appcontext.AppContext,
	libraryServices DomainLibraryServices,
) http.HandlerFunc {
	adapter := NewLibraryRequestAdapter()
	responseAdapter := NewLibraryResponseAdapter()

	return func(w http.ResponseWriter, r *http.Request) {
		appCtx.Logger.Info("LibraryHandler ServeHTTP called", map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
		})

		// 1. Grab request ID
		requestID := r.Header.Get(httputils.XRequestIDHeader)
		appCtx.Logger.Info("requestID", map[string]any{
			"requestID": requestID,
		})

		// 2. Retrieve userID from the request context
		userID, ok := r.Context().Value(constants.UserIDKey).(string)
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
			return
		}
		appCtx.Logger.Info("userID found in request context", map[string]any{
			"user_id": userID,
		})

		// Get the domain parameter, Default to "games" if not provided
		domain := httputils.GetDomainFromRequest(r, "games")

		// 3. Dispatch library service to add item to library
		service, exists := libraryServices[domain]
		if !exists {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("domain not found in libraryServices"),
				http.StatusNotFound,
			)
			return
		}

		// Handle different HTTP methods
		switch r.Method {
		case http.MethodGet:
			appCtx.Logger.Info("GET request received", map[string]any{
				"requestID": requestID,
			})

			// Get user's library from service
			games, physicalLocations, digitalLocations, err := service.GetAllLibraryGames(r.Context(), userID)
			if err != nil {
				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					err,
					http.StatusInternalServerError,
				)
				return
			}

			// Use response adapter to format response
			response := responseAdapter.AdaptToLibraryResponse(
				games,
				physicalLocations,
				digitalLocations,
			)

			httputils.RespondWithJSON(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				http.StatusOK,
				response,
			)
			return

		case http.MethodPut:
			appCtx.Logger.Info("PUT request received", map[string]any{
				"requestID": requestID,
			})

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
					httputils.RespondWithError(
							httputils.NewResponseWriterAdapter(w),
							appCtx.Logger,
							requestID,
							errors.New("invalid ID: must be a number"),
							http.StatusBadRequest,
					)
					return
			}

			// Parse request body
			var tempPutGameRequest types.UpdateLibraryGameRequest
			if err := json.NewDecoder(r.Body).Decode(&tempPutGameRequest); err != nil {
				// Enhanced error logging for more context
				httputils.LogJSONError(appCtx.Logger, requestID, err, bodyBytes)

				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					errors.New("invalid request body"),
					http.StatusBadRequest,
				)
				return
			}

			libraryGame := adapter.AdaptUpdateRequestToLibraryGameModel(tempPutGameRequest)
			libraryGame.GameID = gameIDint64

			// Update game
			if err := service.UpdateLibraryGame(r.Context(), userID, libraryGame); err != nil {
				if errors.Is(err, ErrGameNotFound) {
						httputils.RespondWithError(
								httputils.NewResponseWriterAdapter(w),
								appCtx.Logger,
								requestID,
								err,
								http.StatusNotFound,
						)
						return
				}
				httputils.RespondWithError(
						httputils.NewResponseWriterAdapter(w),
						appCtx.Logger,
						requestID,
						err,
						http.StatusInternalServerError,
				)
				return
			}

			// Get updated game
			game, physicalLocations, digitalLocations, err := service.GetSingleLibraryGame(
				r.Context(),
				userID,
				gameIDint64,
			)
			if err != nil {
					httputils.RespondWithError(
							httputils.NewResponseWriterAdapter(w),
							appCtx.Logger,
							requestID,
							err,
							http.StatusInternalServerError,
					)
					return
			}

			// Return updated game
			response := responseAdapter.AdaptToSingleGameResponse(
				game,
				physicalLocations,
				digitalLocations,
			)

			httputils.RespondWithJSON(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				http.StatusOK,
				response,
			)
			return

		case http.MethodPost:
			appCtx.Logger.Info("POST request received", map[string]any{
				"requestID": requestID,
				"content_type": r.Header.Get("Content-Type"),
			})

			// Log the raw request body for debugging pissy JSON errors
			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Note: NEED TO RESET body AFTER READING

			appCtx.Logger.Info("Request body received", map[string]any{
				"requestID": requestID,
				"body":      string(bodyBytes),
			})

			// Parse request body using a temporary struct that matches the API contract in library-types.ts
			var tempGame types.CreateLibraryGameRequest

			if err := json.NewDecoder(r.Body).Decode(&tempGame); err != nil {
				// Enhanced error logging for more context
				httputils.LogJSONError(appCtx.Logger, requestID, err, bodyBytes)

				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					errors.New("invalid request body"),
					http.StatusBadRequest,
				)
				return
			}

			libraryGame := adapter.AdaptCreateRequestToLibraryGameModel(tempGame)

			// Use service to add game to library
			if err := service.CreateLibraryGame(
				r.Context(),
				userID,
				libraryGame,
			); err != nil {
				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					err,
					http.StatusInternalServerError,
				)
				return
			}

			// Get the newly created game to return in response
			game, physicalLocations, digitalLocations, err := service.GetSingleLibraryGame(
				r.Context(),
				userID,
				libraryGame.GameID,
			)
			if err != nil {
				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					err,
					http.StatusInternalServerError,
				)
				return
			}

			// Use response adapter to format single game response
			response := responseAdapter.AdaptToSingleGameResponse(
				game,
				physicalLocations,
				digitalLocations,
			)

			httputils.RespondWithJSON(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				http.StatusOK,
				response,
			)
			return

		case http.MethodDelete:
			appCtx.Logger.Info("DELETE request received", map[string]any{
				"requestID": requestID,
				"path":      r.URL.Path,
			})

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

			// Call service to delete game from library
			if err := service.DeleteLibraryGame(
				r.Context(),
				userID,
				gameIDint64,
			); err != nil {
				// Check for specific error types
				if errors.Is(err, ErrGameNotFound) {
					httputils.RespondWithError(
						httputils.NewResponseWriterAdapter(w),
						appCtx.Logger,
						requestID,
						err,
						http.StatusNotFound,
					)
				} else {
					httputils.RespondWithError(
						httputils.NewResponseWriterAdapter(w),
						appCtx.Logger,
						requestID,
						err,
						http.StatusInternalServerError,
					)
				}
				return
			}

			// Format response
			response := struct {
				Success    bool    `json:"success"`
				ID         int64   `json:"id"`
			}{
				Success:   true,
				ID:        gameIDint64,
			}

			httputils.RespondWithJSON(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				http.StatusOK,
				response,
			)
			return

		default:
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("method not allowed"),
				http.StatusMethodNotAllowed,
			)
			return
		}
	}
}


