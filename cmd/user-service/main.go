package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/RGITHackathonFall2024/auth/internal/consts"
	grpcserver "github.com/RGITHackathonFall2024/auth/internal/grpc-server"
	"github.com/RGITHackathonFall2024/auth/internal/handlers"
	initdb "github.com/RGITHackathonFall2024/auth/internal/init-db"
	"github.com/RGITHackathonFall2024/auth/internal/server"
	"github.com/RGITHackathonFall2024/auth/pkg/ctxlog"
	"github.com/RGITHackathonFall2024/auth/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

func startServer(s *server.Server) error {
	log := s.Log().WithGroup("start-server")

	log.Info("Setting up middleware")
	s.App().Use(func(c *fiber.Ctx) error {
		c.Locals("server", s)
		return c.Next()
	})
	s.App().Use(requestid.New())
	s.App().Use(func(c *fiber.Ctx) error {
		s := server.FromContext(c)
		if s == nil {
			return fiber.ErrInternalServerError
		}

		log := ctxlog.WithCtx(s.Log(), c).WithGroup("request-middleware")

		log.Info("Request", slog.String("method", c.Method()), slog.String("path", c.Path()))

		return c.Next()
	})
	s.App().Use(cors.New())

	handlers.Setup(s)

	log.Info("Starting server")
	go func() {
		if err := s.App().Listen(fmt.Sprintf("%v:%v", s.Hostname(), s.Port())); err != nil {
			log.Error("Error running server", slog.String("err", err.Error()))
		}
	}()

	grpcServer := grpcserver.From(s)
	log.Info("Starting gRPC server")
	if err := grpcServer.Start(); err != nil {
		log.Error("Error running gRPC server", slog.String("err", err.Error()))
	}

	return nil
}

func main() {
	log := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level: slog.LevelDebug,
	}))
	log.Info("Loading dotenv file")
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
		panic(err)
	}

	log.Info("Connecting to database")
	db, err := initdb.Connect(log)
	if err != nil {
		panic(err)
	}

	if err := initdb.InitDB(db, log); err != nil {
		log.Error("Error initializing database")
		panic(err)
	}

	s := server.New(os.Getenv(consts.EnvHostname), uint(utils.Must(strconv.ParseUint(os.Getenv("PORT"), 10, 64))), nil, log)
	startServer(s)
}
