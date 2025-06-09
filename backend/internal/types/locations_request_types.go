package types

// Physical Location Request Types
// TODO: Add request types for physical locations

// Sublocation Request Types
type CreateSublocationRequest struct {
	Name               string `json:"name"`
	LocationType       string `json:"location_type"`
	PhysicalLocationID string `json:"physical_location_id"`
}

type UpdateSublocationRequest struct {
	Name         string `json:"name,omitempty"`
	LocationType string `json:"location_type,omitempty"`
}

// Game Management Request Types
type MoveGameRequest struct {
	// The ID of the user_game record (from user_games table)
	UserGameID string `json:"user_game_id"`
	// The target sublocation to move the game to
	TargetSublocationID string `json:"target_sublocation_id"`
}

type RemoveGameRequest struct {
	// The ID of the user_game record (from user_games table)
	UserGameID string `json:"user_game_id"`
}

// DeleteSublocationRequest represents a request to delete one or more sublocations
type DeleteSublocationRequest struct {
	// IDs is a comma-separated list of sublocation IDs to delete
	IDs string `json:"ids"`
}

