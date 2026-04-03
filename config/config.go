package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
	ENV  string
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
		return nil, fmt.Errorf("error loading ENVIRONMENT from .env file")
	}

	return &Config{
		PORT: port,
		ENV:  env,
	}, nil
}
