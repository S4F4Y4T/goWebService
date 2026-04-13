package router

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/internal/product"
	"github.com/S4F4Y4T/goWebService/internal/shared/middleware"
	"github.com/S4F4Y4T/goWebService/internal/user"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Router struct {
	UserHandler    *user.Handler
	ProductHandler *product.Handler
}

func NewRouter(uh *user.Handler, ph *product.Handler) *Router {
	return &Router{
		UserHandler:    uh,
		ProductHandler: ph,
	}
}

func (r *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Register Domain Routes
	r.UserHandler.RegisterRoutes(mux)
	r.ProductHandler.RegisterRoutes(mux)

	// Apply Global Middleware
	h := middleware.Apply(
		middleware.Recover,
		middleware.CorrelationID,
		middleware.Logger,
		middleware.Cors,
	)(mux)

	return otelhttp.NewHandler(h, "http-server")
}
