// Package config handles environment-based configuration for the CareerManifest backend.
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values.
type Config struct {
	// Server
	Port    string
	GinMode string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// CORS
	AllowedOrigins string

	// LLM (optional)
	LLMProvider string
	LLMApiKey   string
	LLMModel    string

	// Admin
	AdminEmail    string
	AdminPassword string
}

// Load reads configuration from environment variables (with .env fallback).
func Load() (*Config, error) {
	// Load .env file if it exists (non-fatal if missing)
	_ = godotenv.Load()

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "72"))

	cfg := &Config{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "careermanifest"),

		JWTSecret:      getEnv("JWT_SECRET", "default-secret-change-me"),
		JWTExpiryHours: jwtExpiry,

		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),

		LLMProvider: getEnv("LLM_PROVIDER", ""),
		LLMApiKey:   getEnv("LLM_API_KEY", ""),
		LLMModel:    getEnv("LLM_MODEL", ""),

		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@careermanifest.in"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "Admin@123"),
	}

	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	return cfg, nil
}

// DSN returns the MySQL Data Source Name string.
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// IsLLMEnabled checks if an LLM provider is configured.
func (c *Config) IsLLMEnabled() bool {
	return c.LLMProvider != "" && c.LLMApiKey != ""
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
