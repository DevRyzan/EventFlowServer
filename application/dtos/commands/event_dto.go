package dto

import "eventflow/infra/domain"

// CreateEventRequest is the DTO for POST /events request body.
type CreateEventRequest struct {
	EventId    string         `json:"event_id,omitempty"`
	UserId     string         `json:"user_id"`
	EventName  string         `json:"event_name"`
	Timestamp  int64          `json:"timestamp"`
	Channel    string         `json:"channel,omitempty"`
	CampaignId string         `json:"campaign_id,omitempty"`
	Tags       []string       `json:"tags,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

// ToDomain converts the DTO to a domain Event.
func (r *CreateEventRequest) ToDomain() *domain.Event {
	return &domain.Event{
		EventId:    r.EventId,
		UserId:     r.UserId,
		EventName:  r.EventName,
		Timestamp:  r.Timestamp,
		Channel:    r.Channel,
		CampaignId: r.CampaignId,
		Tags:       r.Tags,
		Metadata:   r.Metadata,
	}
}

// CreateEventResponse is the DTO for POST /events response.
type CreateEventResponse struct {
	Status string `json:"status"`
}
