package library

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	validationErrors "github.com/lokeam/qko-beta/internal/shared/validation"
)

// Validation constants
// TODO: optimize these values for current use case
const (
	MinNameLength      = 1
	MaxNameLength      = 200
	MaxSummaryLength   = 5000
	MinUserIDLength    = 5
	MaxUserIDLength    = 128
	MaxURLLength       = 2048
)

type LibraryValidator struct {
	sanitizer interfaces.Sanitizer
}

// Constructor
func NewLibraryValidator(sanitizer interfaces.Sanitizer) (*LibraryValidator, error) {
	if sanitizer == nil {
		return nil, fmt.Errorf("sanitizer is required")
	}

	return &LibraryValidator{
		sanitizer: sanitizer,
	}, nil
}

// Validate userID for library ops
func (lv *LibraryValidator) ValidateUserID(userID string) error {

	// Check if userID is empty
	if userID == "" {
		return &validationErrors.ValidationError{
			Field: "userID",
			Message: "userID cannot be empty",
		}
	}

	length := utf8.RuneCountInString(userID)
	if length < MinUserIDLength || length > MaxUserIDLength {
		return &validationErrors.ValidationError{
			Field: "userID",
			Message: fmt.Sprintf("userID must be between %d and %d characters", MinUserIDLength, MaxUserIDLength),
		}
	}

	// Sanitize userID
	sanitizedUserID, err := lv.sanitizer.SanitizeString(userID)
	if err != nil {
		return &validationErrors.ValidationError{
			Field:   "userID",
			Message: fmt.Sprintf("invalid userID: %v", err),
		}
	}

	if sanitizedUserID != userID {
		return &validationErrors.ValidationError{
			Field:   "userID",
			Message: "user ID contains invalid characters",
		}
	}

	return nil
}

// Validate game object for adding to library
func (lv *LibraryValidator) ValidateGame(game models.Game) error {
	// Initialize slice to collect validation errors
	var violations []string

	// Validate ID
	if err := lv.validateGameID(game.ID); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate Name
	if err := lv.validateGameName(game.Name); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate Summary (optional)
	if game.Summary != "" {
		if err := lv.validateGameSummary(game.Summary); err != nil {
			violations = append(violations, err.Error())
		}
	}

	// Validate CoverURL (optional)
	if game.CoverURL != "" {
		if err := lv.validateURL(game.CoverURL); err != nil {
			violations = append(violations, err.Error())
		}
	}

	// Validate Rating (optional)
	if game.Rating > 0 {
		if err := lv.validateRating(game.Rating); err != nil {
			violations = append(violations, err.Error())
		}
	}

	if len(violations) > 0 {
		return &validationErrors.ValidationError{
			Field: "game",
			Message: fmt.Sprintf("game validation failed %v", violations),
		}
	}

	return nil
}


// Helpers
func (lv *LibraryValidator) validateGameID(id int64) error {
	if id <= 0 {
		return &validationErrors.ValidationError{
			Field: "id",
			Message: "game id must be positive",
		}
	}

	return nil
}

func (lv *LibraryValidator) validateGameName(name string) error {
	if name == "" {
		return &validationErrors.ValidationError{
			Field: "name",
			Message: "game name cannot be empty",
		}
	}

	length := utf8.RuneCountInString(name)
	if length < MinNameLength || length > MaxNameLength {
		return &validationErrors.ValidationError{
			Field: "name",
			Message: fmt.Sprintf("game name must be between %d and %d characters", MinNameLength, MaxNameLength),
		}
	}

	sanitizedName, err := lv.sanitizer.SanitizeString(name)
	if err != nil {
		return &validationErrors.ValidationError{
			Field: "name",
			Message: fmt.Sprintf("invalid game name: %v", err),
		}
	}

	if sanitizedName != name {
		return &validationErrors.ValidationError{
			Field:     "name",
			Message:   "game name contains invalid characters",
		}
	}

	return nil
}

func (lv *LibraryValidator) validateGameSummary(summary string) error {
	length := utf8.RuneCountInString(summary)
	if length > MaxSummaryLength {
		return &validationErrors.ValidationError{
			Field: "summary",
			Message: fmt.Sprintf("game summary must be less than %d characters", MaxSummaryLength),
		}
	}

	sanitizedSummary, err := lv.sanitizer.SanitizeString(summary)
	if err != nil {
		return &validationErrors.ValidationError{
			Field: "summary",
			Message: fmt.Sprintf("invalid game summary: %v", err),
		}
	}

	if sanitizedSummary != summary {
		return &validationErrors.ValidationError{
			Field:    "summary",
			Message: "game summary contains invalid characters",
		}
	}

	return nil
}

func (lv *LibraryValidator) validateURL (url string) error {
	if url == "" {
		return nil // Note: URL is optional
	}

	length := utf8.RuneCountInString(url)
	if length > MaxURLLength {
		return &validationErrors.ValidationError{
			Field: "coverURL",
			Message: fmt.Sprintf("URL must be less than %d characters", MaxURLLength),
		}
	}

	// Basic URL validation pattern
	// Note: Is there a better way to do this?
	urlPattern := regexp.MustCompile(`^https?://\S+$`)
	if !urlPattern.MatchString(url) {
		return &validationErrors.ValidationError{
			Field:    "coverURL",
			Message:  "invalid URL format",
		}
	}

	return nil
}

func (lv *LibraryValidator) validateRating(rating float64) error {
	if rating < 0 || rating > 100 {
		return &validationErrors.ValidationError{
			Field: "rating",
			Message: "rating must be between 0 and 100",
		}
	}

	return nil
}
