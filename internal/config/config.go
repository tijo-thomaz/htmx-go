package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Port          string
	Env           string
	LogLevel      string
	DatabasePath  string
	SessionSecret    string
	SessionEncKey    string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file (ignore error if not exists)
	_ = godotenv.Load()

	return &Config{
		Port:          getEnv("PORT", "8080"),
		Env:           getEnv("ENV", "development"),
		LogLevel:      getEnv("LOG_LEVEL", "INFO"),
		DatabasePath:  getEnv("DATABASE_PATH", "./data/linkbio.db"),
		SessionSecret:    getEnv("SESSION_SECRET", "change-me-in-production"),
		SessionEncKey:    getEnv("SESSION_ENCRYPTION_KEY", ""),
	}, nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// getEnv retrieves env variable or returns fallback
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

