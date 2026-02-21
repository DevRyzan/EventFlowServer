package contracts

import (
	"context"
	"eventflow/infra/domain"
)

type EventStore interface {
	Exists(ctx context.Context, idempotencyKey string) (bool, error)
	Insert(ctx context.Context, e *domain.Event) error
	InsertBatch(ctx context.Context, events []*domain.Event) error
}