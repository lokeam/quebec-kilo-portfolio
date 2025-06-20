package spend_tracking

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

// Business Intelligence Logic
func (stc *SpendTrackingCalculator) CalculatePercentageChange(
	userID string,
	currentMonth time.Time,
) (float64, error) {
	stc.logger.Debug("CalculatePercentageChange called", map[string]any{
		"userID":       userID,
		"currentMonth": currentMonth,
	})

	// Get current month total
	currentMonthTotal, err := stc.CalculateMonthlyMinimumSpending(userID, currentMonth)
	if err != nil {
			stc.logger.Error("Failed to calculate current month total", map[string]any{
					"error":  err,
					"userID": userID,
			})
			return 0.0, fmt.Errorf("error calculating current month total: %w", err)
	}

	// Get previous month total
	previousMonth := currentMonth.AddDate(0, -1, 0) // Go back one month
	previousMonthTotal, err := stc.CalculateMonthlyMinimumSpending(userID, previousMonth)
	if err != nil {
			stc.logger.Error("Failed to calculate previous month total", map[string]any{
					"error":  err,
					"userID": userID,
			})
			return 0.0, fmt.Errorf("error calculating previous month total: %w", err)
	}

	// Calculate percentage change
	var percentageChange float64
	if previousMonthTotal > 0 {
			percentageChange = ((currentMonthTotal - previousMonthTotal) / previousMonthTotal) * 100
	} else {
			// If previous month was 0, can't calculate percentage change
			if currentMonthTotal > 0 {
					percentageChange = 100.0 // 100% increase from 0
			} else {
					percentageChange = 0.0 // No change if both are 0
			}
	}
		stc.logger.Debug("CalculatePercentageChange completed", map[string]any{
			"userID":              userID,
			"currentMonth":        currentMonth,
			"currentMonthTotal":   currentMonthTotal,
			"previousMonth":       previousMonth,
			"previousMonthTotal":  previousMonthTotal,
			"percentageChange":    percentageChange,
	})

	return percentageChange, nil
}

func (stc *SpendTrackingCalculator) CalculateAnnualSpendingForecast(
	userID string,
	targetYear time.Time,
) (types.AnnualSpendingBFFResponseFINAL, error) {
	stc.logger.Debug("CalculateAnnualSpendingForecast called", map[string]any{
		"userID":      userID,
		"targetYear":  targetYear,
	})

	// Step 1: Get historical monthly spending data
	var monthlyAggregates []models.SpendTrackingMonthlyAggregateDB
	err := stc.dbAdapter.db.SelectContext(
			context.Background(),
			&monthlyAggregates,
			GetMonthlySpendingAggregatesQuery,
			userID,
	)
	if err != nil {
			stc.logger.Error("Failed to get monthly spending aggregates", map[string]any{
					"error":  err,
					"userID": userID,
			})
			return types.AnnualSpendingBFFResponseFINAL{}, fmt.Errorf("error getting monthly spending aggregates: %w", err)
	}

	// Create monthly expenditures array from Jan - Dec
	monthlyExpenditures := make([]types.MonthlyExpenditureBFFResponseFINAL, 12)
    for i := range monthlyExpenditures {
        monthlyExpenditures[i] = types.MonthlyExpenditureBFFResponseFINAL{
            Month:       time.Month(i + 1).String()[:3], // <--- NOTE:"Jan", "Feb", etc.
            Expenditure: 0.0,
        }
  }

	// Fill array with historical data for current year
	currentYear := targetYear.Year()
  currentMonth := int(time.Now().Month())

	for _, agg := range monthlyAggregates {
		if agg.Year == currentYear && int(agg.Month) <= currentMonth {
				monthIndex := int(agg.Month) - 1 // Convert to 0-based index
				if monthIndex >= 0 && monthIndex < 12 {
						// Use dynamic calculation for current and future months
						if agg.TotalAmount > 0 {
								monthlyExpenditures[monthIndex].Expenditure = agg.TotalAmount
						}
				}
		}
	}

	for monthIndex := currentMonth; monthIndex < 12; monthIndex++ {
		targetMonth := time.Date(currentYear, time.Month(monthIndex+1), 1, 0, 0, 0, 0, time.UTC)

		// Calculate dynamic monthly spending for future months
		monthlySpending, err := stc.CalculateMonthlyMinimumSpending(userID, targetMonth)
		if err != nil {
				stc.logger.Error("Failed to calculate dynamic spending for future month", map[string]any{
						"error":      err,
						"userID":    userID,
						"targetMonth": targetMonth,
				})
				// Use average of historical months as fallback
				monthlySpending = stc.calculateAverageHistoricalSpending(monthlyAggregates, currentYear)
		}

		monthlyExpenditures[monthIndex].Expenditure = monthlySpending
	}

	// Calculate median monthly spending
	medianMonthlyCost := stc.CalculateMedianMonthlyCost(monthlyExpenditures)

	// Format date range
	dateRange := fmt.Sprintf("January %d - January %d", currentYear, currentYear + 1)

	response := types.AnnualSpendingBFFResponseFINAL{
		DateRange: dateRange,
		MonthlyExpenditures: monthlyExpenditures,
		MedianMonthlyCost: medianMonthlyCost,
	}

	stc.logger.Debug("CalculateAnnualSpendingForecast completed", map[string]any{
		"userID": userID,
		"targetYear": targetYear,
		"response": response,
	})

	return response, nil
}

func (stc *SpendTrackingCalculator) CalculateCurrentMonthAggregation(
	userID string,
	targetMonth time.Time,
) (types.SpendTrackingCalculatorCurrentMonthData, error) {
	stc.logger.Debug("CalculateCurrentMonthAggregation called", map[string]any{
		"userID":      userID,
		"targetMonth": targetMonth,
	})
	stc.logger.Debug("CalculateCurrentMonthAggregation called", map[string]any{
		"userID":      userID,
		"targetMonth": targetMonth,
	})


	// Get one-time purchases for the target month
	var oneTimePurchases []models.SpendTrackingOneTimePurchaseDB
	err := stc.dbAdapter.db.SelectContext(
			context.Background(),
			&oneTimePurchases,
			GetCurrentMonthOneTimePurchasesQuery,
			userID,
	)
	if err != nil {
		stc.logger.Error("Failed to get one-time purchases", map[string]any{
				"error":  err,
				"userID": userID,
		})
		return types.SpendTrackingCalculatorCurrentMonthData{}, fmt.Errorf("error getting one-time purchases: %w", err)
	}

	// Aggregate one-time purchases by category
	categoryMap := make(map[string]float64)
  var spendingItems []types.SpendTrackingCalculatorSpendingItem

	for _, purchase := range oneTimePurchases {
		// Check if purchase is in target month
		if purchase.PurchaseDate.Year() == targetMonth.Year() &&
			 purchase.PurchaseDate.Month() == targetMonth.Month() {

				// Add to category total
				categoryName := purchase.MediaType
				categoryMap[categoryName] += purchase.Amount

				// Create spending item
				spendingItem := types.SpendTrackingCalculatorSpendingItem{
						SpendingCategoryID:   categoryName,
						SpendingItemName:     purchase.Title,
						SpendingItemAmount:   purchase.Amount,
						SpendingItemCategory: categoryName,
				}
				spendingItems = append(spendingItems, spendingItem)

				stc.logger.Debug("Added one-time purchase to aggregation", map[string]any{
						"purchaseTitle": purchase.Title,
						"purchaseAmount": purchase.Amount,
						"category": categoryName,
						"targetMonth": targetMonth,
				})
		}
	}

	// Get active subscriptions for the month
	var subscriptions []models.SpendTrackingLocationDB
	err = stc.dbAdapter.db.SelectContext(
			context.Background(),
			&subscriptions,
			GetActiveSubscriptionsQuery,
			userID,
	)

	if err != nil {
		stc.logger.Error("Failed to get active subscriptions", map[string]any{
				"error":  err,
				"userID": userID,
		})
		return types.SpendTrackingCalculatorCurrentMonthData{}, fmt.Errorf("error getting active subscriptions: %w", err)
	}

	// Add subscription costs to categories and items
	for _, subscription := range subscriptions {
		// Convert to SpendTrackingSubscriptionDB for calculation
		subscriptionDB := models.SpendTrackingSubscriptionDB{
				ID:               0,
				LocationID:       subscription.ID,
				BillingCycle:     subscription.BillingCycle,
				CostPerCycle:     subscription.CostPerCycle,
				AnchorDate:       subscription.AnchorDate,
				LastPaymentDate:  subscription.LastPaymentDate,
				NextPaymentDate:  subscription.NextPaymentDate,
				PaymentMethod:    subscription.SubscriptionPaymentMethod,
				CreatedAt:        subscription.CreatedAt,
				UpdatedAt:        subscription.UpdatedAt,
		}

		// Check if subscription is due in target month
		isSubscriptionDue, err := stc.IsSubscriptionDueInMonth(subscriptionDB, targetMonth)
		if err != nil {
				stc.logger.Error("Failed to check if subscription is due", map[string]any{
						"error":         err,
						"subscriptionID": subscription.ID,
						"targetMonth":   targetMonth,
				})
				continue // Skip this subscription if calculation fails
		}

		if isSubscriptionDue {
				// Add to subscription category
				categoryName := "subscription"
				categoryMap[categoryName] += subscription.CostPerCycle

				// Create spending item for subscription
				spendingItem := types.SpendTrackingCalculatorSpendingItem{
						SpendingCategoryID:   categoryName,
						SpendingItemName:     subscription.Name,
						SpendingItemAmount:   subscription.CostPerCycle,
						SpendingItemCategory: categoryName,
				}
				spendingItems = append(spendingItems, spendingItem)

				stc.logger.Debug("Added subscription to aggregation", map[string]any{
						"subscriptionName": subscription.Name,
						"subscriptionAmount": subscription.CostPerCycle,
						"category": categoryName,
						"targetMonth": targetMonth,
				})
		}
	}

	// Convert category map to spending categories array
	var spendingCategories []types.SpendTrackingCalculatorSpendingCategory
	for categoryName, categoryValue := range categoryMap {
			spendingCategory := types.SpendTrackingCalculatorSpendingCategory{
					SpendingCategoryName:  categoryName,
					SpendingCategoryValue: categoryValue,
			}
			spendingCategories = append(spendingCategories, spendingCategory)

			stc.logger.Debug("Created spending category", map[string]any{
					"categoryName": categoryName,
					"categoryValue": categoryValue,
			})
	}

	 // Calculate total monthly spending
	 totalMonthlySpending := 0.0
	 for _, category := range spendingCategories {
			 totalMonthlySpending += category.SpendingCategoryValue
	 }

	 // Build and return the response
	 response := types.SpendTrackingCalculatorCurrentMonthData{
			 TotalMonthlySpending: totalMonthlySpending,
			 SpendingCategories:   spendingCategories,
			 SpendingItems:        spendingItems,
	 }

	 stc.logger.Debug("CalculateCurrentMonthAggregation completed", map[string]any{
			 "userID":              userID,
			 "targetMonth":         targetMonth,
			 "totalMonthlySpending": totalMonthlySpending,
			 "categoryCount":       len(spendingCategories),
			 "itemCount":           len(spendingItems),
	 })

	 return response, nil
}

func (stc *SpendTrackingCalculator) CalculateMedianMonthlyCost(
	monthlyExpenditures []types.MonthlyExpenditureBFFResponseFINAL,
) float64 {
	stc.logger.Debug("CalculateMedianMonthlyCost called", map[string]any{
		"expenditureCount": len(monthlyExpenditures),
	})

	// Step 1: Handle edge cases
	if len(monthlyExpenditures) == 0 {
		stc.logger.Debug("No monthly expenditures provided, returning 0", map[string]any{})
		return 0.0
	}

	if len(monthlyExpenditures) == 1 {
		stc.logger.Debug("Single monthly expenditure, returning its value", map[string]any{
				"value": monthlyExpenditures[0].Expenditure,
		})
		return monthlyExpenditures[0].Expenditure
	}

	// Step 2: Extract expenditure values and sort them
	expenditures := make([]float64, len(monthlyExpenditures))
	for i, expenditure := range monthlyExpenditures {
			expenditures[i] = expenditure.Expenditure
	}

	// Sort expenditures in ascending order
	sort.Float64s(expenditures)

	stc.logger.Debug("Sorted expenditures", map[string]any{
			"expenditures": expenditures,
			"count":        len(expenditures),
	})

	// Step 3: Calculate median
	var median float64
	count := len(expenditures)
	middleIndex := count / 2

	if count%2 == 0 {
		// Even number of items: median is average of two middle values
		median = (expenditures[middleIndex-1] + expenditures[middleIndex]) / 2.0

		stc.logger.Debug("Even number of expenditures, calculating average of middle values", map[string]any{
				"middleIndex1": middleIndex - 1,
				"middleIndex2": middleIndex,
				"value1":       expenditures[middleIndex-1],
				"value2":       expenditures[middleIndex],
				"median":       median,
		})
	} else {
		// Odd number of items: median is the middle value
		median = expenditures[middleIndex]

		stc.logger.Debug("Odd number of expenditures, using middle value", map[string]any{
				"middleIndex": middleIndex,
				"median":      median,
		})
	}

	stc.logger.Debug("CalculateMedianMonthlyCost completed", map[string]any{
		"expenditureCount": len(expenditures),
		"median":          median,
		"minValue":        expenditures[0],
		"maxValue":        expenditures[count-1],
	})

	return median
}
