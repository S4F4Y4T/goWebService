package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/S4F4Y4T/goWebService/config"
	"github.com/S4F4Y4T/goWebService/internal/bootstrap"
	"github.com/S4F4Y4T/goWebService/internal/product"
	"github.com/S4F4Y4T/goWebService/internal/router"
	"github.com/S4F4Y4T/goWebService/internal/shared/event"
	"github.com/S4F4Y4T/goWebService/internal/user"
	"github.com/S4F4Y4T/goWebService/pkg/telemetry"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// ── Telemetry ────────────────────────────────────────────────────────
	shutdownTracer, err := telemetry.InitTracer()
	if err != nil {
		slog.Error("Failed to initialize tracer", "error", err)
	}
	defer func() {
		if shutdownTracer != nil {
			shutdownTracer(context.Background())
		}
	}()

	shutdownMetrics, err := telemetry.InitMetrics()
	if err != nil {
		slog.Error("Failed to initialize metrics", "error", err)
	}
	defer func() {
		if shutdownMetrics != nil {
			shutdownMetrics(context.Background())
		}
	}()

	// ── Logging ──────────────────────────────────────────────────────────
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// ── Config ───────────────────────────────────────────────────────────
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		panic("impossible application state: failed to load configuration")
	}

	// ── Database ─────────────────────────────────────────────────────────
	db, err := config.InitDB(cfg)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		panic("impossible application state: database connection failed")
	}

	// ── Migrations ───────────────────────────────────────────────────────
	m, err := migrate.New("file://db/migrations", config.GetDBURL(cfg))
	if err != nil {
		slog.Error("Error initializing migrations", "error", err)
		panic("impossible application state: migration initialization failed")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Error migrating database", "error", err)
		panic("impossible application state: database migration failed")
	}

	// ── Events ───────────────────────────────────────────────────────────
	dispatcher := event.NewDispatcher()
	bootstrap.RegisterEventHandlers(dispatcher)

	// ── Dependency Wiring (DDD Modules) ──────────────────────────────────

	// User Context
	userRepo := user.NewUserRepository(db)
	userService := user.NewService(userRepo, dispatcher)
	userHandler := user.NewHandler(userService)

	// Product Context
	productRepo := product.NewProductRepository(db)
	productService := product.NewService(productRepo, dispatcher)
	productHandler := product.NewHandler(productService)

	// ── HTTP Server ───────────────────────────────────────────────────────
	r := router.NewRouter(userHandler, productHandler)
	mux := r.Setup()

	srv := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Server starting", "port", cfg.PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error starting server", "error", err)
			panic("impossible application state: server failed to start")
		}
	}()

	// ── Graceful Shutdown ─────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	slog.Info("Server exiting cleanly")
}
