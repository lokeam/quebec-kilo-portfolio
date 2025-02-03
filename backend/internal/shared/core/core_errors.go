package core

import "errors"

// Standard error types used across the application
var (
	// ErrValidation represents validation-related errors.
	// Used when request data fails validation requirements.
	// Maps to HTTP 400 Bad Request in HTTP responses.
	ErrValidation     = errors.New("validation error")

	// ErrAuthentication represents authentication-related errors.
	// Used when authentication fails or credentials are invalid.
	// Maps to HTTP 401 Unauthorized in HTTP responses.
	ErrAuthentication = errors.New("authentication error")
)

// ErrorResponse defines the standard error response structure
// used across the application for consistent error handling.
type ErrorResponse struct {
	// Error contains the error message to be returned to the client
	Error       string `json:"error"`

	// RequestID is a unique identifier for the request,
	// used for tracking and debugging purposes
	RequestID   string `json:"requestId"`
}
