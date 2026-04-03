package router

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/internal/handler"
	"github.com/S4F4Y4T/goWebService/internal/middleware"
)

func RegisterProductRoutes(mux *http.ServeMux, h *handler.ProductHandler) {
	productMux := http.NewServeMux()
	productMux.Handle("GET /", middleware.With(h.GetProducts))
	productMux.Handle("GET /{id}", middleware.With(h.GetProduct))
	productMux.Handle("POST /", middleware.With(h.CreateProduct))
	productMux.Handle("PUT /{id}", middleware.With(h.UpdateProduct))
	productMux.Handle("DELETE /{id}", middleware.With(h.DeleteProduct))

	mux.Handle("/products/", http.StripPrefix("/products", productMux))
}
