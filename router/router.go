package router

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/internal/app"
	"github.com/S4F4Y4T/goWebService/internal/handler"
	"github.com/S4F4Y4T/goWebService/internal/middleware"
)

func SetupRoutes(app *app.App) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthHandler)

	// Register routes using handlers from the App container
	RegisterUserRoutes(mux, app.UserHandler)
	RegisterProductRoutes(mux, app.ProductHandler)

	return middleware.Apply(middleware.Logger, middleware.Cors)(mux)
}
