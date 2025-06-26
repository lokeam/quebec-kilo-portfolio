package digital

import (
	"errors"
	"net/http"
)

var (
	ErrDigitalLocationNotFound = errors.New("digital location not found")
	ErrDatabaseError = errors.New("database error")
	ErrValidationFailed = errors.New("validation failed")
	ErrDigitalLocationExists = errors.New("digital location already exists")
	ErrInvalidInput = errors.New("invalid input parameters")
	ErrNotFound = errors.New("one or more locations not found")
	ErrDatabase = errors.New("database error")
	ErrTransaction = errors.New("transaction error")
	ErrPartialDeletion = errors.New("partial deletion occurred")
	ErrRelatedRecords = errors.New("error deleting related records")
	ErrRetryable = errors.New("retryable error occurred")
	ErrInvalidQueryParameter = errors.New("invalid query parameter format")
	ErrEmptyLocationIDs = errors.New("no location IDs provided")
	ErrInvalidLocationID = errors.New("invalid location ID format")
)

func GetStatusCodeForError(err error) int {
	switch {
		case errors.Is(err, ErrDigitalLocationNotFound):
			return http.StatusNotFound
		case errors.Is(err, ErrValidationFailed):
			return http.StatusBadRequest
		case errors.Is(err, ErrDigitalLocationExists):
			return http.StatusConflict
		case errors.Is(err, ErrDatabaseError):
			return http.StatusInternalServerError
		default:
			return http.StatusInternalServerError
	}
}