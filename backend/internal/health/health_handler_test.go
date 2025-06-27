package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	"github.com/lokeam/qko-beta/server"
	"github.com/stretchr/testify/assert"
)

const (
  healthPath = "/api/v1/health"
)

func TestHealthCheckRoute(t *testing.T) {
  mockConfig := mocks.NewMockConfig()
  mockConfig.HealthStatus = "available"
  testLogger := testutils.NewTestLogger()
  mockAppContext := &appcontext.AppContext{
    Config: mockConfig,
    Logger: testLogger,
  }
  mockServices := mocks.NewMockServices()

  t.Run("GET /health returns 200", func(t *testing.T) {
		// Setup - create server without services to avoid real service initialization
		srv := &server.Server{
			Config:     mockConfig,
			Logger:     testLogger,
			AppContext: mockAppContext,
			Router:     chi.NewRouter(),
		}
		req := httptest.NewRequest(http.MethodGet, healthPath, nil)
		rec := httptest.NewRecorder()

		// Execute - setup routes with mock services
		srv.SetupRoutes(mockAppContext, mockServices).ServeHTTP(rec, req)

		// Verify
		assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200")
	})

  t.Run("GET /health returns 500 when service is unavailable", func(t *testing.T) {
    // Setup
    mockConfig := mocks.NewMockConfig()
    mockConfig.HealthStatus = "unavailable"
    mockAppContext := &appcontext.AppContext{
        Config: mockConfig,
        Logger: testutils.NewTestLogger(),
    }

    // Create server without services to avoid real service initialization
    srv := &server.Server{
        Config:     mockConfig,
        Logger:     testLogger,
        AppContext: mockAppContext,
        Router:     chi.NewRouter(),
    }
    req := httptest.NewRequest(http.MethodGet, healthPath, nil)
    rec := httptest.NewRecorder()

    // Execute - setup routes with mock services
    srv.SetupRoutes(mockAppContext, mockServices).ServeHTTP(rec, req)

    // Verify
    assert.Equal(t, http.StatusInternalServerError, rec.Code, "expected status code 500")
})

}
