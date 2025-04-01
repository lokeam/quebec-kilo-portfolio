package interfaces

import "github.com/lokeam/qko-beta/internal/models"

type PhysicalValidator interface {
	ValidatePhysicalLocation(location models.PhysicalLocation) (models.PhysicalLocation, error)
}