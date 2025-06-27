package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

func (m *MockDigitalDbAdapter) GetAllPayments(
	ctx context.Context,
	locationID string,
) ([]models.Payment, error) {
	return m.GetPaymentsFunc(ctx, locationID)
}

func (m *MockDigitalDbAdapter) CreatePayment(
	ctx context.Context,
	payment models.Payment,
) (*models.Payment, error) {
	return m.AddPaymentFunc(ctx, payment)
}

func (m *MockDigitalDbAdapter) GetSinglePayment(
	ctx context.Context,
	paymentID int64,
) (*models.Payment, error) {
	return m.GetPaymentFunc(ctx, paymentID)
}
