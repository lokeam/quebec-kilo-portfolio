package models

import (
	"time"
)

// SpendTrackingCategoryDB represents a spending category in the database
type SpendTrackingCategoryDB struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Type string `db:"type"`
}

// SpendTrackingOneTimePurchaseDB represents a one-time purchase in the database
type SpendTrackingOneTimePurchaseDB struct {
	ID                int       `db:"id"`
	UserID            string    `db:"user_id"`
	Title             string    `db:"title"`
	Amount            float64   `db:"amount"`
	PurchaseDate      time.Time `db:"purchase_date"`
	PaymentMethod     string    `db:"payment_method"`
	CategoryID        int       `db:"spending_category_id"`
	IsDigital         bool      `db:"is_digital"`
	IsWishlisted      bool      `db:"is_wishlisted"`
	MediaType         string    `db:"media_type"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// SpendTrackingLocationDB represents a digital location (subscription) in the database
type SpendTrackingLocationDB struct {
	ID              string    `db:"id"`
	UserID          string    `db:"user_id"`
	Name            string    `db:"name"`
	IsSubscription  bool      `db:"is_subscription"`
	IsActive        bool      `db:"is_active"`
	PaymentMethod   string    `db:"payment_method"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	BillingCycle    string    `db:"billing_cycle"`
	CostPerCycle    float64   `db:"cost_per_cycle"`
	NextPaymentDate time.Time `db:"next_payment_date"`
	SubscriptionPaymentMethod string `db:"subscription_payment_method"`
}

// SpendTrackingSubscriptionDB represents subscription details in the database
type SpendTrackingSubscriptionDB struct {
	ID              int       `db:"id"`
	LocationID      string    `db:"digital_location_id"`
	BillingCycle    string    `db:"billing_cycle"`
	CostPerCycle    float64   `db:"cost_per_cycle"`
	NextPaymentDate time.Time `db:"next_payment_date"`
	PaymentMethod   string    `db:"payment_method"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// SpendTrackingMonthlyAggregateDB represents monthly spending aggregates in the database
type SpendTrackingMonthlyAggregateDB struct {
	ID                  int       `db:"id"`
	UserID              string    `db:"user_id"`
	Year                int       `db:"year"`
	Month               int       `db:"month"`
	TotalAmount         float64   `db:"total_amount"`
	SubscriptionAmount  float64 `db:"subscription_amount"`
	OneTimeAmount       float64   `db:"one_time_amount"`
	CategoryAmounts     []byte    `db:"category_amounts"` // JSONB
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
}

// SpendTrackingYearlyAggregateDB represents yearly spending aggregates in the database
type SpendTrackingYearlyAggregateDB struct {
	ID                  int       `db:"id"`
	UserID              string    `db:"user_id"`
	Year                int       `db:"year"`
	TotalAmount         float64   `db:"total_amount"`
	SubscriptionAmount  float64   `db:"subscription_amount"`
	OneTimeAmount       float64   `db:"one_time_amount"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
}