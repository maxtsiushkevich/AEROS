package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"users/internal/config"
	domainerr "users/internal/domain/errors"
	"users/internal/domain/models"
	"users/internal/domain/outbox"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		&outbox.OutboxEvent{},
	); err != nil {
		return err
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
			return nil, domainerr.ErrUserNotFound
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
			return nil, domainerr.ErrUserNotFound
		}
		return nil, err
	}
	return ToDomain(&dbUser), nil
}

func (s *PostgresUserStorage) UpdateEmailWithOutbox(ctx context.Context, id uuid.UUID, email string, event outbox.OutboxEvent) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		s.logger.Debug("Creating outbox event for email update", "user_id", id, "new_email", email, "event_type", event.EventType)

		result := tx.Model(&User{}).Where("id = ?", id).Update("email", email)
		if result.Error != nil {
			s.logger.Error("Failed to update user email for outbox event", "err", result.Error, "user_id", id)
			return fmt.Errorf("update user email: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			s.logger.Warn("Outbox email update skipped because user was not found", "user_id", id)
			return domainerr.ErrUserNotFound
		}

		if err := tx.Create(&event).Error; err != nil {
			s.logger.Error("Failed to insert outbox event", "err", err, "user_id", id, "event_type", event.EventType)
			return fmt.Errorf("insert outbox event: %w", err)
		}

		s.logger.Info("Outbox event created", "event_id", event.ID, "user_id", id, "event_type", event.EventType, "aggregate_id", event.AggregateID)
		return nil
	})
}

func (s *PostgresUserStorage) ClaimOutboxBatch(ctx context.Context, limit int, lease time.Duration) ([]outbox.OutboxEvent, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("outbox batch size must be positive")
	}
	if lease <= 0 {
		return nil, fmt.Errorf("outbox claim lease must be positive")
	}

	var batch []outbox.OutboxEvent
	result := s.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("status = ?", outbox.StatusPending).
		Order("id ASC").
		Limit(limit).
		Find(&batch)
	if result.Error != nil {
		return nil, fmt.Errorf("select outbox events: %w", result.Error)
	}

	for i := range batch {
		batch[i].Attempts++
		batch[i].Status = outbox.StatusProcessing

		updateResult := s.db.WithContext(ctx).Model(&outbox.OutboxEvent{}).
			Where("id = ?", batch[i].ID).
			Updates(map[string]any{
				"attempts": batch[i].Attempts,
				"status":   outbox.StatusProcessing,
			})
		if updateResult.Error != nil {
			return nil, fmt.Errorf("claim outbox event %s: %w", batch[i].ID, updateResult.Error)
		}
	}

	return batch, nil
}

func (s *PostgresUserStorage) MarkOutboxDelivered(ctx context.Context, id uuid.UUID) error {
	result := s.db.WithContext(ctx).Model(&outbox.OutboxEvent{}).
		Where("id = ? AND status = ?", id, outbox.StatusProcessing).
		Updates(map[string]any{"status": outbox.StatusDelivered})
	if result.Error != nil {
		return fmt.Errorf("mark outbox event %s delivered: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("outbox event %s is not processing", id)
	}
	return nil
}

func (s *PostgresUserStorage) RetryOutboxEvent(ctx context.Context, id uuid.UUID, _ time.Time, _ string) error {
	result := s.db.WithContext(ctx).Model(&outbox.OutboxEvent{}).
		Where("id = ? AND status = ?", id, outbox.StatusProcessing).
		Updates(map[string]any{"status": outbox.StatusPending})
	if result.Error != nil {
		return fmt.Errorf("schedule retry for outbox event %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("outbox event %s is not processing", id)
	}
	return nil
}
