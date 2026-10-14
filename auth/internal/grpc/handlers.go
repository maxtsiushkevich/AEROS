package grpc

import (
	grpc "auth/api/proto"
	"auth/pkg/errors"
	"context"
	"pkg/auth"

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
