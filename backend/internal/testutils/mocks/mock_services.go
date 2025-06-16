package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/stretchr/testify/mock"
)

// NewMockServices creates a new app.Services instance with mocks for testing
func NewMockServices() *MockServices {
	// Create mock services
	return &MockServices{
		Digital:       &MockDigitalService{},
		Physical:      &MockPhysicalService{},
		Sublocation:   &MockSublocationService{},
		Library:       &MockLibraryService{},
		Wishlist:      &MockWishlistService{},
		Search:        &MockSearchService{},
		SpendTracking: &MockSpendTrackingService{},
	}
}

// MockServices contains mock implementations of all application services
type MockServices struct {
	Digital       services.DigitalService
	Physical      services.PhysicalService
	Sublocation   services.SublocationService
	Library       services.LibraryService
	Wishlist      services.WishlistService
	Search        services.SearchService
	SpendTracking services.SpendTrackingService
}

// Mocks for each service type

// MockDigitalService implements services.DigitalService
type MockDigitalService struct {
	GetUserDigitalLocationsFunc    func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetDigitalLocationFunc         func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
	FindDigitalLocationByNameFunc  func(ctx context.Context, userID string, name string) (models.DigitalLocation, error)
	AddDigitalLocationFunc         func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocationFunc      func(ctx context.Context, userID string, location models.DigitalLocation) error
	RemoveDigitalLocationFunc      func(ctx context.Context, userID string, locationIDs []string) (int64, error)
	AddGameToDigitalLocationFunc   func(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocationFunc func(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationIDFunc func(ctx context.Context, userID string, locationID string) ([]models.Game, error)
	GetSubscriptionFunc           func(ctx context.Context, locationID string) (*models.Subscription, error)
	AddSubscriptionFunc           func(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscriptionFunc        func(ctx context.Context, subscription models.Subscription) error
	RemoveSubscriptionFunc        func(ctx context.Context, locationID string) error
	GetPaymentsFunc              func(ctx context.Context, locationID string) ([]models.Payment, error)
	AddPaymentFunc               func(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPaymentFunc               func(ctx context.Context, paymentID int64) (*models.Payment, error)
}

func (m *MockDigitalService) GetAllDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
	if m.GetUserDigitalLocationsFunc != nil {
		return m.GetUserDigitalLocationsFunc(ctx, userID)
	}
	return []models.DigitalLocation{}, nil
}

func (m *MockDigitalService) GetSingleDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
	if m.GetDigitalLocationFunc != nil {
		return m.GetDigitalLocationFunc(ctx, userID, locationID)
	}
	return models.DigitalLocation{}, nil
}

func (m *MockDigitalService) FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
	if m.FindDigitalLocationByNameFunc != nil {
		return m.FindDigitalLocationByNameFunc(ctx, userID, name)
	}
	return models.DigitalLocation{}, nil
}

func (m *MockDigitalService) CreateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
	if m.AddDigitalLocationFunc != nil {
		return m.AddDigitalLocationFunc(ctx, userID, location)
	}
	return models.DigitalLocation{}, nil
}

func (m *MockDigitalService) UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error {
	if m.UpdateDigitalLocationFunc != nil {
		return m.UpdateDigitalLocationFunc(ctx, userID, location)
	}
	return nil
}

func (m *MockDigitalService) DeleteDigitalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error) {
	if m.RemoveDigitalLocationFunc != nil {
		return m.RemoveDigitalLocationFunc(ctx, userID, locationIDs)
	}
	return 0, nil
}

func (m *MockDigitalService) AddGameToDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error {
	if m.AddGameToDigitalLocationFunc != nil {
		return m.AddGameToDigitalLocationFunc(ctx, userID, locationID, gameID)
	}
	return nil
}

func (m *MockDigitalService) RemoveGameFromDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error {
	if m.RemoveGameFromDigitalLocationFunc != nil {
		return m.RemoveGameFromDigitalLocationFunc(ctx, userID, locationID, gameID)
	}
	return nil
}

func (m *MockDigitalService) GetGamesByDigitalLocationID(ctx context.Context, userID string, locationID string) ([]models.Game, error) {
	if m.GetGamesByDigitalLocationIDFunc != nil {
		return m.GetGamesByDigitalLocationIDFunc(ctx, userID, locationID)
	}
	return []models.Game{}, nil
}

func (m *MockDigitalService) GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error) {
	if m.GetSubscriptionFunc != nil {
		return m.GetSubscriptionFunc(ctx, locationID)
	}
	return &models.Subscription{}, nil
}

func (m *MockDigitalService) CreateSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
	if m.AddSubscriptionFunc != nil {
		return m.AddSubscriptionFunc(ctx, subscription)
	}
	return &models.Subscription{}, nil
}

func (m *MockDigitalService) UpdateSubscription(ctx context.Context, subscription models.Subscription) error {
	if m.UpdateSubscriptionFunc != nil {
		return m.UpdateSubscriptionFunc(ctx, subscription)
	}
	return nil
}

func (m *MockDigitalService) DeleteSubscription(ctx context.Context, locationID string) error {
	if m.RemoveSubscriptionFunc != nil {
		return m.RemoveSubscriptionFunc(ctx, locationID)
	}
	return nil
}

func (m *MockDigitalService) GetAllPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	if m.GetPaymentsFunc != nil {
		return m.GetPaymentsFunc(ctx, locationID)
	}
	return []models.Payment{}, nil
}

func (m *MockDigitalService) CreatePayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	if m.AddPaymentFunc != nil {
		return m.AddPaymentFunc(ctx, payment)
	}
	return &models.Payment{}, nil
}

func (m *MockDigitalService) GetSinglePayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
	if m.GetPaymentFunc != nil {
		return m.GetPaymentFunc(ctx, paymentID)
	}
	return &models.Payment{}, nil
}

// MockPhysicalService implements services.PhysicalService
type MockPhysicalService struct {
	GetAllPhysicalLocationsFunc     func(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	GetAllPhysicalLocationsBFFFunc  func(ctx context.Context, userID string) (types.LocationsBFFResponse, error)
	GetSinglePhysicalLocationFunc   func(ctx context.Context, userID, locationID string) (models.PhysicalLocation, error)
	CreatePhysicalLocationFunc      func(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	UpdatePhysicalLocationFunc      func(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	DeletePhysicalLocationFunc      func(ctx context.Context, userID string, locationIDs []string) (int64, error)
	InvalidateCacheFunc             func(ctx context.Context, cacheKey string) error
}

func (m *MockPhysicalService) GetAllPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
	if m.GetAllPhysicalLocationsFunc != nil {
		return m.GetAllPhysicalLocationsFunc(ctx, userID)
	}
	return []models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) GetSinglePhysicalLocation(ctx context.Context, userID, locationID string) (models.PhysicalLocation, error) {
	if m.GetSinglePhysicalLocationFunc != nil {
		return m.GetSinglePhysicalLocationFunc(ctx, userID, locationID)
	}
	return models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) CreatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	if m.CreatePhysicalLocationFunc != nil {
		return m.CreatePhysicalLocationFunc(ctx, userID, location)
	}
	return models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	if m.UpdatePhysicalLocationFunc != nil {
		return m.UpdatePhysicalLocationFunc(ctx, userID, location)
	}
	return models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) DeletePhysicalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error) {
	if m.DeletePhysicalLocationFunc != nil {
		return m.DeletePhysicalLocationFunc(ctx, userID, locationIDs)
	}
	return 0, nil
}

func (m *MockPhysicalService) GetAllPhysicalLocationsBFF(ctx context.Context, userID string) (types.LocationsBFFResponse, error) {
	if m.GetAllPhysicalLocationsBFFFunc != nil {
		return m.GetAllPhysicalLocationsBFFFunc(ctx, userID)
	}
	return types.LocationsBFFResponse{}, nil
}

func (m *MockPhysicalService) InvalidateCache(ctx context.Context, cacheKey string) error {
	if m.InvalidateCacheFunc != nil {
		return m.InvalidateCacheFunc(ctx, cacheKey)
	}
	return nil
}

// MockSublocationService implements services.SublocationService
type MockSublocationService struct {
	GetSublocationsFunc     func(ctx context.Context, userID string) ([]models.Sublocation, error)
	GetSingleSublocationFunc func(ctx context.Context, userID, locationID string) (models.Sublocation, error)
	CreateSublocationFunc    func(ctx context.Context, userID string, req types.CreateSublocationRequest) (models.Sublocation, error)
	UpdateSublocationFunc    func(ctx context.Context, userID string, locationID string, req types.UpdateSublocationRequest) error
	DeleteSublocationFunc    func(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error)
	MoveGameFunc            func(ctx context.Context, userID string, req types.MoveGameRequest) error
	RemoveGameFunc          func(ctx context.Context, userID string, req types.RemoveGameRequest) error
}

func (m *MockSublocationService) GetSublocations(ctx context.Context, userID string) ([]models.Sublocation, error) {
	if m.GetSublocationsFunc != nil {
		return m.GetSublocationsFunc(ctx, userID)
	}
	return []models.Sublocation{}, nil
}

func (m *MockSublocationService) GetSingleSublocation(ctx context.Context, userID, locationID string) (models.Sublocation, error) {
	if m.GetSingleSublocationFunc != nil {
		return m.GetSingleSublocationFunc(ctx, userID, locationID)
	}
	return models.Sublocation{}, nil
}

func (m *MockSublocationService) CreateSublocation(ctx context.Context, userID string, req types.CreateSublocationRequest) (models.Sublocation, error) {
	if m.CreateSublocationFunc != nil {
		return m.CreateSublocationFunc(ctx, userID, req)
	}
	return models.Sublocation{}, nil
}

func (m *MockSublocationService) UpdateSublocation(ctx context.Context, userID string, locationID string, req types.UpdateSublocationRequest) error {
	if m.UpdateSublocationFunc != nil {
		return m.UpdateSublocationFunc(ctx, userID, locationID, req)
	}
	return nil
}

func (m *MockSublocationService) DeleteSublocation(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error) {
	if m.DeleteSublocationFunc != nil {
		return m.DeleteSublocationFunc(ctx, userID, sublocationIDs)
	}
	return types.DeleteSublocationResponse{
		Success: true,
		DeletedCount: len(sublocationIDs),
		SublocationIDs: sublocationIDs,
	}, nil
}

func (m *MockSublocationService) MoveGame(ctx context.Context, userID string, req types.MoveGameRequest) error {
	if m.MoveGameFunc != nil {
		return m.MoveGameFunc(ctx, userID, req)
	}
	return nil
}

func (m *MockSublocationService) RemoveGame(ctx context.Context, userID string, req types.RemoveGameRequest) error {
	if m.RemoveGameFunc != nil {
		return m.RemoveGameFunc(ctx, userID, req)
	}
	return nil
}

// MockLibraryService implements services.LibraryService
type MockLibraryService struct {
	mock.Mock
}

func (m *MockLibraryService) CreateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error {
	return nil
}

func (m *MockLibraryService) DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error {
	return nil
}

func (m *MockLibraryService) GetSingleLibraryGame(ctx context.Context, userID string, gameID int64) (types.LibraryGameItemBFFResponseFINAL, error) {
	args := m.Called(ctx, userID, gameID)
	return args.Get(0).(types.LibraryGameItemBFFResponseFINAL), args.Error(1)
}

func (m *MockLibraryService) UpdateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error {
	return nil
}

func (m *MockLibraryService) GetAllLibraryItemsBFF(ctx context.Context, userID string) (types.LibraryBFFResponseFINAL, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(types.LibraryBFFResponseFINAL), args.Error(1)
}

// InvalidateUserCache mocks the InvalidateUserCache method
func (m *MockLibraryService) InvalidateUserCache(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockLibraryService) IsGameInLibraryBFF(ctx context.Context, userID string, gameID int64) (bool, error) {
	args := m.Called(ctx, userID, gameID)
	return args.Bool(0), args.Error(1)
}

// MockWishlistService implements services.WishlistService
type MockWishlistService struct{}

func (m *MockWishlistService) GetWishlistItems(ctx context.Context, userID string) ([]models.GameToSave, error) {
	return []models.GameToSave{}, nil
}

// MockSearchService implements services.SearchService
type MockSearchService struct {
	mock.Mock
}

func (m *MockSearchService) Search(ctx context.Context, req searchdef.SearchRequest) (*searchdef.SearchResult, error) {
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

// GetAllGameStorageLocationsBFF mocks the GetAllGameStorageLocationsBFF method
func (m *MockSearchService) GetAllGameStorageLocationsBFF(ctx context.Context, userID string) (types.AddGameFormStorageLocationsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return types.AddGameFormStorageLocationsResponse{}, args.Error(1)
	}
	return args.Get(0).(types.AddGameFormStorageLocationsResponse), args.Error(1)
}

type MockSublocationDbAdapter struct {
	GetSublocationFunc func(ctx context.Context, userID string, sublocationID string) (models.Sublocation, error)
	GetUserSublocationsFunc func(ctx context.Context, userID string) ([]models.Sublocation, error)
	AddSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error)
	UpdateSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) error
	DeleteSublocationFunc func(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error)
	CheckGameInAnySublocationFunc func(ctx context.Context, userGameID string) (bool, error)
	CheckGameInSublocationFunc func(ctx context.Context, userGameID string, sublocationID string) (bool, error)
	CheckGameOwnershipFunc func(ctx context.Context, userID string, userGameID string) (bool, error)
	MoveGameToSublocationFunc func(ctx context.Context, userID string, userGameID string, targetSublocationID string) error
	RemoveGameFromSublocationFunc func(ctx context.Context, userID string, userGameID string) error
}

func (m *MockSublocationDbAdapter) CheckGameInAnySublocation(ctx context.Context, userGameID string) (bool, error) {
	if m.CheckGameInAnySublocationFunc != nil {
		return m.CheckGameInAnySublocationFunc(ctx, userGameID)
	}
	return false, nil
}

func (m *MockSublocationDbAdapter) CheckGameInSublocation(ctx context.Context, userGameID string, sublocationID string) (bool, error) {
	if m.CheckGameInSublocationFunc != nil {
		return m.CheckGameInSublocationFunc(ctx, userGameID, sublocationID)
	}
	return false, nil
}

func (m *MockSublocationDbAdapter) CheckGameOwnership(ctx context.Context, userID string, userGameID string) (bool, error) {
	if m.CheckGameOwnershipFunc != nil {
		return m.CheckGameOwnershipFunc(ctx, userID, userGameID)
	}
	return true, nil
}

func (m *MockSublocationDbAdapter) MoveGameToSublocation(ctx context.Context, userID string, userGameID string, targetSublocationID string) error {
	if m.MoveGameToSublocationFunc != nil {
		return m.MoveGameToSublocationFunc(ctx, userID, userGameID, targetSublocationID)
	}
	return nil
}

func (m *MockSublocationDbAdapter) RemoveGameFromSublocation(ctx context.Context, userID string, userGameID string) error {
	if m.RemoveGameFromSublocationFunc != nil {
		return m.RemoveGameFromSublocationFunc(ctx, userID, userGameID)
	}
	return nil
}

func (m *MockSublocationDbAdapter) GetSublocation(ctx context.Context, userID string, sublocationID string) (models.Sublocation, error) {
	if m.GetSublocationFunc != nil {
		return m.GetSublocationFunc(ctx, userID, sublocationID)
	}
	return models.Sublocation{}, nil
}

func (m *MockSublocationDbAdapter) GetAllSublocations(ctx context.Context, userID string) ([]models.Sublocation, error) {
	if m.GetUserSublocationsFunc != nil {
		return m.GetUserSublocationsFunc(ctx, userID)
	}
	return []models.Sublocation{}, nil
}

func (m *MockSublocationDbAdapter) CreateSublocation(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error) {
	if m.AddSublocationFunc != nil {
		return m.AddSublocationFunc(ctx, userID, sublocation)
	}
	return models.Sublocation{}, nil
}

func (m *MockSublocationDbAdapter) UpdateSublocation(ctx context.Context, userID string, sublocation models.Sublocation) error {
	if m.UpdateSublocationFunc != nil {
		return m.UpdateSublocationFunc(ctx, userID, sublocation)
	}
	return nil
}

func (m *MockSublocationDbAdapter) DeleteSublocation(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error) {
	if m.DeleteSublocationFunc != nil {
		return m.DeleteSublocationFunc(ctx, userID, sublocationIDs)
	}
	return types.DeleteSublocationResponse{
		Success: true,
		DeletedCount: len(sublocationIDs),
		SublocationIDs: sublocationIDs,
	}, nil
}

type MockSpendTrackingService struct {
	GetSpendTrackingBFFResponseFunc func(ctx context.Context, userID string) (types.SpendTrackingBFFResponseFINAL, error)
}

func (m *MockSpendTrackingService) GetSpendTrackingBFFResponse(ctx context.Context, userID string) (types.SpendTrackingBFFResponseFINAL, error) {
	if m.GetSpendTrackingBFFResponseFunc != nil {
		return m.GetSpendTrackingBFFResponseFunc(ctx, userID)
	}
	return types.SpendTrackingBFFResponseFINAL{}, nil
}