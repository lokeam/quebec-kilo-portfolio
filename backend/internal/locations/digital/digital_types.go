package digital

// DigitalServiceCatalogItem represents a digital service in the catalog
type DigitalServiceCatalogItem struct {
    ID          string     `json:"id"`    // Unique identifier for the service
    Name        string     `json:"name"`  // Display name of the service
    Logo        string     `json:"logo"`  // Logo identifier used for image lookup
}