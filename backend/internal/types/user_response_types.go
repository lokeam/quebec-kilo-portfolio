package types

import "time"

type UserProfileResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ProfileCompletionResponse struct {
	HasCompleteProfile bool `json:"has_complete_profile"`
}