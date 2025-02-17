package mocks

import (
	"time"

	"github.com/lokeam/qko-beta/config"
)

func NewMockConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
		Env:   "test",
		Debug: true,
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders: []string{"Link"},
			AllowCredentials: true,
			MaxAge: 300,
		},
		IGDB: &config.IGDBConfig{
			ClientID:       "dummy-client-id",
			ClientSecret:   "dummy-client-secret",
			AuthURL:        "https://dummy-auth-url.com",   // Dummy Twitch OAuth URL
			BaseURL:        "https://dummy-base-url.com",    // Dummy IGDB API base URL
			TokenTTL:       24 * time.Hour,                  // Dummy token TTL
			AccessTokenKey: "dummy-access-token-key",
		},
		Redis: config.RedisConfig{
			RedisTimeout: 5 * time.Second,
			RedisTTL: 5 * time.Minute,
		},
	}
}
