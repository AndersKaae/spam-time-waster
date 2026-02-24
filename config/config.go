package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration, primarily environment variables.
type Config struct {
	GeminiAPIKey string
	GeminiModel  string
}

// LoadConfig loads environment variables from a .env file into the Config struct.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, relying on system environment variables.")
	}

	cfg := &Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		GeminiModel:  os.Getenv("GEMINI_MODEL"),
	}

	return cfg, nil
}
