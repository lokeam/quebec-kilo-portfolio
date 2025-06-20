package spend_tracking

import (
	"context"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

// Helper function to calculate average historical spending
func (stc *SpendTrackingCalculator) calculateAverageHistoricalSpending(
	monthlyAggregates []models.SpendTrackingMonthlyAggregateDB,
  currentYear int,
) float64 {
	totalSpending := 0.0
	monthCount := 0

	for _, aggregate := range monthlyAggregates {
		if aggregate.Year == currentYear && aggregate.TotalAmount > 0 {
			totalSpending += aggregate.TotalAmount
			monthCount++
		}
	}

	if monthCount == 0 {
		return totalSpending / float64(monthCount)
	}

	return 0.0
}


// Historical Analysis Logic
func (stc *SpendTrackingCalculator) CalculateThreeYearSubscriptionCosts(
	userID string,
	targetYear time.Time,
) (map[int]float64, error) {
	stc.logger.Debug("CalculateThreeYearSubscriptionCosts called", map[string]any{
		"userID":      userID,
		"targetYear":  targetYear,
	})

	// Get yearly spending aggregates for the last 3 years
	var yearlyAggregates []models.SpendTrackingYearlyAggregateDB
	err := stc.dbAdapter.db.SelectContext(
			context.Background(),
			&yearlyAggregates,
			GetYearlySpendingQuery,
			userID,
	)

	if err != nil {
		stc.logger.Error("Failed to get yearly spending aggregates", map[string]any{
				"error":  err,
				"userID": userID,
		})
		return nil, fmt.Errorf("error getting yearly spending aggregates: %w", err)
	}

	// Initialize result map for 3 years
	result := make(map[int]float64)
	currentYear := targetYear.Year()

	// Initialize with zero values for all 3 years
	for year := currentYear - 2; year <= currentYear; year++ {
			result[year] = 0.0
	}

	// Step 3: Populate result map with historical subscription data
	for _, aggregate := range yearlyAggregates {
		year := aggregate.Year

		// Only include years within our 3-year window
		if year >= currentYear-2 && year <= currentYear {
				result[year] = aggregate.SubscriptionAmount

				stc.logger.Debug("Added historical subscription data", map[string]any{
						"year": year,
						"subscriptionAmount": aggregate.SubscriptionAmount,
						"userID": userID,
				})
		}
	}

	if result[currentYear] == 0.0 {
		// If we don't have historical data for current year, calculate it dynamically
		currentYearSubscriptionTotal, err := stc.calculateCurrentYearSubscriptionCosts(userID, currentYear)
		if err != nil {
				stc.logger.Error("Failed to calculate current year subscription costs", map[string]any{
						"error":  err,
						"userID": userID,
						"year":   currentYear,
				})
				// Keep as 0.0 if calculation fails
		} else {
				result[currentYear] = currentYearSubscriptionTotal

				stc.logger.Debug("Calculated dynamic subscription costs for current year", map[string]any{
						"year": currentYear,
						"subscriptionAmount": currentYearSubscriptionTotal,
						"userID": userID,
				})
		}
	}

	// Log results
	stc.logger.Debug("CalculateThreeYearSubscriptionCosts completed", map[string]any{
		"userID": userID,
		"targetYear": targetYear,
		"result": result,
		"yearCount": len(result),
	})

	return result, nil
}
