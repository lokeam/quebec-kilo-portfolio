package dashboard

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/types"
)

type DashboardService struct {
	dbAdapter      interfaces.DashboardDbAdapter
	cacheWrapper   interfaces.DashboardCacheWrapper
	validator      interfaces.DashboardValidator
	logger         interfaces.Logger
}

func NewDashboardService(
	appContext *appcontext.AppContext,
	dbAdapter interfaces.DashboardDbAdapter,
	cacheWrapper interfaces.DashboardCacheWrapper,
) (*DashboardService, error) {
	if dbAdapter == nil {
		return nil, fmt.Errorf("dbAdapter is required")
	}
	if cacheWrapper == nil {
		return nil, fmt.Errorf("cacheWrapper is required")
	}

	return &DashboardService{
		dbAdapter:     dbAdapter,
		cacheWrapper:  cacheWrapper,
		validator:     NewDashboardValidator(),
		logger:        appContext.Logger,
	}, nil
}

func (ds *DashboardService) GetDashboardBFFResponse(ctx context.Context, userID string) (types.DashboardBFFResponse, error) {
	// 1. Validate userID
	if err := ds.validator.ValidateUserID(userID); err != nil {
		return types.DashboardBFFResponse{}, fmt.Errorf("invalid user ID: %w", err)
	}

	// 2. Try to get from cache first
	cachedResponse, err := ds.cacheWrapper.GetCachedDashboardBFF(ctx, userID)
	if err == nil && cachedResponse.GameStats.Value > 0 {
		ds.logger.Debug("Cache hit for dashboard BFF response", map[string]any{
			"userID": userID,
		})
		return cachedResponse, nil
	}

	// 3. Cache miss, get from database
	ds.logger.Debug("Cache miss for dashboard BFF response, fetching from database", map[string]any{
		"userID": userID,
	})
	response, err := ds.dbAdapter.GetDashboardBFFResponse(ctx, userID)
	if err != nil {
		return types.DashboardBFFResponse{}, fmt.Errorf("failed to get dashboard BFF response: %w", err)
	}

	// 4. Cache the response
	if err := ds.cacheWrapper.SetCachedDashboardBFF(ctx, userID, response); err != nil {
		ds.logger.Error("Failed to cache dashboard BFF response", map[string]any{
			"error":  err,
			"userID": userID,
		})
	}

	return response, nil
}