package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
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
	if err == nil {
		return
	}

	msg := err.Error()

	// 404 Not Found mappings
	if msg == "user not found" || msg == "product not found" || msg == "record not found" {
		NotFound(w, msg)
		return
	}

	// 400 Bad Request mappings (JSON parsing issues, invalid IDs)
	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError
	var numErr *strconv.NumError
	
	if errors.As(err, &syntaxErr) || errors.As(err, &unmarshalErr) || errors.As(err, &numErr) || msg == "EOF" || msg == "unexpected EOF" {
		BadRequest(w, "invalid request data")
		return
	}

	// 500 Internal Server Error (Mask raw DB queries and sensitive errors)
	slog.Error("Internal Server Error", "error", err)
	JSON(w, http.StatusInternalServerError, Response{Success: false, Error: "internal server error"})
}
