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
	"os/signal"
	"pkg/rbac"
	"syscall"
	"time"
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
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

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

	// Setup RBAC service
	rbacService, err := rbac.NewRBACService()
	if err != nil {
		logger.Error("Failed to initialize RBAC service", "err", err)
		return
	}

	// Create HTTP server
	server := http.CreateServer(&cfg, logger, db, cache, rbacService)

	// Start gRPC server
	grpcServer, err := grpc.StartGPRCServer(&cfg, logger, db, rbacService)
	if err != nil {
		logger.Error("Failed to start gRPC server", "err", err)
		return
	}

	// Init server
	if err := server.Start(); err != nil {
		logger.Error("Server failed", "err", err)
	}

	interruptSignal := <-shutdown

	fmt.Printf("\nReceived an interrupt signal (%d)\n", interruptSignal)
	fmt.Println("Shutting down HTTP serve gracefully with 10-second timeout")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	grpcStopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(grpcStopped)
	}()

	if err := server.Shutdown(ctx); err != nil {
		os.Exit(-1)
	}

	select {
	case <-grpcStopped:
		logger.Info("Finished graceful shutdown for the gRPC server")
	case <-ctx.Done():
		grpcServer.Stop()
		logger.Warn("Forced shutdown for the gRPC server after timeout")
	}

	if err := db.Close(); err != nil {
		logger.Error("Failed to close database", "err", err)
	}

	if err := cache.Shutdown(); err != nil {
		logger.Error("Failed to close cache", "err", err)
	}

	if err := rbacService.Close(); err != nil {
		logger.Error("Failed to close RBAC service", "err", err)
	} else {
		logger.Info("RBAC service closed")
	}
}
