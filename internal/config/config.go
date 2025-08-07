package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for our application
type Config struct {
	Port         string
	DatabaseURL  string
	Environment  string
	LogLevel     string
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		DatabaseHost: getEnv("DB_HOST", "localhost"),
		DatabasePort: getEnv("DB_PORT", "5432"),
		DatabaseUser: getEnv("DB_USER", "library_user"),
		DatabasePass: getEnv("DB_PASSWORD", "library_pass"),
		DatabaseName: getEnv("DB_NAME", "library_db"),
	}

	// Build database URL if not provided directly
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.DatabaseURL = dbURL
	} else {
		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DatabaseUser,
			cfg.DatabasePass,
			cfg.DatabaseHost,
			cfg.DatabasePort,
			cfg.DatabaseName,
		)
	}

	return cfg, nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
