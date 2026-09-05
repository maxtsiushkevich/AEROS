package grpc

import (
	grpc "auth/api/proto"
	auth "auth/pkg/auth"
	"auth/pkg/errors"
	"context"

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
			return nil, status.Error(codes.Internal, authErr.Error())
		}
		s.logger.Error("User creation failed", "email", req.Email, "err", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("User created successfully", "id", req.Id, "email", req.Email)

	// Create permissions for user
	s.rbacService.AddUserToRbac(userAuthData.ID)

	access, err := auth.GenerateAccessToken(&userAuthData.ID, &userAuthData.Email, &userAuthData.Version)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	refresh, err := auth.GenerateRefreshToken(&userAuthData.ID, &userAuthData.Email, &userAuthData.Version)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc.AddUserResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
