package contracts

import (
	"context"
	"eventflow/infra/domain"
)

// EventStore extends Repository with event-specific operations.
type EventStore interface {
	BaseContract[*domain.Event]
	Exists(ctx context.Context, idempotencyKey string) (bool, error)
	InsertBatch(ctx context.Context, events []*domain.Event) error
}