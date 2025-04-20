package registry

import (
	"fmt"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/library"
	"github.com/lokeam/qko-beta/internal/search"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/wishlist"
)

// Registry provides access to all domain-specific services
type Registry struct {
    appCtx *appcontext.AppContext

    // Keep concrete service references
    searchServices map[string]services.SearchService
    libraryServices map[string]library.LibraryService
		wishlistServices map[string]wishlist.WishlistService
}

// NewRegistry creates a service registry
func NewRegistry(appCtx *appcontext.AppContext) *Registry {
    registry := &Registry{
        appCtx: appCtx,
        searchServices: make(map[string]services.SearchService),
        libraryServices: make(map[string]library.LibraryService),
				wishlistServices: make(map[string]wishlist.WishlistService),
    }

    registry.registerServices()
    return registry
}

// Register all services
func (r *Registry) registerServices() {
    // Register search services
		gameSearchService, err := search.NewGameSearchService(r.appCtx)
		if err != nil {
			r.appCtx.Logger.Error("Failed to register game search service", map[string]any{
				"error": err.Error(),
			})
		}
		r.searchServices["games"] = gameSearchService

		gameLibraryService, err := library.NewGameLibraryService(r.appCtx)
		if err != nil {
			r.appCtx.Logger.Error("Failed to register game library service", map[string]any{
				"error": err.Error(),
			})
		}
		r.libraryServices["games"] = gameLibraryService

		gameWishlistService, err := wishlist.NewGameWishlistService(r.appCtx)
		if err != nil {
			r.appCtx.Logger.Error("Failed to register game wishlist service", map[string]any{
				"error": err.Error(),
			})
		}
		r.wishlistServices["games"] = gameWishlistService
}

// GetSearchService returns a specific search service
func (r *Registry) GetSearchService(domain string) (services.SearchService, error) {
    service, exists := r.searchServices[domain]
    if !exists {
        return nil, fmt.Errorf("unsupported search domain: %s", domain)
    }
    return service, nil
}

// GetLibraryService returns a specific library service
func (r *Registry) GetLibraryService(domain string) (library.LibraryService, error) {
    service, exists := r.libraryServices[domain]
    if !exists {
        return nil, fmt.Errorf("unsupported library domain: %s", domain)
    }
    return service, nil
}

// GetWishlistService returns a specific wishlist service
func (r *Registry) GetWishlistService(domain string) (wishlist.WishlistService, error) {
	service, exists := r.wishlistServices[domain]
	if !exists {
		return nil, fmt.Errorf("unsupported wishlist domain: %s", domain)
	}

	return service, nil
}
