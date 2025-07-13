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

// Auth0AppMetadataPatchRequest represents the PATCH request body to update app_metadata
// https://auth0.com/docs/api/management/v2#!/Users/patch_users_by_id
// Example:
// {
//   "app_metadata": {
//     "hasCompletedOnboarding": true,
//     "wantsIntroToasts": true
//   }
// }
type Auth0AppMetadataPatchRequest struct {
	AppMetadata map[string]any `json:"app_metadata"`
}

// Auth0AppMetadataPatchResponse represents the response from Auth0 after PATCHing app_metadata
// Only relevant fields included
// https://auth0.com/docs/api/management/v2#!/Users/patch_users_by_id
// Example:
// {
//   "user_id": "auth0|...",
//   "app_metadata": { ... },
//   ...
// }
type Auth0AppMetadataPatchResponse struct {
	UserID      string           `json:"user_id"`
	AppMetadata map[string]any   `json:"app_metadata"`
}