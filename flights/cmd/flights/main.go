package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flights/internal/config"
	"flights/internal/infrastructure/postgres"
	"flights/internal/transport/http"
)

var configPath = flag.String("config", "config/config.yaml", "Path to configuration file")

func main() {

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Load config
	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger := config.SetupLogger(cfg.Env)

	db := postgres.CreateStorage(&cfg, logger)
	if err := db.Open(); err != nil {
		logger.Error("Failed to open database", "err", err)
		return
	}

	// Init server
	server := http.CreateServer(&cfg, logger, db)
	if err := server.Start(); err != nil {
		slog.Error("Server failed", "err", err)
	}

	interruptSignal := <-shutdown

	fmt.Printf("\nReceived an interrupt signal (%d)\n", interruptSignal)
	fmt.Println("Shutting down HTTP serve gracefully with 10-second timeout")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	if err := server.Shutdown(ctx); err != nil {
		os.Exit(-1)
	}

	if err := db.Close(); err != nil {
		logger.Error("Failed to close database", "err", err)
	}

}
