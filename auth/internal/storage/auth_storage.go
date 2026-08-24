package storage

import (
	"auth/internal/config"
	"auth/internal/models"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthStorage interface {
	Open() error
	Close() error
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Read(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID)
}

type AuthPostgresStorage struct {
	config *config.Config
	logger *slog.Logger
	db     *gorm.DB
}

func CreateStorage(cfg *config.Config, l *slog.Logger) AuthStorage {
	return &AuthPostgresStorage{
		config: cfg,
		logger: l,
	}
}

func (s *AuthPostgresStorage) Open() error {
	return nil
}

func (s *AuthPostgresStorage) Close() error {
	return nil
}

func (s *AuthPostgresStorage) Create(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}

func (s *AuthPostgresStorage) Read(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return nil, nil
}

func (s *AuthPostgresStorage) Update(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}

func (s *AuthPostgresStorage) Delete(ctx context.Context, id uuid.UUID) {

}
