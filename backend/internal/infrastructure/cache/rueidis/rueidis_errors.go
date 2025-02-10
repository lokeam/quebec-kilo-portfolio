package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/rueidis"
)

var (
	ErrorKeyNotFound          = errors.New("key not found in redis")
	ErrorConnectionFailed     = errors.New("redis conneciton failed")
	ErrorClientNotReady       = errors.New("redis client not ready")
	ErrorTimeout              = errors.New("redis operations timed out")
)

// Structs / Interfaces
type RueidisError struct {
	RueidisOperation   string
	Key                string
	Err                error
}

// Constructor
func NewRueidisError(
	rueidisOperation string,
	key string,
	err error,
) error {
	return &RueidisError{
		RueidisOperation: rueidisOperation,
		Key:              key,
		Err:              err,
	}
}

// Methods
func (e *RueidisError) Error() string {
	if e.Key != "" {
		return fmt.Sprintf("redis %s operation failed for key '%s': '%v'",
			e.RueidisOperation,
			e.Key,
			e.Err,
		)
	}

	return fmt.Sprintf("redis %s operation failed: '%v'", e.RueidisOperation, e.Err)
}

func (e *RueidisError) Unwrap() error {
	return e.Err
}

// Checks if error is specific to a Redis error type
func IsRedisError(err error) bool {
	if err == nil {
		return false
	}

	// rueidis.IsRedisErr returns (*RedisError, bool)
  // We only care about the bool indicating if it's a Redis error
	_, isRedisErr := rueidis.IsRedisErr(err)
	return isRedisErr
}

// Convert rueidis errors to custom error types
func (c *RueidisClient) ConvertRueidisError(
	err error,
	rueidisOperation string,
) error {
	if err == nil {
		return nil
	}

	// Log the original error for debugging
	c.logger.Debug("redis error", map[string]any{
		"error": err,
		"operation": rueidisOperation,
	})

	switch {
	case rueidis.IsRedisNil(err):
			return ErrorKeyNotFound
	case errors.Is(err, context.DeadlineExceeded):
			return ErrorTimeout
	case !c.IsReady():
			return ErrorClientNotReady
	case IsRedisError(err):
			return ErrorConnectionFailed
	default:
			return err
	}
}
