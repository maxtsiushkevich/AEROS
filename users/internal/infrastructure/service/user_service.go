package service

import (
	"context"
	"errors"
	"log/slog"
	"time"
	auth "users/api/proto"
	"users/internal/domain/cache"
	domainerr "users/internal/domain/errors"
	"users/internal/domain/models"
	"users/internal/domain/otp"
	"users/internal/domain/producer"
	"users/internal/domain/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type UserService struct {
	logger    *slog.Logger
	grpcConn  *grpc.ClientConn
	storage   repository.UserStorage
	publisher producer.Publisher
	otpCache  cache.Cache
}

func NewUsersService(logger *slog.Logger, grpcConn *grpc.ClientConn, storage repository.UserStorage, publisher producer.Publisher, otpCache cache.Cache) *UserService {
	return &UserService{
		logger:    logger,
		grpcConn:  grpcConn,
		storage:   storage,
		publisher: publisher,
		otpCache:  otpCache,
	}
}

func (s *UserService) sendCode(ctx context.Context, email string, purpose otp.OTPPurpose) error {
	code, err := otp.GenerateOTP(6)
	if err != nil {
		s.logger.Error("Failed to generate OTP", "err", err)
		return err
	}

	err = s.otpCache.Set(ctx, email, purpose, []byte(code), 10*time.Minute)
	if err != nil {
		s.logger.Error("Failed to store OTP in cache", "err", err)
		return err
	}

	if s.publisher != nil {
		go s.publisher.Publish(ctx, producer.GetOTPEventMessage(email, code, purpose))
	}
	return nil
}

func (s *UserService) confirmCode(ctx context.Context, email, code string, purpose otp.OTPPurpose) error {
	cachedCode, err := s.otpCache.Get(ctx, email, purpose)
	if err != nil {
		return domainerr.ErrUserActivationFailed
	}

	if string(cachedCode) != code {
		return domainerr.ErrInvalidOTP
	}

	err = s.otpCache.Delete(ctx, email, purpose)
	if err != nil {
		return domainerr.ErrUserActivationFailed
	}

	return nil
}

func (s *UserService) RegisterUser(ctx context.Context, name string, email string, birthday time.Time, password string) (access *string, refresh *string, err error) {
	existingUser, err := s.storage.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domainerr.ErrUserNotFound) {
		return nil, nil, err
	}
	if existingUser != nil {
		return nil, nil, domainerr.ErrUserAlreadyExists
	}

	u, err := models.NewUser(name, email, birthday)
	if err != nil {
		return nil, nil, err
	}

	client := auth.NewAuthClient(s.grpcConn)

	err = s.storage.Create(ctx, u)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.AddUser(ctx, &auth.AddUserRequest{
		Id:       u.Id.String(),
		Password: password,
		Email:    email,
	})
	if err != nil {
		s.storage.Delete(ctx, u.Id)
		return nil, nil, err
	}

	err = s.sendCode(ctx, email, otp.ActivateProfile)
	if err != nil {
		return &resp.AccessToken, &resp.RefreshToken, nil
	}
	return &resp.AccessToken, &resp.RefreshToken, nil
}

func (s *UserService) ActivateUserByAdmin(ctx context.Context, id uuid.UUID) error {
	u, err := s.storage.Read(ctx, id)
	if err != nil {
		return err
	}

	if err := u.Activate(); err != nil {
		return err
	}

	return s.storage.Update(ctx, u)
}

func (s *UserService) ActivateUserByCode(ctx context.Context, id uuid.UUID, code string) error {
	u, err := s.storage.Read(ctx, id)
	if err != nil {
		return err
	}

	err = s.confirmCode(ctx, u.Email, code, otp.ActivateProfile)
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

func (s *UserService) RequestCode(ctx context.Context, id uuid.UUID, email string, purpose otp.OTPPurpose) error {
	err := s.sendCode(ctx, email, purpose)
	if err != nil {
		return err
	}

	return nil
}
