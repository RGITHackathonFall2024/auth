package grpcserver

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/RGITHackathonFall2024/auth/internal/consts"
	grpc_user_service "github.com/RGITHackathonFall2024/auth/internal/grpc-user-service"
	"github.com/RGITHackathonFall2024/auth/internal/logic/auth"
	"github.com/RGITHackathonFall2024/auth/internal/logic/user"
	"github.com/RGITHackathonFall2024/auth/internal/server"
	"google.golang.org/grpc"
)

type Server struct {
	grpc_user_service.UnimplementedUserServer
	server.Server
}

func From(s *server.Server) *Server {
	return &Server{Server: *s}
}

func (s *Server) GetUser(ctx context.Context, req *grpc_user_service.GetUserRequest) (*grpc_user_service.UserData, error) {

	usr, err := user.ByID(&s.Server, req.UserID)
	if err != nil {
		return nil, err
	}

	userData := &grpc_user_service.UserData{
		TelegramID: usr.TelegramID,
		Username:   usr.Username,
		HomeTown:   usr.HomeTown,
		University: usr.University,
	}

	return userData, nil
}

func (s *Server) GetByToken(ctx context.Context, req *grpc_user_service.GetByTokenRequest) (*grpc_user_service.UserData, error) {
	usr, err := auth.GetUserByToken(s.Log(), &s.Server, req.Token)
	if err != nil {
		return nil, err
	}

	return &grpc_user_service.UserData{
		TelegramID: usr.TelegramID,
		Username:   usr.Username,
		HomeTown:   usr.HomeTown,
		University: usr.University,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *grpc_user_service.UpdateUserRequest) (*grpc_user_service.UpdateUserResponse, error) {
	if err := user.Edit(&s.Server, &user.User{
		TelegramID: req.TelegramID,
		Username:   req.Username,
		HomeTown:   req.HomeTown,
		University: req.University,
	}); err != nil {
		return nil, err
	}

	return &grpc_user_service.UpdateUserResponse{}, nil
}

func (s *Server) Start() error {
	grpcPort := os.Getenv(consts.EnvGrpcPort)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	grpc_user_service.RegisterUserServer(grpcServer, s)

	s.Log().Info("Starting gRPC server", slog.String("port", grpcPort))
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
