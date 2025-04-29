package analytics

// Domain names for analytics data categories
const (
	// General statistics about the user's inventory
	DomainGeneral = "general"

	// Financial information about subscriptions and costs
	DomainFinancial = "financial"

	// Storage locations information
	DomainStorage = "storage"

	// Inventory statistics about items by platform, location, etc.
	DomainInventory = "inventory"

	// Wishlist information about desired items
	DomainWishlist = "wishlist"
)

// Cache keys prefixes and formats
const (
	// Format: user:{id}:analytics:{domain}
	CacheKeyFormat = "user:%s:analytics:%s"
)

// Default TTL values in minutes for each domain
const (
	TTLGeneralMinutes   = 15
	TTLFinancialMinutes = 30
	TTLStorageMinutes   = 15
	TTLInventoryMinutes = 10
	TTLWishlistMinutes  = 20
)
