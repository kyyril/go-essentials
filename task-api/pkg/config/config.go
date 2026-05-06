package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Env  string
	Port string
	Name string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig holds JWT token configuration
type JWTConfig struct {
	Secret             string
	ExpirationSeconds  int
	RefreshExpSeconds  int
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (non-blocking if not found)
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
			Name: getEnv("APP_NAME", "Task Management API"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "task_management_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "change-me-in-production"),
			ExpirationSeconds: getEnvInt("JWT_EXPIRATION", 3600),
			RefreshExpSeconds: getEnvInt("JWT_REFRESH_EXPIRATION", 604800),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseOrigins(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		},
	}

	return cfg, nil
}

// DSN returns the PostgreSQL connection string
func (c *Config) DSN() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// getEnv returns an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns an environment variable as an integer or a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// parseOrigins parses a comma-separated string of origins
func parseOrigins(originsStr string) []string {
	// Simple parser - can be enhanced for more complex cases
	if originsStr == "" {
		return []string{"http://localhost:3000"}
	}
	origins := make([]string, 0)
	start := 0
	for i := 0; i < len(originsStr); i++ {
		if originsStr[i] == ',' {
			origin := originsStr[start:i]
			if origin != "" {
				origins = append(origins, origin)
			}
			start = i + 1
		}
	}
	origin := originsStr[start:]
	if origin != "" {
		origins = append(origins, origin)
	}
	return origins
}
