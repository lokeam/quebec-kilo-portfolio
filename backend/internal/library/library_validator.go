package library

import (
	"errors"
	"fmt"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type LibraryValidatorImpl struct{}

func NewLibraryValidator() interfaces.LibraryValidator {
	return &LibraryValidatorImpl{}
}

func (v *LibraryValidatorImpl) ValidateLibraryGame(game models.GameToSave) error {
	if game.GameID <= 0 {
		return errors.New("game ID must be positive")
	}

	if game.GameName == "" {
		return errors.New("game name is required")
	}

	if len(game.PlatformLocations) == 0 {
		return errors.New("at least one platform location is required")
	}

	for i, location := range game.PlatformLocations {
		if location.PlatformID <= 0 {
			return fmt.Errorf("platform ID must be positive at index %d", i)
		}

		if location.PlatformName == "" {
			return fmt.Errorf("platform name is required at index %d", i)
		}

		if location.Type != "physical" && location.Type != "digital" {
			return fmt.Errorf("invalid location type '%s' at index %d", location.Type, i)
		}

		if location.Type == "physical" && location.Location.SublocationID == "" {
			return fmt.Errorf("sublocation ID is required for physical location at index %d", i)
		}

		if location.Type == "digital" && location.Location.DigitalLocationID == "" {
			return fmt.Errorf("digital location ID is required for digital location at index %d", i)
		}
	}

	return nil
}

func (v *LibraryValidatorImpl) ValidateUserID(userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}
	return nil
}

func (v *LibraryValidatorImpl) ValidateGameID(gameID int64) error {
	if gameID <= 0 {
		return errors.New("game ID must be positive")
	}
	return nil
}
