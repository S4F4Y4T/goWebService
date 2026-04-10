package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		correlationID := GetCorrelationID(r.Context())
		
		// Get trace ID from OpenTelemetry
		span := trace.SpanFromContext(r.Context())
		traceID := ""
		if span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
		}
		
		// Wait for next handlers to finish
		next.ServeHTTP(w, r)
		
		slog.Info("HTTP Request", 
			"method", r.Method, 
			"path", r.URL.Path, 
			"correlation_id", correlationID,
			"trace_id", traceID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
