package main

import (
	"log"

	"climia-backend/configs"
	"climia-backend/internal/app"
)

func main() {
	appConfig := configs.LoadConfig()
	app := app.NewApp()

	log.Printf("Servidor iniciando na porta %s...", appConfig.Port)
	log.Fatal(app.FiberApp.Listen(":" + appConfig.Port))
} 