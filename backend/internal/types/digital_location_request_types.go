package types

type DigitalLocationRequest struct {
	ID             string                               `json:"id"`
	Name           string                               `json:"name"`
	IsActive       bool                                 `json:"isActive"`
	URL            string                               `json:"url"`
	IsSubscription bool                                 `json:"isSubscription"`
	PaymentMethod  string                               `json:"payment_method"`
	Subscription   *DigitalLocationRequestSubscription  `json:"subscription,omitempty"` // pointer use for optional field
}

type DigitalLocationRequestSubscription struct {
	BillingCycle   string    `json:"billing_cycle"`
	CostPerCycle   float64   `json:"cost_per_cycle"`
	AnchorDate     string    `json:"anchor_date"`
	PaymentMethod  string    `json:"payment_method"`
}