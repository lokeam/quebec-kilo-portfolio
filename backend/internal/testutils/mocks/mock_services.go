package mocks

import (
	"github.com/lokeam/qko-beta/internal/services"
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
		Dashboard:     &MockDashboardService{},
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
	Dashboard     services.DashboardService
}
