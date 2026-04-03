package main

import (
	"log"
	"net/http"
	"time"

	"github.com/S4F4Y4T/goWebService/config"
	"github.com/S4F4Y4T/goWebService/internal/router"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	mux := router.SetupRoutes()

	srv := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Error starting server: ", err)
	}

}
