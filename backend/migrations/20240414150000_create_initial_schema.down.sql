-- Drop indexes first
DROP INDEX IF EXISTS idx_wishlist_game_id;
DROP INDEX IF EXISTS idx_wishlist_user_id;
DROP INDEX IF EXISTS idx_expenses_user_game_id;
DROP INDEX IF EXISTS idx_expenses_digital_location_id;
DROP INDEX IF EXISTS idx_expenses_user_id;
DROP INDEX IF EXISTS idx_digital_game_locations_digital_location_id;
DROP INDEX IF EXISTS idx_digital_game_locations_user_game_id;
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
DROP TABLE IF EXISTS digital_game_locations;
DROP TABLE IF EXISTS digital_locations;
DROP TABLE IF EXISTS physical_game_locations;
DROP TABLE IF EXISTS sublocations;
DROP TABLE IF EXISTS physical_locations;
DROP TABLE IF EXISTS user_games;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS schema_migrations;

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";