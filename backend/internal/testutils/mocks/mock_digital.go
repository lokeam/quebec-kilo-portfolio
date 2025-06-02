package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockDigitalValidator struct {
	ValidateDigitalLocationFunc func(digitalLocation models.DigitalLocation) (models.DigitalLocation, error)
	ValidateDigitalLocationsBulkFunc func(locations []models.DigitalLocation) ([]models.DigitalLocation, error)
}

func (m *MockDigitalValidator) ValidateDigitalLocation(
	digitalLocation models.DigitalLocation,
) (models.DigitalLocation, error) {
	if m.ValidateDigitalLocationFunc != nil {
		return m.ValidateDigitalLocationFunc(digitalLocation)
	}
	return digitalLocation, nil
}

func (m *MockDigitalValidator) ValidateDigitalLocationsBulk(
	locations []models.DigitalLocation,
) ([]models.DigitalLocation, error) {
	return m.ValidateDigitalLocationsBulkFunc(locations)
}

type MockDigitalDbAdapter struct {
	GetDigitalLocationFunc func(ctx context.Context, userID, digitalLocationID string) (models.DigitalLocation, error)
	GetDigitalLocationsFunc func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	AddDigitalLocationFunc func(ctx context.Context, userID string, digitalLocation models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocationFunc func(ctx context.Context, userID string, digitalLocation models.DigitalLocation) error
	RemoveDigitalLocationFunc func(ctx context.Context, userID string, locationIDs []string) (int64, error)
	FindDigitalLocationByNameFunc func(ctx context.Context, userID string, name string) (models.DigitalLocation, error)

	// Game Management Operations
	AddGameToDigitalLocationFunc func(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocationFunc func(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationIDFunc func(ctx context.Context, userID string, locationID string) ([]models.Game, error)

	// Subscription Operations
	GetSubscriptionFunc func(ctx context.Context, locationID string) (*models.Subscription, error)
	AddSubscriptionFunc func(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscriptionFunc func(ctx context.Context, subscription models.Subscription) error
	RemoveSubscriptionFunc func(ctx context.Context, locationID string) error
	ValidateSubscriptionExistsFunc func(ctx context.Context, locationID string) (*models.Subscription, error)

	// Payment Operations
	GetPaymentsFunc func(ctx context.Context, locationID string) ([]models.Payment, error)
	AddPaymentFunc func(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPaymentFunc func(ctx context.Context, paymentID int64) (*models.Payment, error)
}

// GET
func (m *MockDigitalDbAdapter) GetSingleDigitalLocation(
	ctx context.Context,
	userID,
	digitalLocationID string,
) (models.DigitalLocation, error) {
	digitalLocation, err := m.GetDigitalLocationFunc(ctx, userID, digitalLocationID)
	return digitalLocation, err
}

func (m *MockDigitalDbAdapter) GetAllDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocation, error) {
	return m.GetDigitalLocationsFunc(ctx, userID)
}

// POST
func (m *MockDigitalDbAdapter) CreateDigitalLocation(
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
func (m *MockDigitalDbAdapter) DeleteDigitalLocation(
	ctx context.Context,
	userID string,
	locationIDs []string,
) (int64, error) {
	if m.RemoveDigitalLocationFunc != nil {
		return m.RemoveDigitalLocationFunc(ctx, userID, locationIDs)
	}
	return 0, nil
}

// Find by name
func (m *MockDigitalDbAdapter) FindDigitalLocationByName(
	ctx context.Context,
	userID string,
	name string,
) (models.DigitalLocation, error) {
	return m.FindDigitalLocationByNameFunc(ctx, userID, name)
}

// Game Management Operations
func (m *MockDigitalDbAdapter) AddGameToDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
	gameID int64,
) error {
	if m.AddGameToDigitalLocationFunc != nil {
		return m.AddGameToDigitalLocationFunc(ctx, userID, locationID, gameID)
	}
	return nil
}

func (m *MockDigitalDbAdapter) RemoveGameFromDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
	gameID int64,
) error {
	if m.RemoveGameFromDigitalLocationFunc != nil {
		return m.RemoveGameFromDigitalLocationFunc(ctx, userID, locationID, gameID)
	}
	return nil
}

func (m *MockDigitalDbAdapter) GetGamesByDigitalLocationID(
	ctx context.Context,
	userID string,
	locationID string,
) ([]models.Game, error) {
	if m.GetGamesByDigitalLocationIDFunc != nil {
		return m.GetGamesByDigitalLocationIDFunc(ctx, userID, locationID)
	}
	return nil, nil
}

// Subscription Operations
func (m *MockDigitalDbAdapter) GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error) {
	return m.GetSubscriptionFunc(ctx, locationID)
}

func (m *MockDigitalDbAdapter) CreateSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
	return m.AddSubscriptionFunc(ctx, subscription)
}

func (m *MockDigitalDbAdapter) UpdateSubscription(ctx context.Context, subscription models.Subscription) error {
	return m.UpdateSubscriptionFunc(ctx, subscription)
}

func (m *MockDigitalDbAdapter) DeleteSubscription(ctx context.Context, locationID string) error {
	return m.RemoveSubscriptionFunc(ctx, locationID)
}

func (m *MockDigitalDbAdapter) ValidateSubscriptionExists(
	ctx context.Context,
	locationID string,
) (*models.Subscription, error) {
	if m.ValidateSubscriptionExistsFunc != nil {
		return m.ValidateSubscriptionExistsFunc(ctx, locationID)
	}
	return nil, nil
}

// Payment Operations
func (m *MockDigitalDbAdapter) GetAllPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	return m.GetPaymentsFunc(ctx, locationID)
}

func (m *MockDigitalDbAdapter) CreatePayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	return m.AddPaymentFunc(ctx, payment)
}

func (m *MockDigitalDbAdapter) GetSinglePayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
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
