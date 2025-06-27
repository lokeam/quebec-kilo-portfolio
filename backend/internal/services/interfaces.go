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
	GetAllDigitalLocationsBFF(ctx context.Context, userID string) (types.DigitalLocationsBFFResponse, error)

	CreateDigitalLocation(ctx context.Context, userID string, request types.DigitalLocationRequest) (models.DigitalLocation, error)
	UpdateDigitalLocation(ctx context.Context, userID string, request types.DigitalLocationRequest) error

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
	GetAllPhysicalLocationsBFF(ctx context.Context, userID string) (types.LocationsBFFResponse, error)
	GetSinglePhysicalLocation(ctx context.Context, userID, locationID string) (models.PhysicalLocation, error)
	CreatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	DeletePhysicalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error)
	InvalidateCache(ctx context.Context, cacheKey string) error
}

// SublocationService defines operations for managing sublocations
type SublocationService interface {
	GetSublocations(ctx context.Context, userID string) ([]models.Sublocation, error)
	GetSingleSublocation(ctx context.Context, userID, locationID string) (models.Sublocation, error)
	CreateSublocation(ctx context.Context, userID string, req types.CreateSublocationRequest) (models.Sublocation, error)
	UpdateSublocation(ctx context.Context, userID string, locationID string, req types.UpdateSublocationRequest) error
	DeleteSublocation(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error)
	MoveGame(ctx context.Context, userID string, req types.MoveGameRequest) error
	RemoveGame(ctx context.Context, userID string, req types.RemoveGameRequest) error
}

// LibraryService defines operations for managing the game library
type LibraryService interface {
	CreateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error
	GetAllLibraryItemsBFF(ctx context.Context, userID string) (types.LibraryBFFResponseFINAL, error)
	GetSingleLibraryGame(ctx context.Context, userID string, gameID int64) (types.LibraryGameItemBFFResponseFINAL, error)

	UpdateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error
	DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error
	InvalidateUserCache(ctx context.Context, userID string) error
	IsGameInLibraryBFF(ctx context.Context, userID string, gameID int64) (bool, error)
}

// WishlistService defines operations for managing the wishlist
type WishlistService interface {
	GetWishlistItems(ctx context.Context, userID string) ([]models.GameToSave, error)
}

// SearchService defines operations for searching
type SearchService interface {
	Search(ctx context.Context, req searchdef.SearchRequest) (*searchdef.SearchResult, error)
	GetAllGameStorageLocationsBFF(ctx context.Context, userID string) (types.AddGameFormStorageLocationsResponse, error)
}

// SpendTrackingService defines operations for managing spend tracking
type SpendTrackingService interface {
	GetSpendTrackingBFFResponse(ctx context.Context, userID string) (types.SpendTrackingBFFResponseFINAL, error)
	CreateOneTimePurchase(ctx context.Context, userID string, request types.SpendTrackingRequest) (models.SpendTrackingOneTimePurchaseDB, error)
}

type DashboardService interface {
	GetDashboardBFFResponse(ctx context.Context, userID string) (types.DashboardBFFResponse, error)
}