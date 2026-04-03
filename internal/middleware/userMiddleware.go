package middleware

import (
	"fmt"
	"net/http"
)

func User(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("User Middleware")
		next.ServeHTTP(w, r)
	})
}
