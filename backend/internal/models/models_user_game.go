package models

import (
	"time"
)

type UserGame struct {
	ID      int       `json:"id" db:"id"`
	UserID  string    `json:"user_id" db:"user_id"`
	GameID  int64     `json:"game_id" db:"game_id"`
	AddedAt time.Time `json:"added_at" db:"added_at"`

	// NOTE: Embedded game details (for queries that join tables)
	Game    *Game     `json:"game,omitempty" db:"-"`
}
