package main

import (
	"context"
	"log"

	"climia-backend/config"
	"climia-backend/internal/database"
	"climia-backend/internal/handlers"
	"climia-backend/internal/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

type LambdaHandler struct {
	app *fiber.App
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
		app: app,
	}
}

func (h *LambdaHandler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Lambda request: %s %s", request.HTTPMethod, request.Path)
	
	response, err := adaptor.FiberApp(h.app)(ctx, request)
	if err != nil {
		log.Printf("Erro no handler: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error":"Internal server error"}`,
		}, nil
	}
	
	return response, nil
}

func main() {
	handler := NewLambdaHandler()
	lambda.Start(handler.HandleRequest)
} 