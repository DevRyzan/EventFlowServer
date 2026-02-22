package events

import (
	"context"
	"log/slog"

	commands "eventflow/application/dtos/commands"
	"eventflow/infra/contracts"
)

// Service orchestrates event ingestion use cases.
type Service struct {
	store contracts.EventStore
}

// NewService returns a new EventService.
func NewService(store contracts.EventStore) *Service {
	return &Service{store: store}
}

// Create ingests an event with idempotency. Returns ErrDuplicate if already exists.
func (s *Service) Create(ctx context.Context, req *commands.CreateEventRequest) (*commands.CreateEventResponse, error) {
	event := req.ToDomain()
	key := event.IdempotencyKey()

	exists, err := s.store.Exists(ctx, key)
	if err != nil {
		slog.Error("event service: exists check failed", "key", key, "err", err)
		return nil, err
	}
	if exists {
		slog.Debug("event service: duplicate detected", "key", key)
		return nil, ErrDuplicate
	}

	if err := s.store.Create(ctx, event); err != nil {
		slog.Error("event service: create failed", "key", key, "err", err)
		return nil, err
	}

	slog.Debug("event service: event stored", "key", key)
	return &commands.CreateEventResponse{Status: "accepted"}, nil
}
