package config

import (
	"os"
)

type AppConfig struct {
	Port string
	Env  string
}

func LoadAppConfig() *AppConfig {
	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "development")

	return &AppConfig{
		Port: port,
		Env:  env,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 