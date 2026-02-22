package unit

import (
	"context"
	"errors"
	"testing"

	"eventflow/application/commands/create_event"
	"eventflow/infra/domain"
)

type mockEventStore struct {
	exists    bool
	existsErr error
	createErr error
}

func (m *mockEventStore) Exists(ctx context.Context, idempotencyKey string) (bool, error) {
	if m.existsErr != nil {
		return false, m.existsErr
	}
	return m.exists, nil
}

func (m *mockEventStore) Create(ctx context.Context, e *domain.Event) error {
	return m.createErr
}

func (m *mockEventStore) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	return nil, nil
}

func (m *mockEventStore) Update(ctx context.Context, e *domain.Event) error {
	return nil
}

func (m *mockEventStore) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockEventStore) InsertBatch(ctx context.Context, events []*domain.Event) error {
	return nil
}

func TestCreateEventService_Execute_Success(t *testing.T) {
	store := &mockEventStore{exists: false}
	svc := create_event.NewService(store)
	ctx := context.Background()

	req := &create_event.CreateEventRequest{
		UserId:    "1",
		EventName: "click",
		Timestamp: 1771718400,
	}

	resp, err := svc.Execute(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "accepted" {
		t.Errorf("expected status accepted, got %s", resp.Status)
	}
}

func TestCreateEventService_Execute_Duplicate(t *testing.T) {
	store := &mockEventStore{exists: true}
	svc := create_event.NewService(store)
	ctx := context.Background()

	req := &create_event.CreateEventRequest{
		UserId:    "1",
		EventName: "click",
		Timestamp: 1771718400,
	}

	_, err := svc.Execute(ctx, req)
	if err != create_event.ErrDuplicate {
		t.Errorf("expected ErrDuplicate, got %v", err)
	}
}

func TestCreateEventService_Execute_StoreError(t *testing.T) {
	storeErr := errors.New("db error")
	store := &mockEventStore{exists: false, createErr: storeErr}
	svc := create_event.NewService(store)
	ctx := context.Background()

	req := &create_event.CreateEventRequest{
		UserId:    "1",
		EventName: "click",
		Timestamp: 1771718300,
	}

	_, err := svc.Execute(ctx, req)
	if err != storeErr {
		t.Errorf("expected store error, got %v", err)
	}
}

func TestCreateEventRequest_ToDomain(t *testing.T) {
	req := &create_event.CreateEventRequest{
		EventId:   "evt-1",
		UserId:    "1",
		EventName: "click",
		Timestamp: 1771718300,
		Channel:   "web",
	}

	e := req.ToDomain()
	if e.EventId != "evt-1" || e.UserId != "1" || e.EventName != "click" {
		t.Errorf("ToDomain: unexpected values %+v", e)
	}
	if e.IdempotencyKey() != "evt-1" {
		t.Errorf("expected idempotency key evt-1, got %s", e.IdempotencyKey())
	}
}

func TestCreateEventRequest_ToDomain_IdempotencyKeyFallback(t *testing.T) {
	req := &create_event.CreateEventRequest{
		UserId:    "1",
		EventName: "click",
		Timestamp: 1771718300,
	}

	e := req.ToDomain()
	expected := "1|click|1771718300"
	if e.IdempotencyKey() != expected {
		t.Errorf("expected idempotency key %s, got %s", expected, e.IdempotencyKey())
	}
}
