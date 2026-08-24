package grpc

import (
	auth "auth/api/proto"
	"auth/internal/config"
	"auth/internal/storage"
	"context"
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type AuthGRPCServer struct {
	auth.UnimplementedAuthServer
	config  *config.Config
	storage storage.AuthStorage
	logger  *slog.Logger
}

func NewGRPCServer(cfg *config.Config) *AuthGRPCServer {

	logger := config.SetupLogger(cfg.Env)

	return &AuthGRPCServer{
		config:  cfg,
		logger:  logger,
		storage: storage.CreateStorage(cfg, logger),
	}
}

func StartGPRCServer(ctx context.Context, config *config.Config) {
	lis, err := net.Listen("tcp", config.GRPCServer.Address)
	if err != nil {
		log.Fatalf("Error creating port listener %s: %v", config.GRPCServer.Address, err)
	}

	authServer := NewGRPCServer(config)
	grpcServer := grpc.NewServer()

	auth.RegisterAuthServer(grpcServer, authServer)

	log.Printf("Running a gRPC server on a %s\n", config.GRPCServer.Address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("Error starting gRPC server: %v", err)
	}
}

func (s *AuthGRPCServer) AddUser(ctx context.Context, req *auth.AddUserRequest) (*auth.AddUserResponse, error) {
	// s.users[req.Id] = req.HashedPassword
	// fmt.Println(s.users)

	return &auth.AddUserResponse{
		AccessToken:  "access_token",
		RefreshToken: "refres_token",
	}, nil
}
