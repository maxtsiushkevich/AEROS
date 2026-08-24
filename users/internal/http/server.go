package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	auth "users/api/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"users/internal/config"
)

type Server struct {
	router   *gin.Engine
	config   *config.Config
	grpcConn *grpc.ClientConn
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		router: gin.Default(),
		config: cfg,
	}
}

func (s *Server) ConfigServer() {
	authServerAddr := s.config.AuthServer.Address
	conn, err := grpc.NewClient(authServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating gRPC client: %v", err)
	}
	s.grpcConn = conn

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
			Id:             "1234",
			HashedPassword: "34mf9304mf3940fj43jf34iksdz",
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
