package user

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.GetUsers)
	mux.HandleFunc("GET /by-email", h.GetUserByEmail)
	mux.HandleFunc("GET /{id}", h.GetUser)
	mux.HandleFunc("POST /", h.CreateUser)
	mux.HandleFunc("PUT /{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /{id}", h.DeleteUser)
}
