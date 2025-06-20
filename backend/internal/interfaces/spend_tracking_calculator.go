package interfaces

import (
	"time"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type SpendTrackingCalculator interface {
	// Core Subscription Logic
	CalculateMonthlySubscriptionCosts(userID string, targetMonth time.Time) (float64, error)
	IsSubscriptionDueInMonth(subscription models.SpendTrackingSubscriptionDB, targetMonth time.Time) (bool, error)
	CalculateMonthlyMinimumSpending(userID string, targetMonth time.Time) (float64, error)

	// Business Intelligence Logic
	CalculatePercentageChange(userID string, currentMonth time.Time) (float64, error)
	CalculateAnnualSpendingForecast(userID string, targetYear time.Time) (types.AnnualSpendingBFFResponseFINAL, error)
	CalculateCurrentMonthAggregation(userID string, targetMonth time.Time) (types.SpendTrackingCalculatorCurrentMonthData, error)

	// Historical Analysis Logic
	CalculateThreeYearSubscriptionCosts(userID string, targetYear time.Time) (map[int]float64, error)
	CalculatePerSubscriptionYearlyTotals(userID string, subscriptionID string) ([]types.SingleYearlyTotalBFFResponseFINAL, error)
	CalculateMedianMonthlyCost(monthlyExpenditures []types.MonthlyExpenditureBFFResponseFINAL) float64
}