package pgoutbox

import (
	"context"
	"log/slog"
	"time"
	"users/internal/domain/broker"
	"users/internal/domain/repository"
)

type PostgresOutbox struct {
	logger       *slog.Logger
	publisher    broker.Publisher
	db           repository.UserStorage
	batchSize    int
	pollInterval time.Duration
	maxAttempts  int
}

func NewPostgresOutbox(l *slog.Logger, pub broker.Publisher, db repository.UserStorage, batchSize int, pollInterval time.Duration, maxAttempts int) *PostgresOutbox {
	return &PostgresOutbox{
		logger:       l,
		db:           db,
		publisher:    pub,
		batchSize:    batchSize,
		pollInterval: pollInterval,
		maxAttempts:  maxAttempts,
	}
}

func (o *PostgresOutbox) Run(ctx context.Context) error {
	ticker := time.NewTicker(o.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			o.logger.Info("Outbox closed")
			return ctx.Err()
		case <-ticker.C:
			if err := o.ProcessBatch(ctx); err != nil {
				o.logger.Error("outbox relay batch failed", "err", err)
			}
		}
	}

}
func (o *PostgresOutbox) ProcessBatch(ctx context.Context) error {
	return nil
}
