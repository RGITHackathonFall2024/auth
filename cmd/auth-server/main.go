package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/RGITHackathonFall2024/auth/internal/handlers"
	"github.com/RGITHackathonFall2024/auth/internal/server"
	"github.com/RGITHackathonFall2024/auth/pkg/utils"
	"github.com/joho/godotenv"
)

func startServer(s *server.Server) {
	handlers.Setup(s)
	s.App().Listen(fmt.Sprintf("%v:%v", s.Hostname(), s.Port()))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	log := slog.Default()

	s := server.New(os.Getenv("HOSTNAME"), uint(utils.Must(strconv.ParseUint(os.Getenv("PORT"), 10, 64))), nil, log)
	startServer(s)
}
