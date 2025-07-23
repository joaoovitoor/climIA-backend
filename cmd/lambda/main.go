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

	app := fiber.New()
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
	fiberHandler := adaptor.FiberApp(h.app)

	log.Printf("=== LAMBDA INICIADA ===")
	log.Printf("DEBUG - HTTP Method: %s", request.HTTPMethod)
	log.Printf("DEBUG - Path: %s", request.Path)
	log.Printf("DEBUG - QueryStringParameters: %v", request.QueryStringParameters)
	log.Printf("DEBUG - Headers: %v", request.Headers)

	url := request.Path
	if len(request.QueryStringParameters) > 0 {
		params := make([]string, 0, len(request.QueryStringParameters))
		for key, value := range request.QueryStringParameters {
			params = append(params, key+"="+value)
		}
		url += "?" + strings.Join(params, "&")
		log.Printf("DEBUG - Final URL: %s", url)
	}

	httpReq, err := http.NewRequest(request.HTTPMethod, url, strings.NewReader(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Erro ao processar requisição"}`,
		}, err
	}

	log.Printf("DEBUG - HTTP Request URL: %s", httpReq.URL.String())
	log.Printf("DEBUG - HTTP Request Method: %s", httpReq.Method)

	for key, value := range request.Headers {
		httpReq.Header.Set(key, value)
	}

	response := &responseWriter{
		headers: make(map[string]string),
		body:    &strings.Builder{},
	}

	fiberHandler.ServeHTTP(response, httpReq)

	return events.APIGatewayProxyResponse{
		StatusCode:        response.statusCode,
		Headers:          response.headers,
		Body:             response.body.String(),
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