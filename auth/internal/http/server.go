package http

import (
	"auth/internal/config"
	"auth/internal/handlers"
	mdlwr "auth/internal/middleware"
	"auth/internal/storage"
	"auth/rbac"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/maxtsiushkevich/AEROS/pkg/middleware"
)

type Server struct {
	config    *config.Config
	router    *http.ServeMux
	logger    *slog.Logger
	validator *validator.Validate

	rbacService rbac.AuthorizationService

	auth *handlers.AuthHandler
}

func CreateServer(cfg *config.Config, rbac rbac.AuthorizationService) *Server {
	return &Server{
		config:      cfg,
		router:      http.NewServeMux(),
		validator:   validator.New(),
		rbacService: rbac,
	}
}

// Start web server. Init data storage and router
func (s *Server) Start() error {
	if s.logger == nil {
		s.logger = config.SetupLogger(s.config.Env)
	}

	var err error
	s.auth, err = handlers.NewAuthHandler(storage.CreateStorage(s.config, s.logger), s.logger)

	if err != nil {
		s.logger.Error("Failed to initialize auth handler", "err", err)
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
		mdlwr.AuthMiddleware(s.rbacService),
	}

	s.router.HandleFunc("POST /api/v1/auth/refresh", mw.Apply(s.auth.HandleRefreshTokens()))
	s.router.HandleFunc("POST /api/v1/auth/login", mw.Apply(s.auth.HandleLogin()))
	s.router.HandleFunc("POST /api/v1/auth/logout", mw.Apply(s.auth.HandleLogout()))
	s.router.HandleFunc("POST /api/v1/auth/change-password", mw.Apply(s.auth.HandleChangePassword()))

	s.logger.Info("Router configured")
}
