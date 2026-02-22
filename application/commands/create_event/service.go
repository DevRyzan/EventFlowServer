package create_event

import (
	"context"
	"log/slog"

	"eventflow/infra/contracts"
)

// Service orchestrates the create event use case.
type Service struct {
	store contracts.EventStore
}

// NewService returns a new CreateEventService.
func NewService(store contracts.EventStore) *Service {
	return &Service{store: store}
}

// Execute ingests an event with idempotency. Returns ErrDuplicate if already exists.
func (s *Service) Execute(
	ctx context.Context,
	req *CreateEventRequest,
) (*CreateEventResponse, error) {
	event := req.ToDomain()
	key := event.IdempotencyKey()

	exists, err := s.store.Exists(ctx, key)
	if err != nil {
		slog.Error(
			ExistsCheckFailedError,
			"key", key,
			"err", err,
			"event", event,
		)
		return nil, err
	}
	if exists {
		slog.Debug(
			DuplicateDetectedError,
			"key", key,
			"event", event,
		)
		return nil, ErrDuplicate
	}

	if err := s.store.Create(ctx, event); err != nil {
		slog.Error(
			CreateFailedError,
			"key", key,
			"err", err,
			"event", event)
		return nil, err
	}

	slog.Debug(AcceptedMsg, "key", key)
	return &CreateEventResponse{Status: AcceptedMsg}, nil
}
