package digital

// Common queries used by the DigitalDbAdapter
const (
	GetLocationsWithSubscriptionDataQuery = `
		SELECT
			dl.*,
			'[]'::json as items,
			dls.id as sub_id,
			dls.billing_cycle,
			dls.cost_per_cycle,
			dls.anchor_date,
			dls.last_payment_date,
			dls.next_payment_date,
			dls.payment_method,
			dls.created_at as sub_created_at,
			dls.updated_at as sub_updated_at
		FROM digital_locations dl
		LEFT JOIN digital_location_subscriptions dls ON dls.digital_location_id = dl.id
		WHERE dl.user_id = $1
		ORDER BY dl.created_at
	`

	GetSubscriptionByLocationIDQuery =  `
		SELECT id, digital_location_id, billing_cycle, cost_per_cycle,
		  anchor_date, last_payment_date, next_payment_date, payment_method, created_at, updated_at
		FROM digital_location_subscriptions
		WHERE digital_location_id = $1
	`

	GetAllGamesInDigitalLocationQuery = `
		SELECT g.*
		FROM games g
		JOIN user_games ug ON ug.game_id = g.id
		JOIN digital_game_locations dgl ON dgl.user_game_id = ug.id
		WHERE dgl.digital_location_id = $1 AND ug.user_id = $2
	`

	//retrieves a specific payment by ID
	GetSinglePaymentQuery = `
		SELECT id, digital_location_id, amount, payment_date,
			payment_method, transaction_id, created_at
		FROM digital_location_payments
		WHERE id = $1
	`

	RecordPaymentQuery = `
		INSERT INTO digital_location_payments
			(digital_location_id, amount, payment_date, payment_method, transaction_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	UpdateSubscriptionQuery = `
		UPDATE digital_location_subscriptions
		SET billing_cycle = $1,
			cost_per_cycle = $2,
			anchor_date = $3,
			payment_method = $4,
			updated_at = $5
		WHERE digital_location_id = $6
	`

	UpdateSubscriptionLastPaymentDateQuery = `
		UPDATE digital_location_subscriptions
      SET last_payment_date = $1, updated_at = $2
    WHERE digital_location_id = $3
	`

	SubscriptionAnchorDateQuery = `
		INSERT INTO digital_location_subscriptions
			(digital_location_id, billing_cycle, cost_per_cycle,
			 anchor_date, payment_method, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, digital_location_id, billing_cycle, cost_per_cycle,
		  anchor_date, last_payment_date, next_payment_date, payment_method, created_at, updated_at
	`

	// BFF Response Queries
	GetAllDigitalLocationsBFFQuery = `
        SELECT
            dl.id,
            dl.name,
            dl.is_subscription,
            dl.is_active,
            dl.url,
            dl.payment_method,
            dl.created_at,
            dl.updated_at,
            COALESCE(dls.billing_cycle, '') AS billing_cycle,
            COALESCE(dls.cost_per_cycle, 0) AS cost_per_cycle,
            dls.next_payment_date,
            COUNT(dgl.id) AS stored_items
        FROM digital_locations dl
        LEFT JOIN digital_location_subscriptions dls ON dl.id = dls.digital_location_id
        LEFT JOIN digital_game_locations dgl ON dl.id = dgl.digital_location_id
        LEFT JOIN user_games ug ON dgl.user_game_id = ug.id AND ug.user_id = $1
        WHERE dl.user_id = $1
        GROUP BY dl.id, dl.name, dl.is_subscription, dl.is_active, dl.url, dl.payment_method,
                 dl.created_at, dl.updated_at, dls.billing_cycle, dls.cost_per_cycle, dls.next_payment_date
        ORDER BY dl.name
    `

	GetDigitalLocationGamesBFFQuery = `
			SELECT
					ug.id,
					g.name,
					p.name as platform,
					ug.is_unique_copy,
					EXISTS (
							SELECT 1
							FROM user_games ug2
							WHERE ug2.game_id = ug.game_id
							AND ug2.platform_id = ug.platform_id
							AND ug2.game_type = 'physical'
					) as has_physical_copy
			FROM digital_game_locations dgl
			JOIN user_games ug ON dgl.user_game_id = ug.id
			JOIN games g ON ug.game_id = g.id
			JOIN platforms p ON ug.platform_id = p.id
			WHERE dgl.digital_location_id = $1 AND ug.user_id = $2
			ORDER BY g.name
	`
)