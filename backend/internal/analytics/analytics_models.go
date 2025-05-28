package analytics

import (
	"time"
)

// GeneralStats contains high-level metrics about the user's library
type GeneralStats struct {
	TotalGames              int     `json:"total_games" db:"total_games"`
	MonthlySubscriptionCost float64 `json:"monthly_subscription_cost" db:"monthly_subscription_cost"`
	TotalDigitalLocations   int     `json:"total_digital_locations" db:"total_digital_locations"`
	TotalPhysicalLocations  int     `json:"total_physical_locations" db:"total_physical_locations"`
}

// FinancialStats contains detailed financial information
type FinancialStats struct {
	AnnualSubscriptionCost float64          `json:"annual_subscription_cost" db:"annual_subscription_cost"`
	TotalServices          int              `json:"total_services" db:"total_services"`
	RenewalsThisMonth      int              `json:"renewals_this_month" db:"renewals_this_month"`
	Services               []ServiceDetails `json:"services"`
}

// ServiceDetails contains information about a digital service subscription
type ServiceDetails struct {
	Name         string  `json:"name"`
	MonthlyFee   float64 `json:"monthly_fee"`
	BillingCycle string  `json:"billing_cycle"`
	NextPayment  string  `json:"next_payment"`
}

// StorageStats contains information about storage locations
type StorageStats struct {
	TotalPhysicalLocations int                     `json:"total_physical_locations" db:"total_physical_locations"`
	TotalDigitalLocations  int                     `json:"total_digital_locations" db:"total_digital_locations"`
	DigitalLocations       []DigitalLocation       `json:"digital_locations"`
	PhysicalLocations      []PhysicalLocation      `json:"physical_locations"`
}

type DigitalLocation struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	LocationType  string           `json:"location_type"`
	IsActive      bool             `json:"is_active"`
	URL           string           `json:"url"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	ItemCount     int              `json:"item_count"`
	IsSubscription bool            `json:"is_subscription"`
	MonthlyCost   float64          `json:"monthly_cost"`
	Items         []ItemSummary    `json:"items"`
	// New fields we want to add
	PaymentMethod   string           `json:"payment_method"`
	PaymentDate     *time.Time       `json:"payment_date"`
	BillingCycle    string           `json:"billing_cycle"`
	CostPerCycle    float64          `json:"cost_per_cycle"`
	NextPaymentDate *time.Time       `json:"next_payment_date"`
}

type PhysicalLocation struct {
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	LocationType  string                  `json:"location_type"`
	MapCoordinates string                 `json:"map_coordinates"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
	ItemCount     int                     `json:"item_count"`
	Sublocations  []SublocationSummary    `json:"sublocations"`
}

// LocationSummary represents a summary of a storage location
type LocationSummary struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	LocationType    string    `json:"location_type"`
	ItemCount       int       `json:"item_count"`
	IsSubscription  bool      `json:"is_subscription,omitempty"`
	MonthlyCost     float64   `json:"monthly_cost,omitempty"`
	MapCoordinates  string    `json:"map_coordinates,omitempty"`
	IsActive        bool      `json:"is_active,omitempty"`
	URL             string    `json:"url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Sublocations    []SublocationSummary `json:"sublocations"`
	Items           []ItemSummary        `json:"items"`
	PaymentMethod   string               `json:"payment_method"`
	PaymentDate     time.Time            `json:"payment_date"`
	BillingCycle    string               `json:"billing_cycle"`
	CostPerCycle    float64              `json:"cost_per_cycle"`
	NextPaymentDate time.Time            `json:"next_payment_date"`
}

// SublocationSummary represents a summary of a sublocation
type SublocationSummary struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	LocationType string        `json:"location_type"`
	BgColor     string         `json:"bg_color"`
	StoredItems int           `json:"stored_items"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Items       []ItemSummary `json:"items,omitempty"`
}

// ItemSummary represents a summary of a game item
type ItemSummary struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Platform       string    `json:"platform"`
	PlatformVersion string   `json:"platform_version"`
	AcquiredDate   time.Time `json:"acquired_date"`
}

// InventoryStats contains information about items by various criteria
type InventoryStats struct {
	TotalItemCount int                  `json:"total_item_count" db:"total_item_count"`
	NewItemCount   int                  `json:"new_item_count" db:"new_item_count"`
	PlatformCounts []PlatformItemCount  `json:"platform_counts"`
}

// PlatformItemCount provides item counts per platform
type PlatformItemCount struct {
	Platform  string `json:"platform"`
	ItemCount int    `json:"item_count"`
}

// WishlistStats contains information about wishlisted items
type WishlistStats struct {
	TotalWishlistItems   int     `json:"total_wishlist_items" db:"total_wishlist_items"`
	ItemsOnSale          int     `json:"items_on_sale" db:"items_on_sale"`
	StarredItem          string  `json:"starred_item,omitempty"`
	StarredItemPrice     float64 `json:"starred_item_price,omitempty"`
	CheapestSaleDiscount float64 `json:"cheapest_sale_discount,omitempty"`
}

// AnalyticsResponse is the complete response returned by the analytics service
type AnalyticsResponse struct {
	General   *GeneralStats   `json:"general,omitempty"`
	Financial *FinancialStats `json:"financial,omitempty"`
	Storage   *StorageStats   `json:"storage,omitempty"`
	Inventory *InventoryStats `json:"inventory,omitempty"`
	Wishlist  *WishlistStats  `json:"wishlist,omitempty"`
}
