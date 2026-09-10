package http

import (
	"auth/internal/cache"
	"auth/internal/config"
	"auth/internal/handlers"
	"auth/internal/storage"
	"log/slog"
	"net/http"
	"pkg/middleware"
	"pkg/rbac"

	"github.com/go-playground/validator/v10"
)

type Server struct {
	config      *config.Config
	router      *http.ServeMux
	logger      *slog.Logger
	validator   *validator.Validate
	auth        *handlers.AuthHandler
	storage     storage.AuthStorage
	cache       cache.RevokedTokenCache
	rbacService rbac.AuthorizationService
}

func CreateServer(cfg *config.Config, logger *slog.Logger, db storage.AuthStorage, cache cache.RevokedTokenCache, rbacService rbac.AuthorizationService) *Server {
	return &Server{
		config:      cfg,
		router:      http.NewServeMux(),
		logger:      logger,
		validator:   validator.New(),
		storage:     db,
		cache:       cache,
		rbacService: rbacService,
	}
}

// Start web server. Init data storage and router
func (s *Server) Start() error {
	var err error
	s.auth, err = handlers.NewAuthHandler(s.storage, s.logger, s.cache)

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
	}

	secure_mw := middleware.MiddlewareGroup{
		middleware.LoggingMiddleware(s.logger),
		middleware.AuthMiddleware(s.rbacService),
	}

	auth := http.NewServeMux()
	rbac := http.NewServeMux()

	auth.HandleFunc("POST /refresh", mw.Apply(s.auth.HandleRefreshTokens()))
	auth.HandleFunc("POST /login", mw.Apply(s.auth.HandleLogin()))
	auth.HandleFunc("POST /logout", mw.Apply(s.auth.HandleLogout()))
	auth.HandleFunc("POST /change-password", secure_mw.Apply(s.auth.HandleChangePassword()))

	rbac.Handle("POST /roles", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})) // create role
	rbac.Handle("DELETE /roles/{role}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})) // delete role
	rbac.Handle("POST /actions", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})) // create action
	rbac.Handle("POST /resources", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})) // create resurce
	rbac.Handle("POST /permissions", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})) // create permission
	rbac.Handle("POST /roles/{role_name}/permissions", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})) // grant permission to role
	rbac.Handle("DELETE /roles/{role_name}/permissions", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})) // revoke permission from role
	rbac.Handle("POST /users/{user_id}/roles", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})) // assign role to user
	rbac.Handle("DELETE /users/{user_id}/roles/{role_name}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})) // remove role from user

	s.router.Handle("/api/v1/rbac", http.StripPrefix("/api/v1/rbac", rbac))
	s.router.Handle("/api/v1/auth", http.StripPrefix("/api/v1/auth", auth))

	s.logger.Info("Router configured")
}
