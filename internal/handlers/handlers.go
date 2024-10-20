package handlers

import (
	jwthandler "github.com/RGITHackathonFall2024/auth/internal/handlers/jwt"
	pinghandler "github.com/RGITHackathonFall2024/auth/internal/handlers/ping"
	"github.com/RGITHackathonFall2024/auth/internal/server"
)

func Setup(s *server.Server) {
	jwthandler.Setup(s)
	pinghandler.Setup(s)
}
