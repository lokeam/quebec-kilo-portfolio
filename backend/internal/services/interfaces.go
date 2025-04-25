package services

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
)

// DigitalService defines operations for managing digital locations
type DigitalService interface {
	GetUserDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
	FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error)
	AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error
	RemoveDigitalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error)

	// Game Management Operations
	AddGameToDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationID(ctx context.Context, userID string, locationID string) ([]models.Game, error)

	// Subscription management
	GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error)
	AddSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription models.Subscription) error
	RemoveSubscription(ctx context.Context, locationID string) error

	// Payment management
	GetPayments(ctx context.Context, locationID string) ([]models.Payment, error)
	AddPayment(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPayment(ctx context.Context, paymentID int64) (*models.Payment, error)
}

// PhysicalService defines operations for managing physical locations
type PhysicalService interface {
	GetUserPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	GetPhysicalLocation(ctx context.Context, userID, locationID string) (models.PhysicalLocation, error)
	AddPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	DeletePhysicalLocation(ctx context.Context, userID, locationID string) error
}

// SublocationService defines operations for managing sublocations
type SublocationService interface {
	GetSublocations(ctx context.Context, userID string) ([]models.Sublocation, error)
	GetSublocation(ctx context.Context, userID, locationID string) (models.Sublocation, error)
	AddSublocation(ctx context.Context, userID string, location models.Sublocation) (models.Sublocation, error)
	UpdateSublocation(ctx context.Context, userID string, location models.Sublocation) error
	DeleteSublocation(ctx context.Context, userID, locationID string) error
}

// LibraryService defines operations for managing the game library
type LibraryService interface {
	GetLibraryItems(ctx context.Context, userID string) ([]models.Game, error)
	AddGameToLibrary(ctx context.Context, userID string, game models.Game) error
	DeleteGameFromLibrary(ctx context.Context, userID string, gameID int64) error
	GetGameByID(ctx context.Context, userID string, gameID int64) (models.Game, error)
	UpdateGameInLibrary(ctx context.Context, userID string, game models.Game) error
}

// WishlistService defines operations for managing the wishlist
type WishlistService interface {
	GetWishlistItems(ctx context.Context, userID string) ([]models.Game, error)
}

// SearchService defines operations for searching
type SearchService interface {
	Search(ctx context.Context, req searchdef.SearchRequest) (*searchdef.SearchResult, error)
}

// SearchServiceFactory defines operations for creating search services
type SearchServiceFactory interface {
	GetService(domain string) (SearchService, error)
}

// DomainLibraryServices is a map of domain-specific library services
type DomainLibraryServices map[string]LibraryService

// DomainSearchServices is a map of domain-specific search services
type DomainSearchServices map[string]SearchService