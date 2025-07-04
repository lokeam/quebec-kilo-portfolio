package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lokeam/qko-beta/app"
	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/dashboard"
	"github.com/lokeam/qko-beta/internal/health"
	"github.com/lokeam/qko-beta/internal/library"
	"github.com/lokeam/qko-beta/internal/locations/digital"
	"github.com/lokeam/qko-beta/internal/locations/physical"
	"github.com/lokeam/qko-beta/internal/locations/sublocation"
	"github.com/lokeam/qko-beta/internal/search"
	"github.com/lokeam/qko-beta/internal/shared/logger"
	customMiddleware "github.com/lokeam/qko-beta/internal/shared/middleware"
	"github.com/lokeam/qko-beta/internal/spend_tracking"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) SetupRoutes(appContext *appcontext.AppContext, services interface{}) chi.Router {
	// Type assertion to get services struct
	svc, ok := services.(*app.Services)
	if !ok {
		// Try to assert as mock services
		mockSvc, ok := services.(*mocks.MockServices)
		if !ok {
			appContext.Logger.Error("Invalid services type provided", map[string]any{
				"type": fmt.Sprintf("%T", services),
			})
			// Create a minimal router
			mux := chi.NewRouter()
			mux.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"status":"service initialization failed"}`))
			})
			return mux
		}
		// Convert mock services to app services interface implementation
		svc = &app.Services{
			Digital:       mockSvc.Digital,
			Physical:      mockSvc.Physical,
			Sublocation:   mockSvc.Sublocation,
			Library:       mockSvc.Library,
			Wishlist:      mockSvc.Wishlist,
			Search:        mockSvc.Search,
			SpendTracking: mockSvc.SpendTracking,
			Dashboard:     mockSvc.Dashboard,
		}
		}

	// Create Router
	mux := chi.NewRouter()

	// Setup middleware
	s.setupMiddleware(mux)
	s.setupCORS(mux)

	// Add trailing slash middleware
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log the original path
			appContext.Logger.Debug("Before StripSlashes", map[string]interface{}{
				"method": r.Method,
				"path":   r.URL.Path,
			})

			// Call the StripSlashes middleware
			middleware.StripSlashes(next).ServeHTTP(w, r)

			// Log the path after StripSlashes
			appContext.Logger.Debug("After StripSlashes", map[string]interface{}{
				"method": r.Method,
				"path":   r.URL.Path,
			})
		})
	})

	// Create digital services catalog handler
	digitalServicesCatalogHandler := digital.GetDigitalServicesCatalog(appContext)

	// Initialize handlers using single App Context
	healthHandler := health.NewHealthHandler(s.Config, s.Logger)

	// Create search handler
	gameSearchService, err := search.NewGameSearchService(appContext)
	if err != nil {
			appContext.Logger.Error("Failed to create game search service", map[string]any{
					"error": err,
			})
			// Handle error appropriately
	}

	// Log route setup
	log, err := logger.NewLogger()
	if err == nil {
		log.Debug("Setting up routes", map[string]interface{}{
			"base_path": "/api/v1",
		})
	}

	// Initialize Routes
	mux.Route("/api/v1", func(r chi.Router) {
		// Health
		r.Get("/health", healthHandler)

		// Protect routes with Auth0
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.EnsureValidToken())

			// Search
			r.Route("/search", func(r chi.Router) {
				appContext.Logger.Info("Registering search routes", map[string]any{
					"path": "/api/v1/search",
				})

				search.RegisterSearchRoutes(
					r,
					appContext,
					gameSearchService,
					svc.Library,
					svc.Wishlist,
				)
			})

			// Library
			r.Route("/library", func(r chi.Router) {
				appContext.Logger.Info("Registering library routes", map[string]any{
					"path": "/api/v1/library",
				})

				// Register routes using the new pattern
				library.RegisterLibraryRoutes(r, appContext, svc.Library, svc.Analytics)
			})

			// Physical Locations
			r.Route("/locations/physical", func(r chi.Router) {
				appContext.Logger.Info("Registering physical location routes", map[string]any{
					"path": "/api/v1/locations/physical",
				})

				// Register routes using the new pattern
				physical.RegisterPhysicalRoutes(r, appContext, svc.Physical, svc.Analytics)
			})

			// Sublocations
			r.Route("/locations/sublocations", func(r chi.Router) {
				appContext.Logger.Info("Registering sublocation routes", map[string]any{
					"path": "/api/v1/locations/sublocations",
				})

				// Register routes using the new pattern
				sublocation.RegisterSublocationRoutes(r, appContext, svc.Sublocation, svc.Analytics)
			})

			// Digital Locations
			r.Route("/locations/digital", func(r chi.Router) {
				if err == nil {
					log.Debug("Setting up digital locations routes", map[string]interface{}{
						"path": "/api/v1/locations/digital",
					})
				}

				// Register routes using the new pattern
				digital.RegisterDigitalRoutes(r, appContext, svc.Digital)

				// Services Catalog
				r.Get("/services/catalog", digitalServicesCatalogHandler)
			})

			// Spend Tracking
			r.Route("/spend-tracking", func(r chi.Router) {
				appContext.Logger.Info("Registering spend tracking routes", map[string]any{
					"path": "/api/v1/spend-tracking",
				})

				// Register routes using new pattern
				spend_tracking.RegisterSpendTrackingRoutes(r, appContext, svc.SpendTracking)
			})

			// Dashboard
			r.Route("/dashboard", func(r chi.Router) {
				appContext.Logger.Info("Registering dashboard routes", map[string]any{
					"path": "/api/v1/dashboard",
				})

				dashboard.RegisterDashboardRoutes(r, appContext, svc.Dashboard)
			})

			// Analytics
			r.Route("/analytics", func(r chi.Router) {
				appContext.Logger.Info("Registering analytics routes", map[string]any{
					"path": "/api/v1/analytics",
				})

				// Register routes using the service from the Services struct
				analytics.RegisterRoutes(r, appContext, svc.Analytics)
			})

			appContext.Logger.Info("Routes registered", map[string]any{
				"health":         "/api/v1/health",
				"search":         "/api/v1/search",
				"library":        "/api/v1/library",
				"physical":       "/api/v1/locations/physical",
				"sublocations":   "/api/v1/locations/sublocations",
				"digital":        "/api/v1/locations/digital",
				"spend-tracking": "/api/v1/spend-tracking",
				"dashboard":      "/api/v1/dashboard",
				"analytics":      "/api/v1/analytics",
			})
		})
	})

	// Add metrics endpoint with basic security
	if appContext.Config.Env == "production" {
		// In production, require API key
		mux.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				apiKey := r.Header.Get("X-API-Key")
				validAPIKey := os.Getenv("METRICS_API_KEY")

				if apiKey != validAPIKey {
						http.Error(w, "Forbidden", http.StatusForbidden)
						return
				}

				promhttp.Handler().ServeHTTP(w, r)
		}))
	} else {
		// In development, allow open access
		mux.Handle("/metrics", promhttp.Handler())
	}
	appContext.Logger.Info("Metrics endpoint registered", map[string]any{
    "metrics": "/metrics",
	})


	return mux
}

func (s *Server) setupMiddleware(mux *chi.Mux) {
	mux.Use(middleware.RequestID)  // Trace requests
	mux.Use(middleware.RealIP)     // Get actual client IP

	// Custom request ID middleware
	mux.Use(customMiddleware.EnrichRequestContext)

	// Sentry Middleware
	mux.Use(customMiddleware.SentryMiddleware)

	// Prometheus Middleware
	mux.Use(customMiddleware.MetricsMiddleware)

	mux.Use(middleware.Logger)     // Request logging
	mux.Use(middleware.Recoverer)  // Panic recovery
	mux.Use(middleware.Timeout(60 * time.Second))
}

func (s *Server) setupCORS(mux *chi.Mux) {
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.AppContext.Config.CORS.AllowedOrigins,
		AllowedMethods:   s.AppContext.Config.CORS.AllowedMethods,
		AllowedHeaders:   s.AppContext.Config.CORS.AllowedHeaders,
		ExposedHeaders:   s.AppContext.Config.CORS.ExposedHeaders,
		AllowCredentials: s.AppContext.Config.CORS.AllowCredentials,
		MaxAge:           s.AppContext.Config.CORS.MaxAge,
	}))
}