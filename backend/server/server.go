package server

import (
	"encoding/json"
	"net/http"

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
func (s *Server) ServerHTTP(w http.ResponseWriter, r *http.Request) {
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