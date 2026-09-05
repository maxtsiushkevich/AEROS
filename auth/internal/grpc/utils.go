package grpc

import (
	auth "auth/api/proto"
	"auth/internal/models"
	hash "auth/pkg/auth"
	"auth/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AddUserRequestValidate struct {
	Id       string `validate:"required,uuid"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func getUserAuthData(req *auth.AddUserRequest) (*models.UserAuth, error) {
	validate := validator.New()

	val := AddUserRequestValidate{
		Id:       req.Id,
		Email:    req.Email,
		Password: req.Password,
	}

	err := validate.Struct(val)
	if err != nil {
		payload := errors.NewAuthError("INVALID_REQUEST", err.Error())
		return nil, status.Error(codes.InvalidArgument, payload.Error())
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		payload := errors.NewAuthError("HASH_FAILED", "failed to hash password")
		return nil, status.Error(codes.Internal, payload.Error())
	}

	id, _ := uuid.Parse(req.Id)

	return &models.UserAuth{
		ID:             id,
		Email:          val.Email,
		HashedPassword: hashedPassword,
	}, nil

}
