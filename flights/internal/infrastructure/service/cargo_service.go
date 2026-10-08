package service

import (
	"context"
	"flights/internal/domain/repository"

	"github.com/google/uuid"
)

type CargoDomainService struct {
	repo repository.CargoRepository
}

func NewCargoDomainService(repo repository.CargoRepository) *CargoDomainService {
	return &CargoDomainService{repo: repo}
}

func (s CargoDomainService) AddItem(ctx context.Context, id uuid.UUID) {

}
