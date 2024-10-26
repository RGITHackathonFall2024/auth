package mehandler

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/RGITHackathonFall2024/auth/internal/logic/auth"
	"github.com/RGITHackathonFall2024/auth/internal/logic/user"
	"github.com/RGITHackathonFall2024/auth/internal/server"
	"github.com/RGITHackathonFall2024/auth/pkg/ctxlog"
	"github.com/gofiber/fiber/v2"
)

type GetResponse struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	HomeTown   string `json:"home_town"`
	University string `json:"university"`
}

type PostRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	HomeTown   string `json:"home_town"`
	University string `json:"university"`
}

type PostResponse struct {
}

func Setup(s *server.Server) {
	s.App().Get("/api/v1/me", GetMe)
	s.App().Post("/api/v1/me", UpdateMe)
}

func GetMe(c *fiber.Ctx) error {
	s := server.FromContext(c)
	if s == nil {
		return fiber.ErrInternalServerError
	}

	log := ctxlog.WithCtx(s.Log(), c).WithGroup("get-me")

	log.Info("Getting user by token")
	usr, err := auth.GetUserByToken(log, c)
	if err != nil {
		if errors.Is(err, &auth.ErrInvalidToken{}) || errors.Is(err, &auth.ErrMissingToken{}) {
			return fiber.ErrUnauthorized
		}

		return fiber.ErrInternalServerError
	}
	fmt.Print(err)
	log.Debug("Got user", slog.Any("user", usr))

	return c.JSON(&GetResponse{
		FirstName:  usr.FirstName,
		LastName:   usr.LastName,
		HomeTown:   usr.HomeTown,
		University: usr.University,
	})
}

func UpdateMe(c *fiber.Ctx) error {
	s := server.FromContext(c)
	if s == nil {
		return fiber.ErrInternalServerError
	}

	log := ctxlog.WithCtx(s.Log(), c).WithGroup("update-me")

	var req PostRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error("Error parsing request body",
			slog.String("err", err.Error()),
			slog.String("body", string(c.Body())),
		)
		return fiber.ErrBadRequest
	}

	usr, err := auth.GetUserByToken(s.Log(), c)
	if err != nil {
		if errors.Is(err, &auth.ErrInvalidToken{}) || errors.Is(err, &auth.ErrMissingToken{}) {
			return fiber.ErrUnauthorized
		}

		return fiber.ErrInternalServerError
	}

	if req.FirstName != "" {
		usr.FirstName = req.FirstName
	}
	if req.LastName != "" {
		usr.LastName = req.LastName
	}
	if req.HomeTown != "" {
		usr.HomeTown = req.HomeTown
	}
	if req.University != "" {
		usr.University = req.University
	}

	if err = user.Edit(s, usr); err != nil {
		log.Error("Error editing user", slog.String("err", err.Error()))
		return fiber.ErrInternalServerError
	}

	return c.JSON(&PostResponse{})
}
