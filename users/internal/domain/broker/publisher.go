package broker

import "context"

type Publisher interface {
	Publish(ctx context.Context, payload any) error
	Close() error
}
