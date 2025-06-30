package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/stretchr/testify/mock"
)

type MockLibraryService struct {
	mock.Mock
}

func (m *MockLibraryService) CreateLibraryGame(
	ctx context.Context,
	userID string,
	game models.GameToSave,
) error {
	return nil
}

func (m *MockLibraryService) DeleteLibraryGame(
	ctx context.Context,
	userID string,
	gameID int64,
) error {
	return nil
}

func (m *MockLibraryService) GetSingleLibraryGame(
	ctx context.Context,
	userID string,
	gameID int64,
) (types.LibraryGameItemBFFResponseFINAL, error) {
	args := m.Called(ctx, userID, gameID)
	return args.Get(0).(types.LibraryGameItemBFFResponseFINAL), args.Error(1)
}

func (m *MockLibraryService) UpdateLibraryGame(
	ctx context.Context,
	userID string,
	game models.GameToSave,
) error {
	return nil
}

// LEGACY BFF RESPONSE - MARKED FOR DELETION
func (m *MockLibraryService) GetAllLibraryItemsBFF(
	ctx context.Context,
	userID string,
) (types.LibraryBFFResponseFINAL, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(types.LibraryBFFResponseFINAL), args.Error(1)
}

func (m *MockLibraryService) GetLibraryRefactoredBFFResponse(
	ctx context.Context,
	userID string,
) (types.LibraryBFFRefactoredResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(types.LibraryBFFRefactoredResponse), args.Error(1)
}

// InvalidateUserCache mocks the InvalidateUserCache method
func (m *MockLibraryService) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockLibraryService) IsGameInLibraryBFF(
	ctx context.Context,
	userID string,
	gameID int64,
) (bool, error) {
	args := m.Called(ctx, userID, gameID)
	return args.Bool(0), args.Error(1)
}