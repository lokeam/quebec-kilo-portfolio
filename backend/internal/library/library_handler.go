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

func RegisterLibraryRoutes(
	r chi.Router,
	appCtx *appcontext.AppContext,
	libraryServices services.DomainLibraryServices,
	analyticsService analytics.Service,
) {
	// Base routes
	r.Get("/", GetAllLibraryGames(appCtx, libraryServices))
	r.Post("/", CreateLibraryGame(appCtx, libraryServices, analyticsService))

	// Nested routes with ID
	r.Route("/games/{gameID}", func(r chi.Router) {
		r.Put("/", UpdateLibraryGame(appCtx, libraryServices, analyticsService))
		r.Delete("/", DeleteGameFromLibrary(appCtx, libraryServices, analyticsService))
	})

	// BFF route
	r.Get("/bff", GetAllLibraryItemsBFF(appCtx, libraryServices))
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
	libraryServices services.DomainLibraryServices,
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

		domain := httputils.GetDomainFromRequest(r, "games")
		service, exists := libraryServices[domain]
		if !exists {
			handleError(w, appCtx.Logger, requestID, errors.New("domain not found in libraryServices"))
			return
		}

		libraryItemsPayload, err := service.GetAllLibraryItemsBFF(r.Context(), userID)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "physical" key, DO NOT use a struct{}
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


// GetAllLibraryGames handles GET requests for listing all library games
func GetAllLibraryGames(
	appCtx *appcontext.AppContext,
	libraryServices services.DomainLibraryServices,
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

		appCtx.Logger.Info("Listing all library games", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		domain := httputils.GetDomainFromRequest(r, "games")
		service, exists := libraryServices[domain]
		if !exists {
			handleError(w, appCtx.Logger, requestID, errors.New("domain not found in libraryServices"))
			return
		}

		games, physicalLocations, digitalLocations, err := service.GetAllLibraryGames(r.Context(), userID)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		responseAdapter := NewLibraryResponseAdapter()
		adaptedResponse := responseAdapter.AdaptToLibraryResponse(
			games,
			physicalLocations,
			digitalLocations,
		)


		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"library": adaptedResponse,
		})

		jsonBytes, _ := json.Marshal(response)
		appCtx.Logger.Info("GetAllLibraryGamesResponse:", map[string]any{
				"response": string(jsonBytes),
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// AddGameToLibrary handles POST requests for adding a game to the library
func CreateLibraryGame(
	appCtx *appcontext.AppContext,
	libraryServices services.DomainLibraryServices,
	analyticsService analytics.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)
		userID := httputils.GetUserID(r)

		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request contetx", map[string]any{
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
		libraryGame :=  requestAdapter.AdaptCreateRequestToLibraryGameModel(tempAddGameRequest)

		domain := httputils.GetDomainFromRequest(r, "games")
		service, exists := libraryServices[domain]
		if !exists {
			handleError(w, appCtx.Logger, requestID, errors.New("domain not found in libraryServices"))
			return
		}

		if err := service.CreateLibraryGame(r.Context(), userID, libraryGame); err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		game, physicalLocations, digitalLocations, err := service.GetSingleLibraryGame(
			r.Context(),
			userID,
			libraryGame.GameID,
		)
		if err != nil {
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

		responseAdapter := NewLibraryResponseAdapter()
		adaptedResponse := responseAdapter.AdaptToSingleGameResponse(
			game,
			physicalLocations,
			digitalLocations,
		)

		response := httputils.NewAPIResponse(r, userID, adaptedResponse)

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// UpdateLibraryGame handles PUT requests for updating a library game
func UpdateLibraryGame(
	appCtx *appcontext.AppContext,
	libraryServices services.DomainLibraryServices,
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

		domain := httputils.GetDomainFromRequest(r, "games")
		service, exists := libraryServices[domain]
		if !exists {
			handleError(w, appCtx.Logger, requestID, errors.New("domain not found in libraryServices"))
			return
		}

		game, physicalLocations, digitalLocations, err := service.GetSingleLibraryGame(
			r.Context(),
			userID,
			gameIDint64,
		)
		if err != nil {
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

		responseAdapter := NewLibraryResponseAdapter()
		adaptedResponse := responseAdapter.AdaptToSingleGameResponse(
			game,
			physicalLocations,
			digitalLocations,
		)

		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"library": adaptedResponse,
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
	libraryServices services.DomainLibraryServices,
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

		domain := httputils.GetDomainFromRequest(r, "games")
		service, exists := libraryServices[domain]
		if !exists {
			handleError(w, appCtx.Logger, requestID, errors.New("domain not found in libraryServices"))
			return
		}

		if err := service.DeleteLibraryGame(
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
			adaptedResponse := responseAdapter.AdaptToLibraryResponse(
				games,
				physicalLocations,
				digitalLocations,
			)

			// NOTE: DOUBLE CHECK THIS - LACK OF map wrapping MAY BE THE CAUSE OF FRONTEND ISSUES
			response := httputils.NewAPIResponse(r, userID, adaptedResponse)

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

			adaptedResponse := responseAdapter.AdaptToSingleGameResponse(
				game,
				physicalLocations,
				digitalLocations,
			)

			response := httputils.NewAPIResponse(r, userID, adaptedResponse)

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
			adaptedResponse := responseAdapter.AdaptToSingleGameResponse(
				game,
				physicalLocations,
				digitalLocations,
			)

			response := httputils.NewAPIResponse(r, userID, adaptedResponse)

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
			data := struct {
				ID      int64  `json:"id"`
				Message string `json:"message"`
			}{
				ID:      gameIDint64,
				Message: "Game removed from library successfully",
			}

			response := httputils.NewAPIResponse(r, userID, data)

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

