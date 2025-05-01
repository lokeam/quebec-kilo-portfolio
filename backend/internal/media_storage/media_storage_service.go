package media_storage

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/interfaces"
)

// MediaStorageService defines the interface for retrieving and managing media storage statistics
// This service acts as a facade over the analytics service, providing a focused interface
// for media storage-related operations.
type MediaStorageService interface {
	// GetStorageStats retrieves storage statistics for a user, including both physical
	// and digital storage locations, their counts, and associated metadata.
	GetStorageStats(ctx context.Context, userID string) (*analytics.StorageStats, error)
}

// mediaStorageServiceImpl implements the MediaStorageService interface
// This implementation coordinates with the analytics service to provide storage statistics
// while handling error mapping and logging.
type mediaStorageServiceImpl struct {
	analyticsService analytics.Service
	logger           interfaces.Logger
}

// NewMediaStorageService creates a new instance of the media storage service
// This service is designed to work with the analytics service to provide storage statistics
// in a focused and efficient manner.
//
// Parameters:
//   - analyticsService: The analytics service to use for retrieving storage statistics
//   - logger: The logger to use for logging operations
//
// Returns:
//   - MediaStorageService: A new instance of the media storage service
//   - error: An error if the service could not be created
//
// Errors:
//   - ErrInvalidInput: If analyticsService or logger is nil
func NewMediaStorageService(analyticsService analytics.Service, logger interfaces.Logger) (MediaStorageService, error) {
	// Validate required dependencies
	if analyticsService == nil {
		return nil, fmt.Errorf("%w: analytics service is required", ErrInvalidInput)
	}
	if logger == nil {
		return nil, fmt.Errorf("%w: logger is required", ErrInvalidInput)
	}

	// Create and return the service
	service := &mediaStorageServiceImpl{
		analyticsService: analyticsService,
		logger:          logger,
	}

	// Test the analytics service connection
	ctx := context.Background()
	if _, err := analyticsService.GetStorageStats(ctx, "test-connection"); err != nil {
		// Log the error but don't fail the service creation
		// This allows the service to be created even if the analytics service is temporarily unavailable
		logger.Warn("Analytics service connection test failed", map[string]any{
			"error": err,
		})
	}

	return service, nil
}

// GetStorageStats retrieves storage statistics for a user
func (s *mediaStorageServiceImpl) GetStorageStats(ctx context.Context, userID string) (*analytics.StorageStats, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	s.logger.Debug("Getting storage stats", map[string]any{
		"userID": userID,
	})

	// Get storage stats from analytics service
	stats, err := s.analyticsService.GetStorageStats(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get storage stats", map[string]any{
			"userID": userID,
			"error":  err,
		})

		// Map analytics service errors to our error types
		switch {
		case err.Error() == "storage stats not found":
			return nil, ErrStorageStatsNotFound
		case err.Error() == "analytics service unavailable":
			return nil, ErrAnalyticsServiceUnavailable
		default:
			return nil, fmt.Errorf("failed to get storage stats: %w", err)
		}
	}

	if stats == nil {
		return nil, ErrStorageStatsNotFound
	}

	s.logger.Debug("Successfully retrieved storage stats", map[string]any{
		"userID": userID,
		"stats":  stats,
	})

	return stats, nil
}
