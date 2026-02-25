package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration, primarily environment variables.
type Config struct {
	GeminiAPIKey      string
	GeminiModel       string
	GmailClientID     string
	GmailClientSecret string
	GmailTokenFile    string
	GmailLabel        string
}

// LoadConfig loads environment variables from a .env file into the Config struct.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, relying on system environment variables.")
	}

	cfg := &Config{
		GeminiAPIKey:      mustGetenv("GEMINI_API_KEY"),
		GeminiModel:       mustGetenv("GEMINI_MODEL"),
		GmailClientID:     mustGetenv("GMAIL_CLIENT_ID"),
		GmailClientSecret: mustGetenv("GMAIL_CLIENT_SECRET"),
		GmailTokenFile:    mustGetenv("GMAIL_TOKEN_FILE"),
		GmailLabel:        mustGetenv("GMAIL_LABEL"),
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
