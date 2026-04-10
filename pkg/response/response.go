package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/S4F4Y4T/goWebService/pkg/apperror"
	"gorm.io/gorm"
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

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		switch appErr.Category {
		case apperror.NotFound:
			NotFound(w, appErr.Message)
		case apperror.BadRequest:
			BadRequest(w, appErr.Message)
		case apperror.Conflict:
			BadRequest(w, appErr.Message) // Or specific Conflict response
		case apperror.Unauthorized:
			Unauthorized(w)
		default:
			slog.Error("Uncaught app error", "category", appErr.Category, "error", appErr.Err)
			internalError(w)
		}
		return
	}

	// Fallback for standard errors (e.g. JSON parsing, validation)
	msg := err.Error()

	// Handle GORM record not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		NotFound(w, "record not found")
		return
	}

	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError
	var numErr *strconv.NumError

	if errors.As(err, &syntaxErr) || errors.As(err, &unmarshalErr) || errors.As(err, &numErr) || msg == "EOF" {
		BadRequest(w, "invalid request data")
		return
	}

	// 500 Internal Server Error (Mask raw DB queries and sensitive errors)
	slog.Error("Internal Server Error", "error", err)
	internalError(w)
}

func internalError(w http.ResponseWriter) {
	JSON(w, http.StatusInternalServerError, Response{Success: false, Error: "internal server error"})
}
