package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/S4F4Y4T/goWebService/pkg/response"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Record the panic safely with a stacktrace
				slog.Error("PANIC recovered in handler", "panic_err", err, "stack", string(debug.Stack()))
				
				// Standardize a safe 500 fallback error
				response.Error(w, errors.New("internal server error"))
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}
