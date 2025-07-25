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
		Port:               getEnv("PORT", ""),
		Env:                getEnv("ENV", ""),
		APIToken:           getEnv("API_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
