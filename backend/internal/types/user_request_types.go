package types

type UpdateUserProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateUserRequest struct {
	Auth0UserID string `json:"auth0_userID"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}