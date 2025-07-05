-- =============================
-- QKO SEED DATA
-- =============================
-- This file is (hopefully) structured for easy debugging + maintenance.
-- - Each section is clearly commented
-- - Inserts are in small batches (5-10 records per insert)
-- - The whole file is wrapped in a transaction
-- - Order respects foreign key constraints
-- - If an error occurs, it should be easy to spot and fix

BEGIN;

-- 1. USERS
-- --------
INSERT INTO users (id, email) VALUES
    ('auth0|6866ca863a5f54c1e40be745', 'test@example.com');

-- 2. SPENDING CATEGORIES
-- ----------------------
INSERT INTO spending_categories (name, media_type) VALUES
    ('hardware', 'hardware'),
    ('dlc', 'dlc'),
    ('in_game_purchase', 'in_game_purchase'),
    ('physical_game', 'physical_game'),
    ('digital_game', 'digital_game'),
    ('misc', 'misc');

-- 3. PHYSICAL LOCATIONS (3 records)
-- ---------------------------------
INSERT INTO physical_locations (id, user_id, name, location_type, map_coordinates, bg_color, created_at, updated_at) VALUES
    ('656bb8f2-0b77-444b-be5a-e6d178dad525', 'auth0|6866ca863a5f54c1e40be745', 'John''s House', 'house', '34.390490780304695, 132.50562258911452', 'red', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('8d8ffbf0-1f99-4a90-8d50-db8682982b72', 'auth0|6866ca863a5f54c1e40be745', 'Condo', 'apartment', '34.38147396984524, 132.5150402189977', 'blue', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('6204e0af-91b1-4863-8334-9968732cd0e2', 'auth0|6866ca863a5f54c1e40be745', 'Truck', 'vehicle', NULL, 'green', '2024-01-01 00:00:00', '2024-01-01 00:00:00');

-- 4. SUBLOCATIONS (4 records)
-- ---------------------------
INSERT INTO sublocations (id, user_id, physical_location_id, name, location_type, stored_items, created_at, updated_at) VALUES
    ('9f18041d-1c5b-44f7-bcd5-c0eb91471630', 'auth0|6866ca863a5f54c1e40be745', '656bb8f2-0b77-444b-be5a-e6d178dad525', 'Media Console', 'console', 0, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('2f4fc29e-153c-4ce0-b57c-1905345e7748', 'auth0|6866ca863a5f54c1e40be745', '656bb8f2-0b77-444b-be5a-e6d178dad525', 'Study bookshelf', 'shelf', 0, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('a6f74461-bf70-41ba-a791-885521f4e603', 'auth0|6866ca863a5f54c1e40be745', '8d8ffbf0-1f99-4a90-8d50-db8682982b72', 'Media Console', 'console', 0, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('3eaf36ec-0ac0-4fe9-a9ba-376d6dc67bdd', 'auth0|6866ca863a5f54c1e40be745', '6204e0af-91b1-4863-8334-9968732cd0e2', 'Overlanding Storage Bin', 'box', 0, '2024-01-01 00:00:00', '2024-01-01 00:00:00');

-- 5. DIGITAL LOCATIONS (3 + 5 records)
-- ------------------------------------
-- Non-subscription
INSERT INTO digital_locations (id, user_id, name, is_subscription, is_active, url, payment_method, created_at, updated_at) VALUES
    ('5ab530cb-fdd2-4281-aa90-cdb1abe63d02', 'auth0|6866ca863a5f54c1e40be745', 'Steam', false, true, 'https://store.steampowered.com/', 'paypal', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('f2f71849-f870-4e9e-848e-5be4c09da768', 'auth0|6866ca863a5f54c1e40be745', 'GOG', false, true, 'https://www.gog.com/en/', 'visa', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('c8cac82e-e3df-42db-a22f-19fc861cabf4', 'auth0|6866ca863a5f54c1e40be745', 'Epic Games', false, true, 'https://store.epicgames.com/en-US/', 'amex', '2024-01-01 00:00:00', '2024-01-01 00:00:00');
-- Subscriptions
INSERT INTO digital_locations (id, user_id, name, is_subscription, is_active, url, payment_method, created_at, updated_at) VALUES
    ('97910ffe-4ecd-4bb8-8607-817ab690c331', 'auth0|6866ca863a5f54c1e40be745', 'Apple Arcade', true, true, 'https://www.apple.com/apple-arcade/', 'paypal', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('fbf8cd8f-2f29-4b75-a7cd-22d2658cba4c', 'auth0|6866ca863a5f54c1e40be745', 'Google Play Pass', true, true, 'https://play.google.com/store/pass/getstarted?hl=en', 'visa', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('b7bbc511-1b46-47b5-81d4-a2564bc81700', 'auth0|6866ca863a5f54c1e40be745', 'Nintendo Switch Online', true, true, 'https://www.nintendo.com/', 'paypal', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('bb380774-3ee9-4ab4-a15a-d2b4d1a55e59', 'auth0|6866ca863a5f54c1e40be745', 'Xbox Game Pass', true, true, 'https://www.xbox.com/en-US/xbox-game-pass', 'mastercard', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('bdc08ac6-e6cd-4475-809f-a16b1b13570a', 'auth0|6866ca863a5f54c1e40be745', 'Playstation Plus', true, true, 'https://www.playstation.com/en-us/playstation-network/', 'jcb', '2024-01-01 00:00:00', '2024-01-01 00:00:00');

-- 6. DIGITAL LOCATION SUBSCRIPTIONS (5 records)
-- ---------------------------------------------
INSERT INTO digital_location_subscriptions (digital_location_id, billing_cycle, cost_per_cycle, anchor_date, payment_method, created_at, updated_at) VALUES
    ('97910ffe-4ecd-4bb8-8607-817ab690c331', '1 month', 6.99, '2024-01-09', 'paypal', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('fbf8cd8f-2f29-4b75-a7cd-22d2658cba4c', '1 month', 4.99, '2024-01-12', 'visa', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('b7bbc511-1b46-47b5-81d4-a2564bc81700', '12 month', 49.99, '2024-01-19', 'paypal', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('bb380774-3ee9-4ab4-a15a-d2b4d1a55e59', '12 month', 239.88, '2024-01-20', 'mastercard', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('bdc08ac6-e6cd-4475-809f-a16b1b13570a', '3 month', 49.99, '2024-01-25', 'jcb', '2024-01-01 00:00:00', '2024-01-01 00:00:00');

-- 7. GAMES (batch 1/3)
-- --------------------
INSERT INTO games (id, name, summary, cover_id, cover_url, first_release_date, rating, created_at, updated_at) VALUES
    (119133, 'Elden Ring', 'An action RPG developed by FromSoftware', 12345, '//images.igdb.com/igdb/image/upload/t_thumb/co4jni.jpg', 1645747200, 95.22590412729953, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (325591, 'Elden Ring: Nightreign', 'An action survival game set in the Elden Ring universe', 12346, '/images.igdb.com/igdb/image/upload/t_thumb/co95gk.jpg', 1748563200, 79.19612663287752, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (11133, 'Dark Souls III', 'An action RPG developed by FromSoftware', 12347, '//images.igdb.com/igdb/image/upload/t_thumb/co1vcf.jpg', 1458777600, 87.24220554872059, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (81085, 'Dark Souls: Remastered', 'A remastered version of the original Dark Souls', 12348, '//images.igdb.com/igdb/image/upload/t_thumb/co2uro.jpg', 1527033600, 87.24220554872059, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (17000, 'Stardew Valley', 'A farming simulation RPG', 12349, '//images.igdb.com/igdb/image/upload/t_thumb/xrpmydnu9rpxvxfjkiu7.jpg', 1456444800, 86.67831629089243, '2024-01-01 00:00:00', '2024-01-01 00:00:00');
-- 7. GAMES (batch 2/3)
INSERT INTO games (id, name, summary, cover_id, cover_url, first_release_date, rating, created_at, updated_at) VALUES
    (6710, 'Street Fighter III: 3rd Strike', 'A classic fighting game', 12350, '//images.igdb.com/igdb/image/upload/t_thumb/co6bkh.jpg', 926467200, 87.24220554872059, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (14536, 'Street Fighter Alpha Anthology', 'A collection of Street Fighter Alpha games', 12351, '//images.igdb.com/igdb/image/upload/t_thumb/co6cdy.jpg', 1148515200, 81.10596849215811, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (65017, 'Pocket Planes', 'A mobile flight simulation game', 12352, '//images.igdb.com/igdb/image/upload/t_thumb/t6lxg1120filt021zspc.jpg', 1339632000, 0, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (8965, 'Darkest Dungeon', 'A challenging turn-based RPG', 12353, '//images.igdb.com/igdb/image/upload/t_thumb/co1rc4.jpg', 1453161600, 82.11692565397573, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (1488, 'Gradius V', 'A classic shoot-em-up game', 12354, '//images.igdb.com/igdb/image/upload/t_thumb/co50f9.jpg', 1090454400, 74.82654600301659, '2024-01-01 00:00:00', '2024-01-01 00:00:00');
-- 7. GAMES (batch 3/3)
INSERT INTO games (id, name, summary, cover_id, cover_url, first_release_date, rating, created_at, updated_at) VALUES
    (76, 'Dragon Age: Origins', 'A fantasy RPG developed by BioWare', 12355, '//images.igdb.com/igdb/image/upload/t_thumb/co2mvs.jpg', 1257206400, 86.1523903100883, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (74, 'Mass Effect 2', 'A sci-fi RPG developed by BioWare', 12356, '//images.igdb.com/igdb/image/upload/t_thumb/co20ac.jpg', 1264464000, 91.42965571648986, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (75, 'Mass Effect 3', 'The final chapter of the Mass Effect trilogy', 12357, '//images.igdb.com/igdb/image/upload/t_thumb/co1x7q.jpg', 1330992000, 85.34523963934298, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (115276, 'Super Mario Maker 2', 'A level creation game for Nintendo Switch', 12358, '//images.igdb.com/igdb/image/upload/t_thumb/co21vy.jpg', 1561680000, 81.08166318186613, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (266686, 'F-Zero 99', 'A battle royale racing game', 12359, '//images.igdb.com/igdb/image/upload/t_thumb/co731j.jpg', 1694649600, 70.17683724359068, '2024-01-01 00:00:00', '2024-01-01 00:00:00');
-- 7. GAMES (batch 4/4 - Additional games from legacy)
INSERT INTO games (id, name, summary, cover_id, cover_url, first_release_date, rating, created_at, updated_at) VALUES
    (7346, 'The Legend of Zelda: Breath of the Wild', 'An open-world action-adventure game', 12360, '//images.igdb.com/igdb/image/upload/t_thumb/co3p2d.jpg', 1488499200, 92.73313136791542, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (21162, 'Guwange', 'A classic arcade shoot-em-up', 12361, '//images.igdb.com/igdb/image/upload/t_thumb/co2eyf.jpg', 930182400, 0, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (319924, 'Taiko no Tatsujin: Niji-iro Version', 'A rhythm game', 12362, '//images.igdb.com/igdb/image/upload/t_thumb/co91hx.jpg', 1585008000, 0, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (125642, 'Path of Exile 2', 'An action RPG sequel', 12363, '//images.igdb.com/igdb/image/upload/t_thumb/co8ae0.jpg', 1733443200, 88.31078844648701, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (240009, 'Elden Ring: Shadow of the Erdtree', 'DLC expansion for Elden Ring', 12364, '//images.igdb.com/igdb/image/upload/t_thumb/co7sly.jpg', 1718928000, 82.27669936090378, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (250616, 'Helldivers 2', 'A cooperative shooter game', 12365, '//images.igdb.com/igdb/image/upload/t_thumb/co741o.jpg', 1733443200, 82.27669936090378, '2024-01-01 00:00:00', '2024-01-01 00:00:00');

-- 8. PLATFORMS (batch 1/2)
-- ------------------------
INSERT INTO platforms (id, name, category, model, created_at, updated_at) VALUES
    (6, 'PC (Microsoft Windows)', 'pc', 'Windows PC', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (8, 'PlayStation 2', 'console', 'PS2', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (9, 'PlayStation 3', 'console', 'PS3', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (12, 'Xbox 360', 'console', 'Xbox 360', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (34, 'Mobile', 'mobile', 'Android', '2024-01-01 00:00:00', '2024-01-01 00:00:00');
-- 8. PLATFORMS (batch 2/2)
INSERT INTO platforms (id, name, category, model, created_at, updated_at) VALUES
    (39, 'Mobile', 'mobile', 'iOS', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (48, 'PlayStation 4', 'console', 'PS4', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (49, 'Xbox One', 'console', 'Xbox One', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (52, 'Arcade', 'console', 'Arcade Cabinet', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (130, 'Nintendo Switch', 'console', 'Switch', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    (169, 'Xbox Series X|S', 'console', 'Xbox Series X|S', '2024-01-01 00:00:00', '2024-01-01 00:00:00');

-- 9. USER GAMES (batch 1/6 - Elden Ring complex scenarios)
-- --------------------------------------------------------
INSERT INTO user_games (user_id, game_id, platform_id, game_type, copy_number, is_unique_copy, favorite, created_at) VALUES
    -- Elden Ring (Digital - Steam)
    ('auth0|6866ca863a5f54c1e40be745', 119133, 6, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Elden Ring (Digital - PlayStation Plus)
    ('auth0|6866ca863a5f54c1e40be745', 119133, 48, 'digital', 2, true, false, '2024-01-01 00:00:00'),
    -- Elden Ring (Physical - PS4)
    ('auth0|6866ca863a5f54c1e40be745', 119133, 48, 'physical', 3, true, false, '2024-01-01 00:00:00'),
    -- Elden Ring (Physical - Xbox Series X)
    ('auth0|6866ca863a5f54c1e40be745', 119133, 169, 'physical', 4, true, false, '2024-01-01 00:00:00'),
    -- Elden Ring: Nightreign (Digital - Steam)
    ('auth0|6866ca863a5f54c1e40be745', 325591, 6, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Elden Ring: Nightreign (Physical - Xbox Series X)
    ('auth0|6866ca863a5f54c1e40be745', 325591, 169, 'physical', 2, true, false, '2024-01-01 00:00:00');
-- 9. USER GAMES (batch 2/6 - Dark Souls series)
INSERT INTO user_games (user_id, game_id, platform_id, game_type, copy_number, is_unique_copy, favorite, created_at) VALUES
    -- Dark Souls III (Digital - Steam)
    ('auth0|6866ca863a5f54c1e40be745', 11133, 6, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Dark Souls III (Physical - Xbox One)
    ('auth0|6866ca863a5f54c1e40be745', 11133, 49, 'physical', 2, true, false, '2024-01-01 00:00:00'),
    -- Dark Souls III (Physical - PS4)
    ('auth0|6866ca863a5f54c1e40be745', 11133, 48, 'physical', 3, true, false, '2024-01-01 00:00:00'),
    -- Dark Souls: Remastered (Physical - Xbox One)
    ('auth0|6866ca863a5f54c1e40be745', 81085, 49, 'physical', 1, true, false, '2024-01-01 00:00:00'),
    -- Dark Souls: Remastered (Digital - Nintendo Switch Online)
    ('auth0|6866ca863a5f54c1e40be745', 81085, 130, 'digital', 2, true, false, '2024-01-01 00:00:00');
-- 9. USER GAMES (batch 3/6 - Stardew Valley multi-platform)
INSERT INTO user_games (user_id, game_id, platform_id, game_type, copy_number, is_unique_copy, favorite, created_at) VALUES
    -- Stardew Valley (Digital - Google Play Pass)
    ('auth0|6866ca863a5f54c1e40be745', 17000, 34, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Stardew Valley (Digital - Apple Arcade)
    ('auth0|6866ca863a5f54c1e40be745', 17000, 39, 'digital', 2, true, false, '2024-01-01 00:00:00'),
    -- Stardew Valley (Digital - Nintendo Switch Online)
    ('auth0|6866ca863a5f54c1e40be745', 17000, 130, 'digital', 3, true, false, '2024-01-01 00:00:00'),
    -- Stardew Valley (Physical - Nintendo Switch)
    ('auth0|6866ca863a5f54c1e40be745', 17000, 130, 'physical', 4, true, false, '2024-01-01 00:00:00'),
    -- Stardew Valley (Physical - Xbox One)
    ('auth0|6866ca863a5f54c1e40be745', 17000, 49, 'physical', 5, true, false, '2024-01-01 00:00:00');
-- 9. USER GAMES (batch 4/6 - Fighting games and classics)
INSERT INTO user_games (user_id, game_id, platform_id, game_type, copy_number, is_unique_copy, favorite, created_at) VALUES
    -- Street Fighter III: 3rd Strike (Physical - Arcade)
    ('auth0|6866ca863a5f54c1e40be745', 6710, 52, 'physical', 1, true, false, '2024-01-01 00:00:00'),
    -- Street Fighter Alpha Anthology (Physical - PS2)
    ('auth0|6866ca863a5f54c1e40be745', 14536, 8, 'physical', 1, true, false, '2024-01-01 00:00:00'),
    -- Pocket Planes (Digital - Google Play Pass)
    ('auth0|6866ca863a5f54c1e40be745', 65017, 34, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Pocket Planes (Digital - Apple Arcade)
    ('auth0|6866ca863a5f54c1e40be745', 65017, 39, 'digital', 2, true, false, '2024-01-01 00:00:00'),
    -- Darkest Dungeon (Digital - Steam)
    ('auth0|6866ca863a5f54c1e40be745', 8965, 6, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Darkest Dungeon (Digital - Epic Games)
    ('auth0|6866ca863a5f54c1e40be745', 8965, 6, 'digital', 2, true, false, '2024-01-01 00:00:00'),
    -- Darkest Dungeon (Digital - Apple Arcade)
    ('auth0|6866ca863a5f54c1e40be745', 8965, 39, 'digital', 3, true, false, '2024-01-01 00:00:00'),
    -- Darkest Dungeon (Digital - PlayStation Plus)
    ('auth0|6866ca863a5f54c1e40be745', 8965, 48, 'digital', 4, true, false, '2024-01-01 00:00:00');
-- 9. USER GAMES (batch 5/6 - More classics and RPGs)
INSERT INTO user_games (user_id, game_id, platform_id, game_type, copy_number, is_unique_copy, favorite, created_at) VALUES
    -- Gradius V (Digital - PlayStation Plus)
    ('auth0|6866ca863a5f54c1e40be745', 1488, 9, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Gradius V (Physical - PS2)
    ('auth0|6866ca863a5f54c1e40be745', 1488, 8, 'physical', 2, true, false, '2024-01-01 00:00:00'),
    -- Dragon Age: Origins (Digital - GOG)
    ('auth0|6866ca863a5f54c1e40be745', 76, 6, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Mass Effect 2 (Digital - Steam)
    ('auth0|6866ca863a5f54c1e40be745', 74, 6, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Mass Effect 2 (Digital - Xbox Game Pass)
    ('auth0|6866ca863a5f54c1e40be745', 74, 49, 'digital', 2, true, false, '2024-01-01 00:00:00'),
    -- Mass Effect 2 (Physical - Xbox 360)
    ('auth0|6866ca863a5f54c1e40be745', 74, 12, 'physical', 3, true, false, '2024-01-01 00:00:00'),
    -- Mass Effect 3 (Physical - Xbox 360)
    ('auth0|6866ca863a5f54c1e40be745', 75, 12, 'physical', 1, true, false, '2024-01-01 00:00:00'),
    -- Mass Effect 3 (Digital - Xbox Game Pass)
    ('auth0|6866ca863a5f54c1e40be745', 75, 12, 'digital', 2, true, false, '2024-01-01 00:00:00');
-- 9. USER GAMES (batch 6/6 - Nintendo and additional games)
INSERT INTO user_games (user_id, game_id, platform_id, game_type, copy_number, is_unique_copy, favorite, created_at) VALUES
    -- Super Mario Maker 2 (Digital - Nintendo Switch Online)
    ('auth0|6866ca863a5f54c1e40be745', 115276, 130, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- F-Zero 99 (Digital - Nintendo Switch Online)
    ('auth0|6866ca863a5f54c1e40be745', 266686, 130, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- The Legend of Zelda: Breath of the Wild (Digital - Nintendo Switch Online)
    ('auth0|6866ca863a5f54c1e40be745', 7346, 130, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Guwange (Digital - Xbox Game Pass)
    ('auth0|6866ca863a5f54c1e40be745', 21162, 12, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Guwange (Physical - Arcade)
    ('auth0|6866ca863a5f54c1e40be745', 21162, 52, 'physical', 2, true, false, '2024-01-01 00:00:00'),
    -- Taiko no Tatsujin: Niji-iro Version (Physical - Arcade)
    ('auth0|6866ca863a5f54c1e40be745', 319924, 52, 'physical', 1, true, false, '2024-01-01 00:00:00'),
    -- Path of Exile 2 (Digital - Steam)
    ('auth0|6866ca863a5f54c1e40be745', 125642, 6, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Helldivers 2 (Digital - Steam)
    ('auth0|6866ca863a5f54c1e40be745', 250616, 6, 'digital', 1, true, false, '2024-01-01 00:00:00'),
    -- Elden Ring: Shadow of the Erdtree (Digital - Steam)
    ('auth0|6866ca863a5f54c1e40be745', 240009, 6, 'digital', 1, true, false, '2024-01-01 00:00:00');

-- 10. DIGITAL GAME LOCATIONS (batch 1/4 - Steam games)
-- ----------------------------------------------------
INSERT INTO digital_game_locations (user_game_id, digital_location_id, created_at) VALUES
    -- Elden Ring (Steam)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 119133 AND platform_id = 6 AND game_type = 'digital' AND copy_number = 1), '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', '2024-01-01 00:00:00'),
    -- Elden Ring: Nightreign (Steam)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 325591 AND platform_id = 6 AND game_type = 'digital' AND copy_number = 1), '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', '2024-01-01 00:00:00'),
    -- Dark Souls III (Steam)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 11133 AND platform_id = 6 AND game_type = 'digital' AND copy_number = 1), '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', '2024-01-01 00:00:00'),
    -- Darkest Dungeon (Steam)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 8965 AND platform_id = 6 AND game_type = 'digital' AND copy_number = 1), '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', '2024-01-01 00:00:00'),
    -- Mass Effect 2 (Steam)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 74 AND platform_id = 6 AND game_type = 'digital' AND copy_number = 1), '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', '2024-01-01 00:00:00');
-- 10. DIGITAL GAME LOCATIONS (batch 2/4 - Subscription services)
INSERT INTO digital_game_locations (user_game_id, digital_location_id, created_at) VALUES
    -- Elden Ring (PlayStation Plus)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 119133 AND platform_id = 48 AND game_type = 'digital' AND copy_number = 2), 'bdc08ac6-e6cd-4475-809f-a16b1b13570a', '2024-01-01 00:00:00'),
    -- Dark Souls: Remastered (Nintendo Switch Online)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 81085 AND platform_id = 130 AND game_type = 'digital' AND copy_number = 2), 'b7bbc511-1b46-47b5-81d4-a2564bc81700', '2024-01-01 00:00:00'),
    -- Stardew Valley (Google Play Pass)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 17000 AND platform_id = 34 AND game_type = 'digital' AND copy_number = 1), 'fbf8cd8f-2f29-4b75-a7cd-22d2658cba4c', '2024-01-01 00:00:00'),
    -- Stardew Valley (Apple Arcade)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 17000 AND platform_id = 39 AND game_type = 'digital' AND copy_number = 2), '97910ffe-4ecd-4bb8-8607-817ab690c331', '2024-01-01 00:00:00'),
    -- Stardew Valley (Nintendo Switch Online)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 17000 AND platform_id = 130 AND game_type = 'digital' AND copy_number = 3), 'b7bbc511-1b46-47b5-81d4-a2564bc81700', '2024-01-01 00:00:00');
-- 10. DIGITAL GAME LOCATIONS (batch 3/4 - More subscription services)
INSERT INTO digital_game_locations (user_game_id, digital_location_id, created_at) VALUES
    -- Pocket Planes (Google Play Pass)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 65017 AND platform_id = 34 AND game_type = 'digital' AND copy_number = 1), 'fbf8cd8f-2f29-4b75-a7cd-22d2658cba4c', '2024-01-01 00:00:00'),
    -- Pocket Planes (Apple Arcade)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 65017 AND platform_id = 39 AND game_type = 'digital' AND copy_number = 2), '97910ffe-4ecd-4bb8-8607-817ab690c331', '2024-01-01 00:00:00'),
    -- Darkest Dungeon (Epic Games)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 8965 AND platform_id = 6 AND game_type = 'digital' AND copy_number = 2), 'c8cac82e-e3df-42db-a22f-19fc861cabf4', '2024-01-01 00:00:00'),
    -- Darkest Dungeon (Apple Arcade)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 8965 AND platform_id = 39 AND game_type = 'digital' AND copy_number = 3), '97910ffe-4ecd-4bb8-8607-817ab690c331', '2024-01-01 00:00:00'),
    -- Darkest Dungeon (PlayStation Plus)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 8965 AND platform_id = 48 AND game_type = 'digital' AND copy_number = 4), 'bdc08ac6-e6cd-4475-809f-a16b1b13570a', '2024-01-01 00:00:00');
-- 10. DIGITAL GAME LOCATIONS (batch 4/4 - GOG, Xbox Game Pass, and more)
INSERT INTO digital_game_locations (user_game_id, digital_location_id, created_at) VALUES
    -- Gradius V (PlayStation Plus)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 1488 AND platform_id = 9 AND game_type = 'digital' AND copy_number = 1), 'bdc08ac6-e6cd-4475-809f-a16b1b13570a', '2024-01-01 00:00:00'),
    -- Dragon Age: Origins (GOG)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 76 AND platform_id = 6 AND game_type = 'digital' AND copy_number = 1), 'f2f71849-f870-4e9e-848e-5be4c09da768', '2024-01-01 00:00:00'),
    -- Mass Effect 2 (Xbox Game Pass)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 74 AND platform_id = 49 AND game_type = 'digital' AND copy_number = 2), 'bb380774-3ee9-4ab4-a15a-d2b4d1a55e59', '2024-01-01 00:00:00'),
    -- Mass Effect 3 (Xbox Game Pass)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 75 AND platform_id = 12 AND game_type = 'digital' AND copy_number = 2), 'bb380774-3ee9-4ab4-a15a-d2b4d1a55e59', '2024-01-01 00:00:00'),
    -- Super Mario Maker 2 (Nintendo Switch Online)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 115276 AND platform_id = 130 AND game_type = 'digital' AND copy_number = 1), 'b7bbc511-1b46-47b5-81d4-a2564bc81700', '2024-01-01 00:00:00'),
    -- F-Zero 99 (Nintendo Switch Online)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 266686 AND platform_id = 130 AND game_type = 'digital' AND copy_number = 1), 'b7bbc511-1b46-47b5-81d4-a2564bc81700', '2024-01-01 00:00:00'),
    -- The Legend of Zelda: Breath of the Wild (Nintendo Switch Online)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 7346 AND platform_id = 130 AND game_type = 'digital' AND copy_number = 1), 'b7bbc511-1b46-47b5-81d4-a2564bc81700', '2024-01-01 00:00:00'),
    -- Guwange (Xbox Game Pass)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 21162 AND platform_id = 12 AND game_type = 'digital' AND copy_number = 1), 'bb380774-3ee9-4ab4-a15a-d2b4d1a55e59', '2024-01-01 00:00:00'),
    -- Path of Exile 2 (Steam)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 125642 AND platform_id = 6 AND game_type = 'digital' AND copy_number = 1), '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', '2024-01-01 00:00:00'),
    -- Helldivers 2 (Steam)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 250616 AND platform_id = 6 AND game_type = 'digital' AND copy_number = 1), '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', '2024-01-01 00:00:00'),
    -- Elden Ring: Shadow of the Erdtree (Steam)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 240009 AND platform_id = 6 AND game_type = 'digital' AND copy_number = 1), '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', '2024-01-01 00:00:00');

-- 11. PHYSICAL GAME LOCATIONS (batch 1/3 - House and Condo locations)
-- -------------------------------------------------------------------
INSERT INTO physical_game_locations (user_game_id, sublocation_id, created_at) VALUES
    -- Elden Ring (PS4 - Study Bookshelf)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 119133 AND platform_id = 48 AND game_type = 'physical' AND copy_number = 3), '2f4fc29e-153c-4ce0-b57c-1905345e7748', '2024-01-01 00:00:00'),
    -- Elden Ring (Xbox Series X - Condo Console)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 119133 AND platform_id = 169 AND game_type = 'physical' AND copy_number = 4), 'a6f74461-bf70-41ba-a791-885521f4e603', '2024-01-01 00:00:00'),
    -- Elden Ring: Nightreign (Xbox Series X - Condo Console)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 325591 AND platform_id = 169 AND game_type = 'physical' AND copy_number = 2), 'a6f74461-bf70-41ba-a791-885521f4e603', '2024-01-01 00:00:00'),
    -- Dark Souls: Remastered (Xbox One - Condo Console)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 81085 AND platform_id = 49 AND game_type = 'physical' AND copy_number = 1), 'a6f74461-bf70-41ba-a791-885521f4e603', '2024-01-01 00:00:00'),
    -- Stardew Valley (Nintendo Switch - Study Bookshelf)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 17000 AND platform_id = 130 AND game_type = 'physical' AND copy_number = 4), '2f4fc29e-153c-4ce0-b57c-1905345e7748', '2024-01-01 00:00:00'),
    -- Stardew Valley (Xbox One - Condo Console)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 17000 AND platform_id = 49 AND game_type = 'physical' AND copy_number = 5), 'a6f74461-bf70-41ba-a791-885521f4e603', '2024-01-01 00:00:00');
-- 11. PHYSICAL GAME LOCATIONS (batch 2/3 - Arcade and classic games)
INSERT INTO physical_game_locations (user_game_id, sublocation_id, created_at) VALUES
    -- Street Fighter III: 3rd Strike (Arcade - House Console)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 6710 AND platform_id = 52 AND game_type = 'physical' AND copy_number = 1), '9f18041d-1c5b-44f7-bcd5-c0eb91471630', '2024-01-01 00:00:00'),
    -- Guwange (Arcade - House Console)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 21162 AND platform_id = 52 AND game_type = 'physical' AND copy_number = 2), '9f18041d-1c5b-44f7-bcd5-c0eb91471630', '2024-01-01 00:00:00'),
    -- Taiko no Tatsujin: Niji-iro Version (Arcade - House Console)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 319924 AND platform_id = 52 AND game_type = 'physical' AND copy_number = 1), '9f18041d-1c5b-44f7-bcd5-c0eb91471630', '2024-01-01 00:00:00'),
    -- Street Fighter Alpha Anthology (PS2 - Truck Storage)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 14536 AND platform_id = 8 AND game_type = 'physical' AND copy_number = 1), '3eaf36ec-0ac0-4fe9-a9ba-376d6dc67bdd', '2024-01-01 00:00:00'),
    -- Gradius V (PS2 - Truck Storage)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 1488 AND platform_id = 8 AND game_type = 'physical' AND copy_number = 2), '3eaf36ec-0ac0-4fe9-a9ba-376d6dc67bdd', '2024-01-01 00:00:00');
-- 11. PHYSICAL GAME LOCATIONS (batch 3/3 - Truck storage and Xbox games)
INSERT INTO physical_game_locations (user_game_id, sublocation_id, created_at) VALUES
    -- Dark Souls III (Xbox One - Truck Storage)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 11133 AND platform_id = 49 AND game_type = 'physical' AND copy_number = 2), '3eaf36ec-0ac0-4fe9-a9ba-376d6dc67bdd', '2024-01-01 00:00:00'),
    -- Dark Souls III (PS4 - Truck Storage)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 11133 AND platform_id = 48 AND game_type = 'physical' AND copy_number = 3), '3eaf36ec-0ac0-4fe9-a9ba-376d6dc67bdd', '2024-01-01 00:00:00'),
    -- Mass Effect 2 (Xbox 360 - Condo Console)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 74 AND platform_id = 12 AND game_type = 'physical' AND copy_number = 3), 'a6f74461-bf70-41ba-a791-885521f4e603', '2024-01-01 00:00:00'),
    -- Mass Effect 3 (Xbox 360 - Condo Console)
    ((SELECT id FROM user_games WHERE user_id = 'auth0|6866ca863a5f54c1e40be745' AND game_id = 75 AND platform_id = 12 AND game_type = 'physical' AND copy_number = 1), 'a6f74461-bf70-41ba-a791-885521f4e603', '2024-01-01 00:00:00');

-- 12. ONE TIME PURCHASES (batch 1/4 - 2024 purchases)
-- ----------------------------------------------------
INSERT INTO one_time_purchases (user_id, title, amount, purchase_date, payment_method, spending_category_id, digital_location_id, is_digital, is_wishlisted, created_at, updated_at) VALUES
    -- June 2024
    ('auth0|6866ca863a5f54c1e40be745', 'G.Skill Trident Z5 Neo RGB DDR5-6000', 219.98, '2024-06-16 00:00:00', 'visa', 1, NULL, false, true, '2024-06-16 00:00:00', '2024-06-16 00:00:00'),
    -- July 2024
    ('auth0|6866ca863a5f54c1e40be745', 'ELDEN RING Shadow of the Erdtree', 39.99, '2024-07-16 00:00:00', 'paypal', 2, '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', true, true, '2024-07-16 00:00:00', '2024-07-16 00:00:00'),
    ('auth0|6866ca863a5f54c1e40be745', 'Path of Exile 2', 10.00, '2024-07-16 00:00:00', 'paypal', 3, '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', true, false, '2024-07-16 00:00:00', '2024-07-16 00:00:00'),
    -- August 2024
    ('auth0|6866ca863a5f54c1e40be745', 'Noctua NF-P12 Performance Cooling Fan', 30.90, '2024-08-16 00:00:00', 'visa', 1, NULL, false, false, '2024-08-16 00:00:00', '2024-08-16 00:00:00'),
    -- September 2024
    ('auth0|6866ca863a5f54c1e40be745', 'Dragon Age Origins', 7.99, '2024-09-16 00:00:00', 'paypal', 5, 'f2f71849-f870-4e9e-848e-5be4c09da768', true, true, '2024-09-16 00:00:00', '2024-09-16 00:00:00');
-- 12. ONE TIME PURCHASES (batch 2/4 - Late 2024 purchases)
INSERT INTO one_time_purchases (user_id, title, amount, purchase_date, payment_method, spending_category_id, digital_location_id, is_digital, is_wishlisted, created_at, updated_at) VALUES
    -- November 2024
    ('auth0|6866ca863a5f54c1e40be745', 'Gradius V', 59.99, '2024-11-16 00:00:00', 'paypal', 5, NULL, false, true, '2024-11-16 00:00:00', '2024-11-16 00:00:00'),
    -- December 2024
    ('auth0|6866ca863a5f54c1e40be745', 'Noctua NH-D15 chromax.Black, Dual-Tower CPU Cooler', 139.95, '2024-12-16 00:00:00', 'visa', 1, NULL, false, false, '2024-12-16 00:00:00', '2024-12-16 00:00:00'),
    -- January 2025
    ('auth0|6866ca863a5f54c1e40be745', 'Helldivers 2', 10.00, '2025-01-16 00:00:00', 'paypal', 3, '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', true, false, '2025-01-16 00:00:00', '2025-01-16 00:00:00'),
    -- March 2025
    ('auth0|6866ca863a5f54c1e40be745', 'Path of Exile', 10.00, '2025-03-16 00:00:00', 'paypal', 3, '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', true, false, '2025-03-16 00:00:00', '2025-03-16 00:00:00'),
    ('auth0|6866ca863a5f54c1e40be745', 'Path of Exile', 20.00, '2025-03-16 00:00:00', 'paypal', 3, '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', true, false, '2025-03-16 00:00:00', '2025-03-16 00:00:00');
-- 12. ONE TIME PURCHASES (batch 3/4 - 2025 purchases)
INSERT INTO one_time_purchases (user_id, title, amount, purchase_date, payment_method, spending_category_id, digital_location_id, is_digital, is_wishlisted, created_at, updated_at) VALUES
    -- April 2025
    ('auth0|6866ca863a5f54c1e40be745', 'Dark Souls: Remastered', 59.99, '2025-04-16 00:00:00', 'visa', 5, NULL, false, true, '2025-04-16 00:00:00', '2025-04-16 00:00:00'),
    -- June 2025
    ('auth0|6866ca863a5f54c1e40be745', 'Path of Exile', 10.00, '2025-06-16 00:00:00', 'paypal', 3, '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', true, false, '2025-06-16 00:00:00', '2025-06-16 00:00:00'),
    ('auth0|6866ca863a5f54c1e40be745', 'Elden Ring: Nightreign', 49.00, '2025-06-16 00:00:00', 'paypal', 5, '5ab530cb-fdd2-4281-aa90-cdb1abe63d02', true, false, '2025-06-16 00:00:00', '2025-06-16 00:00:00');

-- 13. WISHLIST (5 records)
-- ------------------------
INSERT INTO wishlist (user_id, game_id, platform_id, release_date, is_on_sale, current_price, sale_price, last_price_check, created_at, updated_at) VALUES
    ('auth0|6866ca863a5f54c1e40be745', 119133, 6, 1645747200, false, 59.99, NULL, '2024-01-01 00:00:00', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('auth0|6866ca863a5f54c1e40be745', 325591, 6, 1748563200, true, 39.99, 29.99, '2024-01-01 00:00:00', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('auth0|6866ca863a5f54c1e40be745', 11133, 6, 1458777600, false, 29.99, NULL, '2024-01-01 00:00:00', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('auth0|6866ca863a5f54c1e40be745', 81085, 6, 1527033600, true, 19.99, 9.99, '2024-01-01 00:00:00', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
    ('auth0|6866ca863a5f54c1e40be745', 17000, 6, 1456444800, false, 14.99, NULL, '2024-01-01 00:00:00', '2024-01-01 00:00:00', '2024-01-01 00:00:00');

-- 14. MONTHLY SPENDING AGGREGATES (10 records - 2024-2025)
-- ---------------------------------------------------------
INSERT INTO monthly_spending_aggregates (user_id, year, month, total_amount, subscription_amount, one_time_amount, category_amounts, created_at, updated_at) VALUES
    -- June 2024
    ('auth0|6866ca863a5f54c1e40be745', 2024, 6, 219.98, 0, 219.98, '{
        "hardware": 219.98,
        "dlc": 0,
        "in_game_purchase": 0,
        "subscription": 0,
        "physical_game": 0,
        "digital_game": 0,
        "misc": 0
    }'::jsonb, '2024-06-16 00:00:00', '2024-06-16 00:00:00'),
    -- July 2024
    ('auth0|6866ca863a5f54c1e40be745', 2024, 7, 49.99, 0, 49.99, '{
        "hardware": 0,
        "dlc": 39.99,
        "in_game_purchase": 10.00,
        "subscription": 0,
        "physical_game": 0,
        "digital_game": 0,
        "misc": 0
    }'::jsonb, '2024-07-16 00:00:00', '2024-07-16 00:00:00'),
    -- August 2024
    ('auth0|6866ca863a5f54c1e40be745', 2024, 8, 30.90, 0, 30.90, '{
        "hardware": 30.90,
        "dlc": 0,
        "in_game_purchase": 0,
        "subscription": 0,
        "physical_game": 0,
        "digital_game": 0,
        "misc": 0
    }'::jsonb, '2024-08-16 00:00:00', '2024-08-16 00:00:00'),
    -- September 2024
    ('auth0|6866ca863a5f54c1e40be745', 2024, 9, 7.99, 0, 7.99, '{
        "hardware": 0,
        "dlc": 0,
        "in_game_purchase": 0,
        "subscription": 0,
        "physical_game": 7.99,
        "digital_game": 0,
        "misc": 0
    }'::jsonb, '2024-09-16 00:00:00', '2024-09-16 00:00:00'),
    -- November 2024
    ('auth0|6866ca863a5f54c1e40be745', 2024, 11, 59.99, 0, 59.99, '{
        "hardware": 0,
        "dlc": 0,
        "in_game_purchase": 0,
        "subscription": 0,
        "physical_game": 59.99,
        "digital_game": 0,
        "misc": 0
    }'::jsonb, '2024-11-16 00:00:00', '2024-11-16 00:00:00'),
    -- December 2024
    ('auth0|6866ca863a5f54c1e40be745', 2024, 12, 139.95, 0, 139.95, '{
        "hardware": 139.95,
        "dlc": 0,
        "in_game_purchase": 0,
        "subscription": 0,
        "physical_game": 0,
        "digital_game": 0,
        "misc": 0
    }'::jsonb, '2024-12-16 00:00:00', '2024-12-16 00:00:00'),
    -- January 2025
    ('auth0|6866ca863a5f54c1e40be745', 2025, 1, 10.00, 0, 10.00, '{
        "hardware": 0,
        "dlc": 0,
        "in_game_purchase": 10.00,
        "subscription": 0,
        "physical_game": 0,
        "digital_game": 0,
        "misc": 0
    }'::jsonb, '2025-01-16 00:00:00', '2025-01-16 00:00:00'),
    -- March 2025
    ('auth0|6866ca863a5f54c1e40be745', 2025, 3, 30.00, 0, 30.00, '{
        "hardware": 0,
        "dlc": 0,
        "in_game_purchase": 30.00,
        "subscription": 0,
        "physical_game": 0,
        "digital_game": 0,
        "misc": 0
    }'::jsonb, '2025-03-16 00:00:00', '2025-03-16 00:00:00'),
    -- April 2025
    ('auth0|6866ca863a5f54c1e40be745', 2025, 4, 59.99, 0, 59.99, '{
        "hardware": 0,
        "dlc": 0,
        "in_game_purchase": 0,
        "subscription": 0,
        "physical_game": 59.99,
        "digital_game": 0,
        "misc": 0
    }'::jsonb, '2025-04-16 00:00:00', '2025-04-16 00:00:00'),
    -- June 2025
    ('auth0|6866ca863a5f54c1e40be745', 2025, 6, 59.00, 0, 59.00, '{
        "hardware": 0,
        "dlc": 0,
        "in_game_purchase": 10.00,
        "subscription": 0,
        "physical_game": 0,
        "digital_game": 49.00,
        "misc": 0
    }'::jsonb, '2025-06-16 00:00:00', '2025-06-16 00:00:00');

-- 15. YEARLY SPENDING AGGREGATES (2 records - 2024-2025)
-- ------------------------------------------------------
INSERT INTO yearly_spending_aggregates (user_id, year, total_amount, subscription_amount, one_time_amount, created_at, updated_at) VALUES
    ('auth0|6866ca863a5f54c1e40be745', 2024, 508.80, 0, 508.80, '2024-12-31 00:00:00', '2024-12-31 00:00:00'),
    ('auth0|6866ca863a5f54c1e40be745', 2025, 159.99, 0, 159.99, '2025-06-30 00:00:00', '2025-06-30 00:00:00');

COMMIT;