package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
)

// Service defines the interface for analytics functionality
type Service interface {
	// Get analytics data for requested domains
	GetAnalytics(ctx context.Context, userID string, domains []string) (map[string]any, error)

	// Individual domain getters
	GetGeneralStats(ctx context.Context, userID string) (*GeneralStats, error)
	GetFinancialStats(ctx context.Context, userID string) (*FinancialStats, error)
	GetStorageStats(ctx context.Context, userID string) (*StorageStats, error)
	GetInventoryStats(ctx context.Context, userID string) (*InventoryStats, error)
	GetWishlistStats(ctx context.Context, userID string) (*WishlistStats, error)

	// Cache invalidation
	InvalidateDomain(ctx context.Context, userID string, domain string) error
	InvalidateDomains(ctx context.Context, userID string, domains []string) error
}

// service implements the Service interface
type service struct {
	repo   Repository
	cache  interfaces.CacheWrapper
	logger interfaces.Logger
	ttl    map[string]time.Duration
}

// NewService creates a new analytics service
func NewService(repo Repository, cache interfaces.CacheWrapper, logger interfaces.Logger) Service {
	// Set default TTLs for different domains
	ttl := map[string]time.Duration{
		DomainGeneral:   time.Duration(TTLGeneralMinutes) * time.Minute,
		DomainFinancial: time.Duration(TTLFinancialMinutes) * time.Minute,
		DomainStorage:   time.Duration(TTLStorageMinutes) * time.Minute,
		DomainInventory: time.Duration(TTLInventoryMinutes) * time.Minute,
		DomainWishlist:  time.Duration(TTLWishlistMinutes) * time.Minute,
	}

	return &service{
		repo:   repo,
		cache:  cache,
		logger: logger,
		ttl:    ttl,
	}
}

// GetAnalytics returns analytics data for the requested domains
func (s *service) GetAnalytics(ctx context.Context, userID string, domains []string) (map[string]any, error) {
	result := make(map[string]any)

	// If no domains specified, use default domain
	if len(domains) == 0 {
			domains = []string{DomainGeneral}
	}

	// Process each requested domain using index-based loop
	for i := 0; i < len(domains); i++ {
			domain := domains[i]
			var data any
			var err error

			// Try to get from cache first
			data, found := s.getFromCache(ctx, userID, domain)
			if found {
					s.logger.Debug("Analytics cache hit", map[string]any{
							"domain": domain,
							"userID": userID,
					})
					result[domain] = data
					continue
			}

			// Cache miss - get from database
			s.logger.Debug("Analytics cache miss", map[string]any{
					"domain": domain,
					"userID": userID,
			})

			switch domain {
			case DomainGeneral:
					data, err = s.GetGeneralStats(ctx, userID)
			case DomainFinancial:
					data, err = s.GetFinancialStats(ctx, userID)
			case DomainStorage:
					data, err = s.GetStorageStats(ctx, userID)
			case DomainInventory:
					data, err = s.GetInventoryStats(ctx, userID)
			case DomainWishlist:
					data, err = s.GetWishlistStats(ctx, userID)
			default:
					s.logger.Warn("Unknown analytics domain requested", map[string]any{
							"domain": domain,
							"userID": userID,
					})
					continue // Skip unknown domains
			}

			if err != nil {
					s.logger.Error("Failed to get analytics data", map[string]any{
							"domain": domain,
							"userID": userID,
							"error":  err,
					})
					return nil, fmt.Errorf("error fetching %s analytics: %w", domain, err)
			}

			// Cache the result
			s.setInCache(ctx, userID, domain, data)
			result[domain] = data
	}

	return result, nil
}

// GetGeneralStats retrieves general statistics
func (s *service) GetGeneralStats(ctx context.Context, userID string) (*GeneralStats, error) {
	// Try to get from cache
	cached, found := s.getFromCache(ctx, userID, DomainGeneral)
	if found {
		if stats, ok := cached.(*GeneralStats); ok {
			return stats, nil
		}
	}

	// Get from repository
	stats, err := s.repo.GetGeneralStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.setInCache(ctx, userID, DomainGeneral, stats)
	return stats, nil
}

// GetFinancialStats retrieves financial statistics
func (s *service) GetFinancialStats(ctx context.Context, userID string) (*FinancialStats, error) {
	// Try to get from cache
	cached, found := s.getFromCache(ctx, userID, DomainFinancial)
	if found {
		if stats, ok := cached.(*FinancialStats); ok {
			return stats, nil
		}
	}

	// Get from repository
	stats, err := s.repo.GetFinancialStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.setInCache(ctx, userID, DomainFinancial, stats)
	return stats, nil
}

// GetStorageStats retrieves storage statistics
func (s *service) GetStorageStats(ctx context.Context, userID string) (*StorageStats, error) {
	// Try to get from cache
	cached, found := s.getFromCache(ctx, userID, DomainStorage)
	if found {
		if stats, ok := cached.(*StorageStats); ok {
			return stats, nil
		}
	}

	// Get from repository
	stats, err := s.repo.GetStorageStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.setInCache(ctx, userID, DomainStorage, stats)
	return stats, nil
}

// GetInventoryStats retrieves inventory statistics
func (s *service) GetInventoryStats(ctx context.Context, userID string) (*InventoryStats, error) {
	// Try to get from cache
	cached, found := s.getFromCache(ctx, userID, DomainInventory)
	if found {
		if stats, ok := cached.(*InventoryStats); ok {
			return stats, nil
		}
	}

	// Get from repository
	stats, err := s.repo.GetInventoryStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.setInCache(ctx, userID, DomainInventory, stats)
	return stats, nil
}

// GetWishlistStats retrieves wishlist statistics
func (s *service) GetWishlistStats(ctx context.Context, userID string) (*WishlistStats, error) {
	// Try to get from cache
	cached, found := s.getFromCache(ctx, userID, DomainWishlist)
	if found {
		if stats, ok := cached.(*WishlistStats); ok {
			return stats, nil
		}
	}

	// Get from repository
	stats, err := s.repo.GetWishlistStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.setInCache(ctx, userID, DomainWishlist, stats)
	return stats, nil
}

// InvalidateDomain removes a specific domain from cache for a user
func (s *service) InvalidateDomain(ctx context.Context, userID string, domain string) error {
	cacheKey := fmt.Sprintf(CacheKeyFormat, userID, domain)
	err := s.cache.DeleteCacheKey(ctx, cacheKey)
	if err != nil {
		s.logger.Error("Failed to invalidate analytics domain", map[string]any{
			"domain":    domain,
			"userID":    userID,
			"error":     err,
			"cache_key": cacheKey,
		})
		return fmt.Errorf("failed to invalidate %s domain: %w", domain, err)
	}

	s.logger.Debug("Invalidated analytics domain", map[string]any{
		"domain":    domain,
		"userID":    userID,
		"cache_key": cacheKey,
	})
	return nil
}

// InvalidateDomains removes multiple domains from cache for a user
func (s *service) InvalidateDomains(ctx context.Context, userID string, domains []string) error {
	var errs []error

	// Use index-based loop
	for i := 0; i < len(domains); i++ {
			domain := domains[i]
			err := s.InvalidateDomain(ctx, userID, domain)
			if err != nil {
					errs = append(errs, err)
			}
	}

	if len(errs) > 0 {
			return fmt.Errorf("failed to invalidate some domains: %v", errs)
	}
	return nil
}

// getFromCache tries to get data from cache
func (s *service) getFromCache(ctx context.Context, userID string, domain string) (any, bool) {
	cacheKey := fmt.Sprintf(CacheKeyFormat, userID, domain)
	var result any
	found, err := s.cache.GetCachedResults(ctx, cacheKey, &result)
	if err != nil {
		s.logger.Error("Failed to get data from cache", map[string]any{
			"domain":    domain,
			"userID":    userID,
			"error":     err,
			"cache_key": cacheKey,
		})
		return nil, false
	}
	return result, found
}

// setInCache stores data in cache
func (s *service) setInCache(ctx context.Context, userID string, domain string, data any) {
	cacheKey := fmt.Sprintf(CacheKeyFormat, userID, domain)
	err := s.cache.SetCachedResults(ctx, cacheKey, data)
	if err != nil {
		s.logger.Error("Failed to cache analytics data", map[string]any{
			"domain":    domain,
			"userID":    userID,
			"error":     err,
			"cache_key": cacheKey,
		})
	} else {
		s.logger.Debug("Cached analytics data", map[string]any{
			"domain":    domain,
			"userID":    userID,
			"cache_key": cacheKey,
		})
	}
}