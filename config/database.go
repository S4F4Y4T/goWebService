package config

import (
	"fmt"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetDBURL returns the standard postgres connection URL.
func GetDBURL(cfg *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
}

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := GetDBURL(cfg)

	// Use custom slog-based logger
	dbLogger := NewDBLogger(time.Duration(cfg.SlowQueryThreshold) * time.Millisecond)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, fmt.Errorf("failed to use otelgorm plugin: %w", err)
	}

	// Extract the underlying *sql.DB object to configure connection pools
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to extract sql.DB: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBMaxLifetime) * time.Minute)

	return db, nil
}
