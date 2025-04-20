package server

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lokeam/qko-beta/app"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/health"
	"github.com/lokeam/qko-beta/internal/library"
	"github.com/lokeam/qko-beta/internal/locations/digital"
	"github.com/lokeam/qko-beta/internal/locations/physical"
	"github.com/lokeam/qko-beta/internal/locations/sublocation"
	"github.com/lokeam/qko-beta/internal/search"
	"github.com/lokeam/qko-beta/internal/shared/logger"
	customMiddleware "github.com/lokeam/qko-beta/internal/shared/middleware"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	authMiddleware "github.com/lokeam/qko-beta/server/middleware"
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
			LibraryMap:    mockSvc.LibraryMap,
			Wishlist:      mockSvc.Wishlist,
			SearchFactory: mockSvc.SearchFactory,
			SearchMap:     mockSvc.SearchMap,
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

	// Add mock authentication middleware
	mux.Use(authMiddleware.MockAuthMiddleware(appContext))

	// Create library handler
	libraryHandler := library.NewLibraryHandler(appContext, svc.LibraryMap)

	// Create physical handler
	physicalHandler := physical.NewPhysicalLocationHandler(appContext, svc.Physical)

	// Create sublocation handler
	sublocationHandler := sublocation.NewSublocationHandler(appContext, svc.Sublocation)

	// Create digital services catalog handler
	digitalServicesCatalogHandler := digital.GetDigitalServicesCatalog(appContext)

	// Initialize handlers using single App Context
	healthHandler := health.NewHealthHandler(s.Config, s.Logger)

	// Create search handler
	searchHandler := search.NewSearchHandler(
		appContext,
		svc.SearchMap,
		svc.Library,
		svc.Wishlist,
	)

	// Debug logging for searchHandler
	if searchHandler == nil {
		appContext.Logger.Error("searchHandler is nil", map[string]any{
			"appContext": appContext,
		})
	} else {
		appContext.Logger.Info("searchHandler initialized successfully", map[string]any{
			"appContext": appContext,
		})
	}

	if reflect.TypeOf(searchHandler) != reflect.TypeOf(http.HandlerFunc(nil)) {
		appContext.Logger.Error("searchHandler is not of type http.HandlerFunc", map[string]any{
			"appContext": appContext,
		})
	} else {
		appContext.Logger.Info("searchHandler is of type http.HandlerFunc", map[string]any{
			"appContext": appContext,
		})
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

		// Search
		r.Post("/search", searchHandler)

		// Library
		// TODO: Refactor this to use the physical location pattern below
		r.Get("/library", libraryHandler)
		r.Post("/library", libraryHandler)
		r.Route("/library/games", func(r chi.Router) {
			r.Delete("/{gameID}", libraryHandler)
		})

		// Physical Locations
		r.Route("/locations/physical", func(r chi.Router) {
			appContext.Logger.Info("Registering physical location routes", map[string]any{
        "path": "/api/v1/locations/physical",
			})

			r.Get("/", physicalHandler)
			r.Post("/", physicalHandler)
			r.Get("/{id}", physicalHandler)
			r.Put("/{id}", physicalHandler)
			r.Delete("/{id}", physicalHandler)
		})

		// Sublocations
		r.Route("/locations/sublocations", func(r chi.Router) {
			r.Get("/", sublocationHandler)
			r.Post("/", sublocationHandler)
			r.Get("/{id}", sublocationHandler)
			r.Put("/{id}", sublocationHandler)
			r.Delete("/{id}", sublocationHandler)
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

		appContext.Logger.Info("Routes registered", map[string]any{
			"health":         "/api/v1/health",
			"search":         "/api/v1/search",
			"library":        "/api/v1/library",
			"physical":       "/api/v1/locations/physical",
			"sublocations":   "/api/v1/locations/sublocations",
			"digital":        "/api/v1/locations/digital",
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