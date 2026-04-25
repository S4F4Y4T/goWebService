package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
	ENV  string

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string

	JWTSecret       string
	UserServiceURL  string
}

func Load() (*Config, error) {
	_ = godotenv.Load() // .env optional; env vars take precedence in Docker

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	return &Config{
		PORT:            port,
		ENV:             getEnv("ENVIRONMENT", "development"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "auth_db"),
		DBPort:          getEnv("DB_PORT", "5432"),
		JWTSecret:       getEnv("JWT_SECRET", "change-me-in-production"),
		UserServiceURL:  getEnv("USER_SERVICE_URL", "http://localhost:8082"),
	}, nil
}

func (c *Config) DBURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
