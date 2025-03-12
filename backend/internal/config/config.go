package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	DB                 DBConfig
	JWTSecret          string
	MediaPath          string
	CorsAllowedOrigins string
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConfig creates a new Config
func NewConfig() *Config {
	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "animestreaming"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWTSecret:          getEnv("JWT_SECRET", "yoursecretkey"),
		MediaPath:          getEnv("MEDIA_PATH", "./media"),
		CorsAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
