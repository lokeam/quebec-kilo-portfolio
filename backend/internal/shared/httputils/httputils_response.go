package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/lokeam/qko-beta/internal/shared/core"
	"github.com/lokeam/qko-beta/internal/shared/logger"
)

// RespondWithJSON writes a JSON response with provided data + status code.
// Creates a standarized way to write responses that include:
//   - Content-Type header management
//   - Response state validation
//   - Comprehensive logging
//   - Performance timing
//   - Error handling
//
// The function is safe for concurrent use + handles cases where headers
// might have been previously written.
//
// Parameters:
//   - w: ResponseWriter interface used to send the HTTP response
//   - logger: Logger instance for tracking response process
//   - status: HTTP status code that will be set in the response
//   - data: The data to be encoded as JSON within the response
//
// Example Usage:
//
//	data := map[string]string{"message": "Success"}
//	err := RespondWithJSON(w, logger, http.StatusOK, data)
//	if err != nil {
//	    // Handle error
//	}
func RespondWithJSON(
	w http.ResponseWriter,
	logger logger.LoggerInterface,
	status int,
	data any,
	) error {

		logger.Debug("preparing response", map[string]any{
				"status": status,
				"contentType": "application/json",
		})
    // Start timing the response
    start := time.Now()

    // Log header state before write attempt
    if rw, ok := w.(interface{ Written() bool }); ok {
			logger.Debug("checking response writer state", map[string]any{
				"status": rw.Written(),
				"contentType": "application/json",
			})
		}

    // Only set Content-Type if not already set
    if w.Header().Get("Content-Type") == "" {
        logger.Debug("setting content-type header", map[string]any{
				"contentType": "application/json",
		})
        w.Header().Set("Content-Type", "application/json")
    }

    // Check if we can write headers
    if rw, ok := w.(interface{ Written() bool }); !ok || !rw.Written() {
        logger.Debug("writing headers", map[string]any{
				"contentType": "application/json",
				"status": status,
		})
        w.WriteHeader(status)
    }

    // Encode data to JSON and write to response
    if err := json.NewEncoder(w).Encode(data); err != nil {
        logger.Error("failed to encode response", map[string]any{
				"error": err,
				"status": status,
				"duration": time.Since(start),
		})
        return fmt.Errorf("failed to encode response: %w", err)
    }

    // Log successful response
    logger.Debug("response sent successfully", map[string]any{
			"status": status,
			"duration": time.Since(start),
		})

    return nil
}

// RespondWithError creates standard error messages across the app.
// It matches error types to the right status codes + makes sure that
// all error messages follow the same format.
//
// Responsibilities include:
//   - Error type classification
//   - Status code mapping
//   - Request tracking
//   - Error logging
//   - Performance monitoring
//
// Parameters:
//   - w: http.ResponseWriter for sending the HTTP response
//   - logger: Logger instance for error tracking and monitoring
//   - requestID: Unique identifier for request tracing
//   - err: The error to be processed and returned
//
// Returns:
//   - error: nil on successful error response, error if response sending fails
//
// Status Code Mapping:
//   - 400 Bad Request: Validation errors (core.ErrValidation)
//   - 401 Unauthorized: Authentication failures (core.ErrAuthentication)
//   - 500 Internal Server Error: All other error types
//
// Response Format:
//   {
//     "error": "error message",
//     "requestId": "unique-request-id"
//   }
//
// Example Usage:
//
//	if err != nil {
//	    return RespondWithError(w, logger, requestID, err)
//	}
func RespondWithError(
	w ResponseWriter,
	logger logger.LoggerInterface,
	requestID string,
	err error,
	) error {
	// Start timing the response
	start := time.Now()

	// Determine HTTP status code based on error type
	status := http.StatusInternalServerError
	if errors.Is(err, core.ErrValidation) {
			status = http.StatusBadRequest
	} else if errors.Is(err, core.ErrAuthentication) {
			status = http.StatusUnauthorized
	}

	// Create error response
	response := core.ErrorResponse{
			Error:     err.Error(),
			RequestID: requestID,
	}

	// Log error with context
	logger.Debug("preparing error response", map[string]any{
		"requestID": requestID,
		"status": status,
		"errorType": fmt.Sprintf("%T", err),
	})

	logger.Error("sending error response", map[string]any{
		"error": err,
		"requestId": requestID,
		"status": time.Since(start),
	})

	// Use existing respondWithJSON to send response
	return RespondWithJSON(w, logger, status, response)
}
