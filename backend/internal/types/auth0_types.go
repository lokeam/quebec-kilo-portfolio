package types

// Auth0TokenResponse represents the response from the Auth0 token endpoint
// when requesting a Management API token
// https://auth0.com/docs/api/management/v2/tokens
// Example response:
// {
//   "access_token": "...",
//   "token_type": "Bearer",
//   "expires_in": 86400
// }
type Auth0TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Auth0UserMetadataPatchRequest represents the PATCH request body to update user_metadata
// https://auth0.com/docs/api/management/v2#!/Users/patch_users_by_id
// Example:
// {
//   "user_metadata": {
//     "firstName": "John",
//     "lastName": "Doe",
//     "hasCompletedOnboarding": true
//   }
// }
type Auth0UserMetadataPatchRequest struct {
	UserMetadata map[string]interface{} `json:"user_metadata"`
}

// Auth0UserMetadataPatchResponse represents the response from Auth0 after PATCHing user_metadata
// Only relevant fields included
// https://auth0.com/docs/api/management/v2#!/Users/patch_users_by_id
// Example:
// {
//   "user_id": "auth0|...",
//   "user_metadata": { ... },
//   ...
// }
type Auth0UserMetadataPatchResponse struct {
	UserID      string                 `json:"user_id"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
}