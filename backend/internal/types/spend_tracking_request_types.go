package types

type SpendTrackingRequest struct {
	ID                    string         `json:"id,omitempty"`
	Title                 string         `json:"title"`
	Amount                float64        `json:"amount"`
	SpendingCategoryID    int            `json:"spending_category_id"`
	PaymentMethod         string         `json:"payment_method"`
	PurchaseDate          string         `json:"purchase_date"`
	DigitalLocationID     *string        `json:"digital_location_id,omitempty"`
	IsWishlisted          *bool          `json:"is_wishlisted,omitempty"`
	IsDigital             *bool          `json:"is_digital,omitempty"`
}