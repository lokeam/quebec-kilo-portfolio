package models

import "time"

type PhysicalLocation struct {
	ID              string  `json:"id" db:"id"`
	UserID          string  `json:"user_id" db:"user_id"`
	Name            string  `json:"name" db:"name"`
	Label           string  `json:"label" db:"label"`
	LocationType    string  `json:"location_type" db:"location_type"`
	MapCoordinates  string  `json:"map_coordinates" db:"map_coordinates"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
  UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
