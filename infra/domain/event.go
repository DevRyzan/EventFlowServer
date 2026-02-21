package domain

import "strconv"
 
type Event struct {
	EventId     string            `json:"event_id,omitempty"`  
	UserId   string            `json:"user_id"`
	EventName   string            `json:"event_name"`
	Timestamp   int64             `json:"timestamp"` // Unix seconds
	Channel     string            `json:"channel,omitempty"`
	CampaignId  string            `json:"campaign_id,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]any    `json:"metadata,omitempty"`
}
 
func (e *Event) IdempotencyKey() string {
	if e.EventId != "" {
		return e.EventId
	}
	return e.UserId + "|" + e.EventName + "|" + strconv.FormatInt(e.Timestamp, 10)
}