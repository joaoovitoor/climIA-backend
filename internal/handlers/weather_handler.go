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

	log.Printf("DEBUG - URL completa: %s", c.Context().URI().String())
	log.Printf("DEBUG - Query string: %s", c.Context().QueryArgs().String())
	log.Printf("DEBUG - Parâmetros recebidos: cidade=%s, estado=%s, data=%s, datainicio=%s, datafim=%s", 
		req.Cidade, req.Estado, req.Data, req.DataInicio, req.DataFim)

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
