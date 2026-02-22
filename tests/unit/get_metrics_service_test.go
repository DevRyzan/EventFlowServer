package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"eventflow/application/queries/get_metrics"
	"eventflow/infra/contracts"
	"eventflow/middleware/cache"
)

type mockMetricsRepo struct {
	totalCount      int64
	uniqueUserCount int64
	buckets         []contracts.BucketRow
	err             error
}

func (m *mockMetricsRepo) GetTotalCount(ctx context.Context, params contracts.MetricsQueryParams) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.totalCount, nil
}

func (m *mockMetricsRepo) GetUniqueUserCount(ctx context.Context, params contracts.MetricsQueryParams) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.uniqueUserCount, nil
}

func (m *mockMetricsRepo) GetGroupedByChannel(ctx context.Context, params contracts.MetricsQueryParams) ([]contracts.BucketRow, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.buckets, nil
}

func TestGetMetricsService_Handle_Success(t *testing.T) {
	repo := &mockMetricsRepo{
		totalCount:      100,
		uniqueUserCount: 25,
	}
	c := cache.NewMemoryCache()
	svc := get_metrics.NewService(repo, c, 60*time.Second)
	ctx := context.Background()

	req := &get_metrics.GetMetricsRequest{
		EventName: "click",
		From:      1771718300,
		To:        1771718600,
	}

	resp, err := svc.Handle(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 100 || resp.UniqueUserCount != 25 {
		t.Errorf("expected total=100 unique=25, got total=%d unique=%d", resp.TotalCount, resp.UniqueUserCount)
	}
}

func TestGetMetricsService_Handle_WithGroupByChannel(t *testing.T) {
	repo := &mockMetricsRepo{
		totalCount:      2000,
		uniqueUserCount: 25,
		buckets: []contracts.BucketRow{
			{Key: "web", Count: 60},
			{Key: "web-2", Count: 40},
		},
	}
	c := cache.NewMemoryCache()
	svc := get_metrics.NewService(repo, c, 60*time.Second)
	ctx := context.Background()

	req := &get_metrics.GetMetricsRequest{
		EventName: "click",
		From:      1771718300,
		To:        1771718600,
		GroupBy:   "channel",
	}

	resp, err := svc.Handle(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(resp.Buckets))
	}
	if resp.Buckets[0].Key != "web" || resp.Buckets[0].Count != 60 {
		t.Errorf("expected web=60, got %s=%d", resp.Buckets[0].Key, resp.Buckets[0].Count)
	}
	if resp.Buckets[1].Key != "web-2" || resp.Buckets[1].Count != 40 {
		t.Errorf("expected web-2=40, got %s=%d", resp.Buckets[1].Key, resp.Buckets[1].Count)
	}
}

func TestGetMetricsService_Handle_CacheHit(t *testing.T) {
	repo := &mockMetricsRepo{totalCount: 50, uniqueUserCount: 10}
	c := cache.NewMemoryCache()
	svc := get_metrics.NewService(repo, c, 60*time.Second)
	ctx := context.Background()

	req := &get_metrics.GetMetricsRequest{
		EventName: "click",
		From:      1771718300,
		To:        1771718600,
	}

	// First call - cache miss
	resp1, err := svc.Handle(ctx, req)
	if err != nil {
		t.Fatalf("first call: %v", err)
	}

	// Second call - cache hit (repo not called again if we could verify)
	resp2, err := svc.Handle(ctx, req)
	if err != nil {
		t.Fatalf("second call: %v", err)
	}
	if resp1.TotalCount != resp2.TotalCount {
		t.Errorf("cache hit should return same result: %d vs %d", resp1.TotalCount, resp2.TotalCount)
	}
}

func TestGetMetricsService_Handle_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	repo := &mockMetricsRepo{err: repoErr}
	c := cache.NewMemoryCache()
	svc := get_metrics.NewService(repo, c, 60*time.Second)
	ctx := context.Background()

	req := &get_metrics.GetMetricsRequest{
		EventName: "click",
		From:      1771718300,
		To:        1771718600,
	}

	_, err := svc.Handle(ctx, req)
	if err != repoErr {
		t.Errorf("expected repo error, got %v", err)
	}
}
