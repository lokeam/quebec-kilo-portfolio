package library

import (
	"context"
	"errors"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
)

type GameLibraryService struct {
	adapter *LibraryDbAdapter
	logger  interfaces.Logger
	// TODO: Add dependencies such as santizer, validator, etc
}

type LibraryService interface {
	GetLibraryItems(ctx context.Context, userID string) ([]models.Game, error)
	AddGameToLibrary(ctx context.Context, userID string, game models.Game) error
	DeleteGameFromLibrary(ctx context.Context, userID string, gameID int64) error
	GetGameByID(ctx context.Context, userID string, gameID int64) (models.Game, error)
	UpdateGameInLibrary(ctx context.Context, userID string, game models.Game) error
}

// Constructor that properly initializes the adapter
func NewGameLibraryService(appCtx *appcontext.AppContext) (*GameLibraryService, error) {
	// Create and initialize the database adapter
	adapter, err := NewLibraryDbAdapter(appCtx)
	if err != nil {
		return nil, err
	}

	return &GameLibraryService{
		adapter: adapter,
		logger:  appCtx.Logger,
	}, nil
}

// GET
func (ls *GameLibraryService) GetLibraryItems(ctx context.Context, userID string) ([]models.Game, error) {
	return ls.adapter.GetLibraryItems(ctx, userID)
}

// POST
func (ls *GameLibraryService) AddGameToLibrary(ctx context.Context, userID string, game models.Game) error {
	return ls.adapter.AddGameToLibrary(ctx, userID, game.ID)
}

// DELETE
func (ls *GameLibraryService) DeleteGameFromLibrary(ctx context.Context, userID string, gameID int64) error {
	return ls.adapter.RemoveGameFromLibrary(ctx, userID, gameID)
}

func (ls *GameLibraryService) GetGameByID(ctx context.Context, userID string, gameID int64) (models.Game, error) {
	// Single database call to get the game while verifying ownership
	game, exists, err := ls.adapter.GetUserGame(ctx, userID, gameID)
	if err != nil {
		return models.Game{}, err
	}

	if !exists {
		return models.Game{}, ErrGameNotFound
	}

	return game, nil
}

func (ls *GameLibraryService) UpdateGameInLibrary(ctx context.Context, userID string, game models.Game) error {
	// NOTE: May need implementation if I want to support updating game metadata
	return errors.New("not implemented")
}