package main

import (
	"flag"
	"fmt"
	"log/slog"

	"flights/internal/config"
	"flights/internal/http"
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

	// Init server
	server := http.CreateServer(&cfg)
	if err := server.Start(); err != nil {
		slog.Error("Server failed", "err", err)
	}
}
