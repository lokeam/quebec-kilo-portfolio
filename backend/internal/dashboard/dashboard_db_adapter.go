package dashboard

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/postgres"
	"github.com/lokeam/qko-beta/internal/spend_tracking"
	"github.com/lokeam/qko-beta/internal/types"
)

// Adapter struct for Dashboard BFF
type DashboardDbAdapter struct {
  client                     *postgres.PostgresClient
  db                         *sqlx.DB
  logger                     interfaces.Logger
  spendTrackingCalculator    *spend_tracking.SpendTrackingCalculator
}

// Constructor for DashboardDbAdapter
func NewDashboardDbAdapter(appContext *appcontext.AppContext) (*DashboardDbAdapter, error) {
  client, err := postgres.NewPostgresClient(appContext)
  if err != nil {
      return nil, fmt.Errorf("failed to create sql connection for dashboard db adapter: %w", err)
  }
  db, err := sqlx.Connect("pgx", appContext.Config.Postgres.ConnectionString)
  if err != nil {
      return nil, fmt.Errorf("failed to create sqlx connection: %w", err)
  }
  db.MapperFunc(strings.ToLower)
  db.SetMaxOpenConns(25)
  db.SetMaxIdleConns(25)
  db.SetConnMaxLifetime(5 * time.Minute)

  spendTrackingDBAdapter, err := spend_tracking.NewSpendTrackingDbAdapter(appContext)
  if err != nil {
    return nil, fmt.Errorf("failed to create spend tracking db adapter: %w", err)
  }

  spendTrackingCalculator, err := spend_tracking.NewSpendTrackingCalculator(appContext, spendTrackingDBAdapter)
  if err != nil {
    return nil, fmt.Errorf("failed to create spend tracking calculator: %w", err)
  }

  return &DashboardDbAdapter{
      client: client,
      db:     db,
      logger: appContext.Logger,
      spendTrackingCalculator: spendTrackingCalculator,
  }, nil
}

// --- SQL QUERY CONSTANTS ---

const (
  // Get total games and last updated
  getGameStatsQuery = `
      SELECT 'Games' AS title, 'games' AS icon, COUNT(*) AS value, MAX(created_at) AS last_updated
      FROM user_games
      WHERE user_id = $1
  `

  // Get total monthly online services costs and last updated
  getSubscriptionStatsQuery = `
      SELECT 'Subscription Costs' AS title, 'coin' AS icon,
           COALESCE(ROUND(SUM(
               CASE
                   WHEN dls.billing_cycle = '1 month' THEN dls.cost_per_cycle
                   WHEN dls.billing_cycle = '3 month' THEN dls.cost_per_cycle / 3
                   WHEN dls.billing_cycle = '6 month' THEN dls.cost_per_cycle / 6
                   WHEN dls.billing_cycle = '12 month' THEN dls.cost_per_cycle / 12
                   ELSE dls.cost_per_cycle
               END
           ), 2), 0) AS value, MAX(dls.updated_at) AS last_updated
      FROM digital_location_subscriptions dls
      JOIN digital_locations dl ON dls.digital_location_id = dl.id
      WHERE dl.user_id = $1 AND dl.is_subscription = true
  `

  // Get digital storage locations count and last updated
  getDigitalLocationStatsQuery = `
      SELECT 'Digital Storage' AS title, 'onlineServices' AS icon,
             COUNT(*) AS value, MAX(updated_at) AS last_updated
      FROM digital_locations
      WHERE user_id = $1
  `

  // Get physical storage locations count and last updated
  getPhysicalLocationStatsQuery = `
      SELECT
        'Physical Storage' AS title,
        'package' AS icon,
        COUNT(DISTINCT pl.id) AS value,
        COUNT(s.id) AS secondary_value,
        MAX(pl.updated_at) AS last_updated
    FROM physical_locations pl
    LEFT JOIN sublocations s ON pl.id = s.physical_location_id
    WHERE pl.user_id = $1
  `

  // Get all digital locations with details
  getDigitalLocationsQuery = `
  SELECT
    dl.name,
    dl.url,
    COALESCE(dls.billing_cycle, '') AS billing_cycle,
    COALESCE(dls.cost_per_cycle, 0) AS monthly_fee,
    COALESCE(stored_games.count, 0) AS stored_items,
    dl.is_subscription,
    dls.next_payment_date
  FROM digital_locations dl
  LEFT JOIN digital_location_subscriptions dls ON dl.id = dls.digital_location_id
  LEFT JOIN (
    SELECT
      dgl.digital_location_id,
      COUNT(*) as count
    FROM digital_game_locations dgl
    JOIN user_games ug ON dgl.user_game_id = ug.id
    WHERE ug.user_id = $1
    GROUP BY dgl.digital_location_id
  ) stored_games ON dl.id = stored_games.digital_location_id
  WHERE dl.user_id = $1
`

  // Get all sublocations with parent location details
  getSublocationsQuery = `
    SELECT s.id AS sublocation_id, s.name AS sublocation_name, s.location_type AS sublocation_type,
           COALESCE(stored_games.count, 0) AS stored_items, pl.id AS parent_location_id, pl.name AS parent_location_name,
           pl.location_type AS parent_location_type, pl.bg_color AS parent_location_bg_color
    FROM sublocations s
    JOIN physical_locations pl ON s.physical_location_id = pl.id
    LEFT JOIN (
      SELECT
        pgl.sublocation_id,
        COUNT(*) as count
      FROM physical_game_locations pgl
      JOIN user_games ug ON pgl.user_game_id = ug.id
      WHERE ug.user_id = $1
      GROUP BY pgl.sublocation_id
    ) stored_games ON s.id = stored_games.sublocation_id
    WHERE s.user_id = $1
  `

  // Get platform distribution
  getPlatformListQuery = `
      SELECT p.name AS platform, COUNT(*) AS item_count
      FROM user_games ug
      JOIN platforms p ON ug.platform_id = p.id
      WHERE ug.user_id = $1
      GROUP BY p.name
  `

  // Get new items added this month
  getNewItemsThisMonthQuery = `
      SELECT COUNT(*) AS new_items
      FROM user_games
      WHERE user_id = $1 AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE)
  `

  // Get monthly expenditures (for the last 12 months)
  getMonthlyExpendituresQuery = `
      SELECT
        TO_CHAR(TO_DATE(CONCAT(year, '-', LPAD(month::text, 2, '0'), '-01'), 'YYYY-MM-DD'), 'YYYY-MM-01') AS date,
        one_time_amount AS one_time_purchase,
        COALESCE((category_amounts->>'hardware')::DECIMAL(10,2), 0) AS hardware,
        COALESCE((category_amounts->>'dlc')::DECIMAL(10,2), 0) AS dlc,
        COALESCE((category_amounts->>'in_game_purchase')::DECIMAL(10,2), 0) AS in_game_purchase
      FROM monthly_spending_aggregates
      WHERE user_id = $1
      ORDER BY year, month
      LIMIT 12
  `
)

// --- TRANSFORMATION LOGIC ---

func (dda *DashboardDbAdapter) transformGameStatsDBToResponse(
  db models.DashboardGameStatsDB,
) types.DashboardStatBFFResponse {
  var lastUpdatedUnix int64
  if db.LastUpdated != nil {
    lastUpdatedUnix = db.LastUpdated.Unix()
  } else {
    lastUpdatedUnix = 0
  }

  return types.DashboardStatBFFResponse{
    Title: db.Title,
    Icon: db.Icon,
    Value: db.Value,
    SecondaryValue: db.SecondaryValue,
    LastUpdated: lastUpdatedUnix,
  }
}

func (dda *DashboardDbAdapter) transformDigitalLocationDBToResponse(
  db models.DashboardDigitalLocationDB,
) types.DashboardDigitalLocationBFFResponse {
  renewsNextMonth := false
  if db.IsSubscription {
    if db.BillingCycle == "1 month" {
      // All monthly subscriptions always renew next month
      renewsNextMonth = true
    } else if db.NextPaymentDate != nil {
      now := time.Now()
      nextMonth := now.AddDate(0, 1, 0)
      if db.NextPaymentDate.Year() == nextMonth.Year() && db.NextPaymentDate.Month() == nextMonth.Month() {
        renewsNextMonth = true
      }
    }
  }

  return types.DashboardDigitalLocationBFFResponse{
    Logo: db.Name,
    Name: db.Name,
    Url:  db.Url,
    BillingCycle: db.BillingCycle,
    MonthlyFee: db.MonthlyFee,
    RenewsNextMonth: renewsNextMonth,
    StoredItems: db.StoredItems,
  }
}

func (dda *DashboardDbAdapter) transformSublocationDBToResponse(
  db models.DashboardSublocationDB,
) types.DashboardSublocationBFFResponse {
  return types.DashboardSublocationBFFResponse{
    SublocationId:         db.SublocationId,
    SublocationName:       db.SublocationName,
    SublocationType:       db.SublocationType,
    StoredItems:           db.StoredItems,
    ParentLocationId:      db.ParentLocationId,
    ParentLocationName:    db.ParentLocationName,
    ParentLocationType:    db.ParentLocationType,
    ParentLocationBgColor: db.ParentLocationBgColor,
  }
}

func (dda *DashboardDbAdapter) transformPlatformDBToResponse(
  db models.DashboardPlatformDB,
) types.DashboardPlatformBFFResponse {
  return types.DashboardPlatformBFFResponse{
      Platform:  db.Platform,
      ItemCount: db.ItemCount,
  }
}

func (dda *DashboardDbAdapter) transformMonthlyExpenditureDBToResponse(
  db models.DashboardMonthlyExpenditureDB,
  subscriptionCost float64,
) types.DashboardMonthlyExpenditureBFFResponse {
  return types.DashboardMonthlyExpenditureBFFResponse{
      Date:            db.Date,
      OneTimePurchase: db.OneTimePurchase,
      Hardware:        db.Hardware,
      Dlc:             db.Dlc,
      InGamePurchase:  db.InGamePurchase,
      Subscription:    subscriptionCost,
  }
}


// --- Helper Methods ---
func (dda *DashboardDbAdapter) calculateSubscriptionCostsForMonths(
  userID string,
  months []time.Time,
) (map[string]float64, error) {
  dda.logger.Debug("calculateSubscriptionCostsForMonths called", map[string]any{
    "userID": userID,
    "monthCount": len(months),
    "months": months,
  })

  subscriptionCosts := make(map[string]float64)

  for _, month := range months {
    // Calculate subscription costs for this month using existing calculator
    // Ensure month is the first day of the month
    firstDayOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
    monthlySubscriptionCost, err := dda.spendTrackingCalculator.CalculateMonthlySubscriptionCosts(userID, firstDayOfMonth)
    if err != nil {
        dda.logger.Error("Failed to calculate subscription costs for month", map[string]any{
            "error": err,
            "userID": userID,
            "month": month,
        })
        // Continue with other months even if one fails
        monthlySubscriptionCost = 0.0
    }

    // Format month key to match the date format from the query
    monthKey := month.Format("2006-01-01")
    subscriptionCosts[monthKey] = monthlySubscriptionCost

    dda.logger.Debug("Calculated subscription costs for month", map[string]any{
        "userID": userID,
        "month": monthKey,
        "subscriptionCost": monthlySubscriptionCost,
    })
  }

  dda.logger.Debug("calculateSubscriptionCostsForMonths completed", map[string]any{
      "userID": userID,
      "subscriptionCosts": subscriptionCosts,
  })

  return subscriptionCosts, nil
}


// Calculate monthly expenditures dynamically from actual purchase data
func (dda *DashboardDbAdapter) calculateMonthlyExpendituresDynamically(
  userID string,
) []models.DashboardMonthlyExpenditureDB {
  dda.logger.Debug("calculateMonthlyExpendituresDynamically called", map[string]any{
    "userID": userID,
  })

  // Get the last 12 months
  currentMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
  monthlyExpenditures := make([]models.DashboardMonthlyExpenditureDB, 0, 12)

  for i := 11; i >= 0; i-- {
    targetMonth := currentMonth.AddDate(0, -i, 0)

    // Calculate monthly spending for this month using the spend tracking calculator
    monthlySpending, err := dda.spendTrackingCalculator.CalculateMonthlyMinimumSpending(userID, targetMonth)
    if err != nil {
      dda.logger.Error("Failed to calculate monthly spending", map[string]any{
        "error": err,
        "userID": userID,
        "targetMonth": targetMonth,
      })
      monthlySpending = 0.0
    }

    // Get one-time purchases for this month to calculate category breakdown
    var oneTimePurchases []models.SpendTrackingOneTimePurchaseDB
    err = dda.db.SelectContext(
      context.Background(),
      &oneTimePurchases,
      `SELECT otp.*, sc.media_type as media_type
      FROM one_time_purchases otp
      LEFT JOIN spending_categories sc ON otp.spending_category_id = sc.id
      WHERE otp.user_id = $1
      AND EXTRACT(YEAR FROM purchase_date) = EXTRACT(YEAR FROM $2::timestamp)
      AND EXTRACT(MONTH FROM purchase_date) = EXTRACT(MONTH FROM $2::timestamp)
      ORDER BY purchase_date DESC`,
      userID,
      targetMonth,
    )
    if err != nil {
      dda.logger.Error("Failed to get one-time purchases for month", map[string]any{
        "error": err,
        "userID": userID,
        "targetMonth": targetMonth,
      })
    }

    // Calculate category breakdown
    hardware := 0.0
    dlc := 0.0
    inGamePurchase := 0.0
    oneTimePurchase := 0.0

    for _, purchase := range oneTimePurchases {
      if purchase.PurchaseDate.Year() == targetMonth.Year() &&
         purchase.PurchaseDate.Month() == targetMonth.Month() {
        oneTimePurchase += purchase.Amount

        // Categorize by media type
        switch purchase.MediaType {
        case "hardware":
          hardware += purchase.Amount
        case "dlc":
          dlc += purchase.Amount
        case "in_game_purchase":
          inGamePurchase += purchase.Amount
        }
      }
    }

    // Format date for response
    dateStr := targetMonth.Format("2006-01-02")

    monthlyExpenditure := models.DashboardMonthlyExpenditureDB{
      Date:            dateStr,
      OneTimePurchase: oneTimePurchase,
      Hardware:        hardware,
      Dlc:             dlc,
      InGamePurchase:  inGamePurchase,
    }

    monthlyExpenditures = append(monthlyExpenditures, monthlyExpenditure)

    dda.logger.Debug("Calculated monthly expenditure", map[string]any{
      "userID": userID,
      "targetMonth": targetMonth,
      "monthlySpending": monthlySpending,
      "oneTimePurchase": oneTimePurchase,
      "hardware": hardware,
      "dlc": dlc,
      "inGamePurchase": inGamePurchase,
    })
  }

  dda.logger.Debug("calculateMonthlyExpendituresDynamically completed", map[string]any{
    "userID": userID,
    "monthlyExpendituresCount": len(monthlyExpenditures),
  })

  return monthlyExpenditures
}


func (dda *DashboardDbAdapter) GetDashboardBFFResponse(
  ctx context.Context,
  userID string,
) (types.DashboardBFFResponse, error) {
  dda.logger.Debug("GetDashboardBFFResponse called", map[string]any{"userID": userID})

  tx, err := dda.db.BeginTxx(ctx, nil)
  if err != nil {
    return types.DashboardBFFResponse{}, fmt.Errorf("failed to start transaction: %w", err)
  }
  defer tx.Rollback()

  // 1. Stats Queries
  var gameStatsDB, subscriptionStatsDB, digitalLocationStatsDB, physicalLocationStatsDB models.DashboardGameStatsDB

  if err := tx.GetContext(ctx, &gameStatsDB, getGameStatsQuery, userID); err != nil {
    dda.logger.Error("Error fetching game stats", map[string]any{
      "error": err,
      "userID": userID,
    })
    return types.DashboardBFFResponse{}, fmt.Errorf("error fetching game stats: %w", err)
  }
  if err := tx.GetContext(ctx, &subscriptionStatsDB, getSubscriptionStatsQuery, userID); err != nil {
    dda.logger.Error("Error fetching subscription stats", map[string]any{
      "error": err,
      "userID": userID,
    })
    return types.DashboardBFFResponse{}, fmt.Errorf("error fetching subscription stats: %w", err)
  }
  if err := tx.GetContext(ctx, &digitalLocationStatsDB, getDigitalLocationStatsQuery, userID); err != nil {
    dda.logger.Error("Error fetching digital location stats", map[string]any{
      "error": err,
      "userID": userID,
    })
    return types.DashboardBFFResponse{}, fmt.Errorf("error fetching digital location stats: %w", err)
  }
  if err := tx.GetContext(ctx, &physicalLocationStatsDB, getPhysicalLocationStatsQuery, userID); err != nil {
    dda.logger.Error("Error fetching physical location stats", map[string]any{
      "error": err,
      "userID": userID,
    })
    return types.DashboardBFFResponse{}, fmt.Errorf("error fetching physical location stats: %w", err)
  }

  // 2. Digital Locations
  var digitalLocationsDB []models.DashboardDigitalLocationDB
  if err := tx.SelectContext(ctx, &digitalLocationsDB, getDigitalLocationsQuery, userID); err != nil {
    dda.logger.Error("Error fetching digital locations", map[string]any{
      "error": err,
      "userID": userID,
    })
    return types.DashboardBFFResponse{}, fmt.Errorf("error fetching digital locations: %w", err)
  }

  // Calculate the annualized total for all active subscriptions
  annualizedSubscriptionTotal := 0.0
  for _, loc := range digitalLocationsDB {
    if loc.MonthlyFee > 0 && loc.BillingCycle != "" {
      var multiplier int
      switch loc.BillingCycle {
      case "1 month":
        multiplier = 12
      case "3 month":
        multiplier = 4
      case "6 month":
        multiplier = 2
      case "12 month":
        multiplier = 1
      default:
        multiplier = 0
      }
      annualizedSubscriptionTotal += loc.MonthlyFee * float64(multiplier)
    }
  }
  subscriptionTotal := annualizedSubscriptionTotal

  // 3. Sublocations
  var sublocationsDB []models.DashboardSublocationDB
  if err := tx.SelectContext(ctx, &sublocationsDB, getSublocationsQuery, userID); err != nil {
    dda.logger.Error("Error fetching sublocations", map[string]any{
      "error": err,
      "userID": userID,
    })
    return types.DashboardBFFResponse{}, fmt.Errorf("error fetching sublocations: %w", err)
  }

  // 4. Platform List
  var platformListDB []models.DashboardPlatformDB
  if err := tx.SelectContext(ctx, &platformListDB, getPlatformListQuery, userID); err != nil {
    dda.logger.Error("Error fetching platform list", map[string]any{
      "error": err,
      "userID": userID,
    })
    return types.DashboardBFFResponse{}, fmt.Errorf("error fetching platform list: %w", err)
  }

  // 5. New Items This Month
  var newItemsThisMonth int
  if err := tx.GetContext(ctx, &newItemsThisMonth, getNewItemsThisMonthQuery, userID); err != nil {
    dda.logger.Error("Error fetching new items this month", map[string]any{
      "error": err,
      "userID": userID,
    })
    return types.DashboardBFFResponse{}, fmt.Errorf("error fetching new items this month: %w", err)
  }

  // 6. Monthly Expenditures - Calculate dynamically instead of reading from aggregates table
  monthlyExpendituresDB := dda.calculateMonthlyExpendituresDynamically(userID)

  // Calculate the current month's total subscription cost using the business logic
  currentMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
  currentMonthSubscriptionCost, err := dda.spendTrackingCalculator.CalculateMonthlySubscriptionCosts(userID, currentMonth)
  if err != nil {
    dda.logger.Error("Failed to calculate current month subscription cost", map[string]any{
      "error": err,
      "userID": userID,
      "currentMonth": currentMonth,
    })
    currentMonthSubscriptionCost = 0.0
  }

  // Extract unique months from the monthlyExpendituresDB instead of hardcoding 2024
  // Calculate subscription costs for ALL months in the response
  uniqueMonths := make([]time.Time, 0, len(monthlyExpendituresDB))
  for _, expenditure := range monthlyExpendituresDB {
      expenditureMonth, err := time.Parse("2006-01-02", expenditure.Date)
      if err != nil {
          dda.logger.Error("Failed to parse expenditure date", map[string]any{
              "error": err,
              "date": expenditure.Date,
          })
          continue
      }
      // Normalize to first day of month
      firstDayOfMonth := time.Date(expenditureMonth.Year(), expenditureMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
      uniqueMonths = append(uniqueMonths, firstDayOfMonth)

      dda.logger.Debug("Added month for subscription calculation", map[string]any{
          "originalDate": expenditure.Date,
          "normalizedDate": firstDayOfMonth.Format("2006-01-01"),
      })
  }

  calculatedSubscriptionCosts, err := dda.calculateSubscriptionCostsForMonths(
    userID,
    uniqueMonths,
  )
  if err != nil {
    dda.logger.Error("Failed to calculate subscription costs per month", map[string]any{
      "error": err,
      "userID": userID,
    })
    // Continue with zero subscription costs if logic fails
    calculatedSubscriptionCosts = make(map[string]float64)
  }

  // Create the final mapping for response - simplified logic
  finalSubscriptionCosts := make(map[string]float64)
  for _, expenditure := range monthlyExpendituresDB {
    expenditureMonth, err := time.Parse("2006-01-02", expenditure.Date)
    if err != nil {
      dda.logger.Error("Failed to parse expenditure date", map[string]any{
        "error": err,
        "date": expenditure.Date,
      })
      continue
    }

    // Use the same format as calculatedSubscriptionCosts for direct lookup
    monthKey := expenditureMonth.Format("2006-01-01")
    if cost, exists := calculatedSubscriptionCosts[monthKey]; exists {
      finalSubscriptionCosts[expenditure.Date] = cost
      dda.logger.Debug("Found subscription cost for month", map[string]any{
        "expenditureDate": expenditure.Date,
        "monthKey": monthKey,
        "subscriptionCost": cost,
      })
    } else {
      finalSubscriptionCosts[expenditure.Date] = 0.0
      dda.logger.Debug("No subscription cost found for month", map[string]any{
        "expenditureDate": expenditure.Date,
        "monthKey": monthKey,
        "availableKeys": calculatedSubscriptionCosts,
      })
    }
  }

  dda.logger.Debug("Final subscription costs mapping", map[string]any{
    "finalSubscriptionCosts": finalSubscriptionCosts,
    "calculatedSubscriptionCosts": calculatedSubscriptionCosts,
  })

  // Transformations
  gameStats := dda.transformGameStatsDBToResponse(gameStatsDB)
  subscriptionStats := dda.transformGameStatsDBToResponse(subscriptionStatsDB)
  // Override the value with the calculated current month subscription cost
  subscriptionStats.Value = currentMonthSubscriptionCost
  digitalLocationStats := dda.transformGameStatsDBToResponse(digitalLocationStatsDB)
  physicalLocationStats := dda.transformGameStatsDBToResponse(physicalLocationStatsDB)

  digitalLocations := make([]types.DashboardDigitalLocationBFFResponse, len(digitalLocationsDB))
  for i, db := range digitalLocationsDB {
      digitalLocations[i] = dda.transformDigitalLocationDBToResponse(db)
  }

  sublocations := make([]types.DashboardSublocationBFFResponse, len(sublocationsDB))
  for i, db := range sublocationsDB {
      sublocations[i] = dda.transformSublocationDBToResponse(db)
  }

  platformList := make([]types.DashboardPlatformBFFResponse, len(platformListDB))
  for i, db := range platformListDB {
      platformList[i] = dda.transformPlatformDBToResponse(db)
  }

  monthlyExpenditures := make([]types.DashboardMonthlyExpenditureBFFResponse, len(monthlyExpendituresDB))
  for i, db := range monthlyExpendituresDB {
      currentMonthSubscriptionCost := finalSubscriptionCosts[db.Date]
      monthlyExpenditures[i] = dda.transformMonthlyExpenditureDBToResponse(db, currentMonthSubscriptionCost)
  }

  // 8. Response assembly
  response := types.DashboardBFFResponse{
    GameStats:                   gameStats,
    SubscriptionStats:           subscriptionStats,
    DigitalLocationStats:        digitalLocationStats,
    PhysicalLocationStats:       physicalLocationStats,
    SubscriptionTotal:           subscriptionTotal,
    DigitalLocations:            digitalLocations,
    Sublocations:                sublocations,
    NewItemsThisMonth:           newItemsThisMonth,
    PlatformList:                platformList,
    MediaTypeDomains:            []string{"oneTimePurchase", "hardware", "dlc", "inGamePurchase", "subscription"}, // or fetch dynamically if needed
    MonthlyExpenditures:         monthlyExpenditures,
  }

  if err := tx.Commit(); err != nil {
    dda.logger.Error("Error committing transaction", map[string]any{"error": err})
    return types.DashboardBFFResponse{}, fmt.Errorf("failed to commit transaction: %w", err)
  }

  dda.logger.Debug("GetDashboardBFFResponse response", map[string]any{"response": response})
  return response, nil
}