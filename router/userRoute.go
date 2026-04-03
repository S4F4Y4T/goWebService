package router

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/internal/handler"
	"github.com/S4F4Y4T/goWebService/internal/middleware"
)

func RegisterUserRoutes(mux *http.ServeMux) {
	userMux := http.NewServeMux()
	userMux.Handle("GET /", middleware.With(handler.GetUsers))
	userMux.Handle("GET /{id}", middleware.With(handler.GetUser, middleware.Single))
	userMux.Handle("POST /", middleware.With(handler.CreateUser))
	userMux.Handle("PUT /{id}", middleware.With(handler.UpdateUser))
	userMux.Handle("DELETE /{id}", middleware.With(handler.DeleteUser))

	userMiddleware := middleware.Apply(middleware.User)

	mux.Handle("/users/", http.StripPrefix("/users", userMiddleware(userMux)))
}
