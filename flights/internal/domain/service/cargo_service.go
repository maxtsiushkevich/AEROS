package service

import (
	"context"

	"github.com/google/uuid"
)

type CargoService interface {
	AddItem(ctx context.Context, id uuid.UUID)
}
