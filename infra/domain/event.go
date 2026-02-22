package domain

import "strconv"

// dto to domain event
type Event struct {
	EventId    string
	UserId     string
	EventName  string
	Timestamp  int64
	Channel    string
	CampaignId string
	Tags       []string
	Metadata   map[string]any
}

// for unique key for deduplication.
func (e *Event) IdempotencyKey() string {
	if e.EventId != "" {
		return e.EventId
	}
	return e.UserId + "|" + e.EventName + "|" + strconv.FormatInt(e.Timestamp, 10)
}
