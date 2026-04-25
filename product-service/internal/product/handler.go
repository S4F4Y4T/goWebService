package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/S4F4Y4T/goWebService/pkg/response"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) *Handler {
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
