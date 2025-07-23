package main

import (
	"context"
	"log"
	"net/http"
	"strings"

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

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Erro na aplicação: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		},
	})

	// Configurar rotas
	app.Get("/", weatherHandler.CalculateForecast)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "ClimIA API is running",
		})
	})

	return &LambdaHandler{app: app}
}

func (h *LambdaHandler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log para debug
	log.Printf("DEBUG - HTTP Method: %s", request.HTTPMethod)
	log.Printf("DEBUG - Path: %s", request.Path)
	log.Printf("DEBUG - QueryStringParameters: %+v", request.QueryStringParameters)
	log.Printf("DEBUG - Headers: %+v", request.Headers)

	// Converter APIGatewayProxyRequest para Fiber Context
	fiberHandler := adaptor.FiberApp(h.app)

	// Criar HTTP request
	httpReq, err := http.NewRequest(request.HTTPMethod, request.Path, strings.NewReader(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Erro ao processar requisição"}`,
		}, err
	}

	// Adicionar headers
	for key, value := range request.Headers {
		httpReq.Header.Set(key, value)
	}

	// Adicionar query parameters
	q := httpReq.URL.Query()
	for key, value := range request.QueryStringParameters {
		q.Add(key, value)
	}
	httpReq.URL.RawQuery = q.Encode()

	log.Printf("DEBUG - Final URL: %s", httpReq.URL.String())

	// Criar response writer
	response := &responseWriter{
		headers: make(map[string]string),
		body:    &strings.Builder{},
	}

	// Processar requisição
	fiberHandler.ServeHTTP(response, httpReq)

	// Converter resposta
	responseBody := response.body.String()

	return events.APIGatewayProxyResponse{
		StatusCode:        response.statusCode,
		Headers:          response.headers,
		Body:             responseBody,
		IsBase64Encoded:  false,
	}, nil
}

type responseWriter struct {
	statusCode int
	headers    map[string]string
	body       *strings.Builder
}

func (w *responseWriter) Header() http.Header {
	return make(http.Header)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func main() {
	handler := NewLambdaHandler()
	lambda.Start(handler.HandleRequest)
} 