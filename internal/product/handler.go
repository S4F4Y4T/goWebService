package product

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/S4F4Y4T/goWebService/pkg/response"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ProductService defines the interface for product application logic.
// Using an interface here allows handlers to be unit-tested with mocks.
type ProductService interface {
	Create(ctx context.Context, req *CreateProductRequest) (*Product, error)
	Update(ctx context.Context, req *UpdateProductRequest) (*Product, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Product, error)
	FindAll(ctx context.Context, limit, offset int) (*GetProductsResponse, error)
}

type Handler struct {
	srv ProductService
}

func NewHandler(srv ProductService) *Handler {
	return &Handler{srv: srv}
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.srv.FindAll(r.Context(), 10, 0)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, products)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	p, err := h.srv.FindByID(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, p)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, "Validation failed: "+err.Error())
		return
	}

	p, err := h.srv.Create(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	// Fix: POST that creates a resource must return 201 Created
	response.Created(w, p)
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}
	req.ID = uint(id)

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, "Validation failed: "+err.Error())
		return
	}

	p, err := h.srv.Update(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, p)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.srv.Delete(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.Message(w, "Deleted Product")
}
