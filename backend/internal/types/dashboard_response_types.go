package types

// DashboardStatBFFResponse represents a single statistics card (games, subscriptions, locations)
type DashboardStatBFFResponse struct {
    Title           string   `json:"title"`
    Icon            string   `json:"icon"`
    Value           float64  `json:"value"`
    SecondaryValue  float64  `json:"secondary_value,omitempty"`
    LastUpdated     int64    `json:"lastUpdated"`
}

// DashboardDigitalLocationBFFResponse represents a digital storage location in the dashboard response
type DashboardDigitalLocationBFFResponse struct {
    Logo            string  `json:"logo"`
    Name            string  `json:"name"`
    Url             string  `json:"url"`
    BillingCycle    string  `json:"billingCycle"`
    MonthlyFee      float64 `json:"monthlyFee"`
    StoredItems     int     `json:"storedItems"`
    RenewsNextMonth bool    `json:"renewsNextMonth"`
}

// DashboardSublocationBFFResponse represents a physical sublocation in the dashboard response
type DashboardSublocationBFFResponse struct {
    SublocationId         string `json:"sublocationId"`
    SublocationName       string `json:"sublocationName"`
    SublocationType       string `json:"sublocationType"`
    StoredItems           int    `json:"storedItems"`
    ParentLocationId      string `json:"parentLocationId"`
    ParentLocationName    string `json:"parentLocationName"`
    ParentLocationType    string `json:"parentLocationType"`
    ParentLocationBgColor string `json:"parentLocationBgColor"`
}

// DashboardPlatformBFFResponse represents a platform and its item count
type DashboardPlatformBFFResponse struct {
    Platform  string `json:"platform"`
    ItemCount int    `json:"itemCount"`
}

// DashboardMonthlyExpenditureBFFResponseFINAL represents a single month's expenditures
type DashboardMonthlyExpenditureBFFResponse struct {
    Date            string  `json:"date"`
    OneTimePurchase float64 `json:"oneTimePurchase"`
    Hardware        float64 `json:"hardware"`
    Dlc             float64 `json:"dlc"`
    InGamePurchase  float64 `json:"inGamePurchase"`
    Subscription    float64 `json:"subscription"`
}

// DashboardBFFResponseFINAL is the top-level dashboard response type
// This should match the frontend DashboardResponse structure
type DashboardBFFResponse struct {
    GameStats                   DashboardStatBFFResponse                   `json:"gameStats"`
    SubscriptionStats           DashboardStatBFFResponse                   `json:"subscriptionStats"`
    DigitalLocationStats        DashboardStatBFFResponse                   `json:"digitalLocationStats"`
    PhysicalLocationStats       DashboardStatBFFResponse                   `json:"physicalLocationStats"`
    SubscriptionTotal           float64                                    `json:"subscriptionTotal"`
    DigitalLocations            []DashboardDigitalLocationBFFResponse      `json:"digitalLocations"`
    Sublocations                []DashboardSublocationBFFResponse          `json:"sublocations"`
    NewItemsThisMonth           int                                        `json:"newItemsThisMonth"`
    PlatformList                []DashboardPlatformBFFResponse             `json:"platformList"`
    MediaTypeDomains            []string                                   `json:"mediaTypeDomains"`
    MonthlyExpenditures         []DashboardMonthlyExpenditureBFFResponse   `json:"monthlyExpenditures"`
}
