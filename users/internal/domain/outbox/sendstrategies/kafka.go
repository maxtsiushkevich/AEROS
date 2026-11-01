package sendstrategies

import (
	"context"
	"users/internal/domain/broker"
	"users/internal/domain/outbox"
)

type KafkaStrategy struct {
	pub broker.Publisher
}

func NewKafkaStrategy() (*KafkaStrategy, error) {
	return &KafkaStrategy{}, nil
}

func (s *KafkaStrategy) Exec(ctx context.Context, event outbox.OutboxEvent) error {
	err := s.pub.Publish(ctx, event.Payload)
	if err != nil {
		return err
	}
	return nil
}
