package outbox

import (
	"context"
)

type OutboxSendStrategy interface {
	Exec(ctx context.Context, event OutboxEvent) error
}
