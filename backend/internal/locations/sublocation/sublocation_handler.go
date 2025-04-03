package sublocation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

type SublocationRequest struct {
	Name            string `json:"name"`
	LocationType    string `json:"location_type"`
	BgColor         string `json:"bg_color"`
	Capacity        int    `json:"capacity"`
}

func NewSublocationHandler(
	appCtx *appcontext.AppContext,
	sublocationService *GameSublocationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appCtx.Logger.Info("SublocationHandler ServeHTTP called", map[string]any{
			"method": r.Method,
			"path": r.URL.Path,
		})

		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

		// Retrieve user ID from request context

		// Extract location ID from URL
		userID := httputils.GetUserID(r)
		if userID == "" {
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

		var locationID string
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) > 4 {
			locationID = parts[len(parts) - 1]
		}

		// Handle different HTTP Methods
		switch r.Method {
		case http.MethodGet:
			if locationID != "" {
				handleGetSublocation(w, r, appCtx, sublocationService, userID, locationID, requestID)
			} else {
				handleListSublocations(w, r, appCtx, sublocationService, userID, requestID)
			}

		case http.MethodPost:
			handleCreateSublocation(w, r, appCtx, sublocationService, userID, requestID)

		case http.MethodPut:
			if locationID == "" {
				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					errors.New("location ID is required"),
					http.StatusBadRequest,
				)
				return
			}
			handleUpdateSublocation(w, r, appCtx, sublocationService, userID, locationID, requestID)

		case http.MethodDelete:
			if locationID == "" {
				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					errors.New("location ID is required"),
					http.StatusBadRequest,
				)
				return
			}
			handleDeleteSublocation(w, r, appCtx, sublocationService, userID, locationID, requestID)

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

// Helper fns
func handleListSublocations(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service *GameSublocationService,
	userID string,
	requestID string,
) {
	appCtx.Logger.Info("listing physical locations", map[string]any{
		"requestID": requestID,
		"userID": userID,
	})

	locations, err := service.GetSublocations(r.Context(),userID)
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

	response := struct {
		Success        bool                    `json:"success"`
		Locations      []models.Sublocation    `json:"sublocations"`
	} {
		Success: true,
		Locations: locations,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusOK,
		response,
	)
}

func handleGetSublocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service SublocationService,
	userID string,
	locationID string,
	requestID string,
) {
	appCtx.Logger.Info("Getting sublocation", map[string]any{
		"requestID":  requestID,
		"userID":     userID,
		"locationID": locationID,
	})

	sublocation, err := service.GetSublocation(r.Context(), userID, locationID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrSublocationNotFound) {
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

	response := struct {
		Success    bool                 `json:"success"`
		Location   models.Sublocation   `json:"location"`
	} {
		Success:  true,
		Location: sublocation,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusOK,
		response,
	)
}

func handleCreateSublocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service SublocationService,
	userID string,
	requestID string,
) {
	appCtx.Logger.Info("Getting sublocation", map[string]any{
		"requestID":  requestID,
		"userID":     userID,
	})

	var locationRequest SublocationRequest
	if err := json.NewDecoder(r.Body).Decode(&locationRequest); err != nil {
		httputils.RespondWithError(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			requestID,
			errors.New("invalid request body"),
			http.StatusBadRequest,
		)
		return
	}

	sublocation := models.Sublocation{
		Name:             locationRequest.Name,
		LocationType:     locationRequest.LocationType,
		BgColor:          locationRequest.BgColor,
		Capacity:         locationRequest.Capacity,
		UserID:           userID,
	}

	err := service.AddSublocation(r.Context(), userID, sublocation)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrValidationFailed) {
			statusCode = http.StatusBadRequest
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

	response := struct {
		Success     bool                  `json:"success"`
		Location    models.Sublocation    `json:"location"`
	} {
		Success:    true,
		Location:   sublocation,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusCreated,
		response,
	)
}

func handleUpdateSublocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service SublocationService,
	userID string,
	locationID string,
	requestID string,
) {
	appCtx.Logger.Info("Updating physical location", map[string]any{
		"requestID":    requestID,
		"userID":       userID,
		"locationID":   locationID,
	})

	var locationRequest SublocationRequest
	if err := json.NewDecoder(r.Body).Decode(&locationRequest); err != nil {
		httputils.RespondWithError(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			requestID,
			errors.New("invalid request body"),
			http.StatusBadRequest,
		)
		return
	}

	location := models.Sublocation{
		ID:             locationID,
		Name:           locationRequest.Name,
		LocationType:   locationRequest.LocationType,
		BgColor:        locationRequest.BgColor,
		Capacity:       locationRequest.Capacity,
		UserID:         userID,
	}

	err := service.UpdateSublocation(r.Context(), userID, location)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrSublocationNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, ErrValidationFailed) {
			statusCode = http.StatusBadRequest
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

	response := struct {
		Success     bool                `json:"success"`
		Location    models.Sublocation  `json:"location"`
	} {
		Success:    true,
		Location:   location,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusOK,
		response,
	)
}

func handleDeleteSublocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service SublocationService,
	userID string,
	locationID string,
	requestID string,
) {
	appCtx.Logger.Info("Deleting sublocation location", map[string]any{
		"requestID":    requestID,
		"userID":       userID,
		"locationID":   locationID,
	})

	err := service.DeleteSublocation(r.Context(), userID, locationID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrSublocationNotFound) {
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

	response := struct {
		Success    bool      `json:"success"`
		ID         string    `json:"id"`
	} {
		Success:   true,
		ID:        locationID,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusOK,
		response,
	)
}
