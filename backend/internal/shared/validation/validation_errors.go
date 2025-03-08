package validation

import "fmt"

type ValidationError struct {
	Field    string
	Message  string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", ve.Field, ve.Message)
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}