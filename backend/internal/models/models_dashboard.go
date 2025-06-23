package models

import (
	"time"
)

// DashboardGameStatsDB represents the raw DB result for a statistics card (games, subscriptions, locations)
type DashboardGameStatsDB struct {
	Title           string      `db:"title"`
	Icon            string      `db:"icon"`
	Value           float64     `db:"value"`
	SecondaryValue  float64     `db:"secondary_value"`
	LastUpdated     time.Time   `db:"last_updated"`
}

// DashboardDigitalLocationDB represents a digital storage location from the DB
// (joined from digital_locations, digital_location_subscriptions, etc)
type DashboardDigitalLocationDB struct {
	Logo             string      `db:"logo"`
	Name             string      `db:"name"`
	Url              string      `db:"url"`
	BillingCycle     string      `db:"billing_cycle"`
	MonthlyFee       float64     `db:"monthly_fee"`
	StoredItems      int         `db:"stored_items"`
	IsSubscription   bool        `db:"is_subscription"`
	NextPaymentDate  *time.Time  `db:"next_payment_date"`
}

// DashboardSublocationDB represents a physical sublocation from the DB
// (joined from sublocations, physical_locations, etc)
type DashboardSublocationDB struct {
	SublocationId         string `db:"sublocation_id"`
	SublocationName       string `db:"sublocation_name"`
	SublocationType       string `db:"sublocation_type"`
	StoredItems           int    `db:"stored_items"`
	ParentLocationId      string `db:"parent_location_id"`
	ParentLocationName    string `db:"parent_location_name"`
	ParentLocationType    string `db:"parent_location_type"`
	ParentLocationBgColor string `db:"parent_location_bg_color"`
}

// DashboardPlatformDB represents a platform and its item count from the DB
type DashboardPlatformDB struct {
	Platform  string `db:"platform"`
	ItemCount int    `db:"item_count"`
}

// DashboardMonthlyExpenditureDB represents a single month's expenditures from the DB
type DashboardMonthlyExpenditureDB struct {
	Date            string    `db:"date"`
	OneTimePurchase float64   `db:"one_time_purchase"`
	Hardware        float64   `db:"hardware"`
	Dlc             float64   `db:"dlc"`
	InGamePurchase  float64   `db:"in_game_purchase"`
	Subscription    float64   `db:"subscription"`
}