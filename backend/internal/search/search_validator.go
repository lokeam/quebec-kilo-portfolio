package search

import (
	"fmt"
	"unicode/utf8"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
)

// Validation consts
const (
	MinQueryLength = 2
	MaxQueryLength = 100
	MinResultLimit = 1
	MaxResultLimit = 50
	MaxResultOffset = 500
)


// ValidationError struct
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}


// SearchValidator struct
type SearchValidator struct {
	maxQueryLength    int
	maxResultLimit    int
	maxResultOffset   int
	sanitizer         interfaces.Sanitizer
}

func NewSearchValidator(sanitizer interfaces.Sanitizer) (*SearchValidator, error) {
	return &SearchValidator{
		maxQueryLength: MaxQueryLength,
		maxResultLimit: MaxResultLimit,
		maxResultOffset: MaxResultOffset,
		sanitizer: sanitizer,
	}, nil
}

func (v *SearchValidator) ValidateQuery(query searchdef.SearchQuery) error {
	var searchQueryViolations []string

	// Query validation
	if sanitized, err := v.validateQueryString(query.Query); err != nil {
		searchQueryViolations = append(searchQueryViolations, err.Error())
	} else {
		// Use sanitized query for further processing
		query.Query = sanitized
	}

	// Result limit validation
	if err := v.validateResultLimit(query.Limit, 0); err != nil {
		searchQueryViolations = append(searchQueryViolations, err.Error())
	}

	// Field validation
	// if err := v.validateFields(query.Fields); err != nil {
	// 	searchQueryViolations = append(searchQueryViolations, err.Error())
	// }

	if len(searchQueryViolations) > 0 {
		return &ValidationError{
			Field: "query",
			Message: fmt.Sprintf("Search query validation failed: %v", searchQueryViolations),
		}
	}

	return nil
}

// Individual validation rules
func (v *SearchValidator) validateQueryString(query string) (string, error) {
	length := utf8.RuneCountInString(query)

	// Check if query is at least 2 characters long
	if length < MinQueryLength {
		return "", &ValidationError{
			Field: "query",
			Message: fmt.Sprintf("query must be at least %d characters", MinQueryLength),
		}
	}

	// Check if query length is shorter than max allowed
	if length > v.maxQueryLength {
		return "", &ValidationError{
			Field: "query",
			Message: fmt.Sprintf("search query must be less than %d characters", v.maxQueryLength),
		}
	}

	// Sanitize to prevent SQL injection/XSS
	sanitized, err := v.sanitizer.SanitizeSearchQuery(query)
	if err != nil {
		return "", &ValidationError{
			Field: "query",
			Message: fmt.Sprintf("invalid query content: %v", err),
		}
	}

	return sanitized,nil
}

func (v *SearchValidator) validateResultLimit(limit, offset int) error {
	if limit < MinResultLimit || limit > v.maxResultLimit {
		return &ValidationError{
			Field: "limit",
			Message: fmt.Sprintf("result limit must be between %d and %d", MinResultLimit, v.maxResultLimit),
		}
	}

	if offset < 0 || offset > v.maxResultOffset {
		return &ValidationError{
			Field: "offset",
			Message: fmt.Sprintf("result offset must be less than %d", v.maxResultOffset),
		}
	}

	return nil
}

