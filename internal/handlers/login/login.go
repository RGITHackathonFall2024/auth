package loginhandler

import (
	"log/slog"
	"strconv"

	"github.com/RGITHackathonFall2024/auth/internal/logic/auth"
	"github.com/RGITHackathonFall2024/auth/internal/logic/user"
	"github.com/RGITHackathonFall2024/auth/internal/server"
	"github.com/RGITHackathonFall2024/auth/pkg/ctxlog"
	"github.com/gofiber/fiber/v2"
)

type Request struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoUrl  string `json:"photo_url"`
	AuthDate  uint64 `json:"auth_date"`
	Hash      string `json:"hash"`
}

type Response struct {
	Token string `json:"token"`
}

func Setup(s *server.Server) {
	s.App().Post("/api/v1/login", Login)
}

func Login(c *fiber.Ctx) error {
	s := server.FromContext(c)
	if s == nil {
		return fiber.ErrInternalServerError
	}

	log := ctxlog.WithCtx(s.Log(), c).WithGroup("login")

	log.Info("Parsing request body")
	var req Request
	if err := c.BodyParser(&req); err != nil {
		log.Error("Error parsing request body",
			slog.String("body", string(c.Body())),
			slog.String("err", err.Error()),
		)
		return fiber.ErrInternalServerError
	}

	log.Debug("Parsed request", slog.Any("request", req))

	log.Info("Verifying hash")
	if err := auth.VerifyHash(log,
		req.ID,
		req.FirstName,
		req.LastName,
		req.Username,
		req.PhotoUrl,
		req.AuthDate,
		req.Hash,
	); err != nil {
		if _, ok := err.(*auth.ErrInvalidHash); ok {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid hash")
		}

		return fiber.ErrInternalServerError
	}
	log.Info("Verified hash")

	log.Info("Trying to find user")
	_, err := user.ByID(s, req.ID)
	if err != nil {
		if _, ok := err.(*user.ErrNoSuchUser); ok {
			return user.Create(s, req.ID, req.FirstName, req.LastName)
		}

		log.Error("Error finding user",
			slog.Int64("id", req.ID),
			slog.String("err", err.Error()),
		)
		return fiber.ErrInternalServerError
	}
	log.Info("Found user")

	log.Info("Generating token")
	token, err := auth.GenerateToken(strconv.FormatInt(req.ID, 10))
	if err != nil {
		log.Error("Error generating token",
			slog.Int64("id", req.ID),
			slog.String("err", err.Error()),
		)
		return fiber.ErrInternalServerError
	}
	log.Debug("Generated token", slog.String("token", token))

	return c.Status(fiber.StatusOK).JSON(&Response{Token: token})
}
