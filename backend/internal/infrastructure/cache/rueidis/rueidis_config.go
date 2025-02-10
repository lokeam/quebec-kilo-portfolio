package cache

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/redis/rueidis"
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
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return fmt.Errorf("REDIS_HOST environment variable is required")
	}
	c.Host = host

	// Required: Port
	portStr := os.Getenv("REDIS_PORT")
	if portStr == "" {
		return fmt.Errorf("REDIS_PORT environment variable is required")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid REDIS_PORT value: %w", err)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("REDIS_PORT must be between 1 and 65535")
	}
	c.Port = port

	// Optional: Password
	if pass := os.Getenv("REDIS_PASSWORD"); pass != "" {
		c.Password = pass
	}

	// NOTE: Save this debug for later when we re-implement Redis ACL
	// fmt.Printf("DEBUG: REDIS_PASSWORD=%q\n", c.Password)

	// Optional: Database number
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		db, err := strconv.Atoi(dbStr)
		if err != nil {
			return fmt.Errorf("invalid REDIS_DB value: %w", err)
		}
		if db < 0 {
			return fmt.Errorf("REDIS_DB must be non-negative")
		}
		c.DB = db
	}

	// Load connection timeouts from environment (optional)
	if writeTimeout := os.Getenv("REDIS_CONN_WRITE_TIMEOUT"); writeTimeout != "" {
		duration, err := time.ParseDuration(writeTimeout)
		if err != nil {
			return fmt.Errorf("invalid REDIS_CONN_WRITE_TIMEOUT value: %w", err)
		}
		c.ConnWriteTimeout = duration
	}

	if readTimeout := os.Getenv("REDIS_CONN_READ_TIMEOUT"); readTimeout != "" {
		duration, err := time.ParseDuration(readTimeout)
		if err != nil {
			return fmt.Errorf("invalid REDIS_CONN_READ_TIMEOUT value: %w", err)
		}
		c.ConnReadTimeout = duration
	}

	// Load operation-specific timeouts
	if opReadTimeout := os.Getenv("REDIS_TIMEOUT_READ"); opReadTimeout != "" {
		duration, err := time.ParseDuration(opReadTimeout)
		if err != nil {
			return fmt.Errorf("invalid REDIS_TIMEOUT_READ value: %w", err)
		}
		c.TimeoutConfig.Read = duration
	}

	if opWriteTimeout := os.Getenv("REDIS_TIMEOUT_WRITE"); opWriteTimeout != "" {
		duration, err := time.ParseDuration(opWriteTimeout)
		if err != nil {
			return fmt.Errorf("invalid REDIS_TIMEOUT_WRITE value: %w", err)
		}
		c.TimeoutConfig.Write = duration
	}

	if opDefaultTimeout := os.Getenv("REDIS_TIMEOUT_DEFAULT"); opDefaultTimeout != "" {
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
	if c.Port == 0 || c.Port < 1 || c.Port > 65535 {
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
			Username:         "default",
			Password:         c.Password, // Use the environment-loaded value
			SelectDB:         c.DB,
			ConnWriteTimeout: c.ConnWriteTimeout,
			Dialer: net.Dialer{
					Timeout: c.ConnWriteTimeout,
			},
			BlockingPoolSize: c.BlockingOpsPoolSize,
	}
}