package main

import (
	"flag"
	"fmt"
	"os"
	"pkg/rbac"

	"users/internal/config"
	"users/internal/http"
)

var configPath = flag.String("config", "config/config.yaml", "Path to configuration file")

func main() {
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

	rbacService, err := rbac.NewRBACServiceFromEnv()
	if err != nil {
		logger.Error("Failed to initialize RBAC service", "err", err)
		return
	}

	srv := http.NewServer(&cfg, logger, rbacService)
	srv.ConfigServer()
}
