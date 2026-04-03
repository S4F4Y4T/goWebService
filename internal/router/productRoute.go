package router

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/internal/handler"
)

func RegisterProductRoutes(mux *http.ServeMux) {
	productMux := http.NewServeMux()
	productMux.Handle("GET /", http.HandlerFunc(handler.GetProducts))
	productMux.Handle("GET /{id}", http.HandlerFunc(handler.GetProduct))
	productMux.Handle("POST /", http.HandlerFunc(handler.CreateProduct))
	productMux.Handle("PUT /{id}", http.HandlerFunc(handler.UpdateProduct))
	productMux.Handle("DELETE /{id}", http.HandlerFunc(handler.DeleteProduct))

	mux.Handle("/products/", http.StripPrefix("/products", productMux))
}
