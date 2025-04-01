package physical

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

type PhysicalLocationRequest struct {
	Name           string `json:"name"`
	Label          string `json:"label"`
	LocationType   string `json:"location_type"`
	MapCoordinates string `json:"map_coordinates"`
}

func NewPhysicalLocationHandler(
  appCtx *appcontext.AppContext,
   physicalService *GamePhysicalService,
  ) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      appCtx.Logger.Info("PhysicalHandler ServeHTTP called", map[string]any{
        "method": r.Method,
        "path": r.URL.Path,
      })

      // Get Request ID for tracing
      requestID := httputils.GetRequestID(r)

      // Retrieve user ID from request context
      //domain := httputils.GetDomainFromRequest(r, "locations")

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

      // Extract location ID from URL if present
      var locationID string
      parts := strings.Split(r.URL.Path, "/")
      if len(parts) > 4 {
        locationID = parts[len(parts)-1]
      }

    // Handle different HTTP methods
    switch r.Method {
    case http.MethodGet:
      if locationID != "" {
        handleGetLocation(w, r, appCtx, physicalService, userID, locationID, requestID)
      } else {
        handleListLocations(w, r, appCtx, physicalService, userID, requestID)
      }

    case http.MethodPost:
      handleCreateLocation(w, r, appCtx, physicalService, userID, requestID)

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
      handleUpdateLocation(w, r, appCtx, physicalService, userID, locationID, requestID)

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
      handleDeleteLocation(w, r, appCtx, physicalService, userID, locationID, requestID)

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
// handles GET request to list all physical locations
func handleListLocations(
  w http.ResponseWriter,
  r *http.Request,
  appCtx *appcontext.AppContext,
  service *GamePhysicalService,
  userID string,
  requestID string,
) {
  appCtx.Logger.Info("listing physical locations", map[string]any{
    "requestID": requestID,
    "userID": userID,
  })

  locations, err := service.GetPhysicalLocations(r.Context(), userID)
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
    Success      bool               `json:"success"`
    Locations    []models.PhysicalLocation  `json:"locations"`
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
func handleGetLocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service PhysicalService,
	userID string,
	locationID string,
	requestID string,
) {
	appCtx.Logger.Info("Getting physical location", map[string]any{
		"requestID":  requestID,
		"userID":     userID,
		"locationID": locationID,
	})

	location, err := service.GetPhysicalLocations(r.Context(), userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrLocationNotFound) {
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
		Success  bool                    `json:"success"`
		Location []models.PhysicalLocation `json:"location"`
	}{
		Success:  true,
		Location: location,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusOK,
		response,
	)
}
func handleCreateLocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service PhysicalService,
	userID string,
	requestID string,
) {
	appCtx.Logger.Info("Creating physical location", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	var locationRequest PhysicalLocationRequest
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

	location := models.PhysicalLocation{
		Name:           locationRequest.Name,
		Label:          locationRequest.Label,
		LocationType:   locationRequest.LocationType,
		MapCoordinates: locationRequest.MapCoordinates,
		UserID:         userID,
	}

	err := service.AddPhysicalLocation(r.Context(), userID, location)
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
		Success  bool                    `json:"success"`
		Location models.PhysicalLocation `json:"location"`
	}{
		Success:  true,
		Location: location,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusCreated,
		response,
	)
}
func handleUpdateLocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service PhysicalService,
	userID string,
	locationID string,
	requestID string,
) {
	appCtx.Logger.Info("Updating physical location", map[string]any{
		"requestID":  requestID,
		"userID":     userID,
		"locationID": locationID,
	})

	var locationRequest PhysicalLocationRequest
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

	location := models.PhysicalLocation{
		ID:             locationID,
		Name:           locationRequest.Name,
		Label:          locationRequest.Label,
		LocationType:   locationRequest.LocationType,
		MapCoordinates: locationRequest.MapCoordinates,
		UserID:         userID,
	}

	err := service.UpdatePhysicalLocation(r.Context(), userID, location)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrLocationNotFound) {
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
		Success  bool                    `json:"success"`
		Location models.PhysicalLocation `json:"location"`
	}{
		Success:  true,
		Location: location,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusOK,
		response,
	)
}
func handleDeleteLocation(
	w http.ResponseWriter,
	r *http.Request,
	appCtx *appcontext.AppContext,
	service PhysicalService,
	userID string,
	locationID string,
	requestID string,
) {
	appCtx.Logger.Info("Deleting physical location", map[string]any{
		"requestID":  requestID,
		"userID":     userID,
		"locationID": locationID,
	})

	err := service.DeletePhysicalLocation(r.Context(), userID, locationID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrLocationNotFound) {
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
		Success bool   `json:"success"`
		ID      string `json:"id"`
	}{
		Success: true,
		ID:      locationID,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		appCtx.Logger,
		http.StatusOK,
		response,
	)
}
