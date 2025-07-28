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
	HealthStatus = "health_status"
)

type Auth0Config struct {
	Domain              string
	ClientID            string
	ClientSecret        string
	Audience            string
	ManagementAudience  string
}

type Config struct {
	Server ServerConfig
	Env    string
	Debug  bool
	CORS   CORSConfig
	IGDB   *IGDBConfig
	Redis  RedisConfig
	Postgres *PostgresConfig
	Email  *EmailConfig
	HealthStatus string
	Auth0  Auth0Config
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

type PostgresConfig struct {
	Host             string
	Port             int
	User             string
	Password         string
	Database         string
	SSLMode          string
	ConnectionString string
	MaxConnections   int
	MaxIdleTime      time.Duration
	MaxLifetime      time.Duration
}

type EmailConfig struct {
	ResendAPIKey string
	FromAddress  string
	FromName     string
	TemplateDir  string
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

	healthStatus := os.Getenv(HealthStatus)
	if healthStatus == "" {
		healthStatus = "available"
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

	// PostgresSQL Config
	postgresConfig := &PostgresConfig{
		Host:               getEnvOrDefault("POSTGRES_HOST", "localhost"),
		Port:               getEnvIntOrDefault("POSTGRES_PORT", 5432),
		User:               getEnvOrDefault("POSTGRES_USER", ""),
		Password:           getEnvOrDefault("POSTGRES_PASSWORD", ""),
		Database:           getEnvOrDefault("POSTGRES_DB", ""),
		SSLMode:            getEnvOrDefault("POSTGRES_SSLMODE", "disable"),
		ConnectionString:   os.Getenv("DATABASE_URL"),
		MaxConnections:     getEnvIntOrDefault("DATABASE_MAX_CONNECTIONS", 500),
		MaxIdleTime:        time.Duration(getEnvIntOrDefault("DATABASE_MAX_IDLE_TIME", 60)) * time.Second,
		MaxLifetime:        time.Duration(getEnvIntOrDefault("DATABASE_MAX_LIFETIME", 30)) * time.Minute,
	}

	// Validate required PostgresSQL Settings
	if postgresConfig.ConnectionString == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	// Email Configuration
	emailConfig := &EmailConfig{
		ResendAPIKey: os.Getenv(EnvResendAPIKey),
		FromAddress:  getEnvOrDefault(EnvEmailFromAddress, "noreply@resend.dev"),
		FromName:     getEnvOrDefault(EnvEmailFromName, "Quebec Kilo"),
		TemplateDir:  "internal/email/templates",
	}

	auth0Config := Auth0Config{
		Domain:              os.Getenv("AUTH0_DOMAIN"),
		ClientID:            os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret:        os.Getenv("AUTH0_CLIENT_SECRET"),
		Audience:            os.Getenv("AUTH0_AUDIENCE"),
		ManagementAudience:  os.Getenv("AUTH0_MANAGEMENT_AUDIENCE"),
	}

	return &Config{
		Server: ServerConfig{
				Port: port,
				Host: host,
		},
		Env:          env,
		Debug:        debug,
		CORS:         corsConfig,
		IGDB:         &igdbConfig,
		Redis:        redisConfig,
		Postgres:     postgresConfig,
		Email:        emailConfig,
		HealthStatus: healthStatus,
		Auth0:        auth0Config,
	}, nil
}

func (cfg *IGDBConfig) GetAccessTokenKey() (string, error) {
	if cfg.AccessTokenKey == "" {
		return "", errors.New("IGDB access token key is missing")
	}
	return cfg.AccessTokenKey, nil
}

// GetConnectionString builds a connection string from components or returns the existing one
func (pc *PostgresConfig) GetConnectionString() string {
	if pc.ConnectionString != "" {
			return pc.ConnectionString
	}

	// Build connection string from individual components
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			pc.User, pc.Password, pc.Host, pc.Port, pc.Database, pc.SSLMode)
}
