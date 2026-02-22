package get_metrics

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"eventflow/infra/contracts"
	"eventflow/middleware/cache"
)

// Service orchestrates the get metrics use case.
type Service struct {
	repo   contracts.MetricsRepository
	cache  cache.Cache
	cacheTTL time.Duration
}

// NewService returns a new GetMetricsService.
func NewService(repo contracts.MetricsRepository, c cache.Cache, cacheTTL time.Duration) *Service {
	return &Service{
		repo:     repo,
		cache:    c,
		cacheTTL: cacheTTL,
	}
}

// cacheKey builds a unique key for the metrics request.
func cacheKey(req *GetMetricsRequest) string {
	return fmt.Sprintf("metrics:%s|%d|%d|%s", req.EventName, req.From, req.To, req.GroupBy)
}

// Handle returns aggregated metrics for the given request.
func (s *Service) Handle(ctx context.Context, req *GetMetricsRequest) (*GetMetricsResponse, error) {
	key := cacheKey(req)
	if v, ok := s.cache.Get(key); ok {
		if resp, ok := v.(*GetMetricsResponse); ok {
			slog.Debug("get_metrics: cache hit", "key", key)
			return resp, nil
		}
	}

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

	s.cache.Set(key, resp, s.cacheTTL)
	return resp, nil
}
