package httpd

import (
	"fmt"

	"github.com/aritradevelops/billbharat/backend/product/internal/core/jwtutil"
	"github.com/aritradevelops/billbharat/backend/product/internal/ports/httpd/handlers"
	"github.com/aritradevelops/billbharat/backend/shared/translation"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	host       string
	port       int
	app        *fiber.App
	handlers   *handlers.Handler
	jwtManager *jwtutil.JwtManager
}

func NewServer(host string, port int, handlers *handlers.Handler, jwtManager *jwtutil.JwtManager) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
	}))
	app.Use(recover.New(
		recover.Config{
			EnableStackTrace: true,
		},
	))
	app.Use(logger.New())
	app.Use(translation.New())
	server := &Server{
		host:       host,
		port:       port,
		app:        app,
		handlers:   handlers,
		jwtManager: jwtManager,
	}
	return server
}

func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf("%s:%d", s.host, s.port))
}
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
