package mocks

import "github.com/lokeam/qko-beta/internal/models"


type MockPhysicalValidator struct {
	ValidatePhysicalLocationFunc func(location models.PhysicalLocation) (models.PhysicalLocation, error)
	ValidatePhysicalLocationCreationFunc func(location models.PhysicalLocation) (models.PhysicalLocation, error)
	ValidatePhysicalLocationUpdateFunc func(update, existing models.PhysicalLocation) (models.PhysicalLocation, error)
	ValidateRemovePhysicalLocationFunc func(userID string, locationIDs []string) ([]string, error)
}

func (m *MockPhysicalValidator) ValidatePhysicalLocation(location models.PhysicalLocation) (models.PhysicalLocation, error) {
	if m.ValidatePhysicalLocationFunc != nil {
		return m.ValidatePhysicalLocationFunc(location)
	}
	return location, nil
}

func (m *MockPhysicalValidator) ValidatePhysicalLocationCreation(location models.PhysicalLocation) (models.PhysicalLocation, error) {
	if m.ValidatePhysicalLocationCreationFunc != nil {
		return m.ValidatePhysicalLocationCreationFunc(location)
	}
	return location, nil
}

func (m *MockPhysicalValidator) ValidatePhysicalLocationUpdate(update, existing models.PhysicalLocation) (models.PhysicalLocation, error) {
	if m.ValidatePhysicalLocationUpdateFunc != nil {
		return m.ValidatePhysicalLocationUpdateFunc(update, existing)
	}
	return update, nil
}

func (m *MockPhysicalValidator) ValidateRemovePhysicalLocation(userID string, locationIDs []string) ([]string, error) {
	if m.ValidateRemovePhysicalLocationFunc != nil {
		return m.ValidateRemovePhysicalLocationFunc(userID, locationIDs)
	}
	return locationIDs, nil
}