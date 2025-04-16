package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	interfaces "github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/monitoring"
	"github.com/redis/rueidis"
)

type ClientStatus int32

const (
	StatusReady  ClientStatus = iota
	StatusError
	StatusClosed
)

type RueidisClient struct {
	client rueidis.Client
	logger interfaces.Logger
	config *RueidisConfig

	// Metrics go here
	stats *ClientStats

	status atomic.Int32
}

type ClientStats struct {
	Operations     atomic.Int64
	Errors         atomic.Int64
	LastOperation  atomic.Value
	StartTime      time.Time
}

type Stats struct {
	Operations    int64
	Errors        int64
	LastOperation time.Time
	StartTime     time.Time
	Uptime        time.Duration
}

// Constructor
func NewRueidisClient(
	rcfg *RueidisConfig,
	logger interfaces.Logger,
) (*RueidisClient, error) {
	// Guard clauses
	if rcfg == nil {
		return nil, fmt.Errorf("rueidis config cannot be nil")
	}

	// Convert config to rueidis options
	rueidisOptions := rcfg.GetRueidisOptions()

	// Create client
	rueidisClient, err := rueidis.NewClient(rueidisOptions)
	if err != nil {
		logger.Error("failed to create rueidis client", map[string]any{
			"error": err,
		})
		return nil, err
	}

	// initialize metric stats
	stats := &ClientStats {
		StartTime: time.Now(),
	}
	stats.LastOperation.Store(time.Now())

	// initialize client wrapper
	clientWrapper := &RueidisClient{
		client: rueidisClient,
		logger: logger,
		config: rcfg,
		stats: stats,
	}

	// set initial status to ready
	clientWrapper.status.Store(int32(StatusReady))

	return clientWrapper, nil
}

// Methods
func (c *RueidisClient) Get(ctx context.Context, key string) (string, error) {
	start := time.Now()
	c.logger.Debug("attempting to run cache get", map[string]any{
    "key":       key,
    "operation": "GET",
    "timestamp": start,
	})

	defer func() {
		c.stats.Operations.Add(1)
		c.stats.LastOperation.Store(time.Now())
	}()

	if !c.IsReady() {
		return "", ErrorClientNotReady
	}

	// Build and execute command
	cmd := c.client.B().Get().Key(key).Build()
	result, err := c.client.Do(ctx, cmd).ToString()

	duration := time.Since(start).Seconds()

	if err != nil {
		c.stats.Errors.Add(1)
		c.logger.Error("redis get failed", map[string]any{
				"key": key,
				"error": err,
				"duration": duration,
		})
		// Record error metric with duration - use Inc() for counters
		monitoring.RedisOperations.WithLabelValues("GET", "error").Inc()
		monitoring.RedisOperationDuration.WithLabelValues("GET", "error").Observe(duration)
		return "", c.ConvertRueidisError(err, "GET")
	}

	// Record success metric with duration - using Inc() for counters
	monitoring.RedisOperations.WithLabelValues("GET", "success").Inc()
	monitoring.RedisOperationDuration.WithLabelValues("GET", "success").Observe(duration)

	c.logger.Debug("redis get completed", map[string]any{
		"key": key,
		"duration": time.Since(start),
		"operation": "GET",
		"timestamp": time.Now(),
		"hit": result != "",
		"size": len(result),
	})

	return result, nil
}

func (c *RueidisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Convert value to string based on type
	start := time.Now()
	var strValue string
	switch v := value.(type) {
	case string:
			strValue = v
	case []byte:
			strValue = string(v)
	default:
			// For other types, use JSON marshaling
			jsonBytes, err := json.Marshal(value)
			if err != nil {
					c.logger.Error("failed to marshal value", map[string]any{
						"key": key,
						"valueType": fmt.Sprintf("%T", value),
						"error": err,
					})
					return fmt.Errorf("failed to marshal value: %w", err)
			}
			strValue = string(jsonBytes)
	}

	c.logger.Debug("attempting redis SET", map[string]any{
		"key": key,
		"valueType": fmt.Sprintf("%T", value),
		"valueSize": len(strValue),
		"ttl": expiration,
	})

	// Build and execute command with properly serialized value
	cmd := c.client.B().Set().Key(key).Value(strValue).Ex(expiration).Build()
	err := c.client.Do(ctx, cmd).Error()

	if err != nil {
		c.stats.Errors.Add(1)
		c.logger.Error("redis set failed", map[string]any{
				"key": key,
				"error": err,
				"duration": time.Since(start),
		})
		return c.ConvertRueidisError(err, "SET")
	}

	c.logger.Debug("redis set successful", map[string]any{
		"key": key,
		"duration": time.Since(start),
	})

	return nil
}

func (c *RueidisClient) Delete(ctx context.Context, key string) error {
	start := time.Now()
	defer func() {
		c.stats.Operations.Add(1)
		c.stats.LastOperation.Store(time.Now())
	}()

	if !c.IsReady() {
		return ErrorClientNotReady
	}

	// Build and execute command
	cmd := c.client.B().Del().Key(key).Build()
	err := c.client.Do(ctx, cmd).Error()

	if err != nil {
		c.stats.Errors.Add(1)
		c.logger.Error("redis delete failed", map[string]any{
			"key": key,
			"error": err,
			"duration": time.Since(start),
		})
		return c.ConvertRueidisError(err, "DELETE")
	}

	c.logger.Debug("redis delete successful", map[string]any{
		"key": key,
		"duration": time.Since(start),
	})

	return nil
}

func (c *RueidisClient) Close() error {
	c.status.Store(int32(StatusClosed))
	c.client.Close()
	c.logger.Info("redis client closed", map[string]any{
		"duration": time.Since(c.stats.StartTime),
	})

	return nil
}

func (c *RueidisClient) Ping(ctx context.Context) error {
	start := time.Now()
	defer func() {
		c.stats.Operations.Add(1)
		c.stats.LastOperation.Store(time.Now())

	}()

	err := c.client.Do(ctx, c.client.B().Ping().Build()).Error()
	if err != nil {
			c.stats.Errors.Add(1)
			c.logger.Error("redis ping failed", map[string]any{
				"error": err,
				"duration": time.Since(start),
			})
			return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

func (c *RueidisClient) GetStatus() ClientStatus {
	return ClientStatus(c.status.Load())
}

// Returns true if the client is ready to use
func (c *RueidisClient) IsReady() bool {
	return c.GetStatus() == StatusReady
}

// Returns current client statistics
func (c *RueidisClient) GetStats() Stats {
	now := time.Now()
	lastOperation, _ := c.stats.LastOperation.Load().(time.Time)

	return Stats{
		Operations:     c.stats.Operations.Load(),
		Errors:         c.stats.Errors.Load(),
		LastOperation:  lastOperation,
		StartTime:      c.stats.StartTime,
		Uptime:         now.Sub(c.stats.StartTime),
	}
}

func (c *RueidisClient) GetConfig() *RueidisConfig {
	return c.config
}

// Export client stats to Prometheus Metrics
func (c *RueidisClient) UpdateMetrics() {
	// Only update metrics if client is ready
	if c.IsReady() {
		// Get current stats
		stats := c.GetStats()

		// Update Prometheus metrics
		monitoring.RedisOperations.WithLabelValues("total", "all").Add(float64(stats.Operations))
		monitoring.RedisOperations.WithLabelValues("error", "all").Add(float64(stats.Errors))

		// Calculate success rate (avoid division by zero)
		if stats.Operations > 0 {
			successRate := float64(stats.Operations-stats.Errors) / float64(stats.Operations)
			monitoring.RedisSuccessRate.Set(successRate * 100) // as percentage
	}

		// Update uptime metric
		monitoring.RedisUptime.Set(float64(stats.Uptime.Seconds()))
	}
}

// Start goroutine to periodically update metrics
func (c *RueidisClient) StartMetricsCollection() {
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		// Update metrics immediately on start
		c.UpdateMetrics()

		// Then update on ticker intervals
		for range ticker.C {
			c.UpdateMetrics()
		}
	}()

	c.logger.Info("Redis metrics collection started", map[string]any{
		"interval": "15s",
	})
}
// Private methods
//func (c *RueidisClient) updateStats(err error) {}
