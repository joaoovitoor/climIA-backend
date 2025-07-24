package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnectionString string
	Port              string
	Env               string
	APIToken          string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}
	return &Config{
		DBConnectionString: getEnv("DB_CONNECTION_STRING", "root:senha123@tcp(localhost:3306)/climia?parseTime=true"),
		Port:              getEnv("PORT", "8080"),
		Env:               getEnv("ENV", "development"),
		APIToken:          getEnv("API_TOKEN", "climia-secret-token-2025"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 