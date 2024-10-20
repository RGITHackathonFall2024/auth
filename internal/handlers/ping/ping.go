package pinghandler

import (
	"github.com/RGITHackathonFall2024/auth/internal/server"
	"github.com/gofiber/fiber/v2"
)

func Setup(s *server.Server) {
	s.App().Get("/api/v1/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong!")
	})
}
