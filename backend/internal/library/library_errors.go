package library

import (
	"errors"
	"net/http"
)

// Package errors with errors.Is
var (
	ErrGameNotFound = errors.New("game not found in library")
	ErrDatabaseError = errors.New("database error")
	ErrValidationFailed = errors.New("validation failed")
	ErrInvalidGameData = errors.New("invalid game data")
	ErrUnauthorizedGame = errors.New("unauthorized: game does not belong to user")
	ErrInvalidGameID = errors.New("invalid game ID")
	ErrInvalidUserID = errors.New("invalid user ID")
	ErrDuplicateGame = errors.New("game already exists in library")
)

// GetStatusCodeForError returns the appropriate HTTP status code for a given error
func GetStatusCodeForError(err error) int {
	switch {
	case errors.Is(err, ErrGameNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrValidationFailed):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorizedGame):
		return http.StatusForbidden
	case errors.Is(err, ErrDuplicateGame):
		return http.StatusConflict
	case errors.Is(err, ErrDatabaseError):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}