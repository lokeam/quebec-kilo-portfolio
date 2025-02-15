package searchdef

import (
	"fmt"
)

// Error message constants to avoid magic strings
const (
	ErrMsgConnectionFailed = "failed to connect to IGDB"
	ErrMsgInvalidResponse  = "invalid response from IGDB"
	ErrMsgNoResults       = "no results found"
)

// Sentinel errors
var (
	ErrIGDBConnection = fmt.Errorf(ErrMsgConnectionFailed)
	ErrIGDBRequest = fmt.Errorf(ErrMsgInvalidResponse)
	ErrEmptyResponse = fmt.Errorf(ErrMsgNoResults)
)

// ErrorCode represents specific error condtions
type ErrorCode int

const (
	CodeSuccess ErrorCode = iota // iota makes these auto-increment starting from 0
	CodeInvalidRequest          // 1
	CodeAPIError                // 2
	CodeNotFound                // 3
	CodeTimeout                 // 4
	CodeRateLimit               // IGDB rate limiting
	CodeServerError             // IGDB server errors
)

// Provides detailed error info providing rich context
type RepositoryError struct {
	FailedOperation  string    // Specific operation that failed
	Code             ErrorCode  // Error classification
	ErrorMessage     string     // Human readable error message
	Err              error      // Original error, if any
}

// String representations of code errors for logging + debugging
func (c ErrorCode) String() string {
	switch c {
	case CodeSuccess:
			return "SUCCESS"
	case CodeInvalidRequest:
			return "INVALID_REQUEST"
	case CodeAPIError:
			return "API_ERROR"
	case CodeNotFound:
			return "NOT_FOUND"
	case CodeTimeout:
			return "TIMEOUT"
	case CodeRateLimit:
			return "RATE_LIMIT"
	case CodeServerError:
			return "SERVER_ERROR"
	default:
			return "UNKNOWN"
	}
}

func (e *RepositoryError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.FailedOperation, e.ErrorMessage, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.FailedOperation, e.ErrorMessage)
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// Constructor for creating new SearchRepository errors
func NewSearchRepositoryError(
	op string,
	code ErrorCode,
	message string,
	err error,
) *RepositoryError {
	return &RepositoryError{
		FailedOperation: op,
		Code:            code,
		ErrorMessage:    message,
		Err:             err,
	}
}
