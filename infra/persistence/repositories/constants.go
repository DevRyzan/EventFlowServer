package repositories

const (
	// queris
	GetTotalCountQuery       = "event_name = ? AND event_timestamp >= ? AND event_timestamp <= ?"
	GetUniqueUserCountQuery  = "SELECT COUNT(DISTINCT user_id) FROM events WHERE event_name = ? AND event_timestamp >= ? AND event_timestamp <= ?"
	GetGroupedByChannelQuery = "SELECT COALESCE(channel, '(empty)') as key, COUNT(*) as count FROM events WHERE event_name = ? AND event_timestamp >= ? AND event_timestamp <= ? GROUP BY COALESCE(channel, '(empty)')"

	// err
	GetGroupedByChannelFailedError = "get grouped by channel failed"
	GetbyIdFailedError             = "get by id failed idempotency key"
	DeleteFailedError              = "delete failed idempotency key"
)
