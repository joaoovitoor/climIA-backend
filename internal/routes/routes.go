package routes

import (
	"climia-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, weatherHandler *handlers.WeatherHandler) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "ClimIA API is running",
		})
	})

	api := app.Group("/", weatherHandler.AuthMiddleware)
	api.Get("/", weatherHandler.CalculateForecast)
} 