package outbox

import "context"

type Outbox interface {
	Run(ctx context.Context) error
	ProcessBatch(ctx context.Context) error
}
