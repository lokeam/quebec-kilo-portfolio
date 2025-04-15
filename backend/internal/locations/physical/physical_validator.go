package physical

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	validationErrors "github.com/lokeam/qko-beta/internal/shared/validation"
)

// Validation constants
const (
	MaxNameLength     = 100
	MaxLabelLength    = 50
	CoordinatePattern = `^-?\d+(\.\d+)?,\s*-?\d+(\.\d+)?$`
)

// Valid location types
var ValidLocationTypes = []string{"house", "apartment", "office", "warehouse"}

// Valid background colors
var ValidBgColors = []string{"red", "green", "blue", "orange", "gold", "purple", "brown", "gray"}

// PhysicalValidator struct
type PhysicalValidator struct {
	sanitizer interfaces.Sanitizer
}

func NewPhysicalValidator(sanitizer interfaces.Sanitizer) (*PhysicalValidator, error) {
	return &PhysicalValidator{
		sanitizer: sanitizer,
	}, nil
}

func (v *PhysicalValidator) ValidatePhysicalLocation(location models.PhysicalLocation) (models.PhysicalLocation, error) {
	var validatedLocation models.PhysicalLocation
	var violations []string

	// Validate name
	if sanitizedName, err := v.validateName(location.Name); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedLocation.Name = sanitizedName
	}

	// Validate label
	if sanitizedLabel, err := v.validateLabel(location.Label); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedLocation.Label = sanitizedLabel
	}

	// Validate location type
	if sanitizedType, err := v.validateLocationType(location.LocationType); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedLocation.LocationType = sanitizedType
	}

	// Validate map coordinates
	if sanitizedCoords, err := v.validateMapCoordinates(location.MapCoordinates); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedLocation.MapCoordinates = sanitizedCoords
	}

	// Validate background color
	if sanitizedColor, err := v.validateBgColor(location.BgColor); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedLocation.BgColor = sanitizedColor
	}

	// Copy other fields that don't need validation
	validatedLocation.ID = location.ID
	validatedLocation.UserID = location.UserID
	validatedLocation.CreatedAt = location.CreatedAt
	validatedLocation.UpdatedAt = location.UpdatedAt

	if len(violations) > 0 {
		return models.PhysicalLocation{}, &validationErrors.ValidationError{
			Field:   "location",
			Message: fmt.Sprintf("Physical location validation failed: %v", violations),
		}
	}

	return validatedLocation, nil
}

// Individual validation rules
func (v *PhysicalValidator) validateName(name string) (string, error) {
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

func (v *PhysicalValidator) validateLabel(label string) (string, error) {
	// Empty label is allowed
	if label == "" {
		return "", nil
	}

	// Check label length
	length := utf8.RuneCountInString(label)
	if length > MaxLabelLength {
		return "", &validationErrors.ValidationError{
			Field:   "label",
			Message: fmt.Sprintf("label must be less than %d characters", MaxLabelLength),
		}
	}

	// Sanitize label
	sanitized, err := v.sanitizer.SanitizeString(label)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "label",
			Message: fmt.Sprintf("invalid label content: %v", err),
		}
	}

	return sanitized, nil
}

func (v *PhysicalValidator) validateLocationType(locationType string) (string, error) {
	// Check if location type is valid
	isValid := false
	for _, validType := range ValidLocationTypes {
		if locationType == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		return "", &validationErrors.ValidationError{
			Field:   "locationType",
			Message: fmt.Sprintf("location type must be one of: %s", strings.Join(ValidLocationTypes, ", ")),
		}
	}

	// Sanitize location type
	sanitized, err := v.sanitizer.SanitizeString(locationType)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "locationType",
			Message: fmt.Sprintf("invalid location type content: %v", err),
		}
	}

	return sanitized, nil
}

func (v *PhysicalValidator) validateMapCoordinates(coordinates string) (string, error) {
	// Empty coordinates are allowed
	if coordinates == "" {
		return "", nil
	}

	// Check format using regex
	re := regexp.MustCompile(CoordinatePattern)
	if !re.MatchString(coordinates) {
		return "", &validationErrors.ValidationError{
			Field:   "mapCoordinates",
			Message: "map coordinates must be in format latitude,longitude",
		}
	}

	// Parse and validate coordinate values
	parts := strings.Split(coordinates, ",")
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "mapCoordinates",
			Message: "invalid latitude value",
		}
	}

	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "mapCoordinates",
			Message: "invalid longitude value",
		}
	}

	// Check latitude range (-90 to 90)
	if lat < -90 || lat > 90 {
		return "", &validationErrors.ValidationError{
			Field:   "mapCoordinates",
			Message: "latitude must be between -90 and 90",
		}
	}

	// Check longitude range (-180 to 180)
	if lng < -180 || lng > 180 {
		return "", &validationErrors.ValidationError{
			Field:   "mapCoordinates",
			Message: "longitude must be between -180 and 180",
		}
	}

	// Sanitize coordinates
	sanitized, err := v.sanitizer.SanitizeString(coordinates)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "mapCoordinates",
			Message: fmt.Sprintf("invalid coordinates content: %v", err),
		}
	}

	return sanitized, nil
}

func (v *PhysicalValidator) validateBgColor(color string) (string, error) {
	// If color is empty, it's valid (optional field)
	if color == "" {
		return "", nil
	}

	// Convert to lowercase for case-insensitive comparison
	color = strings.ToLower(color)

	// Check if color is valid
	isValid := false
	for _, validColor := range ValidBgColors {
		if color == validColor {
			isValid = true
			break
		}
	}

	if !isValid {
		return "", &validationErrors.ValidationError{
			Field:   "bgColor",
			Message: fmt.Sprintf("background color must be one of: %s", strings.Join(ValidBgColors, ", ")),
		}
	}

	// Sanitize color
	sanitized, err := v.sanitizer.SanitizeString(color)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "bgColor",
			Message: fmt.Sprintf("invalid background color content: %v", err),
		}
	}

	return sanitized, nil
}