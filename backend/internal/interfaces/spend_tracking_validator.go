package interfaces

// Todo: add validate onetime purchase
type SpendTrackingValidator interface {
	ValidateUserID(userID string) error
}