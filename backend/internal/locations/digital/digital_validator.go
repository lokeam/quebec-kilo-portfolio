package digital

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	validationErrors "github.com/lokeam/qko-beta/internal/shared/validation"
)

const (
	MaxNameLength = 100
	MaxURLLength = 2048
)

type DigitalValidator struct {
	sanitizer interfaces.Sanitizer
}

func NewDigitalValidator(sanitizer interfaces.Sanitizer) (*DigitalValidator, error) {
	return &DigitalValidator{
		sanitizer: sanitizer,
	}, nil
}

func (v *DigitalValidator) ValidateDigitalLocation(
	location models.DigitalLocation,
) (models.DigitalLocation, error) {
	var validatedLocation models.DigitalLocation
	var violations []string

	// Valiate name
	if sanitizedName, err := v.validateName(location.Name); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedLocation.Name = sanitizedName
	}

	// Validate URL
	if sanitizedURL, err := v.validateURL(location.URL); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedLocation.URL = sanitizedURL
	}

	validatedLocation.ID = location.ID
	validatedLocation.UserID = location.UserID
	validatedLocation.IsActive = location.IsActive
	validatedLocation.CreatedAt = location.CreatedAt
	validatedLocation.UpdatedAt = location.UpdatedAt

	if len(violations) > 0 {
		return models.DigitalLocation{}, &validationErrors.ValidationError{
			Field:   "location",
			Message: fmt.Sprintf("Digital location validation failed: %v", violations),
		}
	}
	return validatedLocation, nil
}

func (v *DigitalValidator) validateName(name string) (string, error) {
	// Check if name is empty
	if name == "" {
		return "", &validationErrors.ValidationError{
			Field:   "name",
			Message: "name cannot be empty",
		}
	}

	// Check name length
	length := utf8.RuneCountInString(name)
	if length > MaxNameLength {
		return "", &validationErrors.ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("name must be less than %d characters", MaxNameLength),
		}
	}

	// Sanitize name
	sanitized, err := v.sanitizer.SanitizeString(name)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("invalid name content: %v", err),
		}
	}

	return sanitized, nil
}

func (v *DigitalValidator) validateURL(urlStr string) (string, error) {
	// Check URL is empty
	if urlStr == "" {
		return "", &validationErrors.ValidationError{
			Field:     "url",
			Message:   "URL cannot be empty",
		}
	}

	// Check URL length
	length := utf8.RuneCountInString(urlStr)
	if length > MaxURLLength {
		return "", &validationErrors.ValidationError{
			Field:     "url",
			Message:   fmt.Sprintf("URL must be less than %d characters", MaxURLLength),
		}
	}

	// Validate URL format
	parsedURL, err := url.Parse(urlStr)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == ""{
		return "", &validationErrors.ValidationError{
			Field:      "url",
			Message:    "invalid URL format",
		}
	}

	// Ensure URL has http or https scheme
	if !strings.HasPrefix(parsedURL.Scheme, "http") {
		return "", &validationErrors.ValidationError{
			Field:     "url",
			Message:   "URL must use http or https protocol",
		}
	}

	// Sanitize URL
	sanitized, err := v.sanitizer.SanitizeString(urlStr)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:     "url",
			Message:   fmt.Sprintf("invalid URL content: %v", err),
		}
	}

	return sanitized, nil
}
