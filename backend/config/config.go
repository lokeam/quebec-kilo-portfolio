package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	TwitchAccessTokenKey = "twitch:access_token"
)

type Config struct {
	Server ServerConfig
	Env    string
	Debug  bool
	CORS   CORSConfig
	IGDB   *IGDBConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type CORSConfig struct {
	AllowedOrigins       []string
	AllowedMethods       []string
	AllowedHeaders       []string
	ExposedHeaders       []string
	AllowCredentials     bool
	MaxAge               int
}

type IGDBConfig struct {
	ClientID       string
	ClientSecret   string
	AuthURL        string           // Twitch OAuth URL
	BaseURL        string           // IGDB API base URL
	TokenTTL       time.Duration    // Duration until we need a token refresh
	AccessTokenKey string           // Key for storing access token in Redis + Memcache
}

type RedisConfig struct {
	RedisTimeout time.Duration
	RedisTTL     time.Duration
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

	// CORS Configuration + defaults
	corsConfig := CORSConfig{
		AllowedOrigins:   strings.Split(getEnvOrDefault("CORS_ALLOWED_ORIGINS", "https://*,http://*"), ","),
		AllowedMethods:   strings.Split(getEnvOrDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
		AllowedHeaders:   strings.Split(getEnvOrDefault("CORS_ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-CSRF-Token"), ","),
		ExposedHeaders:   strings.Split(getEnvOrDefault("CORS_EXPOSED_HEADERS", "Link"), ","),
		AllowCredentials: getEnvBoolOrDefault("CORS_ALLOW_CREDENTIALS", true),
		MaxAge:          getEnvIntOrDefault("CORS_MAX_AGE", 300),
	}

	// IGDB Configuration
	igdbConfig := IGDBConfig{
		ClientID:        os.Getenv(EnvIGDBClientID),
		ClientSecret:    os.Getenv(EnvIGDBClientSecret),
		AuthURL:         IGDBAuthURL,
		BaseURL:         IGDBBaseURL,
		TokenTTL:        24 * time.Hour,
		AccessTokenKey:  TwitchAccessTokenKey,
	}
	redisConfig := RedisConfig{
		RedisTimeout: 5 * time.Second,
		RedisTTL:     5 * time.Minute,
	}

	return &Config{
		Server: ServerConfig{
				Port: port,
				Host: host,
		},
		Env: env,
		Debug: debug,
		CORS:  corsConfig,
		IGDB:  &igdbConfig,
		Redis: redisConfig,
	}, nil
}

func (cfg *IGDBConfig) GetAccessTokenKey() (string, error) {
	if cfg.AccessTokenKey == "" {
		return "", errors.New("IGDB access token key is missing")
	}
	return cfg.AccessTokenKey, nil
}
