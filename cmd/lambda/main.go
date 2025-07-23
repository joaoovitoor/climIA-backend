package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"climia-backend/config"
	"climia-backend/internal/database"
	"climia-backend/internal/handlers"
	"climia-backend/internal/models"
	"climia-backend/internal/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
)

type LambdaHandler struct {
	app            *fiber.App
	weatherService *services.WeatherService
}

func NewLambdaHandler() *LambdaHandler {
	dbConfig := config.LoadConfig()
	db := database.NewConnection(dbConfig)
	weatherRepo := database.NewWeatherRepository(db)
	weatherService := services.NewWeatherService(weatherRepo)
	weatherHandler := handlers.NewWeatherHandler(weatherService)

	app := fiber.New()
	app.Get("/", weatherHandler.CalculateForecast)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "ClimIA API is running",
		})
	})

	return &LambdaHandler{
		app:            app,
		weatherService: weatherService,
	}
}

func (h *LambdaHandler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("=== LAMBDA INICIADA ===")
	log.Printf("DEBUG - HTTP Method: %s", request.HTTPMethod)
	log.Printf("DEBUG - Path: %s", request.Path)
	log.Printf("DEBUG - QueryStringParameters: %v", request.QueryStringParameters)

	if request.Path == "/health" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"status":"ok","message":"ClimIA API is running"}`,
		}, nil
	}

	if request.Path == "/" {
		var req models.WeatherRequest
		for key, value := range request.QueryStringParameters {
			switch key {
			case "cidade":
				req.Cidade = value
			case "estado":
				req.Estado = value
			case "data":
				req.Data = value
			case "datainicio":
				req.DataInicio = value
			case "datafim":
				req.DataFim = value
			}
		}

		log.Printf("DEBUG - Parâmetros extraídos: cidade=%s, estado=%s, data=%s", req.Cidade, req.Estado, req.Data)

		forecasts, err := h.weatherService.CalculateForecast(req)
		if err != nil {
			log.Printf("Erro ao calcular previsão: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       fmt.Sprintf(`{"error":"%s"}`, err.Error()),
			}, nil
		}

		responseBody, _ := json.Marshal(forecasts)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(responseBody),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"error":"Not found"}`,
	}, nil
}

func main() {
	handler := NewLambdaHandler()
	lambda.Start(handler.HandleRequest)
} 