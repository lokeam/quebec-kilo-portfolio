package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	memcache "github.com/lokeam/qko-beta/internal/infrastructure/cache/memorycache"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/shared/connectionutil"
	"github.com/lokeam/qko-beta/internal/shared/redisclient"
	"github.com/lokeam/qko-beta/internal/shared/token"
	"github.com/lokeam/qko-beta/internal/shared/twitch"
	"github.com/lokeam/qko-beta/internal/testutils"
)

// mockRueidisClient is our test implementation for the redis connectivity check
// Provides specific methods used by InitIGDBJob.
type mockRueidisClient struct{}

func (f *mockRueidisClient) Ping(ctx context.Context) error {
	return nil // Always succeed
}

func (f *mockRueidisClient) IsReady() bool {
	return true // Always ready
}

func (m *mockRueidisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil // Always succeed
}

func TestInitIGDBJob_HappyPath(t *testing.T) {
	t.Run(
		"Happy path - All dependencies succeed immediately, when InitIGDBJob is called it returns nil",
		func(t *testing.T) {
			// Note: We're overriding the net connectivity check to always succeed
			overrideIsOnline := connectionutil.IsOnline
			connectionutil.IsOnline = func(
				host string,
				port int,
				timeout time.Duration,
			) bool {
				return true
			}
			defer func() { connectionutil.IsOnline = overrideIsOnline }()

			// Note: We're overriding twitch token refresh so that it immediately succeeds
			overrideRefreshToken := twitch.RefreshToken
			twitch.RefreshToken = func(
				ctx context.Context,
				clientID,
				clientSecret,
				authURL string,
				logger interfaces.Logger,
			) (*twitch.TokenResponse,error) {
				return &twitch.TokenResponse{
					AccessToken: "test-access-token",
					TokenType: "test-token-type",
					ExpiresIn: 360,
				}, nil
			}
			defer func() { twitch.RefreshToken = overrideRefreshToken }()

			// Note: We're overriding the token job so that it returns no error
			overrideTokenJob := UpdateTwitchTokenJob
			UpdateTwitchTokenJob = func(
				ctx context.Context,
				redisKey string,
				redisClient redisclient.RedisClient,
				memCache *memcache.MemoryCache,
				tokenInfo token.TokenInfo,
				log interfaces.Logger,
			) error {
				return nil
			}
			defer func() { UpdateTwitchTokenJob = overrideTokenJob }()

			ctx := context.Background()
			redisKey := "test-redis-key"

			// Create a fake redis client
			mockRedis := &mockRueidisClient{}

			// For memcache, we're assuming that a dummy instance is acceptable
			mockMemCache := &memcache.MemoryCache{}
			clientID := "test-client-id"
			clientSecret := "test-client-secret"
			authURL := "http://test-auth-url"
			testLogger := testutils.NewTestLogger()

			// ACT
			err := InitIGDBJob(
				ctx,
				redisKey,
				mockRedis,
				mockMemCache,
				clientID,
				clientSecret,
				authURL,
				testLogger,
			)

			// ASSERT
			if err != nil {
				t.Fatalf("Expected no error no happy path igdb test, but we got: %v", err)
			}

			expectedMessages := []string{
				"Starting INIT_IGDB job",
				"Internet connection available",
				"Redis connection ready",
				"Successfully refreshed Twitch token",
				"Successfully saved Twitch token in Redis + Memcache",
			}

			for _, msg := range expectedMessages {
				foundMsg := false
				for _, loggedMsg := range testLogger.InfoCalls {
					if loggedMsg == msg {
						foundMsg = true
						break
					}
				}
				if !foundMsg {
					t.Errorf("Expected message '%s' not found in info logs", msg)
				}
			}
		},
	)
}

// Sad path: context is cancelled, job exits early
func TestInitIGDBJob_ContextCancelled(t *testing.T) {
	t.Run("Sad path: context is cancelled before dependencies are ready, job exits early",
		func(t *testing.T) {
			// Force internet connectivity check to always fail
			overrideIsOnline := connectionutil.IsOnline
			connectionutil.IsOnline = func(host string, port int, timeout time.Duration) bool {
				return false
			}
			defer func() { connectionutil.IsOnline = overrideIsOnline }()

			// Force a context cancellation
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			redisKey := "test-redis-key"
			mockRedis := &mockRueidisClient{}
			mockMemCache := &memcache.MemoryCache{}
			clientID := "test-client-id"
			clientSecret := "test-client-secret"
			authURL := "http://test-auth-url"
			testLogger := testutils.NewTestLogger()

			// ACT
			err := InitIGDBJob(ctx, redisKey, mockRedis, mockMemCache, clientID, clientSecret, authURL, testLogger)

			// ASSERT
			if err == nil {
				t.Fatalf("Expected an error due to canceled context, but got nil")
			}

			if !errors.Is(err, context.Canceled) {
				t.Fatalf("Expected error context.Canceled but got: %v", err)
			}
		},
	)
}
