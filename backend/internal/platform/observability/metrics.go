// Package observability wires Prometheus metrics and OpenTelemetry tracing into
// the platform. Both degrade gracefully: /metrics is always available, while
// tracing is a no-op unless OTEL_EXPORTER_OTLP_ENDPOINT is configured.
package observability

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests by method, route, and status code.",
	}, []string{"method", "route", "status"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency by method and route.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route"})

	httpInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "In-flight HTTP requests.",
	})

	cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "Cache-aside hits.",
	})
	cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "Cache-aside misses.",
	})
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, httpInFlight, cacheHits, cacheMisses)
}

// Handler returns the Prometheus exposition handler for /metrics.
func Handler() http.Handler { return promhttp.Handler() }

// RecordCacheHit / RecordCacheMiss are called by the cache layer to track hit
// rate. They are always safe to call.
func RecordCacheHit()  { cacheHits.Inc() }
func RecordCacheMiss() { cacheMisses.Inc() }

// statusRecorder captures the response status code for metrics.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Unwrap exposes the underlying writer so http.ResponseController and otelhttp
// can reach optional interfaces (Flusher, Hijacker, …).
func (r *statusRecorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

// MetricsMiddleware records RED metrics (rate, errors, duration) per route. It
// must wrap the ServeMux directly so it can read the matched route pattern
// (r.Pattern) after the request is served, keeping label cardinality bounded.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpInFlight.Inc()
		defer httpInFlight.Dec()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rec, r)
		elapsed := time.Since(start).Seconds()

		route := r.Pattern
		if route == "" {
			route = "unmatched"
		}
		httpRequestDuration.WithLabelValues(r.Method, route).Observe(elapsed)
		httpRequestsTotal.WithLabelValues(r.Method, route, strconv.Itoa(rec.status)).Inc()
	})
}

// RegisterDBStats exposes connection-pool gauges sourced from sql.DB.Stats().
func RegisterDBStats(db *sql.DB) {
	gauge := func(name, help string, read func(sql.DBStats) float64) {
		prometheus.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: name, Help: help,
		}, func() float64 { return read(db.Stats()) }))
	}
	gauge("db_connections_open", "Open DB connections (in use + idle).",
		func(s sql.DBStats) float64 { return float64(s.OpenConnections) })
	gauge("db_connections_in_use", "DB connections currently in use.",
		func(s sql.DBStats) float64 { return float64(s.InUse) })
	gauge("db_connections_idle", "Idle DB connections.",
		func(s sql.DBStats) float64 { return float64(s.Idle) })
	gauge("db_wait_count_total", "Total number of connections waited for.",
		func(s sql.DBStats) float64 { return float64(s.WaitCount) })
	gauge("db_wait_seconds_total", "Total time blocked waiting for a connection.",
		func(s sql.DBStats) float64 { return s.WaitDuration.Seconds() })
}
