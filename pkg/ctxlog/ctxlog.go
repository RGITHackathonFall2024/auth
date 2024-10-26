package ctxlog

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func WithCtx(logger *slog.Logger, ctx *fiber.Ctx) *slog.Logger {
	return logger.With(slog.String("req_id", ctx.Locals("requestid").(string)))
}
