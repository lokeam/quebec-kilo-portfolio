package app

import (
	"fmt"

	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/library"
	"github.com/lokeam/qko-beta/internal/locations/digital"
	"github.com/lokeam/qko-beta/internal/locations/physical"
	"github.com/lokeam/qko-beta/internal/locations/sublocation"
	"github.com/lokeam/qko-beta/internal/media_storage"
	"github.com/lokeam/qko-beta/internal/search"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/wishlist"
)

// Services contains all application services
type Services struct {
	Digital       services.DigitalService
	Physical      services.PhysicalService
	Sublocation   services.SublocationService
	Library       services.LibraryService
	LibraryMap    services.DomainLibraryServices
	Wishlist      services.WishlistService
	SearchFactory services.SearchServiceFactory
	SearchMap     services.DomainSearchServices
	Analytics     analytics.Service
	MediaStorage  media_storage.MediaStorageService
}

// NewServices initializes all application services
func NewServices(appCtx *appcontext.AppContext) (*Services, error) {
	servicesObj := &Services{}
	var err error

	// Initialize digital service
	digitalService, err := digital.NewGameDigitalService(appCtx)
	if err != nil {
		return nil, fmt.Errorf("initializing digital service: %w", err)
	}
	servicesObj.Digital = digitalService

	// Initialize physical service
	physicalService, err := physical.NewGamePhysicalService(appCtx)
	if err != nil {
		return nil, fmt.Errorf("initializing physical service: %w", err)
	}
	servicesObj.Physical = physicalService

	// Initialize sublocation service - Pass physical service for cache refresh upon crud operation
	sublocationService, err := sublocation.NewGameSublocationService(appCtx, physicalService)
	if err != nil {
		return nil, fmt.Errorf("initializing sublocation service: %w", err)
	}
	servicesObj.Sublocation = sublocationService

	// Initialize library service
	libraryService, err := library.NewGameLibraryService(appCtx)
	if err != nil {
		return nil, fmt.Errorf("initializing library service: %w", err)
	}
	servicesObj.Library = libraryService

	// Create library services map
	servicesObj.LibraryMap = make(services.DomainLibraryServices)
	servicesObj.LibraryMap["games"] = libraryService

	// Initialize wishlist service
	wishlistService, err := wishlist.NewGameWishlistService(appCtx)
	if err != nil {
		return nil, fmt.Errorf("initializing wishlist service: %w", err)
	}
	servicesObj.Wishlist = wishlistService

	// Initialize search services
	searchFactory := search.NewSearchServiceFactory(appCtx)
	servicesObj.SearchFactory = searchFactory
	servicesObj.SearchMap = make(services.DomainSearchServices)

	gameSearchService, err := searchFactory.GetService("games")
	if err == nil {
		servicesObj.SearchMap["games"] = gameSearchService
	} else {
		appCtx.Logger.Warn("Game search service not available", map[string]any{
			"error": err,
		})
	}

	// Initialize analytics service
	analyticsService, err := analytics.NewAnalyticsService(appCtx)
	if err != nil {
		return nil, fmt.Errorf("initializing analytics service: %w", err)
	}
	servicesObj.Analytics = analyticsService

	// Initialize media storage service
	mediaStorageService, err := media_storage.NewMediaStorageService(analyticsService, appCtx.Logger)
	if err != nil {
		return nil, fmt.Errorf("initializing media storage service: %w", err)
	}
	servicesObj.MediaStorage = mediaStorageService

	return servicesObj, nil
}