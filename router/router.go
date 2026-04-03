package router

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/internal/handler"
	"github.com/S4F4Y4T/goWebService/internal/middleware"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthHandler)

	RegisterUserRoutes(mux)
	RegisterProductRoutes(mux)

	middleware := middleware.Apply(middleware.Logger, middleware.Cors)

	return middleware(mux)
}
