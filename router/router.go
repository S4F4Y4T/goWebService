package router

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/internal/app"
	"github.com/S4F4Y4T/goWebService/internal/handler"
	"github.com/S4F4Y4T/goWebService/internal/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupRoutes(app *app.App) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthHandler)

	// Register routes using handlers from the App container
	RegisterUserRoutes(mux, app.UserHandler)
	RegisterProductRoutes(mux, app.ProductHandler)

	handler := middleware.Apply(middleware.Recover, middleware.CorrelationID, middleware.Logger, middleware.Cors)(mux)
	return otelhttp.NewHandler(handler, "http-server")
}
