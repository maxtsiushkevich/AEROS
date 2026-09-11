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
	config *config.Config
	router *http.ServeMux
	logger *slog.Logger

	validator *validator.Validate
	auth      *handlers.AuthHandler
	rbac      *handlers.RbacHandler

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
	s.auth = handlers.NewAuthHandler(s.storage, s.logger, s.cache)
	s.rbac = handlers.NewRbacHandler(s.rbacService, s.logger)

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

	s.router.Handle("/api/v1/rbac/", http.StripPrefix("/api/v1/rbac", rbac))
	s.router.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", auth))

	s.logger.Info("Router configured")
}

// 123e4567-e89b-12d3-a456-426614175000
