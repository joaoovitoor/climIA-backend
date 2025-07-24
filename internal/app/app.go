package app

import (
	"climia-backend/configs"
	"climia-backend/internal/modules/weather"
	"climia-backend/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type App struct {
	FiberApp *fiber.App
}

func NewApp() *App {
	appConfig := configs.LoadConfig()
	
	db := database.NewConnection(appConfig)
	weatherRepo := weather.NewRepository(db)
	weatherService := weather.NewService(weatherRepo)
	weatherHandler := weather.NewHandler(weatherService, appConfig)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "America/Sao_Paulo",
	}))

	setupRoutes(app, weatherHandler)

	return &App{
		FiberApp: app,
	}
}

func setupRoutes(app *fiber.App, weatherHandler *weather.Handler) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "ClimIA API is running",
		})
	})

	api := app.Group("/", weatherHandler.AuthMiddleware)
	api.Get("/", weatherHandler.CalculateForecast)
} 