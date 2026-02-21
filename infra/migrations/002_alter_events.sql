 
CREATE TABLE IF NOT EXISTS events (
    idempotency_key TEXT PRIMARY KEY,
    name      TEXT        NOT NULL,
    user_id         TEXT        NOT NULL,
    event_timestamp BIGINT      NOT NULL,
    channel         TEXT,
    campaign_id     TEXT,
    tags            JSONB,
    metadata        JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
 -- 002_alter_events.sql (real alter)
ALTER TABLE events ADD COLUMN IF NOT EXISTS event_id TEXT;
CREATE INDEX IF NOT EXISTS idx_events_event_id ON events (event_id) WHERE event_id IS NOT NULL;