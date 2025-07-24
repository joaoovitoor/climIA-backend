package server

import (
	"log"

	"climia-backend/config"
	"climia-backend/internal/handlers"
	"climia-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	app    *fiber.App
	config *config.Config
}

func NewServer(appConfig *config.Config) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Erro na aplicação: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		},
	})

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "America/Sao_Paulo",
	}))

	return &Server{
		app:    app,
		config: appConfig,
	}
}

func (s *Server) Setup(weatherHandler *handlers.WeatherHandler) {
	routes.SetupRoutes(s.app, weatherHandler)
}

func (s *Server) Start() error {
	log.Printf("Servidor iniciando na porta %s...", s.config.Port)
	return s.app.Listen(":" + s.config.Port)
} 