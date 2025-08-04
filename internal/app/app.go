package app

import (
	"log"

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
	weatherMySQLRepo := weather.NewMySQLRepository(db)
	weatherService := weather.NewService(weatherMySQLRepo)
	weatherHandler := weather.NewHandler(weatherService, appConfig)

	var dynamoHandler *weather.DynamoDBHandler
	dynamoRepo, err := weather.NewDynamoDBRepository(appConfig)
	if err != nil {
		log.Printf("Erro ao inicializar DynamoDB: %v", err)
		// Se não conseguir inicializar o DynamoDB, não criar o handler
		dynamoHandler = nil
	} else {
		dynamoService := weather.NewDynamoDBService(dynamoRepo)
		dynamoHandler = weather.NewDynamoDBHandler(dynamoService, appConfig)
	}

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

	setupRoutes(app, weatherHandler, dynamoHandler)

	return &App{
		FiberApp: app,
	}
}

func setupRoutes(app *fiber.App, weatherHandler *weather.Handler, dynamoHandler *weather.DynamoDBHandler) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "ClimIA API is running",
		})
	})

	api := app.Group("/", dynamoHandler.AuthMiddleware)
	api.Get("/", dynamoHandler.GetProcessedForecast)

	mysqlApi := app.Group("/mysql", weatherHandler.AuthMiddleware)
	mysqlApi.Get("/", weatherHandler.CalculateForecast)
}
