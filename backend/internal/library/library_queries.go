package library

const (
	CheckIfGameIsInLibraryQuery = `
		SELECT EXISTS(
			SELECT 1 FROM user_games
			WHERE user_id = $1 AND game_id = $2
		)
	`

	GetSingleLibraryGameQuery = `
		SELECT DISTINCT ON (g.id)
			g.id,
			g.name,
			g.cover_url,
			ug.game_type as game_type_display_text,
			LOWER(ug.game_type) as game_type_normalized_text,
			ug.favorite as is_favorite,
			ug.created_at
		FROM games g
		JOIN user_games ug ON g.id = ug.game_id
		WHERE ug.user_id = $1 AND g.id = $2
		ORDER BY g.id, ug.id
	`

	GetLibraryGamesBFFQuery = `
		SELECT DISTINCT ON (g.id)
			g.id,
			g.name,
			g.cover_url,
			ug.game_type as game_type_display_text,
			LOWER(ug.game_type) as game_type_normalized_text,
			ug.favorite as is_favorite,
			ug.created_at
		FROM games g
		JOIN user_games ug ON g.id = ug.game_id
		WHERE ug.user_id = $1
		ORDER BY g.id, ug.id
	`

	GetLibraryLocationsQuery = `
		SELECT
			ug.game_id,
			p.id as platform_id,
			p.name as platform_name,
			p.category,
			ug.created_at,
			pl.id as parent_location_id,
			pl.name as parent_location_name,
			pl.location_type as parent_location_type,
			pl.bg_color as parent_location_bg_color,
			sl.id as sublocation_id,
			sl.name as sublocation_name,
			sl.location_type as sublocation_type
		FROM user_games ug
		JOIN platforms p ON ug.platform_id = p.id
		LEFT JOIN physical_game_locations pgl ON ug.id = pgl.user_game_id
		LEFT JOIN sublocations sl ON pgl.sublocation_id = sl.id
		LEFT JOIN physical_locations pl ON sl.physical_location_id = pl.id
		WHERE ug.user_id = $1
	`

	UpdateLibraryGameQuery = `
		UPDATE games
		SET name = $2,
				cover_url = $3,
				game_type_display_text = $4,
				game_type_normalized_text = $5
		WHERE id = $1
	`

	DeleteLibraryLocationsQuery = `
		DELETE FROM physical_game_locations
		WHERE user_game_id IN (
			SELECT id FROM user_games WHERE game_id = $1
		)
	`

	// Note: Due to ON DELETE CASCADE in our schema, this will automatically
	// remove related records from physical_game_locations and digital_game_locations
	CascadingDeleteLibraryGameQuery = `
		DELETE FROM user_games
		WHERE user_id = $1 AND game_id = $2
	`

	InsertLibraryLocationQuery = `
		INSERT INTO physical_game_locations (user_game_id, sublocation_id)
		SELECT id, $5
		FROM user_games
		WHERE game_id = $1 AND platform_id = $2 AND platform_name = $3 AND game_type = $4
	`

	GetPhysicalLocationsRefactoredQuery = `
		SELECT
				ug.game_id,
				p.id as platform_id,
				p.name as platform_name,
				p.category,
				ug.created_at,
				pl.id as parent_location_id,
				pl.name as parent_location_name,
				pl.location_type as parent_location_type,
				pl.bg_color as parent_location_bg_color,
				sl.id as sublocation_id,
				sl.name as sublocation_name,
				sl.location_type as sublocation_type
		FROM user_games ug
		JOIN platforms p ON ug.platform_id = p.id
		LEFT JOIN physical_game_locations pgl ON ug.id = pgl.user_game_id
		LEFT JOIN sublocations sl ON pgl.sublocation_id = sl.id
		LEFT JOIN physical_locations pl ON sl.physical_location_id = pl.id
		WHERE ug.user_id = $1 AND ug.game_type = 'physical'
		ORDER BY ug.game_id, p.id
	`

	GetDigitalLocationsRefactoredQuery = `
		SELECT
				ug.game_id,
				p.id as platform_id,
				p.name as platform_name,
				p.category,
				ug.created_at,
				dl.id as digital_location_id,
				dl.name as digital_location_name
		FROM user_games ug
		JOIN platforms p ON ug.platform_id = p.id
		LEFT JOIN digital_game_locations dgl ON ug.id = dgl.user_game_id
		LEFT JOIN digital_locations dl ON dgl.digital_location_id = dl.id
		WHERE ug.user_id = $1 AND ug.game_type = 'digital'
		ORDER BY ug.game_id, p.id
	`

	GetGamesMetadataRefactoredQuery = `
		SELECT DISTINCT ON (g.id)
					g.id,
					g.name,
					g.cover_url,
					g.first_release_date,
					g.rating,
					ug.game_type as game_type_display_text,
					LOWER(ug.game_type) as game_type_normalized_text,
					ug.favorite as favorite,
					ug.created_at,
					EXISTS(SELECT 1 FROM wishlist w WHERE w.user_id = $1 AND w.game_id = g.id) as is_in_wishlist,
					COALESCE(ARRAY_AGG(DISTINCT gen.name) FILTER (WHERE gen.name IS NOT NULL), ARRAY[]::text[]) as genre_names
			FROM games g
			JOIN user_games ug ON g.id = ug.game_id
			LEFT JOIN game_genres gg ON g.id = gg.game_id
			LEFT JOIN genres gen ON gg.genre_id = gen.id
			WHERE ug.user_id = $1
			GROUP BY g.id, g.name, g.cover_url, g.first_release_date, g.rating, ug.game_type, ug.favorite, ug.created_at, ug.id
			ORDER BY g.id, ug.id
	`

	// Batch deletion queries
	// Get all versions of a game for a user (used when deleteAll = true)
	GetGameVersionsForBatchDeleteQuery = `
		SELECT
			ug.id as user_game_id,
			ug.game_id,
			ug.platform_id,
			p.name as platform_name,
			ug.game_type as type,
			COALESCE(pgl.sublocation_id, dgl.digital_location_id) as location_id
		FROM user_games ug
		JOIN platforms p ON ug.platform_id = p.id
		LEFT JOIN physical_game_locations pgl ON ug.id = pgl.user_game_id
		LEFT JOIN digital_game_locations dgl ON ug.id = dgl.user_game_id
		WHERE ug.user_id = $1 AND ug.game_id = $2
	`

	// Delete specific user_game entries (this will cascade to location tables)
	DeleteSpecificGameVersionsQuery = `
		DELETE FROM user_games
		WHERE user_id = $1
		AND game_id = $2
		AND platform_id = ANY($3)
		AND (
			(game_type = 'physical' AND id IN (
				SELECT pgl.user_game_id
				FROM physical_game_locations pgl
				WHERE pgl.sublocation_id = ANY($4)
			))
			OR
			(game_type = 'digital' AND id IN (
				SELECT dgl.user_game_id
				FROM digital_game_locations dgl
				WHERE dgl.digital_location_id = ANY($5)
			))
		)
	`

		// Count deleted versions for response
	CountDeletedGameVersionsQuery = `
		SELECT COUNT(*)
		FROM user_games
		WHERE user_id = $1 AND game_id = $2
	`
)