package broker

import "context"

type Reader interface {
	Read(ctx context.Context) any
	Close() error
}
