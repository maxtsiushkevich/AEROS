package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"users/internal/config"
	domain_err "users/internal/domain/errors"
	"users/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Realizes users/internal/domain/repository UsersStorage interface
type PostgresUserStorage struct {
	config *config.Config
	logger *slog.Logger
	db     *gorm.DB
}

func NewStorage(cfg *config.Config, l *slog.Logger) *PostgresUserStorage {
	return &PostgresUserStorage{
		config: cfg,
		logger: l,
	}
}

func (s *PostgresUserStorage) Open() error {
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
		&User{},
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

func (s *PostgresUserStorage) Close() error {
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

func (s *PostgresUserStorage) Create(ctx context.Context, u *models.User) error {
	dbUser := ToPersistence(u)
	if err := s.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		s.logger.Error("Failed to create user", "err", err, "id", u.Id)
		return err
	}
	return nil
}

func (s *PostgresUserStorage) Read(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var dbUser User
	if err := s.db.WithContext(ctx).First(&dbUser, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain_err.ErrUserNotFound
		}
		s.logger.Error("Failed to read user", "err", err, "id", id)
		return nil, err
	}
	return ToDomain(&dbUser), nil
}

func (s *PostgresUserStorage) Update(ctx context.Context, u *models.User) error {
	dbUser := ToPersistence(u)
	if err := s.db.WithContext(ctx).Save(dbUser).Error; err != nil {
		s.logger.Error("Failed to update user", "err", err, "id", u.Id)
		return err
	}
	return nil
}

func (s *PostgresUserStorage) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.db.WithContext(ctx).Delete(&User{}, "id = ?", id).Error; err != nil {
		s.logger.Error("Failed to delete user", "err", err, "id", id)
		return err
	}
	return nil
}

func (s *PostgresUserStorage) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var dbUser User
	if err := s.db.WithContext(ctx).First(&dbUser, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain_err.ErrUserNotFound
		}
		return nil, err
	}
	return ToDomain(&dbUser), nil
}
