package handler

import (
	"encoding/json"
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
	var req model.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}
	product, err := h.srv.Create(&req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	var req model.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}
	req.ID = uint(id)

	product, err := h.srv.Update(&req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.srv.Delete(&model.DeleteProductRequest{ID: uint(id)}); err != nil {
		response.Error(w, err)
		return
	}
	response.Message(w, "Deleted Product")
}
