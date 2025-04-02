package sublocation

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	validationErrors "github.com/lokeam/qko-beta/internal/shared/validation"
)

// Validation constants
const (
	MaxNameLength = 100
	MinCapacity   = 1
	MaxCapacity   = 1000
)

// Valid sublocation types
var ValidSublocationTypes = []string{"shelf", "cabinet", "drawer", "box", "display", "other"}

// Valid background colors
var ValidBgColors = []string{
	"red", "blue", "green", "yellow", "purple",
	"orange", "black", "white", "gray",
}

// SublocationValidator struct
type SublocationValidator struct {
	sanitizer interfaces.Sanitizer
}

// NewSublocationValidator creates a new sublocation validator
func NewSublocationValidator(sanitizer interfaces.Sanitizer) (*SublocationValidator, error) {
	if sanitizer == nil {
		return nil, fmt.Errorf("sanitizer cannot be nil")
	}

	return &SublocationValidator{
		sanitizer: sanitizer,
	}, nil
}

// Validation Sublocation validates a sublocation model
func (v *SublocationValidator) ValidateSublocation(sublocation *models.Sublocation) (models.Sublocation, error) {
	var validatedSublocation models.Sublocation
	var violations []string

	// Validate name
	if sanitizedName, err := v.validateName(sublocation.Name); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedSublocation.Name = sanitizedName
	}

	// Validate location type
	if sanitizedType, err := v.validateLocationType(sublocation.LocationType); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedSublocation.LocationType = sanitizedType
	}

	// Validate background color
	if sanitizedColor, err := v.validateBgColor(sublocation.BgColor); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedSublocation.BgColor = sanitizedColor
	}

	// Validate capacity
	if validatedCapacity, err := v.validateCapacity(sublocation.Capacity); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedSublocation.Capacity = validatedCapacity
	}

	// Copy other fields that don't need validation
	validatedSublocation.ID = sublocation.ID
	validatedSublocation.UserID = sublocation.UserID
	validatedSublocation.CreatedAt = sublocation.CreatedAt
	validatedSublocation.UpdatedAt = sublocation.UpdatedAt
	validatedSublocation.Items = sublocation.Items

	if len(violations) > 0 {
		return models.Sublocation{}, &validationErrors.ValidationError{
			Field:   "sublocation",
			Message: strings.Join(violations, "; "),
		}
	}

	return validatedSublocation, nil
}

func (v *SublocationValidator) validateName(name string) (string, error) {
	// Check if name is empty
	if name == "" {
		return "", &validationErrors.ValidationError{
			Field:     "name",
			Message:   "name cannot be empty",
		}
	}

	// Check name length
	if utf8.RuneCountInString(name) > MaxNameLength {
		return "", &validationErrors.ValidationError{
			Field:     "name",
			Message:   fmt.Sprintf("name must be less than %d characters", MaxNameLength),
		}
	}

	// Sanitize name
	sanitized, err := v.sanitizer.SanitizeString(name)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:     "name",
			Message:   fmt.Sprintf("invalid name content: %v", err),
		}
	}

	return sanitized, nil
}

func (v *SublocationValidator) validateLocationType(locationType string) (string, error) {
	// First sanitize the location type
	sanitized, err := v.sanitizer.SanitizeString(locationType)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "locationType",
			Message: fmt.Sprintf("invalid location type content: %v", err),
		}
	}

	// Then check if location type is valid
	isValid := false
	for _, validType := range ValidSublocationTypes {
		if strings.EqualFold(sanitized, validType) {
			isValid = true
			break
		}
	}

	if !isValid {
		return "", &validationErrors.ValidationError{
			Field:   "locationType",
			Message: fmt.Sprintf("location type must be one of %v", ValidSublocationTypes),
		}
	}

	return sanitized, nil
}

func (v *SublocationValidator) validateBgColor(bgColor string) (string, error) {
	// First sanitize the background color
	sanitized, err := v.sanitizer.SanitizeString(bgColor)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "bgColor",
			Message: fmt.Sprintf("invalid background color content: %v", err),
		}
	}

	// Next check if background color is valid
	isValid := false
	for _, validColor := range ValidBgColors {
		if strings.EqualFold(bgColor, validColor) {
			isValid = true
			break
		}
	}

	if !isValid {
		return "", &validationErrors.ValidationError{
			Field:     "bgColor",
			Message:   fmt.Sprintf("background color must be one of %v", ValidBgColors),
		}
	}

	return sanitized, nil
}

func (v *SublocationValidator) validateCapacity(capacity int) (int, error) {
	// Check if capacity is positive
	if capacity <= 0 {
		return 0, &validationErrors.ValidationError{
			Field:   "capacity",
			Message: "capacity must be a positive number",
		}
	}

	// Check if capacity is within reasonable limits
	if capacity > MaxCapacity {
		return 0, &validationErrors.ValidationError{
			Field:   "capacity",
			Message: fmt.Sprintf("capacity must not exceed %d", MaxCapacity),
		}
	}

	return capacity, nil
}