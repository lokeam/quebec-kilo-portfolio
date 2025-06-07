package types

import (
	"time"
)

// GeneralStats contains high-level metrics about the user's library
type AnalyticsGeneralStatsResponse struct {
	TotalGames              int     `json:"total_games" db:"total_games"`
	MonthlySubscriptionCost float64 `json:"monthly_subscription_cost" db:"monthly_subscription_cost"`
	TotalDigitalLocations   int     `json:"total_digital_locations" db:"total_digital_locations"`
	TotalPhysicalLocations  int     `json:"total_physical_locations" db:"total_physical_locations"`
}

// FinancialStats contains detailed financial information
type AnalyticsFinancialStatsResponse struct {
	AnnualSubscriptionCost float64                            `json:"annual_subscription_cost" db:"annual_subscription_cost"`
	TotalServices          int                                `json:"total_services" db:"total_services"`
	RenewalsThisMonth      int                                `json:"renewals_this_month" db:"renewals_this_month"`
	Services               []AnalyticsServiceDetailsResponse  `json:"services"`
}

// ServiceDetails contains information about a digital service subscription
type AnalyticsServiceDetailsResponse struct {
	Name         string  `json:"name"`
	MonthlyFee   float64 `json:"monthly_fee"`
	BillingCycle string  `json:"billing_cycle"`
	NextPayment  string  `json:"next_payment"`
}

// StorageStats contains information about storage locations
type AnalyticsStorageStatsResponse struct {
	TotalPhysicalLocations int                                  `json:"total_physical_locations" db:"total_physical_locations"`
	TotalDigitalLocations  int                                  `json:"total_digital_locations" db:"total_digital_locations"`
	DigitalLocations       []AnalyticsDigitalLocationResponse   `json:"digital_locations"`
	PhysicalLocations      []AnalyticsPhysicalLocationResponse  `json:"physical_locations"`
}

type AnalyticsDigitalLocationResponse struct {
	ID                              string                                  `json:"id"`
	Name                            string                                  `json:"name"`
	LocationType                    string                                  `json:"location_type"`
	IsActive                        bool                                    `json:"is_active"`
	URL                             string                                  `json:"url"`
	CreatedAt                       time.Time                               `json:"created_at"`
	UpdatedAt                       time.Time                               `json:"updated_at"`
	ItemCount                       int                                     `json:"item_count"`
	IsSubscription                  bool                                    `json:"is_subscription"`
	MonthlyCost                     float64                                 `json:"monthly_cost"`
	Items                           []AnalyticsGameLibraryItemSummary       `json:"items"`
	PaymentMethod                   string                                  `json:"payment_method"`
	PaymentDate                     *time.Time                              `json:"payment_date"`
	BillingCycle                    string                                  `json:"billing_cycle"`
	CostPerCycle                    float64                                 `json:"cost_per_cycle"`
	NextPaymentDate                 *time.Time                              `json:"next_payment_date"`
	ParentPhysicalLocationName      string                                  `json:"parent_physical_location_name"`
	ParentPhysicalLocationBgColor   string                                  `json:"parent_physical_location_bg_color"`
	ParentPhysicalLocationType      string                                  `json:"parent_physical_location_type"`
}

type AnalyticsPhysicalLocationResponse struct {
	ID               string                                     `json:"id"`
	Name             string                                     `json:"name"`
	LocationType     string                                     `json:"location_type"`
	MapCoordinates   AnalyticsPhysicalMapCoordinatesResponse    `json:"map_coordinates"`
	BgColor          string                                     `json:"bg_color"`
	CreatedAt        time.Time                                  `json:"created_at"`
	UpdatedAt        time.Time                                  `json:"updated_at"`
	ItemCount        int                                        `json:"item_count"`
	Sublocations    []AnalyticsSublocationSummaryResponse       `json:"sublocations"`
}

type AnalyticsPhysicalMapCoordinatesResponse struct {
	Coords          string     `json:"coords"`
	GoogleMapsLink  string     `json:"google_maps_link"`
}

// LocationSummary represents a summary of a storage location
type AnalyticsLocationSummaryResponse struct {
	ID              string                                    `json:"id"`
	Name            string                                    `json:"name"`
	LocationType    string                                    `json:"location_type"`
	ItemCount       int                                       `json:"item_count"`
	IsSubscription  bool                                      `json:"is_subscription,omitempty"`
	MonthlyCost     float64                                   `json:"monthly_cost,omitempty"`
	MapCoordinates  AnalyticsPhysicalMapCoordinatesResponse   `json:"map_coordinates,omitempty"`
	BgColor         string                                    `json:"bg_color,omitempty"`
	IsActive        bool                                      `json:"is_active,omitempty"`
	URL             string                                    `json:"url,omitempty"`
	CreatedAt       time.Time                                 `json:"created_at"`
	UpdatedAt       time.Time                                 `json:"updated_at"`
	Sublocations    []AnalyticsSublocationSummaryResponse     `json:"sublocations"`
	Items           []AnalyticsGameLibraryItemSummary         `json:"payment_method"`
	PaymentDate     time.Time                                 `json:"payment_date"`
	BillingCycle    string                                    `json:"billing_cycle"`
	CostPerCycle    float64                                   `json:"cost_per_cycle"`
	NextPaymentDate time.Time                                 `json:"next_payment_date"`
}

// SublocationSummary represents a summary of a sublocation
type AnalyticsSublocationSummaryResponse struct {
	ID          string                             `json:"id"`
	Name        string                             `json:"name"`
	LocationType string                            `json:"location_type"`
	StoredItems int                                `json:"stored_items"`
	CreatedAt   time.Time                          `json:"created_at"`
	UpdatedAt   time.Time                          `json:"updated_at"`
	Items       []AnalyticsGameLibraryItemSummary  `json:"items,omitempty"`
}

// ItemSummary represents a summary of a game item
type AnalyticsGameLibraryItemSummary struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Platform       string    `json:"platform"`
	PlatformVersion string   `json:"platform_version"`
	AcquiredDate   time.Time `json:"acquired_date"`
}

// InventoryStats contains information about items by various criteria
type AnalyticsInventoryStatsResponse struct {
	TotalItemCount int                                    `json:"total_item_count" db:"total_item_count"`
	NewItemCount   int                                    `json:"new_item_count" db:"new_item_count"`
	PlatformCounts []AnalyticsPlatformItemCountResponse   `json:"platform_counts"`
}

// PlatformItemCount provides item counts per platform
type AnalyticsPlatformItemCountResponse struct {
	Platform  string `json:"platform"`
	ItemCount int    `json:"item_count"`
}

// WishlistStats contains information about wishlisted items
type AnalyticsWishlistStatsResponse struct {
	TotalWishlistItems   int     `json:"total_wishlist_items" db:"total_wishlist_items"`
	ItemsOnSale          int     `json:"items_on_sale" db:"items_on_sale"`
	StarredItem          string  `json:"starred_item,omitempty"`
	StarredItemPrice     float64 `json:"starred_item_price,omitempty"`
	CheapestSaleDiscount float64 `json:"cheapest_sale_discount,omitempty"`
}

// AnalyticsResponse is the complete response returned by the analytics service
type AnalyticsResponse struct {
	General   *AnalyticsGeneralStatsResponse   `json:"general,omitempty"`
	Financial *AnalyticsFinancialStatsResponse `json:"financial,omitempty"`
	Storage   *AnalyticsStorageStatsResponse   `json:"storage,omitempty"`
	Inventory *AnalyticsInventoryStatsResponse `json:"inventory,omitempty"`
	Wishlist  *AnalyticsWishlistStatsResponse  `json:"wishlist,omitempty"`
}


/*
==================
	Homepage:
	- No CRUD operations, page only displays data
	- This is a dashboard that shows the following high level data:

  // Top row of dashboard cards show the following data:
	- Total number of games in library and last time this number was updated (this number needs to include platform specific games)
	  Example: PS5 Elden Ring and PS4 Elden Ring count as 2 games even though I only record one game in
		my database with two different platform versions.
	- Total monthly subscription cost for digital lcoations this current month
	- Total number of digital locations and last time this number was updated
	- Total number of sublocations where games are stored

	// Large left column card shows the following digital location data
	- Sum of all digital locations annual subscription cost
	- Total number of digital locations
	- The digital locations subscriptions that renew (are charged) this month
	- Listing of all digital locations that render the following data:
		* Digital location name
		* Digital location logo
		* Digital location URL
		* Digital location billing cycle
		* Digital location cost per cycle

	// Large right column card shows the following storage location data
	- Total number sublocations
	- Total number of digital locations
	- Listing of all sublocations that render the following data:
		* Sublocation name
		* Sublocation type
		* Number of games stored in this specific sublocation
		* Digital location name
		* Digital location logo
		* Number of games stored in this specific sublocation

	// Small left column card shows the following data:
	- Total number of games in library
	- Total number of games by platform
	  Example: 38 PC games, 4 PS5 games, 10 Switch games, 2 XBox games, etc

	// Small middle column card shows the following data:
	- Name of starred wishlist item
	- Best current price of wishlist item
	- Number of wishlist items on sale
	- The cheapest sale discount for a wishlist item on sale

	// Large right column card shows the following data:
	- Total monthly spending (subscription costs + one-time purchases) across games

==================
	SpendTracking
		- CRUD operations for One-time expenses
		// Large left column card shows the following data:
		* Total spending this month
		* Percentage rise or increase compared to last month
		* Total amount spent this month in the specific one-type purchase categories:
		  * Hardware
			* In Game Purchases
			* Physical Games (discs)
			* Digital Games (via a digital location)
			* Digital location Subscriptions

		// Large right column card shows the following data:
		- Total spending this year, month-by-month
		- Average monthly spending

	 // Accordion component shows a collapsible, month-by-month listing of:
			* Digital location name
			* Digital location logo
			* Digital location billing cycle
			* Digital location subscription payment date
			* Digital location subscription cost per cycle
			* One-time expense name
			* One-time expense icon
			* One-time expense amount
			* One-time expense date
			* One-time expense type

		// List of subscription expenses that will be charged next month
		* Digital location name
		* Digital location logo
		* Digital location billing cycle
		* Digital location subscription cost per cycle
		* Digital location subscription payment date
		* Digital location subscription cost per cycle


==================
	OnlineServices:
	- CRUD operations for DigitalLocations
	- Page shows:
	 * Total number of digital locations
	 * Total number of subscription digital locations
	 * Total number of non-subscription digital locations
	 * Listing of Digital locations defined in AnalyticsDigitalLocationResponse

	PhysicalLocations:
	- CRUD operations for Physical and Sublocations
	- Page shows:
		* Total number of (parent) physical locations
		* Total number of (child) sublocations

	  * Listing of SUBLOCATIONS
			- Each sublocation item renders:
				* Parent Location Name
				* Parent Location BgColor (sublocation shares parent location bg color)
				* Parent Location Type
				* Parent location map coordinates
				* Sublocation Name
				* Sublocation Type
				* Parent Location BgColor (sublocation shares parent location bg color)
				* Parent Location Type

	Library:
	- CRUD operations for LibraryGames
	- Page shows a listing of LibraryGames

	Wishlist:
	- No CRUD operations, page only displays data
	- Listing of wishlisted items
	- Name of wis

	- Media Storage
	- No CRUD operations, page only displays data
	- Accordion list of all digital locations and all games stored in each digital location
	- Accordion list of all parent physical locations. Accordion list of all sublocations in each parent physical location
	- Child accordion list of all sublocations and all games stored in each sublocation

*/