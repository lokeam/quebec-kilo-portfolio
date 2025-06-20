package spend_tracking

import (
	"context"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/shared/utils"
	"github.com/lokeam/qko-beta/internal/types"
)

// Helper functions to calculate current year and yearly subscription costs
func (stc *SpendTrackingCalculator) calculateCurrentYearSubscriptionCosts(
	userID string,
	currentYear int,
) (float64, error) {
	stc.logger.Debug("calculateCurrentYearSubscriptionCosts called", map[string]any{
		"userID":      userID,
		"currentYear": currentYear,
	})

	// Get active subscriptions
	var subscriptions []models.SpendTrackingLocationDB
	err := stc.dbAdapter.db.SelectContext(
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
			return 0.0, fmt.Errorf("error getting active subscriptions: %w", err)
	}

	// Calculate total subscription costs for the entire year
	totalYearlyCost := 0.0
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

			// Calculate yearly cost based on billing cycle
			yearlyCost, err := stc.calculateSubscriptionYearlyCost(subscriptionDB, currentYear)
			if err != nil {
					stc.logger.Error("Failed to calculate yearly cost for subscription", map[string]any{
							"error":         err,
							"subscriptionID": subscription.ID,
							"currentYear":   currentYear,
					})
					continue // Skip this subscription if calculation fails
			}

			totalYearlyCost += yearlyCost

			stc.logger.Debug("Added subscription yearly cost", map[string]any{
					"subscriptionName": subscription.Name,
					"yearlyCost": yearlyCost,
					"currentYear": currentYear,
			})
	}

	stc.logger.Debug("calculateCurrentYearSubscriptionCosts completed", map[string]any{
		"userID": userID,
		"currentYear": currentYear,
		"totalYearlyCost": totalYearlyCost,
	})

	return totalYearlyCost, nil
}


func (stc *SpendTrackingCalculator) calculateSubscriptionYearlyCost(
	subscription models.SpendTrackingSubscriptionDB,
	targetYear int,
) (float64, error) {
	stc.logger.Debug("calculateSubscriptionYearlyCost called", map[string]any{
		"subscriptionID": subscription.ID,
		"targetYear":     targetYear,
	})

	// Calculate how many times this subscription will be charged in the target year
	paymentCount := 0

	// Start from whichever is later; the anchor date or beginning of target year
	startDate := subscription.AnchorDate
	if startDate.Year() < targetYear {
			// If anchor date is before target year, start from beginning of target year
			startDate = time.Date(targetYear, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	// If anchor date is after target year, there are no payments in this year
	if startDate.Year() > targetYear {
		return 0.0, nil
	}

	// Calculate billing cycle in months <---- NOTE: SEPARATE THIS OUT INTO A UTILITY FN
	billingCycleMonths, err := utils.GetBillingCycleMonths(subscription.BillingCycle)
	if err != nil {
		return 0.0, fmt.Errorf("error getting billing cycle months: %w", err)
	}

	// Calculate payments for the target year
	currentDate := startDate
	endOfYear := time.Date(targetYear, 12, 31, 23, 59, 59, 999999999, time.UTC)

	for currentDate.Before(endOfYear) || currentDate.Equal(endOfYear) {
			if currentDate.Year() == targetYear {
					paymentCount++
			}

			// Move to next payment date
			currentDate = currentDate.AddDate(0, billingCycleMonths, 0)
	}

	// Calculate total yearly cost
	yearlyCost := float64(paymentCount) * subscription.CostPerCycle

	stc.logger.Debug("calculateSubscriptionYearlyCost completed", map[string]any{
			"subscriptionID": subscription.ID,
			"targetYear":     targetYear,
			"paymentCount":   paymentCount,
			"costPerCycle":   subscription.CostPerCycle,
			"yearlyCost":     yearlyCost,
	})

	return yearlyCost, nil
}


// Core Subscription Logic
func (stc *SpendTrackingCalculator) CalculatePerSubscriptionYearlyTotals(
	userID string,
	subscriptionID string,
) ([]types.SingleYearlyTotalBFFResponseFINAL, error) {
	stc.logger.Debug("CalculatePerSubscriptionYearlyTotals called", map[string]any{
		"userID":          userID,
		"subscriptionID":  subscriptionID,
	})

	// Get subscription details
	var subscription models.SpendTrackingLocationDB
	err := stc.dbAdapter.db.GetContext(
			context.Background(),
			&subscription,
			GetSubscriptionByIDQuery,
			subscriptionID,
			userID,
	)
	if err != nil {
		stc.logger.Error("Failed to get subscription details", map[string]any{
				"error":         err,
				"subscriptionID": subscriptionID,
				"userID":        userID,
		})
		return nil, fmt.Errorf("error getting subscription details: %w", err)
	}

	// Convert to SpendTrackingSubscriptionDB for calculations
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

	// Calculate yearly totals for the last 3 years
	currentYear := time.Now().Year()
  yearlyTotals := make([]types.SingleYearlyTotalBFFResponseFINAL, 0, 3)

	for year := currentYear - 2; year <= currentYear; year++ {
			// Calculate yearly cost for this specific subscription
			yearlyCost, err := stc.calculateSubscriptionYearlyCost(subscriptionDB, year)
			if err != nil {
					stc.logger.Error("Failed to calculate yearly cost for subscription", map[string]any{
							"error":         err,
							"subscriptionID": subscriptionID,
							"year":          year,
					})
					// Continue with other years if one fails
					continue
			}

			// Create yearly total entry
			yearlyTotal := types.SingleYearlyTotalBFFResponseFINAL{
					Year:   year,
					Amount: yearlyCost,
			}
			yearlyTotals = append(yearlyTotals, yearlyTotal)

			stc.logger.Debug("Calculated yearly total for subscription", map[string]any{
					"subscriptionID": subscriptionID,
					"subscriptionName": subscription.Name,
					"year":          year,
					"yearlyCost":    yearlyCost,
			})
	}

	stc.logger.Debug("CalculatePerSubscriptionYearlyTotals completed", map[string]any{
		"userID":          userID,
		"subscriptionID":  subscriptionID,
		"yearlyTotals":    yearlyTotals,
		"totalYears":      len(yearlyTotals),
	})

	return yearlyTotals, nil
}

func (stc *SpendTrackingCalculator) IsSubscriptionDueInMonth(
	subscription models.SpendTrackingSubscriptionDB,
	targetMonth time.Time,
) (bool, error) {
	// Calculate months since anchor date
	stc.logger.Debug("IsSubscriptionDueInMonth called", map[string]any{
		"subscriptionID": subscription.LocationID,
		"billingCycle":   subscription.BillingCycle,
		"anchorDate":     subscription.AnchorDate,
		"targetMonth":    targetMonth,
	})

	// Normalize target month to first day for a consistent comparison
	targetMonthStart := time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, targetMonth.Location())

	// Use anchor date as the base for calculations
	anchorDate := subscription.AnchorDate

	// Calculate months since anchor date
	monthsSinceAnchor := int(
		(targetMonthStart.Year() - anchorDate.Year()) * 12 + int(targetMonthStart.Month() - anchorDate.Month()),
	)

	// Determine billing cycle in months
	billingCycleInMonths, err := utils.GetBillingCycleMonths(subscription.BillingCycle)
	if err != nil {
		return false, fmt.Errorf("error getting billing cycle months: %w", err)
	}

	isPaymentDue := monthsSinceAnchor % billingCycleInMonths == 0

	stc.logger.Debug("Subscription due calculation result", map[string]any{
		"subscriptionID":      subscription.LocationID,
		"monthsSinceAnchor":   monthsSinceAnchor,
		"billingCycleMonths":  billingCycleInMonths,
		"isDue":              isPaymentDue,
		"targetMonth":        targetMonth,
	})

	return isPaymentDue, nil
}
