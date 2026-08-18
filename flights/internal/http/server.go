package http

import (
	"flights/internal/config"
	"flights/internal/service"
	"flights/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Server struct {
	config    *config.Config
	router    *http.ServeMux
	storage   storage.Storage
	logger    *slog.Logger
	service   *service.FlightService
	validator *validator.Validate
}

func CreateServer(cfg *config.Config) *Server {
	return &Server{
		config:    cfg,
		router:    http.NewServeMux(),
		validator: validator.New(),
	}
}

// Start web server. Init data storage and router
func (s *Server) Start() error {
	if s.logger == nil {
		s.logger = config.SetupLogger(s.config.Env)
	}

	err := s.configureStorage()
	if err != nil {
		return err
	}
	s.configureRouter()
	s.configureService()

	s.logger.Info("Start server", "env", s.config.Env)
	s.logger.Debug("Serve on", "addr", "http://"+s.config.HTTPServer.Address)
	return http.ListenAndServe(s.config.HTTPServer.Address, s.router)
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("GET /flights/", s.HandleGetFlights())
	s.router.HandleFunc("POST /flights/", s.HandleCreateFlight())

	s.logger.Info("Router configured")
}

func (s *Server) configureStorage() error {
	s.storage = storage.CreateStorage(s.config, s.logger)
	if err := s.storage.Open(); err != nil {
		s.logger.Error("Connection to DB failed", "err", err)
		return err
	}

	s.logger.Info("Init db", "env", s.config.Env)
	return nil
}

func (s *Server) configureService() {
	s.service = service.CreateFlightService(s.storage)
}
