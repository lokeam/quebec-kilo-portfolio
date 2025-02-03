package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server ServerConfig
	Env    string
	Debug  bool
	IGDB   IGDBConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type IGDBConfig struct {
	ClientID       string
	ClientSecret   string
	AuthURL        string           // Twitch OAuth URL
	BaseURL        string           // IGDB API base URL
	TokenTTL       time.Duration    // Duration until we need a token refresh
}

func Load() (*Config, error) {
	env := os.Getenv(EnvEnvironment)
	if env == "" {
		env = EnvDevelopment
	} else if env != EnvDevelopment && env != EnvTest && env != EnvProduction {
		return nil, fmt.Errorf("invalid environment: must be one of %s, %s or %s",
			EnvDevelopment, EnvTest, EnvProduction)
	}

	// Debug mode handling
	debug := false
	if debugStr := os.Getenv(EnvDebug); debugStr == "true" {
		debug = true
	}

	// Port handling
	portStr := os.Getenv(EnvPort)
	port := DefaultPort
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port number format")
		}
		if port < 1 || port > 65535 {
			return nil, fmt.Errorf("port number must be between 1 and 65535")
		}
	}

	// Host handling
	host := os.Getenv(EnvHost)
	// Check if HOST was explicitly set to empty
	if _, exists := os.LookupEnv(EnvHost); exists && host == "" {
			return nil, fmt.Errorf("host cannot be empty")
	}
	// Use default if HOST not set
	if host == "" {
		host = DefaultHost
	}

	// IGDB Configuration
	igdbConfig := IGDBConfig{
		ClientID:        os.Getenv(EnvIGDBClientID),
		ClientSecret:    os.Getenv(EnvIGDBClientSecret),
		AuthURL:         IGDBAuthURL,
		BaseURL:         IGDBBaseURL,
		TokenTTL:        24 * time.Hour,
	}

	return &Config{
		Server: ServerConfig{
				Port: port,
				Host: host,
		},
		Env: env,
		Debug: debug,
		IGDB: igdbConfig,
	}, nil
}