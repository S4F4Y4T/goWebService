package product

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/internal/shared/middleware"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	productMux := http.NewServeMux()
	productMux.Handle("GET /", middleware.With(h.GetProducts))
	productMux.Handle("GET /{id}", middleware.With(h.GetProduct))
	productMux.Handle("POST /", middleware.With(h.CreateProduct))
	productMux.Handle("PUT /{id}", middleware.With(h.UpdateProduct))
	productMux.Handle("DELETE /{id}", middleware.With(h.DeleteProduct))

	mux.Handle("/products/", http.StripPrefix("/products", productMux))
}
