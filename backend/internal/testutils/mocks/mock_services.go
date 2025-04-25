package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/services"
)

// NewMockServices creates a new app.Services instance with mocks for testing
func NewMockServices() *MockServices {
	// Create mock services
	libraryMap := make(services.DomainLibraryServices)
	libraryMap["games"] = &MockLibraryService{}

	searchMap := make(services.DomainSearchServices)
	searchMap["games"] = &MockSearchService{}

	return &MockServices{
		Digital:       &MockDigitalService{},
		Physical:      &MockPhysicalService{},
		Sublocation:   &MockSublocationService{},
		Library:       &MockLibraryService{},
		LibraryMap:    libraryMap,
		Wishlist:      &MockWishlistService{},
		SearchFactory: &MockSearchServiceFactory{},
		SearchMap:     searchMap,
	}
}

// MockServices contains mock implementations of all application services
type MockServices struct {
	Digital       services.DigitalService
	Physical      services.PhysicalService
	Sublocation   services.SublocationService
	Library       services.LibraryService
	LibraryMap    services.DomainLibraryServices
	Wishlist      services.WishlistService
	SearchFactory services.SearchServiceFactory
	SearchMap     services.DomainSearchServices
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

func (m *MockDigitalService) GetUserDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
	if m.GetUserDigitalLocationsFunc != nil {
		return m.GetUserDigitalLocationsFunc(ctx, userID)
	}
	return []models.DigitalLocation{}, nil
}

func (m *MockDigitalService) GetDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
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

func (m *MockDigitalService) AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
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

func (m *MockDigitalService) RemoveDigitalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error) {
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

func (m *MockDigitalService) AddSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
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

func (m *MockDigitalService) RemoveSubscription(ctx context.Context, locationID string) error {
	if m.RemoveSubscriptionFunc != nil {
		return m.RemoveSubscriptionFunc(ctx, locationID)
	}
	return nil
}

func (m *MockDigitalService) GetPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	if m.GetPaymentsFunc != nil {
		return m.GetPaymentsFunc(ctx, locationID)
	}
	return []models.Payment{}, nil
}

func (m *MockDigitalService) AddPayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	if m.AddPaymentFunc != nil {
		return m.AddPaymentFunc(ctx, payment)
	}
	return &models.Payment{}, nil
}

func (m *MockDigitalService) GetPayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
	if m.GetPaymentFunc != nil {
		return m.GetPaymentFunc(ctx, paymentID)
	}
	return &models.Payment{}, nil
}

// MockPhysicalService implements services.PhysicalService
type MockPhysicalService struct{}

func (m *MockPhysicalService) GetUserPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
	return []models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) GetPhysicalLocation(ctx context.Context, userID, locationID string) (models.PhysicalLocation, error) {
	return models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) AddPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	return models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	return models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) DeletePhysicalLocation(ctx context.Context, userID, locationID string) error {
	return nil
}

// MockSublocationService implements services.SublocationService
type MockSublocationService struct{}

func (m *MockSublocationService) GetSublocations(ctx context.Context, userID string) ([]models.Sublocation, error) {
	return []models.Sublocation{}, nil
}

func (m *MockSublocationService) GetSublocation(ctx context.Context, userID, locationID string) (models.Sublocation, error) {
	return models.Sublocation{}, nil
}

func (m *MockSublocationService) AddSublocation(ctx context.Context, userID string, location models.Sublocation) (models.Sublocation, error) {
	return models.Sublocation{}, nil
}

func (m *MockSublocationService) UpdateSublocation(ctx context.Context, userID string, location models.Sublocation) error {
	return nil
}

func (m *MockSublocationService) DeleteSublocation(ctx context.Context, userID, locationID string) error {
	return nil
}

// MockLibraryService implements services.LibraryService
type MockLibraryService struct{}

func (m *MockLibraryService) GetLibraryItems(ctx context.Context, userID string) ([]models.Game, error) {
	return []models.Game{}, nil
}

func (m *MockLibraryService) AddGameToLibrary(ctx context.Context, userID string, game models.Game) error {
	return nil
}

func (m *MockLibraryService) DeleteGameFromLibrary(ctx context.Context, userID string, gameID int64) error {
	return nil
}

func (m *MockLibraryService) GetGameByID(ctx context.Context, userID string, gameID int64) (models.Game, error) {
	return models.Game{}, nil
}

func (m *MockLibraryService) UpdateGameInLibrary(ctx context.Context, userID string, game models.Game) error {
	return nil
}

// MockWishlistService implements services.WishlistService
type MockWishlistService struct{}

func (m *MockWishlistService) GetWishlistItems(ctx context.Context, userID string) ([]models.Game, error) {
	return []models.Game{}, nil
}

// MockSearchService implements services.SearchService
type MockSearchService struct{}

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

// MockSearchServiceFactory implements services.SearchServiceFactory
type MockSearchServiceFactory struct{}

func (m *MockSearchServiceFactory) GetService(domain string) (services.SearchService, error) {
	return &MockSearchService{}, nil
}