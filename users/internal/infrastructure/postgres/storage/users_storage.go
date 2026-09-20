package storage

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"users/internal/config"
	"users/internal/domain/storage"
	"users/internal/infrastructure/postgres/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
 
// Realizes users/internal/domain/storage UsersStorage interface
type UsersPostgresStorage struct {
	config *config.Config
	logger *slog.Logger
	db     *gorm.DB
}

func CreateStorage(cfg *config.Config, l *slog.Logger) storage.UsersStorage {
	return &UsersPostgresStorage{
		config: cfg,
		logger: l,
	}
}

func (s *UsersPostgresStorage) Open() error {
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

	if err := s.db.AutoMigrate(
		&models.User{},
	); err != nil {
		return nil
	}

	pool, _ := s.db.DB()

	pool.SetMaxOpenConns(5)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(30 * time.Second)
	pool.SetConnMaxIdleTime(15 * time.Second)

	return nil
}

func (s *UsersPostgresStorage) Close() error {
	database, err := s.db.DB()
	if err != nil {
		return err
	}

	err = database.Close()
	if err != nil {
		s.logger.Error("Failed to close database connection", "err", err)
		return err
	}

	s.logger.Info("Database connection closed")
	return nil
}

func (s *UsersPostgresStorage) Create(ctx context.Context) error {
	return nil
}

func (s *UsersPostgresStorage) Read(ctx context.Context) error {
	return nil
}

func (s *UsersPostgresStorage) ReadByID(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *UsersPostgresStorage) Update(ctx context.Context) error {
	return nil
}

func (s *UsersPostgresStorage) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
