package get_metrics

import (
	"context"
	"log/slog"

	"eventflow/infra/contracts"
)

// Service orchestrates the get metrics use case.
type Service struct {
	repo contracts.MetricsRepository
}

// NewService returns a new GetMetricsService.
func NewService(repo contracts.MetricsRepository) *Service {
	return &Service{repo: repo}
}

// Handle returns aggregated metrics for the given request.
func (s *Service) Handle(ctx context.Context, req *GetMetricsRequest) (*GetMetricsResponse, error) {
	slog.Debug(
		QueryMsg,
		"event_name", req.EventName,
		"from", req.From,
		"to", req.To,
	)

	params := contracts.MetricsQueryParams{
		EventName: req.EventName,
		From:      req.From,
		To:        req.To,
		GroupBy:   req.GroupBy,
	}

	totalCount, err := s.repo.GetTotalCount(ctx, params)
	if err != nil {
		slog.Error(
			GetTotalCountFailedError,
			"event_name", req.EventName,
			"err", err,
			"params", params,
		)
		return nil, err
	}

	uniqueUserCount, err := s.repo.GetUniqueUserCount(ctx, params)
	if err != nil {
		slog.Error(
			GetUniqueUserCountFailedError,
			"event_name", req.EventName,
			"err", err,
			"params", params,
		)
		return nil, err
	}

	resp := &GetMetricsResponse{
		TotalCount:      totalCount,
		UniqueUserCount: uniqueUserCount,
		GroupBy:         req.GroupBy,
	}

	if req.GroupBy == "channel" {
		buckets, err := s.repo.GetGroupedByChannel(ctx, params)
		if err != nil {
			slog.Error(
				GetGroupedByChannelFailedError,
				"event_name", req.EventName,
				"err", err,
				"params", params,
			)
			return nil, err
		}
		for _, b := range buckets {
			resp.Buckets = append(resp.Buckets, Bucket{Key: b.Key, Count: b.Count})
		}
	}

	return resp, nil
}
