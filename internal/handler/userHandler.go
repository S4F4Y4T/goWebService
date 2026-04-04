package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/S4F4Y4T/goWebService/internal/model"
	"github.com/S4F4Y4T/goWebService/internal/service"
	"github.com/S4F4Y4T/goWebService/pkg/response"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type UserHandler struct {
	srv *service.UserService
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	return &UserHandler{srv: srv}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.srv.FindAll(r.Context(), &model.GetUsersRequest{Limit: 10, Offset: 0})

	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	user, err := h.srv.FindByID(r.Context(), &model.GetUserRequest{ID: uint(id)})
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, "Validation failed: "+err.Error())
		return
	}
	user, err := h.srv.Create(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}
	req.ID = uint(id)

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, "Validation failed: "+err.Error())
		return
	}

	user, err := h.srv.Update(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.srv.Delete(r.Context(), &model.DeleteUserRequest{ID: uint(id)}); err != nil {
		response.Error(w, err)
		return
	}
	response.Message(w, "Deleted User")
}
