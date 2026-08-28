package main

import (
	"auth/internal/cache"
	"auth/internal/config"
	"auth/internal/grpc"
	"auth/internal/http"
	"auth/internal/storage"
	"auth/rbac"
	"context"
	"flag"
	"fmt"
	"os"
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
	rbacService := rbac.NewCasbinService(&cfg.Casbin.ConfigPath, logger)

	// Create HTTP server
	server := http.CreateServer(&cfg, logger, rbacService, db, cache)

	// Start gRPC server
	go grpc.StartGPRCServer(context.Background(), &cfg, logger, rbacService, db)

	// Init server
	if err := server.Start(); err != nil {
		logger.Error("Server failed", "err", err)
	}
}
