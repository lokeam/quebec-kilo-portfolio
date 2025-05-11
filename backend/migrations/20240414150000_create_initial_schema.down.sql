-- Drop indexes first
DROP INDEX IF EXISTS idx_wishlist_game_id;
DROP INDEX IF EXISTS idx_wishlist_user_id;
DROP INDEX IF EXISTS idx_expenses_user_game_id;
DROP INDEX IF EXISTS idx_expenses_digital_location_id;
DROP INDEX IF EXISTS idx_expenses_user_id;
DROP INDEX IF EXISTS idx_digital_game_locations_digital_location_id;
DROP INDEX IF EXISTS idx_digital_game_locations_user_game_id;
DROP INDEX IF EXISTS idx_digital_location_payments_payment_date;
DROP INDEX IF EXISTS idx_digital_location_payments_digital_location_id;
DROP INDEX IF EXISTS idx_digital_location_subscriptions_next_payment;
DROP INDEX IF EXISTS idx_digital_locations_service_type;
DROP INDEX IF EXISTS idx_digital_locations_user_id;
DROP INDEX IF EXISTS idx_physical_game_locations_sublocation_id;
DROP INDEX IF EXISTS idx_physical_game_locations_user_game_id;
DROP INDEX IF EXISTS idx_sublocations_physical_location_id;
DROP INDEX IF EXISTS idx_sublocations_user_id;
DROP INDEX IF EXISTS idx_physical_locations_user_id;
DROP INDEX IF EXISTS idx_user_games_game_id;
DROP INDEX IF EXISTS idx_user_games_user_id;

-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS wishlist;
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS digital_location_payments;
DROP TABLE IF EXISTS digital_location_subscriptions;
DROP TABLE IF EXISTS digital_game_locations;
DROP TABLE IF EXISTS digital_locations;
DROP TABLE IF EXISTS physical_game_locations;
DROP TABLE IF EXISTS sublocations;
DROP TABLE IF EXISTS physical_locations;
DROP TABLE IF EXISTS user_games;
DROP TABLE IF EXISTS game_themes;
DROP TABLE IF EXISTS game_genres;
DROP TABLE IF EXISTS game_platforms;
DROP TABLE IF EXISTS themes;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS platforms;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS schema_migrations;

-- Drop types
DROP TYPE IF EXISTS payment_method_type;
DROP TYPE IF EXISTS digital_service_type;

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";