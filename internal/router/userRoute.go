package router

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/internal/handler"
	"github.com/S4F4Y4T/goWebService/internal/middleware"
)

func RegisterUserRoutes(mux *http.ServeMux) {
	userMux := http.NewServeMux()
	userMux.Handle("GET /", http.HandlerFunc(handler.GetUsers))
	userMux.Handle("GET /{id}", middleware.With(handler.GetUser, middleware.Single))
	userMux.Handle("POST /", http.HandlerFunc(handler.CreateUser))
	userMux.Handle("PUT /{id}", http.HandlerFunc(handler.UpdateUser))
	userMux.Handle("DELETE /{id}", http.HandlerFunc(handler.DeleteUser))

	userMiddleware := middleware.Apply(middleware.User)

	mux.Handle("/users/", http.StripPrefix("/users", userMiddleware(userMux)))
}
