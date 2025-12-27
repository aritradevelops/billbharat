package httpd

import (
	"fmt"

	"github.com/aritradeveops/billbharat/backend/auth/internal/pkg/translation"
	"github.com/aritradeveops/billbharat/backend/auth/internal/ports/httpd/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	host     string
	port     int
	app      *fiber.App
	handlers *handlers.Handler
}

func NewServer(host string, port int, handlers *handlers.Handler) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})
	app.Use(recover.New(
		recover.Config{
			EnableStackTrace: true,
		},
	))
	app.Use(logger.New())
	app.Use(translation.New())
	server := &Server{
		host:     host,
		port:     port,
		app:      app,
		handlers: handlers,
	}
	return server
}

func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf("%s:%d", s.host, s.port))
}
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
