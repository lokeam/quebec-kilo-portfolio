package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)


type MockDigitalCacheWrapper struct {
	GetCachedDigitalLocationsFunc           func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	SetCachedDigitalLocationsFunc           func(ctx context.Context, userID string, locations []models.DigitalLocation) error
	GetSingleCachedDigitalLocationFunc      func(ctx context.Context, userID, digitalLocationID string) (*models.DigitalLocation, bool, error)
	SetSingleCachedDigitalLocationFunc      func(ctx context.Context, userID string, location models.DigitalLocation) error
	InvalidateUserCacheFunc                 func(ctx context.Context, userID string) error
	InvalidateDigitalLocationCacheFunc      func(ctx context.Context, userID, digitalLocationID string) error
	InvalidateDigitalLocationsBulkFunc      func(ctx context.Context, userID string, locationIDs []string) error

	// Subscription caching
	GetCachedSubscriptionFunc               func(ctx context.Context, locationID string) (*models.Subscription, bool, error)
	SetCachedSubscriptionFunc               func(ctx context.Context, locationID string, subscription models.Subscription) error
	InvalidateSubscriptionCacheFunc         func(ctx context.Context, locationID string) error

	// Payment caching
	GetCachedPaymentsFunc                   func(ctx context.Context, locationID string) ([]models.Payment, error)
	SetCachedPaymentsFunc                   func(ctx context.Context, locationID string, payments []models.Payment) error
	InvalidatePaymentsCacheFunc             func(ctx context.Context, locationID string) error
}

// GET
func (m *MockDigitalCacheWrapper) GetCachedDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocation, error) {
	return m.GetCachedDigitalLocationsFunc(ctx, userID)
}

func (m *MockDigitalCacheWrapper) GetSingleCachedDigitalLocation(
	ctx context.Context,
	userID,
	digitalLocationID string,
) (*models.DigitalLocation, bool, error) {
	return m.GetSingleCachedDigitalLocationFunc(ctx, userID, digitalLocationID)
}

// SET
func (m *MockDigitalCacheWrapper) SetCachedDigitalLocations(
	ctx context.Context,
	userID string,
	digitalLocations []models.DigitalLocation,
) error {
	return m.SetCachedDigitalLocationsFunc(ctx, userID, digitalLocations)
}

func (m *MockDigitalCacheWrapper) SetSingleCachedDigitalLocation(
	ctx context.Context,
	userID string,
	digitalLocation models.DigitalLocation,
) error {
	return m.SetSingleCachedDigitalLocationFunc(ctx, userID, digitalLocation)
}

// CLEAR
func (m *MockDigitalCacheWrapper) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	return m.InvalidateUserCacheFunc(ctx, userID)
}

func (m *MockDigitalCacheWrapper) InvalidateDigitalLocationCache(
	ctx context.Context,
	userID,
	digitalLocationID string,
) error {
	return m.InvalidateDigitalLocationCacheFunc(ctx, userID, digitalLocationID)
}

func (m *MockDigitalCacheWrapper) InvalidateDigitalLocationsBulk(
	ctx context.Context,
	userID string,
	locationIDs []string,
) error {
	return m.InvalidateDigitalLocationsBulkFunc(ctx, userID, locationIDs)
}

// Subscription caching
func (m *MockDigitalCacheWrapper) GetCachedSubscription(
	ctx context.Context,
	locationID string,
) (*models.Subscription, bool, error) {
	return m.GetCachedSubscriptionFunc(ctx, locationID)
}

func (m *MockDigitalCacheWrapper) SetCachedSubscription(
	ctx context.Context,
	locationID string,
	subscription models.Subscription,
) error {
	return m.SetCachedSubscriptionFunc(ctx, locationID, subscription)
}

func (m *MockDigitalCacheWrapper) InvalidateSubscriptionCache(
	ctx context.Context,
	locationID string,
) error {
	return m.InvalidateSubscriptionCacheFunc(ctx, locationID)
}

// Payment caching
func (m *MockDigitalCacheWrapper) GetCachedPayments(
	ctx context.Context,
	locationID string,
) ([]models.Payment, error) {
	return m.GetCachedPaymentsFunc(ctx, locationID)
}

func (m *MockDigitalCacheWrapper) SetCachedPayments(
	ctx context.Context,
	locationID string,
	payments []models.Payment,
) error {
	return m.SetCachedPaymentsFunc(ctx, locationID, payments)
}

func (m *MockDigitalCacheWrapper) InvalidatePaymentsCache(
	ctx context.Context,
	locationID string,
) error {
	return m.InvalidatePaymentsCacheFunc(ctx, locationID)
}
