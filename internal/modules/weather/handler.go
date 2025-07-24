package weather

import (
	"log"
	"strings"

	"climia-backend/configs"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
	config  *configs.Config
}

func NewHandler(service *Service, config *configs.Config) *Handler {
	return &Handler{
		service: service,
		config:  config,
	}
}

func (h *Handler) AuthMiddleware(c *fiber.Ctx) error {
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

func (h *Handler) CalculateForecast(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in CalculateForecast: %v", r)
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

	forecasts, err := h.service.CalculateForecast(req)
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