package config

import (
	"os"
	"strings"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port               string
	CORSAllowedOrigins []string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "tictactoe"),
			Password: getEnv("DB_PASSWORD", "tictactoe123"),
			DBName:   getEnv("DB_NAME", "tictactoe"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		},
		Server: ServerConfig{
			Port:               getEnv("PORT", "8080"),
			CORSAllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost"), ","),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
