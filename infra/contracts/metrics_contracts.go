package contracts

import "context"

// MetricsQueryParams holds parameters for metrics aggregation.
type MetricsQueryParams struct {
	EventName string
	From      int64
	To        int64
	GroupBy   string
}

// MetricsResult holds aggregated metrics from the repository.
type MetricsResult struct {
	TotalCount      int64
	UniqueUserCount int64
	Buckets         []BucketRow
}

// BucketRow represents a single bucket from group-by aggregation.
type BucketRow struct {
	Key   string
	Count int64
}

// MetricsRepository abstracts metrics data access.
type MetricsRepository interface {
	GetTotalCount(ctx context.Context, params MetricsQueryParams) (int64, error)
	GetUniqueUserCount(ctx context.Context, params MetricsQueryParams) (int64, error)
	GetGroupedByChannel(ctx context.Context, params MetricsQueryParams) ([]BucketRow, error)
}
