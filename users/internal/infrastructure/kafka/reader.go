package kafka

import (
	"context"
	"log/slog"
	"time"
	"users/internal/domain/broker"

	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	reader *kafka.Reader
	logger *slog.Logger
}

func NewKafkaReader(brokers []string, topic string, l *slog.Logger) broker.Reader {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: "registration-service-client",
	}

	c := kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		Dialer:         dialer,
		CommitInterval: 0,
	}

	return &KafkaReader{
		reader: kafka.NewReader(c),
		logger: l,
	}
}

func (r *KafkaReader) Close() error {
	err := r.reader.Close()
	if err != nil {
		r.logger.Error("Failed to close Kafka reader", "err", err)
		return err
	}
	r.logger.Info("Kafka reader closed")
	return nil
}

func (r *KafkaReader) Read(ctx context.Context) any {
	return nil
}
