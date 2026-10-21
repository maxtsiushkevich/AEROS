package grpc

import (
	grpc "auth/api/proto"
	"auth/internal/models"
	"auth/pkg/errors"
	"context"
	"pkg/auth"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Auth) AddUser(ctx context.Context, req *grpc.AddUserRequest) (*grpc.AddUserResponse, error) {
	userAuthData, err := getUserAuthData(req)
	if err != nil {
		return nil, err
	}

	// Create user
	userAuthData, err = s.storage.Create(ctx, userAuthData)

	if err != nil {
		if authErr, ok := err.(*errors.AuthError); ok {
			if authErr.Code == "USER_ALREADY_EXISTS" {
				return nil, status.Error(codes.AlreadyExists, authErr.Error())
			}
			s.storage.Delete(ctx, userAuthData.ID)
			return nil, status.Error(codes.Internal, authErr.Error())
		}
		s.logger.Error("User creation failed", "email", req.Email, "err", err)
		s.storage.Delete(ctx, userAuthData.ID)
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("User created successfully", "id", req.Id, "email", req.Email)

	// Create permissions for user
	if s.rbacService != nil {
		if err := s.rbacService.AssignRoleToUser(ctx, userAuthData.ID, "unconfirmed_user"); err != nil {
			s.logger.Error("Failed to assign default RBAC role", "user_id", userAuthData.ID, "err", err)
			s.storage.Delete(ctx, userAuthData.ID)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	access, err := auth.GenerateAccessToken(&userAuthData.ID, &userAuthData.Email, &userAuthData.Version)
	if err != nil {
		s.storage.Delete(ctx, userAuthData.ID)
		return nil, status.Error(codes.Internal, err.Error())
	}

	refresh, err := auth.GenerateRefreshToken(&userAuthData.ID, &userAuthData.Email, &userAuthData.Version)
	if err != nil {
		s.storage.Delete(ctx, userAuthData.ID)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc.AddUserResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *Auth) ChangeEmail(ctx context.Context, req *grpc.ChangeEmailRequest) (*grpc.ChangeEmailResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if err := validator.New().Var(req.NewEmail, "required,email"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	}

	user, err := s.storage.ReadByID(ctx, id)
	if err != nil {
		if authErr, ok := err.(*errors.AuthError); ok && authErr.Code == "USER_NOT_FOUND" {
			return nil, status.Error(codes.NotFound, authErr.Error())
		}
		s.logger.Error("Failed to read user for email change", "id", id, "err", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.Email == req.NewEmail {
		return &grpc.ChangeEmailResponse{IsChanged: false}, nil
	}

	_, err = s.storage.Update(ctx, &models.UserAuthUpdate{
		ID:       id,
		NewEmail: req.NewEmail,
	})
	if err != nil {
		if authErr, ok := err.(*errors.AuthError); ok {
			switch authErr.Code {
			case "USER_NOT_FOUND":
				return nil, status.Error(codes.NotFound, authErr.Error())
			case "EMAIL_ALREADY_EXISTS":
				return nil, status.Error(codes.AlreadyExists, authErr.Error())
			}
		}
		s.logger.Error("Failed to change user email", "id", id, "err", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("User email changed successfully", "id", id)
	return &grpc.ChangeEmailResponse{IsChanged: true}, nil
}
