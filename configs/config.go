package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnectionString string
	Port               string
	Env                string
	APIToken           string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}
	return &Config{
		DBConnectionString: getEnv("DB_CONNECTION_STRING", ""),
		Port:               getEnv("PORT", "8080"),
		Env:                getEnv("ENV", "development"),
		APIToken:           getEnv("API_TOKEN", "c58c5a0d964e9301df9a09900c3be55e6b03f78bef593dea650f55f357f206d4"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
