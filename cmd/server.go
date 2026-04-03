package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello World")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("POST /", postHandler)
	mux.HandleFunc("PUT /{id}", putHandler)
	mux.Handle("GET /health", logger(auth(http.HandlerFunc(healthHandler))))

	srv := &http.Server{
		Addr:    ":6969",
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Error starting server: ", err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home handler")
	w.Write([]byte("Hello World"))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post handler")
	w.Write([]byte("Post handler"))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Put Handler")
	id := r.PathValue("id")
	w.Write([]byte("Put Handler " + id))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Health handler")
	w.Write([]byte("OK"))
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logger Handler: ", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth handler")
		next.ServeHTTP(w, r)
	})
}
