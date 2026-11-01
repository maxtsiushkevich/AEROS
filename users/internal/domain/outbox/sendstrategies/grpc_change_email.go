package sendstrategies

import (
	"context"
	auth "users/api/proto"
	domainerr "users/internal/domain/errors"
	"users/internal/domain/outbox"

	"google.golang.org/grpc"
)

type GrpcChangeEmailStrategy struct {
	strategy *GrpcStrategy[*auth.ChangeEmailRequest]
}

func NewGrpcChangeEmailStrategy(conn *grpc.ClientConn) (*GrpcChangeEmailStrategy, error) {
	client := auth.NewAuthClient(conn)

	return &GrpcChangeEmailStrategy{
		strategy: NewGrpcStrategy(
			func() *auth.ChangeEmailRequest {
				return &auth.ChangeEmailRequest{}
			},
			func(ctx context.Context, request *auth.ChangeEmailRequest) error {
				response, err := client.ChangeEmail(ctx, request)
				if err != nil {
					return err
				}
				if !response.IsChanged {
					return domainerr.ErrUserAlreadyExists
				}

				return nil
			},
		),
	}, nil
}

func (s *GrpcChangeEmailStrategy) Exec(ctx context.Context, event outbox.OutboxEvent) error {
	return s.strategy.Exec(ctx, event)
}
