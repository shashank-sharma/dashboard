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
	DatabaseOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pocketbase_database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"collection", "operation"},
	)

	ActiveSessions = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "pocketbase_active_sessions",
			Help: "Number of active user sessions",
		},
	)

	CronJobExecutions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pocketbase_cronjob_executions_total",
			Help: "Total number of cron job executions",
		},
		[]string{"job_name"},
	)

	CronJobDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pocketbase_cronjob_duration_seconds",
			Help:    "Cron job execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"job_name"},
	)
)

// Initialize metrics with default values so they show up in Prometheus
func initializeMetrics() {	
	for _, collection := range []string{"users", "track_devices", "tokens"} {
		for _, op := range []string{"create", "update", "delete", "auth"} {
			DatabaseOperations.WithLabelValues(collection, op).Add(0)
		}
	}
	
	// TODO: Initialize session metrics
	ActiveSessions.Set(0)
	
	// Initialize cron job metrics
	for _, job := range []string{"track-device"} {
		CronJobExecutions.WithLabelValues(job).Add(0)
		CronJobDuration.WithLabelValues(job).Observe(0)
	}
}

// RegisterPrometheusMetrics registers Prometheus metrics with the application
func RegisterPrometheusMetrics(app *pocketbase.PocketBase) {
	initializeMetrics()
		
	// Check if metrics server should be enabled
	enabled, enabledOk := app.Store().Get("METRICS_ENABLED").(bool)
	if !enabledOk || !enabled {
		logger.LogInfo("Prometheus metrics server is disabled")
		return
	}
	
	logger.LogInfo("Metrics server enabled from runtime store")
	
	// Get port from store
	port, portOk := app.Store().Get("METRICS_PORT").(int)
	if !portOk || port <= 0 {
		port = 9091
	}

	logger.LogInfo(fmt.Sprintf("Enabling Prometheus metrics server on port %d", port))

	app.Store().Set("METRICS_PORT", port)
	
	setupBasicHooks(app)
}

func setupBasicHooks(app *pocketbase.PocketBase) {
	collections := []string{"users", "track_devices", "tokens"}
	
	for _, collection := range collections {
		collName := collection
		
		app.OnRecordCreate(collName).BindFunc(func(e *core.RecordEvent) error {
			DatabaseOperations.WithLabelValues(collName, "create").Inc()
			return e.Next()
		})
		
		app.OnRecordUpdate(collName).BindFunc(func(e *core.RecordEvent) error {
			DatabaseOperations.WithLabelValues(collName, "update").Inc()
			return e.Next()
		})
		
		app.OnRecordDelete(collName).BindFunc(func(e *core.RecordEvent) error {
			DatabaseOperations.WithLabelValues(collName, "delete").Inc()
			return e.Next()
		})
	}
}

// StartMetricsServer starts the Prometheus metrics HTTP server
func StartMetricsServer(app *pocketbase.PocketBase) error {
	// Always start the metrics server regardless of the enabled flag
	// This ensures metrics are always available for Grafana
	
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

// TrackDatabaseOperation manually tracks a database operation
func TrackDatabaseOperation(collection, operation string) {
	DatabaseOperations.WithLabelValues(collection, operation).Inc()
}