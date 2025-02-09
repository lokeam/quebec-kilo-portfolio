package worker

import (
	"context"
	"time"

	memcache "github.com/lokeam/qko-beta/internal/infrastructure/cache/memorycache"
	cache "github.com/lokeam/qko-beta/internal/infrastructure/cache/rueidis"
	"github.com/lokeam/qko-beta/internal/shared/connectionutil"
	"github.com/lokeam/qko-beta/internal/shared/logger"
	"github.com/lokeam/qko-beta/internal/shared/token"
	"github.com/lokeam/qko-beta/internal/shared/twitch"
)

// InitIGDBJob is the job that initializes tasks for IGDB access
func InitIGDBJob(
	ctx context.Context,
	redisKey string,
	rueidisClient *cache.RueidisClient,
	memCache *memcache.MemoryCache,
	clientID string,
	clientSecret string,
	authURL string,
	log logger.Logger,
) error {

	log.Info("Starting INIT_IGDB job", nil)

	// Check internet connectivity every 30 seconds unti online
	for {
		if connectionutil.IsOnline("www.google.com", 80, 30*time.Second) {
			log.Info("Internet connection available", nil)
      break // Connection is available; proceed to next steps.
		}
		log.Warn("Internet connection not available; retrying in 30 seconds", nil)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(30 * time.Second):
		}
	}

	// Check Redis connectivity every 3 seconds until ready
	for {
		if rueidisClient.Ping(ctx) == nil && rueidisClient.IsReady() {
			log.Info("Redis connection ready", nil)
			break
		}
		log.Warn("Redis not available; retrying in 3 seconds", nil)

		select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(3 * time.Second):
		}
	}

	// Call Refresh Twitch token every 30 seconds until success
	var twitchTokenResponse *twitch.TokenResponse
	for {
		var err error

		twitchTokenResponse, err = twitch.RefreshToken(
			ctx,
			clientID,
			clientSecret,
			authURL,
			&log,
		)
		if err == nil {
			log.Info("Successfully refreshed Twitch token", nil)
			break
		}
		log.Error("Failed to refresh Twitch token; retrying in 30 seconds", nil)

		select {
			case <- ctx.Done():
				return ctx.Err()
			case <-time.After(30 * time.Second):
		}
	}

	// Compute expiration time (separate this out into a fn)
	tokenExpiration := time.Now().Add(time.Duration(twitchTokenResponse.ExpiresIn) * time.Second)
	tokenInfo := token.TokenInfo {
		AccessToken: twitchTokenResponse.AccessToken,
		ExpiresAt: tokenExpiration,
	}

	for {
		err := UpdateTwitchTokenJob(
			ctx,
			redisKey,
			rueidisClient,
			memCache,
			tokenInfo,
			&log,
		)
		if err == nil {
			log.Info("Successfully saved Twitch token in Redis + Memcache", nil)
			break
		}
		log.Error("Failed to save Twitch token in Redis; retrying in 15 seconds", nil)

		select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(15 * time.Second):
		}
	}

	// Attempt to save token in memcache + Redis until success, retrying every 15 seconds

	return  nil
}