package handler

import (
	"net/http"

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
	w.Write([]byte("Get User"))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create User"))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update User"))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete User"))
}
