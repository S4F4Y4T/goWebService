package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/S4F4Y4T/goWebService/pkg/response"
	"github.com/S4F4Y4T/goWebService/pkg/telemetry"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) *Handler {
	return &Handler{srv: srv}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	client := telemetry.NewHTTPClient()
	req, _ := http.NewRequestWithContext(r.Context(), "GET", "https://jsonplaceholder.typicode.com/users", nil)
	_, _ = client.Do(req)

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
	response.OK(w, u)
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
