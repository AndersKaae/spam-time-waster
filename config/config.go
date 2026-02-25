package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration, primarily environment variables.
type Config struct {
	GeminiAPIKey         string
	GeminiModel          string
	GmailCredentialsFile string
}

// LoadConfig loads environment variables from a .env file into the Config struct.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found, relying on system environment variables.")

	}

	cfg := &Config{
		GeminiAPIKey:         mustGetenv("GEMINI_API_KEY"),
		GeminiModel:          mustGetenv("GEMINI_MODEL"),
		GmailCredentialsFile: mustGetenv("GMAIL_CREDENTIALS_FILE"),
	}

	return cfg, nil
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}
