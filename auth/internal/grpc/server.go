package grpc

import (
	authGrpc "auth/api/proto"
	"auth/internal/config"
	"auth/internal/storage"
	"context"
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Auth struct {
	authGrpc.UnimplementedAuthServer
	config  *config.Config
	storage storage.AuthStorage
	logger  *slog.Logger
}

func NewGRPCServer(cfg *config.Config) *Auth {

	logger := config.SetupLogger(cfg.Env)

	return &Auth{
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

	auth := NewGRPCServer(config)

	if err := auth.storage.Open(); err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	grpcServer := grpc.NewServer()

	authGrpc.RegisterAuthServer(grpcServer, auth)

	auth.logger.Info("Running a gRPC server on a %s\n", "addr", config.GRPCServer.Address)

	defer func() {
		if err := auth.storage.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}
}
