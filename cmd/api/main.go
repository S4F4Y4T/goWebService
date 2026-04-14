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
	"github.com/S4F4Y4T/goWebService/internal/product"
	"github.com/S4F4Y4T/goWebService/internal/router"
	"github.com/S4F4Y4T/goWebService/internal/shared/domain"
	"github.com/S4F4Y4T/goWebService/internal/shared/event"
	"github.com/S4F4Y4T/goWebService/internal/user"
	"github.com/S4F4Y4T/goWebService/pkg/telemetry"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Initialize Telemetry
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
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

	// ── Infrastructure ───────────────────────────────────────────────────

	dispatcher := event.NewDispatcher()

	// Register Domain Event Handlers
	dispatcher.Subscribe(user.UserCreatedTopic, func(ctx context.Context, ev domain.DomainEvent) error {
		if e, ok := ev.(user.UserCreated); ok {
			slog.Info("[EVENT] User Created",
				"userID", e.UserID,
				"email", e.Email,
				"occurredAt", e.OccurredAt())
		}
		return nil
	})

	dispatcher.Subscribe(product.ProductCreatedTopic, func(ctx context.Context, ev domain.DomainEvent) error {
		if e, ok := ev.(product.ProductCreated); ok {
			slog.Info("[EVENT] Product Created",
				"productID", e.ProductID,
				"name", e.Name,
				"occurredAt", e.OccurredAt())
		}
		return nil
	})

	// ── Dependency Wiring (DDD Modules) ───────────────────────────────────
	
	// User Context
	userRepo := user.NewUserRepository(db)
	userService := user.NewService(userRepo, dispatcher)
	userHandler := user.NewHandler(userService)

	// Product Context
	productRepo := product.NewProductRepository(db)
	productService := product.NewService(productRepo, dispatcher)
	productHandler := product.NewHandler(productService)

	// Unified Router
	r := router.NewRouter(userHandler, productHandler)
	mux := r.Setup()

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
