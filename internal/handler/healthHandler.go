package handler

import (
	"net/http"

	"github.com/S4F4Y4T/goWebService/pkg/response"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response.Message(w, "Api is Running")
}
