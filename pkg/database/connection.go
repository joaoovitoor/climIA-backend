package database

import (
	"climia-backend/configs"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnection(config *configs.Config) *gorm.DB {
	dsn := config.DBConnectionString
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar com o banco de dados:", err)
	}

	fmt.Println("Conectado ao banco de dados com sucesso!")
	return db
} 