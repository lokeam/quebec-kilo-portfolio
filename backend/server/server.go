package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/shared/logger"
)

type Server struct {
	config *config.Config
	logger *logger.Logger
}

func NewServer(
	config *config.Config,
	logger *logger.Logger,
) *Server {

	return &Server{
		config: config,
		logger: logger,
	}
}

// Add http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Use router to handle requests
	s.Routes().ServeHTTP(w, r)
}

// Utility - write JSON
func (s *Server) writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers http.Header,
) error {
	w.Header().Set("Content-Type", "application/json")

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Gracefully shut down server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server", map[string]any{
		"time": time.Now().Format(time.RFC3339),
	})

	// Clean up tasts
	return nil
}

// Force an immediate shutdown
func (s *Server) Close() error {
	s.logger.Info("Forcefully closing server", map[string]any{
		"time": time.Now().Format(time.RFC3339),
	})

	return nil
}