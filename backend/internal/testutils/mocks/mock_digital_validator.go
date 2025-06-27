package mocks

import "github.com/lokeam/qko-beta/internal/models"

type MockDigitalValidator struct {
	ValidateDigitalLocationFunc func(digitalLocation models.DigitalLocation) (models.DigitalLocation, error)
	ValidateDigitalLocationsBulkFunc func(locations []models.DigitalLocation) ([]models.DigitalLocation, error)
}

func (m *MockDigitalValidator) ValidateDigitalLocation(
	digitalLocation models.DigitalLocation,
) (models.DigitalLocation, error) {
	if m.ValidateDigitalLocationFunc != nil {
		return m.ValidateDigitalLocationFunc(digitalLocation)
	}
	return digitalLocation, nil
}

func (m *MockDigitalValidator) ValidateDigitalLocationsBulk(
	locations []models.DigitalLocation,
) ([]models.DigitalLocation, error) {
	return m.ValidateDigitalLocationsBulkFunc(locations)
}