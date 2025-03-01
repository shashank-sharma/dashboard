# Prometheus Metrics for PocketBase

This package provides Prometheus metrics integration for the PocketBase dashboard application.

## Enabling Metrics

Metrics are disabled by default. You can enable them in several ways:

1. **Command Line Flag**: Start PocketBase with the `--metrics` flag
   ```
   go run cmd/dashboard/main.go serve --metrics
   ```

2. **Environment Variable**: Set `METRICS_ENABLED=true` in your `.env` file or environment

3. **Custom Port**: To use a custom port, set `METRICS_PORT=<your_port>` (default is 9091)

4. **Runtime API**: Enable or disable metrics at runtime via API endpoints:
   ```
   # Get current metrics status
   GET /api/admin/metrics
   
   # Enable metrics (optionally with custom port)
   POST /api/admin/metrics/enable?port=9095
   
   # Disable metrics
   POST /api/admin/metrics/disable
   ```

## Features

- HTTP request metrics (count, duration)
- Database operation metrics (creates, updates, deletes)
- Active session tracking
- Cron job execution metrics

## Metrics Exposed

The following metrics are available:

1. `pocketbase_http_requests_total` - Counter of HTTP requests by method, path, and status
2. `pocketbase_http_request_duration_seconds` - Histogram of HTTP request durations
3. `pocketbase_database_operations_total` - Counter of database operations by collection and operation type
4. `pocketbase_active_sessions` - Gauge of active user sessions
5. `pocketbase_cronjob_executions_total` - Counter of cron job executions by job name
6. `pocketbase_cronjob_duration_seconds` - Histogram of cron job execution durations

## How It Works

The metrics integration works by:

1. Setting up hooks on PocketBase to track database operations
2. Starting a dedicated HTTP server on port 9091 (or your configured port) to expose metrics to Prometheus
3. Providing utility functions to track custom metrics (cron jobs, sessions)

## Testing Locally

To test the Prometheus integration:

1. Start your PocketBase application with metrics enabled:
   ```
   go run cmd/dashboard/main.go serve --metrics
   ```

2. In a separate terminal, run the provided Prometheus test script:
   ```
   ./scripts/test_prometheus.sh
   ```

3. Access the Prometheus UI at http://localhost:9090
   - Go to the "Graph" tab
   - Enter a query like `pocketbase_http_requests_total` to see metrics
   - Use the "Execute" button to run the query

4. Generate some traffic by using the PocketBase application

5. The raw metrics are available at http://localhost:9091/metrics (or your configured port)

## Integrating with an Existing Prometheus Server

If you already have Prometheus running, add the following to your `prometheus.yml` configuration:

```yaml
scrape_configs:
  - job_name: 'pocketbase'
    static_configs:
      - targets: ['your-pocketbase-host:9091']
```

## Grafana Dashboard

For a better visualization experience, you can set up Grafana with the following steps:

1. Install and run Grafana
2. Add your Prometheus server as a data source
3. Import a dashboard for PocketBase metrics (sample JSON configuration coming soon) 