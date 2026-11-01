package pgoutbox

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"users/internal/domain/events"
	"users/internal/domain/outbox"
	"users/internal/domain/repository"
)

const claimLease = 5 * time.Minute

type PostgresOutbox struct {
	logger       *slog.Logger
	db           repository.UserStorage
	batchSize    int
	pollInterval time.Duration
	maxAttempts  int
	strategies   map[events.Event]outbox.OutboxSendStrategy
}

func NewPostgresOutbox(
	l *slog.Logger,
	db repository.UserStorage,
	batchSize int,
	pollInterval time.Duration,
	maxAttempts int,
	strategies map[events.Event]outbox.OutboxSendStrategy,
) *PostgresOutbox {
	return &PostgresOutbox{
		logger:       l,
		db:           db,
		batchSize:    batchSize,
		pollInterval: pollInterval,
		maxAttempts:  maxAttempts,
		strategies:   strategies,
	}
}

func (o *PostgresOutbox) Run(ctx context.Context) error {
	o.logger.Info("Outbox worker started", "poll_interval", o.pollInterval, "batch_size", o.batchSize, "max_attempts", o.maxAttempts)
	ticker := time.NewTicker(o.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			o.logger.Info("Outbox worker stopped", "reason", ctx.Err())
			return ctx.Err()
		case <-ticker.C:
			o.logger.Debug("Outbox polling tick")
			if err := o.ProcessBatch(ctx); err != nil && !errors.Is(err, context.Canceled) {
				o.logger.Error("outbox relay batch failed", "err", err)
			}
		}
	}
}

func (o *PostgresOutbox) ProcessBatch(ctx context.Context) error {
	o.logger.Debug("outbox batch processing started")
	if ctx.Err() != nil {
		o.logger.Info("outbox batch processing skipped", "reason", ctx.Err())
		return ctx.Err()
	}
	if o.db == nil {
		o.logger.Debug("outbox batch processing finished", "pending_events", 0, "reason", "db_not_configured")
		return nil
	}

	batch, err := o.db.ClaimOutboxBatch(ctx, o.batchSize, claimLease)
	if err != nil {
		return fmt.Errorf("claim outbox batch: %w", err)
	}
	if len(batch) == 0 {
		o.logger.Debug("outbox batch processing finished", "pending_events", 0)
		return nil
	}

	for _, event := range batch {
		strategy, ok := o.strategies[events.Event(event.EventType)]
		if !ok {
			o.logger.Warn("No outbox strategy registered for event type", "event_id", event.ID, "event_type", event.EventType)
			if retryErr := o.db.RetryOutboxEvent(ctx, event.ID, time.Now().Add(claimLease), "no strategy registered"); retryErr != nil {
				o.logger.Error("Failed to retry outbox event without strategy", "event_id", event.ID, "err", retryErr)
			}
			continue
		}

		o.logger.Info("Processing outbox event", "event_id", event.ID, "event_type", event.EventType, "attempt", event.Attempts)
		if err := strategy.Exec(ctx, event); err != nil {
			o.logger.Error("Outbox event delivery failed", "event_id", event.ID, "event_type", event.EventType, "attempt", event.Attempts, "err", err)
			if retryErr := o.db.RetryOutboxEvent(ctx, event.ID, time.Now().Add(claimLease), err.Error()); retryErr != nil {
				o.logger.Error("Failed to reschedule outbox event", "event_id", event.ID, "err", retryErr)
			}
			continue
		}

		if err := o.db.MarkOutboxDelivered(ctx, event.ID); err != nil {
			o.logger.Error("Failed to mark outbox event delivered", "event_id", event.ID, "err", err)
			continue
		}
		o.logger.Info("Outbox event delivered", "event_id", event.ID, "event_type", event.EventType)
	}

	o.logger.Debug("outbox batch processing finished", "pending_events", len(batch))
	return nil
}
