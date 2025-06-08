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

// SublocationValidator implements the interfaces.SublocationValidator interface
type SublocationValidator struct {
	sanitizer interfaces.Sanitizer
	cacheWrapper interfaces.SublocationCacheWrapper
	logger interfaces.Logger
}

// Ensure SublocationValidator implements interfaces.SublocationValidator
var _ interfaces.SublocationValidator = (*SublocationValidator)(nil)

// NewSublocationValidator creates a new sublocation validator
func NewSublocationValidator(
	sanitizer interfaces.Sanitizer,
	cacheWrapper interfaces.SublocationCacheWrapper,
	logger interfaces.Logger,
) (*SublocationValidator, error) {
	if sanitizer == nil {
		return nil, fmt.Errorf("sanitizer cannot be nil")
	}

	if cacheWrapper == nil {
		return nil, fmt.Errorf("cacheWrapper cannot be nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	return &SublocationValidator{
		sanitizer: sanitizer,
		cacheWrapper: cacheWrapper,
		logger: logger,
	}, nil
}

// ValidateSublocation validates a sublocation for creation
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

	// Validate stored items
	if err := sv.validateStoredItems(sublocation.StoredItems); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate physical location ID
	if err := sv.validatePhysicalLocationID(sublocation.PhysicalLocationID); err != nil {
		violations = append(violations, err.Error())
	}

	if len(violations) > 0 {
		return models.Sublocation{}, fmt.Errorf("validation failed: %s", strings.Join(violations, ", "))
	}

	return sublocation, nil
}

// ValidateSublocationUpdate validates a sublocation update by only checking fields that have changed
func (sv *SublocationValidator) ValidateSublocationUpdate(update, existing models.Sublocation) (models.Sublocation, error) {
	// Start with the existing sublocation
	validated := existing

	// Only validate and update fields that have changed
	if update.Name != "" && update.Name != existing.Name {
		// Check for duplicate name only if the name is being changed
		cachedSublocations, err := sv.cacheWrapper.GetCachedSublocations(context.Background(), update.UserID)
		if err != nil {
				sv.logger.Warn("Cache error during duplicate name check", map[string]any{
						"error": err,
						"userID": update.UserID,
				})
		} else {
			for i := 0; i < len(cachedSublocations); i++ {
				if cachedSublocations[i].PhysicalLocationID == update.PhysicalLocationID &&
						cachedSublocations[i].ID != update.ID && // Don't compare with self
						strings.EqualFold(cachedSublocations[i].Name, update.Name) {

						return models.Sublocation{}, &validationErrors.ValidationError{
								Field: "name",
								Message: "a sublocation with this name already exists in this physical location",
						}
				}
			}
		}
		validated.Name = update.Name
	}

	if update.LocationType != "" {
		validated.LocationType = update.LocationType
	}

	if update.StoredItems != existing.StoredItems {
		validated.StoredItems = update.StoredItems
	}

	return validated, nil
}

// ValidateSublocationCreation validates a sublocation for creation (duplicate name check)
func (sv *SublocationValidator) ValidateSublocationCreation(sublocation models.Sublocation) (models.Sublocation, error) {
	// First sanitize the name
	sanitizedName, err := sv.sanitizer.SanitizeString(sublocation.Name)
	if err != nil {
		return models.Sublocation{}, &validationErrors.ValidationError{
			Field: "name",
			Message: fmt.Sprintf("invalid name content: %v", err),
		}
	}

	// Create a copy of the sublocation with sanitized for duplicate checking
	sublocationToCheck := sublocation
	sublocationToCheck.Name = sanitizedName

	// Check for duplicate name within the physical location
	cachedSublocations, err := sv.cacheWrapper.GetCachedSublocations(
		context.Background(),
		sublocation.UserID,
	)
	if err != nil {
		sv.logger.Warn("Cache error during duplicate name check", map[string]any{
			"error": err,
			"userID": sublocation.UserID,
		})
	} else {
		for i := 0; i < len(cachedSublocations); i++ {
			if cachedSublocations[i].PhysicalLocationID == sublocation.PhysicalLocationID &&
				 strings.EqualFold(cachedSublocations[i].Name, sanitizedName) {

					return models.Sublocation{}, &validationErrors.ValidationError{
							Field: "name",
							Message: "a sublocation with this name already exists in this physical location",
					}
			}
		}
	}

	// Perform regular validation
	return sv.ValidateSublocation(sublocationToCheck)
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

// ValidateGameOwnership validates the format of user ID and user game ID
func (sv *SublocationValidator) ValidateGameOwnership(userID string, userGameID string) error {
	// Sanitize and validate user ID
	sanitizedUserID, err := sv.sanitizer.SanitizeString(userID)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "user_id",
			Message: fmt.Sprintf("invalid user_id content: %v", err),
		}
	}

	// Validate user ID is not empty
	if sanitizedUserID == "" {
		return &validationErrors.ValidationError{
			Field:   "user_id",
			Message: "user_id cannot be empty",
		}
	}

	// Sanitize and validate user game ID
	sanitizedUserGameID, err := sv.sanitizer.SanitizeString(userGameID)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "user_game_id",
			Message: fmt.Sprintf("invalid user_game_id content: %v", err),
		}
	}

	// Validate user game ID is not empty
	if sanitizedUserGameID == "" {
		return &validationErrors.ValidationError{
			Field:   "user_game_id",
			Message: "user_game_id cannot be empty",
		}
	}

	return nil
}

// ValidateSublocationOwnership validates the format of user ID and sublocation ID
func (sv *SublocationValidator) ValidateSublocationOwnership(userID string, sublocationID string) error {
	// Sanitize and validate user ID
	sanitizedUserID, err := sv.sanitizer.SanitizeString(userID)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "user_id",
			Message: fmt.Sprintf("invalid user_id content: %v", err),
		}
	}

	// Validate user ID is not empty
	if sanitizedUserID == "" {
		return &validationErrors.ValidationError{
			Field:   "user_id",
			Message: "user_id cannot be empty",
		}
	}

	// Sanitize and validate sublocation ID
	sanitizedSublocationID, err := sv.sanitizer.SanitizeString(sublocationID)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "sublocation_id",
			Message: fmt.Sprintf("invalid sublocation_id content: %v", err),
		}
	}

	// Validate sublocation ID is not empty
	if sanitizedSublocationID == "" {
		return &validationErrors.ValidationError{
			Field:   "sublocation_id",
			Message: "sublocation_id cannot be empty",
		}
	}

	return nil
}

// ValidateGameNotInSublocation validates the format of user game ID and sublocation ID
func (sv *SublocationValidator) ValidateGameNotInSublocation(userGameID string, sublocationID string) error {
	// Sanitize and validate user game ID
	sanitizedUserGameID, err := sv.sanitizer.SanitizeString(userGameID)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "user_game_id",
			Message: fmt.Sprintf("invalid user_game_id content: %v", err),
		}
	}

	// Validate user game ID is not empty
	if sanitizedUserGameID == "" {
		return &validationErrors.ValidationError{
			Field:   "user_game_id",
			Message: "user_game_id cannot be empty",
		}
	}

	// Sanitize and validate sublocation ID
	sanitizedSublocationID, err := sv.sanitizer.SanitizeString(sublocationID)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "sublocation_id",
			Message: fmt.Sprintf("invalid sublocation_id content: %v", err),
		}
	}

	// Validate sublocation ID is not empty
	if sanitizedSublocationID == "" {
		return &validationErrors.ValidationError{
			Field:   "sublocation_id",
			Message: "sublocation_id cannot be empty",
		}
	}

	return nil
}