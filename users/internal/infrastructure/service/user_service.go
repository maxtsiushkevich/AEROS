package service

import (
	"context"
	"time"
	auth "users/api/proto"
	"users/internal/domain/models"
	"users/internal/domain/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type UserService struct {
	grpcConn *grpc.ClientConn
	storage  repository.UserStorage
}

func NewUsersService(grpcConn *grpc.ClientConn, storage repository.UserStorage) *UserService {
	return &UserService{
		grpcConn: grpcConn,
		storage:  storage,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, name string, email string, birthday time.Time, password string) (access *string, refresh *string, err error) {
	u, err := models.NewUser(name, email, birthday)
	if err != nil {
		return nil, nil, err
	}

	client := auth.NewAuthClient(s.grpcConn)

	err = s.storage.Create(ctx, u)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.AddUser(context.Background(), &auth.AddUserRequest{
		Id:       u.Id.String(),
		Password: password,
		Email:    email,
	})

	if err != nil {
		s.storage.Delete(ctx, u.Id)
		return nil, nil, err
	}

	return &resp.AccessToken, &resp.RefreshToken, nil
}

func (s *UserService) ActivateUser(ctx context.Context, id uuid.UUID) error {
	u, err := s.storage.Read(ctx, id)
	if err != nil {
		return err
	}

	if err := u.Activate(); err != nil {
		return err
	}

	return s.storage.Update(ctx, u)
}

func (s *UserService) Profile(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.storage.Read(ctx, id)
}
