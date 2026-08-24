package http

import (
	"flights/internal/config"
	"flights/internal/handlers"
	"flights/internal/storage"
	"flights/pkg/middleware"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Server struct {
	config    *config.Config
	router    *http.ServeMux
	logger    *slog.Logger
	validator *validator.Validate

	flights *handlers.FlightHandler
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

	var err error
	s.flights, err = handlers.NewFlightHandler(
		storage.CreateStorage(s.config, s.logger), s.validator)

	if err != nil {
		s.logger.Error("Connection to DB failed", "err", err)
		return err
	}

	s.configureRouter()

	s.logger.Info("Start server", "env", s.config.Env)
	s.logger.Debug("Serve on", "addr", "http://"+s.config.HTTPServer.Address)
	return http.ListenAndServe(s.config.HTTPServer.Address, s.router)
}

func (s *Server) configureRouter() {

	mw := middleware.MiddlewareGroup{
		middleware.LoggingMiddleware(s.logger),
	}

	s.router.HandleFunc("GET /api/v1/flights", mw.Apply(s.flights.HandleGetFlights()))
	s.router.HandleFunc("POST /api/v1/flights", mw.Apply(s.flights.HandleCreateFlight()))
	s.router.HandleFunc("PATCH /api/v1/flights", mw.Apply(s.flights.HandlePatchFlight()))
	s.router.HandleFunc("DELETE /api/v1/flights", mw.Apply(s.flights.HandleDeleteFlight()))

	s.logger.Info("Router configured")
}
