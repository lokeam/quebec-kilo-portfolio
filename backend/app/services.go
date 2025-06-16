package app

import (
	"fmt"

	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/infrastructure/cache"
	"github.com/lokeam/qko-beta/internal/library"
	"github.com/lokeam/qko-beta/internal/locations/digital"
	"github.com/lokeam/qko-beta/internal/locations/physical"
	"github.com/lokeam/qko-beta/internal/locations/sublocation"
	"github.com/lokeam/qko-beta/internal/media_storage"
	"github.com/lokeam/qko-beta/internal/search"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/spend_tracking"
	"github.com/lokeam/qko-beta/internal/wishlist"
)

// Services contains all application services
type Services struct {
	Digital       services.DigitalService
	Physical      services.PhysicalService
	Sublocation   services.SublocationService
	Library       services.LibraryService
	Wishlist      services.WishlistService
	Search        services.SearchService
	SpendTracking services.SpendTrackingService
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
	// Create cache wrapper
	cacheWrapper, err := cache.NewCacheWrapper(
		appCtx.RedisClient,
		appCtx.Config.Redis.RedisTTL,
		appCtx.Config.Redis.RedisTimeout,
		appCtx.Logger,
	)
	if err != nil {
		return nil, fmt.Errorf("initializing cache wrapper: %w", err)
	}

	// Create library cache adapter
	libraryCacheAdapter, err := library.NewLibraryCacheAdapter(cacheWrapper)
	if err != nil {
		return nil, fmt.Errorf("initializing library cache adapter: %w", err)
	}

	// Create library db adapter
	libraryDbAdapter, err := library.NewLibraryDbAdapter(appCtx)
	if err != nil {
		return nil, fmt.Errorf("initializing library db adapter: %w", err)
	}

	// Initialize library service with dependencies
	libraryService, err := library.NewGameLibraryService(appCtx, libraryDbAdapter, libraryCacheAdapter)
	if err != nil {
		return nil, fmt.Errorf("initializing library service: %w", err)
	}
	servicesObj.Library = libraryService

	// Initialize spend tracking service

	// Create spend tracking cache adapter
	spendTrackingCacheAdapter, err := spend_tracking.NewSpendTrackingCacheAdapter(cacheWrapper)
	if err != nil {
		return nil, fmt.Errorf("initializing spend tracking cache adapter: %w", err)
	}

	// Create spend tracking db adapter
	spendTrackingDbAdapter, err := spend_tracking.NewSpendTrackingDbAdapter(appCtx)
	if err != nil {
		return nil, fmt.Errorf("initializing spend tracking db adapter: %w", err)
	}

	spendTrackingService, err := spend_tracking.NewSpendTrackingService(appCtx, spendTrackingDbAdapter, spendTrackingCacheAdapter)
	if err != nil {
		return nil, fmt.Errorf("initializing spend tracking service: %w", err)
	}
	servicesObj.SpendTracking = spendTrackingService

	// Initialize wishlist service
	wishlistService, err := wishlist.NewGameWishlistService(appCtx)
	if err != nil {
		return nil, fmt.Errorf("initializing wishlist service: %w", err)
	}
	servicesObj.Wishlist = wishlistService

	// Initialize search services
	gameSearchService, err := search.NewGameSearchService(appCtx)
	if err != nil {
			return nil, fmt.Errorf("initializing game search service: %w", err)
	}
	servicesObj.Search = gameSearchService

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