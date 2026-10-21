package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
	auth "users/api/proto"
	"users/internal/domain/cache"
	domainerr "users/internal/domain/errors"
	"users/internal/domain/models"
	"users/internal/domain/otp"
	"users/internal/domain/producer"
	"users/internal/domain/repository"
	"users/internal/domain/service"
	userService "users/internal/infrastructure/service"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type UserUseCases struct {
	logger      *slog.Logger
	grpcConn    *grpc.ClientConn
	storage     repository.UserStorage
	publisher   producer.Publisher
	otpCache    cache.Cache
	userService service.UserService
}

type changeEmailOTP struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Code   string    `json:"code"`
}

func NewUserUsecases(logger *slog.Logger, grpcConn *grpc.ClientConn, storage repository.UserStorage, publisher producer.Publisher, otpCache cache.Cache) *UserUseCases {
	return &UserUseCases{
		logger:      logger,
		grpcConn:    grpcConn,
		storage:     storage,
		publisher:   publisher,
		otpCache:    otpCache,
		userService: userService.NewUserService(storage),
	}
}

func (s *UserUseCases) sendCode(ctx context.Context, email string, purpose otp.OTPPurpose) error {
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

func (s *UserUseCases) confirmCode(ctx context.Context, email, code string, purpose otp.OTPPurpose) error {
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

func (s *UserUseCases) RegisterUser(ctx context.Context, name string, email string, birthday time.Time, password string) (access *string, refresh *string, err error) {
	err = s.userService.CheckEmailAvailability(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	u, err := models.NewUser(name, email, birthday)
	if err != nil {
		return nil, nil, err
	}

	err = s.storage.Create(ctx, u)
	if err != nil {
		return nil, nil, err
	}

	client := auth.NewAuthClient(s.grpcConn)
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

func (s *UserUseCases) ActivateUserByAdmin(ctx context.Context, id uuid.UUID) error {
	u, err := s.storage.Read(ctx, id)
	if err != nil {
		return err
	}

	if err := u.Activate(); err != nil {
		return err
	}

	return s.storage.Update(ctx, u)
}

func (s *UserUseCases) ActivateUserByCode(ctx context.Context, id uuid.UUID, code string) error {
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

func (s *UserUseCases) Profile(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.storage.Read(ctx, id)
}

func (s *UserUseCases) RequestCode(ctx context.Context, id uuid.UUID, email string, purpose otp.OTPPurpose) error {
	err := s.sendCode(ctx, email, purpose)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserUseCases) RequestChangeEmail(ctx context.Context, id uuid.UUID, newEmail string) error {
	if _, err := s.storage.Read(ctx, id); err != nil {
		return err
	}

	err := s.userService.CheckEmailAvailability(ctx, newEmail)
	if err != nil {
		return err
	}

	code, err := otp.GenerateOTP(6)
	if err != nil {
		s.logger.Error("Failed to generate OTP", "err", err)
		return err
	}

	payload, err := json.Marshal(changeEmailOTP{UserID: id, Email: newEmail, Code: code})
	if err != nil {
		return err
	}

	if err := s.otpCache.Set(ctx, id.String(), otp.ChangeEmail, payload, 10*time.Minute); err != nil {
		s.logger.Error("Failed to store OTP in cache", "err", err)
		return err
	}

	if s.publisher != nil {
		go s.publisher.Publish(ctx, producer.GetOTPEventMessage(newEmail, code, otp.ChangeEmail))
	}
	return nil
}

func (s *UserUseCases) ConfirmChangeEmail(ctx context.Context, id uuid.UUID, code string) error {
	cachedPayload, err := s.otpCache.Get(ctx, id.String(), otp.ChangeEmail)
	if err != nil {
		return domainerr.ErrUserActivationFailed
	}

	var cachedOTP changeEmailOTP
	if err := json.Unmarshal(cachedPayload, &cachedOTP); err != nil {
		return domainerr.ErrUserActivationFailed
	}
	if cachedOTP.UserID != id || cachedOTP.Code != code || cachedOTP.Email == "" {
		return domainerr.ErrInvalidOTP
	}
	if err := s.otpCache.Delete(ctx, id.String(), otp.ChangeEmail); err != nil {
		return domainerr.ErrUserActivationFailed
	}

	newEmail := cachedOTP.Email
	if err := s.userService.CheckEmailAvailability(ctx, newEmail); err != nil {
		return err
	}

	u, err := s.storage.Read(ctx, id)
	if err != nil {
		return err
	}

	client := auth.NewAuthClient(s.grpcConn)
	changeResp, err := client.ChangeEmail(ctx, &auth.ChangeEmailRequest{
		Id:       id.String(),
		NewEmail: newEmail,
	})
	if err != nil {
		return err
	}
	if !changeResp.IsChanged {
		return domainerr.ErrUserAlreadyExists
	}

	oldEmail := u.Email
	u.Email = newEmail
	if err := s.storage.Update(ctx, u); err != nil {
		storageErr := fmt.Errorf("failed to update user email in storage: %w", err)
		rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()

		rollbackResp, rollbackErr := client.ChangeEmail(rollbackCtx, &auth.ChangeEmailRequest{
			Id:       id.String(),
			NewEmail: oldEmail,
		})
		if rollbackErr != nil {
			return errors.Join(storageErr, fmt.Errorf("failed to roll back auth email: %w", rollbackErr))
		}
		if !rollbackResp.IsChanged {
			return errors.Join(storageErr, errors.New("auth service rejected email rollback"))
		}
		return storageErr
	}
	return nil
}
