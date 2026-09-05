package storage

import (
	"auth/internal/config"
	"auth/internal/models"
	"auth/pkg/errors"
	"context"
	errors_pkg "errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthStorage interface {
	Open() error
	Close() error
	Create(ctx context.Context, user *models.UserAuth) (*models.UserAuth, error)
	Read(ctx context.Context, email *string) (*models.UserAuth, error)
	ReadByID(ctx context.Context, id uuid.UUID) (*models.UserAuth, error)
	Update(ctx context.Context, user *models.UserAuthUpdate) (*models.UserAuth, error)
	Delete(ctx context.Context, id uuid.UUID) error
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

	pool, _ := s.db.DB()

	pool.SetMaxOpenConns(5)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(30 * time.Second)
	pool.SetConnMaxIdleTime(15 * time.Second)

	return nil
}

func (s *AuthPostgresStorage) Close() error {
	database, err := s.db.DB()
	if err != nil {
		return err
	}
	database.Close()
	s.logger.Info("Database connection closed", "err", err)
	return nil
}

func (s *AuthPostgresStorage) Create(ctx context.Context, user *models.UserAuth) (*models.UserAuth, error) {
	err := s.db.WithContext(ctx).Create(user).Error
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "duplicate key") || strings.Contains(errMsg, "UNIQUE constraint") {
			s.logger.Warn("User already exists", "id", user.ID, "email", user.Email)
			return nil, errors.UserAlreadyExistsError
		}
		s.logger.Error("Failed to create user", "email", user.Email, "err", err)
		return nil, errors.NewAuthError("CREATE_FAILED", "failed to create user")
	}

	return user, nil
}

func (s *AuthPostgresStorage) Read(ctx context.Context, email *string) (*models.UserAuth, error) {
	var authData models.UserAuth
	err := s.db.WithContext(ctx).Where("email = ?", email).First(&authData).Error

	if err != nil {
		if errors_pkg.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("User not found", "email", email)
			return nil, errors.UserNotFoundError
		}
		s.logger.Error("Failed to read user", "email", email, "err", err)
		return nil, errors.NewAuthError("READ_FAILED", "failed to read user")
	}

	return &authData, nil
}

func (s *AuthPostgresStorage) ReadByID(ctx context.Context, id uuid.UUID) (*models.UserAuth, error) {
	var authData models.UserAuth
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&authData).Error
	if err != nil {
		if errors_pkg.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAuthError("USER_NOT_FOUND", fmt.Sprintf("auth data with id %s not found", id))
		}
		return nil, errors.NewAuthError("READ_FAILED", "failed to read user")
	}

	return &authData, nil
}

func (s *AuthPostgresStorage) Update(ctx context.Context, update *models.UserAuthUpdate) (*models.UserAuth, error) {
	var existing models.UserAuth

	err := s.db.WithContext(ctx).Where("id = ?", update.ID).First(&existing).Error
	if err != nil {
		if errors_pkg.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("User not found for update", "id", update.ID)
			return nil, errors.NewAuthError("USER_NOT_FOUND", fmt.Sprintf("auth data with id %s not found", update.ID))
		}
		s.logger.Error("Failed to fetch user for update", "id", update.ID, "err", err)
		return nil, errors.NewAuthError("UPDATE_FAILED", "failed to fetch user for update")
	}

	if update.NewEmail != "" {
		existing.Email = update.NewEmail
	}
	if update.NewHashedPassword != "" {
		existing.HashedPassword = update.NewHashedPassword
	}

	existing.Version++

	if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
		s.logger.Error("Failed to update user", "id", update.ID, "err", err)
		return nil, errors.NewAuthError("SAVE_FAILED", "failed to save updated user")
	}

	return &existing, nil
}

func (s *AuthPostgresStorage) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.db.WithContext(ctx).Where("id = ?", id).Delete(&models.UserAuth{}).Error
	if err != nil {
		s.logger.Error("Failed to delete user", "id", id, "err", err)
		return errors.NewAuthError("DELETE_FAILED", "failed to delete user")
	}
	return nil
}
