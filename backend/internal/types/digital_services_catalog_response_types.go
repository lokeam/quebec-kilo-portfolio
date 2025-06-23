package types

// DigitalServicesCatalogResponse represents the response format for the digital services catalog
type DigitalServicesCatalogResponse struct {
	Success bool                   `json:"success"`
	Catalog []DigitalServiceItem   `json:"catalog"`
}

// DigitalServiceItem represents a single digital service in the catalog
type DigitalServiceItem struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Logo                  string `json:"logo"`
	IsSubscriptionService bool   `json:"is_subscription_service"`
	URL                   string `json:"url"`
}