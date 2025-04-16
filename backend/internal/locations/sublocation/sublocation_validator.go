package sublocation

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	validationErrors "github.com/lokeam/qko-beta/internal/shared/validation"
)

// Validation constants
const (
	MaxNameLength = 100
	MinNameLength = 1
	MinStoredItems   = 0
	MaxStoredItems   = 1000
)

// ValidSublocationTypes lists all valid sublocation types
var ValidSublocationTypes = []string{
	"shelf",
	"console",
	"cabinet",
	"drawer",
	"box",
	"closet",
}

// ValidBgColors lists all valid background colors
var ValidBgColors = []string{
	"red",
	"blue",
	"green",
	"gold",
	"purple",
	"orange",
	"brown",
	"white",
	"gray",
}

// SublocationValidator handles validation of sublocation data
type SublocationValidator struct {
	sanitizer interfaces.Sanitizer
	dbAdapter interfaces.SublocationDbAdapter
}

// NewSublocationValidator creates a new sublocation validator
func NewSublocationValidator(sanitizer interfaces.Sanitizer, dbAdapter interfaces.SublocationDbAdapter) (*SublocationValidator, error) {
	if sanitizer == nil {
		return nil, fmt.Errorf("sanitizer cannot be nil")
	}

	if dbAdapter == nil {
		return nil, fmt.Errorf("dbAdapter cannot be nil")
	}

	return &SublocationValidator{
		sanitizer: sanitizer,
		dbAdapter: dbAdapter,
	}, nil
}

// Validation Sublocation validates a sublocation model
func (sv *SublocationValidator) ValidateSublocation(sublocation models.Sublocation) (models.Sublocation, error) {
	var violations []string

	// Validate name
	if err := sv.validateName(sublocation.Name); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate location type
	if err := sv.validateLocationType(sublocation.LocationType); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate background color
	if err := sv.validateBackgroundColor(sublocation.BgColor); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate stored items
	if err := sv.validateStoredItems(sublocation.StoredItems); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate physical location ID
	if err := sv.validatePhysicalLocationID(sublocation.PhysicalLocationID); err != nil {
		violations = append(violations, err.Error())
	}

	// Check for duplicate sublocation
	if exists, err := sv.dbAdapter.CheckDuplicateSublocation(
		context.Background(),
		sublocation.UserID,
		sublocation.PhysicalLocationID,
		sublocation.Name,
	); err != nil {
		violations = append(violations, "error checking for duplicate sublocation")
	} else if exists {
		violations = append(violations, "a sublocation with this name already exists in this physical location")
	}

	if len(violations) > 0 {
		return models.Sublocation{}, fmt.Errorf("validation failed: %s", strings.Join(violations, ", "))
	}

	return sublocation, nil
}

func (v *SublocationValidator) validateName(name string) error {
	// First sanitize the name
	sanitized, err := v.sanitizer.SanitizeString(name)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("invalid name content: %v", err),
		}
	}

	// Check name length
	if len(sanitized) < MinNameLength {
		return &validationErrors.ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("name must be at least %d characters long", MinNameLength),
		}
	}

	if len(sanitized) > MaxNameLength {
		return &validationErrors.ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("name must not exceed %d characters", MaxNameLength),
		}
	}

	return nil
}

func (v *SublocationValidator) validateLocationType(locationType string) error {
	// Check if location type is valid
	isValid := false
	for _, validType := range ValidSublocationTypes {
		if locationType == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		return &validationErrors.ValidationError{
			Field:   "locationType",
			Message: fmt.Sprintf("location type must be one of %v", ValidSublocationTypes),
		}
	}

	return nil
}

func (v *SublocationValidator) validateBackgroundColor(bgColor string) error {
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

func (v *SublocationValidator) validateStoredItems(storedItems int) error {
	// Check if stored items is negative
	if storedItems < 0 {
		return &validationErrors.ValidationError{
			Field:   "stored_items",
			Message: "stored_items cannot be negative",
		}
	}

	// Check if stored items is within reasonable limits
	if storedItems > MaxStoredItems {
		return &validationErrors.ValidationError{
			Field:   "stored_items",
			Message: fmt.Sprintf("stored_items must not exceed %d", MaxStoredItems),
		}
	}

	return nil
}

func (v *SublocationValidator) validatePhysicalLocationID(physicalLocationID string) error {
	// Check if physical location ID is empty
	if physicalLocationID == "" {
		return &validationErrors.ValidationError{
			Field:   "physical_location_id",
			Message: "physical_location_id cannot be empty",
		}
	}

	// Sanitize physical location ID
	sanitized, err := v.sanitizer.SanitizeString(physicalLocationID)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "physical_location_id",
			Message: fmt.Sprintf("invalid physical_location_id content: %v", err),
		}
	}

	// Validate UUID format
	_, err = uuid.Parse(sanitized)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "physical_location_id",
			Message: "physical_location_id must be a valid UUID",
		}
	}

	return nil
}