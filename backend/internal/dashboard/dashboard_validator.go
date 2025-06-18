package dashboard

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

type DashboardValidatorImpl struct{}

func NewDashboardValidator() interfaces.DashboardValidator {
	return &DashboardValidatorImpl{}
}

func (v *DashboardValidatorImpl) ValidateUserID(userID string) error {
	if userID == "" {
		return errors.New("userID is required")
	}
	return nil
}