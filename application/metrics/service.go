package metrics

import (
	"eventflow/infra/domain"
	"eventflow/infra/domain/dbmodels"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Service implements metrics aggregation.
type Service struct {
	db *gorm.DB
}

// NewService returns a new metrics Service.
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Get returns aggregated metrics for the given request.
func (s *Service) Get(c echo.Context, req *domain.MetricsRequest) (*domain.MetricsResponse, error) {
	ctx := c.Request().Context()

	var totalCount int64
	err := s.db.WithContext(ctx).Model(&dbmodels.EventModel{}).
		Where("event_name = ? AND event_timestamp >= ? AND event_timestamp <= ?", req.EventName, req.From, req.To).
		Count(&totalCount).Error
	if err != nil {
		return nil, err
	}

	var uniqueUserCount int64
	err = s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM events WHERE event_name = ? AND event_timestamp >= ? AND event_timestamp <= ?",
		req.EventName, req.From, req.To,
	).Scan(&uniqueUserCount).Error
	if err != nil {
		return nil, err
	}

	resp := &domain.MetricsResponse{
		TotalCount:      totalCount,
		UniqueUserCount: uniqueUserCount,
		GroupBy:         req.GroupBy,
	}

	if req.GroupBy == "channel" {
		type channelRow struct {
			Key   string
			Count int64
		}
		var rows []channelRow
		err = s.db.WithContext(ctx).Raw(
			"SELECT COALESCE(channel, '(empty)') as key, COUNT(*) as count FROM events WHERE event_name = ? AND event_timestamp >= ? AND event_timestamp <= ? GROUP BY COALESCE(channel, '(empty)')",
			req.EventName, req.From, req.To,
		).Scan(&rows).Error
		if err != nil {
			return nil, err
		}
		for _, r := range rows {
			resp.Buckets = append(resp.Buckets, domain.Bucket{Key: r.Key, Count: r.Count})
		}
	}

	return resp, nil
}
