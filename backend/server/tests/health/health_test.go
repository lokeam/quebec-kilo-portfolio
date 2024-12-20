package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/shared/logger"
	"github.com/lokeam/qko-beta/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
  healthPath = "/api/v1/health"
  apiVersion = "1.0.0"
  testEnv = "test"
  defaultTestPort = 8080
)

/* Helper fns for setup, cleanup, custom assertions */

// Shared server creation logic
func createServer(tb testing.TB) *server.Server {
  tb.Helper()

  // Set up test logger
  log, err := logger.NewLogger(
    logger.WithEnv(logger.EnvTest),
    logger.WithAlertLevel(logger.LevelInfo),
  )
  require.NoError(tb, err)

  // Set up test config
  cfg := &config.Config{
    Server: config.ServerConfig{
      Port: defaultTestPort,
      Host: "localhost",
    },
    Env: testEnv,
    Debug: false,
  }

  srv := server.NewServer(cfg, log)

  // Only cleanup once per test, not per subtest
  if tb.Name() == "TestHealthCheckRoute" || tb.Name() == "TestHealthCheckRoute_EdgeCases" {
    tb.Cleanup(func() {
        if err := log.Cleanup(); err != nil {
            // Just log the error, don't fail the test
            tb.Logf("logger cleanup warning: %v", err)
        }
    })
}

  return srv
}

// Helper to create test server
func newTestServer(t *testing.T) *server.Server {
  t.Helper()
  return createServer(t)
}

// Helper to create benchmark server
func newBenchmarkServer(b *testing.B) *server.Server {
  b.Helper()
  return createServer(b)
}


// Test Health Check Route
func TestHealthCheckRoute(test *testing.T) {
  testCases := []struct {
    name                 string
    method               string
    path                 string
    wantHTTPStatusCode   int
    wantResponseBody     map[string]string
  }{
    {
      name: "GET /health returns 200 and the correct response",
      method: http.MethodGet,
      path: healthPath,
      wantHTTPStatusCode: http.StatusOK,
      wantResponseBody: map[string]string{
        "status": "available",
        "env": testEnv,
        "version": apiVersion,
      },
    },
  }

  // Test runner loop
  for _, tc := range testCases {
    test.Run(tc.name, func(t *testing.T) {
      // GIVEN
      srv := newTestServer(t)

      // Create test request
      // Mocks request
      request := httptest.NewRequest(tc.method, tc.path, nil)

      // Replacement for response writer, processes and compare the HTTP response with the expected output
      recorder := httptest.NewRecorder()

      // WHEN
      srv.Routes().ServeHTTP(recorder, request)

      // THEN
      assert.Equal(t,
        tc.wantHTTPStatusCode,
        recorder.Code,
        "Status code mismatch for %s request made to %s\nWant: %d\nGot: %d",
        tc.method,
        tc.path,
        tc.wantHTTPStatusCode,
        recorder.Code,
      )

      var actualResponse map[string]string
      err := json.NewDecoder(recorder.Body).Decode(&actualResponse)
      require.NoError(t, err)

      assert.Equal(t,
        tc.wantResponseBody,
        actualResponse,
        "Response body mismatch for %s request made to %s\nWant: %v\nGot: %v",
        tc.method,
        tc.path,
        tc.wantResponseBody,
        actualResponse,
      )
    })
  }
}

// Benchmark tests
func BenchmarkRoutes(benchmark *testing.B) {
  benchmarks := []struct {
    name     string
    method   string
    path     string
  }{
    {
      name: "Health Check",
      method: http.MethodGet,
      path: healthPath,
    },
  }

  // Test runner loop
  for _, bm := range benchmarks {
    benchmark.Run(bm.name, func(b *testing.B) {
      srv := newBenchmarkServer(b)
      req := httptest.NewRequest(bm.method, bm.path, nil)
      rec := httptest.NewRecorder()

      b.ResetTimer() // Reset timer after setup
      b.ReportAllocs() // Report memory allocations

      b.ResetTimer()
      for index := 0; index < b.N; index++ {
        srv.Routes().ServeHTTP(rec, req)
      }
    })
  }
}

