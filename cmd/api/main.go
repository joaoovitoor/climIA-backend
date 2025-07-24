package main

import (
	"log"

	"climia-backend/config"
	"climia-backend/internal/database"
	"climia-backend/internal/handlers"
	"climia-backend/internal/server"
	"climia-backend/internal/services"
)

func main() {
	appConfig := config.LoadConfig()

	db := database.NewConnection(appConfig)
	weatherRepo := database.NewWeatherRepository(db)
	weatherService := services.NewWeatherService(weatherRepo)
	weatherHandler := handlers.NewWeatherHandler(weatherService, appConfig)

	server := server.NewServer(appConfig)
	server.Setup(weatherHandler)

	log.Fatal(server.Start())
} 