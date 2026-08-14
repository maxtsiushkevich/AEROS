package main

import (
	"flag"
	"fmt"

	"flights/internal/config"
	"flights/internal/http"
	"flights/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	server := http.CreateServer(&cfg)
	server.Start()

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Address,
		cfg.Database.DbName)

	// Init DB
	db, _ := gorm.Open(postgres.Open(connString), &gorm.Config{})
	db.AutoMigrate(&models.Flight{})
	flight := models.Flight{FlightNumber: "B2-2555", Origin: "MSQ", Destination: "DBX", Aircraft: "Boeing 737 Max 7"}
	db.Create(&flight)
}
