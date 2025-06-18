package dashboard

import (
	"errors"
	"net/http"
)

var (
	ErrDatabaseError = errors.New("database error")
	ErrInvalidUserID = errors.New("invalid userID")
	ErrValidationFailed = errors.New("validation failed")
	ErrDashboardItemNotFound = errors.New("dashboard item not found")
)

func GetStatusCodeForError(err error) int {
	switch {
	case errors.Is(err, ErrDashboardItemNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInvalidUserID):
		return http.StatusBadRequest
	case errors.Is(err, ErrValidationFailed):
		return http.StatusBadRequest
	case errors.Is(err, ErrDatabaseError):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}