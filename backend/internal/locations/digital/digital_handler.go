package digital

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

// DigitalLocationRequest represents the request payload for digital location operations
type DigitalLocationRequest struct {
	ID             string               `json:"id,omitempty" db:"id"`
	Name           string               `json:"name" db:"name"`
	ServiceType    string               `json:"service_type" db:"service_type"`
	IsActive       bool                 `json:"is_active" db:"is_active"`
	URL            string               `json:"url" db:"url"`
	CreatedAt      time.Time            `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty" db:"updated_at"`
	Subscription   *models.Subscription `json:"subscription,omitempty"`
}

// RegisterDigitalRoutes registers all digital location routes
func RegisterDigitalRoutes(r chi.Router, appCtx *appcontext.AppContext, service services.DigitalService) {
	// Base routes
	r.Get("/", GetUserDigitalLocations(appCtx, service))
	r.Post("/", AddDigitalLocation(appCtx, service))

	// Nested routes with ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", GetDigitalLocation(appCtx, service))
		r.Put("/", UpdateDigitalLocation(appCtx, service))
		r.Delete("/", RemoveDigitalLocation(appCtx, service))
	})
}

// GetUserDigitalLocations handles GET requests for listing all digital locations
func GetUserDigitalLocations(appCtx *appcontext.AppContext, service services.DigitalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

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

		appCtx.Logger.Info("Listing digital locations", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		locations, err := service.GetUserDigitalLocations(r.Context(), userID)
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
			Success   bool                     `json:"success"`
			UserID    string                   `json:"user_id"`
			Locations []models.DigitalLocation `json:"locations"`
		}{
			Success:   true,
			UserID:    userID,
			Locations: locations,
	}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
		http.StatusOK,
		response,
	)
	}
}

// GetDigitalLocation handles GET requests for a single digital location
func GetDigitalLocation(appCtx *appcontext.AppContext, service services.DigitalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

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

		locationID := chi.URLParam(r, "id")
		appCtx.Logger.Info("Getting digital location", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": locationID,
	})

	if locationID == "" {
		httputils.RespondWithError(
			httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
			errors.New("id is required"),
			http.StatusBadRequest,
		)
		return
	}

	var location models.DigitalLocation
	var err error

	if _, parseErr := uuid.Parse(locationID); parseErr == nil {
			location, err = service.GetDigitalLocation(r.Context(), userID, locationID)
	} else {
			location, err = service.FindDigitalLocationByName(r.Context(), userID, locationID)
	}

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
			Success  bool                   `json:"success"`
			Location models.DigitalLocation `json:"location"`
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
}

// AddDigitalLocation handles POST requests for creating a new digital location
func AddDigitalLocation(appCtx *appcontext.AppContext, service services.DigitalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

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

		appCtx.Logger.Info("Creating digital location", map[string]any{
			"requestID": requestID,
			"userID":    userID,
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

	now := time.Now()
	digitalLocation := models.DigitalLocation{
			ID:          uuid.New().String(),
			UserID:      userID,
			Name:        locationRequest.Name,
			ServiceType: locationRequest.ServiceType,
			IsActive:    locationRequest.IsActive,
			URL:         locationRequest.URL,
			CreatedAt:   now,
			UpdatedAt:   now,
	}

	// Handle subscription if provided
	if locationRequest.Subscription != nil {
		subscription := *locationRequest.Subscription
		subscription.LocationID = digitalLocation.ID
		subscription.CreatedAt = now
		subscription.UpdatedAt = now
		digitalLocation.Subscription = &subscription
	}

		createdLocation, err := service.AddDigitalLocation(r.Context(), userID, digitalLocation)
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
			Success  bool                   `json:"success"`
			Location models.DigitalLocation `json:"location"`
		}{
			Success:  true,
			Location: createdLocation,
		}

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
		http.StatusCreated,
			response,
	)
	}
}

// UpdateDigitalLocation handles PUT requests for updating a digital location
func UpdateDigitalLocation(appCtx *appcontext.AppContext, service services.DigitalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

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

	locationID := chi.URLParam(r, "id")
		appCtx.Logger.Info("Updating digital location", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": locationID,
		})

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

	var updateReq DigitalLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			appCtx.Logger.Error("Failed to decode request body", map[string]any{"error": err})
		httputils.RespondWithError(
			httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
			errors.New("invalid request body"),
			http.StatusBadRequest,
		)
		return
	}

	// Get existing location
		location, err := service.GetDigitalLocation(r.Context(), userID, locationID)
	if err != nil {
			appCtx.Logger.Error("Failed to get existing location", map[string]any{"error": err})
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

	// Update fields
	location.Name = updateReq.Name
	location.ServiceType = updateReq.ServiceType
	location.IsActive = updateReq.IsActive
	location.URL = updateReq.URL
	location.UpdatedAt = time.Now()

	// Ensure the ID is set from the URL
	location.ID = locationID

		// If subscription is provided, update it
	if updateReq.Subscription != nil {
		subscription := &models.Subscription{
			ID:              updateReq.Subscription.ID,
			LocationID:      location.ID,
				BillingCycle:    updateReq.Subscription.BillingCycle,
				CostPerCycle:    updateReq.Subscription.CostPerCycle,
				NextPaymentDate: updateReq.Subscription.NextPaymentDate,
				PaymentMethod:   updateReq.Subscription.PaymentMethod,
			UpdatedAt:       time.Now(),
		}

		// Only set CreatedAt if this is a new subscription
		if location.Subscription == nil {
			subscription.CreatedAt = time.Now()
		} else {
			subscription.CreatedAt = location.Subscription.CreatedAt
		}

		location.Subscription = subscription
	}

		err = service.UpdateDigitalLocation(r.Context(), userID, location)
	if err != nil {
			appCtx.Logger.Error("Failed to update location", map[string]any{"error": err})
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
			Success  bool                   `json:"success"`
			Location models.DigitalLocation `json:"location"`
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
}

// RemoveDigitalLocation handles DELETE requests for removing a digital location
func RemoveDigitalLocation(appCtx *appcontext.AppContext, service services.DigitalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

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

	locationID := chi.URLParam(r, "id")
		appCtx.Logger.Info("Deleting digital location", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": locationID,
		})

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

		err := service.RemoveDigitalLocation(r.Context(), userID, locationID)
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
}

// For backward compatibility - wrapper for the digital service catalog
func GetDigitalServicesCatalog(appCtx *appcontext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create catalog data
		catalog := []DigitalServiceCatalogItem{
			{
				ID:                    "steam",
				Name:                  "Steam",
				Logo:                  "steam_logo.png",
				IsSubscriptionService: false,
				URL:                   "https://store.steampowered.com/",
			},
			{
				ID:                    "epic",
				Name:                  "Epic Games Store",
				Logo:                  "epic_logo.png",
				IsSubscriptionService: false,
				URL:                   "https://www.epicgames.com/store/",
			},
			{
				ID:                    "gog",
				Name:                  "GOG.com",
				Logo:                  "gog_logo.png",
				IsSubscriptionService: false,
				URL:                   "https://www.gog.com/",
			},
			{
				ID:                    "xbox",
				Name:                  "Xbox Game Pass",
				Logo:                  "xbox_logo.png",
				IsSubscriptionService: true,
				URL:                   "https://www.xbox.com/xbox-game-pass",
			},
			{
				ID:                    "playstation",
				Name:                  "PlayStation Plus",
				Logo:                  "playstation_logo.png",
				IsSubscriptionService: true,
				URL:                   "https://www.playstation.com/playstation-plus/",
			},
			{
				ID:                    "ea_play",
				Name:                  "EA Play",
				Logo:                  "ea_play_logo.png",
				IsSubscriptionService: true,
				URL:                   "https://www.ea.com/ea-play",
			},
			{
				ID:                    "ubisoft_plus",
				Name:                  "Ubisoft+",
				Logo:                  "ubisoft_plus_logo.png",
				IsSubscriptionService: true,
				URL:                   "https://store.ubi.com/ubisoft-plus/",
			},
			{
				ID:                    "nintendo_switch_online",
				Name:                  "Nintendo Switch Online",
				Logo:                  "nintendo_switch_online_logo.png",
				IsSubscriptionService: true,
				URL:                   "https://www.nintendo.com/switch/online-service/",
			},
		}

		response := struct {
			Success bool                       `json:"success"`
			Catalog []DigitalServiceCatalogItem `json:"catalog"`
		}{
			Success: true,
			Catalog: catalog,
		}

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}
