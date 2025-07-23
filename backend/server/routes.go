package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
	"github.com/lokeam/qko-beta/internal/users"
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

	// Create user service for middleware and routes
	appContext.Logger.Debug("Creating user service for middleware", nil)
	userService, err := users.NewUserService(appContext)
	if err != nil {
		appContext.Logger.Error("Failed to create user service for middleware", map[string]any{
			"error": err,
			"error_type": fmt.Sprintf("%T", err),
		})
		// Continue without user service - middleware will handle gracefully
		userService = nil
	} else {
		appContext.Logger.Info("User service created successfully for middleware", nil)
	}

	// Initialize Routes
	mux.Route("/api/v1", func(r chi.Router) {
		// Health
		r.Get("/health", healthHandler)

		// Protect routes with Auth0
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.EnsureValidToken())
			r.Use(customMiddleware.RequireUserExists(userService))

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
				appContext.Logger.Debug("Setting up digital locations routes", map[string]interface{}{
					"path": "/api/v1/locations/digital",
				})

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
				appContext.Logger.Debug("ENTERED dashboard route block", nil)
				appContext.Logger.Info("Registering dashboard routes", map[string]any{
					"path": "/api/v1/dashboard",
				})

				appContext.Logger.Debug("Before RegisterDashboardRoutes", nil)
				dashboard.RegisterDashboardRoutes(r, appContext, svc.Dashboard)
				appContext.Logger.Debug("After RegisterDashboardRoutes", nil)
			})

			appContext.Logger.Debug("After dashboard, before users", nil)
			appContext.Logger.Debug("About to register /users route", nil)
			// Users - Profile Management & Account Deletion
			r.Route("/users", func(r chi.Router) {
				appContext.Logger.Debug("ENTERED users route block", nil)
				appContext.Logger.Info("Registering user routes", map[string]any{
					"path": "/api/v1/users",
				})

				// Create user deletion service for account deletion
				appContext.Logger.Debug("Attempting to create user deletion service", nil)
				userDeletionService, err := users.NewUserDeletionService(appContext)
				if err != nil {
					appContext.Logger.Error("Failed to create user deletion service", map[string]any{
						"error": err,
						"error_type": fmt.Sprintf("%T", err),
					})
				} else {
					appContext.Logger.Info("User deletion service created successfully", nil)
				}

				// Register unified user routes (profile + deletion)
				if userService != nil && userDeletionService != nil {
					users.RegisterUserRoutes(r, appContext, userService, userDeletionService)
				}
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
				"users":          "/api/v1/users",
				"analytics":      "/api/v1/analytics",
			})
		})
	})

	// Add Sentry tunnel endpoint to avoid ad blockers
	mux.HandleFunc("/api/events", SentryTunnelHandler)

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
		// In "dev", allow open access
		mux.Handle("/metrics", promhttp.Handler())
	}
	appContext.Logger.Info("Metrics endpoint registered", map[string]any{
    "metrics": "/metrics",
	})


	return mux
}

// extractSentryKey extracts the key from a Sentry DSN
// DSN format: https://key@host/project
func extractSentryKey(dsn string) string {
	// Remove the https:// prefix
	if strings.HasPrefix(dsn, "https://") {
		dsn = dsn[8:] // Remove "https://"
	}

	// Find the @ symbol and extract the key part
	for i, char := range dsn {
		if char == '@' {
			return dsn[:i]
		}
	}
	return ""
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

// SentryTunnelHandler forwards Sentry envelopes from the frontend to Sentry's API.
func SentryTunnelHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}

	sentryDSN := os.Getenv("SENTRY_DSN_FRNT")
	if sentryDSN == "" {
			http.Error(w, "Sentry DSN not configured", http.StatusServiceUnavailable)
			return
	}

	// Parse DSN: "https://key@host/projectId"
	dsnNoProtocol := strings.TrimPrefix(sentryDSN, "https://")
	dsnParts := strings.SplitN(dsnNoProtocol, "/", 2)
	if len(dsnParts) != 2 {
			log.Printf("Invalid Sentry DSN format: %q", sentryDSN)
			http.Error(w, "Invalid Sentry DSN format", http.StatusInternalServerError)
			return
	}
	hostPart, projectId := dsnParts[0], dsnParts[1]
	keyAndHost := strings.SplitN(hostPart, "@", 2)
	if len(keyAndHost) != 2 {
			log.Printf("Invalid Sentry DSN host: %q", hostPart)
			http.Error(w, "Invalid Sentry DSN host", http.StatusInternalServerError)
			return
	}
	host := keyAndHost[1]

	// Build Sentry envelope URL
	sentryURL := "https://" + host + "/api/" + projectId + "/envelope/"

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", sentryURL, r.Body)
	if err != nil {
			log.Printf("Failed to create Sentry request: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	req.Header.Set("User-Agent", r.Header.Get("User-Agent"))

	// Forward to Sentry
	resp, err := client.Do(req)
	if err != nil {
			log.Printf("Failed to forward Sentry request: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
	}
	defer resp.Body.Close()

	// Return Sentry's response
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}