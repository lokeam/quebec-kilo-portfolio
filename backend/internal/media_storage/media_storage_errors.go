package media_storage

import (
	"errors"
	"net/http"
)

var (
	// ErrStorageStatsNotFound is returned when storage statistics cannot be found
	ErrStorageStatsNotFound = errors.New("storage statistics not found")

	// ErrInvalidUserID is returned when the user ID is invalid or missing
	ErrInvalidUserID = errors.New("invalid user ID")

	// ErrAnalyticsServiceUnavailable is returned when the analytics service is not available
	ErrAnalyticsServiceUnavailable = errors.New("analytics service unavailable")

	// ErrCacheOperationFailed is returned when a cache operation fails
	ErrCacheOperationFailed = errors.New("cache operation failed")

	// ErrDatabaseError is returned for database-related errors
	ErrDatabaseError = errors.New("database error")

	// ErrInvalidInput is returned when input parameters are invalid
	ErrInvalidInput = errors.New("invalid input parameters")

	// ErrRetryable is returned for errors that might succeed on retry
	ErrRetryable = errors.New("retryable error occurred")

	// ErrPartialData is returned when only partial data could be retrieved
	ErrPartialData = errors.New("partial data retrieved")
)

// GetStatusCodeForError returns the appropriate HTTP status code for a given error
func GetStatusCodeForError(err error) int {
	switch {
	case errors.Is(err, ErrStorageStatsNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInvalidUserID):
		return http.StatusUnauthorized
	case errors.Is(err, ErrAnalyticsServiceUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, ErrCacheOperationFailed):
		return http.StatusInternalServerError
	case errors.Is(err, ErrDatabaseError):
		return http.StatusInternalServerError
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, ErrRetryable):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrPartialData):
		return http.StatusPartialContent
	default:
		return http.StatusInternalServerError
	}
}
