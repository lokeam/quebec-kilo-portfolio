package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

// Game Management Operations
func (m *MockDigitalDbAdapter) AddGameToDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
	gameID int64,
) error {
	if m.AddGameToDigitalLocationFunc != nil {
		return m.AddGameToDigitalLocationFunc(ctx, userID, locationID, gameID)
	}
	return nil
}

func (m *MockDigitalDbAdapter) RemoveGameFromDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
	gameID int64,
) error {
	if m.RemoveGameFromDigitalLocationFunc != nil {
		return m.RemoveGameFromDigitalLocationFunc(ctx, userID, locationID, gameID)
	}
	return nil
}

func (m *MockDigitalDbAdapter) GetGamesByDigitalLocationID(
	ctx context.Context,
	userID string,
	locationID string,
) ([]models.Game, error) {
	if m.GetGamesByDigitalLocationIDFunc != nil {
		return m.GetGamesByDigitalLocationIDFunc(ctx, userID, locationID)
	}
	return nil, nil
}