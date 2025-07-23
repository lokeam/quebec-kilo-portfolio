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


// Helper function to calculate current year total spending (subscriptions + one-time purchases)
func (stc *SpendTrackingCalculator) calculateCurrentYearTotalSpending(
	userID string,
	currentYear int,
) (float64, error) {
	stc.logger.Debug("calculateCurrentYearTotalSpending called", map[string]any{
		"userID":      userID,
		"currentYear": currentYear,
	})

	// Calculate subscription costs for the year
	subscriptionCosts, err := stc.calculateCurrentYearSubscriptionCosts(userID, currentYear)
	if err != nil {
		stc.logger.Error("Failed to calculate subscription costs", map[string]any{
			"error":  err,
			"userID": userID,
		})
		subscriptionCosts = 0.0
	}

	// Calculate one-time purchase costs for the year
	oneTimeCosts, err := stc.calculateCurrentYearOneTimeCosts(userID, currentYear)
	if err != nil {
		stc.logger.Error("Failed to calculate one-time costs", map[string]any{
			"error":  err,
			"userID": userID,
		})
		oneTimeCosts = 0.0
	}

	totalSpending := subscriptionCosts + oneTimeCosts

	stc.logger.Debug("calculateCurrentYearTotalSpending completed", map[string]any{
		"userID":          userID,
		"currentYear":     currentYear,
		"subscriptionCosts": subscriptionCosts,
		"oneTimeCosts":    oneTimeCosts,
		"totalSpending":   totalSpending,
	})

	return totalSpending, nil
}

// Helper function to calculate current year one-time purchase costs
func (stc *SpendTrackingCalculator) calculateCurrentYearOneTimeCosts(
	userID string,
	currentYear int,
) (float64, error) {
	stc.logger.Debug("calculateCurrentYearOneTimeCosts called", map[string]any{
		"userID":      userID,
		"currentYear": currentYear,
	})

	// Get all one-time purchases for the current year
	var oneTimePurchases []models.SpendTrackingOneTimePurchaseDB
	err := stc.dbAdapter.db.SelectContext(
		context.Background(),
		&oneTimePurchases,
		`SELECT otp.*, sc.media_type as media_type
		FROM one_time_purchases otp
		LEFT JOIN spending_categories sc ON otp.spending_category_id = sc.id
		WHERE otp.user_id = $1
		AND EXTRACT(YEAR FROM purchase_date) = $2
		ORDER BY purchase_date DESC`,
		userID,
		currentYear,
	)
	if err != nil {
		stc.logger.Error("Failed to get one-time purchases for year", map[string]any{
			"error":  err,
			"userID": userID,
			"year":   currentYear,
		})
		return 0.0, fmt.Errorf("error getting one-time purchases: %w", err)
	}

	// Calculate total one-time purchases for the year
	totalOneTimeCosts := 0.0
	for _, purchase := range oneTimePurchases {
		if purchase.PurchaseDate.Year() == currentYear {
			totalOneTimeCosts += purchase.Amount
			stc.logger.Debug("Added one-time purchase to yearly total", map[string]any{
				"purchaseTitle": purchase.Title,
				"purchaseAmount": purchase.Amount,
				"currentYear": currentYear,
			})
		}
	}

	stc.logger.Debug("calculateCurrentYearOneTimeCosts completed", map[string]any{
		"userID":            userID,
		"currentYear":       currentYear,
		"totalOneTimeCosts": totalOneTimeCosts,
	})

	return totalOneTimeCosts, nil
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

	// Step 3: Populate result map with historical total spending data
	for _, aggregate := range yearlyAggregates {
		year := aggregate.Year

		// Only include years within our 3-year window
		if year >= currentYear-2 && year <= currentYear {
				// Use total amount (subscription + one-time) instead of just subscription
				result[year] = aggregate.TotalAmount

				stc.logger.Debug("Added historical total spending data", map[string]any{
						"year": year,
						"totalAmount": aggregate.TotalAmount,
						"subscriptionAmount": aggregate.SubscriptionAmount,
						"oneTimeAmount": aggregate.OneTimeAmount,
						"userID": userID,
				})
		}
	}

	if result[currentYear] == 0.0 {
		// If we don't have historical data for current year, calculate it dynamically
		currentYearTotal, err := stc.calculateCurrentYearTotalSpending(userID, currentYear)
		if err != nil {
				stc.logger.Error("Failed to calculate current year total spending", map[string]any{
						"error":  err,
						"userID": userID,
						"year":   currentYear,
				})
				// Keep as 0.0 if calculation fails
		} else {
				result[currentYear] = currentYearTotal

				stc.logger.Debug("Calculated dynamic total spending for current year", map[string]any{
						"year": currentYear,
						"totalAmount": currentYearTotal,
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
