package config

import (
	"os"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	OllamaURL      string
	HuggingFaceKey string
}

func Load() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://user:password@localhost/studypartner?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		OllamaURL:      getEnv("OLLAMA_URL", "http://localhost:11434"),
		HuggingFaceKey: getEnv("HUGGINGFACE_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
