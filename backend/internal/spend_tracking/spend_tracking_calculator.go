package spend_tracking

import (
	"context"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
)

type SpendTrackingCalculator struct {
	dbAdapter   SpendTrackingDbAdapter
	logger      interfaces.Logger
}

func NewSpendTrackingCalculator(appContext *appcontext.AppContext) (*SpendTrackingCalculator, error) {
	appContext.Logger.Debug("Creating SpendTracking Calculator", map[string]any{
		"appContext": appContext,
	})

	dbAdapter, err := NewSpendTrackingDbAdapter(appContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create SpendTrackingDbAdapter: %w", err)
	}

	return &SpendTrackingCalculator{
		dbAdapter: *dbAdapter,
		logger:    appContext.Logger,
	}, nil
}


// Interface methods (delegates work to other files)
func (stc *SpendTrackingCalculator) CalculateMonthlySubscriptionCosts(
	userID string,
	targetMonth time.Time,
) (float64, error) {
	stc.logger.Debug("CalculateMonthlySubscriptionCosts called", map[string]any{
		"userID": userID,
		"targetMonth": targetMonth,
	})

	// Get active subscriptions for user
	var activeSubscriptions []models.SpendTrackingLocationDB
	err := stc.dbAdapter.db.SelectContext(
		context.Background(),
		&activeSubscriptions,
		GetActiveSubscriptionsQuery,
		userID,
	)
	if err != nil {
		stc.logger.Error("Failed to get active subscriptions", map[string]any{
			"error": err,
			"userID": userID,
		})
		return 0.0, fmt.Errorf("error getting active subscriptions: %w", err)
	}

	// Calculate total subscription costs for target month
	totalSubscriptionCosts := 0.0
	for _, subscription := range activeSubscriptions {
		// Convert to SpendTrackingSubscriptionDB for calculation
		subscriptionDB := models.SpendTrackingSubscriptionDB{
			ID:               0, // Will be set from subscription data
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
		isPaymentDue, err := stc.IsSubscriptionDueInMonth(subscriptionDB, targetMonth)
		if err != nil {
				stc.logger.Error("Failed to check if subscription is due", map[string]any{
						"error":         err,
						"subscriptionID": subscription.ID,
						"targetMonth":   targetMonth,
				})
				continue // Skip this subscription if calculation fails
		}

		if isPaymentDue {
			totalSubscriptionCosts += subscription.CostPerCycle
			stc.logger.Debug("Subscription due in target month", map[string]any{
					"subscriptionID": subscription.ID,
					"subscriptionName": subscription.Name,
					"costPerCycle":   subscription.CostPerCycle,
					"targetMonth":    targetMonth,
			})
		}
	}

	stc.logger.Debug("CalculateMonthlySubscriptionCosts completed", map[string]any{
		"userID":     userID,
		"targetMonth": targetMonth,
		"totalCost":  totalSubscriptionCosts,
	})

	return totalSubscriptionCosts, nil
}

func (stc *SpendTrackingCalculator) CalculateMonthlyMinimumSpending(
	userID string,
	targetMonth time.Time,
) (float64, error) {
		stc.logger.Debug("CalculateMonthlyMinimumSpending called", map[string]any{
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
        return 0.0, fmt.Errorf("error getting one-time purchases: %w", err)
    }

		// Calculate total one-time purchases for the month
    oneTimeTotal := 0.0
    for _, purchase := range oneTimePurchases {
        // Check if purchase is in target month
        if purchase.PurchaseDate.Year() == targetMonth.Year() &&
           purchase.PurchaseDate.Month() == targetMonth.Month() {
            oneTimeTotal += purchase.Amount
            stc.logger.Debug("Added one-time purchase to monthly total", map[string]any{
                "purchaseTitle": purchase.Title,
                "purchaseAmount": purchase.Amount,
                "targetMonth": targetMonth,
            })
        }
    }

		// Calculate dynamic subscription costs for the month
    totalSubscriptionCosts, err := stc.CalculateMonthlySubscriptionCosts(userID, targetMonth)
    if err != nil {
        stc.logger.Error("Failed to calculate subscription costs", map[string]any{
            "error":  err,
            "userID": userID,
        })
        // Continue with one-time purchases only if subscription calculation fails
        totalSubscriptionCosts = 0.0
    }

		// Calculate total monthly minimum spending
    totalMinimumMonthlySpending := oneTimeTotal + totalSubscriptionCosts

		stc.logger.Debug("CalculateMonthlyMinimumSpending completed", map[string]any{
			"userID": userID,
			"targetMonth": targetMonth,
			"oneTimeTotal": oneTimeTotal,
			"totalSubscriptionCosts": totalSubscriptionCosts,
			"totalMinimumMonthlySpending": totalMinimumMonthlySpending,
		})

		return totalMinimumMonthlySpending, nil
}