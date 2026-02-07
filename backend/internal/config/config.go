package config

import (
	"os"
	"strings"
)

type Config struct {
	Port           string
	Env            string
	DatabaseURL    string
	AllowedOrigins []string
	FirebaseProjectID string
	FirebaseCredentials string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		Env:            getEnv("ENV", "development"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", ""),
		FirebaseCredentials: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
