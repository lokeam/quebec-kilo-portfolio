package spend_tracking

import (
	"errors"
	"net/http"
)

var (
	ErrDatabaseError = errors.New("database error")
	ErrSpendTrackingItemNotFound = errors.New("spend tracking item not found")
	ErrInvalidUserID = errors.New("invalid userID")
	ErrInvalidSpendTrackingItemData = errors.New("invalid spend tracking item data")
	ErrUnauthorizedSpendTrackingItem = errors.New("unauthorized: spend tracking item does not belong to user")
	ErrInvalidSpendTrackingItem = errors.New("invalid spend tracking item ID")
	ErrValidationFailed = errors.New("validation failed")
	ErrEmptySpendTrackingIDs = errors.New("no spend tracking IDs provided")
)

// GetStatusCodeForError returns the appropriate HTTP status code for a given error
func GetStatusCodeForError(err error) int {
	switch {
	case errors.Is(err, ErrSpendTrackingItemNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInvalidUserID):
		return http.StatusBadRequest
	case errors.Is(err, ErrInvalidSpendTrackingItemData):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorizedSpendTrackingItem):
		return http.StatusForbidden
	case errors.Is(err, ErrInvalidSpendTrackingItem):
		return http.StatusBadRequest
	case errors.Is(err, ErrValidationFailed):
		return http.StatusBadRequest
	case errors.Is(err, ErrDatabaseError):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}