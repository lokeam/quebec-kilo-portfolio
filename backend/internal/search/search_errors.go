package search

import (
	"errors"
	"net/http"
)

var (
	ErrLocationNotFound = errors.New("location not found")
	ErrDatabaseError = errors.New("database error")
	ErrValidationFailed = errors.New("validation failed")
	ErrUnauthorizedLocation = errors.New("unauthorized: location does not belong to user")
	ErrDuplicateLocation    = errors.New("a physical location with this name already exists")
	ErrEmptyLocationIDs = errors.New("no location IDs provided")
)

// GetStatusCodeForError returns the appropriate HTTP status code for a given error
func GetStatusCodeForError(err error) int {
	switch {
	case errors.Is(err, ErrLocationNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrValidationFailed):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorizedLocation):
		return http.StatusForbidden
	case errors.Is(err, ErrDuplicateLocation):
		return http.StatusConflict
	case errors.Is(err, ErrDatabaseError):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}