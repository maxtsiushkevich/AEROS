package http

import (
	"auth/internal/cache"
	"auth/internal/config"
	"auth/internal/handlers"
	"auth/internal/storage"
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"pkg/middleware"
	"pkg/rbac"
)

type Server struct {
	config *config.Config
	server *http.Server
	logger *slog.Logger
	auth   *handlers.AuthHandler
	rbac   *handlers.RbacHandler

	storage     storage.AuthStorage
	cache       cache.RevokedTokenCache
	rbacService rbac.AuthorizationService
}

func CreateServer(cfg *config.Config, l *slog.Logger, db storage.AuthStorage, c cache.RevokedTokenCache, rbac rbac.AuthorizationService) *Server {
	return &Server{
		config:      cfg,
		logger:      l,
		storage:     db,
		cache:       c,
		rbacService: rbac,
	}
}

// Start web server. Init data storage and router
func (s *Server) Start() error {
	s.auth = handlers.NewAuthHandler(s.storage, s.logger, s.cache)
	s.rbac = handlers.NewRbacHandler(s.rbacService, s.logger)

	router := http.NewServeMux()
	s.configureRouter(router)

	s.logger.Info("Start server", "env", s.config.Env)
	s.logger.Debug("Serve on", "addr", "http://"+s.config.HTTPServer.Address)

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

	rbac.Handle("PUT /roles", mw.Apply(s.rbac.CreateRole()))           // create role
	rbac.Handle("DELETE /roles/{role}", mw.Apply(s.rbac.DeleteRole())) // delete role

	rbac.Handle("PUT /actions", mw.Apply(s.rbac.CreateAction()))         // create action
	rbac.Handle("PUT /resources", mw.Apply(s.rbac.CreateResource()))     // create resource
	rbac.Handle("PUT /permissions", mw.Apply(s.rbac.CreatePermission())) // create permission

	rbac.Handle("PUT /roles/{role_name}/permissions", mw.Apply(s.rbac.GrantPermissionToRole())) // grant permission to role

	rbac.Handle("DELETE /roles/{role_name}/permissions", mw.Apply(s.rbac.RevokePermissionFromRole())) // revoke permission from role
	// DELETE /api/v1/rbac/roles/admin/permissions?resource=rbac/resources&action=read

	rbac.Handle("PUT /users/{user_id}/roles", mw.Apply(s.rbac.AssignRoleToUser()))                  // assign role to user
	rbac.Handle("DELETE /users/{user_id}/roles/{role_name}", mw.Apply(s.rbac.RemoveRoleFromUser())) // remove role from user

	router.Handle("/api/v1/rbac/", http.StripPrefix("/api/v1/rbac", rbac))
	router.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", auth))

	s.logger.Info("Router configured")
}
