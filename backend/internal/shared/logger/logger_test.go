package logger

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Logger Creation
func TestLogger_Creation(test *testing.T) {
	testCases := []struct {
		name       string
		options    []Option
		wantEnv    Environment
		wantLevel  AlertLevel
	}{
		{
			name:      "Default settings create production logger",
			options:    nil,
			wantEnv:    EnvProd,
			wantLevel:  LevelInfo,
		},
		{
			name:     "Development mode creates development logger",
			options:  []Option{
				WithEnv(EnvDev),
			},
			wantEnv:   EnvDev,
			wantLevel: LevelInfo,
		},
		{
			name:     "Custom alert level is set correctly",
			options:  []Option{
				WithAlertLevel(LevelDebug),
			},
			wantEnv:   EnvProd,
			wantLevel: LevelDebug,
		},
		{
			name:     "Multiple options are applied correctly",
			options:  []Option{
				WithEnv(EnvDev),
				WithAlertLevel(LevelDebug),
			},
			wantEnv:   EnvDev,
			wantLevel: LevelDebug,
		},
	}

	// Test runner loop
	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			// GIVEN
			// Create RAM space for current test
			buffer := setupTestBuffer(test)

			// Ensure RAM is cleared after test
			test.Cleanup(func() {
				cleanupLogger(test, buffer)
			})

			// WHEN
			// Create logger with test options and RAM output
			options := append(tc.options, WithOutput(buffer))
			logger, err := NewLogger(options...)

			// THEN
			// Check creation succeeded
			require.NoError(test, err)
			require.NotNil(test, logger)

			// Verify settings
			assert.Equal(test, tc.wantEnv, logger.env)
			assert.Equal(test, tc.wantLevel, logger.level)
			assert.NotNil(test, logger.zap)

			// Write test log to verify that logger is working
			logger.Info("test message", nil)
			output := buffer.String()
			assert.Contains(test, output, "test message")

		})
	}
}

// Test Basic Logging Functionality
func TestLogger_Suite(test *testing.T) {
	testCases := []struct {
		name         string
		description  string
		message      string
		fields       map[string]any
		wantColor    bool
	}{
		{
			name: "Given simple message, logger formats correctly",
			description: "Tests basic message logging",
			message: "test message",
			fields: nil,
			wantColor: false,
		},
		{
			name: "Given message with fields, logger includes all fields",
			description: "Tests structured logging with fields",
			message: "user action",
			fields: map[string]interface{}{
				"user_id": "123",
				"action": "login",
			},
			wantColor: false,
		},
		{
			name: "Given development mode, logger uses color output",
			description: "Validates that Zap's pretty printing includes color codes",
			message: "test message",
			fields: map[string]interface{}{
				"user_id": "123",
				"action": "login",
			},
			wantColor: true,
		},
	}

	// Test runner loop
	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			// GIVEN
			// Create RAM space for current test
			buffer := setupTestBuffer(test)

			// Ensure RAM is cleared after test
			test.Cleanup(func() {
				cleanupLogger(test, buffer)
			})

			// WHEN
			// Tell logger to write to RAM
			logger, err := NewLogger(
				WithOutput(buffer),
				WithAlertLevel("info"),
				WithEnv("dev"),
			)
			require.NoError(test, err)

			// Write test log
			logger.Info(tc.message, tc.fields)

			// THEN
			// Get saved log from RAM
			output := buffer.String()

			// Check if output matches stuff we want:
			if tc.wantColor {
				// Check color codes
				assert.Contains(test, output, "\x1b[")
			}

			// Check msg
			assert.Contains(test, output, tc.message)

			// Check fields
			for key, value := range tc.fields {
				assert.Contains(test, output, key)
				assert.Contains(test, output, fmt.Sprint(value))
			}
		})
	}
}

// Test Specific Alert Level Behavior
func TestLogger_AlertLevel(test *testing.T) {
	testCases := []struct {
		name          string
		level         string
		message       string
		fields        map[string]any
		logFn         func(*Logger, string, map[string]any)
		wantText      string // Text that must be in output
	}{
		{
			name:     "Info level logs correctly",
			level:    "info",
			message:  "test info",
			fields:   map[string]any{"key": "value"},
			logFn:    (*Logger).Info,
			wantText: "info",
		},
		{
			name:     "Debug level logs correctly",
			level:    "debug",
			message:  "test debug",
			fields:   map[string]any{"key": "value"},
			logFn:    (*Logger).Debug,
			wantText: "debug",
		},
		{
			name:     "Warn level logs correctly",
			level:    "warn",
			message:  "test warn",
			fields:   map[string]any{"key": "value"},
			logFn:    (*Logger).Warn,
			wantText: "warn",
		},
		{
			name:     "Error level logs correctly",
			level:    "error",
			message:  "test error",
			fields:   map[string]any{"key": "value"},
			logFn:    (*Logger).Error,
			wantText: "error",
		},
	}

	// Test runner loop
	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			// GIVEN
			// Create RAM space for current test
			buffer := setupTestBuffer(test)

			// Ensure RAM is cleard after test
			test.Cleanup(func() {
				cleanupLogger(test, buffer)
			})

			// WHEN
			// Tell logger to write to RAM
			logger, err := NewLogger(
				WithOutput(buffer),
				WithAlertLevel(AlertLevel(tc.level)),
				WithEnv("dev"),
			)
			require.NoError(test, err)

			// Write test log, call specific log function
			tc.logFn(logger, tc.message, tc.fields)

			// THEN
			// Get saved log from RAM
			output := buffer.String()

			// Check that log level is correct
			assert.Contains(test, output, tc.wantText)

			// Check that message is correct
			assert.Contains(test, output, tc.message)

			// Check that fields are correct
			for key, value := range tc.fields {
				assert.Contains(test, output, key)
				assert.Contains(test, output, fmt.Sprint(value))
			}
		})
	}
}

// Test that we can verify request logging works as expected
func TestLogger_Middleware(test *testing.T) {
	testCases := []struct {
		name          string
		description   string
		method        string
		path          string
		query         string
		statusCode    int
		wantLevel     AlertLevel
	}{
		{
			name:        "Success request logs as info alert",
			method:      "GET",
			path:        "/test",
			query:       "key=value",
			statusCode:  200,
			wantLevel:   LevelInfo,
		},
		{
			name:        "Client error logs as warn alert",
			method:      "POST",
			path:        "/test",
			statusCode:  400,
			wantLevel:   LevelWarn,
		},
		{
			name:        "Server error logs as error alert",
			method:      "GET",
			path:        "/test",
			statusCode:  500,
			wantLevel:   LevelError,
		},
	}

	// Test runner loop
	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			// GIVEN
			// Create RAM space for current test
			buffer := setupTestBuffer(test)

			// Ensure RAM is cleared after test
			test.Cleanup(func() {
				cleanupLogger(test, buffer)
			})

			logger, err := NewLogger(
				WithOutput(buffer),
			)
			require.NoError(test, err)

			// Create test http handler that returns the desired status code
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			})

			// WHEN
			// Create test request
			url := tc.path
			if tc.query != "" {
				url += "?" + tc.query
			}
			request := httptest.NewRequest(tc.method, url, nil)
			recorder := httptest.NewRecorder()

			// Run middleware
			logger.LogMiddleware(handler).ServeHTTP(recorder, request)

			// THEN
			output := buffer.String()

			// Check alert level
			assert.Contains(test, output, tc.wantLevel)

			// Check logged request details
			assert.Contains(test, output, tc.method)
			assert.Contains(test, output, tc.path)
			if tc.query != "" {
				assert.Contains(test, output, tc.query)
			}
			assert.Contains(test, output, fmt.Sprint(tc.statusCode))

			// Check timing and client info
			assert.Contains(test, output, "duration")
			assert.Contains(test, output, "remote_addr")
			assert.Contains(test, output, "user_agent")

			// Verify that all fields from requestFields in middleware are logged
			assert.Contains(test, output, "method")
			assert.Contains(test, output, "path")
			assert.Contains(test, output, "query")
			assert.Contains(test, output, "status")
			assert.Contains(test, output, "duration")
			assert.Contains(test, output, "remote_addr")
			assert.Contains(test, output, "user_agent")
		})
	}
}

// Test panic logging works as expected
func TestLoggerMiddleware_PanicRecovery(test *testing.T) {
	testCases := []struct {
		name          string
		panicValue    any
		wantFields    []string
	}{
		{
			name:        "Given panic, logger recovers and logs string panic",
			panicValue:  "oops. I am error. Something broke.",
			wantFields:  []string{
				"error",
				"method",
				"path",
				"status",
			},
		},
		{
			name:        "Given panic, logger recovers and logs error panic",
			panicValue:  errors.New("oops. I am a new error. Something broke."),
			wantFields:  []string{
				"error",
				"method",
				"path",
				"status",
			},
		},
	}

	// Test runner loop
	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			// GIVEN
			// Create RAM space for current test
			buffer := setupTestBuffer(test)

			// Ensure RAM is cleared after test
			test.Cleanup(func() {
				cleanupLogger(test, buffer)
			})

			logger, err := NewLogger(
				WithOutput(buffer),
			)
			require.NoError(test, err)

			// Create handler that panics
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(tc.panicValue)
			})

			// WHEN
			request := httptest.NewRequest("GET", "/test", nil)
			recorder := httptest.NewRecorder()

			// THEN
			// Verify that the middleware recovers and logs the panic
			assert.Panics(test, func() {
				logger.LogMiddleware(handler).ServeHTTP(recorder, request)
			}, "Middleware should re-panic after logging")

			// Check the log output
			output := buffer.String()

			// Verify that the error was logged at error level
			assert.Contains(test, output, "error")

			// Verify that all expected fields were logged
			for _, field := range tc.wantFields {
				assert.Contains(test, output, field)
			}

			// Verify that the panic value was logged
			assert.Contains(test, output, fmt.Sprint(tc.panicValue))
		})
	}
}


// Test Edge Cases
func TestLogger_EdgeCases(test *testing.T) {
	edgeCases := []struct {
		name     string
		message  string
		fields   map[string]any
		wantErr bool
		wantColor bool
	} {
		{
			name: "Given empty message, logger handles gracefully",
			message: "",
			fields: nil,
			wantErr: false, // Most loggers accept empty msgs
			wantColor: false,
		},
		{
			name: "Given nil fields map, logger handles gracefully",
			message: "test",
			fields: nil,
			wantErr: false,
			wantColor: false,
		},
		{
			name: "Given large message, logger handles w/o panicking",
			message: strings.Repeat("a", 1024*1024), // 1MB string
			fields: nil,
			wantErr: false,
			wantColor: true,
		},
		{
			name: "Given fields with nil values, logger handles gracefully",
			message: "test",
			fields: map[string]any{"key": nil},
			wantErr: false,
			wantColor: true,
		},
		{
			name: "Given fields with complex types, logger handles gracefully",
			message: "test",
			fields: map[string]any{
				"slice": []string{"a", "b", "c"},
				"map": map[string]int{"a": 1, "b": 1},
				"struct": struct{ Name string }{"test"},
			},
			wantErr: false,
			wantColor: true,
		},
	}

	// Test runner loop
	for _, tc := range edgeCases {
		test.Run(tc.name, func(t *testing.T) {
			// GIVEN
			// Create RAM space for current test
			buffer := setupTestBuffer(test)

			// Ensure RAM is cleared after test
			test.Cleanup(func() {
				cleanupLogger(test, buffer)
			})

			// WHEN
			// Tell logger to write to RAM
			logger, err := NewLogger(
				WithOutput(buffer),
				WithAlertLevel("info"),
				WithEnv("dev"),
			)
			require.NoError(test, err)

			// Write test log
			logger.Info(tc.message, tc.fields)

			// THEN
			// Get saved log from RAM
			output := buffer.String()

			// Check if output matches stuff we want
			// Check msg
			assert.Contains(test, output, tc.message)

			// Check fields
			for key, value := range tc.fields {
				assert.Contains(test, output, key)
				assert.Contains(test, output, fmt.Sprint(value))
			}
		})
	}
}

// Setup tests, cleanup env, custom assertions
func setupTestBuffer(test testing.TB) *bytes.Buffer {
	// Mark this as a helper fn for better failure reporting
	test.Helper()

	// Create new expandable RAM space for capturing log output
	return new(bytes.Buffer)
}

func cleanupLogger(test testing.TB, buffer *bytes.Buffer) {
	test.Helper()

	// Clear RAM space we just used for testing
	buffer.Reset()
}


// Fn - Benchmarks
func BenchmarkLogger(benchmark *testing.B) {
	benchmarks := []struct {
		name      string
		message   string
		fields    map[string]any
		wantColor bool
	} {
		{
			name:       "Basic message without fields",
			message:    "test message",
			fields:     nil,
			wantColor:  false,
		},
		{
			name:       "Message with single field",
			message:    "test message",
			fields:     map[string]any{
				"key": "value",
			},
			wantColor: true,
		},
		{
			name:        "Message with multiple fields",
			message:     "test message",
			fields:      map[string]any{
				"string": "value",
				"number": 123,
				"bool":   true,
			},
			wantColor: true,
		},
		{
			name:        "Message with nested fields",
			message:     "test message",
			fields: map[string]interface{}{
				"user": map[string]interface{}{
						"id":   "123",
						"name": "test",
				},
				"metadata": map[string]interface{}{
						"region": "us-east-1",
						"env":    "prod",
				},
			},
			wantColor: true,
		},
	}

	// Test runner loop
	for _, bm := range benchmarks {
		benchmark.Run(bm.name, func(b *testing.B) {
			// Setup
			buffer := setupTestBuffer(b)
			b.Cleanup(func() {
				cleanupLogger(b, buffer)
			})

			logger, err := NewLogger(
				WithOutput(buffer),
				WithAlertLevel("info"),
				WithEnv("dev"),
			)
			require.NoError(b, err)

			// Reset timer before actual benchmark
			b.ResetTimer()

			// Run benchmark
			for index := 0; index < b.N; index++ {
				logger.Info(bm.message, bm.fields)
				buffer.Reset()
			}
		})
	}

}