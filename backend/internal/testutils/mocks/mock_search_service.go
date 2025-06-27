package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/stretchr/testify/mock"
)

type MockSearchService struct {
	mock.Mock
}

func (m *MockSearchService) Search(
	ctx context.Context,
	req searchdef.SearchRequest,
) (*searchdef.SearchResult, error) {
	// Return empty search result
	return &searchdef.SearchResult{
		Games: []models.Game{},
		Meta: searchdef.SearchMeta{
			Total: 0,
			CurrentPage: 1,
			ResultsPerPage: 20,
		},
	}, nil
}

func (m *MockSearchService) GetAllGameStorageLocationsBFF(
	ctx context.Context,
	userID string,
) (types.AddGameFormStorageLocationsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return types.AddGameFormStorageLocationsResponse{}, args.Error(1)
	}
	return args.Get(0).(types.AddGameFormStorageLocationsResponse), args.Error(1)
}