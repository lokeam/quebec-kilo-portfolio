package digital

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

type DigitalLocationRequest struct {
	ID             string     `json:"id" db:"id"`
	Name           string     `json:"name" db:"id"`
	LocationType   string     `json:"location_type" db:"user_id"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	URL            string     `json:"url" db:"url"`
	CreateAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

func NewDigitalLocationHandler(
	appCtx *appcontext.AppContext,
	digitalService *GameDigitalService,
) http.HandlerFunc {
	return func( w http.ResponseWriter, r *http.Request) {
		appCtx.Logger.Info("DigitalLocationHandler ServeHTTP called", map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
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
				handleGetDigitalLocation(w, r, appCtx, digitalService, userID, locationID, requestID)
			} else {
				handleListDigitalLocations(w, r, appCtx, digitalService, userID, requestID)
			}

		case http.MethodPost:
			handleCreateDigitalLocation(w, r, appCtx, digitalService, userID, requestID)

		case http.MethodPut:
		if locationID == "" {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("locationID is required"),
				http.StatusBadRequest,
			)
			return
		}
		handleUpdateDigitalLocation(w, r, appCtx, digitalService, userID, locationID, requestID)

		case http.MethodDelete:
			if locationID == "" {
				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					errors.New("locationID is required"),
					http.StatusBadRequest,
				)
				return
			}
			handleDeleteDigitalLocation(w, r, appCtx, digitalService, userID, locationID, requestID)

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
func handleListDigitalLocations(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service DigitalService,
	userID string,
	requestID string,
) {
	appCtx.Logger.Info("listing digital locations", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	locations, err := service.GetDigitalLocations(r.Context(), userID)
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
		Success     bool                        `json:"success"`
		Locations   []models.DigitalLocation    `json:"sublocations"`
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

func handleGetDigitalLocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service DigitalService,
	userID string,
	locationID string,
	requestID string,
) {
	appCtx.Logger.Info("Getting digital location", map[string]any{
		"requestID":   requestID,
		"userID":      userID,
		"locationID":  locationID,
	})

	digitalLocation, err := service.GetDigitalLocation(r.Context(), userID, locationID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrDigitalLocationNotFound) {
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
		Success    bool                     `json:"success"`
		Location   models.DigitalLocation   `json:"location"`
	} {
		Success:  true,
		Location: *digitalLocation,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusOK,
		response,
	)
}

func handleCreateDigitalLocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service DigitalService,
	userID string,
	requestID string,
) {
	appCtx.Logger.Info("Getting digital location", map[string]any{
		"requestID":    requestID,
		"userID":       userID,
	})

	var locationRequest DigitalLocationRequest
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

	digitalLocation := models.DigitalLocation{
		ID:             locationRequest.ID,
		Name:           locationRequest.Name,
		IsActive:       locationRequest.IsActive,
		URL:            locationRequest.URL,
		CreatedAt:      locationRequest.CreateAt,
		UpdatedAt:      locationRequest.UpdatedAt,
	}

	err := service.AddDigitalLocation(r.Context(), userID, digitalLocation)
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
		Success     bool                   `json:"success"`
		Location   models.DigitalLocation  `json:"location"`
	} {
		Success:  true,
		Location: digitalLocation,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusCreated,
		response,
	)
}

func handleUpdateDigitalLocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service DigitalService,
	userID string,
	locationID string,
	requestID string,
) {
	appCtx.Logger.Info("Updating digital location", map[string]any{
		"requestID":    requestID,
		"userID":       userID,
		"locationID":   locationID,
	})

	var locationRequest DigitalLocationRequest
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

	location := models.DigitalLocation{
		ID:            locationID,
		Name:          locationRequest.Name,
		IsActive:      locationRequest.IsActive,
		URL:           locationRequest.URL,
		CreatedAt:     locationRequest.CreateAt,
		UpdatedAt:     locationRequest.UpdatedAt,
	}

	err := service.UpdateDigitalLocation(r.Context(), userID, location)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrDigitalLocationNotFound) {
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
		Success    bool                    `json:"success"`
		Location   models.DigitalLocation  `json:"location"`
	} {
		Success:   true,
		Location:  location,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusOK,
		response,
	)
}

func handleDeleteDigitalLocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service DigitalService,
	userID string,
	locationID string,
	requestID string,
) {
	appCtx.Logger.Info("Deleting digital location", map[string]any{
		"requestID":    requestID,
		"userID":       userID,
		"locationID":   locationID,
	})

	err := service.DeleteDigitalLocation(r.Context(), userID, locationID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrDigitalLocationNotFound) {
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
