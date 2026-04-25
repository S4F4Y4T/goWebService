package auth

import (
	"encoding/json"
	"net/http"

	"github.com/S4F4Y4T/goWebService/pkg/response"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Register godoc
// POST /auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, "validation failed: "+err.Error())
		return
	}

	res, err := h.svc.Register(&req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, res)
}

// Login godoc
// POST /auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, "validation failed: "+err.Error())
		return
	}

	res, err := h.svc.Login(&req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, res)
}
