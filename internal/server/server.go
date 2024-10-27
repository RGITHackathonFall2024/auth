package server

import (
	"log/slog"
	"reflect"

	"github.com/RGITHackathonFall2024/auth/pkg/ctxlog"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Server struct {
	hostname string
	port     uint

	db     *gorm.DB
	logger *slog.Logger

	app *fiber.App
}

func New(hostname string, port uint, db *gorm.DB, logger *slog.Logger) *Server {
	server := Server{
		hostname: hostname,
		port:     port,

		db:     db,
		logger: logger,

		app: fiber.New(),
	}

	DB = server.db

	return &server
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) Hostname() string {
	return s.hostname
}

func (s *Server) Port() uint {
	return s.port
}

func (s *Server) Log() *slog.Logger {
	return s.logger
}

func (s *Server) DB() *gorm.DB {
	if s.db == nil {
		return DB
	}

	return s.db
}

func FromContext(ctx *fiber.Ctx) *Server {
	log := ctxlog.WithCtx(slog.Default(), ctx)

	iServer := ctx.Locals("server")
	if iServer == nil {
		log.Error("No local server in context")
		return nil
	}

	server, ok := iServer.(*Server)
	if !ok {
		log.Error("Invalid type of the local server in context",
			slog.String("got_type", reflect.TypeOf(iServer).String()),
		)
		return nil
	}

	if server.db == nil {
		log.Error("No database connection in context")
		return nil
	}

	return server
}
