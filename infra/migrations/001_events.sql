-- Events table models
CREATE TABLE IF NOT EXISTS events (
    idempotency_key TEXT PRIMARY KEY,
    event_name      TEXT        NOT NULL,
    user_id         TEXT        NOT NULL,
    event_timestamp BIGINT      NOT NULL,
    channel         TEXT,
    campaign_id     TEXT,
    tags            JSONB,
    metadata        JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_events_name_ts ON events (event_name, event_timestamp);

CREATE INDEX IF NOT EXISTS idx_events_user ON events (event_name, user_id);
