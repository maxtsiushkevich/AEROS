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
	"users/internal/domain/cache"
	"users/internal/domain/producer"
	"users/internal/domain/repository"
	"users/internal/infrastructure/kafka"
	"users/internal/infrastructure/service"
	"users/internal/transport/handler"
)

type applicationProducers struct {
	producer producer.Publisher
}

type Server struct {
	server       *http.Server
	config       *config.Config
	grpcConn     *grpc.ClientConn
	logger       *slog.Logger
	rbacService  rbac.AuthorizationService
	userHandler  *handler.UserHandler
	storage      repository.UserStorage
	appProducers applicationProducers
	otpCache     cache.Cache
}

func NewServer(cfg *config.Config, logger *slog.Logger, rbacService rbac.AuthorizationService, db repository.UserStorage, otpCache cache.Cache) *Server {
	return &Server{
		config:       cfg,
		logger:       logger,
		rbacService:  rbacService,
		storage:      db,
		appProducers: applicationProducers{},
		otpCache:     otpCache,
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

	s.appProducers.producer = kafka.NewKafkaProducer(s.config.Kafka.Brokers, s.config.Kafka.Topic)

	service := service.NewUsersService(s.logger, s.grpcConn, s.storage, s.appProducers.producer, s.otpCache)

	s.userHandler = handler.NewUserHandler(service)

	router := gin.Default()

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

	protected.GET("", s.userHandler.Profile)
	protected.POST("/activate", s.userHandler.ActivateUserByCode)       // endpoint for self-activation with code
	protected.POST("/:userId/activate", s.userHandler.AdminUserByAdmin) // endpoint for admin activation
	protected.POST("/request-activation", s.userHandler.RequestActivationCode)
}
