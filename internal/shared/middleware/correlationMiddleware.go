package middleware

import (
	"context"
	"net/http"

	"github.com/S4F4Y4T/goWebService/pkg/correlation"
	"github.com/google/uuid"
)

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

		// Inject the correlation ID into the context using the shared pkg key
		ctx := context.WithValue(r.Context(), correlation.CorrelationIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCorrelationID extracts the correlation ID from the context.
// Delegates to the shared pkg/correlation package.
func GetCorrelationID(ctx context.Context) string {
	return correlation.GetCorrelationID(ctx)
}
