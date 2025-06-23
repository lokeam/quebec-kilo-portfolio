package digital

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/postgres"
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

	query := `
		SELECT id, digital_location_id, amount, payment_date,
		       payment_method, transaction_id, created_at
		FROM digital_location_payments
		WHERE digital_location_id = $1
		ORDER BY payment_date DESC
	`

	var payments []models.Payment
	err := da.db.SelectContext(ctx, &payments, query, locationID)
	if err != nil {
		return nil, fmt.Errorf("error getting payments: %w", err)
	}

	return payments, nil
}


// CreatePayment records a new payment for a digital location
func (da *DigitalDbAdapter) CreatePayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	da.logger.Debug("CreatePayment called", map[string]any{
		"payment": payment,
	})

	query := `
		INSERT INTO digital_location_payments
			(digital_location_id, amount, payment_date,
			 payment_method, transaction_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, digital_location_id, amount, payment_date,
		          payment_method, transaction_id, created_at
	`

	payment.CreatedAt = time.Now()

	err := da.db.QueryRowxContext(
		ctx,
		query,
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


func (da *DigitalDbAdapter) RecordPayment(
	ctx context.Context,
	payment models.Payment,
) error {
		return postgres.WithTransaction(
			ctx,
			da.db,
			da.logger,
			func(tx *sqlx.Tx) error {
				// 1. Record the payment
				_, err := tx.ExecContext(
					ctx,
					RecordPaymentQuery,
					payment.LocationID,
					payment.Amount,
					payment.PaymentDate,
					payment.PaymentMethod,
					payment.TransactionID,
					time.Now(),
				)
				if err != nil {
					return fmt.Errorf("error recording payment: %w", err)
				}

				// 2. Update the subscription's last payment date
				result, err := tx.ExecContext(
					ctx,
					UpdateSubscriptionLastPaymentDateQuery,
					payment.PaymentDate,
					time.Now(),
					payment.LocationID,
				)
				if err != nil {
					return fmt.Errorf("error updating subscription last payment date: %w", err)
				}

				rowsAffected, err := result.RowsAffected()
				if err != nil {
					return fmt.Errorf("error getting rows affected: %w", err)
				}

				if rowsAffected == 0 {
					return fmt.Errorf("subscription not found for location: %s", payment.LocationID)
				}

				return nil
			})
}
