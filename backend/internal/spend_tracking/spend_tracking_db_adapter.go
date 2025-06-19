package spend_tracking

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/postgres"
	"github.com/lokeam/qko-beta/internal/types"
)


type SpendTrackingDbAdapter struct {
	client  *postgres.PostgresClient
	db      *sqlx.DB
	logger  interfaces.Logger
}

func NewSpendTrackingDbAdapter(
	appContext *appcontext.AppContext,
) (*SpendTrackingDbAdapter, error) {
	// Log that we're creating the SpendTrackingDbAdapter
	appContext.Logger.Debug("Creating LibraryDbAdapter", map[string]any{"appContext": appContext})

	// Create sql db from px pool
	client, err := postgres.NewPostgresClient(appContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create sql connection for spend tracking db adapter: %w", err)
	}

	// Create sqlx db from px pool
	db, err := sqlx.Connect("pgx", appContext.Config.Postgres.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlx connection: %w", err)
	}

	// Register custom types? Not sure if needed
	db.MapperFunc(strings.ToLower)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &SpendTrackingDbAdapter{
		client:  client,
		db:      db,
		logger:  appContext.Logger,
	}, nil
}


// --- QUERIES ---
const (
	getCurrentAndPreviousMonthSpendingQuery = `
    SELECT * FROM monthly_spending_aggregates
    WHERE user_id = $1
    AND (year, month) IN (
        (EXTRACT(YEAR FROM CURRENT_DATE)::int, EXTRACT(MONTH FROM CURRENT_DATE)::int),
        (EXTRACT(YEAR FROM CURRENT_DATE - INTERVAL '1 month')::int,
         EXTRACT(MONTH FROM CURRENT_DATE - INTERVAL '1 month')::int)
    )
    ORDER BY year DESC, month DESC
	`

	// SELECT year, month, total_amount
	getMonthlySpendingAggregatesQuery = `
    SELECT id, user_id, year, month, total_amount, subscription_amount, one_time_amount, category_amounts, created_at, updated_at
    FROM monthly_spending_aggregates
    WHERE user_id = $1
    ORDER BY year, month
	`

	getYearlySpendingQuery = `
		SELECT * FROM yearly_spending_aggregates
		WHERE user_id = $1
		AND year >= EXTRACT(YEAR FROM CURRENT_DATE)::int - 2
		ORDER BY year DESC
	`

	getCurrentMonthOneTimePurchasesQuery = `
		SELECT otp.*, sc.media_type as media_type
		FROM one_time_purchases otp
		LEFT JOIN spending_categories sc ON otp.spending_category_id = sc.id
		WHERE otp.user_id = $1
		AND EXTRACT(YEAR FROM purchase_date) = EXTRACT(YEAR FROM CURRENT_DATE)
		AND EXTRACT(MONTH FROM purchase_date) = EXTRACT(MONTH FROM CURRENT_DATE)
		ORDER BY purchase_date DESC
	`

	getActiveSubscriptionsQuery = `
		SELECT
				dl.id,
				dl.user_id,
				dl.name,
				dl.is_subscription,
				dl.is_active,
				dl.payment_method,
				dl.created_at,
				dl.updated_at,
				dls.billing_cycle,
				dls.cost_per_cycle,
				dls.anchor_date,
				dls.last_payment_date,
				dls.next_payment_date,
				dls.payment_method as subscription_payment_method
		FROM digital_locations dl
		LEFT JOIN digital_location_subscriptions dls ON dl.id = dls.digital_location_id
		WHERE dl.user_id = $1
		AND dl.is_subscription = true
		AND dl.is_active = true
	ORDER BY dls.next_payment_date ASC
	`


)

// --- TRANSFORMATION LOGIC ---
func (sta *SpendTrackingDbAdapter) transformMonthlySpendingDBToResponse(
	currentMonth models.SpendTrackingMonthlyAggregateDB,
	previousMonth models.SpendTrackingMonthlyAggregateDB,
) types.MonthlySpendingBFFResponseFINAL {
	// Log function call
	sta.logger.Debug("transformMonthlySpendingDBToResponse called", map[string]any{
		"currentMonth": currentMonth,
		"previousMonth": previousMonth,
	})

	// Calc percentage change
	percentageChange := 0.0
	if previousMonth.TotalAmount > 0 {
		percentageChange = ((currentMonth.TotalAmount - previousMonth.TotalAmount) / previousMonth.TotalAmount) * 100
	}

	// Parse category amounts from JSONB
	var categoryAmounts map[string]float64
	sta.logger.Debug("Attempting to unmarshal category_amounts", map[string]any{
		"categoryAmounts": string(currentMonth.CategoryAmounts),
	})

	if err := json.Unmarshal(currentMonth.CategoryAmounts, &categoryAmounts); err != nil {
		sta.logger.Error("Failed to unmarshal category_amounts", map[string]any{
			"error": err,
			"categoryAmounts": string(currentMonth.CategoryAmounts),
		})
		categoryAmounts = make(map[string]float64)
	} else {
		sta.logger.Debug("Successfully unmarshaled category_amounts", map[string]any{
			"categoryAmounts": categoryAmounts,
		})
	}

	// Transform category amounts to response format
	spendingCategories := make([]types.SpendingCategoryBFFResponseFINAL, 0, len(categoryAmounts))
	for name, value := range categoryAmounts {
		sta.logger.Debug("Processing category", map[string]any{
			"name": name,
			"value": value,
		})
		spendingCategories = append(spendingCategories, types.SpendingCategoryBFFResponseFINAL{
			Name:  name,
			Value: value,
		})
	}

	// Format the date range in response
	currentDate := time.Date(currentMonth.Year, time.Month(currentMonth.Month), 1, 0, 0, 0, 0, time.UTC)
	previousDate := time.Date(previousMonth.Year, time.Month(previousMonth.Month), 1, 0, 0, 0, 0, time.UTC)
	dateRange := fmt.Sprintf("%s - %s",
		previousDate.Format("Jan 2"),
		currentDate.Format("Jan 2, 2006"),
	)

	response := types.MonthlySpendingBFFResponseFINAL{
		CurrentMonthTotal:    currentMonth.TotalAmount,
		LastMonthTotal:      previousMonth.TotalAmount,
		PercentageChange:    percentageChange,
		ComparisonDateRange: dateRange,
		SpendingCategories:  spendingCategories,
	}

	sta.logger.Debug("Returning monthly spending response", map[string]any{
		"response": response,
	})

	return response
}

func tranformYearlySpendingDBToResponse(
	yearlySpendingAggregates []models.SpendTrackingYearlyAggregateDB,
	monthlySpendingAggregates []models.SpendTrackingMonthlyAggregateDB,
) types.AnnualSpendingBFFResponseFINAL {
	// Initialize monthly expenditures array
	monthlyExpenditures := make([]types.MonthlyExpenditureBFFResponseFINAL, 12)
	for i := range monthlyExpenditures {
		monthlyExpenditures[i] = types.MonthlyExpenditureBFFResponseFINAL{
			Month:       time.Month(i + 1).String()[:3],
			Expenditure: 0,
		}
	}

	// Get current year and month for proper data filtering
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	// Populate monthly expenditures with actual data
	for _, agg := range monthlySpendingAggregates {
		// Only process data for the current year
		if agg.Year == currentYear && int(agg.Month) <= currentMonth {
				monthIndex := int(agg.Month) - 1 // Convert to 0-based index
				if monthIndex >= 0 && monthIndex < 12 {
						// Only update if we have data for this month
						if agg.TotalAmount > 0 {
								monthlyExpenditures[monthIndex].Expenditure = agg.TotalAmount
						}
				}
		}
}

	// Calculate median monthly cost from yearly aggregates
	var monthlyCosts []float64
	for _, agg := range yearlySpendingAggregates {
		if agg.Year == currentYear {
			monthlyCosts = append(monthlyCosts, agg.TotalAmount/12)
		}
	}
	sort.Float64s(monthlyCosts)
	medianMonthlyCost := 0.0
	if len(monthlyCosts) > 0 {
		if len(monthlyCosts)%2 == 0 {
			medianMonthlyCost = (monthlyCosts[len(monthlyCosts)/2-1] + monthlyCosts[len(monthlyCosts)/2]) / 2
		} else {
			medianMonthlyCost = monthlyCosts[len(monthlyCosts)/2]
		}
	}

	// Format date range
	dateRange := fmt.Sprintf("January %d - January %d", currentYear, currentYear+1)

	return types.AnnualSpendingBFFResponseFINAL{
		DateRange:          dateRange,
		MonthlyExpenditures: monthlyExpenditures,
		MedianMonthlyCost:  medianMonthlyCost,
	}
}

func transformOneTimePurchasesDBToResponse(
	oneTimePurchase models.SpendTrackingOneTimePurchaseDB,
) types.SpendingItemBFFResponseFINAL{
	return types.SpendingItemBFFResponseFINAL{
		ID:                fmt.Sprintf("one-%d", oneTimePurchase.ID),
        Title:            oneTimePurchase.Title,
        Amount:           oneTimePurchase.Amount,
        SpendTransactionType: "ONE_TIME",
        PaymentMethod:    oneTimePurchase.PaymentMethod,
        MediaType:        oneTimePurchase.MediaType,
        CreatedAt:        oneTimePurchase.CreatedAt.Unix(),
        UpdatedAt:        oneTimePurchase.UpdatedAt.Unix(),
        IsActive:         true,
        IsDigital:        oneTimePurchase.IsDigital,
        IsWishlisted:     oneTimePurchase.IsWishlisted,
        PurchaseDate:     oneTimePurchase.PurchaseDate.Unix(),
	}
}

func transformSubscriptionsDBToResponse(
	digitalLocation models.SpendTrackingLocationDB,
	subscriptionDetails models.SpendTrackingSubscriptionDB,
) types.SpendingItemBFFResponseFINAL{
	// Calc yearly spending
	annualAmount:= subscriptionDetails.CostPerCycle
	switch subscriptionDetails.BillingCycle {
	case "1 month":
		annualAmount = subscriptionDetails.CostPerCycle * 12
	case "3 month":
		annualAmount = subscriptionDetails.CostPerCycle * 4
	case "6 month":
		annualAmount = subscriptionDetails.CostPerCycle * 2
	case "12 month":
		annualAmount = subscriptionDetails.CostPerCycle
	}

	yearlySpending := []types.SingleYearlyTotalBFFResponseFINAL{
		{Year: time.Now().Year() - 2, Amount: annualAmount},
		{Year: time.Now().Year() - 1, Amount: annualAmount},
		{Year: time.Now().Year(), Amount: annualAmount},
	}

	return types.SpendingItemBFFResponseFINAL{
		ID:                     fmt.Sprintf("sub-%d", subscriptionDetails.ID),
		Title:                  digitalLocation.Name,
		Amount:                 subscriptionDetails.CostPerCycle,
		SpendTransactionType:   "SUBSCRIPTION",
		PaymentMethod:          subscriptionDetails.PaymentMethod,
		MediaType:              "SUBSCRIPTION",
		Provider:               strings.ToLower(digitalLocation.Name),
		CreatedAt:              digitalLocation.CreatedAt.Unix(),
		UpdatedAt:              digitalLocation.UpdatedAt.Unix(),
		IsActive:               digitalLocation.IsActive,
		BillingCycle:           subscriptionDetails.BillingCycle,
		NextBillingDate:        subscriptionDetails.NextPaymentDate.Unix(),
		YearlySpending:         yearlySpending,
}
}



// GET
func (sta *SpendTrackingDbAdapter) GetSpendTrackingBFFResponse(
	ctx context.Context,
	userID string,
) (types.SpendTrackingBFFResponseFINAL, error) {
	// Log function call
	sta.logger.Debug("GetSpendTrackingBFFResponse called", map[string]any{
		"userID": userID,
	})

	// Start transaction
	tx, err := sta.db.BeginTxx(ctx, nil)
	if err != nil {
		return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("failed to start transaction for GetSpendTrackingBFFResponse: %w", err)
	}
	defer tx.Rollback()

	// Get monthly spending data
	var monthlySpendingAggregates []models.SpendTrackingMonthlyAggregateDB
	if err := tx.SelectContext(
		ctx,
		&monthlySpendingAggregates,
		getMonthlySpendingAggregatesQuery,
		userID,
	); err != nil {
		return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error querying monthly spending: %w", err)
	}

	// Log the raw data from the database
	sta.logger.Debug("Raw monthly spending data from database", map[string]any{
		"count": len(monthlySpendingAggregates),
		"data": monthlySpendingAggregates,
	})

	// Get yearly spending data
	var yearlySpendingAggregates []models.SpendTrackingYearlyAggregateDB
	if err := tx.SelectContext(
		ctx, &yearlySpendingAggregates,
		getYearlySpendingQuery,
		userID,
	); err != nil {
		return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error querying yearly spending: %w", err)
	}

	// Get one-time purchases
	var oneTimePurchases []models.SpendTrackingOneTimePurchaseDB
	if err := tx.SelectContext(
		ctx,
		&oneTimePurchases,
		getCurrentMonthOneTimePurchasesQuery,
		userID,
	); err != nil {
		return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error querying one-time purchases: %w", err)
	}

	// Get subscriptions
	var subscriptions []models.SpendTrackingLocationDB
	if err := tx.SelectContext(
		ctx,
		&subscriptions,
		getActiveSubscriptionsQuery,
		userID,
	); err != nil {
		return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("error querying subscriptions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return types.SpendTrackingBFFResponseFINAL{}, fmt.Errorf("failed to commit transaction for GetSpendTrackingBFFResponse: %w", err)
	}

	// Transform monthly and annual spending data
	var monthlySpendingResponse types.MonthlySpendingBFFResponseFINAL
	if len(monthlySpendingAggregates) >= 2 {
		sta.logger.Debug("Transforming monthly spending data", map[string]any{
			"currentMonth": monthlySpendingAggregates[0],
			"previousMonth": monthlySpendingAggregates[1],
		})
		monthlySpendingResponse = sta.transformMonthlySpendingDBToResponse(monthlySpendingAggregates[1], monthlySpendingAggregates[0])
	} else {
		sta.logger.Debug("Insufficient monthly spending data, using empty state", map[string]any{
			"count": len(monthlySpendingAggregates),
		})
		// Empty state matching mock data structure
		monthlySpendingResponse = types.MonthlySpendingBFFResponseFINAL{
			CurrentMonthTotal:    0,
			LastMonthTotal:      0,
			PercentageChange:    0,
			ComparisonDateRange: "No data available",
			SpendingCategories:  []types.SpendingCategoryBFFResponseFINAL{},
		}
	}

	// Log the transformed response
	sta.logger.Debug("Transformed monthly spending response", map[string]any{
		"response": monthlySpendingResponse,
	})

	annualSpendingResponse := tranformYearlySpendingDBToResponse(yearlySpendingAggregates, monthlySpendingAggregates)

	// Transform one-time purchases
	oneTimePurchasesResponse := make([]types.SpendingItemBFFResponseFINAL, len(oneTimePurchases))
	for i := 0; i < len(oneTimePurchases); i++ {
		oneTimePurchasesResponse[i] = transformOneTimePurchasesDBToResponse(oneTimePurchases[i])
	}

	// Transform subscriptions
	subscriptionsResponse := make([]types.SpendingItemBFFResponseFINAL, len(subscriptions))
	for i := 0; i < len(subscriptions); i++ {
		subscriptionsResponse[i] = transformSubscriptionsDBToResponse(subscriptions[i], models.SpendTrackingSubscriptionDB{
			ID:              0, // Not needed for response
			LocationID:      subscriptions[i].ID,
			BillingCycle:    subscriptions[i].BillingCycle,
			CostPerCycle:    subscriptions[i].CostPerCycle,
			AnchorDate:      subscriptions[i].AnchorDate,
			LastPaymentDate: subscriptions[i].LastPaymentDate,
			NextPaymentDate: subscriptions[i].NextPaymentDate,
			PaymentMethod:   subscriptions[i].SubscriptionPaymentMethod,
			CreatedAt:       subscriptions[i].CreatedAt,
			UpdatedAt:       subscriptions[i].UpdatedAt,
		})
	}

	// Build yearly totals
	yearlyTotals := types.AllYearlyTotalsBFFResponseFINAL{
		SubscriptionTotal: make([]types.SingleYearlyTotalBFFResponseFINAL, 0, len(yearlySpendingAggregates)),
		OneTimeTotal:      make([]types.SingleYearlyTotalBFFResponseFINAL, 0, len(yearlySpendingAggregates)),
		CombinedTotal:     make([]types.SingleYearlyTotalBFFResponseFINAL, 0, len(yearlySpendingAggregates)),
	}

	// Only process yearly totals if we have data
	if len(yearlySpendingAggregates) > 0 {
		// Sort yearly aggregates by year
		sort.Slice(yearlySpendingAggregates, func(i, j int) bool {
			return yearlySpendingAggregates[i].Year < yearlySpendingAggregates[j].Year
		})

		// Construct yearly totals
		for _, agg := range yearlySpendingAggregates {
			yearlyTotals.SubscriptionTotal = append(yearlyTotals.SubscriptionTotal, types.SingleYearlyTotalBFFResponseFINAL{
				Year:   agg.Year,
				Amount: agg.SubscriptionAmount,
			})
			yearlyTotals.OneTimeTotal = append(yearlyTotals.OneTimeTotal, types.SingleYearlyTotalBFFResponseFINAL{
				Year:   agg.Year,
				Amount: agg.OneTimeAmount,
			})
			yearlyTotals.CombinedTotal = append(yearlyTotals.CombinedTotal, types.SingleYearlyTotalBFFResponseFINAL{
				Year:   agg.Year,
				Amount: agg.TotalAmount,
			})
		}
	}

	// Return
	response := types.SpendTrackingBFFResponseFINAL{
		TotalMonthlySpending: monthlySpendingResponse,
		TotalAnnualSpending:  annualSpendingResponse,
		CurrentTotalThisMonth: append(oneTimePurchasesResponse, subscriptionsResponse...),
		OneTimeThisMonth:     oneTimePurchasesResponse,
		RecurringNextMonth:   subscriptionsResponse,
		YearlyTotals:         yearlyTotals,
	}

	sta.logger.Debug("GetSpendTrackingBFFResponse response", map[string]any{
		"response": response,
	})

	return response, nil
}