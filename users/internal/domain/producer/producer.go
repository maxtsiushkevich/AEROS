package producer

import "context"

type Publisher interface {
	Publish(ctx context.Context, payload any) error
}
