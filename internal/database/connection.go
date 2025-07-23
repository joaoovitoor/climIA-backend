package database

import (
	"database/sql"
	"log"

	"climia-backend/config"
	_ "github.com/go-sql-driver/mysql"
)

func NewConnection(config *config.Config) *sql.DB {
	db, err := sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Erro ao fazer ping no banco de dados:", err)
	}

	log.Println("Conectado ao banco de dados MySQL")
	return db
} 