package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const CorrelationIDKey contextKey = "correlation_id"

// CorrelationID is a middleware that injects a unique ID into the request context.
// It looks for an existing X-Correlation-ID header, and generates one if not found.
func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Correlation-ID")
		if id == "" {
			id = uuid.New().String()
		}

		// Set the correlation ID in the response header
		w.Header().Set("X-Correlation-ID", id)

		// Inject the correlation ID into the context
		ctx := context.WithValue(r.Context(), CorrelationIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCorrelationID extracts the correlation ID from the context.
func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return id
	}
	return ""
}
