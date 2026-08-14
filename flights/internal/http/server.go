package http

import (
	"flights/internal/config"
	"io"
	"log/slog"
	"net/http"
)

type Server struct {
	config *config.Config
	logger *slog.Logger
	router *http.ServeMux
}

func CreateServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		logger: nil,
		router: http.NewServeMux(),
	}
}

func (s *Server) Start() error {
	s.logger = config.SetupLogger(s.config.Env)
	s.logger.Info("Start server", "env", s.config.Env)

	s.configureRouter()

	return http.ListenAndServe(s.config.HTTPServer.Address, s.router)
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/hello", s.HandleHello())
}

func (s *Server) HandleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}
