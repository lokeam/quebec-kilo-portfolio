package physical

import "errors"

var (
	ErrLocationNotFound = errors.New("physical location not found")
	ErrDatabaseError = errors.New("database error")
	ErrValidationFailed = errors.New("validation failed")
)