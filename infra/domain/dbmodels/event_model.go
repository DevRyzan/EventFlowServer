package dbmodels

import (
	"time"

	"gorm.io/datatypes"
)

// EventModel is the GORM/persistence model for the events table.
type EventModel struct {
	IdempotencyKey string         `gorm:"column:idempotency_key;primaryKey"`
	EventName      string         `gorm:"column:event_name;not null"`
	UserID         string         `gorm:"column:user_id;not null"`
	EventTimestamp int64          `gorm:"column:event_timestamp;not null"`
	Channel        string         `gorm:"column:channel"`
	CampaignID     string         `gorm:"column:campaign_id"`
	Tags           datatypes.JSON `gorm:"column:tags;type:jsonb"`
	Metadata       datatypes.JSON `gorm:"column:metadata;type:jsonb"`
	CreatedAt      time.Time      `gorm:"column:created_at;->"` // read-only, DB default
}

// TableName overrides GORM's default table name.
func (EventModel) TableName() string {
	return "events"
}
