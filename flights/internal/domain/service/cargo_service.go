package service

import (
	"context"

	"github.com/google/uuid"
)

type CargoService interface {
	AddItem(ctx context.Context, flightID uuid.UUID)
}
