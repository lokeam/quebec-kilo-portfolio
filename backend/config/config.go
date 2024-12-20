package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server ServerConfig
	Env    string
	Debug  bool
}

type ServerConfig struct {
	Port int
	Host string
}

func Load() (*Config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	} else if env != "dev" && env != "test" && env != "prod" {
		return nil, fmt.Errorf("invalid environment: must be one of dev, test or prod")
	}

	// Debug mode handling
	debug := false
	if debugStr := os.Getenv("APP_DEBUG"); debugStr == "true" {
		debug = true
	}

	// Port handling
	portStr := os.Getenv("PORT")
	port := 8080 // default
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
	host := os.Getenv("HOST")
	// Check if HOST was explicitly set to empty
	if _, exists := os.LookupEnv("HOST"); exists && host == "" {
			return nil, fmt.Errorf("host cannot be empty")
	}
	// Use default if HOST not set
	if host == "" {
		host = "localhost"
	}

	return &Config{
		Server: ServerConfig{
				Port: port,
				Host: host,
		},
		Env: env,
		Debug: debug,
	}, nil
}