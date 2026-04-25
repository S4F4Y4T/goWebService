package product

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.GetProducts)
	mux.HandleFunc("GET /{id}", h.GetProduct)
	mux.HandleFunc("POST /", h.CreateProduct)
	mux.HandleFunc("PUT /{id}", h.UpdateProduct)
	mux.HandleFunc("DELETE /{id}", h.DeleteProduct)
}
