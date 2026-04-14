package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/S4F4Y4T/goWebService/pkg/response"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// UserService defines the interface for user application logic.
// Using an interface here allows handlers to be unit-tested with mocks.
type UserService interface {
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	Update(ctx context.Context, req *UpdateUserRequest) (*User, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindAll(ctx context.Context, limit, offset int) (*GetUsersResponse, error)
}

type Handler struct {
	srv UserService
}

func NewHandler(srv UserService) *Handler {
	return &Handler{srv: srv}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.srv.FindAll(r.Context(), 10, 0)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, users)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	u, err := h.srv.FindByID(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, u)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, "Validation failed: "+err.Error())
		return
	}

	u, err := h.srv.Create(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	// Fix: POST that creates a resource must return 201 Created
	response.Created(w, u)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}
	req.ID = uint(id)

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, "Validation failed: "+err.Error())
		return
	}

	u, err := h.srv.Update(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, u)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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
	response.Message(w, "Deleted User")
}
