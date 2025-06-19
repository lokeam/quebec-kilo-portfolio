package models

import (
	"encoding/json"
	"time"
)

type DigitalLocation struct {
	ID               string            `json:"id" db:"id"`
	UserID           string            `json:"user_id" db:"user_id"`
	Name             string            `json:"name" db:"name"`
	IsSubscription   bool              `json:"is_subscription" db:"is_subscription"`
	IsActive         bool              `json:"is_active" db:"is_active"`
	URL              string            `json:"url" db:"url"`
	PaymentMethod    string            `json:"payment_method" db:"payment_method"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at" db:"updated_at"`
	Items            []Game            `json:"items" db:"items"`
	Subscription     *Subscription     `json:"subscription,omitempty"`
	DiskSizeValue    *float64          `json:"disk_size_value" db:"disk_size_value"`
	DiskSizeUnit     *string           `json:"disk_size_unit" db:"disk_size_unit"`
}

// Payment model
type Payment struct {
	ID            int64     `json:"id" db:"id"`
	LocationID    string    `json:"location_id" db:"digital_location_id"`
	Amount        float64     `json:"amount" db:"amount"`
	PaymentDate   time.Time `json:"payment_date" db:"payment_date"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	TransactionID string    `json:"transaction_id,omitempty" db:"transaction_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// IsSubscriptionService returns true if the location's service type is 'subscription'
func (dl *DigitalLocation) IsSubscriptionService() bool {
	return dl.IsSubscription || dl.Subscription != nil
}

// MarshalJSON implements json.Marshaler for DigitalLocation
// to include computed fields like isSubscriptionService
func (dl DigitalLocation) MarshalJSON() ([]byte, error) {
	type Alias DigitalLocation
	return json.Marshal(&struct {
		Alias
		IsSubscriptionService bool `json:"isSubscriptionService"`
	}{
		Alias:                Alias(dl),
		IsSubscriptionService: dl.IsSubscriptionService(),
	})
}

// Subscription model
type Subscription struct {
	ID               int64       `json:"id" db:"id"`
	LocationID       string      `json:"location_id" db:"digital_location_id"`
	BillingCycle     string      `json:"billing_cycle" db:"billing_cycle"`
	CostPerCycle     float64     `json:"cost_per_cycle" db:"cost_per_cycle"`
	AnchorDate       time.Time   `json:"anchor_date" db:"anchor_date"`
	LastPaymentDate  *time.Time  `json:"last_payment_date,omitempty" db:"last_payment_date"`
	NextPaymentDate  time.Time   `json:"next_payment_date" db:"next_payment_date"` // Computed
	PaymentMethod    string      `json:"payment_method" db:"payment_method"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
}

// UnmarshalJSON implements json.Unmarshaler for Subscription
func (s *Subscription) UnmarshalJSON(data []byte) error {
	type Alias Subscription
	aux := &struct {
		NextPaymentDate string `json:"next_payment_date"`
		AnchorDate      string `json:"anchor_date"`
		LastPaymentDate string `json:"last_payment_date"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse NextPaymentDate
	if aux.NextPaymentDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05Z", aux.NextPaymentDate)
		if err != nil {
			return err
		}
		s.NextPaymentDate = t
	}

	// Parse AnchorDate
	if aux.AnchorDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05Z", aux.AnchorDate)
		if err != nil {
			return err
		}
		s.AnchorDate = t
	}

	// Parse LastPaymentDate (nullable)
	if aux.LastPaymentDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05Z", aux.LastPaymentDate)
		if err != nil {
			return err
		}
		s.LastPaymentDate = &t
	}

	return nil
}
