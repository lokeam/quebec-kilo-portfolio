package config

// Environment Variables
const (
	// Server
	EnvEnvironment = "ENV"
	EnvDebug       = "APP_DEBUG"
	EnvPort        = "PORT"
	EnvHost        = "HOST"

	// IGDB
	EnvIGDBClientID       = "IGDB_CLIENT_ID"
	EnvIGDBClientSecret   = "IGDB_CLIENT_SECRET"
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
	DefaultPort = 8080
	DefaultHost = "localhost"
)