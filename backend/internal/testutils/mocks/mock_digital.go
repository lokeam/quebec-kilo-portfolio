package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockDigitalValidator struct {
	ValidateDigitalLocationFunc func(digitalLocation models.DigitalLocation) (models.DigitalLocation, error)
}

func (m *MockDigitalValidator) ValidateDigitalLocation(
	digitalLocation models.DigitalLocation,
) (models.DigitalLocation, error) {
	if m.ValidateDigitalLocationFunc != nil {
		return m.ValidateDigitalLocationFunc(digitalLocation)
	}
	return digitalLocation, nil
}

type MockDigitalDbAdapter struct {
	GetDigitalLocationFunc func(ctx context.Context, userID, digitalLocationID string) (models.DigitalLocation, error)
	GetDigitalLocationsFunc func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	AddDigitalLocationFunc func(ctx context.Context, userID string, digitalLocation models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocationFunc func(ctx context.Context, userID string, digitalLocation models.DigitalLocation) error
	DeleteDigitalLocationFunc func(ctx context.Context, userID, digitalLocationID string) error
	FindDigitalLocationByNameFunc func(ctx context.Context, userID string, name string) (models.DigitalLocation, error)

	// Subscription Operations
	GetSubscriptionFunc func(ctx context.Context, locationID string) (*models.Subscription, error)
	AddSubscriptionFunc func(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscriptionFunc func(ctx context.Context, subscription models.Subscription) error
	RemoveSubscriptionFunc func(ctx context.Context, locationID string) error

	// Payment Operations
	GetPaymentsFunc func(ctx context.Context, locationID string) ([]models.Payment, error)
	AddPaymentFunc func(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPaymentFunc func(ctx context.Context, paymentID int64) (*models.Payment, error)
}

// GET
func (m *MockDigitalDbAdapter) GetDigitalLocation(
	ctx context.Context,
	userID,
	digitalLocationID string,
) (models.DigitalLocation, error) {
	digitalLocation, err := m.GetDigitalLocationFunc(ctx, userID, digitalLocationID)
	return digitalLocation, err
}

func (m *MockDigitalDbAdapter) GetUserDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocation, error) {
	return m.GetDigitalLocationsFunc(ctx, userID)
}

// POST
func (m *MockDigitalDbAdapter) AddDigitalLocation(
	ctx context.Context,
	userID string,
	digitalLocation models.DigitalLocation,
) (models.DigitalLocation, error) {
	return m.AddDigitalLocationFunc(ctx, userID, digitalLocation)
}

// PUT
func (m *MockDigitalDbAdapter) UpdateDigitalLocation(
	ctx context.Context,
	userID string,
	digitalLocation models.DigitalLocation,
) error {
	return m.UpdateDigitalLocationFunc(ctx, userID, digitalLocation)
}

// DELETE
func (m *MockDigitalDbAdapter) RemoveDigitalLocation(
	ctx context.Context,
	userID string,
	digitalLocationID string,
) error {
	return m.DeleteDigitalLocationFunc(ctx, userID, digitalLocationID)
}

// Find by name
func (m *MockDigitalDbAdapter) FindDigitalLocationByName(
	ctx context.Context,
	userID string,
	name string,
) (models.DigitalLocation, error) {
	return m.FindDigitalLocationByNameFunc(ctx, userID, name)
}

// Subscription Operations
func (m *MockDigitalDbAdapter) GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error) {
	return m.GetSubscriptionFunc(ctx, locationID)
}

func (m *MockDigitalDbAdapter) AddSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
	return m.AddSubscriptionFunc(ctx, subscription)
}

func (m *MockDigitalDbAdapter) UpdateSubscription(ctx context.Context, subscription models.Subscription) error {
	return m.UpdateSubscriptionFunc(ctx, subscription)
}

func (m *MockDigitalDbAdapter) RemoveSubscription(ctx context.Context, locationID string) error {
	return m.RemoveSubscriptionFunc(ctx, locationID)
}

// Payment Operations
func (m *MockDigitalDbAdapter) GetPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	return m.GetPaymentsFunc(ctx, locationID)
}

func (m *MockDigitalDbAdapter) AddPayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	return m.AddPaymentFunc(ctx, payment)
}

func (m *MockDigitalDbAdapter) GetPayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
	return m.GetPaymentFunc(ctx, paymentID)
}

// ---------

type MockDigitalCacheWrapper struct {
	GetCachedDigitalLocationsFunc           func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	SetCachedDigitalLocationsFunc           func(ctx context.Context, userID string, locations []models.DigitalLocation) error
	GetSingleCachedDigitalLocationFunc      func(ctx context.Context, userID, digitalLocationID string) (*models.DigitalLocation, bool, error)
	SetSingleCachedDigitalLocationFunc      func(ctx context.Context, userID string, location models.DigitalLocation) error
	InvalidateUserCacheFunc                 func(ctx context.Context, userID string) error
	InvalidateDigitalLocationCacheFunc      func(ctx context.Context, userID, digitalLocationID string) error

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
