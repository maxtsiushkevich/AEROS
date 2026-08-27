package grpc

import (
	authGrpc "auth/api/proto"
	"auth/internal/config"
	"auth/internal/storage"
	"auth/rbac"
	"context"
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Auth struct {
	authGrpc.UnimplementedAuthServer
	config      *config.Config
	storage     storage.AuthStorage
	logger      *slog.Logger
	rbacService rbac.AuthorizationService
}

func NewGRPCServer(cfg *config.Config, logger *slog.Logger, rbac rbac.AuthorizationService, db storage.AuthStorage) *Auth {
	return &Auth{
		config:      cfg,
		storage:     db,
		logger:      logger,
		rbacService: rbac,
	}
}

func StartGPRCServer(ctx context.Context, config *config.Config, logger *slog.Logger, rbac rbac.AuthorizationService, db storage.AuthStorage) {
	lis, err := net.Listen("tcp", config.GRPCServer.Address)
	if err != nil {
		log.Fatalf("Error creating port listener %s: %v", config.GRPCServer.Address, err)
	}

	auth := NewGRPCServer(config, logger, rbac, db)

	grpcServer := grpc.NewServer()

	authGrpc.RegisterAuthServer(grpcServer, auth)

	auth.logger.Info("Running a gRPC server on a %s\n", "addr", config.GRPCServer.Address)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}
}
