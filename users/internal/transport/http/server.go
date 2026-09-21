package http

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"pkg/middleware"
	"pkg/rbac"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"users/internal/config"
	"users/internal/domain/repository"
	"users/internal/infrastructure/service"
	"users/internal/transport/handler"
)

type Server struct {
	// router      *gin.Engine
	server      *http.Server
	config      *config.Config
	grpcConn    *grpc.ClientConn
	logger      *slog.Logger
	rbacService rbac.AuthorizationService
	userHandler *handler.UserHandler
	storage     repository.UserStorage
}

func NewServer(cfg *config.Config, logger *slog.Logger, rbacService rbac.AuthorizationService, db repository.UserStorage) *Server {
	return &Server{
		config:      cfg,
		logger:      logger,
		rbacService: rbacService,
		storage:     db,
	}
}

func (s *Server) Start() error {
	authServerAddr := s.config.AuthServer.Address

	// should use secure connection
	conn, err := grpc.NewClient(authServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating gRPC client: %v", err)
		return err
	}

	s.grpcConn = conn

	service := service.NewUsersService(s.grpcConn, s.storage)

	s.userHandler = handler.NewUserHandler(service)

	router := gin.Default()
	// router.Use(middleware.GinAuthMiddleware(s.rbacService))

	s.configRoutes(router)

	s.server = &http.Server{
		Addr:        s.config.HTTPServer.Address,
		ReadTimeout: s.config.HTTPServer.Timeout,
		IdleTimeout: s.config.HTTPServer.IdleTimeout,
		Handler:     router,
	}

	go func() {
		s.logger.Info("Starting HTTP server", "addr", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP server error", "error", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	err = s.grpcConn.Close()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	s.logger.Info("Finished graceful shutdown for the HTTP server")
	s.logger.Info("gRPC client connection closed")

	return nil
}

func (s *Server) configRoutes(router *gin.Engine) {
	g := router.Group("/api/v1/users")

	g.POST("/registration", s.userHandler.Registration)

	protected := g.Group("")
	protected.Use(middleware.GinAuthMiddleware(s.rbacService))

	protected.GET("/", s.userHandler.Profile)
	protected.POST("/:userId/activate", s.userHandler.Activate)
}
