package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wait for next handlers to finish
		next.ServeHTTP(w, r)
		
		slog.Info("HTTP Request", 
			"method", r.Method, 
			"path", r.URL.Path, 
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
