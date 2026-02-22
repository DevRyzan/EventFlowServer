package metrics

import (
	"context"
	"log/slog"

	"eventflow/infra/contracts"
	"eventflow/infra/domain"
)

// Service orchestrates metrics aggregation use cases.
type Service struct {
	repo contracts.MetricsRepository
}

// NewService returns a new metrics Service.
func NewService(repo contracts.MetricsRepository) *Service {
	return &Service{repo: repo}
}

// Get returns aggregated metrics for the given request.
func (s *Service) Get(ctx context.Context, req *domain.MetricsRequest) (*domain.MetricsResponse, error) {
	slog.Debug("metrics service: query", "event_name", req.EventName, "from", req.From, "to", req.To)
	params := contracts.MetricsQueryParams{
		EventName: req.EventName,
		From:      req.From,
		To:        req.To,
		GroupBy:   req.GroupBy,
	}

	totalCount, err := s.repo.GetTotalCount(ctx, params)
	if err != nil {
		slog.Error("metrics service: get total count failed", "event_name", req.EventName, "err", err)
		return nil, err
	}

	uniqueUserCount, err := s.repo.GetUniqueUserCount(ctx, params)
	if err != nil {
		slog.Error("metrics service: get unique user count failed", "event_name", req.EventName, "err", err)
		return nil, err
	}

	resp := &domain.MetricsResponse{
		TotalCount:      totalCount,
		UniqueUserCount: uniqueUserCount,
		GroupBy:         req.GroupBy,
	}

	if req.GroupBy == "channel" {
		buckets, err := s.repo.GetGroupedByChannel(ctx, params)
		if err != nil {
			slog.Error("metrics service: get grouped by channel failed", "event_name", req.EventName, "err", err)
			return nil, err
		}
		for _, b := range buckets {
			resp.Buckets = append(resp.Buckets, domain.Bucket{Key: b.Key, Count: b.Count})
		}
	}

	return resp, nil
}
