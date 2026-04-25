package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// New creates a reverse proxy that strips the given prefix before forwarding.
func New(target, stripPrefix string) http.Handler {
	targetURL, err := url.Parse(target)
	if err != nil {
		slog.Error("invalid proxy target", "target", target, "error", err)
		panic(err)
	}

	rp := httputil.NewSingleHostReverseProxy(targetURL)

	// Custom error handler so we return JSON-ish errors instead of plain text.
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("proxy error", "target", target, "path", r.URL.Path, "error", err)
		http.Error(w, `{"success":false,"error":"upstream service unavailable"}`, http.StatusBadGateway)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip the routing prefix (e.g. /users → /) before forwarding.
		if stripPrefix != "" {
			r2 := r.Clone(r.Context())
			r2.URL.Path = r.URL.Path[len(stripPrefix):]
			if r2.URL.Path == "" {
				r2.URL.Path = "/"
			}
			r2.URL.RawPath = ""
			rp.ServeHTTP(w, r2)
			return
		}
		rp.ServeHTTP(w, r)
	})
}
