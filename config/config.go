package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
	ENV  string

	DBHost         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBPort         string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBMaxLifetime  int // in minutes

	SlowQueryThreshold int // in milliseconds
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
		PORT:           port,
		ENV:            env,
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "go_web_service"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 100),
		DBMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxLifetime:  getEnvInt("DB_CONN_MAX_LIFETIME", 60),

		SlowQueryThreshold: getEnvInt("SLOW_QUERY_THRESHOLD", 200),
	}, nil
}

func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		var i int
		fmt.Sscanf(value, "%d", &i)
		return i
	}
	return defaultVal
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
