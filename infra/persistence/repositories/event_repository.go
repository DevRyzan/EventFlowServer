package repositories

import (
	"context"
	"encoding/json"

	"eventflow/infra/contracts"
	"eventflow/infra/domain"

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
	err := s.db.WithContext(ctx).Model(&domain.EventModel{}).
		Where("idempotency_key = ?", idempotencyKey).
		Count(&count).Error
	return count > 0, err
}

func (s *PostgresEventStore) Insert(ctx context.Context, e *domain.Event) error {
	model := domainToModel(e)
	return s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&model).Error
}

func (s *PostgresEventStore) InsertBatch(ctx context.Context, events []*domain.Event) error {
	if len(events) == 0 {
		return nil
	}
	models := make([]domain.EventModel, len(events))
	for i, e := range events {
		models[i] = domainToModel(e)
	}
	return s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(models, 100).Error
}

func domainToModel(e *domain.Event) domain.EventModel {
	tagsJSON, _ := json.Marshal(e.Tags)
	metadataJSON, _ := json.Marshal(e.Metadata)
	return domain.EventModel{
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
