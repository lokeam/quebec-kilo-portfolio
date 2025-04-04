package interfaces

import "github.com/lokeam/qko-beta/internal/models"

type DigitalValidator interface {
	ValidateDigitalLocation(location models.DigitalLocation) (models.DigitalLocation, error)
}
