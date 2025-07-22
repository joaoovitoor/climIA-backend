package server

import (
	"log"

	"climia-backend/internal/config"
	"climia-backend/internal/handlers"
	"climia-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app    *fiber.App
	config *config.AppConfig
}

func NewServer(appConfig *config.AppConfig) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Erro na aplicação: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		},
	})

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