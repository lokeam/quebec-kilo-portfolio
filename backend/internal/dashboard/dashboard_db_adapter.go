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
	"github.com/lokeam/qko-beta/internal/types"
)

// Adapter struct for Dashboard BFF
type DashboardDbAdapter struct {
  client *postgres.PostgresClient
  db     *sqlx.DB
  logger interfaces.Logger
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

  return &DashboardDbAdapter{
      client: client,
      db:     db,
      logger: appContext.Logger,
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
      SELECT 'Monthly Online Services Costs' AS title, 'coin' AS icon,
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
      WHERE dl.user_id = $1 AND dl.is_subscription = true AND dl.is_active = true
  `

  // Get digital storage locations count and last updated
  getDigitalLocationStatsQuery = `
      SELECT 'Digital Storage Locations' AS title, 'onlineServices' AS icon,
             COUNT(*) AS value, MAX(updated_at) AS last_updated
      FROM digital_locations
      WHERE user_id = $1
  `

  // Get physical storage locations count and last updated
  getPhysicalLocationStatsQuery = `
      SELECT 'Physical Storage Locations' AS title, 'package' AS icon,
          COUNT(*) AS value, MAX(updated_at) AS last_updated
      FROM physical_locations
      WHERE user_id = $1
  `

  // Get all digital locations with details
  getDigitalLocationsQuery = `
    SELECT
      dl.name,
      dl.url,
      COALESCE(dls.billing_cycle, '') AS billing_cycle,
      COALESCE(dls.cost_per_cycle, 0) AS monthly_fee,
      COUNT(dgl.id) AS stored_items,
      dl.is_subscription,
      dls.next_payment_date
    FROM digital_locations dl
    LEFT JOIN digital_location_subscriptions dls ON dl.id = dls.digital_location_id
    LEFT JOIN digital_game_locations dgl ON dl.id = dgl.digital_location_id
    LEFT JOIN user_games ug ON dgl.user_game_id = ug.id AND ug.user_id = $1
    WHERE dl.user_id = $1
    GROUP BY dl.id, dl.name, dl.url, dls.billing_cycle, dls.cost_per_cycle, dl.is_subscription, dls.next_payment_date
  `

  // Get all sublocations with parent location details
  getSublocationsQuery = `
      SELECT s.id AS sublocation_id, s.name AS sublocation_name, s.location_type AS sublocation_type,
             s.stored_items, pl.id AS parent_location_id, pl.name AS parent_location_name,
             pl.location_type AS parent_location_type, pl.bg_color AS parent_location_bg_color
      FROM sublocations s
      JOIN physical_locations pl ON s.physical_location_id = pl.id
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
        COALESCE((category_amounts->>'one_time_purchase')::DECIMAL(10,2), 0) AS one_time_purchase,
        COALESCE((category_amounts->>'hardware')::DECIMAL(10,2), 0) AS hardware,
        COALESCE((category_amounts->>'dlc')::DECIMAL(10,2), 0) AS dlc,
        COALESCE((category_amounts->>'in_game_purchase')::DECIMAL(10,2), 0) AS in_game_purchase,
        COALESCE((category_amounts->>'subscription')::DECIMAL(10,2), 0) AS subscription
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
  return types.DashboardStatBFFResponse{
    Title: db.Title,
    Icon: db.Icon,
    Value: db.Value,
    LastUpdated: db.LastUpdated.Unix(),
  }
}

func (dda *DashboardDbAdapter) transformDigitalLocationDBToResponse(
  db models.DashboardDigitalLocationDB,
) types.DashboardDigitalLocationBFFResponse {
  renewsNextMonth := false
  if db.IsSubscription && db.NextPaymentDate != nil {
    now := time.Now()
    nextMonth := now.AddDate(0, 1, 0)

    // Double check if next_payment_date col is indeed next month
    if db.NextPaymentDate.Year() == nextMonth.Year() && db.NextPaymentDate.Month() == nextMonth.Month() {
      renewsNextMonth = true
    }
  }

  return types.DashboardDigitalLocationBFFResponse{
    Logo: db.Name,
    Name: db.Name,
    Url:  db.Url,
    BillingCycle: db.BillingCycle,
    MonthlyFee: db.MonthlyFee,
    RenewsNextMonth: renewsNextMonth,
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
) types.DashboardMonthlyExpenditureBFFResponse {
  return types.DashboardMonthlyExpenditureBFFResponse{
      Date:            db.Date,
      OneTimePurchase: db.OneTimePurchase,
      Hardware:        db.Hardware,
      Dlc:             db.Dlc,
      InGamePurchase:  db.InGamePurchase,
      Subscription:    db.Subscription,
  }
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
  subscriptionTotal := subscriptionStatsDB.Value

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

  // 6. Monthly Expenditures
  var monthlyExpendituresDB []models.DashboardMonthlyExpenditureDB
  if err := tx.SelectContext(ctx, &monthlyExpendituresDB, getMonthlyExpendituresQuery, userID); err != nil {
    dda.logger.Error("Error fetching monthly expenditures", map[string]any{"error": err})
    return types.DashboardBFFResponse{}, fmt.Errorf("error fetching monthly expenditures: %w", err)
  }

  // 7. Transformation - THIS NEEDS TO BE ADJUSTED ACCORDING TO ITEMS 1-6 AND THE TRANSFORMATION HELPER FNS.
  gameStats := dda.transformGameStatsDBToResponse(gameStatsDB)
  subscriptionStats := dda.transformGameStatsDBToResponse(subscriptionStatsDB)
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
      monthlyExpenditures[i] = dda.transformMonthlyExpenditureDBToResponse(db)
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