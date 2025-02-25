package library

import (
	"context"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/types"
)

type GameLibraryService struct {
	logger interfaces.Logger
	// TODO: Add dependencies such as santizer, validator, etc
}

type LibraryService interface {
	GetLibraryItems(ctx context.Context, userID string) ([]types.Game, error)
}

func NewGameLibraryService(appContext *appcontext.AppContext) *GameLibraryService {
	return &GameLibraryService{
		logger: appContext.Logger,
	}
}

func (ls *GameLibraryService) GetLibraryItems(ctx context.Context, userID string) ([]types.Game, error) {
	// TODO: Implement, return an empty list for testing
	return []types.Game{}, nil
}
