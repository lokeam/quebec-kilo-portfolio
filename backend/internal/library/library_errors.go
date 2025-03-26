package library

import "errors"

// Package errors with errors.Is
var (
	ErrGameNotFound = errors.New("game not found in library")
	ErrDatabaseConnection = errors.New("database connection error")
	ErrInvalidGameData = errors.New("invalid game data")
	ErrInvalidGameID = errors.New("invalid game ID")
	ErrInvalidUserID = errors.New("invalid user ID")
)
