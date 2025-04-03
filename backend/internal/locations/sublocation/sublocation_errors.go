package sublocation

import "errors"

var (
	ErrSublocationNotFound = errors.New("sublocation not found")
	ErrDatabaseError = errors.New("database error")
	ErrValidationFailed = errors.New("validation failed")
)