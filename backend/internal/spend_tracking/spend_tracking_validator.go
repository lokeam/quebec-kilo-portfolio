package spend_tracking

import (
	"errors"
	"fmt"

	"github.com/lokeam/qko-beta/internal/interfaces"
)

type ValidationError struct {
	Field string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type SpendTrackingValidatorImpl struct{}

func NewSpendTrackingValidator() interfaces.SpendTrackingValidator {
	return &SpendTrackingValidatorImpl{}
}

func (v *SpendTrackingValidatorImpl) ValidateUserID(userID string) error {
	if userID == "" {
		return errors.New("userID is required")
	}
	return nil
}