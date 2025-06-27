package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

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
