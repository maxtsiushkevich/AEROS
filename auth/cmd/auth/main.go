package main

import (
	"auth/internal/config"
	"auth/internal/grpc"
	"auth/internal/http"
	"auth/rbac"
	"context"
	"flag"
	"fmt"
	"log/slog"
)

var configPath = flag.String("config", "config/config.yaml", "Path to configuration file")

func main() {

	// Load config
	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start gRPC server
	go grpc.StartGPRCServer(context.Background(), &cfg)

	// Init server

	rbacService := rbac.NewCasbinService(cfg.Casbin.ConfigPath)
	server := http.CreateServer(&cfg, rbacService)

	if err := server.Start(); err != nil {
		slog.Error("Server failed", "err", err)
	}
}
