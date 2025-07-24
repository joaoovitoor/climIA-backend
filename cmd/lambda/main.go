package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"climia-backend/config"
	"climia-backend/internal/database"
	"climia-backend/internal/models"
	"climia-backend/internal/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LambdaHandler struct {
	weatherService *services.WeatherService
	config         *config.Config
}

func NewLambdaHandler() *LambdaHandler {
	dbConfig := config.LoadConfig()
	db := database.NewConnection(dbConfig)
	weatherRepo := database.NewWeatherRepository(db)
	weatherService := services.NewWeatherService(weatherRepo)

	return &LambdaHandler{
		weatherService: weatherService,
		config:         dbConfig,
	}
}

// validateAuth valida o Bearer Token
func (h *LambdaHandler) validateAuth(request events.APIGatewayProxyRequest) bool {
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		return false
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token == h.config.APIToken
}

func (h *LambdaHandler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Lambda request: %s %s", request.HTTPMethod, request.Path)

	// Headers CORS
	corsHeaders := map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-Requested-With",
		"Access-Control-Allow-Credentials": "true",
	}

	// Handle OPTIONS request (preflight)
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    corsHeaders,
			Body:       "",
		}, nil
	}

	if request.Path == "/health" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    corsHeaders,
			Body:       `{"status":"ok","message":"ClimIA API is running"}`,
		}, nil
	}

	if request.Path == "/" {
		// Valida autenticação (exceto para health check)
		if !h.validateAuth(request) {
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
				Headers:    corsHeaders,
				Body:       `{"error":"Invalid API token. Use: Bearer <token>"}`,
			}, nil
		}

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

		if req.Cidade == "" || req.Estado == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers:    corsHeaders,
				Body:       `{"error":"cidade e estado são obrigatórios"}`,
			}, nil
		}

		forecasts, err := h.weatherService.CalculateForecast(req)
		if err != nil {
			log.Printf("Erro ao calcular previsão: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers:    corsHeaders,
				Body:       fmt.Sprintf(`{"error":"%s"}`, err.Error()),
			}, nil
		}

		responseBody, _ := json.Marshal(forecasts)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    corsHeaders,
			Body:       string(responseBody),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Headers:    corsHeaders,
		Body:       `{"error":"Not found"}`,
	}, nil
}

func main() {
	handler := NewLambdaHandler()
	lambda.Start(handler.HandleRequest)
} 