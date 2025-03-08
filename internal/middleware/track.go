package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shashank-sharma/backend/internal/logger"
)

var (
	// RequestCounter tracks HTTP requests by method, path group, and status code
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pocketbase_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route_group", "path", "status"},
	)

	// RequestDuration tracks HTTP request duration
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pocketbase_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
		},
		[]string{"method", "route_group", "path"},
	)
)

func LogRequest(e *core.RequestEvent) error {
	println("Request received:", e.Request.Method, e.Request.URL.Path)
	logger.LogDebug("Request received:", "method", e.Request.Method, "path", e.Request.URL.Path)
	return e.Next()
}

// RouteMetricsMiddleware automatically determines the route group from the URL path
// and records metrics with that information
func RouteMetricsMiddleware(e *core.RequestEvent) error {
	start := time.Now()
	
	routeGroup := extractRouteGroup(e.Request.URL.Path)
	
	normalizedPath := normalizePath(e.Request.URL.Path)
	
	responseWriter := newResponseWriterWrapper(e.Response)
	e.Response = responseWriter
	
	err := e.Next()
	
	status := responseWriter.StatusCode
	
	if err != nil && status == http.StatusOK {
		status = http.StatusInternalServerError
	}
	
	RequestCounter.WithLabelValues(
		e.Request.Method,
		routeGroup,
		normalizedPath,
		fmt.Sprintf("%d", status),
	).Inc()
	
	RequestDuration.WithLabelValues(
		e.Request.Method,
		routeGroup,
		normalizedPath,
	).Observe(time.Since(start).Seconds())
	
	return err
}

type responseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
}

func newResponseWriterWrapper(original http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: original,
		StatusCode:     http.StatusOK,
	}
}

func (rww *responseWriterWrapper) WriteHeader(statusCode int) {
	rww.StatusCode = statusCode
	rww.ResponseWriter.WriteHeader(statusCode)
}

func extractRouteGroup(path string) string {
	trimmedPath := strings.TrimPrefix(path, "/")
	segments := strings.Split(trimmedPath, "/")
	
	if len(segments) == 0 || segments[0] == "" {
		return "root"
	}
	
	if len(segments) == 1 {
		return segments[0]
	}
	
	if len(segments) >= 2 {
		return segments[0] + "_" + segments[1]
	}
	
	// Fallback
	return "unknown"
}

func normalizePath(path string) string {
	segments := strings.Split(path, "/")
	normalizedSegments := make([]string, 0, len(segments))
	
	for i, segment := range segments {
		if segment == "" {
			normalizedSegments = append(normalizedSegments, segment)
			continue
		}
		
		if i <= 2 {
			normalizedSegments = append(normalizedSegments, segment)
			continue
		}
		
		if isLikelyID(segment) {
			normalizedSegments = append(normalizedSegments, "{id}")
		} else {
			normalizedSegments = append(normalizedSegments, segment)
		}
	}
	
	return strings.Join(normalizedSegments, "/")
}

func isLikelyID(segment string) bool {
	if isNumeric(segment) {
		return true
	}
	
	if len(segment) > 32 && strings.Count(segment, "-") >= 4 {
		return true
	}
	
	if len(segment) >= 15 && len(segment) <= 20 {
		hasLettersAndNumbers := false
		hasOnlyValidChars := true
		
		for _, c := range segment {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c == '-') {
				hasOnlyValidChars = false
				break
			}
			if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				hasLettersAndNumbers = true
			}
		}
		
		if hasLettersAndNumbers && hasOnlyValidChars {
			return true
		}
	}
	
	return false
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}