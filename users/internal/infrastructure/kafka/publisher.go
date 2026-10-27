package kafka

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"time"
	"users/internal/domain/broker"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

type KafkaPublisher struct {
	writer *kafka.Writer
	logger *slog.Logger
}

func (p *KafkaPublisher) Publish(ctx context.Context, payload any) error {
	message, err := p.encodeMessage(payload)
	if err != nil {
		return err
	}
	log.Printf("publishing message to %s topic...", p.writer.Topic)
	return p.writer.WriteMessages(ctx, message)
}

func NewKafkaProducer(brokers []string, topic string, l *slog.Logger) broker.Publisher {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: "registration-service-client",
	}

	c := kafka.WriterConfig{
		Brokers:          brokers,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}

	return &KafkaPublisher{
		writer: kafka.NewWriter(c),
		logger: l,
	}
}

func (p *KafkaPublisher) Close() error {
	err := p.writer.Close()
	if err != nil {
		p.logger.Error("Failed to close Kafka connection", "err", err)
		return err
	}
	p.logger.Info("Kafka connection closed")
	return nil
}

func (p *KafkaPublisher) encodeMessage(payload any) (kafka.Message, error) {
	m, err := json.Marshal(payload)
	if err != nil {
		return kafka.Message{}, err
	}
	key := uuid.New().String()

	return kafka.Message{
		Key:   []byte(key),
		Value: m,
	}, nil
}
