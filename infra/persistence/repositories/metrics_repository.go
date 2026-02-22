package repositories

import (
	"context"

	"eventflow/infra/contracts"
	"eventflow/infra/domain/dbmodels"

	"gorm.io/gorm"
)

// PostgresMetricsRepository implements MetricsRepository using GORM.
type PostgresMetricsRepository struct {
	db *gorm.DB
}

// NewPostgresMetricsRepository returns a new PostgresMetricsRepository.
func NewPostgresMetricsRepository(db *gorm.DB) contracts.MetricsRepository {
	return &PostgresMetricsRepository{db: db}
}

// GetTotalCount returns the total event count for the given params.
func (r *PostgresMetricsRepository) GetTotalCount(ctx context.Context, params contracts.MetricsQueryParams) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&dbmodels.EventModel{}).
		Where("event_name = ? AND event_timestamp >= ? AND event_timestamp <= ?", params.EventName, params.From, params.To).
		Count(&count).Error
	return count, err
}

// GetUniqueUserCount returns the distinct user count for the given params.
func (r *PostgresMetricsRepository) GetUniqueUserCount(ctx context.Context, params contracts.MetricsQueryParams) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM events WHERE event_name = ? AND event_timestamp >= ? AND event_timestamp <= ?",
		params.EventName, params.From, params.To,
	).Scan(&count).Error
	return count, err
}

// GetGroupedByChannel returns channel buckets for the given params.
func (r *PostgresMetricsRepository) GetGroupedByChannel(ctx context.Context, params contracts.MetricsQueryParams) ([]contracts.BucketRow, error) {
	var rows []struct {
		Key   string
		Count int64
	}
	err := r.db.WithContext(ctx).Raw(
		"SELECT COALESCE(channel, '(empty)') as key, COUNT(*) as count FROM events WHERE event_name = ? AND event_timestamp >= ? AND event_timestamp <= ? GROUP BY COALESCE(channel, '(empty)')",
		params.EventName, params.From, params.To,
	).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	buckets := make([]contracts.BucketRow, len(rows))
	for i, r := range rows {
		buckets[i] = contracts.BucketRow{Key: r.Key, Count: r.Count}
	}
	return buckets, nil
}
