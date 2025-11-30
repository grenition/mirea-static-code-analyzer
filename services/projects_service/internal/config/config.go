package config

import "os"

type Config struct {
	DatabaseURL     string
	JWTSecret      string
	AnalyzerBaseURL string
}

func Load() *Config {
	return &Config{
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/projects_db?sslmode=disable"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		AnalyzerBaseURL: getEnv("ANALYZER_BASE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

