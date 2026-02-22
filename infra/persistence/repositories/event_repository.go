package repositories

import (
	"context"
	"encoding/json"
	"errors"

	"eventflow/infra/contracts"
	"eventflow/infra/domain"
	"eventflow/infra/domain/dbmodels"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresEventStore struct {
	db *gorm.DB
}

func NewPostgresEventStore(db *gorm.DB) contracts.EventStore {
	return &PostgresEventStore{db: db}
}

func (s *PostgresEventStore) Exists(ctx context.Context, idempotencyKey string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&dbmodels.EventModel{}).
		Where("idempotency_key = ?", idempotencyKey).
		Count(&count).Error
	return count > 0, err
}

func (s *PostgresEventStore) Create(ctx context.Context, e *domain.Event) error {
	model := domainToModel(e)
	return s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&model).Error
}

func (s *PostgresEventStore) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	var m dbmodels.EventModel
	err := s.db.WithContext(ctx).Where("idempotency_key = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contracts.ErrNotFound
		}
		return nil, err
	}
	return modelToDomain(&m), nil
}

func (s *PostgresEventStore) Update(ctx context.Context, e *domain.Event) error {
	return contracts.ErrNotSupported
}

func (s *PostgresEventStore) Delete(ctx context.Context, id string) error {
	result := s.db.WithContext(ctx).Where("idempotency_key = ?", id).Delete(&dbmodels.EventModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return contracts.ErrNotFound
	}
	return nil
}

func (s *PostgresEventStore) InsertBatch(ctx context.Context, events []*domain.Event) error {
	if len(events) == 0 {
		return nil
	}
	models := make([]dbmodels.EventModel, len(events))
	for i, e := range events {
		models[i] = domainToModel(e)
	}
	return s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(models, 100).Error
}

func domainToModel(e *domain.Event) dbmodels.EventModel {
	tagsJSON, _ := json.Marshal(e.Tags)
	metadataJSON, _ := json.Marshal(e.Metadata)
	return dbmodels.EventModel{
		IdempotencyKey: e.IdempotencyKey(),
		EventName:      e.EventName,
		UserID:         e.UserId,
		EventTimestamp: e.Timestamp,
		Channel:        e.Channel,
		CampaignID:     e.CampaignId,
		Tags:           datatypes.JSON(tagsJSON),
		Metadata:       datatypes.JSON(metadataJSON),
	}
}

func modelToDomain(m *dbmodels.EventModel) *domain.Event {
	var tags []string
	_ = json.Unmarshal(m.Tags, &tags)
	var metadata map[string]any
	_ = json.Unmarshal(m.Metadata, &metadata)
	return &domain.Event{
		EventId:    m.IdempotencyKey,
		UserId:     m.UserID,
		EventName:  m.EventName,
		Timestamp:  m.EventTimestamp,
		Channel:    m.Channel,
		CampaignId: m.CampaignID,
		Tags:       tags,
		Metadata:   metadata,
	}
}
