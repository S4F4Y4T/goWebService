package handler

import (
	"net/http"
	"strconv"

	"github.com/S4F4Y4T/goWebService/internal/model"
	"github.com/S4F4Y4T/goWebService/internal/service"
	"github.com/S4F4Y4T/goWebService/pkg/response"
)

type ProductHandler struct {
	srv *service.ProductService
}

func NewProductHandler(srv *service.ProductService) *ProductHandler {
	return &ProductHandler{srv: srv}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {

	products, err := h.srv.FindAll()
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	product, err := h.srv.FindByID(&model.GetProductRequest{ID: uint(id)})
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, product)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	response.Message(w, "Create Product")
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	response.Message(w, "Update Product")
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	response.Message(w, "Delete Product")
}
