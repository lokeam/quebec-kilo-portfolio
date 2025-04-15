package digital

// DigitalServiceCatalogItem represents a digital service in the catalog
type DigitalServiceCatalogItem struct {
    ID                        string     `json:"id"`    // Unique identifier for the service
    Name                      string     `json:"name"`  // Display name of the service
    Logo                      string     `json:"logo"`  // Logo identifier used for image lookup
	IsSubscriptionService     bool       `json:"is_subscription_service"` // Whether the service is a subscription service
    URL                       string     `json:"url"`                      // login URL of the service
}