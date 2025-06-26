package digital

import (
	"context"
	"fmt"
	"html"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

// BFF Response
func (da *DigitalDbAdapter) GetAllDigitalLocationsBFF(
	ctx context.Context,
	userID string,
) (types.DigitalLocationsBFFResponse, error) {
	da.logger.Debug("GetAllDigitalLocationsBFF called", map[string]any{"userID": userID})

	// Start transaction
	tx, err := da.db.BeginTxx(ctx, nil)
	if err != nil {
			return types.DigitalLocationsBFFResponse{
				DigitalLocations: []types.SingleDigitalLocationBFFResponse{},
			}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Get all digital locations with basic info
	var locationsDB []models.DigitalLocationBFFDB
	if err := tx.SelectContext(
		ctx,
		&locationsDB,
		GetAllDigitalLocationsBFFQuery,
		userID,
	); err != nil {
			return types.DigitalLocationsBFFResponse{
				DigitalLocations: []types.SingleDigitalLocationBFFResponse{},
			}, fmt.Errorf("failed to get digital locations: %w", err)
	}

	// For each location, get its games
	digitalLocations := make([]types.SingleDigitalLocationBFFResponse, len(locationsDB))
	for i, locationDB := range locationsDB {
			// Get games for this location
			var gamesDB []models.DigitalLocationGameDB
			if err := tx.SelectContext(
				ctx,
				&gamesDB,
				GetDigitalLocationGamesBFFQuery,
				locationDB.ID,
				userID,
			); err != nil {
					return types.DigitalLocationsBFFResponse{
						DigitalLocations: []types.SingleDigitalLocationBFFResponse{},
					}, fmt.Errorf("failed to get games for location %s: %w", locationDB.ID, err)
			}

			// Transform to response format
			digitalLocations[i] = da.transformDigitalLocationDBToResponse(locationDB, gamesDB)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
			return types.DigitalLocationsBFFResponse{
				DigitalLocations: []types.SingleDigitalLocationBFFResponse{},
			}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return types.DigitalLocationsBFFResponse{
			DigitalLocations: digitalLocations,
	}, nil
}

func (da *DigitalDbAdapter) transformDigitalLocationDBToResponse(
	db models.DigitalLocationBFFDB,
	gamesDB []models.DigitalLocationGameDB,
) types.SingleDigitalLocationBFFResponse {
	// Calculate monthly cost based on billing cycle
	monthlyCost := 0.0
	if db.BillingCycle != "" && db.CostPerCycle > 0 {
			switch db.BillingCycle {
			case "1 month":
					monthlyCost = db.CostPerCycle
			case "3 month":
					monthlyCost = db.CostPerCycle / 3
			case "6 month":
					monthlyCost = db.CostPerCycle / 6
			case "12 month":
					monthlyCost = db.CostPerCycle / 12
			}
	}

	// Transform games
	storedGames := make([]types.DigitalLocationGameResponse, len(gamesDB))
	for i, gameDB := range gamesDB {
			storedGames[i] = types.DigitalLocationGameResponse{
					ID:              gameDB.ID,
					Name:            gameDB.Name,
					Platform:        gameDB.Platform,
					IsUniqueCopy:    gameDB.IsUniqueCopy,
					HasPhysicalCopy: gameDB.HasPhysicalCopy,
			}
	}

	return types.SingleDigitalLocationBFFResponse{
			ID: db.ID,
			Name:              html.UnescapeString(db.Name),
			IsSubscription:    db.IsSubscription,
			IsActive:          db.IsActive,
			URL:               db.URL,
			PaymentMethod:     db.PaymentMethod,
			MonthlyCost:       monthlyCost,
			BillingCycle:      db.BillingCycle,
			CostPerCycle:      db.CostPerCycle,
			NextPaymentDate:   db.NextPaymentDate,
			ItemCount:         db.StoredItems,
			StoredGames:       storedGames,
			CreatedAt:         db.CreatedAt,
			UpdatedAt:         db.UpdatedAt,
	}
}