package digital

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lokeam/qko-beta/internal/models"
)

// AddGameToDigitalLocation adds a game to a digital location
func (da *DigitalDbAdapter) AddGameToDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error {
	da.logger.Debug("AddGameToDigitalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
		"gameID":     gameID,
	})

	// First, get the user_game_id for this user and game
	var userGameID int
	err := da.db.GetContext(
		ctx,
		&userGameID,
		GetUserIDGameQuery,
		userID,
		gameID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("game not found in user's library")
		}
		return fmt.Errorf("error getting user game: %w", err)
	}

	// Then, add the game to the digital location
	_, err = da.db.ExecContext(
		ctx,
		AddGameToDigitalLocationQuery,
		userGameID,
		locationID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "digital_game_locations_user_game_id_digital_location_id_key") {
			return fmt.Errorf("game already exists in this digital location")
		}
		return fmt.Errorf("error adding game to digital location: %w", err)
	}

	return nil
}


// RemoveGameFromDigitalLocation removes a game from a digital location
func (da *DigitalDbAdapter) RemoveGameFromDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
	gameID int64,
) error {
	da.logger.Debug("RemoveGameFromDigitalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
		"gameID":     gameID,
	})

	// Get the user_game_id for this user and game
	var userGameID int
	err := da.db.GetContext(
		ctx,
		&userGameID,
		GetUserIDGameQuery,
		userID,
		gameID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("game not found in user's library")
		}
		return fmt.Errorf("error getting user game: %w", err)
	}

	// Then, remove the game from the digital location
	result, err := da.db.ExecContext(
		ctx,
		RemoveGameFromDigitalLocationQuery,
		userGameID,
		locationID,
	)
	if err != nil {
		return fmt.Errorf("error removing game from digital location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("game not found in digital location")
	}

	return nil
}


// GetGamesByDigitalLocationID gets all games in a digital location
func (da *DigitalDbAdapter) GetGamesByDigitalLocationID(
	ctx context.Context,
	userID string,
	locationID string,
) ([]models.Game, error) {
	da.logger.Debug("GetGamesByDigitalLocationID called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
	})

	var games []models.Game
	err := da.db.SelectContext(
		ctx,
		&games,
		GetAllGamesInDigitalLocationQuery,
		locationID,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting games for digital location: %w", err)
	}

	return games, nil
}
