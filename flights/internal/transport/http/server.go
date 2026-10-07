package http

import (
	"context"
	"errors"
	"flights/internal/config"
	"flights/internal/domain/repository"
	infra "flights/internal/infrastructure/usecase"
	"flights/internal/transport/handlers"
	"log/slog"
	"net"
	"net/http"

	"pkg/middleware"
)

type Server struct {
	config         *config.Config
	server         *http.Server
	logger         *slog.Logger
	flightsStorage repository.FlightRepository
	cargoStorage   repository.CargoRepository

	flights *handlers.FlightHandler
	cargo   *handlers.CargoHandler
}

func CreateServer(cfg *config.Config, l *slog.Logger, flight_db repository.FlightRepository, cargo_db repository.CargoRepository) *Server {
	return &Server{
		config:         cfg,
		logger:         l,
		flightsStorage: flight_db,
		cargoStorage:   cargo_db,
	}
}

// Start web server. Init data storage and router
func (s *Server) Start() error {
	if s.logger == nil {
		s.logger = config.SetupLogger(s.config.Env)
	}

	var err error

	flightUseCase := infra.NewFlightUseCase(s.flightsStorage, s.cargoStorage)
	cargoUseCase := infra.NewCargoUseCase(s.cargoStorage)

	s.flights, err = handlers.NewFlightHandler(flightUseCase)
	s.cargo, err = handlers.NewCargoHandler(cargoUseCase)

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
	router.HandleFunc("DELETE /api/v1/flights/{id}", mw.Apply(s.flights.HandleDeleteFlight()))

	router.HandleFunc("POST /api/v1/flights/{id}/cancel", mw.Apply(s.flights.HandleCancelFlight()))
	router.HandleFunc("POST /api/v1/flights/{id}/reschedule", mw.Apply(s.flights.HandleRescheduleFlight()))
	router.HandleFunc("POST /api/v1/flights/{id}/status", mw.Apply(s.flights.HandleChangeStatus()))

	router.HandleFunc("GET /api/v1/cargo/{flight_id}", mw.Apply(s.cargo.HandleGetCargoManifest()))

	s.logger.Info("Router configured")
}
