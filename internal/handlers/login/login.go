package login

import (
	"github.com/RGITHackathonFall2024/auth/internal/server"
	"github.com/gofiber/fiber/v2"
)

type Request struct {
	username string
	password string
}

func Setup(s *server.Server) {
	s.App().Post("/api/v1/login", func(c *fiber.Ctx) error {

	})
}
