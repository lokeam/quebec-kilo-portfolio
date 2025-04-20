package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/app"
	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
)

type Server struct {
	Config      *config.Config
	AppContext  *appcontext.AppContext
	Logger      interfaces.Logger
	Router      chi.Router
}

func NewServer(
	config *config.Config,
	logger interfaces.Logger,
	appContext *appcontext.AppContext,
) *Server {

	s := &Server{
		Config:      config,
		Logger:      logger,
		AppContext:  appContext,
		Router:      chi.NewRouter(),
	}

	// NOTE: may need to use mock services for tests
	var services any

	// Initialize services
	appServices, err := app.NewServices(appContext)
	if err != nil {
		logger.Error("Failed to initialize services", map[string]any{
			"error": err.Error(),
		})
		if config.Env == "production" {
			panic(err)
		}

		// In test environment, attempt to useuse mocks
		if config.Env == "test" {
			// Create mock services for testing
			mockServices := mocks.NewMockServices()
			services = mockServices
			logger.Info("Using mock services for testing", nil)
		} else {
			// Return a server with minimal functionality for non-prod environments
			s.Router = chi.NewRouter()
			s.Router.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"status":"service initialization failed"}`))
			})
			return s
		}
	} else {
		services = appServices
	}

	s.Router = s.SetupRoutes(appContext, services)
	return s
}

// Add http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// Gracefully shut down server
func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("Shutting down server", map[string]any{
		"time": time.Now().Format(time.RFC3339),
	})

	// Clean up tasks
	return nil
}

// Force an immediate shutdown
func (s *Server) Close() error {
	s.Logger.Info("Forcefully closing server", map[string]any{
		"time": time.Now().Format(time.RFC3339),
	})

	return nil
}
