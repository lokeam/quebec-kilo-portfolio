package library

import "errors"

// Package errors with errors.Is
var (
	ErrGameNotFound = errors.New("game not found in library")
	ErrInvalidGameID = errors.New("invalid game ID")
	ErrInvalidUserID = errors.New("invalid user ID")
)