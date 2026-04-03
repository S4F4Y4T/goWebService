package middleware

import (
	"fmt"
	"net/http"
)

func Single(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Single Middleware")
		next.ServeHTTP(w, r)
	})
}
