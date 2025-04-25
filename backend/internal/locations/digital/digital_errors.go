package digital

import "errors"

var (
	ErrDigitalLocationNotFound = errors.New("digital location not found")
	ErrDatabaseError = errors.New("database error")
	ErrValidationFailed = errors.New("validation failed")
	ErrDigitalLocationExists = errors.New("digital location already exists")

	// ErrInvalidInput is returned when input parameters are invalid
	ErrInvalidInput = errors.New("invalid input parameters")

	// ErrNotFound is returned when one or more locations are not found
	ErrNotFound = errors.New("one or more locations not found")

	// ErrDatabase is returned for database-related errors
	ErrDatabase = errors.New("database error")

	// ErrTransaction is returned for transaction-related errors
	ErrTransaction = errors.New("transaction error")

	// ErrPartialDeletion is returned when some but not all locations were deleted
	ErrPartialDeletion = errors.New("partial deletion occurred")

	// ErrRelatedRecords is returned when there are issues deleting related records
	ErrRelatedRecords = errors.New("error deleting related records")

	// ErrRetryable is returned for errors that might succeed on retry
	ErrRetryable = errors.New("retryable error occurred")
)