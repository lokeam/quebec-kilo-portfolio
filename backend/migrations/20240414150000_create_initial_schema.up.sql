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

-- Create user_games table
CREATE TABLE user_games (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    game_id BIGINT REFERENCES games(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, game_id)
);

-- Create physical_locations table
CREATE TABLE physical_locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    label VARCHAR(255),
    location_type VARCHAR(50) NOT NULL CHECK (location_type IN ('house', 'apartment', 'office', 'warehouse', 'vehicle')),
    map_coordinates VARCHAR(255),
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
    bg_color VARCHAR(50) NOT NULL CHECK (bg_color IN ('red', 'green', 'blue', 'orange', 'gold', 'purple', 'brown', 'gray')),
    capacity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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
CREATE TABLE IF NOT EXISTS digital_locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    url TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name)
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
CREATE INDEX idx_physical_locations_user_id ON physical_locations(user_id);
CREATE INDEX idx_sublocations_user_id ON sublocations(user_id);
CREATE INDEX idx_sublocations_physical_location_id ON sublocations(physical_location_id);
CREATE INDEX idx_physical_game_locations_user_game_id ON physical_game_locations(user_game_id);
CREATE INDEX idx_physical_game_locations_sublocation_id ON physical_game_locations(sublocation_id);
CREATE INDEX idx_digital_locations_user_id ON digital_locations(user_id);
CREATE INDEX idx_digital_game_locations_user_game_id ON digital_game_locations(user_game_id);
CREATE INDEX idx_digital_game_locations_digital_location_id ON digital_game_locations(digital_location_id);
CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_digital_location_id ON expenses(digital_location_id);
CREATE INDEX idx_expenses_user_game_id ON expenses(user_game_id);
CREATE INDEX idx_wishlist_user_id ON wishlist(user_id);
CREATE INDEX idx_wishlist_game_id ON wishlist(game_id);