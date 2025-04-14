package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Config holds database migration configuration
type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	SSLMode      string
	MaxRetries   int
	RetryTimeout time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Host:         os.Getenv("POSTGRES_HOST"),
		Port:         "5432",
		User:         os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		Database:     os.Getenv("POSTGRES_DB"),
		SSLMode:      "disable",
		MaxRetries:   5,
		RetryTimeout: 5 * time.Second,
	}
}

// DSN returns the database connection string
func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}

// RunMigrations applies all pending migrations with retries
func RunMigrations(ctx context.Context, config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	var db *sql.DB
	var err error

	// Retry connection with backoff
	for i := 0; i < config.MaxRetries; i++ {
		db, err = sql.Open("postgres", config.DSN())
		if err != nil {
			log.Printf("Attempt %d: Failed to open database: %v", i+1, err)
			time.Sleep(config.RetryTimeout)
			continue
		}

		// Test the connection
		if err = db.PingContext(ctx); err == nil {
			break
		}

		log.Printf("Attempt %d: Failed to ping database: %v", i+1, err)
		db.Close()
		time.Sleep(config.RetryTimeout)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database after %d attempts: %w",
			config.MaxRetries, err)
	}
	defer db.Close()

	// Create a new Postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating postgres driver: %w", err)
	}

	// Get the migrations directory path
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		// Default to ./migrations relative to the working directory
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting working directory: %w", err)
		}
		migrationsPath = fmt.Sprintf("%s/migrations", wd)
	}

	// Create a new migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}

	// Apply all pending migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}