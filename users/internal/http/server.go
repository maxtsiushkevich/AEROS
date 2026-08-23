package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	return &Server{
		router: gin.Default(),
	}
}

func (s *Server) ConfigServer() {
	s.configRoutes()
	server := &http.Server{
		Addr:    ":3003",
		Handler: s.router,
	}

	server.ListenAndServe()

}

func (s *Server) configRoutes() {
	s.router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
