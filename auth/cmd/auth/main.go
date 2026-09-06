package main

import (
	"auth/internal/cache"
	"auth/internal/config"
	"auth/internal/grpc"
	"auth/internal/http"
	"auth/internal/storage"
	"context"
	"flag"
	"fmt"
	"os"
	"pkg/rbac"
)

var configPath = flag.String("config", "config/config.yaml", "Path to configuration file")

func ensureJWTEnv() {
	if _, ok := os.LookupEnv("JWT_SECRET"); !ok || os.Getenv("JWT_SECRET") == "" {
		_ = os.Setenv("JWT_SECRET", "dev-access-secret")
	}
	if _, ok := os.LookupEnv("JWT_REFRESH_SECRET"); !ok || os.Getenv("JWT_REFRESH_SECRET") == "" {
		_ = os.Setenv("JWT_REFRESH_SECRET", "dev-refresh-secret")
	}
}

func main() {
	ensureJWTEnv()

	if _, ok := os.LookupEnv("RBAC_CONFIG_PATH"); !ok || os.Getenv("RBAC_CONFIG_PATH") == "" {
		_ = os.Setenv("RBAC_CONFIG_PATH", "../rbac_config/config.yaml")
	}

	// Load config
	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger := config.SetupLogger(cfg.Env)

	// Setup database
	db := storage.CreateStorage(&cfg, logger)
	if err := db.Open(); err != nil {
		logger.Error("Failed to open database", "err", err)
		return
	}

	// Setup cache
	cache, err := cache.NewRedis(&cfg, logger)
	if err != nil {
		logger.Error("Failed to open Redis", "err", err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database", "err", err)
		}
	}()

	// Setup RBAC service
	rbacService, err := rbac.NewRBACServiceFromEnv()
	if err != nil {
		logger.Error("Failed to initialize RBAC service", "err", err)
		return
	}

	// Create HTTP server
	server := http.CreateServer(&cfg, logger, db, cache, rbacService)

	// Start gRPC server
	go grpc.StartGPRCServer(context.Background(), &cfg, logger, db, rbacService)

	// Init server
	if err := server.Start(); err != nil {
		logger.Error("Server failed", "err", err)
	}
}
