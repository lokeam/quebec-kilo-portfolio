package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
  testLogger := testutils.NewTestLogger()
  mockAppContext := &appcontext.AppContext{
    Config: mockConfig,
    Logger: testLogger,
  }

  t.Run(`GET /health returns 200 and the correct response`, func(t *testing.T) {
    testServer := server.NewServer(mockConfig, testLogger, mockAppContext)
    testRequest := httptest.NewRequest(http.MethodGet, healthPath, nil)
    testResponseRecorder := httptest.NewRecorder()

    mockServices := mocks.NewMockServices()
    testServer.SetupRoutes(mockAppContext, mockServices).ServeHTTP(testResponseRecorder, testRequest)

    assert.Equal(t, http.StatusOK, testResponseRecorder.Code, "expected status code 200")

    var response map[string]string
    err := json.NewDecoder(testResponseRecorder.Body).Decode(&response)
    assert.NoError(t, err, "health check test failed to decode resposne")

    expectedResponse := map[string]string{
      "status": "available",
    }

    assert.Equal(t, expectedResponse, response, "response body does not match expected response")
  })

  t.Run(`GET /health returns 500 when service is unavailable`, func(t *testing.T) {
    mockConfig := mocks.NewMockConfig()
    mockAppContext := &appcontext.AppContext{
      Config: mockConfig,
      Logger: testutils.NewTestLogger(),
    }

    testServer := server.NewServer(mockConfig, testLogger, mockAppContext)
    testRequest := httptest.NewRequest(http.MethodGet, healthPath, nil)
    testResponseRecorder := httptest.NewRecorder()

    mockServices := mocks.NewMockServices()
    testServer.SetupRoutes(mockAppContext, mockServices).ServeHTTP(testResponseRecorder, testRequest)

    assert.Equal(t, http.StatusInternalServerError, testResponseRecorder.Code, "expected status code 500")

    var response map[string]string
		err := json.NewDecoder(testResponseRecorder.Body).Decode(&response)
		assert.NoError(t, err, "failed to decode response")

		expectedResponse := map[string]string{
			"status":  "unavailable",
			"env":     mockConfig.Env,
		}
		assert.Equal(t, expectedResponse, response, "response body mismatch")
  })

}
