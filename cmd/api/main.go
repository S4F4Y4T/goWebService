package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	os.MkdirAll("tmp", os.ModePerm)
	logFile, err := os.OpenFile("tmp/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("impossible application state: failed to create log file")
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logFile), nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		panic("impossible application state: failed to load configuration")
	}

	db, err := config.InitDB(cfg)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		panic("impossible application state: database connection failed")
	}

	// Run migrations
	migrationURL := "file://db/migrations"
	dbURL := config.GetDBURL(cfg)

	m, err := migrate.New(migrationURL, dbURL)
	if err != nil {
		slog.Error("Error initializing migrations", "error", err)
		panic("impossible application state: migration initialization failed")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Error migrating database", "error", err)
		panic("impossible application state: database migration failed")
	}

	// ── Dependency Wiring ───────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	productRepo := repository.NewProductRepository(db)
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

	// Run server in a goroutine so it doesn't block
	go func() {
		slog.Info("Server starting", "port", cfg.PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error starting server", "error", err)
			panic("impossible application state: server failed to start")
		}
	}()

	// Wait for OS interruption signals to safely bring down the architecture
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block the main thread until signal is received
	<-quit
	slog.Info("Shutting down server...")

	// Give active HTTP connections 5 seconds to finish their work
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	// Close database connections
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	slog.Info("Server exiting cleanly")
}
