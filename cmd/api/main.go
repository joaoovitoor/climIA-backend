package main

import (
	"log"

	"climia-backend/config"
	appconfig "climia-backend/internal/config"
	"climia-backend/internal/database"
	"climia-backend/internal/handlers"
	"climia-backend/internal/server"
	"climia-backend/internal/services"
)

func main() {
	dbConfig := config.LoadConfig()
	appConfig := appconfig.LoadAppConfig()

	db := database.NewConnection(dbConfig)

	weatherRepo := database.NewWeatherRepository(db)
	weatherService := services.NewWeatherService(weatherRepo)
	weatherHandler := handlers.NewWeatherHandler(weatherService)

	srv := server.NewServer(appConfig)
	srv.Setup(weatherHandler)

	log.Fatal(srv.Start())
} 