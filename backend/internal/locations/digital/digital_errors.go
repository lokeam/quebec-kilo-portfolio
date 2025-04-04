package digital

import "errors"

var (
	ErrDigitalLocationNotFound = errors.New("digital location not found")
	ErrDatabaseError = errors.New("database error")
	ErrValidationFailed = errors.New("validation failed")
)