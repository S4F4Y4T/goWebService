package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/S4F4Y4T/goWebService/user-service/config"
	"github.com/S4F4Y4T/goWebService/user-service/internal/user"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := gorm.Open(postgres.Open(cfg.DBURL()), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	// ── Migrations ────────────────────────────────────────────────────────────
	m, err := migrate.New("file://db/migrations", cfg.DBURL())
	if err != nil {
		slog.Error("failed to init migrations", "error", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// ── Wiring ────────────────────────────────────────────────────────────────
	repo := user.NewUserRepository(db)
	svc := user.NewService(repo)
	handler := user.NewHandler(svc)

	mux := http.NewServeMux()
	userMux := http.NewServeMux()
	handler.RegisterRoutes(userMux)
	// Gateway forwards /users/* as /* to the service. But just in case, we also serve on /* directly.
	mux.Handle("/", userMux)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"user-service"}`))
	})

	srv := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("User service starting", "port", cfg.PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("User service failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("User service shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
}
