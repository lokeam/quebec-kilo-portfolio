-- Seed spending categories
INSERT INTO spending_categories (name, type) VALUES
    ('hardware', 'hardware'),
    ('dlc', 'dlc'),
    ('in_game', 'in_game'),
    ('subscription', 'misc'),
    ('physical', 'disc'),
    ('disc', 'disc');

-- Seed one-time purchases
INSERT INTO one_time_purchases (
    user_id,
    title,
    amount,
    purchase_date,
    payment_method,
    spending_category_id,
    is_digital,
    is_wishlisted
) VALUES
    -- Hardware
    ('00000000-0000-0000-0000-000000000001', 'G.Skill Trident Z5 Neo RGB DDR5-6000', 219.98, '2024-03-20', 'visa', 1, false, true),

    -- DLC
    ('00000000-0000-0000-0000-000000000001', 'Sid Meier\'s Civilization VI', 19.99, '2024-03-21', 'paypal', 2, true, true),
    ('00000000-0000-0000-0000-000000000001', 'ELDEN RING Shadow of the Erdtree', 39.99, '2024-03-19', 'paypal', 2, true, true),

    -- In-game purchases
    ('00000000-0000-0000-0000-000000000001', 'Helldivers 2', 10.99, '2024-03-18', 'paypal', 3, true, false),
    ('00000000-0000-0000-0000-000000000001', 'Path of Exile 2', 20.00, '2024-03-20', 'paypal', 3, true, true),

    -- Physical/Disc
    ('00000000-0000-0000-0000-000000000001', 'The Legend of Zelda: Breath of the Wild', 29.99, '2024-03-22', 'visa', 5, false, true),
    ('00000000-0000-0000-0000-000000000001', 'Gradius V', 59.99, '2024-03-25', 'visa', 6, false, true);

-- Seed subscriptions
INSERT INTO digital_locations (
    id,
    user_id,
    name,
    is_subscription,
    is_active,
    payment_method
) VALUES
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'Playstation Plus', true, true, 'visa'),
    ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'Xbox Game Pass Ultimate', true, true, 'mastercard'),
    ('00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000001', 'Nintendo Switch Online', true, true, 'visa'),
    ('00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000001', 'Apple Arcade', true, true, 'mastercard'),
    ('00000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000001', 'Google Play Pass', true, true, 'mastercard');

INSERT INTO digital_location_subscriptions (
    digital_location_id,
    billing_cycle,
    cost_per_cycle,
    next_payment_date,
    payment_method
) VALUES
    ('00000000-0000-0000-0000-000000000002', '3 month', 6.66, '2024-03-01', 'visa'),
    ('00000000-0000-0000-0000-000000000003', '1 month', 14.99, '2024-03-01', 'mastercard'),
    ('00000000-0000-0000-0000-000000000004', '12 month', 3.99, '2024-03-01', 'visa'),
    ('00000000-0000-0000-0000-000000000005', '1 month', 6.99, '2024-03-01', 'mastercard'),
    ('00000000-0000-0000-0000-000000000006', '1 month', 5.99, '2024-03-01', 'mastercard');

-- Seed monthly spending aggregates
INSERT INTO monthly_spending_aggregates (
    user_id,
    year,
    month,
    total_amount,
    subscription_amount,
    one_time_amount,
    category_amounts
) VALUES
    ('00000000-0000-0000-0000-000000000001', 2024, 3, 1784.04, 356.02, 1428.02, '{
        "hardware": 534.04,
        "dlc": 267.02,
        "in_game": 178.01,
        "subscription": 356.02,
        "physical": 267.02,
        "disc": 178.01
    }'::jsonb);

-- Seed yearly spending aggregates
INSERT INTO yearly_spending_aggregates (
    user_id,
    year,
    total_amount,
    subscription_amount,
    one_time_amount
) VALUES
    ('00000000-0000-0000-0000-000000000001', 2022, 1299.83, 399.88, 899.95),
    ('00000000-0000-0000-0000-000000000001', 2023, 1759.93, 459.96, 1299.97),
    ('00000000-0000-0000-0000-000000000001', 2024, 1299.88, 499.92, 799.96);