package handlers

import (
	"log"
	"strings"

	"climia-backend/config"
	"climia-backend/internal/models"
	"climia-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type WeatherHandler struct {
	weatherService *services.WeatherService
	config         *config.Config
}

func NewWeatherHandler(service *services.WeatherService, config *config.Config) *WeatherHandler {
	return &WeatherHandler{
		weatherService: service,
		config:         config,
	}
}

// AuthMiddleware valida o Bearer Token
func (h *WeatherHandler) AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	// Verifica se é Bearer Token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization format. Use: Bearer <token>",
		})
	}

	// Extrai o token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != h.config.APIToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid API token",
		})
	}

	return c.Next()
}

func (h *WeatherHandler) CalculateForecast(c *fiber.Ctx) error {
	var req models.WeatherRequest
	c.QueryParser(&req)

	if req.Cidade == "" || req.Estado == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cidade e estado são obrigatórios",
		})
	}

	forecasts, err := h.weatherService.CalculateForecast(req)
	if err != nil {
		log.Printf("Erro ao calcular previsão: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(forecasts) == 1 {
		return c.JSON(forecasts[0])
	}

	return c.JSON(forecasts)
}
