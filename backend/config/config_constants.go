package config

// Environment Variables
const (
	// Server
	EnvEnvironment = "API_ENV"
	EnvDebug       = "APP_DEBUG"
	EnvPort        = "PORT"
	EnvHost        = "HOST"

	// IGDB
	EnvIGDBClientID       = "IGDB_CLIENT_ID"
	EnvIGDBClientSecret   = "IGDB_CLIENT_SECRET"

	// Auth0
	EnvAuth0Domain = "AUTH0_DOMAIN"
	EnvAuth0Audience = "AUTH0_AUDIENCE"
)

// IGDB API endpoints
const (
	IGDBAuthURL = "https://id.twitch.tv/oauth2/token"
	IGDBBaseURL = "https://api.igdb.com/v4"
)

// Environment values
const (
	EnvDevelopment = "dev"
	EnvTest        = "test"
	EnvProduction  = "prod"
)

// Defaults
const (
	DefaultPort = 8000
	DefaultHost = "localhost"
)