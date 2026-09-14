package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	auth "users/api/proto"

	"pkg/middleware"
	"pkg/rbac"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"users/internal/config"
)

type Server struct {
	router   *gin.Engine
	config   *config.Config
	grpcConn *grpc.ClientConn

	logger      *slog.Logger
	rbacService rbac.AuthorizationService
}

func NewServer(cfg *config.Config, logger *slog.Logger, rbacService rbac.AuthorizationService) *Server {
	return &Server{
		router:      gin.Default(),
		config:      cfg,
		logger:      logger,
		rbacService: rbacService,
	}
}

func (s *Server) ConfigServer() {
	authServerAddr := s.config.AuthServer.Address

	// should use secure connection
	conn, err := grpc.NewClient(authServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating gRPC client: %v", err)
	}
	s.grpcConn = conn

	s.router.Use(middleware.GinAuthMiddleware(s.rbacService))

	s.configRoutes()

	server := &http.Server{
		Addr:        s.config.HTTPServer.Address,
		ReadTimeout: s.config.HTTPServer.Timeout,
		IdleTimeout: s.config.HTTPServer.IdleTimeout,
		Handler:     s.router,
	}

	server.ListenAndServe()

}

func (s *Server) configRoutes() {
	s.router.GET("/ping", func(c *gin.Context) {
		client := auth.NewAuthClient(s.grpcConn)

		resp, err := client.AddUser(context.Background(), &auth.AddUserRequest{
			Id:       "00000000-0000-0000-0000-100000000000",
			Password: "34mf9304mf3940fj43jf34iksdz",
			Email:    "max@gmail.com",
		})

		if err != nil {
			c.String(http.StatusInternalServerError, "Error calling gRPC service: %v", err)
			return
		}

		// В resp будут refresh и access токены, которые нужно записать в куки
		fmt.Println(resp)

		c.String(http.StatusOK, "pong")
	})
}
