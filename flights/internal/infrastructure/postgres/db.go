package postgres

import (
	"flights/internal/config"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(cfg *config.Config, l *slog.Logger) (*gorm.DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Address,
		cfg.Database.DbName)

	var err error

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	l.Debug("Open db connect", "connString", connString)
	if err != nil {
		return nil, err
	}

	pool, _ := db.DB()

	pool.SetMaxOpenConns(5)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(30 * time.Second)
	pool.SetConnMaxIdleTime(15 * time.Second)

	return db, nil
}

func Close(db *gorm.DB, l *slog.Logger) error {
	database, err := db.DB()
	if err != nil {
		return err
	}
	database.Close()
	l.Info("Database connection closed", "err", err)
	return nil
}
