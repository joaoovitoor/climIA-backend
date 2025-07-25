package database

import (
	"climia-backend/configs"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(config *configs.Config) *gorm.DB {
	dsn := config.DBConnectionString
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar com o banco de dados:", err)
	}

	fmt.Println("Conectado ao banco de dados PostgreSQL com sucesso!")
	return db
} 