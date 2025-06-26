package digital

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

// GetSinglePayment retrieves a specific payment by ID
func (da *DigitalDbAdapter) GetSinglePayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
	da.logger.Debug("GetSinglePayment called", map[string]any{
		"paymentID": paymentID,
	})

	var payment models.Payment
	err := da.db.GetContext(
		ctx,
		&payment,
		GetSinglePaymentQuery,
		paymentID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("error getting payment: %w", err)
	}

	return &payment, nil
}


// GetAllPayments retrieves all payments for a digital location
func (da *DigitalDbAdapter) GetAllPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	da.logger.Debug("GetAllPayments called", map[string]any{
		"locationID": locationID,
	})

	var payments []models.Payment
	err := da.db.SelectContext(
		ctx,
		&payments,
		GetAllPaymentsQuery,
		locationID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting payments: %w", err)
	}

	return payments, nil
}


// CreatePayment records a new payment for a digital location
func (da *DigitalDbAdapter) CreatePayment(
	ctx context.Context, payment models.Payment,
) (*models.Payment, error) {
	da.logger.Debug("CreatePayment called", map[string]any{
		"payment": payment,
	})

	payment.CreatedAt = time.Now()

	err := da.db.QueryRowxContext(
		ctx,
		CreatePaymentQuery,
		payment.LocationID,
		payment.Amount,
		payment.PaymentDate,
		payment.PaymentMethod,
		payment.TransactionID,
		payment.CreatedAt,
	).StructScan(&payment)

	if err != nil {
		return nil, fmt.Errorf("error adding payment: %w", err)
	}

	return &payment, nil
}


func (da *DigitalDbAdapter) UpdatePayment(
	ctx context.Context,
	payment models.Payment,
) error {
	payment.UpdatedAt = time.Now()

	result, err := da.db.ExecContext(
		ctx,
		UpdatePaymentQuery,
		payment.Amount,
		payment.PaymentDate,
		payment.PaymentMethod,
		payment.TransactionID,
		payment.UpdatedAt,
		payment.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating payment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}
