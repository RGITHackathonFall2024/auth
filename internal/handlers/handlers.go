package handlers

import (
	loginhandler "github.com/RGITHackathonFall2024/auth/internal/handlers/login"
	mehandler "github.com/RGITHackathonFall2024/auth/internal/handlers/me"
	pinghandler "github.com/RGITHackathonFall2024/auth/internal/handlers/ping"
	"github.com/RGITHackathonFall2024/auth/internal/server"
)

func Setup(s *server.Server) {
	pinghandler.Setup(s)
	loginhandler.Setup(s)
	mehandler.Setup(s)
}
