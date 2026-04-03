package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, res Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, Response{Success: true, Data: data})
}

func Message(w http.ResponseWriter, message string) {
	JSON(w, http.StatusOK, Response{Success: true, Message: message})
}

func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, Response{Success: true, Data: data})
}

func BadRequest(w http.ResponseWriter, message string) {
	JSON(w, http.StatusBadRequest, Response{Success: false, Error: message})
}

func Unauthorized(w http.ResponseWriter) {
	JSON(w, http.StatusUnauthorized, Response{Success: false, Error: "unauthorized"})
}

func NotFound(w http.ResponseWriter, message string) {
	JSON(w, http.StatusNotFound, Response{Success: false, Error: message})
}

func Error(w http.ResponseWriter, err error) {
	JSON(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
}
