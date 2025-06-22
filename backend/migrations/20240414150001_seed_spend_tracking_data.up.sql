-- Create test user
INSERT INTO users (id, email) VALUES
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 'test@example.com');

-- First ensure tables exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'spending_categories') THEN
        RAISE EXCEPTION 'Required tables do not exist. Please run the initial schema migration first.';
    END IF;
END $$;

-- Seed spending categories
INSERT INTO spending_categories (name, media_type) VALUES
    ('hardware', 'hardware'),
    ('dlc', 'dlc'),
    ('in_game_purchase', 'in_game_purchase'),
    ('subscription', 'subscription'),
    ('physical_game', 'physical_game'),
    ('digital_game', 'digital_game'),
    ('misc', 'misc');

-- Seed monthly spending aggregates
INSERT INTO monthly_spending_aggregates (
    user_id,
    year,
    month,
    total_amount,
    subscription_amount,
    one_time_amount,
    category_amounts,
    created_at,
    updated_at
) VALUES
    -- January 2025
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 2025, 1, 201.65, 201.65, 0, '{
        "hardware": 0,
        "dlc": 0,
        "in_game_purchase": 0,
        "subscription": 201.65,
        "physical": 0,
        "disc": 0
    }'::jsonb, '2025-01-16 00:00:00', '2025-01-16 00:00:00'),

    -- February 2025
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 2025, 2, 10.99, 10.99, 0, '{
        "hardware": 0,
        "dlc": 0,
        "in_game_purchase": 10.99,
        "subscription": 0,
        "physical": 0,
        "disc": 0
    }'::jsonb, '2025-02-16 00:00:00', '2025-02-16 00:00:00'),

    -- March 2025
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 2025, 3, 39.55, 39.55, 0, '{
        "hardware": 0,
        "dlc": 39.55,
        "in_game_purchase": 0,
        "subscription": 0,
        "physical": 0,
        "disc": 0
    }'::jsonb, '2025-03-16 00:00:00', '2025-03-16 00:00:00'),

    -- April 2025
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 2025, 4, 25.99, 25.99, 0, '{
        "hardware": 0,
        "dlc": 0,
        "in_game_purchase": 0,
        "subscription": 0,
        "physical": 25.99,
        "disc": 0
    }'::jsonb, '2025-04-16 00:00:00', '2025-04-16 00:00:00'),

    -- May 2025
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 2025, 5, 879.37, 356.02, 523.35, '{
        "hardware": 0,
        "dlc": 0,
        "in_game_purchase": 0,
        "subscription": 356.02,
        "physical_game": 0,
        "digital_game": 0
    }'::jsonb, '2025-05-16 00:00:00', '2025-05-16 00:00:00'),

    -- June 2025
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 2025, 6, 1784.0, 356.02, 1428.02, '{
        "hardware": 534.04,
        "dlc": 267.02,
        "in_game_purchase": 178.01,
        "subscription": 356.02,
        "physical_game": 267.02,
        "digital_game": 178.01
    }'::jsonb, '2025-06-16 00:00:00', '2025-06-16 00:00:00');

-- Seed yearly spending aggregates
INSERT INTO yearly_spending_aggregates (
    user_id,
    year,
    total_amount,
    subscription_amount,
    one_time_amount,
    created_at,
    updated_at
) VALUES
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 2023, 1299.83, 399.88, 899.95, '2025-06-16 00:00:00', '2025-06-16 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 2024, 1759.93, 459.96, 1299.97, '2025-06-16 00:00:00', '2025-06-16 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 2025, 1299.88, 499.92, 799.96, '2025-06-16 00:00:00', '2025-06-16 00:00:00');

-- Seed digital locations (subscriptions AND Steam)
INSERT INTO digital_locations (
    id,
    user_id,
    name,
    is_subscription,
    is_active,
    payment_method,
    created_at,
    updated_at
) VALUES
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', '9a4aeee6-fb31-4839-a921-f61b0525046d', 'Playstation Plus', true, true, 'visa', '2025-06-01 00:00:00', '2025-06-01 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b0525046e', '9a4aeee6-fb31-4839-a921-f61b0525046d', 'Xbox Game Pass Ultimate', true, true, 'mastercard', '2025-06-02 00:00:00', '2025-06-02 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b0525046f', '9a4aeee6-fb31-4839-a921-f61b0525046d', 'Nintendo Switch Online', true, true, 'visa', '2025-06-03 00:00:00', '2025-06-03 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b05250470', '9a4aeee6-fb31-4839-a921-f61b0525046d', 'Apple Arcade', true, true, 'mastercard', '2025-06-04 00:00:00', '2025-06-04 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b05250471', '9a4aeee6-fb31-4839-a921-f61b0525046d', 'Steam', false, true, 'paypal', '2025-06-16 00:00:00', '2025-06-16 00:00:00');

-- Seed one-time purchases
INSERT INTO one_time_purchases (
    user_id,
    title,
    amount,
    purchase_date,
    payment_method,
    spending_category_id,
    digital_location_id,
    is_digital,
    is_wishlisted,
    created_at,
    updated_at
) VALUES
    -- Hardware (no digital_location_id needed)
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 'G.Skill Trident Z5 Neo RGB DDR5-6000', 219.98, '2025-06-16', 'visa', 1, NULL, false, true, '2025-06-16 00:00:00', '2025-06-16 00:00:00'),

    -- DLC (Steam)
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 'Sid Meier''s Civilization VI', 19.99, '2025-06-17', 'paypal', 2, '9a4aeee6-fb31-4839-a921-f61b05250471', true, true, '2025-06-17 00:00:00', '2025-06-17 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 'ELDEN RING Shadow of the Erdtree', 39.99, '2025-06-18', 'paypal', 2, '9a4aeee6-fb31-4839-a921-f61b05250471', true, true, '2025-06-18 00:00:00', '2025-06-18 00:00:00'),

    -- In-game purchases (Steam)
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 'Helldivers 2', 10.99, '2025-06-19', 'paypal', 3, '9a4aeee6-fb31-4839-a921-f61b05250471', true, false, '2025-06-19 00:00:00', '2025-06-19 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 'Path of Exile 2', 20.00, '2025-06-19', 'paypal', 3, '9a4aeee6-fb31-4839-a921-f61b05250471', true, true, '2025-06-19 00:00:00', '2025-06-19 00:00:00'),

    -- Physical/Disc (no digital_location_id needed)
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 'The Legend of Zelda: Breath of the Wild', 29.99, '2025-06-20', 'visa', 5, NULL, false, true, '2025-06-20 00:00:00', '2025-06-20 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', 'Gradius V', 59.99, '2025-06-21', 'visa', 6, NULL, false, true, '2025-06-21 00:00:00', '2025-06-21 00:00:00');

-- Seed subscription details
INSERT INTO digital_location_subscriptions (
    digital_location_id,
    billing_cycle,
    cost_per_cycle,
    next_payment_date,
    payment_method,
    created_at,
    updated_at
) VALUES
    ('9a4aeee6-fb31-4839-a921-f61b0525046d', '3 month', 6.66, '2025-07-01', 'visa', '2025-06-01 00:00:00', '2025-06-01 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b0525046e', '1 month', 14.99, '2025-07-02', 'mastercard', '2025-06-02 00:00:00', '2025-06-02 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b0525046f', '12 month', 3.99, '2025-07-03', 'visa', '2025-06-03 00:00:00', '2025-06-03 00:00:00'),
    ('9a4aeee6-fb31-4839-a921-f61b05250470', '1 month', 6.99, '2025-07-04', 'mastercard', '2025-06-04 00:00:00', '2025-06-04 00:00:00');