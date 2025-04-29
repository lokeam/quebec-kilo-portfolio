package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for analytics data retrieval
type Repository interface {
	// General stats methods
	GetGeneralStats(ctx context.Context, userID string) (*GeneralStats, error)

	// Financial stats methods
	GetFinancialStats(ctx context.Context, userID string) (*FinancialStats, error)

	// Storage stats methods
	GetStorageStats(ctx context.Context, userID string) (*StorageStats, error)

	// Inventory stats methods
	GetInventoryStats(ctx context.Context, userID string) (*InventoryStats, error)

	// Wishlist stats methods
	GetWishlistStats(ctx context.Context, userID string) (*WishlistStats, error)
}

// repository implements Repository
type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new analytics repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// GetGeneralStats retrieves general statistics for a user
func (r *repository) GetGeneralStats(ctx context.Context, userID string) (*GeneralStats, error) {
	stats := &GeneralStats{}

	// Get physical locations count
	var physicalLocationsCount int
	err := r.db.GetContext(ctx, &physicalLocationsCount,
		`SELECT COUNT(*) FROM physical_locations WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get physical locations count: %w", err)
	}
	stats.TotalPhysicalLocations = physicalLocationsCount

	// Get digital locations count
	var digitalLocationsCount int
	err = r.db.GetContext(ctx, &digitalLocationsCount,
		`SELECT COUNT(*) FROM digital_locations WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get digital locations count: %w", err)
	}
	stats.TotalDigitalLocations = digitalLocationsCount

	// Get monthly subscription cost
	var monthlyCost float64
	err = r.db.GetContext(ctx, &monthlyCost, `
		SELECT COALESCE(SUM(CASE
			WHEN s.billing_cycle = 'monthly' THEN s.cost_per_cycle
			WHEN s.billing_cycle = 'quarterly' THEN s.cost_per_cycle / 3
			WHEN s.billing_cycle = 'annually' THEN s.cost_per_cycle / 12
			ELSE 0
		END), 0)
		FROM digital_location_subscriptions s
		JOIN digital_locations l ON s.digital_location_id = l.id
		WHERE l.user_id = $1 AND l.is_active = true`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly subscription cost: %w", err)
	}
	stats.MonthlySubscriptionCost = monthlyCost

	// Get total games count (if user_games table exists)
	var gamesCount int
	err = r.db.GetContext(ctx, &gamesCount,
		`SELECT COUNT(*) FROM user_games WHERE user_id = $1`, userID)
	if err != nil {
		// This is expected to fail if user_games doesn't exist yet
		// Just log and continue with 0
		gamesCount = 0
	}
	stats.TotalGames = gamesCount

	return stats, nil
}

// GetFinancialStats retrieves financial statistics for a user
func (r *repository) GetFinancialStats(ctx context.Context, userID string) (*FinancialStats, error) {
	stats := &FinancialStats{}

	// Get annual subscription cost
	err := r.db.GetContext(ctx, &stats.AnnualSubscriptionCost, `
		SELECT COALESCE(SUM(CASE
			WHEN s.billing_cycle = 'monthly' THEN s.cost_per_cycle * 12
			WHEN s.billing_cycle = 'quarterly' THEN s.cost_per_cycle * 4
			WHEN s.billing_cycle = 'annually' THEN s.cost_per_cycle
			ELSE 0
		END), 0)
		FROM digital_location_subscriptions s
		JOIN digital_locations l ON s.digital_location_id = l.id
		WHERE l.user_id = $1 AND l.is_active = true`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get annual subscription cost: %w", err)
	}

	// Get renewals this month
	currentMonth := time.Now().Format("2006-01")
	err = r.db.GetContext(ctx, &stats.RenewalsThisMonth, `
		SELECT COUNT(*)
		FROM digital_location_subscriptions s
		JOIN digital_locations l ON s.digital_location_id = l.id
		WHERE l.user_id = $1
		AND l.is_active = true
		AND TO_CHAR(s.next_payment_date, 'YYYY-MM') = $2`, userID, currentMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get renewals this month: %w", err)
	}

	// Get total services count
	err = r.db.GetContext(ctx, &stats.TotalServices, `
		SELECT COUNT(*)
		FROM digital_locations
		WHERE user_id = $1 AND service_type = 'subscription'`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total services count: %w", err)
	}

	// Get service details
	services := []ServiceDetails{}
	rows, err := r.db.QueryxContext(ctx, `
		SELECT
			l.name,
			CASE
				WHEN s.billing_cycle = 'monthly' THEN s.cost_per_cycle
				WHEN s.billing_cycle = 'quarterly' THEN s.cost_per_cycle / 3
				WHEN s.billing_cycle = 'annually' THEN s.cost_per_cycle / 12
				ELSE 0
			END as monthly_fee,
			s.billing_cycle,
			TO_CHAR(s.next_payment_date, 'YYYY-MM-DD') as next_payment
		FROM digital_locations l
		JOIN digital_location_subscriptions s ON l.id = s.digital_location_id
		WHERE l.user_id = $1 AND l.is_active = true
		ORDER BY monthly_fee DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service details: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var service ServiceDetails
		err := rows.Scan(&service.Name, &service.MonthlyFee, &service.BillingCycle, &service.NextPayment)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service details: %w", err)
		}
		services = append(services, service)
	}

	stats.Services = services
	return stats, nil
}

// GetStorageStats retrieves storage location statistics for a user
func (r *repository) GetStorageStats(ctx context.Context, userID string) (*StorageStats, error) {
	stats := &StorageStats{}

	// Get physical and digital location counts
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(CASE WHEN location_type = 'physical' THEN 1 END) as physical_count,
			COUNT(CASE WHEN location_type = 'digital' THEN 1 END) as digital_count
		FROM (
			SELECT 'physical' as location_type FROM physical_locations WHERE user_id = $1
			UNION ALL
			SELECT 'digital' as location_type FROM digital_locations WHERE user_id = $1
		) locations`, userID).Scan(&stats.TotalPhysicalLocations, &stats.TotalDigitalLocations)
	if err != nil {
		return nil, fmt.Errorf("failed to get location counts: %w", err)
	}

	// Get digital locations with item counts
	digitalLocations := []LocationSummary{}
	rows, err := r.db.QueryxContext(ctx, `
		SELECT
			l.id,
			l.name,
			COUNT(dgl.id) as item_count,
			l.service_type as location_type,
			CASE WHEN s.id IS NOT NULL THEN true ELSE false END as is_subscription,
			COALESCE(CASE
				WHEN s.billing_cycle = 'monthly' THEN s.cost_per_cycle
				WHEN s.billing_cycle = 'quarterly' THEN s.cost_per_cycle / 3
				WHEN s.billing_cycle = 'annually' THEN s.cost_per_cycle / 12
				ELSE 0
			END, 0) as monthly_cost
		FROM digital_locations l
		LEFT JOIN digital_game_locations dgl ON l.id = dgl.digital_location_id
		LEFT JOIN digital_location_subscriptions s ON l.id = s.digital_location_id
		WHERE l.user_id = $1
		GROUP BY l.id, l.name, l.service_type, s.id, s.billing_cycle, s.cost_per_cycle
		ORDER BY item_count DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get digital locations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var location LocationSummary
		err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.ItemCount,
			&location.LocationType,
			&location.IsSubscription,
			&location.MonthlyCost,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan digital location: %w", err)
		}
		digitalLocations = append(digitalLocations, location)
	}
	stats.DigitalLocations = digitalLocations

	// Get physical locations with item counts
	physicalLocations := []LocationSummary{}
	rows, err = r.db.QueryxContext(ctx, `
		SELECT
			l.id,
			l.name,
			COUNT(pgl.id) as item_count,
			l.location_type
		FROM physical_locations l
		LEFT JOIN sublocations sl ON l.id = sl.physical_location_id
		LEFT JOIN physical_game_locations pgl ON sl.id = pgl.sublocation_id
		WHERE l.user_id = $1
		GROUP BY l.id, l.name, l.location_type
		ORDER BY item_count DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get physical locations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var location LocationSummary
		err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.ItemCount,
			&location.LocationType,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan physical location: %w", err)
		}
		physicalLocations = append(physicalLocations, location)
	}
	stats.PhysicalLocations = physicalLocations

	return stats, nil
}

// GetInventoryStats retrieves inventory statistics for a user
func (r *repository) GetInventoryStats(ctx context.Context, userID string) (*InventoryStats, error) {
	stats := &InventoryStats{}

	// Get total item count
	err := r.db.GetContext(ctx, &stats.TotalItemCount, `
		SELECT COUNT(*) FROM user_games WHERE user_id = $1`, userID)
	if err != nil {
		// This might fail if user_games table doesn't exist yet
		stats.TotalItemCount = 0
	}

	// Get new items this month
	currentMonth := time.Now().Format("2006-01")
	err = r.db.GetContext(ctx, &stats.NewItemCount, `
		SELECT COUNT(*)
		FROM user_games
		WHERE user_id = $1
		AND TO_CHAR(added_at, 'YYYY-MM') = $2`, userID, currentMonth)
	if err != nil {
		// This might fail if user_games table doesn't exist yet
		stats.NewItemCount = 0
	}

	// Get platform counts if user_games and game_platforms exist
	platformCounts := []PlatformItemCount{}
	rows, err := r.db.QueryxContext(ctx, `
		SELECT
			p.name as platform,
			COUNT(ug.id) as item_count
		FROM user_games ug
		JOIN game_platforms gp ON ug.game_id = gp.game_id
		JOIN platforms p ON gp.platform_id = p.id
		WHERE ug.user_id = $1
		GROUP BY p.name
		ORDER BY item_count DESC`, userID)
	if err == nil {
		defer rows.Close()

		for rows.Next() {
			var platformCount PlatformItemCount
			err := rows.Scan(&platformCount.Platform, &platformCount.ItemCount)
			if err != nil {
				return nil, fmt.Errorf("failed to scan platform count: %w", err)
			}
			platformCounts = append(platformCounts, platformCount)
		}
	}
	// If error, just use empty slice
	stats.PlatformCounts = platformCounts

	return stats, nil
}

// GetWishlistStats retrieves wishlist statistics for a user
func (r *repository) GetWishlistStats(ctx context.Context, userID string) (*WishlistStats, error) {
	stats := &WishlistStats{}

	// Check if wishlist table exists
	var tableExists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_name = 'wishlist'
		)`).Scan(&tableExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if wishlist table exists: %w", err)
	}

	if !tableExists {
		// Wishlist table doesn't exist yet, return empty stats
		return stats, nil
	}

	// Get total wishlist items
	err = r.db.GetContext(ctx, &stats.TotalWishlistItems, `
		SELECT COUNT(*) FROM wishlist WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist count: %w", err)
	}

	// Get items on sale
	err = r.db.GetContext(ctx, &stats.ItemsOnSale, `
		SELECT COUNT(*)
		FROM wishlist
		WHERE user_id = $1 AND is_on_sale = true`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items on sale count: %w", err)
	}

	// Get starred item if any
	var starredItem struct {
		Name  string
		Price float64
	}
	err = r.db.GetContext(ctx, &starredItem, `
		SELECT g.name, w.current_price
		FROM wishlist w
		JOIN games g ON w.game_id = g.id
		WHERE w.user_id = $1
		ORDER BY w.created_at
		LIMIT 1`, userID)
	if err == nil {
		stats.StarredItem = starredItem.Name
		stats.StarredItemPrice = starredItem.Price
	}

	// Get cheapest sale discount
	err = r.db.GetContext(ctx, &stats.CheapestSaleDiscount, `
		SELECT COALESCE(MAX((current_price - sale_price) / current_price * 100), 0)
		FROM wishlist
		WHERE user_id = $1 AND is_on_sale = true`, userID)
	if err != nil {
		stats.CheapestSaleDiscount = 0
	}

	return stats, nil
}