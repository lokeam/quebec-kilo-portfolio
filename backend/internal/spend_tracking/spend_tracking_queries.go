package spend_tracking

// Shared queries used by both SpendTrackingCalculator and SpendTrackingDbAdapter
const (
	GetActiveSubscriptionsQuery = `
		SELECT
				dl.id,
				dl.user_id,
				dl.name,
				dl.is_subscription,
				dl.is_active,
				dl.payment_method,
				dl.created_at,
				dl.updated_at,
				dls.billing_cycle,
				dls.cost_per_cycle,
				dls.anchor_date,
				dls.last_payment_date,
				dls.next_payment_date,
				dls.payment_method as subscription_payment_method
			FROM digital_locations dl
			LEFT JOIN digital_location_subscriptions dls ON dl.id = dls.digital_location_id
			WHERE dl.user_id = $1
			AND dl.is_subscription = true
			AND dl.is_active = true
		ORDER BY dls.next_payment_date ASC
	`

	GetCurrentMonthOneTimePurchasesQuery = `
		SELECT otp.*, sc.media_type as media_type
    FROM one_time_purchases otp
    LEFT JOIN spending_categories sc ON otp.spending_category_id = sc.id
    WHERE otp.user_id = $1
    AND EXTRACT(YEAR FROM purchase_date) = EXTRACT(YEAR FROM $2::timestamp)
    AND EXTRACT(MONTH FROM purchase_date) = EXTRACT(MONTH FROM $2::timestamp)
    ORDER BY purchase_date DESC
  `

	GetMonthlySpendingAggregatesQuery = `
		SELECT
			id,
			user_id,
			year,
			month,
			total_amount,
			subscription_amount,
			one_time_amount,
			category_amounts,
			created_at,
			updated_at
		FROM monthly_spending_aggregates
		WHERE user_id = $1
		ORDER BY year, month
  `

	GetYearlySpendingQuery = `
		SELECT * FROM yearly_spending_aggregates
		WHERE user_id = $1
		AND year >= EXTRACT(YEAR FROM CURRENT_DATE)::int - 2
		ORDER BY year DESC
	`

	GetSubscriptionByIDQuery = `
    SELECT
        dl.id,
        dl.user_id,
        dl.name,
        dl.is_subscription,
        dl.is_active,
        dl.payment_method,
        dl.created_at,
        dl.updated_at,
        dls.billing_cycle,
        dls.cost_per_cycle,
        dls.anchor_date,
        dls.last_payment_date,
        dls.next_payment_date,
        dls.payment_method as subscription_payment_method
    FROM digital_locations dl
    LEFT JOIN digital_location_subscriptions dls ON dl.id = dls.digital_location_id
    WHERE dl.id = $1 AND dl.user_id = $2 AND dl.is_subscription = true
`
)