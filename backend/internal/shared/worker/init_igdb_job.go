package worker

import (
	"context"
	"time"

	memcache "github.com/lokeam/qko-beta/internal/infrastructure/cache/memorycache"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/shared/connectionutil"
	"github.com/lokeam/qko-beta/internal/shared/redisclient"
	"github.com/lokeam/qko-beta/internal/shared/token"
	"github.com/lokeam/qko-beta/internal/shared/twitch"
)

// InitIGDBJob is the job that initializes tasks for IGDB access
func InitIGDBJob(
	ctx context.Context,
	redisKey string,
	rueidisClient redisclient.RedisClient,
	memCache *memcache.MemoryCache,
	clientID string,
	clientSecret string,
	authURL string,
	log interfaces.Logger,
) error {
	log.Info("Starting INIT_IGDB job", nil)

	// Check internet connectivity with retries
	if err := waitForInternetConnection(ctx, log); err != nil {
		return err
	}

	// Check Redis connectivity with retries
	if err := waitForRedisConnection(ctx, rueidisClient, log); err != nil {
		return err
	}

	// Get Twitch token with retries
	tokenInfo, err := GetTwitchTokenWithRetry(ctx, clientID, clientSecret, authURL, log)
	if err != nil {
		return err
	}

	// Save token with retries
	if err := saveTokenWithRetry(ctx, redisKey, rueidisClient, memCache, tokenInfo, log); err != nil {
		return err
	}

	log.Info("Successfully initialized IGDB access", nil)
	return nil
}

func waitForInternetConnection(ctx context.Context, log interfaces.Logger) error {
	for {
		if connectionutil.IsOnline("www.google.com", 80, 30*time.Second) {
			log.Info("Internet connection available", nil)
			return nil
		}
		log.Warn("Internet connection not available; retrying in 30 seconds", nil)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(30 * time.Second):
		}
	}
}

func waitForRedisConnection(ctx context.Context, client redisclient.RedisClient, log interfaces.Logger) error {
	for {
		if client.Ping(ctx) == nil && client.IsReady() {
			log.Info("Redis connection ready", nil)
			return nil
		}
		log.Warn("Redis not available; retrying in 3 seconds", nil)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(3 * time.Second):
		}
	}
}

func GetTwitchTokenWithRetry(ctx context.Context, clientID, clientSecret, authURL string, log interfaces.Logger) (*token.TokenInfo, error) {
	for {
		twitchTokenResponse, err := twitch.RefreshToken(
			ctx,
			clientID,
			clientSecret,
			authURL,
			log,
		)
		if err == nil {
			log.Info("Successfully refreshed Twitch token", nil)
			return &token.TokenInfo{
				AccessToken: twitchTokenResponse.AccessToken,
				ExpiresAt:   time.Now().Add(time.Duration(twitchTokenResponse.ExpiresIn) * time.Second),
			}, nil
		}
		log.Error("Failed to refresh Twitch token; retrying in 30 seconds", nil)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(30 * time.Second):
		}
	}
}

func saveTokenWithRetry(ctx context.Context, redisKey string, client redisclient.RedisClient, memCache *memcache.MemoryCache, tokenInfo *token.TokenInfo, log interfaces.Logger) error {
	for {
		err := UpdateTwitchTokenJob(
			ctx,
			redisKey,
			client,
			memCache,
			*tokenInfo,
			log,
		)
		if err == nil {
			log.Info("Successfully saved Twitch token in Redis + Memcache", nil)
			return nil
		}
		log.Error("Failed to save Twitch token in Redis; retrying in 15 seconds", nil)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(15 * time.Second):
		}
	}
}