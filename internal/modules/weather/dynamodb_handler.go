package weather

import (
	"log"
	"strings"

	"climia-backend/configs"

	"github.com/gofiber/fiber/v2"
)

type DynamoDBHandler struct {
	service *DynamoDBService
	config  *configs.Config
}

func NewDynamoDBHandler(service *DynamoDBService, config *configs.Config) *DynamoDBHandler {
	return &DynamoDBHandler{
		service: service,
		config:  config,
	}
}

func (h *DynamoDBHandler) AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization format. Use: Bearer <token>",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != h.config.APIToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid API token",
		})
	}

	return c.Next()
}

func (h *DynamoDBHandler) GetProcessedForecast(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in GetProcessedForecast: %v", r)
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		}
	}()

	var req WeatherRequest
	if err := c.QueryParser(&req); err != nil {
		log.Printf("Erro ao fazer parse dos parâmetros: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Parâmetros inválidos",
		})
	}

	forecasts, err := h.service.GetProcessedForecast(req)
	if err != nil {
		log.Printf("Erro ao buscar previsão processada: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(forecasts) == 1 {
		return c.JSON(forecasts[0])
	}

	return c.JSON(forecasts)
}
