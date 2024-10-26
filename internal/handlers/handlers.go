package handlers

import (
	pinghandler "github.com/RGITHackathonFall2024/auth/internal/handlers/ping"
	"github.com/RGITHackathonFall2024/auth/internal/server"
)

func Setup(s *server.Server) {
	pinghandler.Setup(s)
}
