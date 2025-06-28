package interfaces

import "github.com/lokeam/qko-beta/internal/types"

// Todo: add validate onetime purchase
type SpendTrackingValidator interface {
	ValidateUserID(userID string) error
	ValidateOneTimePurchase(request types.SpendTrackingRequest) error
	ValidateDeleteOneSpendTrackingItems(userID string, itemIDs []string) ([]string, error)
}