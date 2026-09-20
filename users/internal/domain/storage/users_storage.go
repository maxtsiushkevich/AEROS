package storage

import (
	"context"

	"github.com/google/uuid"
)

type UsersStorage interface {
	Open() error
	Close() error
	Create(ctx context.Context) error
	Read(ctx context.Context) error
	ReadByID(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context) error
	Delete(ctx context.Context, id uuid.UUID) error
}
