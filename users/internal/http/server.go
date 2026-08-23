package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"users/internal/config"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		router: gin.Default(),
		config: cfg,
	}
}

func (s *Server) ConfigServer() {
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
		c.String(http.StatusOK, "pong")
	})
}
