// File: internal/appcontext_test/app_context_test_helper.go
package appcontext_test

import (
	"time"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/testutils"
)

// NewTestingAppContext creates an AppContext for testing.
func NewTestingAppContext(tokenKey string, tokenErr error) *appcontext.AppContext {
    return &appcontext.AppContext{
        Logger: testutils.NewTestLogger(), // Using TestLogger from testutils (which now has no prod deps)
        Config: &config.Config{
            // Adjust according to your production Config type.
            IGDB: &config.IGDBConfig{
							AccessTokenKey:   tokenKey,
							ClientID:         tokenKey,
							ClientSecret:     tokenKey,
							AuthURL:          tokenKey,
							BaseURL:          tokenKey,
							TokenTTL:         24 * time.Hour,
					},
        },
        // Other fields (MemCache, RedisClient, etc.) can be left nil or initialized as needed.
    }
}