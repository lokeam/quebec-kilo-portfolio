package spend_tracking

import (
	"errors"
	"fmt"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/types"
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

func (v *SpendTrackingValidatorImpl) ValidateOneTimePurchase(request types.SpendTrackingRequest) error {
	if request.Title == "" {
		return &ValidationError{Field: "title", Message: "title is required"}
	}
	if request.Amount <= 0 {
		return &ValidationError{Field: "amount", Message: "amount must be greater than 0"}
	}
	if request.PurchaseDate == "" {
		return &ValidationError{Field: "purchase_date", Message: "purchase_date is required"}
	}
	if request.PaymentMethod == "" {
		return &ValidationError{Field: "payment_method", Message: "payment_method is required"}
	}
	if request.SpendingCategoryID <= 0 {
		return &ValidationError{Field: "spending_category_id", Message: "spending_category_id is required"}
	}
	return nil
}