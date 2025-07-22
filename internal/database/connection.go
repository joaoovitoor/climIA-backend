package database

import (
	"database/sql"
	"fmt"
	"log"

	"climia-backend/config"
	_ "github.com/go-sql-driver/mysql"
)

func NewConnection(config *config.Config) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Erro ao fazer ping no banco de dados:", err)
	}

	log.Println("Conectado ao banco de dados MySQL")
	return db
} 