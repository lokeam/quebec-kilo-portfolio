package interfaces

import "github.com/lokeam/qko-beta/internal/models"

type SublocationValidator interface {
	ValidateSublocation(sublocation models.Sublocation) (models.Sublocation, error)
	ValidateSublocationUpdate(update, existing models.Sublocation) (models.Sublocation, error)
}
