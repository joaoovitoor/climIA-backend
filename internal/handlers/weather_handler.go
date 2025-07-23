package handlers

import (
	"log"
	"strings"

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

	queryArgs := c.Context().QueryArgs()
	
	req.Cidade = strings.TrimSpace(string(queryArgs.Peek("cidade")))
	req.Estado = strings.ToUpper(strings.TrimSpace(string(queryArgs.Peek("estado"))))
	req.Data = string(queryArgs.Peek("data"))
	req.DataInicio = string(queryArgs.Peek("datainicio"))
	req.DataFim = string(queryArgs.Peek("datafim"))

	if req.Cidade == "" || req.Estado == "" {
		return c.JSON(fiber.Map{
			"message": "ClimIA API - Previsão Meteorológica",
			"version": "1.0.0",
			"status":  "running",
			"usage": "Adicione parâmetros: ?cidade=Guarulhos&estado=SP&data=2025-11-01",
			"example": "http://localhost:8080/?cidade=Guarulhos&estado=SP&data=2025-11-01",
		})
	}

	forecasts, err := h.weatherService.CalculateForecast(req)
	if err != nil {
		log.Printf("Erro ao calcular previsão: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(forecasts) == 1 {
		return c.JSON(forecasts[0])
	}

	return c.JSON(forecasts)
} 