-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create schema_migrations table to track applied migrations
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create games table
CREATE TABLE games (
    id BIGINT PRIMARY KEY,  -- IGDB ID
    name VARCHAR(255) NOT NULL,
    summary TEXT,
    cover_id BIGINT,
    cover_url VARCHAR(255),
    first_release_date BIGINT,
    rating FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create platforms table
CREATE TABLE platforms (
    id BIGINT PRIMARY KEY,  -- IGDB ID
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('console', 'pc', 'mobile')),
    model VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create genres table
CREATE TABLE genres (
    id BIGINT PRIMARY KEY,  -- IGDB ID
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create themes table
CREATE TABLE themes (
    id BIGINT PRIMARY KEY,  -- IGDB ID
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create game_genres junction table
CREATE TABLE game_genres (
    game_id BIGINT REFERENCES games(id) ON DELETE CASCADE,
    genre_id BIGINT REFERENCES genres(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (game_id, genre_id)
);

-- Create game_themes junction table
CREATE TABLE game_themes (
    game_id BIGINT REFERENCES games(id) ON DELETE CASCADE,
    theme_id BIGINT REFERENCES themes(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (game_id, theme_id)
);

-- Create user_games table
CREATE TABLE user_games (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    game_id BIGINT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    platform_id BIGINT NOT NULL REFERENCES platforms(id) ON DELETE CASCADE,
    game_type VARCHAR(50) NOT NULL CHECK (game_type IN ('physical', 'digital')),
    favorite BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, game_id, platform_id)
);

-- Create physical_locations table
CREATE TABLE physical_locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    label VARCHAR(255),
    location_type VARCHAR(50) NOT NULL CHECK (location_type IN ('house', 'apartment', 'office', 'warehouse', 'vehicle')),
    map_coordinates VARCHAR(255),
    bg_color VARCHAR(50) NOT NULL CHECK (bg_color IN ('red', 'green', 'blue', 'orange', 'gold', 'purple', 'brown', 'pink', 'gray')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create sublocations table
CREATE TABLE sublocations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    physical_location_id UUID REFERENCES physical_locations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    location_type VARCHAR(50) NOT NULL CHECK (location_type IN ('shelf', 'console', 'cabinet', 'closet', 'drawer', 'box', 'device')),
    stored_items INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(physical_location_id, name)
);

-- Create physical_game_locations table
CREATE TABLE physical_game_locations (
    id SERIAL PRIMARY KEY,
    user_game_id INTEGER REFERENCES user_games(id) ON DELETE CASCADE,
    sublocation_id UUID REFERENCES sublocations(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_game_id, sublocation_id)
);

-- Create digital_locations table
CREATE TABLE digital_locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    is_subscription BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    -- Disk size fields with proper constraints
    disk_size_value DECIMAL(10,2) CHECK (disk_size_value >= 0),
    disk_size_unit VARCHAR(10) CHECK (disk_size_unit IN ('KB', 'MB', 'GB', 'TB')),
    url TEXT,
    payment_method VARCHAR(50) CHECK (payment_method IN ('alipay', 'amex', 'diners', 'discover', 'elo', 'generic', 'hiper', 'hipercard', 'jcb', 'maestro', 'mastercard', 'mir', 'paypal', 'unionpay', 'visa')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name)
);

-- Create digital_location_subscriptions table
CREATE TABLE digital_location_subscriptions (
    id SERIAL PRIMARY KEY,
    digital_location_id UUID NOT NULL REFERENCES digital_locations(id) ON DELETE CASCADE,
    billing_cycle VARCHAR(50) NOT NULL CHECK (billing_cycle IN ('1 month', '3 month', '6 month', '12 month')),
    cost_per_cycle DECIMAL(10,2) NOT NULL CHECK (cost_per_cycle > 0),
    next_payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    payment_method VARCHAR(50) NOT NULL CHECK (payment_method IN ('alipay', 'amex', 'diners', 'discover', 'elo', 'generic', 'hiper', 'hipercard', 'jcb', 'maestro', 'mastercard', 'mir', 'paypal', 'unionpay', 'visa')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(digital_location_id)  -- One subscription per location
);

-- Create digital_location_payments table
CREATE TABLE digital_location_payments (
    id SERIAL PRIMARY KEY,
    digital_location_id UUID NOT NULL REFERENCES digital_locations(id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    payment_method VARCHAR(50) NOT NULL CHECK (payment_method IN ('alipay', 'amex', 'diners', 'discover', 'elo', 'generic', 'hiper', 'hipercard', 'jcb', 'maestro', 'mastercard', 'mir', 'paypal', 'unionpay', 'visa')),
    transaction_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create digital_game_locations table
CREATE TABLE digital_game_locations (
    id SERIAL PRIMARY KEY,
    user_game_id INTEGER REFERENCES user_games(id) ON DELETE CASCADE,
    digital_location_id UUID REFERENCES digital_locations(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_game_id, digital_location_id)
);

-- Create expenses table
CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    expense_type VARCHAR(50) NOT NULL CHECK (expense_type IN ('subscription', 'purchase', 'dlc', 'in_game')),
    digital_location_id UUID REFERENCES digital_locations(id) ON DELETE SET NULL,
    user_game_id INTEGER REFERENCES user_games(id) ON DELETE SET NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create wishlist table
CREATE TABLE wishlist (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    game_id BIGINT REFERENCES games(id) ON DELETE CASCADE,
    platform_id BIGINT NOT NULL,
    release_date BIGINT,
    is_on_sale BOOLEAN DEFAULT false,
    current_price DECIMAL(10,2),
    sale_price DECIMAL(10,2),
    last_price_check TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, game_id)
);

-- Create foreign key indexes
CREATE INDEX idx_user_games_user_id ON user_games(user_id);
CREATE INDEX idx_user_games_game_id ON user_games(game_id);
CREATE INDEX idx_user_games_platform_id ON user_games(platform_id);
CREATE INDEX idx_physical_locations_user_id ON physical_locations(user_id);
CREATE INDEX idx_sublocations_user_id ON sublocations(user_id);
CREATE INDEX idx_sublocations_physical_location_id ON sublocations(physical_location_id);
CREATE INDEX idx_physical_game_locations_user_game_id ON physical_game_locations(user_game_id);
CREATE INDEX idx_physical_game_locations_sublocation_id ON physical_game_locations(sublocation_id);
CREATE INDEX idx_digital_locations_user_id ON digital_locations(user_id);
CREATE INDEX idx_digital_locations_is_subscription ON digital_locations(is_subscription);
CREATE INDEX idx_digital_location_subscriptions_next_payment ON digital_location_subscriptions(next_payment_date);
CREATE INDEX idx_digital_location_payments_digital_location_id ON digital_location_payments(digital_location_id);
CREATE INDEX idx_digital_location_payments_payment_date ON digital_location_payments(payment_date);
CREATE INDEX idx_digital_game_locations_user_game_id ON digital_game_locations(user_game_id);
CREATE INDEX idx_digital_game_locations_digital_location_id ON digital_game_locations(digital_location_id);
CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_digital_location_id ON expenses(digital_location_id);
CREATE INDEX idx_expenses_user_game_id ON expenses(user_game_id);
CREATE INDEX idx_wishlist_user_id ON wishlist(user_id);
CREATE INDEX idx_wishlist_game_id ON wishlist(game_id);