package sendstrategies

import (
	"context"
	"fmt"
	"users/internal/domain/outbox"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type GrpcStrategy[Request proto.Message] struct {
	newRequest func() Request
	call       func(context.Context, Request) error
}

func NewGrpcStrategy[Request proto.Message](
	newRequest func() Request,
	call func(context.Context, Request) error,
) *GrpcStrategy[Request] {
	return &GrpcStrategy[Request]{
		newRequest: newRequest,
		call:       call,
	}
}

func (s *GrpcStrategy[Request]) Exec(ctx context.Context, event outbox.OutboxEvent) error {
	request := s.newRequest()
	if err := protojson.Unmarshal(event.Payload, request); err != nil {
		return fmt.Errorf("decode gRPC request payload: %w", err)
	}

	return s.call(ctx, request)
}
