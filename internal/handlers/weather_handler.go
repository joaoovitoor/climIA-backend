package handlers

import (
	"log"

	"climia-backend/internal/models"
	"climia-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type WeatherHandler struct {
	weatherService *services.WeatherService
}

func NewWeatherHandler(service *services.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: service,
	}
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
