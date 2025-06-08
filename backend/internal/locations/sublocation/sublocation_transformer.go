package sublocation

import (
	"time"

	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

// TransformCreateRequestToModel converts a CreateSublocationRequest to a models.Sublocation
// This is a pure function that only handles the transformation of data.
// No validation, error handling, or business logic should be included here.
func TransformCreateRequestToModel(req types.CreateSublocationRequest, userID string) models.Sublocation {
	return models.Sublocation{
		ID:                 uuid.New().String(),
		UserID:             userID,
		PhysicalLocationID: req.PhysicalLocationID,
		Name:               req.Name,
		LocationType:       req.LocationType,
		StoredItems:        0, // Initialize with 0 stored items
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Items:              []models.Game{}, // Initialize empty items slice
	}
}

// TransformUpdateRequestToModel converts an UpdateSublocationRequest to a models.Sublocation
// by merging the request fields with an existing sublocation.
// This is a pure function that only handles the transformation of data.
// No validation, error handling, or business logic should be included here.
func TransformUpdateRequestToModel(req types.UpdateSublocationRequest, existing models.Sublocation) models.Sublocation {
	// Start with the existing sublocation
	updated := existing

	// Only update fields that are provided in the request
	if req.Name != "" {
		updated.Name = req.Name
	}
	if req.LocationType != "" {
		updated.LocationType = req.LocationType
	}

	// Always update the UpdatedAt timestamp
	updated.UpdatedAt = time.Now()

	return updated
}

// TransformMoveGameRequest validates and prepares a move game request
// This is a pure function that only handles the transformation of data.
// No validation, error handling, or business logic should be included here.
func TransformMoveGameRequest(req types.MoveGameRequest) models.GameLocationMove {
	return models.GameLocationMove{
		UserGameID:         req.UserGameID,
		TargetSublocationID: req.TargetSublocationID,
	}
}

// TransformRemoveGameRequest validates and prepares a remove game request
// This is a pure function that only handles the transformation of data.
// No validation, error handling, or business logic should be included here.
func TransformRemoveGameRequest(req types.RemoveGameRequest) models.GameLocationRemove {
	return models.GameLocationRemove{
		UserGameID: req.UserGameID,
	}
}
