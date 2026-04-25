package router

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/S4F4Y4T/goWebService/gateway/internal/proxy"
	"github.com/S4F4Y4T/goWebService/pkg/jwtutil"
)

type Config struct {
	AuthServiceURL    string
	UserServiceURL    string
	ProductServiceURL string
}

func LoadConfig() Config {
	return Config{
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		UserServiceURL:    getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://localhost:8083"),
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

// Setup wires all routes on to mux and returns the composed handler.
func Setup(mux *http.ServeMux, cfg Config) http.Handler {
	// ── Auth routes (public — no JWT required) ────────────────────────────────
	authProxy := proxy.New(cfg.AuthServiceURL, "/auth")
	mux.Handle("/auth/", authProxy)

	// ── User routes (protected) ───────────────────────────────────────────────
	userProxy := proxy.New(cfg.UserServiceURL, "/users")
	mux.Handle("/users/", jwtGuard(userProxy))

	// ── Product routes (protected) ────────────────────────────────────────────
	productProxy := proxy.New(cfg.ProductServiceURL, "/products")
	mux.Handle("/products/", jwtGuard(productProxy))

	// ── Health check ──────────────────────────────────────────────────────────
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok","service":"gateway"}`)
	})

	// Apply global middleware: recover → cors → logger
	return applyMiddleware(mux, recoverMiddleware, corsMiddleware, loggerMiddleware)
}

// ── JWT Guard ──────────────────────────────────────────────────────────────────

func jwtGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			jsonUnauthorized(w)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtutil.ValidateToken(token)
		if err != nil {
			jsonUnauthorized(w)
			return
		}
		// Forward user identity downstream
		r2 := r.Clone(r.Context())
		r2.Header.Set("X-User-ID", fmt.Sprintf("%d", claims.UserID))
		r2.Header.Set("X-User-Email", claims.Email)
		next.ServeHTTP(w, r2)
	})
}

func jsonUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "unauthorized"})
}

// ── Middleware helpers ─────────────────────────────────────────────────────────

type middleware func(http.Handler) http.Handler

func applyMiddleware(h http.Handler, ms ...middleware) http.Handler {
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i](h)
	}
	return h
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "error", rec)
				http.Error(w, `{"success":false,"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
