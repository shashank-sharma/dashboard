package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shashank-sharma/backend/internal/logger"
)

var (
	// RequestCounter tracks HTTP requests by method, path, and status code
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pocketbase_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// RequestDuration tracks HTTP request duration
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pocketbase_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// DatabaseOperations tracks database operations
	DatabaseOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pocketbase_database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"collection", "operation"},
	)

	// ActiveSessions tracks the number of active user sessions
	ActiveSessions = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "pocketbase_active_sessions",
			Help: "Number of active user sessions",
		},
	)

	// CronJobExecutions tracks cron job executions
	CronJobExecutions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pocketbase_cronjob_executions_total",
			Help: "Total number of cron job executions",
		},
		[]string{"job_name"},
	)

	// CronJobDuration tracks cron job execution duration
	CronJobDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pocketbase_cronjob_duration_seconds",
			Help:    "Cron job execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"job_name"},
	)
)

// metricsResponseWriter is a wrapper around http.ResponseWriter to capture the status code
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before delegating to the wrapped ResponseWriter
func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}

// newMetricsResponseWriter creates a new metricsResponseWriter
func newMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{w, http.StatusOK}
}

// RegisterPrometheusMetrics registers Prometheus metrics with the application
func RegisterPrometheusMetrics(app *pocketbase.PocketBase) {
	// Check if metrics are enabled from store
	enabled, enabledOk := app.Store().Get("METRICS_ENABLED").(bool)
	if !enabledOk || !enabled {
		logger.LogInfo("Prometheus metrics are disabled")
		return
	}
	
	logger.LogInfo("Metrics enabled from runtime store")
	
	// Get port from store
	port, portOk := app.Store().Get("METRICS_PORT").(int)
	if !portOk || port <= 0 {
		// Use default port if not found or invalid
		port = 9091
		logger.LogInfo(fmt.Sprintf("Using default metrics port %d", port))
	} else {
		logger.LogInfo(fmt.Sprintf("Using metrics port %d from runtime store", port))
	}

	logger.LogInfo(fmt.Sprintf("Enabling Prometheus metrics with port %d", port))

	// Set up the metrics port in the store for use by StartMetricsServer
	app.Store().Set("METRICS_PORT", port)
	
	// Track record operations for all collections
	for _, collection := range []string{"users", "track_devices", "tokens"} {
		// Use a separate variable for each closure to avoid issues with the loop variable
		collectionName := collection
		
		// Track creates
		app.OnRecordCreate(collectionName).BindFunc(func(e *core.RecordEvent) error {
			DatabaseOperations.WithLabelValues(collectionName, "create").Inc()
			return e.Next()
		})
		
		// Track updates
		app.OnRecordUpdate(collectionName).BindFunc(func(e *core.RecordEvent) error {
			DatabaseOperations.WithLabelValues(collectionName, "update").Inc()
			return e.Next()
		})
		
		// Track deletes
		app.OnRecordDelete(collectionName).BindFunc(func(e *core.RecordEvent) error {
			DatabaseOperations.WithLabelValues(collectionName, "delete").Inc()
			return e.Next()
		})
	}
}

// StartMetricsServer starts the Prometheus metrics HTTP server
func StartMetricsServer(app *pocketbase.PocketBase) error {
	enabled, ok := app.Store().Get("METRICS_ENABLED").(bool)
	if !ok || !enabled {
		return nil
	}
	
	// Get the port
	port, ok := app.Store().Get("METRICS_PORT").(int)
	if !ok || port <= 0 {
		port = 9091
	}
	
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: promhttp.Handler(),
	}

	go func() {
		logger.LogInfo(fmt.Sprintf("Starting Prometheus metrics server on :%d", port))
		err := metricsServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.LogError("Metrics server error", "error", err)
		}
	}()
	
	return nil
}

// TrackCronJobExecution records cron job execution metrics
func TrackCronJobExecution(jobName string, duration time.Duration) {
	CronJobExecutions.WithLabelValues(jobName).Inc()
	CronJobDuration.WithLabelValues(jobName).Observe(duration.Seconds())
}

// UpdateActiveSessions updates the active sessions gauge
func UpdateActiveSessions(count int) {
	ActiveSessions.Set(float64(count))
}