package main

import (
	"flag"
	"fmt"

	"users/internal/config"
	"users/internal/http"
)

var configPath = flag.String("config", "config/config.yaml", "Path to configuration file")

func main() {

	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	srv := http.NewServer(&cfg)
	srv.ConfigServer()
}
