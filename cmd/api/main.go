package main

import (
	"log"
	"net/http"
	"time"

	"github.com/S4F4Y4T/goWebService/config"
	"github.com/S4F4Y4T/goWebService/internal/app"
	"github.com/S4F4Y4T/goWebService/internal/handler"
	"github.com/S4F4Y4T/goWebService/internal/repository"
	"github.com/S4F4Y4T/goWebService/internal/service"
	"github.com/S4F4Y4T/goWebService/router"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	// Run migrations
	migrationURL := "file://db/migrations"
	dbURL := config.GetDBURL(cfg)

	m, err := migrate.New(migrationURL, dbURL)
	if err != nil {
		log.Fatal("Error initializing migrations: ", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error migrating database: ", err)
	}

	// ── Dependency Wiring ───────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	productRepo := repository.NewProductRepository()
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	appInstance := &app.App{
		UserHandler:    userHandler,
		ProductHandler: productHandler,
	}
	mux := router.SetupRoutes(appInstance)

	srv := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", cfg.PORT)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Error starting server: ", err)
	}
}
