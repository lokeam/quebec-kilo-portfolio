package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type SublocationDbAdapter interface {
	GetUserSublocations(ctx context.Context, userID string) ([]models.Sublocation, error)
  GetSublocation(ctx context.Context, userID, locationID string) (models.Sublocation, error)
  AddSublocation(ctx context.Context, userID string, location models.Sublocation) (models.Sublocation, error)
  UpdateSublocation(ctx context.Context, userID string, location models.Sublocation) error
  RemoveSublocation(ctx context.Context, userID, locationID string) error
}
