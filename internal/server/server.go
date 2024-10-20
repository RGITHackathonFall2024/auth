package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	hostname string
	port     uint
	db       *gorm.DB
	app      *fiber.App
}

func New(hostname string, port uint, db *gorm.DB) *Server {
	server := Server{
		hostname: hostname,
		port:     port,
		db:       db,
		app:      fiber.New(),
	}
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
