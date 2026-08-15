package storage

import (
	"flights/internal/config"
	"flights/internal/models"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	config *config.Config
	logger *slog.Logger
	db     *gorm.DB
}

func CreateStorage(cfg *config.Config, l *slog.Logger) *Storage {
	return &Storage{
		config: cfg,
		logger: l,
	}
}

func (s *Storage) Open() error {
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s",
		s.config.Database.User,
		s.config.Database.Password,
		s.config.Database.Address,
		s.config.Database.DbName)

	var err error

	s.db, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
	s.logger.Debug("Open db connect", "connString", connString)
	if err != nil {
		return err
	}

	if err := s.autoMigrateModels(); err != nil {
		return err
	}
	s.logger.Debug("Init schemas")

	return nil
}

func (s *Storage) autoMigrateModels() error {
	if err := s.db.AutoMigrate(&models.Flight{}); err != nil {
		s.logger.Error("Failed to auto-migrate models", "err", err)
		return err
	}

	flight := models.Flight{FlightNumber: "B2-2555", Origin: "MSQ", Destination: "DBX", Aircraft: "Boeing 737 Max 7"}
	s.db.Create(&flight)

	return nil
}

func (s *Storage) Close() error {
	database, err := s.db.DB()
	if err != nil {
		return err
	}
	database.Close()
	s.logger.Info("Database connection closed", "err", err)
	return nil
}

func (s *Storage) Create(flight models.Flight) {

}

func (s *Storage) Read(id uuid.UUID) {

}

func (s *Storage) ReadByFlightNumber(flightNumber string) (*models.Flight, error) {
	var flight models.Flight

	result := s.db.Where("flight_number = ?", flightNumber).First(&flight)

	if result.Error != nil {
		s.logger.Error("Failed to read flight", "flightNumber", flightNumber, "err", result.Error)
		return nil, result.Error
	}

	s.logger.Debug("Flight found", "flightNumber", flightNumber)
	return &flight, nil
}

func (s *Storage) Update(flight models.Flight) {

}

func (s *Storage) Delete(id uuid.UUID) {

}
