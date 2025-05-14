package igdb

import (
	"errors"
	"fmt"
)

const (
	QueryBuildOperation = "QueryBuild"
	QueryValidateOperation = "QueryValidate"
)

// QueryError represents an error that occurred during query building
type IGDBQueryError struct {
	Op  string // The operation that failed (e.g., "QueryBuild", "Validate")
	Err error  // The underlying error
}

// Error implements the error interface for QueryError.
// It returns a string representation of the error in the format "operation: error".
func (e *IGDBQueryError) Error() string {
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error, allowing for error inspection and comparison.
// This implements the errors.Unwrap interface introduced in Go 1.13.
func (e *IGDBQueryError) Unwrap() error {
	return e.Err
}

// NewQueryError creates a new QueryError with the given operation and error.
// If err is nil, it will be treated as an empty error.
func NewIGDBQueryError(op string, err error) error {
	return &IGDBQueryError{
		Op:  op,
		Err: err,
	}
}


// Common query error types
var (
	// ErrEmptyQuery is returned when a query is built without any components
	ErrEmptyQuery = NewIGDBQueryError(QueryBuildOperation, fmt.Errorf("query is empty - at least one component (fields, search, where) is require"))

	// ErrNoFields is returned when a query is built without any fields specified
	ErrNoFields = NewIGDBQueryError(QueryBuildOperation, fmt.Errorf("no fields specified - at least one field must be selected"))

	// ErrInvalidSearchTerm is returned when a search term is empty or invalid
	ErrInvalidSearchTerm = NewIGDBQueryError(QueryValidateOperation, fmt.Errorf("invalid search term - search term must be a valid, non-empty string"))

	// ErrInvalidWhereCondition is returned when a where condition is invalid
	ErrInvalidWhereCondition = NewIGDBQueryError(QueryValidateOperation, fmt.Errorf("invalid where condition"))
)

// NewInvalidLimitError creates a new error for invalid limit values
// This is a fn instead of a const because the error msg needs to include the actual limit value
func NewInvalidLimitError(limit int) error {
	return NewIGDBQueryError(QueryValidateOperation, fmt.Errorf("invalid query limit: %d (must be between %d and %d)",
		limit, MinLimit, MaxLimit))
}

// Helper fn for error type checking - checks if the given error is a IGDBQueryError
func IsIGDBQueryError(err error) bool {
	var igdbQueryErr *IGDBQueryError
	return errors.As(err, &igdbQueryErr)
}

// Returns the operation that failed if the error is an IGDBQueryError
// Returns an empty string if error is not an IGDBQueryError
func GetIGDBQueryErrorOperation(err error) string {
	var igdbQueryErr *IGDBQueryError
	if errors.As(err, &igdbQueryErr) {
		return igdbQueryErr.Op
	}
	return ""
}
