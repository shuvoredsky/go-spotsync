package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Dsn           string
	JwtSecret     string
	AllowedOrigin string
}

func LoadEnv() *Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	return &Config{
		Port:          port,
		Dsn:           os.Getenv("DSN"),
		JwtSecret:     os.Getenv("JWT_SECRET"),
		AllowedOrigin: allowedOrigin,
	}
}
