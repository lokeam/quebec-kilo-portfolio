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
func (lv *LibraryValidator) ValidateGame(game models.LibraryGame) error {
	// Initialize slice to collect validation errors
	var violations []string

	// Validate IGDB GameID
	if err := lv.ValidateGameID(game.GameID); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate GameName
	if err := lv.ValidateGameName(game.GameName); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate GameCoverURL
	if game.GameCoverURL != "" {
		if err := lv.validateURL(game.GameCoverURL); err != nil {
			violations = append(violations, err.Error())
		}
	}

	// Validate GameType
	if err := lv.validateGameType(game.GameType); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate CoverURL
	if game.GameCoverURL != "" {
		if err := lv.validateURL(game.GameCoverURL); err != nil {
			violations = append(violations, err.Error())
		}
	}

	// Validate Rating
	if err := lv.validateGameRating(game.GameRating); err != nil {
		violations = append(violations, err.Error())
}

	// Validate GameThemeNames
	if err := lv.validateGameThemeNames(game.GameThemeNames); err != nil {
		violations = append(violations, err.Error())
	}

	if err := lv.validateGameReleaseDate(game.GameFirstReleaseDate); err != nil {
		violations = append(violations, err.Error())
	}

	// Validate PlatformLocations
	if err := lv.validatePlatformLocations(game.PlatformLocations); err != nil {
		violations = append(violations, err.Error())
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
func (lv *LibraryValidator) ValidateGameID(id int64) error {
	if id <= 0 {
		return &validationErrors.ValidationError{
			Field: "id",
			Message: "game id must be positive",
		}
	}

	return nil
}

func (lv *LibraryValidator) ValidateGameName(name string) error {
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

func (lv *LibraryValidator) ValidateGameSummary(summary string) error {
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

func (lv *LibraryValidator) validateGameType(gameType models.LibraryGameType) error {
	if gameType.DisplayText == "" {
		return &validationErrors.ValidationError{
			Field: "gameType.displayText",
			Message: "game type display text cannot be empty",
		}
	}

	if gameType.NormalizedText == "" {
		return &validationErrors.ValidationError{
			Field: "gameType.normalizedText",
			Message: "game type normalized text cannot be empty",
		}
	}

	// Sanitize both fields
	sanitizedDisplayText, err := lv.sanitizer.SanitizeString(gameType.DisplayText)
	if err != nil {
		return &validationErrors.ValidationError{
			Field: "gameType.displayText",
			Message: fmt.Sprintf("invalid game type display text: %v", err),
		}
	}

	sanitizedNormalizedText, err := lv.sanitizer.SanitizeString(gameType.NormalizedText)
	if err != nil {
		return &validationErrors.ValidationError{
			Field: "gameType.normalizedText",
			Message: fmt.Sprintf("invalid game type normalized text: %v", err),
		}
	}

	if sanitizedDisplayText != gameType.DisplayText || sanitizedNormalizedText != gameType.NormalizedText {
		return &validationErrors.ValidationError{
			Field: "gameType",
			Message: "game type contains invalid characters",
		}
	}

	return nil
}

func (lv *LibraryValidator) validateGameThemeNames(themes []string) error {
	if len(themes) == 0 {
			return nil // Themes are optional
	}

	for i, theme := range themes {
			if theme == "" {
					return &validationErrors.ValidationError{
							Field: fmt.Sprintf("gameThemeNames[%d]", i),
							Message: "theme name cannot be empty",
					}
			}

			sanitizedTheme, err := lv.sanitizer.SanitizeString(theme)
			if err != nil {
					return &validationErrors.ValidationError{
							Field: fmt.Sprintf("gameThemeNames[%d]", i),
							Message: fmt.Sprintf("invalid theme name: %v", err),
					}
			}

			if sanitizedTheme != theme {
					return &validationErrors.ValidationError{
							Field: fmt.Sprintf("gameThemeNames[%d]", i),
							Message: "theme name contains invalid characters",
					}
			}
	}

	return nil
}

func (lv *LibraryValidator) validatePlatformLocations(locations []models.CreateLibraryGameLocation) error {
	if len(locations) == 0 {
			return &validationErrors.ValidationError{
					Field: "platformLocations",
					Message: "at least one platform location is required",
			}
	}

	for i, loc := range locations {
			// Validate PlatformName
			if loc.PlatformName == "" {
					return &validationErrors.ValidationError{
							Field: fmt.Sprintf("platformLocations[%d].platformName", i),
							Message: "platform name cannot be empty",
					}
			}

			// Validate Type
			if loc.Type != "digital" && loc.Type != "physical" {
					return &validationErrors.ValidationError{
							Field: fmt.Sprintf("platformLocations[%d].type", i),
							Message: "type must be either 'digital' or 'physical'",
					}
			}

			// Validate LocationID
			if loc.Location.SublocationID == "" && loc.Location.DigitalLocationID == "" {
					return &validationErrors.ValidationError{
							Field: fmt.Sprintf("platformLocations[%d].locationId", i),
							Message: "location ID cannot be empty",
					}
			}

			// Sanitize fields
			sanitizedPlatform, err := lv.sanitizer.SanitizeString(loc.PlatformName)
			if err != nil {
					return &validationErrors.ValidationError{
							Field: fmt.Sprintf("platformLocations[%d].platformName", i),
							Message: fmt.Sprintf("invalid platform name: %v", err),
					}
			}

			sanitizedSublocation, err := lv.sanitizer.SanitizeString(loc.Location.SublocationID)
			if err != nil {
					return &validationErrors.ValidationError{
							Field: fmt.Sprintf("platformLocations[%d].locationId", i),
							Message: fmt.Sprintf("invalid location ID: %v", err),
					}
			}

			sanitizedDigitalLocation, err := lv.sanitizer.SanitizeString(loc.Location.DigitalLocationID)
			if err != nil {
					return &validationErrors.ValidationError{
							Field: fmt.Sprintf("platformLocations[%d].locationId", i),
							Message: fmt.Sprintf("invalid location ID: %v", err),
					}
			}

			if sanitizedPlatform != loc.PlatformName || sanitizedSublocation != loc.Location.SublocationID || sanitizedDigitalLocation != loc.Location.DigitalLocationID {
					return &validationErrors.ValidationError{
							Field: fmt.Sprintf("platformLocations[%d]", i),
							Message: "platform location contains invalid characters",
					}
			}
	}

	return nil
}

func (lv *LibraryValidator) validateGameRating(rating float64) error {
	if rating < 0 || rating > 100 {
		return &validationErrors.ValidationError{
			Field: "gameRating",
			Message: "rating must be between 0 and 100",
		}
	}

	return nil
}

func (lv *LibraryValidator) validateGameReleaseDate(releaseDate int64) error {
	if releaseDate < 0 {
		return &validationErrors.ValidationError{
			Field: "gameFirstReleaseDate",
			Message: "release date cannot be negative",
		}
	}

	return nil
}