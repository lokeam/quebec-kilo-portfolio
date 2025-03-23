package library

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
	authMiddleware "github.com/lokeam/qko-beta/server/middleware"
)

type DomainLibraryServices map[string]LibraryService

type LibraryRequestBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

type AddGameRequest struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Summary    string `json:"summary"`
	CoverURL   string `json:"cover_url"`
	ReleaseDate int64 `json:"release_date"`
}

func NewLibraryHandler(
	appCtx *appcontext.AppContext,
	libraryServices DomainLibraryServices,
) http.HandlerFunc {

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
		}

		switch r.Method {
		case http.MethodPost:
			appCtx.Logger.Info("POST request received", map[string]any{
				"requestID": requestID,
			})

			// Parse request body
			var gameRequest AddGameRequest
			if err := json.NewDecoder(r.Body).Decode(&gameRequest); err != nil {
				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					errors.New("invalid request body"),
					http.StatusBadRequest,
				)
				return
			}

			// Create game object from the request
			gameObj := types.Game{
				ID:   gameRequest.ID,
				Name: gameRequest.Name,
				Summary: gameRequest.Summary,
				CoverURL: gameRequest.CoverURL,
				FirstReleaseDate: gameRequest.ReleaseDate,
			}

			// Use service to add game to library
			if err := service.AddGameToLibrary(r.Context(), userID, gameObj); err != nil {
				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					err,
					http.StatusInternalServerError,
				)
				return
			}

			// Format response to match frontend expectations
			response := struct {
				Success bool `json:"success"`
				Game struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
				} `json:"game"`
			}{
				Success: true,
				Game: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
				}{
					ID:   gameObj.ID,
					Name: gameObj.Name,
				},
			}

			httputils.RespondWithJSON(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				http.StatusOK,
				response,
			)
		case http.MethodGet:
			appCtx.Logger.Info("GET request received", map[string]any{
				"requestID": requestID,
			})

			// Get user's library from service
			games, err := service.GetLibraryItems(r.Context(), userID)
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

			// Format response to game objects
			gameResponses := make([]struct{
				ID    int64  `json:"id"`
				Name  string `json:"name"`
			}, len(games))

			for i, game := range games {
				gameResponses[i].ID = game.ID
				gameResponses[i].Name = game.Name
			}

			response := struct {
				Success bool `json:"success"`
				Games any `json:"games"`
			}{
				Success: true,
				Games: gameResponses,
			}

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
        "path": r.URL.Path,
    })

    // Extract ID parameter using Chi's URL param extraction
    idParam := chi.URLParam(r, "id")
    appCtx.Logger.Info("Delete request with ID param", map[string]any{
        "idParam": idParam,
    })

    if idParam == "" {
        httputils.RespondWithError(
            httputils.NewResponseWriterAdapter(w),
            appCtx.Logger,
            requestID,
            errors.New("missing game ID parameter"),
            http.StatusBadRequest,
        )
        return
    }

    gameID, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        httputils.RespondWithError(
            httputils.NewResponseWriterAdapter(w),
            appCtx.Logger,
            requestID,
            errors.New("invalid game ID"),
            http.StatusBadRequest,
        )
        return
    }

    // Use service to delete game from library
    err = service.DeleteGameFromLibrary(r.Context(), userID, gameID)
    if err != nil {
        statusCode := http.StatusInternalServerError

        // Check specifically for not found errors
        if errors.Is(err, ErrGameNotFound) ||
           strings.Contains(strings.ToLower(err.Error()), "not found") {
            statusCode = http.StatusNotFound
        }

        httputils.RespondWithError(
            httputils.NewResponseWriterAdapter(w),
            appCtx.Logger,
            requestID,
            err,
            statusCode,
        )
        return
    }

    // Successful deletion response
    response := struct {
        Success bool  `json:"success"`
        ID      int64 `json:"id"`
    }{
        Success: true,
        ID:      gameID,
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
		}
	}
}



// Helper methods
func handleDelete(
	appCtx *appcontext.AppContext,
	w http.ResponseWriter,
	r *http.Request,
	service LibraryService,
	userID string,
	requestID string,
) {
	appCtx.Logger.Info("DELETE request received", map[string]any{
		"requestID": requestID,
		"path": r.URL.Path,
	})

	// Extract ID parameter from URL
	idParam := chi.URLParam(r, "id")
	appCtx.Logger.Debug("Delete request with ID param", map[string]any{
			"idParam": idParam,
	})

	if idParam == "" {
			httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					errors.New("missing game ID parameter"),
					http.StatusBadRequest,
			)
			return
	}

	gameID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
			httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					errors.New("invalid game ID"),
					http.StatusBadRequest,
			)
			return
	}

	// Delete game from library
	err = service.DeleteGameFromLibrary(r.Context(), userID, gameID)
	if err != nil {
			statusCode := http.StatusInternalServerError

			// Check specifically for "not found" errors
			if errors.Is(err, ErrGameNotFound) ||
				strings.Contains(strings.ToLower(err.Error()), "not found") {
					statusCode = http.StatusNotFound
			}

			httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					err,
					statusCode,
			)
			return
	}

	// Successful deletion response
	response := struct {
			Success bool  `json:"success"`
			ID      int64 `json:"id"`
	}{
			Success: true,
			ID:      gameID,
	}

	httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
	)
}