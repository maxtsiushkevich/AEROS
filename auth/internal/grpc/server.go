package grpc

import (
	authGrpc "auth/api/proto"
	"auth/internal/config"
	"auth/internal/storage"
	"log/slog"
	"net"
	"pkg/rbac"

	"google.golang.org/grpc"
)

type Auth struct {
	authGrpc.UnimplementedAuthServer
	config      *config.Config
	storage     storage.AuthStorage
	logger      *slog.Logger
	rbacService rbac.AuthorizationService
}

func NewGRPCServer(cfg *config.Config, logger *slog.Logger, db storage.AuthStorage, rbacService rbac.AuthorizationService) *Auth {
	return &Auth{
		config:      cfg,
		storage:     db,
		logger:      logger,
		rbacService: rbacService,
	}
}

func StartGPRCServer(config *config.Config, logger *slog.Logger, db storage.AuthStorage, rbacService rbac.AuthorizationService) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", config.GRPCServer.Address)
	if err != nil {
		return nil, err
	}

	auth := NewGRPCServer(config, logger, db, rbacService)

	grpcServer := grpc.NewServer()

	authGrpc.RegisterAuthServer(grpcServer, auth)

	auth.logger.Info("Running a gRPC server on a %s\n", "addr", config.GRPCServer.Address)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC server stopped", "err", err)
		}
	}()

	return grpcServer, nil
}
