package types

// SpendingCategoryBFFResponseFINAL represents a spending category in the BFF response
type SpendingCategoryBFFResponseFINAL struct {
    Name  string  `json:"name"`
    Value float64 `json:"value"`
}

// MonthlySpendingBFFResponseFINAL represents monthly spending data in the BFF response
type MonthlySpendingBFFResponseFINAL struct {
    CurrentMonthTotal     float64                            `json:"currentMonthTotal"`
    LastMonthTotal        float64                            `json:"lastMonthTotal"`
    PercentageChange      float64                            `json:"percentageChange"`
    ComparisonDateRange   string                             `json:"comparisonDateRange"`
    SpendingCategories    []SpendingCategoryBFFResponseFINAL `json:"spendingCategories"`
}

// AnnualSpendingBFFResponseFINAL represents annual spending data in the BFF response
type AnnualSpendingBFFResponseFINAL struct {
    DateRange            string                                `json:"dateRange"`
    MonthlyExpenditures  []MonthlyExpenditureBFFResponseFINAL  `json:"monthlyExpenditures"`
    MedianMonthlyCost    float64                               `json:"medianMonthlyCost"`
}

// MonthlyExpenditureBFFResponseFINAL represents monthly expenditure in the BFF response
type MonthlyExpenditureBFFResponseFINAL struct {
    Month         string    `json:"month"`
    Expenditure   float64   `json:"expenditure"`
}

// SingleYearlyTotalBFFResponseFINAL represents yearly total in the BFF response
type SingleYearlyTotalBFFResponseFINAL struct {
    Year     int       `json:"year"`
    Amount   float64   `json:"amount"`
}

// AllYearlyTotalsBFFResponseFINAL represents all yearly totals in the BFF response
type AllYearlyTotalsBFFResponseFINAL struct {
    SubscriptionTotal  []SingleYearlyTotalBFFResponseFINAL `json:"subscriptionTotal"`
    OneTimeTotal       []SingleYearlyTotalBFFResponseFINAL `json:"oneTimeTotal"`
    CombinedTotal      []SingleYearlyTotalBFFResponseFINAL `json:"combinedTotal"`
}

// SpendingItemBFFResponseFINAL represents a spending item in the BFF response
type SpendingItemBFFResponseFINAL struct {
    ID                    string                                `json:"id"`
    Title                 string                                `json:"title"`
    Amount                float64                               `json:"amount"`
    SpendTransactionType  string                                `json:"spendTransactionType"`
    PaymentMethod         string                                `json:"paymentMethod"`
    MediaType             string                                `json:"mediaType"`
    CreatedAt             int64                                 `json:"createdAt"`
    UpdatedAt             int64                                 `json:"updatedAt"`
    IsActive              bool                                  `json:"isActive"`
	Provider              string                                `json:"provider,omitempty"`

		// Subscription specific fields
    BillingCycle          string                                `json:"billingCycle,omitempty"`
    NextBillingDate       int64                                 `json:"nextBillingDate,omitempty"`
    YearlySpending        []SingleYearlyTotalBFFResponseFINAL   `json:"yearlySpending,omitempty"`

    // One-time purchase specific fields
    IsDigital             bool                                  `json:"isDigital,omitempty"`
    IsWishlisted          bool                                  `json:"isWishlisted,omitempty"`
    PurchaseDate          int64                                 `json:"purchaseDate,omitempty"`
}

// SpendTrackingBFFResponseFINAL represents the complete BFF response
type SpendTrackingBFFResponseFINAL struct {
    TotalMonthlySpending    MonthlySpendingBFFResponseFINAL    `json:"totalMonthlySpending"`
    TotalAnnualSpending     AnnualSpendingBFFResponseFINAL     `json:"totalAnnualSpending"`
    CurrentTotalThisMonth   []SpendingItemBFFResponseFINAL     `json:"currentTotalThisMonth"`
    OneTimeThisMonth        []SpendingItemBFFResponseFINAL     `json:"oneTimeThisMonth"`
    RecurringNextMonth      []SpendingItemBFFResponseFINAL     `json:"recurringNextMonth"`
    YearlyTotals            AllYearlyTotalsBFFResponseFINAL    `json:"yearlyTotals"`
}

type SpendTrackingCalculatorCurrentMonthData struct {
    TotalMonthlySpending     float64
    SpendingCategories       []SpendTrackingCalculatorSpendingCategory
    SpendingItems            []SpendTrackingCalculatorSpendingItem
}

type SpendTrackingCalculatorSpendingCategory struct {
    SpendingCategoryName   string
    SpendingCategoryValue  float64
}

type SpendTrackingCalculatorSpendingItem struct {
    SpendingCategoryID      string
    SpendingItemName        string
    SpendingItemAmount      float64
    SpendingItemCategory    string
}