package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/S4F4Y4T/goWebService/internal/model"
	"github.com/S4F4Y4T/goWebService/internal/service"
	"github.com/S4F4Y4T/goWebService/pkg/response"
)

type UserHandler struct {
	srv *service.UserService
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	return &UserHandler{srv: srv}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.srv.FindAll(&model.GetUsersRequest{Limit: 10, Offset: 0})

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

	user, err := h.srv.FindByID(&model.GetUserRequest{ID: uint(id)})
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
	user, err := h.srv.Create(&req)
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

	user, err := h.srv.Update(&req)
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

	if err := h.srv.Delete(&model.DeleteUserRequest{ID: uint(id)}); err != nil {
		response.Error(w, err)
		return
	}
	response.Message(w, "Deleted User")
}
