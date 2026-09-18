package http

import (
	"context"
	"errors"
	"flights/internal/config"
	"flights/internal/handlers"
	"flights/internal/storage"
	"log/slog"
	"net"
	"net/http"

	"pkg/middleware"
)

type Server struct {
	config  *config.Config
	server  *http.Server
	logger  *slog.Logger
	storage storage.FlightsStorage

	flights *handlers.FlightHandler
}

func CreateServer(cfg *config.Config, l *slog.Logger, db storage.FlightsStorage) *Server {
	return &Server{
		config:  cfg,
		logger:  l,
		storage: db,
	}
}

// Start web server. Init data storage and router
func (s *Server) Start() error {
	if s.logger == nil {
		s.logger = config.SetupLogger(s.config.Env)
	}

	var err error
	s.flights, err = handlers.NewFlightHandler(s.storage)

	if err != nil {
		s.logger.Error("Connection to DB failed", "err", err)
		return err
	}

	s.logger.Info("Start server", "env", s.config.Env)
	s.logger.Debug("Serve on", "addr", "http://"+s.config.HTTPServer.Address)

	router := http.NewServeMux()
	s.configureRouter(router)

	s.server = &http.Server{
		Addr:    s.config.HTTPServer.Address,
		Handler: router,
	}

	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	go func() {
		if err := s.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Failed to start HTTP server", "err", err)
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
	s.logger.Info("Finished graceful shutdown for the HTTP server")
	return nil
}

func (s *Server) configureRouter(router *http.ServeMux) {
	mw := middleware.MiddlewareGroup{
		middleware.LoggingMiddleware(s.logger),
	}

	router.HandleFunc("GET /api/v1/flights", mw.Apply(s.flights.HandleGetFlights()))
	router.HandleFunc("POST /api/v1/flights", mw.Apply(s.flights.HandleCreateFlight()))
	router.HandleFunc("PATCH /api/v1/flights", mw.Apply(s.flights.HandlePatchFlight()))
	router.HandleFunc("DELETE /api/v1/flights", mw.Apply(s.flights.HandleDeleteFlight()))

	s.logger.Info("Router configured")
}
