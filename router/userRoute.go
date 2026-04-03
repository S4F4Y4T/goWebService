package router

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/internal/handler"
	"github.com/S4F4Y4T/goWebService/internal/middleware"
)

func RegisterUserRoutes(mux *http.ServeMux, h *handler.UserHandler) {
	userMux := http.NewServeMux()
	userMux.Handle("GET /", middleware.With(h.GetUsers))
	userMux.Handle("GET /{id}", middleware.With(h.GetUser))
	userMux.Handle("POST /", middleware.With(h.CreateUser))
	userMux.Handle("PUT /{id}", middleware.With(h.UpdateUser))
	userMux.Handle("DELETE /{id}", middleware.With(h.DeleteUser))

	mux.Handle("/users/", http.StripPrefix("/users", userMux))
}
