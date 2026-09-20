package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"pkg/rbac"
	"syscall"
	"time"

	"users/internal/config"
	"users/internal/infrastructure/postgres/storage"
	"users/internal/transport/http"
)

var configPath = flag.String("config", "config/config.yaml", "Path to configuration file")

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	if _, ok := os.LookupEnv("RBAC_CONFIG_PATH"); !ok || os.Getenv("RBAC_CONFIG_PATH") == "" {
		_ = os.Setenv("RBAC_CONFIG_PATH", "../rbac_config/config.yaml")
	}

	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger := config.SetupLogger(cfg.Env)

	rbacService, err := rbac.NewRBACService()
	if err != nil {
		logger.Error("Failed to initialize RBAC service", "err", err)
		return
	}

	db := storage.CreateStorage(&cfg, logger)
	if err := db.Open(); err != nil {
		logger.Error("Failed to open database", "err", err)
		return
	}

	srv := http.NewServer(&cfg, logger, rbacService, db)
	srv.Start()

	<-shutdown
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	if err := srv.Shutdown(ctx); err != nil {
		os.Exit(-1)
	}

}
