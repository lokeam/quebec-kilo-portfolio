package interfaces

import "github.com/lokeam/qko-beta/internal/models"

type PhysicalValidator interface {
	ValidatePhysicalLocation(location models.PhysicalLocation) (models.PhysicalLocation, error)
	ValidatePhysicalLocationCreation(location models.PhysicalLocation) (models.PhysicalLocation, error)
	ValidatePhysicalLocationUpdate(update, existing models.PhysicalLocation) (models.PhysicalLocation, error)
	ValidateRemovePhysicalLocation(userID string, locationIDs []string) ([]string, error)
}