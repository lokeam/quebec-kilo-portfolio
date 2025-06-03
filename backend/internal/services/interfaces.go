package services

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/types"
)

// DigitalService defines operations for managing digital locations
type DigitalService interface {
	GetAllDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetSingleDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
	FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error)
	CreateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error
	DeleteDigitalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error)

	// Game Management Operations
	AddGameToDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationID(ctx context.Context, userID string, locationID string) ([]models.Game, error)

	// Subscription management
	GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error)
	CreateSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription models.Subscription) error
	DeleteSubscription(ctx context.Context, locationID string) error

	// Payment management
	GetAllPayments(ctx context.Context, locationID string) ([]models.Payment, error)
	CreatePayment(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetSinglePayment(ctx context.Context, paymentID int64) (*models.Payment, error)
}

// PhysicalService defines operations for managing physical locations
type PhysicalService interface {
	GetAllPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	GetSinglePhysicalLocation(ctx context.Context, userID, locationID string) (models.PhysicalLocation, error)
	CreatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	DeletePhysicalLocation(ctx context.Context, userID, locationID string) error
}

// SublocationService defines operations for managing sublocations
type SublocationService interface {
	GetSublocations(ctx context.Context, userID string) ([]models.Sublocation, error)
	GetSingleSublocation(ctx context.Context, userID, locationID string) (models.Sublocation, error)
	CreateSublocation(ctx context.Context, userID string, location models.Sublocation) (models.Sublocation, error)
	UpdateSublocation(ctx context.Context, userID string, location models.Sublocation) error
	DeleteSublocation(ctx context.Context, userID, locationID string) error
}

// LibraryService defines operations for managing the game library
type LibraryService interface {
	CreateLibraryGame(ctx context.Context, userID string, game models.LibraryGame) error
	GetAllLibraryGames(
		ctx context.Context,
		userID string,
	) (
		[]types.LibraryGameDBResult,
		[]types.LibraryGamePhysicalLocationDBResponse,
		[]types.LibraryGameDigitalLocationDBResponse,
		error,
	)
	GetSingleLibraryGame(
		ctx context.Context,
		userID string,
		gameID int64,
	) (
		types.LibraryGameDBResult,
		[]types.LibraryGamePhysicalLocationDBResponse,
		[]types.LibraryGameDigitalLocationDBResponse,
		error,
	)
	UpdateLibraryGame(ctx context.Context, userID string, game models.LibraryGame) error
	DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error
}

// WishlistService defines operations for managing the wishlist
type WishlistService interface {
	GetWishlistItems(ctx context.Context, userID string) ([]models.LibraryGame, error)
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