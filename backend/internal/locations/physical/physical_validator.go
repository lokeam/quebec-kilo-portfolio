package physical

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	validationErrors "github.com/lokeam/qko-beta/internal/shared/validation"
)

// Validation constants
const (
	MaxNameLength     = 100
	MinNameLength     = 1
	MaxLabelLength    = 50
	CoordinatePattern = `^-?\d+(\.\d+)?,\s*-?\d+(\.\d+)?$`
)

// Valid location types
var ValidLocationTypes = []string{"house", "apartment", "office", "warehouse"}

// ValidBgColors lists all valid background colors
var ValidBgColors = []string{
	"red",
	"blue",
	"green",
	"gold",
	"purple",
	"orange",
	"brown",
	"pink",
	"gray",
}

// PhysicalValidator struct
type PhysicalValidator struct {
	sanitizer interfaces.Sanitizer
	cacheWrapper interfaces.PhysicalCacheWrapper
	logger interfaces.Logger
}

func NewPhysicalValidator(
	sanitizer interfaces.Sanitizer,
	cacheWrapper interfaces.PhysicalCacheWrapper,
	logger interfaces.Logger,
) (*PhysicalValidator, error) {
	if sanitizer == nil {
		return nil, fmt.Errorf("sanitizer cannot be nil")
	}
	if cacheWrapper == nil {
		return nil, fmt.Errorf("cacheWrapper cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	return &PhysicalValidator{
		sanitizer: sanitizer,
		cacheWrapper: cacheWrapper,
		logger: logger,
	}, nil
}

// ValidatePhysicalLocation validates a physical location for updates
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
	if err := v.validateBackgroundColor(location.BgColor); err != nil {
		violations = append(violations, err.Error())
	}

	// Copy other fields that don't need validation
	validatedLocation.ID = location.ID
	validatedLocation.UserID = location.UserID
	validatedLocation.CreatedAt = location.CreatedAt
	validatedLocation.UpdatedAt = location.UpdatedAt

	if len(violations) > 0 {
		return models.PhysicalLocation{}, &validationErrors.ValidationError{
			Field:   "location",
			Message: fmt.Sprintf("Physical location validation failed: %s", strings.Join(violations, "; ")),
		}
	}

	return validatedLocation, nil
}

// ValidatePhysicalLocationCreation validates a physical location for creation, including duplicate name check
func (v *PhysicalValidator) ValidatePhysicalLocationCreation(location models.PhysicalLocation) (models.PhysicalLocation, error) {
	// First sanitize the name
	sanitizedName, err := v.sanitizer.SanitizeString(location.Name)
	if err != nil {
		return models.PhysicalLocation{}, &validationErrors.ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("invalid name content: %v", err),
		}
	}

	// Create a copy of the location with sanitized name for duplicate checking
	locationToCheck := location
	locationToCheck.Name = sanitizedName

	// Check for duplicate name
	cachedLocations, err := v.cacheWrapper.GetCachedPhysicalLocations(context.Background(), location.UserID)
	if err != nil {
		// Log the cache error but continue with validation
		// This ensures we don't block creation if cache is temporarily unavailable
		v.logger.Warn("Cache error during duplicate name check", map[string]any{
			"error": err,
			"userID": location.UserID,
		})
	} else {
		// Only check for duplicates if we successfully got cached locations
		for i := 0; i < len(cachedLocations); i++ {
			if strings.EqualFold(cachedLocations[i].Name, sanitizedName) {
				return models.PhysicalLocation{}, &validationErrors.ValidationError{
					Field: "name",
					Message: "a location with this name already exists",
				}
			}
		}
	}

	// Then perform regular validation with the sanitized name
	return v.ValidatePhysicalLocation(locationToCheck)
}

// ValidatePhysicalLocationUpdate validates a physical location update by only checking fields that have changed
func (v *PhysicalValidator) ValidatePhysicalLocationUpdate(update, existing models.PhysicalLocation) (models.PhysicalLocation, error) {
	// Start with the existing location
	validated := existing

	// Only validate and update fields that have changed
	if update.Name != "" && update.Name != existing.Name {
			// First sanitize the new name
			sanitizedName, err := v.sanitizer.SanitizeString(update.Name)
			if err != nil {

				// Log the sanitization error
				v.logger.Error("Failed to sanitize name", map[string]any{
					"error": err,
					"name": update.Name,
				})

				// Return validation error with sanitization error message
				return models.PhysicalLocation{}, &validationErrors.ValidationError{
						Field:   "name",
						Message: fmt.Sprintf("invalid name content: %v", err),
				}
			}

			// Check for duplicate name only if the name is being changed
			cachedLocations, err := v.cacheWrapper.GetCachedPhysicalLocations(context.Background(), update.UserID)
			if err != nil {
				// Log the cache error but continue with validation
				v.logger.Warn("Cache error during duplicate name check", map[string]any{
						"error": err,
						"userID": update.UserID,
				})
			} else {
				for i := 0; i < len(cachedLocations); i++ {
						if cachedLocations[i].ID != update.ID && // Don't compare with self
								strings.EqualFold(cachedLocations[i].Name, sanitizedName) {
								return models.PhysicalLocation{}, &validationErrors.ValidationError{
										Field: "name",
										Message: "a location with this name already exists",
								}
						}
				}
			}
			validated.Name = sanitizedName
	}

	// Validate label if changed
	if update.Label != existing.Label {
			if sanitizedLabel, err := v.validateLabel(update.Label); err != nil {
					return models.PhysicalLocation{}, err
			} else {
					validated.Label = sanitizedLabel
			}
	}

	// Validate location type if changed
	if update.LocationType != existing.LocationType {
			if sanitizedType, err := v.validateLocationType(update.LocationType); err != nil {
					return models.PhysicalLocation{}, err
			} else {
					validated.LocationType = sanitizedType
			}
	}

	// Validate map coordinates if changed
	if update.MapCoordinates != existing.MapCoordinates {
			if sanitizedCoords, err := v.validateMapCoordinates(update.MapCoordinates); err != nil {
					return models.PhysicalLocation{}, err
			} else {
					validated.MapCoordinates = sanitizedCoords
			}
	}

	// Validate background color if changed
	if update.BgColor != existing.BgColor {
			if err := v.validateBackgroundColor(update.BgColor); err != nil {
					return models.PhysicalLocation{}, err
			}
			validated.BgColor = update.BgColor
	}

	// Update the timestamp
	validated.UpdatedAt = time.Now()

	return validated, nil
}


// Individual validation rules
func (v *PhysicalValidator) validateName(name string) (string, error) {
	// Check if name is empty
	if name == "" {
		return "", &validationErrors.ValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}

	// Check name length
	length := utf8.RuneCountInString(name)
	if length < MinNameLength {
		return "", &validationErrors.ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("name must be at least %d characters", MinNameLength),
		}
	}
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
	// Check if location type is empty
	if locationType == "" {
		return "", &validationErrors.ValidationError{
			Field:   "locationType",
			Message: "location type is required",
		}
	}

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
			Message: "map coordinates must be in format latitude,longitude (e.g. 45.5017,-73.5673)",
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

func (v *PhysicalValidator) validateBackgroundColor(bgColor string) error {
	// First sanitize the background color
	sanitized, err := v.sanitizer.SanitizeString(bgColor)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "bgColor",
			Message: fmt.Sprintf("invalid background color content: %v", err),
		}
	}

	// Next check if background color is valid
	isValid := false
	for _, validColor := range ValidBgColors {
		if strings.EqualFold(sanitized, validColor) {
			isValid = true
			break
		}
	}

	if !isValid {
		return &validationErrors.ValidationError{
			Field:     "bgColor",
			Message:   fmt.Sprintf("background color must be one of %v", ValidBgColors),
		}
	}

	return nil
}