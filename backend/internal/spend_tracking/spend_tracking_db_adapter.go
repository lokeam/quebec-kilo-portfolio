package spend_tracking

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)


type SpendTrackingDbAdapter struct {
	calculator  *SpendTrackingCalculator
	db          *sqlx.DB
	logger      interfaces.Logger
}

func NewSpendTrackingDbAdapter(
	appContext *appcontext.AppContext,
) (*SpendTrackingDbAdapter, error) {
	// Log that we're creating the SpendTrackingDbAdapter
	appContext.Logger.Debug("Creating LibraryDbAdapter", map[string]any{"appContext": appContext})

	// Use shared DB pool
	db := appContext.DB

	// Create the adapter first
	adapter := &SpendTrackingDbAdapter{
		db:      db,
		logger:  appContext.Logger,
	}

	// Create the spend tracking calculator with the adapter
	calculator, err := NewSpendTrackingCalculator(appContext, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create spend tracking calculator: %w", err)
	}

	// Set the calculator on the adapter
	adapter.calculator = calculator

	return adapter, nil
}

// --- HELPER METHODS ---
func (sta *SpendTrackingDbAdapter) isDateInMonth(date, targetMonth time.Time) bool {
	return date.Year() == targetMonth.Year() && date.Month() == targetMonth.Month()
}

func (sta *SpendTrackingDbAdapter) calculateYearlySpendingForSubscription(
	subscription models.SpendTrackingLocationDB,
) []types.SingleYearlyTotalBFFResponseFINAL {
	currentYear := time.Now().Year()

	// Calculate annual amount based on billing cycle
	annualAmount := subscription.CostPerCycle
	switch subscription.BillingCycle {
	case "1 month":
			annualAmount = subscription.CostPerCycle * 12
	case "3 month":
			annualAmount = subscription.CostPerCycle * 4
	case "6 month":
			annualAmount = subscription.CostPerCycle * 2
	case "12 month":
			annualAmount = subscription.CostPerCycle
	}

	return []types.SingleYearlyTotalBFFResponseFINAL{
			{Year: currentYear - 2, Amount: annualAmount},
			{Year: currentYear - 1, Amount: annualAmount},
			{Year: currentYear, Amount: annualAmount},
	}
}

// --- TRANSFORMATION LOGIC ---
func (sta *SpendTrackingDbAdapter) transformCalculatorCategoriesToBFFResponse(
	calculatorCategories []types.SpendTrackingCalculatorSpendingCategory,
) []types.SpendingCategoryBFFResponseFINAL {
	bffCategories := make([]types.SpendingCategoryBFFResponseFINAL,len(calculatorCategories))
	for i, category := range calculatorCategories {
		bffCategories[i] = types.SpendingCategoryBFFResponseFINAL{
			Name:  category.SpendingCategoryName,
			Value: category.SpendingCategoryValue,
		}
	}
	return bffCategories
}

func (sta *SpendTrackingDbAdapter) transformThreeYearTotalsToBFFResponse(
	threeYearTotals map[int]float64,
) types.AllYearlyTotalsBFFResponseFINAL {
	currentYear := time.Now().Year()

    subscriptionTotal := make([]types.SingleYearlyTotalBFFResponseFINAL, 0, 3)
    oneTimeTotal := make([]types.SingleYearlyTotalBFFResponseFINAL, 0, 3)
    combinedTotal := make([]types.SingleYearlyTotalBFFResponseFINAL, 0, 3)

    for year := currentYear - 2; year <= currentYear; year++ {
        totalAmount := threeYearTotals[year]

        // For now, we're using total amounts for all categories
        // In the future, we could separate subscription vs one-time amounts
        subscriptionTotal = append(subscriptionTotal, types.SingleYearlyTotalBFFResponseFINAL{
            Year:   year,
            Amount: totalAmount,
        })

        oneTimeTotal = append(oneTimeTotal, types.SingleYearlyTotalBFFResponseFINAL{
            Year:   year,
            Amount: totalAmount,
        })

        combinedTotal = append(combinedTotal, types.SingleYearlyTotalBFFResponseFINAL{
            Year:   year,
            Amount: totalAmount,
        })
    }

    return types.AllYearlyTotalsBFFResponseFINAL{
        SubscriptionTotal: subscriptionTotal,
        OneTimeTotal:      oneTimeTotal,
        CombinedTotal:     combinedTotal,
    }
}

func (sta *SpendTrackingDbAdapter) buildTransactionArraysFromDatabaseRecords(
	oneTimePurchases []models.SpendTrackingOneTimePurchaseDB,
	subscriptions []models.SpendTrackingLocationDB,
	currentMonth time.Time,
) ([]types.SpendingItemBFFResponseFINAL, []types.SpendingItemBFFResponseFINAL, []types.SpendingItemBFFResponseFINAL) {
	sta.logger.Debug("transformDetailedTransactionsToBFF called", map[string]any{
			"oneTimePurchaseCount": len(oneTimePurchases),
			"subscriptionCount":    len(subscriptions),
			"currentMonth":         currentMonth,
	})

	var currentTotalThisMonth []types.SpendingItemBFFResponseFINAL
	var oneTimeThisMonth []types.SpendingItemBFFResponseFINAL
	var recurringNextMonth []types.SpendingItemBFFResponseFINAL

	// Transform one-time purchases with real transaction data
	for _, purchase := range oneTimePurchases {
			// Check if purchase made in current month
			if sta.isDateInMonth(purchase.PurchaseDate, currentMonth) {
					transaction := types.SpendingItemBFFResponseFINAL{
							ID:                   fmt.Sprintf("one-%d", purchase.ID),
							Title:                purchase.Title,
							Amount:               purchase.Amount,
							SpendTransactionType: "one time purchase",
							PaymentMethod:        purchase.PaymentMethod,
							MediaType:            purchase.MediaType,
							CreatedAt:            purchase.CreatedAt.Unix(),
							UpdatedAt:            purchase.UpdatedAt.Unix(),
							IsActive:             true,
							IsDigital:            purchase.IsDigital,
							IsWishlisted:         purchase.IsWishlisted,
							PurchaseDate:         purchase.PurchaseDate.Unix(),
					}

					// Add to current month total (filtered transactions)
					currentTotalThisMonth = append(currentTotalThisMonth, transaction)

					// Add to one-time purchases array
					oneTimeThisMonth = append(oneTimeThisMonth, transaction)

					sta.logger.Debug("Added one-time purchase to detailed transactions", map[string]any{
							"purchaseID":    purchase.ID,
							"title":         purchase.Title,
							"amount":        purchase.Amount,
							"paymentMethod": purchase.PaymentMethod,
					})
			}
	}

	// Transform subscriptions with real transaction data
	for _, subscription := range subscriptions {
			// ✅ LEVERAGE CALCULATOR LOGIC DIRECTLY
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

			// ✅ USE CALCULATOR'S METHOD DIRECTLY
			isDue, err := sta.calculator.IsSubscriptionDueInMonth(subscriptionDB, currentMonth)
			if err != nil {
					sta.logger.Error("Failed to check if subscription is due", map[string]any{
							"error":         err,
							"subscriptionID": subscription.ID,
							"currentMonth":   currentMonth,
					})
					continue // Skip this subscription if calculation fails
			}

			if isDue {
					transaction := types.SpendingItemBFFResponseFINAL{
							ID:                   fmt.Sprintf("sub-%s", subscription.ID),
							Title:                subscription.Name,
							Amount:               subscription.CostPerCycle,
							SpendTransactionType: "subscription",
							PaymentMethod:        subscription.SubscriptionPaymentMethod,
							MediaType:            "subscription",
							Provider:             strings.ToLower(subscription.Name),
							CreatedAt:            subscription.CreatedAt.Unix(),
							UpdatedAt:            subscription.UpdatedAt.Unix(),
							IsActive:             subscription.IsActive,
							BillingCycle:         subscription.BillingCycle,
							NextBillingDate:      subscription.NextPaymentDate.Unix(),
							YearlySpending:       sta.calculateYearlySpendingForSubscription(subscription),
					}

					// Add to current month total (filtered transactions)
					currentTotalThisMonth = append(currentTotalThisMonth, transaction)

					// Add to recurring subscriptions array
					recurringNextMonth = append(recurringNextMonth, transaction)

					sta.logger.Debug("Added subscription to detailed transactions", map[string]any{
							"subscriptionID": subscription.ID,
							"title":          subscription.Name,
							"amount":         subscription.CostPerCycle,
							"paymentMethod":  subscription.SubscriptionPaymentMethod,
					})
			}
	}

	sta.logger.Debug("transformDetailedTransactionsToBFF completed", map[string]any{
			"currentTotalCount": len(currentTotalThisMonth),
			"oneTimeCount":      len(oneTimeThisMonth),
			"recurringCount":    len(recurringNextMonth),
	})

	return currentTotalThisMonth, oneTimeThisMonth, recurringNextMonth
}



// --- MAIN RESPONSE LOGIC -- GET - Send backend for frontend Spend Tracking Response ---
func (sta *SpendTrackingDbAdapter) GetSpendTrackingBFFResponse(
	ctx context.Context,
	userID string,
) (types.SpendTrackingBFFResponseFINAL, error) {
	sta.logger.Debug("GetSpendTrackingBFFResponse called", map[string]any{
		"userID": userID,
	})

	// Calculate Total Monthly Spending
	currentMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	currentMonthTotal, err := sta.calculator.CalculateMonthlyMinimumSpending(userID, currentMonth)
	if err != nil {
			sta.logger.Error("Failed to calculate current month total", map[string]any{
					"error":  err,
					"userID": userID,
			})
			return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error calculating current month total: %w", err)
	}

	lastMonth := currentMonth.AddDate(0, -1, 0)
	lastMonthTotal, err := sta.calculator.CalculateMonthlyMinimumSpending(userID, lastMonth)
	if err != nil {
			sta.logger.Error("Failed to calculate previous month total", map[string]any{
					"error":  err,
					"userID": userID,
			})
			return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error calculating previous month total: %w", err)
	}

	percentageChange, err := sta.calculator.CalculatePercentageChange(userID, currentMonth)
	if err != nil {
			sta.logger.Error("Failed to calculate percentage change", map[string]any{
					"error":  err,
					"userID": userID,
			})
			return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error calculating percentage change: %w", err)
	}

	// Calculate Total Annual Spending with dynamic forecasts
	annualSpendingForecast, err := sta.calculator.CalculateAnnualSpendingForecast(userID, currentMonth)
	if err != nil {
		sta.logger.Error("Failed to calculate annual spending forecast", map[string]any{
				"error":  err,
				"userID": userID,
		})
		return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error calculating annual spending forecast: %w", err)
	}

	// Calculate CurrentTotalThisMonth with dynamic data aggregation
	currentMonthAggregation, err := sta.calculator.CalculateCurrentMonthAggregation(userID, currentMonth)
	if err != nil {
			sta.logger.Error("Failed to calculate current month aggregation", map[string]any{
					"error":  err,
					"userID": userID,
			})
			return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error calculating current month aggregation: %w", err)
	}

	// Calculate YearlyTotals with dynamic subscription costs
	threeYearSubscriptionTotals, err := sta.calculator.CalculateThreeYearSubscriptionCosts(userID, currentMonth)
	if err != nil {
		sta.logger.Error("Failed to calculate three year subscription costs", map[string]any{
			"error":  err,
			"userID": userID,
		})
		return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error calculating three year subscription costs: %w", err)
	}

	// Get detailed db records for frontend (purchase type and payment method)
	var oneTimePurchases []models.SpendTrackingOneTimePurchaseDB
	if err := sta.db.SelectContext(
		ctx,
		&oneTimePurchases,
		GetCurrentMonthOneTimePurchasesQuery,
		userID,
		currentMonth,
	); err != nil {
			sta.logger.Error("Failed to get one-time purchases", map[string]any{
					"error":  err,
					"userID": userID,
			})
			return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error getting one-time purchases: %w", err)
	}

	var subscriptions []models.SpendTrackingLocationDB
	if err := sta.db.SelectContext(
			ctx,
			&subscriptions,
			GetActiveSubscriptionsQuery,
			userID,
	); err != nil {
			sta.logger.Error("Failed to get subscriptions", map[string]any{
					"error":  err,
					"userID": userID,
			})
			return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error getting subscriptions: %w", err)
	}

	currentTotalThisMonth, oneTimeThisMonth, recurringNextMonth := sta.buildTransactionArraysFromDatabaseRecords(
    oneTimePurchases,
    subscriptions,
    currentMonth,
)

	// Transform calculated data to BFF response
	monthlySpendingResponse := types.MonthlySpendingBFFResponseFINAL{
		CurrentMonthTotal:           currentMonthTotal,
		LastMonthTotal:              lastMonthTotal,
		PercentageChange:            percentageChange,
		ComparisonDateRange:         fmt.Sprintf("%s - %s", lastMonth.Format("Jan 2"), currentMonth.Format("Jan 2, 2006")),
		SpendingCategories:          sta.transformCalculatorCategoriesToBFFResponse(currentMonthAggregation.SpendingCategories),
	}

	// Transform calculated annual spending data
	annualSpendingResponse := annualSpendingForecast

	// Transform calculated yearly totals
	yearlyTotals := sta.transformThreeYearTotalsToBFFResponse(threeYearSubscriptionTotals)

	// Build FINAL BFF response
	response := types.SpendTrackingBFFResponseFINAL{
		TotalMonthlySpending:  monthlySpendingResponse,
		TotalAnnualSpending:   annualSpendingResponse,
		CurrentTotalThisMonth: currentTotalThisMonth,
		OneTimeThisMonth:      oneTimeThisMonth,
		RecurringNextMonth:    recurringNextMonth,
		YearlyTotals:          yearlyTotals,
	}

	sta.logger.Debug("GetSpendTrackingBFFResponse completed with calculated data", map[string]any{
    "userID":              userID,
    "currentMonthTotal":   currentMonthTotal,
    "previousMonthTotal":  lastMonthTotal,
    "percentageChange":    percentageChange,
    "response":            response,
	})

	return response, nil
}

// --- SINGLE GET OPERATION ---
func (sta *SpendTrackingDbAdapter) GetSingleSpendTrackingItem(
	ctx context.Context,
	userID string,
	itemID string,
) (models.SpendTrackingOneTimePurchaseDB, error) {
	sta.logger.Debug("GetSingleSpendTrackingItem called", map[string]any{
		"userID": userID,
		"itemID": itemID,
	})

	// Extract numeric ID from frontend format (e.g., "one-16" -> "16")
	if !strings.HasPrefix(itemID, "one-") {
		return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("invalid purchase ID format: must start with 'one-'")
	}
	numericID := strings.TrimPrefix(itemID, "one-")

	// Convert string ID to int for database query
	purchaseID, err := strconv.Atoi(numericID)  // ✅ NOW converts "16" to int 16
	if err != nil {
		return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("invalid purchase ID format: %w", err)
	}

	// Execute SELECT query
	var purchase models.SpendTrackingOneTimePurchaseDB
	err = sta.db.GetContext(
		ctx,
		&purchase,
		GetSingleSpendTrackingItemQuery,
		purchaseID,
		userID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("one-time purchase not found")
		}
		return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("failed to get one-time purchase: %w", err)
	}

	sta.logger.Debug("GetSingleSpendTrackingItem success", map[string]any{
		"purchase": purchase,
	})

	return purchase, nil
}



// --- WRITE OPERATIONS ---
func (sta *SpendTrackingDbAdapter) CreateOneTimePurchase(
	ctx context.Context,
	userID string,
	request models.SpendTrackingOneTimePurchaseDB,
) (models.SpendTrackingOneTimePurchaseDB, error) {
	sta.logger.Debug("CreateOneTimePurchase called", map[string]any{
		"userID": userID,
		"request": request,
	})

	// Start transaction
	tx, err := sta.db.BeginTxx(ctx, nil)
	if err != nil {
			return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Execute INSERT query
	var newPurchase models.SpendTrackingOneTimePurchaseDB
	err = tx.GetContext(
		ctx,
		&newPurchase,
		CreateOneTimePurchaseQuery,
		userID,
		request.Title,
		request.Amount,
		request.PurchaseDate,
		request.PaymentMethod,
		request.CategoryID,
		request.DigitalLocationID,
		request.IsDigital,
		request.IsWishlisted,
	)

	if err != nil {
		tx.Rollback()
		return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("failed to create one-time purchase: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	sta.logger.Debug("CreateOneTimePurchase success", map[string]any{
		"oneTimePurchase": newPurchase,
	})

	return newPurchase, nil
}

func (sta *SpendTrackingDbAdapter) UpdateOneTimePurchase(
	ctx context.Context,
	userID string,
	request models.SpendTrackingOneTimePurchaseDB,
) (models.SpendTrackingOneTimePurchaseDB, error) {
	sta.logger.Debug("UpdateOneTimePurchase called", map[string]any{
		"userID": userID,
		"request": request,
	})

	// Start transaction
	tx, err := sta.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Execute UPDATE query
	var updatedPurchase models.SpendTrackingOneTimePurchaseDB
	err = tx.GetContext(
		ctx,
		&updatedPurchase,
		UpdateOneTimePurchaseQuery,
		request.Title,
		request.Amount,
		request.PurchaseDate,
		request.PaymentMethod,
		request.CategoryID,
		request.DigitalLocationID,
		request.IsDigital,
		request.IsWishlisted,
		request.ID,
		userID,
	)

	if err != nil {
		tx.Rollback()
		return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("failed to update one-time purchase: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	sta.logger.Debug("UpdateOneTimePurchase success", map[string]any{
		"updatedPurchase": updatedPurchase,
	})

	return updatedPurchase, nil
}

func (sta *SpendTrackingDbAdapter) DeleteSpendTrackingItems(
	ctx context.Context,
	userID string,
	itemIDs []string,
) (int64, error) {
	sta.logger.Debug("DeleteSpendTrackingItems called", map[string]any{
		"userID": userID,
		"itemIDs": itemIDs,
	})

	// Validate input parameters
	if userID == "" {
		return 0, fmt.Errorf("user ID cannot be empty")
	}

	if len(itemIDs) == 0 {
		return 0, fmt.Errorf("no item IDs provided for deletion")
	}

	// Start transaction
	tx, err := sta.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Convert itemIDs from ["one-15"] to [15] before passing to SQL
	var numericIDs []int
	for _, itemID := range itemIDs {
		if !strings.HasPrefix(itemID, "one-") {
			return 0, fmt.Errorf("invalid item ID format: must start with 'one-'")
		}
		numericID := strings.TrimPrefix(itemID, "one-")
		id, err := strconv.Atoi(numericID)
		if err != nil {
			return 0, fmt.Errorf("invalid item ID format: %w", err)
		}
		numericIDs = append(numericIDs, id)
	}

	// Verify all items exist and belong to the user
	var count int
	err = tx.QueryRowxContext(
		ctx,
		CheckIfAllOneTimePurchasesExistForUserQuery,
		pq.Array(numericIDs),  // ← Now passing [15] to integer array
		userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error verifying one-time purchases: %w", err)
	}
	if count != len(itemIDs) {
		return 0, fmt.Errorf("one or more one-time purchases not found or do not belong to user")
	}

	// Delete all items in one go
	result, err := tx.ExecContext(
		ctx,
		DeleteOneTimePurchasesQuery,
		pq.Array(numericIDs),  // ← Use numericIDs instead of itemIDs
		userID,
	)
	if err != nil {
		return 0, fmt.Errorf("error executing delete: %w", err)
	}

	// Get the number of deleted rows
	totalDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %w", err)
	}

	// If not all items were deleted, return an error
	if totalDeleted < int64(len(itemIDs)) {
		return totalDeleted, fmt.Errorf("partial deletion: %d of %d items deleted", totalDeleted, len(itemIDs))
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	sta.logger.Debug("DeleteOneTimePurchases success", map[string]any{
		"totalDeleted": totalDeleted,
		"isBulk": len(itemIDs) > 1,
	})

	return totalDeleted, nil
}