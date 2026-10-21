package service

import (
	"context"
)

type UserService interface {
	CheckEmailAvailability(ctx context.Context, email string) error
}
