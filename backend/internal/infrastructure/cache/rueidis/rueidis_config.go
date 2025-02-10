package cache

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/redis/rueidis"
)

const (
	MIN_REDIS_PORT = 1
	MAX_REDIS_PORT = 65535
	MAX_REDIS_DB = 15

	REDIS_HOST_ENV_VAR = "REDIS_HOST"
	REDIS_PORT_ENV_VAR = "REDIS_PORT"
	REDIS_PASSWORD_ENV_VAR = "REDIS_PASSWORD"
	REDIS_DB_ENV_VAR = "REDIS_DB"
	REDIS_CONN_WRITE_TIMEOUT_ENV_VAR = "REDIS_CONN_WRITE_TIMEOUT"
	REDIS_CONN_READ_TIMEOUT_ENV_VAR = "REDIS_CONN_READ_TIMEOUT"
	REDIS_TIMEOUT_READ_ENV_VAR = "REDIS_TIMEOUT_READ"
	REDIS_TIMEOUT_WRITE_ENV_VAR = "REDIS_TIMEOUT_WRITE"
	REDIS_TIMEOUT_DEFAULT_ENV_VAR = "REDIS_TIMEOUT_DEFAULT"
	RUEIDIS_CLIENT_OPTION_USERNAME = "default"
)

// Structs / Interfaces
type RueidisConfig struct {
	// Connection settings (required, loaded from environment)

	// Redis server hostname
	Host      string

	// Redis server port
	Port      int

	// Redis server pwd (optional, loaded from environment)
	Password  string

	// Redis db number (optional, loaded from environment)
	DB        int

	// Timeout settings for Redis operations

	// ConnectionWriteTimeout
	ConnWriteTimeout     time.Duration

	// ConnectionReadTimeout
	ConnReadTimeout      time.Duration

	// Size of connection pool for blocking operations
	BlockingOpsPoolSize  int

	// Internal Operation-specific timeouts
	TimeoutConfig        TimeoutConfig

	// Migration specific timeouts
	EnableMetrics         bool
	EnableTracing         bool

	// Aplication-specific timeouts
	CacheConfig           CacheConfig
}

type CacheConfig struct {
	// User-related cache settings
	UserData     time.Duration

	// Default time to live for cache entries
	DefaultTTL   time.Duration
}

type TimeoutConfig struct {
	Read     time.Duration  // 100ms for Redis read ops
	Write    time.Duration  // 100ms for Redis write ops
	Default  time.Duration  // 150ms for general internal ops
}

// Constructor
func NewRueidisConfig() *RueidisConfig {
	return &RueidisConfig{

		// Set lower connection timeouts for fast in-memory ops
		ConnWriteTimeout:     200 * time.Millisecond,
		ConnReadTimeout:      200 * time.Millisecond,
		BlockingOpsPoolSize:  10,

		TimeoutConfig: TimeoutConfig{
			Read:    100 * time.Millisecond,
			Write:   100 * time.Millisecond,
			Default: 150 * time.Millisecond,
		},

		EnableMetrics: true,
		EnableTracing: true,

		CacheConfig: CacheConfig{
			// NOTE: Adjust these values based on UAT data update findings
			UserData:     30 * time.Minute,
			DefaultTTL:   15 * time.Minute,
		},
	}
}


// Methods
func (c *RueidisConfig) LoadFromEnv() error {
	// Required: Host
	host := os.Getenv(REDIS_HOST_ENV_VAR)
	if host == "" {
		return fmt.Errorf("REDIS_HOST environment variable is required")
	}
	c.Host = host

	// Required: Port
	portStr := os.Getenv(REDIS_PORT_ENV_VAR)
	if portStr == "" {
		return fmt.Errorf("REDIS_PORT environment variable is required")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid REDIS_PORT value: %w", err)
	}
	if port < MIN_REDIS_PORT || port > MAX_REDIS_PORT {
		return fmt.Errorf("REDIS_PORT must be between 1 and 65535")
	}
	c.Port = port

	// Optional: Password
	if pass := os.Getenv(REDIS_PASSWORD_ENV_VAR); pass != "" {
		c.Password = pass
	}

	// NOTE: Save this debug for later when we re-implement Redis ACL
	// fmt.Printf("DEBUG: REDIS_PASSWORD=%q\n", c.Password)

	// Optional: Database number
	if dbStr := os.Getenv(REDIS_DB_ENV_VAR); dbStr != "" {
		db, err := strconv.Atoi(dbStr)
		if err != nil {
			return fmt.Errorf("invalid REDIS_DB value: %w", err)
		}
		if db < 0 || db > MAX_REDIS_DB{
			return fmt.Errorf("REDIS_DB must be between 0 and %d", MAX_REDIS_DB)
		}
		c.DB = db
	}

	// Load connection timeouts from environment (optional)
	if writeTimeout := os.Getenv(REDIS_CONN_WRITE_TIMEOUT_ENV_VAR); writeTimeout != "" {
		duration, err := time.ParseDuration(writeTimeout)
		if err != nil {
			return fmt.Errorf("invalid REDIS_CONN_WRITE_TIMEOUT value: %w", err)
		}
		c.ConnWriteTimeout = duration
	}

	if readTimeout := os.Getenv(REDIS_CONN_READ_TIMEOUT_ENV_VAR); readTimeout != "" {
		duration, err := time.ParseDuration(readTimeout)
		if err != nil {
			return fmt.Errorf("invalid REDIS_CONN_READ_TIMEOUT value: %w", err)
		}
		c.ConnReadTimeout = duration
	}

	// Load operation-specific timeouts
	if opReadTimeout := os.Getenv(REDIS_TIMEOUT_READ_ENV_VAR); opReadTimeout != "" {
		duration, err := time.ParseDuration(opReadTimeout)
		if err != nil {
			return fmt.Errorf("invalid REDIS_TIMEOUT_READ value: %w", err)
		}
		c.TimeoutConfig.Read = duration
	}

	if opWriteTimeout := os.Getenv(REDIS_TIMEOUT_WRITE_ENV_VAR); opWriteTimeout != "" {
		duration, err := time.ParseDuration(opWriteTimeout)
		if err != nil {
			return fmt.Errorf("invalid REDIS_TIMEOUT_WRITE value: %w", err)
		}
		c.TimeoutConfig.Write = duration
	}

	if opDefaultTimeout := os.Getenv(REDIS_TIMEOUT_DEFAULT_ENV_VAR); opDefaultTimeout != "" {
		duration, err := time.ParseDuration(opDefaultTimeout)
		if err != nil {
			return fmt.Errorf("invalid REDIS_TIMEOUT_DEFAULT value: %w", err)
		}
		c.TimeoutConfig.Default = duration
	}

	return nil
}

func (c *RueidisConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("redis host cannot be empty")
	}
	if c.Port == 0 || c.Port < MIN_REDIS_PORT || c.Port > MAX_REDIS_PORT {
		return fmt.Errorf("invalid Redis port: %d", c.Port)
	}

	// Validate timeout values
	if c.TimeoutConfig.Read <= 0 {
		return fmt.Errorf("TimeoutConfig.Read must be positive")
	}
	if c.TimeoutConfig.Write <= 0 {
		return fmt.Errorf("TimeoutConfig.Write must be positive")
	}
	if c.TimeoutConfig.Default <= 0 {
		return fmt.Errorf("TimeoutConfig.Default must be positive")
	}

	if c.ConnWriteTimeout <= 0 {
		return fmt.Errorf("ConnWriteTimeout must be positive")
	}
	if c.ConnReadTimeout <= 0 {
		return fmt.Errorf("ConnReadTimeout must be positive")
	}
	if c.BlockingOpsPoolSize < 1 {
		return fmt.Errorf("BlockingOpsPoolSize must be at least 1")
	}
	return nil
}

func (c *RueidisConfig) GetRueidisOptions() rueidis.ClientOption {
	return rueidis.ClientOption{
			InitAddress:      []string{fmt.Sprintf("%s:%d", c.Host, c.Port)},
			Username:         RUEIDIS_CLIENT_OPTION_USERNAME,
			Password:         c.Password, // Use the environment-loaded value
			SelectDB:         c.DB,
			ConnWriteTimeout: c.ConnWriteTimeout,
			Dialer: net.Dialer{
					Timeout: c.ConnWriteTimeout,
			},
			BlockingPoolSize: c.BlockingOpsPoolSize,
	}
}
