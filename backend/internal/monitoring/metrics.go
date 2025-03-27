package monitoring

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Postgres Database metrics
	DBConnectionsOpen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "qko_db_connections_open",
		Help: "The current number of open Postgres connections",
	})

	DBConnectionsMax = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "qko_db_connections_max",
		Help: "The maximum number of Postgres connections",
	})

	// HTTP Metrics
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "qko_http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// Redis Metrics
	RedisOperations = promauto.NewCounterVec(
    prometheus.CounterOpts{
      Name: "qko_redis_operations_total",
      Help: "Total number of Redis operations by type and result",
    },
    []string{"operation", "status"},
	)

	// Redis operation gauges - for overall stats
	RedisOperationGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_operation_current",
			Help: "Current count of Redis operations by type",
		},
		[]string{"type"},
	)

	// Redis operation duration
	RedisOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "redis_operation_duration_seconds",
			Help: "Duration of Redis operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "result"},
	)

	// Redis success rate
	RedisSuccessRate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_success_rate_percentage",
			Help: "Percentage of successful Redis operations",
		},
	)

	// Redis uptime
	RedisUptime = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_client_uptime_seconds",
			Help: "Uptime of the Redis client in seconds",
		},
	)

	// IGDB API Metrics
	IGDBRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qko_igdb_requests_total",
			Help: "Total number of IGDB requests",
		},
		[]string{"endpoint", "status"},
	)
)

// Register /metrics endpoint for Prometheus scraping
func SetupMetricsEndpoint(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}
