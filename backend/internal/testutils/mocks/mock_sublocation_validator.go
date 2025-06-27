package mocks

import "github.com/lokeam/qko-beta/internal/models"

type MockSublocationValidator struct {
	ValidateSublocationFunc func(sublocation models.Sublocation) (models.Sublocation, error)
	ValidateSublocationUpdateFunc func(update, existing models.Sublocation) (models.Sublocation, error)
	ValidateSublocationCreationFunc func(sublocation models.Sublocation) (models.Sublocation, error)
	ValidateGameOwnershipFunc func(userID string, userGameID string) error
	ValidateSublocationOwnershipFunc func(userID string, sublocationID string) error
	ValidateGameNotInSublocationFunc func(userGameID string, sublocationID string) error
	ValidateDeleteSublocationRequestFunc func(userID string, sublocationIDs []string) error
}

func (m *MockSublocationValidator) ValidateSublocation(
	sublocation models.Sublocation,
) (models.Sublocation, error) {
	if m.ValidateSublocationFunc != nil {
		return m.ValidateSublocationFunc(sublocation)
	}
	return sublocation, nil
}

func (m *MockSublocationValidator) ValidateSublocationUpdate(
	update, existing models.Sublocation,
) (models.Sublocation, error) {
	if m.ValidateSublocationUpdateFunc != nil {
		return m.ValidateSublocationUpdateFunc(update, existing)
	}
	return update, nil
}

func (m *MockSublocationValidator) ValidateSublocationCreation(
	sublocation models.Sublocation,
) (models.Sublocation, error) {
	if m.ValidateSublocationCreationFunc != nil {
		return m.ValidateSublocationCreationFunc(sublocation)
	}
	return sublocation, nil
}

func (m *MockSublocationValidator) ValidateGameOwnership(
	userID string,
	userGameID string,
) error {
	if m.ValidateGameOwnershipFunc != nil {
		return m.ValidateGameOwnershipFunc(userID, userGameID)
	}
	return nil
}

func (m *MockSublocationValidator) ValidateSublocationOwnership(
	userID string,
	sublocationID string,
) error {
	if m.ValidateSublocationOwnershipFunc != nil {
		return m.ValidateSublocationOwnershipFunc(userID, sublocationID)
	}
	return nil
}

func (m *MockSublocationValidator) ValidateGameNotInSublocation(
	userGameID string,
	sublocationID string,
) error {
	if m.ValidateGameNotInSublocationFunc != nil {
		return m.ValidateGameNotInSublocationFunc(userGameID, sublocationID)
	}
	return nil
}

func (m *MockSublocationValidator) ValidateDeleteSublocationRequest(
	userID string,
	sublocationIDs []string,
) error {
	if m.ValidateDeleteSublocationRequestFunc != nil {
		return m.ValidateDeleteSublocationRequestFunc(userID, sublocationIDs)
	}
	return nil
}