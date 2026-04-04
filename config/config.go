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
}

func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("error loading PORT from .env file")
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	return &Config{
		PORT:       port,
		ENV:        env,
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "go_web_service"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
